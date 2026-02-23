/**
 * Smoke test: verify all API endpoints work correctly with minimal load.
 * 1 VU, 30 seconds - basic sanity check.
 */
import http from 'k6/http'
import { check, sleep } from 'k6'
import { Rate, Trend } from 'k6/metrics'

const errorRate = new Rate('errors')
const dashboardDuration = new Trend('dashboard_duration', true)
const orgsDuration = new Trend('organizations_duration', true)
const projectsDuration = new Trend('projects_duration', true)
const issuesDuration = new Trend('issues_duration', true)

export const options = {
  vus: 1,
  duration: '30s',
  thresholds: {
    // All requests must succeed
    errors: ['rate<0.01'],
    // p95 response times
    http_req_duration: ['p(95)<500'],
    dashboard_duration: ['p(95)<500'],
    organizations_duration: ['p(95)<300'],
    projects_duration: ['p(95)<500'],
    issues_duration: ['p(95)<500'],
  },
}

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'

export default function () {
  // 1. Dashboard summary
  const dashRes = http.get(`${BASE_URL}/api/v1/dashboard/summary`)
  dashboardDuration.add(dashRes.timings.duration)
  const dashOk = check(dashRes, {
    'dashboard: status 200': (r) => r.status === 200,
    'dashboard: has total_projects': (r) => JSON.parse(r.body).total_projects !== undefined,
  })
  errorRate.add(!dashOk)

  sleep(0.5)

  // 2. Organizations list
  const orgsRes = http.get(`${BASE_URL}/api/v1/organizations`)
  orgsDuration.add(orgsRes.timings.duration)
  const orgsOk = check(orgsRes, {
    'organizations: status 200': (r) => r.status === 200,
    'organizations: returns array': (r) => Array.isArray(JSON.parse(r.body)),
  })
  errorRate.add(!orgsOk)

  sleep(0.5)

  // 3. Projects list (default pagination)
  const projRes = http.get(`${BASE_URL}/api/v1/projects?page=1&per_page=25`)
  projectsDuration.add(projRes.timings.duration)
  const projOk = check(projRes, {
    'projects: status 200': (r) => r.status === 200,
    'projects: has data array': (r) => Array.isArray(JSON.parse(r.body).data),
    'projects: has pagination': (r) => JSON.parse(r.body).pagination !== undefined,
  })
  errorRate.add(!projOk)

  sleep(0.5)

  // 4. Issues list (default)
  const issRes = http.get(`${BASE_URL}/api/v1/issues?page=1&per_page=25`)
  issuesDuration.add(issRes.timings.duration)
  const issOk = check(issRes, {
    'issues: status 200': (r) => r.status === 200,
    'issues: has data array': (r) => Array.isArray(JSON.parse(r.body).data),
  })
  errorRate.add(!issOk)

  sleep(0.5)

  // 5. Issues filtered by RED
  const issRedRes = http.get(`${BASE_URL}/api/v1/issues?delay_status=RED&page=1&per_page=25`)
  check(issRedRes, {
    'issues(RED): status 200': (r) => r.status === 200,
  })

  sleep(1)
}
