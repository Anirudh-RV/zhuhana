import algorithm
import zhuhana

from algorithm.zhuhana_algorithm import ZhuhanaStrategy
from zhuhana.types import OHLCData
from .portfolio import BacktestPortfolio

class BacktestEngine:
    def __init__(self, bars: list[OHLCData], strategy: ZhuhanaStrategy):
        self.bars = bars
        self.strategy = strategy
        self.portfolio = BacktestPortfolio()
    
    def run(self):
        for bar in self.bars:
            self.portfolio.update_price(bar.Close)
            self.strategy.on_data(bar)

