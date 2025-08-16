class SimpleExecutionModel:
    """MVP Execution Model"""
    def __init__(self, fee_rate: float = 0.0, min_fee: float = 0.0, slippage: float = 0.0):
        self.fee_rate = fee_rate
        self.min_fee = min_fee
        self.slippage = slippage  # Price slippage as a fraction (e.g., 0.001 for 0.1% slippage)

    def fill_price(self, side: str, ref_price: float) -> float:
        # Buy: add slippage; Sell: subtract slippage
        if side == "BUY":
            return ref_price * (1.0 + self.slippage)
        else:
            return ref_price * (1.0 - self.slippage)

    def fee(self, gross: float) -> float:
        return max(abs(gross) * self.fee_rate, self.min_fee)