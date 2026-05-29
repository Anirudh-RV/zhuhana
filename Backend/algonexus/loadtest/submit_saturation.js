import http from 'k6/http';
import { Trend, Rate, Counter } from 'k6/metrics';

// Saturation test: step the arrival rate up and HOLD at each plateau.
// The ceiling = the highest plateau where achieved RPS still tracks the
// target and http_req_failed stays ~0 with bounded latency.
const TARGET = __ENV.TARGET || 'http://algonexus:8080';
const URL = `${TARGET}/v1/ordermanager/submit`;

const lat = new Trend('submit_latency', true);
const okRate = new Rate('submit_ok');
const reqCount = new Counter('submit_count');

const payload = JSON.stringify({
  symbol: 'AAPL', mode: 'INTRADAY', side: 'BUY', type: 'MARKET',
  domain: 'BACKTEST', time_in_force: 'DAY', quantity: 10, price: 150.25, priority: 1,
});
const params = { headers: { 'Content-Type': 'application/json' } };

export const options = {
  discardResponseBodies: true,
  scenarios: {
    saturation: {
      executor: 'ramping-arrival-rate',
      startRate: 1000,
      timeUnit: '1s',
      preAllocatedVUs: 1000,
      maxVUs: 8000,
      stages: [
        { target: 2000,  duration: '5s'  }, { target: 2000,  duration: '10s' },
        { target: 5000,  duration: '5s'  }, { target: 5000,  duration: '10s' },
        { target: 10000, duration: '5s'  }, { target: 10000, duration: '10s' },
        { target: 15000, duration: '5s'  }, { target: 15000, duration: '10s' },
        { target: 20000, duration: '5s'  }, { target: 20000, duration: '10s' },
        { target: 30000, duration: '5s'  }, { target: 30000, duration: '10s' },
      ],
    },
  },
};

export default function () {
  const res = http.post(URL, payload, params);
  const ok = res.status === 200;
  okRate.add(ok);
  reqCount.add(1);
  lat.add(res.timings.duration);
}
