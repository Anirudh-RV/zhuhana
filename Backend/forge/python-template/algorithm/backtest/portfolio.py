class BacktestPortfolio:
    def __init__(self, initial_cash: float = 100_000):
        self.cash = initial_cash
        self.position = 0
        self.avg_price = 0
        self.last_price = 0 
        self.history = []

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