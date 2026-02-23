/**
 * Stress test: find the breaking point by gradually increasing load.
 * Identifies bottlenecks and maximum capacity.
 */
import http from 'k6/http'
import { check, sleep } from 'k6'
import { Rate, Trend } from 'k6/metrics'

const errorRate = new Rate('errors')
const responseTimes = new Trend('response_times', true)

export const options = {
  stages: [
    { duration: '20s', target: 5 },   // Warm up
    { duration: '30s', target: 10 },  // Normal load
    { duration: '30s', target: 20 },  // High load
    { duration: '30s', target: 40 },  // Stress
    { duration: '20s', target: 50 },  // Near-breaking
    { duration: '30s', target: 0 },   // Recovery
  ],
  thresholds: {
    // Stress test: allow higher error rate as we're finding limits
    http_req_failed: ['rate<0.10'],
    response_times: ['p(95)<2000'],
  },
}

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'

// Focus on the heaviest endpoints
export default function () {
  const endpoints = [
    `${BASE_URL}/api/v1/issues?page=1&per_page=25`,
    `${BASE_URL}/api/v1/projects?page=1&per_page=25`,
    `${BASE_URL}/api/v1/dashboard/summary`,
  ]

  const url = endpoints[Math.floor(Math.random() * endpoints.length)]
  const res = http.get(url)
  responseTimes.add(res.timings.duration)

  const ok = check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 2s': (r) => r.timings.duration < 2000,
  })
  errorRate.add(!ok)

  sleep(0.5)
}
