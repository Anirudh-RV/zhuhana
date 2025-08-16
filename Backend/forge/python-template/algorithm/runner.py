import os

import zhuhana
from backtest.datastream import fake_ohlc_stream
from example_SMA_crossover import SMA2Strategy
from backtest.portfolio import ZhuhanaBacktestPortfolio
from backtest.engine import ZhuhanaBacktestEngine
from backtest.execution import SimpleExecutionModel


def main():
    # Swap this to your real API stream when ready:
    USER_ALGORITHM_TOKEN = str(os.getenv("USER_ALGORITHM_TOKEN"))
    API_ENDPOINT = str(os.getenv("API_ENDPOINT"))
    zhuhana_sdk: zhuhana.ZhuhanaClass = zhuhana.init(api_endpoint=API_ENDPOINT, token=USER_ALGORITHM_TOKEN)

    bars = fake_ohlc_stream(symbol="TEST", bars=300, start_price=100.0, noise=0.8)
    strategy = SMA2Strategy(short_window=5, long_window=20, trade_qty=20, zhuhana_sdk=zhuhana_sdk)
    portfolio = ZhuhanaBacktestPortfolio(init_cash=50_000.0)
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

if __name__ == "__main__":
    main()