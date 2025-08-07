from dataclasses import dataclass
from typing import Optional
from enum import Enum

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


@dataclass
class OrderInstruction:
    """
    Represents an order to be placed by the strategy.
    """
    side: OrderSide              #: Direction of the order (BUY, SELL, SHORT, etc.)
    type: OrderType              #: Type of order (MARKET, LIMIT, etc.)
    mode: OrderMode              #: Whether INTRADAY or DELIVERY
    tif: OrderTIF                #: Time in force (DAY, IOC, GTC)
    domain: OrderDomain          #: The execution context (e.g., BACKTEST)
    quantity: float              #: Number of units
    price: Optional[float] = None  #: Price (optional for MARKET orders)


@dataclass
class OHLCData:
    Symbol: str
    Market: str
    Date_Time: str
    Open: float
    High: float
    Low: float
    Close: float
    Volume: int
    Day: int
    Weekday: int
    Week: int
    Month: int
    Year: int

@dataclass
class OHLCResponse:
    status: int
    status_description: str
    ohlc_data: OHLCData
    next_url: Optional[str]
