import http from 'k6/http';
import { sleep, check } from 'k6';

// k6 configuration for 1000+ concurrent virtual users (VUs)
export const options = {
  stages: [
    { duration: '30s', target: 500 },  // Ramp up to 500 users
    { duration: '1m', target: 1000 },  // Ramp up to 1000 users and hold
    { duration: '30s', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],   // Error rate must be less than 1%
    http_req_duration: ['p(95)<500'], // 95% of requests must complete under 500ms
  },
};

export default function () {
  // Target our backend health endpoint
  const url = 'http://localhost:5000/health';
  
  const res = http.get(url);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.5); // wait 500ms between requests for each VU
}
