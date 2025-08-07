# zhuhana/zhuhana.py

import requests
from .types import OHLCData, OHLCResponse

class ZhuhanaClass:
    def __init__(self, api_endpoint: str, token: str):
        self.api_endpoint = api_endpoint
        self.token = token

    def get_ohlc_data_with_next(self, current_time: str, end_time: str, symbol: str, market: str, next_step: int) -> OHLCResponse:
        url = self.api_endpoint + "/v1/backtest/ohlc/next/"
        params = {
            "current_time": current_time,
            "end_time": end_time,
            "symbol": symbol,
            "market": market,
            "next_step": next_step
        }
        headers = {
            "USER_ALGORITHM_TOKEN": self.token
        }

        response = requests.get(url, headers=headers, params=params)
        data = response.json()

        if response.status_code != 200 or data.get("status") != 1:
            raise Exception(f"Failed to fetch data: {data}")

        ohlc_raw = data["ohlc_data"]
        ohlc = OHLCData(**ohlc_raw)

        return OHLCResponse(
            status=data["status"],
            status_description=data["status_description"],
            ohlc_data=ohlc,
            next_url=data.get("next_url")
        )
