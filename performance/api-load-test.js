import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
  stages: [
    { duration: '30s', target: 10 },  // Ramp up to 10 users
    { duration: '1m', target: 50 },   // Ramp up to 50 users
    { duration: '2m', target: 50 },   // Stay at 50 users
    { duration: '30s', target: 100 }, // Spike to 100 users
    { duration: '1m', target: 100 },  // Stay at 100 users
    { duration: '30s', target: 0 },   // Ramp down to 0
  ],
  thresholds: {
    'http_req_duration': ['p(95)<500'],  // 95% of requests must complete below 500ms
    'http_req_duration{name:dashboard}': ['p(95)<500'],
    'http_req_duration{name:projects}': ['p(95)<300'],
    'http_req_duration{name:issues}': ['p(95)<500'],
    'http_req_failed': ['rate<0.01'],    // Error rate must be less than 1%
    'errors': ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Login and get token
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
  // Get auth token for authenticated requests
  const token = getAuthToken();
  return { token };
}

export default function(data) {
  const token = data.token;
  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };

  // Test 1: Health check
  {
    const res = http.get(`${BASE_URL}/health`);
    check(res, {
      'health check status is 200': (r) => r.status === 200,
    }) || errorRate.add(1);
  }

  sleep(0.5);

  // Test 2: Dashboard API
  {
    const res = http.get(`${BASE_URL}/api/v1/dashboard/summary`, {
      headers,
      tags: { name: 'dashboard' },
    });

    check(res, {
      'dashboard status is 200': (r) => r.status === 200,
      'dashboard response time < 500ms': (r) => r.timings.duration < 500,
      'dashboard has data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body && typeof body === 'object';
        } catch {
          return false;
        }
      },
    }) || errorRate.add(1);
  }

  sleep(1);

  // Test 3: Organizations list
  {
    const res = http.get(`${BASE_URL}/api/v1/organizations`, {
      headers,
      tags: { name: 'organizations' },
    });

    check(res, {
      'organizations status is 200': (r) => r.status === 200,
      'organizations response time < 300ms': (r) => r.timings.duration < 300,
    }) || errorRate.add(1);
  }

  sleep(0.5);

  // Test 4: Projects list
  {
    const res = http.get(`${BASE_URL}/api/v1/projects`, {
      headers,
      tags: { name: 'projects' },
    });

    check(res, {
      'projects status is 200': (r) => r.status === 200,
      'projects response time < 300ms': (r) => r.timings.duration < 300,
      'projects has array': (r) => {
        try {
          const body = JSON.parse(r.body);
          return Array.isArray(body);
        } catch {
          return false;
        }
      },
    }) || errorRate.add(1);
  }

  sleep(1);

  // Test 5: Issues search with filters
  {
    const res = http.get(`${BASE_URL}/api/v1/issues?delay_status=RED&limit=20`, {
      headers,
      tags: { name: 'issues' },
    });

    check(res, {
      'issues status is 200': (r) => r.status === 200,
      'issues response time < 500ms': (r) => r.timings.duration < 500,
      'issues has array': (r) => {
        try {
          const body = JSON.parse(r.body);
          return Array.isArray(body);
        } catch {
          return false;
        }
      },
    }) || errorRate.add(1);
  }

  sleep(1);

  // Test 6: Get specific project
  {
    const res = http.get(`${BASE_URL}/api/v1/projects/1`, {
      headers,
      tags: { name: 'project_detail' },
    });

    check(res, {
      'project detail status is 200 or 404': (r) => r.status === 200 || r.status === 404,
      'project detail response time < 200ms': (r) => r.timings.duration < 200,
    }) || errorRate.add(1);
  }

  sleep(2);
}

export function teardown(data) {
  // Cleanup if needed
  console.log('Test completed');
}
