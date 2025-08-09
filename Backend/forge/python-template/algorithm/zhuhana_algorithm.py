import zhuhana
from zhuhana.types import OrderInstruction, OrderSide, OrderType, OrderMode, OrderTIF, OrderDomain, OHLCData



class ZhuhanaStrategy:
    def __init__(self, zhuhana_sdk: zhuhana.ZhuhanaClass):
        self.zhuhana_sdk: zhuhana.ZhuhanaClass = zhuhana_sdk


    def on_data(self, current_data: OHLCData):
        pass

    def condition_for_sell(self, current_data: OHLCData) -> OrderInstruction:
        return OrderInstruction(
            side=OrderSide.SELL,
            type=OrderType.MARKET,
            mode=OrderMode.INTRADAY,
            tif=OrderTIF.DAY,
            domain=OrderDomain.BACKTEST,
            quantity=100
        )

    def condition_for_buy(self, current_data: OHLCData) -> OrderInstruction:
        return OrderInstruction(
            side=OrderSide.BUY,
            type=OrderType.MARKET,
            mode=OrderMode.INTRADAY,
            tif=OrderTIF.DAY,
            domain=OrderDomain.BACKTEST,
            quantity=100
        )
