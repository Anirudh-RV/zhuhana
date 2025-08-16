from typing import Iterable
from zhuhana.strategy import ZhuhanaStrategy
from zhuhana.types import OHLCData
from .portfolio import ZhuhanaBacktestPortfolio
from .execution import SimpleExecutionModel


class ZhuhanaBacktestEngine:
    """
    MVP Backtest Engine
    - execution='NEXT_OPEN': process the previous bar's data, then execute orders at the open price of the next bar.
      For the first bar, only feed data without placing orders.
    - execution='CLOSE': process the current bar's data, then execute orders at the close price of the current bar.
    """
    def __init__(
        self,
        bars: Iterable[OHLCData],
        strategy: ZhuhanaStrategy,
        portfolio: ZhuhanaBacktestPortfolio | None = None,
        execution: str = "NEXT_OPEN",  # or "CLOSE"
        exec_model: SimpleExecutionModel | None = None,
    ):
        self.bars = iter(bars)
        self.strategy = strategy
        self.portfolio = portfolio or ZhuhanaBacktestPortfolio()
        self.execution = execution
        self.exec_model = exec_model or SimpleExecutionModel()

    def _update_marks(self, bar: OHLCData):
        # Update mark-to-market prices using close price
        self.portfolio.update_price(bar.Symbol, float(bar.Close))

    def _execute_orders(self, bar: OHLCData, orders) -> None:
        """
        Execute a batch of orders on the given bar.
        - Aggregate net quantity per symbol
        - Execute sells before buys to avoid inflating position size in the same bar
        """
        # 1) Filter out invalid orders
        valid = [o for o in orders if o and getattr(o, "quantity", 0)]
        if not valid:
            return
        # 2) Aggregate net quantities by symbol
        net_by_symbol = {}
        for o in valid:
            qty = int(o.quantity)
            if o.side == "SELL":
                qty = -qty
            net_by_symbol.setdefault(o.symbol, 0)
            net_by_symbol[o.symbol] += qty

        # 3) Execute net sells (negative qty) first, then net buys (positive qty)
        for side_group in ("SELL", "BUY"):
            for symbol, net_qty in net_by_symbol.items():
                if side_group == "SELL" and net_qty < 0:
                    qty = -net_qty
                    ref = bar.Open if self.execution == "NEXT_OPEN" else bar.Close
                    px = self.exec_model.fill_price("SELL", float(ref))
                    gross = px * qty
                    fee = self.exec_model.fee(gross)
                    ok = self.portfolio.sell(symbol, px, qty, fee=fee)
                    # Optional: if insufficient position, adjust to sell available quantity
                    # Here we simply ignore failures
                elif side_group == "BUY" and net_qty > 0:
                    qty = net_qty
                    ref = bar.Open if self.execution == "NEXT_OPEN" else bar.Close
                    px = self.exec_model.fill_price("BUY", float(ref))
                    gross = px * qty
                    fee = self.exec_model.fee(gross)
                    ok = self.portfolio.buy(symbol, px, qty, fee=fee)
                    # Optional: if insufficient cash, adjust to buy affordable quantity

    def run(self):
        """
        - NEXT_OPEN: Flow = (generate signals using prev_bar) -> execute orders at current bar's open -> update marks using current bar.
        - CLOSE: Flow = (generate signals using current bar) -> execute orders at current bar's close -> update marks using current bar.
        """
        first_bar = True
        prev_bar = None
        for bar in self.bars:
            # Always update marks first (so unrealized PnL is correct)
            self._update_marks(bar)

            if self.execution == "NEXT_OPEN":
                if first_bar:
                    # For the first bar: feed data only, no orders (orders can be placed on the next bar)
                    self.strategy.on_data(bar)
                    first_bar = False
                    prev_bar = bar
                    continue
                
                assert prev_bar is not None
                # 1) Generate orders from the previous bar
                self.strategy.on_data(prev_bar)
                orders = getattr(self.strategy, "generate_orders", None)
                if callable(orders):
                    orders = orders(prev_bar)
                else:
                    sell_ins = self.strategy.condition_for_sell(prev_bar)
                    buy_ins = self.strategy.condition_for_buy(prev_bar)
                    orders = []
                    if sell_ins: orders.append(sell_ins)
                    if buy_ins: orders.append(buy_ins)

                # 2) Execute orders at the open of the current bar
                self._execute_orders(bar, orders)

                prev_bar = bar

            elif self.execution == "CLOSE":
                # 1) Generate orders from the current bar
                self.strategy.on_data(bar)
                orders = getattr(self.strategy, "generate_orders", None)
                if callable(orders):
                    orders = orders(bar)
                else:
                    sell_ins = self.strategy.condition_for_sell(bar)
                    buy_ins = self.strategy.condition_for_buy(bar)
                    orders = []
                    if sell_ins: orders.append(sell_ins)
                    if buy_ins: orders.append(buy_ins)

                # 2) Execute orders at the close of the current bar
                self._execute_orders(bar, orders)

        # End of backtest: return a portfolio snapshot/metrics
        return self.portfolio.snapshot()