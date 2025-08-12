from __future__ import annotations
import requests
from typing import Iterator, Optional, NoReturn
from zhuhana.types import OHLCData, OHLCResponse
from zhuhana import ZhuhanaClass

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