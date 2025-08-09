import time
import requests
from typing import Optional

from zhuhana.types import OHLCData, OHLCResponse

class ZhuhanaBacktestDataStream:
    """
    Sync version of OHLC data stream.
      - First: sdk.get_ohlc_data_with_next(...)
      - Then: next_url from response to fetch next data.
    """

    def __init__(
        self,
        zhuhana_sdk,
        *,
        start_time,
        end_time,
        symbol: str,
        market: str,
        frequency: str,
        token: str,
        session: Optional[requests.Session] = None,
        timeout: float = 10.0,
        max_retries: int = 3,
        backoff_base: float = 0.5,
    ):
        self.zhuhana_sdk = zhuhana_sdk
        self.start_time = start_time
        self.end_time = end_time
        self.symbol = symbol
        self.market = market
        self.frequency = frequency
        self.token = token

        self._session = session or requests.Session()
        self._timeout = timeout
        self._max_retries = max_retries
        self._backoff_base = backoff_base

        self._initial_done = False
        self._next_url: Optional[str] = None
        self._exhausted = False


    def __iter__(self):
        return self

    def __next__(self):
        bar = self.step()
        if bar is None:
            raise StopIteration
        return bar

    """
    Get next OHLCData
    Return: OHLCData or None if no more data
    """
    def step(self):
        if self._exhausted:
            return None
        # First fetch
        if not self._initial_done:
            resp = self._first_fetch()
            self._initial_done = True
            self._next_url = getattr(resp, "next_url", None)
            if resp and getattr(resp, "ohlc_data", None) is not None:
                if not self._next_url:
                    self._exhausted = True
                return resp.ohlc_data
            else:
                self._exhausted = True
                return None
        if not self._next_url:
            self._exhausted = True
            return None

        resp = self._fetch_next(self._next_url)
        self._next_url = getattr(resp, "next_url", None)
        if resp and getattr(resp, "ohlc_data", None) is not None:
            if not self._next_url:
                self._exhausted = True
            return resp.ohlc_data

        self._exhausted = True
        return None

    def skip(self, n: int) -> int:
        """
        skip n steps
        :param n: number of steps to skip
        :return: int, actual skipped steps (may be less than n if end reached)
        """
        skipped = 0
        while skipped < n:
            if self.step() is None:
                break
            skipped += 1
        return skipped

    def reset(self):
        """skip to the beginning of the stream."""
        self._initial_done = False
        self._next_url = None
        self._exhausted = False

    def close(self):
        """close the stream."""
        if self._session and not isinstance(self._session, requests.Session):
            return
        pass

    def _first_fetch(self):
        resp = self.zhuhana_sdk.get_ohlc_data_with_next(
            current_time=self.start_time,
            end_time=self.end_time,
            symbol=self.symbol,
            market=self.market,
            next_step=self.frequency,
        )
        return resp

    def _fetch_next(self, url: str):
        for attempt in range(1, self._max_retries + 1):
            try:
                r = self._session.get(
                    url,
                    headers={"USER_ALGORITHM_TOKEN": self.token},
                    timeout=self._timeout,
                )
                # network error
                r.raise_for_status()
                d = r.json()

                if d.get("status") != 1:
                    raise RuntimeError(f"API status not ok: {d}")

                # dict to OHLCData and OHLCResponse
                ohlc_data = OHLCData(**d["ohlc_data"])
                resp = OHLCResponse(
                    status=d["status"],
                    status_description=d["status_description"],
                    ohlc_data=ohlc_data,
                    next_url=d.get("next_url"),
                )
                return resp

            except (requests.RequestException, ValueError, RuntimeError) as e:
                if attempt >= self._max_retries:
                    raise
                time.sleep(self._backoff_base * (2 ** (attempt - 1)))