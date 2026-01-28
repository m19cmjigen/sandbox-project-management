import http from 'k6/http';
import { check, sleep } from 'k6';

// Smoke test: minimal load to verify basic functionality
export const options = {
  vus: 1,
  duration: '30s',
  thresholds: {
    'http_req_duration': ['p(95)<1000'],
    'http_req_failed': ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
  // Login to get auth token
  const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({
    username: 'admin',
    password: 'admin123',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  if (loginRes.status !== 200) {
    throw new Error(`Login failed: ${loginRes.status}`);
  }

  const body = JSON.parse(loginRes.body);
  return { token: body.token };
}

export default function(data) {
  const headers = {
    'Authorization': `Bearer ${data.token}`,
    'Content-Type': 'application/json',
  };

  // Test critical paths only
  const tests = [
    {
      name: 'Health Check',
      request: () => http.get(`${BASE_URL}/health`),
      checks: {
        'status is 200': (r) => r.status === 200,
      },
    },
    {
      name: 'Dashboard',
      request: () => http.get(`${BASE_URL}/api/v1/dashboard/summary`, { headers }),
      checks: {
        'status is 200': (r) => r.status === 200,
        'has data': (r) => r.body.length > 0,
      },
    },
    {
      name: 'Projects List',
      request: () => http.get(`${BASE_URL}/api/v1/projects`, { headers }),
      checks: {
        'status is 200': (r) => r.status === 200,
      },
    },
    {
      name: 'Issues List',
      request: () => http.get(`${BASE_URL}/api/v1/issues?limit=10`, { headers }),
      checks: {
        'status is 200': (r) => r.status === 200,
      },
    },
  ];

  tests.forEach(test => {
    const res = test.request();
    const passed = check(res, test.checks);

    if (!passed) {
      console.error(`❌ ${test.name} failed`);
    } else {
      console.log(`✅ ${test.name} passed`);
    }

    sleep(1);
  });
}

export function teardown(data) {
  console.log('Smoke test completed');
}
