import http from 'k6/http';
import { check } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Target is overridable: -e TARGET=http://algonexus:8080
const TARGET = __ENV.TARGET || 'http://algonexus:8080';
const URL = `${TARGET}/v1/ordermanager/submit`;

const okRate = new Rate('submit_ok');
const latency = new Trend('submit_latency', true);

const payload = JSON.stringify({
  symbol: 'AAPL',
  mode: 'INTRADAY',
  side: 'BUY',
  type: 'MARKET',
  domain: 'BACKTEST',
  time_in_force: 'DAY',
  quantity: 10,
  price: 150.25,
  priority: 1,
});

const params = { headers: { 'Content-Type': 'application/json' } };

// Ramp toward ~1000 RPS using a constant-arrival-rate executor.
export const options = {
  scenarios: {
    submit_load: {
      executor: 'ramping-arrival-rate',
      startRate: 100,
      timeUnit: '1s',
      preAllocatedVUs: 200,
      maxVUs: 1000,
      stages: [
        { target: 200, duration: '10s' },
        { target: 500, duration: '10s' },
        { target: 1000, duration: '20s' },
        { target: 1000, duration: '20s' }, // hold at ~1k RPS
      ],
    },
  },
  thresholds: {
    submit_ok: ['rate>0.99'],
    submit_latency: ['p(95)<50'],
  },
};

export default function () {
  const res = http.post(URL, payload, params);
  okRate.add(res.status === 200);
  latency.add(res.timings.duration);
  check(res, { 'status is 200': (r) => r.status === 200 });
}
