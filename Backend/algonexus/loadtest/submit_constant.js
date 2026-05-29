import http from 'k6/http';

// Constant-arrival-rate end-to-end probe.
// All knobs are env-driven so run_loadtest.py owns the configuration; defaults
// keep the script standalone-runnable.
//
//   k6 run -e TARGET=http://algonexus:8080 -e RATE=5000 -e DUR=30s submit_constant.js
//
// Drain proof (redis lag==0, pending==0, entries-read==accepted, 0 errors) is
// asserted by run_loadtest.py against XINFO; this script just generates load.

const env = __ENV;
const intEnv = (k, d) => parseInt(env[k] || String(d), 10);

const TARGET   = env.TARGET   || 'http://algonexus:8080';
const PATH     = env.PATH_    || '/v1/ordermanager/submit';
const URL      = `${TARGET}${PATH}`;
const RATE     = intEnv('RATE', 5000);
const DUR      = env.DUR      || '30s';
const PRE_VUS  = intEnv('PRE_VUS', Math.max(200, Math.ceil(RATE * 0.4)));
const MAX_VUS  = intEnv('MAX_VUS', Math.max(PRE_VUS, RATE * 2));

const payload  = env.PAYLOAD || JSON.stringify({
  symbol: 'AAPL', mode: 'INTRADAY', side: 'BUY', type: 'MARKET',
  domain: 'BACKTEST', time_in_force: 'DAY', quantity: 10, price: 150.25, priority: 1,
});
const params = { headers: { 'Content-Type': 'application/json' } };

export const options = {
  discardResponseBodies: true,
  scenarios: {
    constant: {
      executor: 'constant-arrival-rate',
      rate: RATE,
      timeUnit: '1s',
      duration: DUR,
      preAllocatedVUs: PRE_VUS,
      maxVUs: MAX_VUS,
    },
  },
};

export default function () {
  http.post(URL, payload, params);
}
