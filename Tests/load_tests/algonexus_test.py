from locust import HttpUser, task, between

class OrderUser(HttpUser):
    wait_time = between(1.1,1.2)
    @task
    def submit_order(self):
        self.client.post("/v1/algonexus/ordermanager/submit", json={
                                                 "symbol": "test",
                                                 "mode": "BACKTEST",
                                                 "side": "BUY",
                                                 "type": "MARKET",
                                                 "domain": "BACKTEST",
                                                 "time_in_force": "DAY",
                                                 "quantity": 100,
                                                 "price": 999,
                                                 "priority": 10
                                               })