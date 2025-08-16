from __future__ import annotations
import requests
import math
import random

from typing import Iterator, Optional, NoReturn
import zhuhana
from zhuhana import ZhuhanaClass
from zhuhana.types import OHLCData, OHLCResponse



class ZhuhanaBacktestDataStream:
    """
    Iterator OHLC bars for backtesting.
    - First bar sdk.get_ohlc_data_with_next(...)
    - Next bars: next_url + requests.get(...)
    - EOF/Error: StopIteration
    """
    def __init__(
        self,
        zhuhana_sdk: ZhuhanaClass,          # zhuhana.init(...)
        *,
        market: str,
        symbol: str,
        start_time: str,
        end_time: str,
        frequency: int,
    ) -> None:
        self.sdk = zhuhana_sdk
        self.market = market
        self.symbol = symbol
        self.start_time = start_time
        self.end_time = end_time
        self.frequency = frequency

        self._initial_fetch_done: bool = False
        self._next_url: Optional[str] = None
        self._exhausted: bool = False

    def __iter__(self) -> Iterator[OHLCData]:
        return self

    def _stop(self) -> NoReturn:
        self._exhausted = True
        raise StopIteration

    def __next__(self) -> OHLCData:
        if self._exhausted:
            self._stop()

        if not self._initial_fetch_done:
            self._initial_fetch_done = True
            try:
                resp: OHLCResponse = self.sdk.get_ohlc_data_with_next(
                    current_time=self.start_time,
                    end_time=self.end_time,
                    symbol=self.symbol,
                    market=self.market,
                    next_step=self.frequency,
                )
            except Exception:
                self._stop()

            if not resp or not resp.ohlc_data:
                self._stop()

            self._next_url = resp.next_url
            return resp.ohlc_data

        if not self._next_url:
            self._stop()

        try:
            headers = {"USER_ALGORITHM_TOKEN": self.sdk.token}
            r = requests.get(self._next_url, headers=headers, timeout=10)
            data = r.json()
            if r.status_code != 200 or data.get("status") != 1:
                self._stop()

            ohlc = OHLCData(**data["ohlc_data"])
            self._next_url = data.get("next_url")
            if not ohlc:
                self._stop()

            return ohlc
        except Exception:
            self._stop()


def fake_ohlc_stream(
    symbol: str = "FAKE",
    bars: int = 100,
    start_price: float = 100.0,
    noise: float = 0.5
) -> Iterator[OHLCData]:
    """
    Generate a fake OHLC data stream for testing.
    :param symbol: Ticker symbol
    :param bars: Number of bars to generate
    :param start_price: Initial price
    :param noise: Random noise level
    """
    price = start_price
    for i in range(bars):
        # base trend using sine wave
        base_change = math.sin(i / 10) * 0.5
        # random noise
        rnd = random.uniform(-noise, noise)
        open_price = price
        close_price = max(1.0, price + base_change + rnd)  # avoid negative prices
        high_price = max(open_price, close_price) + random.uniform(0, 0.3)
        low_price = min(open_price, close_price) - random.uniform(0, 0.3)

        data = OHLCData(
            Symbol=symbol,
            Open=round(open_price, 2),
            High=round(high_price, 2),
            Low=round(low_price, 2),
            Close=round(close_price, 2),
            Market="FAKE",
            Date_Time=str(i),  # use index as datetime
            Volume=1000,  # constant volume for simplicity
            Day= i % 30 + 1,
            Weekday=i % 7,
            Week=i // 7 + 1,
            Month=(i // 30) + 1,
            Year=2025
        )

        print(f"Generated OHLC: {data}")

        yield data

        price = close_price