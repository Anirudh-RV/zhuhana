import zhuhana
from zhuhana.types import (
    OHLCData,
    OrderDomain,
    OrderInstruction,
    OrderMode,
    OrderSide,
    OrderTIF,
    OrderType,
)


class ZhuhanaStrategy:
    def __init__(self, zhuhana_sdk: zhuhana.ZhuhanaClass):
      """
      Init function for the Strategy to initialize any variables you want to.
      """
      self.zhuhana_sdk: zhuhana.ZhuhanaClass = zhuhana_sdk

    def on_data(self, current_data: OHLCData):
      """
      Use this function to describe what you want to do when you get a new data point.
      For example,
      Describe all the variables and setup you want to do for the logic iteration for this data point.
      """
      pass

    def condition_for_sell(self, current_data: OHLCData) -> OrderInstruction:
      """
      Use this function to describe the logic required for a Sell condition
      """
      return OrderInstruction(
            side=OrderSide.SELL,
            type=OrderType.MARKET,
            mode=OrderMode.INTRADAY,
            tif=OrderTIF.DAY,
            domain=OrderDomain.BACKTEST,
            quantity=100,
        )

    def condition_for_buy(self, current_data: OHLCData) -> OrderInstruction:
      """
      Use this function to describe the logic required for a Buy condition
      """
      return OrderInstruction(
            side=OrderSide.BUY,
            type=OrderType.MARKET,
            mode=OrderMode.INTRADAY,
            tif=OrderTIF.DAY,
            domain=OrderDomain.BACKTEST,
            quantity=100,
        )
