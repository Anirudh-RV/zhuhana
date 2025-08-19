import os

import zhuhana
from backtest.datastream import fake_ohlc_stream
from example_SMA_crossover import SMA2Strategy
from backtest.portfolio import ZhuhanaBacktestPortfolio
from backtest.engine import ZhuhanaBacktestEngine
from backtest.execution import SimpleExecutionModel

class ZhuhanaStrategyRunner:
    def __init__(self, USER_ALGORITHM_TOKEN, ORDER_DOMAIN, API_ENDPOINT, strategy):
        self.zhuhana_sdk: zhuhana.ZhuhanaClass = zhuhana.init(api_endpoint=API_ENDPOINT, token=USER_ALGORITHM_TOKEN)
        self.order_domain = ORDER_DOMAIN
        self.user_algorithm_token = USER_ALGORITHM_TOKEN
        self.api_endpoint = API_ENDPOINT
        self.strategy = strategy

        if ORDER_DOMAIN == "Backtest":
            self.market: str = os.getenv("MARKET", "")
            self.symbol: str = os.getenv("SYMBOL", "")
            self.start_time: str = os.getenv("START_TIME", "")
            self.end_time: str = os.getenv("END_TIME", "")
            self.portfolio_size: int = int(os.getenv("PORTFOLIO_SIZE", 0))
            self.frequency: int = int(os.getenv("FREQUENCY", 0))
            self.next_url: str = ""  # initial state
            self.initial_fetch_done: bool = False
        else:
            raise ValueError("Unsupported ORDER_DOMAIN")



    def run(self, bars, portfolio, strategy, exec_model, execution="NEXT_OPEN"):
        # Swap this to your real API stream when ready:
        USER_ALGORITHM_TOKEN = str(os.getenv("USER_ALGORITHM_TOKEN"))
        API_ENDPOINT = str(os.getenv("API_ENDPOINT"))
        zhuhana_sdk: zhuhana.ZhuhanaClass = zhuhana.init(api_endpoint=API_ENDPOINT, token=USER_ALGORITHM_TOKEN)

        bars = fake_ohlc_stream(symbol="TEST", bars=300, start_price=100.0, noise=0.8)
        portfolio = ZhuhanaBacktestPortfolio(init_cash=50_000.0)

        strategy = SMA2Strategy(short_window=5, long_window=20, trade_qty=20, zhuhana_sdk=zhuhana_sdk)

        engine = ZhuhanaBacktestEngine(
            bars=bars,
            strategy=strategy,
            portfolio=portfolio,
            execution="NEXT_OPEN",  # or "CLOSE"
            exec_model=SimpleExecutionModel(fee_rate=0.0005, min_fee=0.0, slippage=0.0005),
        )

        print("=== Starting Backtest ===")
        result = engine.run()

        print("=== Final Snapshot ===")
        for k, v in result.items():
            print(f"{k}: {v}")
        print("\n=== Trades (first 10) ===")
        for rec in portfolio.history[:10]:
            print(rec)
        print(f"... total trades: {len(portfolio.history)}")