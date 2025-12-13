/**
 * k6 Load Test - Simple API Check
 * 
 * A simpler load test for initial validation.
 * Run with: k6 run tests/load/api-health.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
  stages: [
    { duration: '30s', target: 10 },  // Ramp up to 10 users
    { duration: '1m', target: 20 },   // Ramp up to 20 users
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'],  // 95% of requests under 1s
    http_req_failed: ['rate<0.05'],     // Error rate under 5%
    errors: ['rate<0.05'],
  },
  // Skip TLS verification for staging environments with self-signed certs
  insecureSkipTLSVerify: true,
};

const BASE_URL = __ENV.BASE_URL || 'http://api.yousoon.com';

// GraphQL health check query
const healthQuery = `
  query Health {
    __typename
  }
`;

// Simple offers query
const offersQuery = `
  query GetOffers {
    offers(first: 10) {
      edges {
        node {
          id
          title
        }
      }
    }
  }
`;

export default function () {
  // Test 1: GraphQL Health Check
  const healthRes = http.post(`${BASE_URL}/graphql`, JSON.stringify({
    query: healthQuery,
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  const healthCheck = check(healthRes, {
    'health: status is 200': (r) => r.status === 200,
    'health: response time < 500ms': (r) => r.timings.duration < 500,
  });
  
  errorRate.add(!healthCheck);
  
  sleep(1);

  // Test 2: Get Offers
  const offersRes = http.post(`${BASE_URL}/graphql`, JSON.stringify({
    query: offersQuery,
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  const offersCheck = check(offersRes, {
    'offers: status is 200': (r) => r.status === 200,
    'offers: response time < 1000ms': (r) => r.timings.duration < 1000,
    'offers: has data': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data !== undefined;
      } catch {
        return false;
      }
    },
  });
  
  errorRate.add(!offersCheck);
  
  sleep(1);
}

export function handleSummary(data) {
  return {
    'stdout': JSON.stringify({
      checks_passed: data.metrics.checks.values.passes,
      checks_failed: data.metrics.checks.values.fails,
      http_req_duration_p95: data.metrics.http_req_duration.values['p(95)'],
      error_rate: data.metrics.errors.values.rate,
    }, null, 2),
  };
}
