import algorithm
import zhuhana

from typing import Iterable
from algorithm.zhuhana_algorithm import ZhuhanaStrategy
from zhuhana.types import OHLCData
from .portfolio import ZhuhanaBacktestPortfolio
from .datastream import ZhuhanaBacktestDataStream

class ZhuhanaBacktestEngine:

    # Initialize the backtest engine with OHLC data and strategy
    # bars_stream: OHLCDataflow object (generateor of OHLCData)
    # strategy: instance of ZhuhanaStrategy
    # dataflow: "BACKTEST" or other modes
    def __init__(self, bars: Iterable["OHLCData"], strategy: ZhuhanaStrategy):
        self.bars = bars
        self.strategy = strategy
        self.portfolio = ZhuhanaBacktestPortfolio()
    
    def run(self):
        for bar in self.bars:
            self.portfolio.update_price(bar.Close)
            self.strategy.on_data(bar)
            sell_instruction = self.strategy.condition_for_sell(bar)
            buy_instruction = self.strategy.condition_for_buy(bar)
            
            if sell_instruction:
                self.portfolio.sell(bar.Close, sell_instruction.quantity)

            if buy_instruction:
                self.portfolio.buy(bar.Close, buy_instruction.quantity)
        
    

