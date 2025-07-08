import os
import requests
import zhuhana
from zhuhana.types import OHLCData, OHLCResponse

from algorithm.zhuhana_algorithm import ZhuhanaStrategy

class ZhuhanaStrategyRunner:
    def __init__(self, USER_ALGORITHM_TOKEN, ORDER_DOMAIN, API_ENDPOINT):
        self.zhuhana_sdk: zhuhana.ZhuhanaClass = zhuhana.init(api_endpoint=API_ENDPOINT, token=USER_ALGORITHM_TOKEN)
        self.order_domain = ORDER_DOMAIN
        self.user_algorithm_token = USER_ALGORITHM_TOKEN
        self.api_endpoint = API_ENDPOINT

        if ORDER_DOMAIN == "Backtest":
            self.market: str = os.getenv("MARKET")
            self.symbol: str = os.getenv("SYMBOL")
            self.start_time: str = os.getenv("START_TIME")
            self.end_time: str = os.getenv("END_TIME")
            self.portfolio_size: int = int(os.getenv("PORTFOLIO_SIZE"))
            self.frequency: int = int(os.getenv("FREQUENCY"))
            self.next_url: str = None  # initial state
            self.initial_fetch_done: bool = False
        else:
            raise ValueError("Unsupported ORDER_DOMAIN")

    def OnData(self) -> bool:
        if self.order_domain == "Backtest":
            try:
                if not self.initial_fetch_done:
                    # First fetch
                    response = self.zhuhana_sdk.get_ohlc_data_with_next(
                        current_time=self.start_time,
                        end_time=self.end_time,
                        symbol=self.symbol,
                        market=self.market,
                        next_step=self.frequency
                    )
                    self.initial_fetch_done = True
                elif self.next_url:
                    headers = {"USER_ALGORITHM_TOKEN": self.zhuhana_sdk.token}
                    response_raw = requests.get(self.next_url, headers=headers)
                    response_dict = response_raw.json()
                    if response_raw.status_code != 200 or response_dict.get("status") != 1:
                        raise Exception(f"Failed to fetch next_url data: {response_dict}")

                    # Convert dict to OHLCData and OHLCResponse
                    ohlc_data = OHLCData(**response_dict["ohlc_data"])
                    response: OHLCResponse = OHLCResponse(
                        status=response_dict["status"],
                        status_description=response_dict["status_description"],
                        ohlc_data=ohlc_data,
                        next_url=response_dict.get("next_url")
    )
                else:
                    return False  # No more data

                ohlc_data: OHLCData = response.ohlc_data
                self.next_url = response.next_url

                if not ohlc_data:
                    return False

                # Run strategy
                zhuhanaStrategy = ZhuhanaStrategy(zhuhana_sdk=self.zhuhana_sdk)
                zhuhanaStrategy.on_data(ohlc_data)

                sellInstruction = zhuhanaStrategy.condition_for_sell(ohlc_data)
                buyInstruction = zhuhanaStrategy.condition_for_buy(ohlc_data)

                if sellInstruction:
                    print(f"[SELL] {sellInstruction}")
                elif buyInstruction:
                    print(f"[BUY] {buyInstruction}")
                else:
                    print("[HOLD]")

                return True

            except Exception as e:
                print(f"Error during backtest OnData: {e}")
                return False
        else:
            raise ValueError("Unsupported ORDER_DOMAIN")
