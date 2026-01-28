import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

// Stress test: gradually increase load to find breaking point
export const options = {
  stages: [
    { duration: '2m', target: 50 },   // Ramp up to 50 users
    { duration: '3m', target: 50 },   // Stay at 50 for 3 minutes
    { duration: '2m', target: 100 },  // Ramp up to 100
    { duration: '3m', target: 100 },  // Stay at 100 for 3 minutes
    { duration: '2m', target: 200 },  // Ramp up to 200 (stress)
    { duration: '3m', target: 200 },  // Stay at 200 for 3 minutes
    { duration: '2m', target: 100 },  // Scale down to 100
    { duration: '2m', target: 50 },   // Scale down to 50
    { duration: '2m', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    'http_req_duration': ['p(95)<1000'],  // More lenient for stress test
    'http_req_failed': ['rate<0.05'],     // Allow up to 5% errors under stress
    'errors': ['rate<0.05'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

function getAuthToken() {
  const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({
    username: 'admin',
    password: 'admin123',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  if (loginRes.status === 200) {
    const body = JSON.parse(loginRes.body);
    return body.token;
  }
  return null;
}

export function setup() {
  const token = getAuthToken();
  console.log('Stress test starting...');
  return { token };
}

export default function(data) {
  const token = data.token;
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };

  // Mix of different endpoints to simulate real usage
  const endpoints = [
    { url: `${BASE_URL}/api/v1/dashboard/summary`, name: 'dashboard' },
    { url: `${BASE_URL}/api/v1/projects`, name: 'projects' },
    { url: `${BASE_URL}/api/v1/issues?limit=50`, name: 'issues' },
    { url: `${BASE_URL}/api/v1/organizations`, name: 'organizations' },
  ];

  // Randomly select an endpoint
  const endpoint = endpoints[Math.floor(Math.random() * endpoints.length)];

  const res = http.get(endpoint.url, {
    headers,
    tags: { name: endpoint.name },
  });

  const success = check(res, {
    'status is 2xx': (r) => r.status >= 200 && r.status < 300,
    'response time acceptable': (r) => r.timings.duration < 2000,
  });

  if (!success) {
    errorRate.add(1);
    console.log(`Error at VUs=${__VU}, endpoint=${endpoint.name}, status=${res.status}`);
  }

  // Random sleep between 0.5 and 2 seconds to simulate real user behavior
  sleep(Math.random() * 1.5 + 0.5);
}

export function teardown(data) {
  console.log('Stress test completed');
  console.log('Check the results to identify the breaking point');
}
