from models import OrderInstruction

class ZhuhanaStrategy:
    def __init__(self, zhuhana_sdk):
        self.zhuhana_sdk = zhuhana_sdk

    def on_data(self, current_data):
       pass

    def condition_for_sell(self, current_data) -> OrderInstruction:
        pass

    def condition_for_buy(self, current_data) -> OrderInstruction:
        pass
