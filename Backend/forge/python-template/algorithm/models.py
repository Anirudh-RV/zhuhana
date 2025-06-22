from enum import Enum
from pydantic import BaseModel, Field
from typing import Optional


class OrderSide(str, Enum):
    BUY = "BUY"
    SELL = "SELL"
    SHORT = "SHORT"
    INVALID = "INVALID"


class OrderType(str, Enum):
    MARKET = "MARKET"
    LIMIT = "LIMIT"
    STOP = "STOP"
    STOP_LIMIT = "STOP_LIMIT"
    FILL_OR_KILL = "FILL_OR_KILL"
    IMMEDIATE_OR_CANCEL = "IMMEDIATE_OR_CANCEL"
    ALL_OR_NONE = "ALL_OR_NONE"
    INVALID = "INVALID"


class OrderMode(str, Enum):
    INTRADAY = "INTRADAY"
    DELIVERY = "DELIVERY"


class OrderDomain(str, Enum):
    BACKTEST = "BACKTEST"


class OrderTIF(str, Enum):
    DAY = "DAY"
    GTC = "GTC"
    IOC = "IOC"


class OrderInstruction(BaseModel):
    symbol: str = Field(..., description="Stock symbol, e.g., 'AAPL', 'GOOGL'")
    mode: OrderMode
    side: OrderSide
    type: OrderType
    domain: Optional[OrderDomain] = None
    time_in_force: Optional[OrderTIF] = None
    quantity: float
    price: float = Field(..., description="Price per share for limit orders")
    priority: Optional[int] = None

    def __init__(
        self,
        symbol: str,
        mode: OrderMode,
        side: OrderSide,
        type: OrderType,
        quantity: float,
        price: float,
        domain: Optional[OrderDomain] = None,
        time_in_force: Optional[OrderTIF] = None,
        priority: Optional[int] = None,
    ):
        super().__init__(
            symbol=symbol,
            mode=mode,
            side=side,
            type=type,
            quantity=quantity,
            price=price,
            domain=domain,
            time_in_force=time_in_force,
            priority=priority,
        )
