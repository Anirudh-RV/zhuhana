from collections import deque
from zhuhana.types import OrderInstruction, OrderSide, OrderType, OrderMode, OrderTIF, OrderDomain
from zhuhana import ZhuhanaStrategy

class SMA2Strategy(ZhuhanaStrategy):
    """
    SMA Crossover Strategy:
    - BUY when short MA crosses above long MA (golden cross)
    - SELL when short MA crosses below long MA (death cross)
    Debounce with last_signal to avoid repeated orders.
    """
    def __init__(self, zhuhana_sdk, short_window: int = 5, long_window: int = 20, trade_qty: int = 100):
        self.zhuhana_sdk = zhuhana_sdk
        self.short_window = short_window
        self.long_window = long_window
        self.trade_qty = trade_qty

        self._closes = deque(maxlen=long_window)
        self.short_ma: float | None = None
        self.long_ma: float | None = None

        self._prev_diff: float | None = None  # previous (short_ma - long_ma)
        self._curr_diff: float | None = None  # current  (short_ma - long_ma)

        self.last_signal: str | None = None   # "BUY" / "SELL" / None

    def on_data(self, current_data):
        """Update indicators with the newest bar."""
        self._closes.append(float(current_data.Close))

        # compute short/long MAs if enough data
        if len(self._closes) >= self.short_window:
            tail = list(self._closes)[-self.short_window:]
            self.short_ma = sum(tail) / self.short_window

        if len(self._closes) >= self.long_window:
            self.long_ma = sum(self._closes) / self.long_window

        # roll diff
        self._prev_diff = self._curr_diff
        if self.short_ma is not None and self.long_ma is not None:
            self._curr_diff = self.short_ma - self.long_ma
        else:
            self._curr_diff = None

    # --- helpers -------------------------------------------------------------
    def _make_order(self, side: str, qty: int, current_data):
        """Build OrderInstruction; pass a symbol if supported, else attach dynamically."""
        return OrderInstruction(
            symbol=current_data.Symbol,
            side=OrderSide.BUY if side == "BUY" else OrderSide.SELL,
            type=OrderType.MARKET,
            mode=OrderMode.INTRADAY,
            tif=OrderTIF.DAY,
            domain=OrderDomain.BACKTEST,
            quantity=qty,
        )


    # --- signal methods used by engine ---
    def condition_for_buy(self, current_data):
        """Return BUY order on golden cross; otherwise None."""
        if self._prev_diff is None or self._curr_diff is None:
            return None

        crossed_up = self._prev_diff <= 0 and self._curr_diff > 0
        if crossed_up and self.last_signal != "BUY":
            self.last_signal = "BUY"
            return self._make_order("BUY", self.trade_qty, current_data)
        return None

    def condition_for_sell(self, current_data):
        """Return SELL order on death cross; otherwise None."""
        if self._prev_diff is None or self._curr_diff is None:
            return None
        
        crossed_down = self._prev_diff >= 0 and self._curr_diff < 0
        if crossed_down and self.last_signal != "SELL":
            self.last_signal = "SELL"
            return self._make_order("SELL", self.trade_qty, current_data)
        return None