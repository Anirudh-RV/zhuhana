# algorithm/apps.py

from django.apps import AppConfig
import threading
import os

class AlgorithmConfig(AppConfig):
    default_auto_field = 'django.db.models.BigAutoField'
    name = 'algorithm'

    def ready(self):
        # Optional: only skip if manage.py runserver to avoid double-start
        if os.environ.get("RUN_MAIN") == "true":
            print("Inside Django autoreload process - skipping runner to avoid double start")
            return

        from algorithm.algorithm_runner import ZhuhanaStrategyRunner

        def run_strategy():
            USER_ALGORITHM_TOKEN = os.getenv("USER_ALGORITHM_TOKEN")
            ORDER_DOMAIN = os.getenv("ORDER_DOMAIN")
            API_ENDPOINT = os.getenv("API_ENDPOINT")
            print(f"Starting strategy")
            match ORDER_DOMAIN:
                case "Backtest":
                    runner = ZhuhanaStrategyRunner(
                        USER_ALGORITHM_TOKEN=USER_ALGORITHM_TOKEN,
                        ORDER_DOMAIN=ORDER_DOMAIN,
                        API_ENDPOINT=API_ENDPOINT
                    )
                    while runner.OnData():
                        pass
                case _:
                    print("Unsupported ORDER_DOMAIN")

            print("Container run completed")

        threading.Thread(target=run_strategy, daemon=True).start()
