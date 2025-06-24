from dataclasses import dataclass
from typing import Optional

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
