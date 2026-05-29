import http from 'k6/http';
import { check } from 'k6';

const TARGET = __ENV.TARGET || 'http://algonexus:8080';
const URL = `${TARGET}/v1/ordermanager/submit`;

const payload = JSON.stringify({
  symbol: 'SPY', mode: 'INTRADAY', side: 'BUY', type: 'MARKET',
  domain: 'BACKTEST', time_in_force: 'DAY', quantity: 10, price: 150.25, priority: 1,
});
const params = { headers: { 'Content-Type': 'application/json' } };

// Moderate fixed rate: enough orders to make a handle leak obvious, gentle on 1 vCPU.
export const options = {
  discardResponseBodies: true,
  scenarios: {
    loop: {
      executor: 'constant-arrival-rate',
      rate: 150, timeUnit: '1s', duration: '20s',
      preAllocatedVUs: 100, maxVUs: 400,
    },
  },
};

export default function () {
  const res = http.post(URL, payload, params);
  check(res, { 'status is 200': (r) => r.status === 200 });
}
