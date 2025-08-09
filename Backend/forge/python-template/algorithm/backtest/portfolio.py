class ZhuhanaBacktestPortfolio:
    def __init__(self, init_cash: float = 100.000):
        self.cash = init_cash
        self.positions: dict[str, int] = {}         # symbol -> volume
        self.avg_prices: dict[str, float] = {}      # symbol -> avg cost
        self.last_prices: dict[str, float] = {}     # symbol -> latest close
        self.history: list[tuple] = []              # (symbol, side, price, volume)

    def update_price(self, price: float):
        self.last_price = price

    def buy(self, price: float, volume: int):
        cost = price * volume
        if self.cash >= cost:
            if self.position == 0:
                self.avg_price = price
            else:
                self.avg_price = (self.avg_price * self.position + cost) / (self.position + volume)

            self.position += volume
            self.cash -= cost
            self.history.append(("BUY", price, volume))
        else:
            print("Not enough cash to buy.")

    def sell(self, price: float, volume: int):
        if self.position >= volume:
            self.position -= volume
            self.cash += price * volume
            self.history.append(("SELL", price, volume))
        else:
            print("Not enough position to sell.")

    def total_value(self):
        return self.cash + self.position * self.last_price