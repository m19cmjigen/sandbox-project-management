/**
 * Load test: simulate normal production traffic.
 * Ramp up to 10 VUs, hold for 1 minute, ramp down.
 */
import http from 'k6/http'
import { check, sleep } from 'k6'
import { Rate, Counter } from 'k6/metrics'

const errorRate = new Rate('errors')
const requestCount = new Counter('total_requests')

export const options = {
  stages: [
    { duration: '15s', target: 5 },   // Ramp up
    { duration: '60s', target: 10 },  // Hold at 10 VUs
    { duration: '15s', target: 0 },   // Ramp down
  ],
  thresholds: {
    errors: ['rate<0.01'],
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.01'],
  },
}

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'

// Weighted scenario: simulate realistic user navigation
const scenarios = [
  { weight: 30, fn: loadDashboard },
  { weight: 25, fn: loadProjects },
  { weight: 30, fn: loadIssues },
  { weight: 15, fn: loadOrganizations },
]

function weightedRandom() {
  const total = scenarios.reduce((s, sc) => s + sc.weight, 0)
  let r = Math.random() * total
  for (const sc of scenarios) {
    r -= sc.weight
    if (r <= 0) return sc.fn
  }
  return scenarios[0].fn
}

function loadDashboard() {
  const res = http.get(`${BASE_URL}/api/v1/dashboard/summary`)
  requestCount.add(1)
  check(res, { 'dashboard ok': (r) => r.status === 200 })
  errorRate.add(res.status !== 200)
}

function loadProjects() {
  const page = Math.floor(Math.random() * 3) + 1
  const res = http.get(`${BASE_URL}/api/v1/projects?page=${page}&per_page=25`)
  requestCount.add(1)
  check(res, { 'projects ok': (r) => r.status === 200 })
  errorRate.add(res.status !== 200)
}

function loadIssues() {
  const filters = ['', 'delay_status=RED', 'delay_status=YELLOW', 'delay_status=GREEN']
  const filter = filters[Math.floor(Math.random() * filters.length)]
  const url = filter
    ? `${BASE_URL}/api/v1/issues?${filter}&page=1&per_page=25`
    : `${BASE_URL}/api/v1/issues?page=1&per_page=25`
  const res = http.get(url)
  requestCount.add(1)
  check(res, { 'issues ok': (r) => r.status === 200 })
  errorRate.add(res.status !== 200)
}

function loadOrganizations() {
  const res = http.get(`${BASE_URL}/api/v1/organizations`)
  requestCount.add(1)
  check(res, { 'organizations ok': (r) => r.status === 200 })
  errorRate.add(res.status !== 200)
}

export default function () {
  const fn = weightedRandom()
  fn()
  sleep(Math.random() * 2 + 0.5) // 0.5–2.5s think time
}
