from zhuhana_algorithm import ZhuhanaStrategy

class ZhuhanaStrategyRunner():
    def __init__(self):
        self.zhuhana_sdk = {} # Initialize Zhuhana SDK

    def OnData(self):
        zhuhanaStrategy = ZhuhanaStrategy(zhuhana_sdk=self.zhuhana_sdk)

        current_data = self.get_currrent_data()

        zhuhanaStrategy.on_data(current_data)

        sellCondition, sellInstruction = zhuhanaStrategy.condition_for_sell(current_data)
        buyCondition, buyInstruction = zhuhanaStrategy.condition_for_buy(current_data)

        # if sellCondition:
        #     OrderSystem.Sell(sellInstruction)
        # elif buyCondition:
        #     OrderSystem.Buy(buyInstruction)
        # else:
        #     pass

    def get_current_data(self):
        current_data = self.zhuhana_sdk.get_next_data(
            "USER_DEFINED_DATA_PROVIDER",
            "USER_DEFINED_STOCK_MARKET",
            "USER_DEFINED_STOCK_TICKER")
        return current_data
