from dataclasses import dataclass, field
from typing import Dict, List, Tuple, Optional

@dataclass
class Trade:
    symbol: str
    side: str             # "BUY" or "SELL"
    price: float
    volume: int
    fee: float
    cash_after: float
    position_after: int

class ZhuhanaBacktestPortfolio:
    def __init__(self, init_cash: float = 100.0):
        self.cash: float = float(init_cash)
        self.positions: Dict[str, int] = {}          # symbol -> shares
        self.avg_prices: Dict[str, float] = {}       # symbol -> average cost
        self.last_prices: Dict[str, float] = {}      # symbol -> last close/mark
        self.history: List[Trade] = []
        self.realized_pnl: float = 0.0

    def update_price(self, symbol: str, price: float) -> None:
        self.last_prices[symbol] = float(price)

    def _ensure_symbol(self, symbol: str) -> None:
        if symbol not in self.positions:
            self.positions[symbol] = 0
            self.avg_prices[symbol] = 0.0

    def buy(self, symbol: str, price: float, volume: int, fee: float = 0.0) -> bool:
        assert volume > 0, "buy volume must be > 0"
        self._ensure_symbol(symbol)
        cost = price * volume + fee
        if self.cash < cost:
            return False  # or raise ValueError("Not enough cash")
        pos = self.positions[symbol]
        avg = self.avg_prices[symbol]
        new_pos = pos + volume
        new_avg = price if pos == 0 else (avg * pos + price * volume) / new_pos
        # commit
        self.positions[symbol] = new_pos
        self.avg_prices[symbol] = new_avg
        self.cash -= cost
        self.history.append(Trade(symbol, "BUY", price, volume, fee, self.cash, new_pos))
        return True

    def sell(self, symbol: str, price: float, volume: int, fee: float = 0.0) -> bool:
        assert volume > 0, "sell volume must be > 0"
        self._ensure_symbol(symbol)
        pos = self.positions[symbol]
        if pos < volume:
            return False  # or raise ValueError("Not enough position")
        avg = self.avg_prices[symbol]
        proceeds = price * volume - fee
        self.cash += proceeds
        new_pos = pos - volume
        # realized PnL on closed shares
        self.realized_pnl += (price - avg) * volume - fee
        # if flat, reset avg cost
        if new_pos == 0:
            self.avg_prices[symbol] = 0.0
        self.positions[symbol] = new_pos
        self.history.append(Trade(symbol, "SELL", price, volume, fee, self.cash, new_pos))
        return True

    def position_value(self, symbol: str) -> float:
        shares = self.positions.get(symbol, 0)
        mark = self.last_prices.get(symbol, 0.0)
        return shares * mark

    def unrealized_pnl(self, symbol: Optional[str] = None) -> float:
        if symbol is not None:
            pos = self.positions.get(symbol, 0)
            if pos == 0:
                return 0.0
            return (self.last_prices.get(symbol, 0.0) - self.avg_prices.get(symbol, 0.0)) * pos
        # portfolio-wide
        return sum(self.unrealized_pnl(sym) for sym in self.positions.keys())

    def total_value(self) -> float:
        # cash + current market value of all positions
        return self.cash + sum(self.position_value(sym) for sym in self.positions.keys())

    def snapshot(self) -> Dict[str, float]:
        return {
            "cash": self.cash,
            "equity": self.total_value(),
            "realized_pnl": self.realized_pnl,
            "unrealized_pnl": self.unrealized_pnl(),
        }