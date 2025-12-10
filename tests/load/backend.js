/**
 * k6 Load Test Configuration for Yousoon Backend
 * 
 * Run with: k6 run tests/load/backend.js
 * 
 * Environment variables:
 * - BASE_URL: API endpoint (default: http://localhost:8080)
 * - VUS: Number of virtual users (default: 50)
 * - DURATION: Test duration (default: 5m)
 */

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const latencyTrend = new Trend('latency');

// Test configuration
export const options = {
  stages: [
    { duration: '1m', target: 10 },   // Ramp up to 10 users
    { duration: '3m', target: 50 },   // Ramp up to 50 users
    { duration: '5m', target: 50 },   // Stay at 50 users
    { duration: '2m', target: 100 },  // Spike to 100 users
    { duration: '3m', target: 100 },  // Stay at 100 users
    { duration: '1m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% of requests under 500ms
    http_req_failed: ['rate<0.01'],     // Error rate under 1%
    errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// GraphQL queries
const queries = {
  getOffers: `
    query GetOffers($first: Int!, $location: GeoLocationInput) {
      offers(first: $first, location: $location) {
        edges {
          node {
            id
            title
            description
            discount {
              type
              value
            }
            establishment {
              name
              address {
                city
              }
            }
          }
        }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
    }
  `,
  
  getOffer: `
    query GetOffer($id: ID!) {
      offer(id: $id) {
        id
        title
        description
        discount {
          type
          value
          originalPrice
        }
        establishment {
          name
          address {
            street
            city
            postalCode
          }
          location {
            latitude
            longitude
          }
        }
        category {
          name
        }
        images {
          url
        }
        schedule {
          allDay
          slots {
            dayOfWeek
            startTime
            endTime
          }
        }
      }
    }
  `,
  
  getCategories: `
    query GetCategories {
      categories {
        id
        name
        slug
        icon
        color
      }
    }
  `,
  
  searchOffers: `
    query SearchOffers($query: String!, $first: Int!) {
      searchOffers(query: $query, first: $first) {
        edges {
          node {
            id
            title
            establishment {
              name
            }
          }
        }
      }
    }
  `,
};

// Helper to send GraphQL request
function graphqlRequest(query, variables = {}, authToken = null) {
  const headers = {
    'Content-Type': 'application/json',
  };
  
  if (authToken) {
    headers['Authorization'] = `Bearer ${authToken}`;
  }
  
  const payload = JSON.stringify({
    query,
    variables,
  });
  
  const response = http.post(`${BASE_URL}/graphql`, payload, { headers });
  
  const success = check(response, {
    'status is 200': (r) => r.status === 200,
    'no errors in response': (r) => {
      const body = JSON.parse(r.body);
      return !body.errors || body.errors.length === 0;
    },
  });
  
  errorRate.add(!success);
  latencyTrend.add(response.timings.duration);
  
  return response;
}

// Main test scenario
export default function() {
  // Random location in Paris area
  const location = {
    latitude: 48.8566 + (Math.random() - 0.5) * 0.1,
    longitude: 2.3522 + (Math.random() - 0.5) * 0.1,
  };
  
  group('Public API - Offers', () => {
    // Get offers list
    const offersResponse = graphqlRequest(queries.getOffers, {
      first: 20,
      location: location,
    });
    
    // If we got offers, get details of one
    if (offersResponse.status === 200) {
      const body = JSON.parse(offersResponse.body);
      if (body.data?.offers?.edges?.length > 0) {
        const offerId = body.data.offers.edges[0].node.id;
        graphqlRequest(queries.getOffer, { id: offerId });
      }
    }
    
    sleep(1);
  });
  
  group('Public API - Categories', () => {
    graphqlRequest(queries.getCategories);
    sleep(0.5);
  });
  
  group('Public API - Search', () => {
    const searchTerms = ['restaurant', 'bar', 'cinema', 'spa', 'brunch'];
    const randomTerm = searchTerms[Math.floor(Math.random() * searchTerms.length)];
    
    graphqlRequest(queries.searchOffers, {
      query: randomTerm,
      first: 10,
    });
    
    sleep(0.5);
  });
  
  // Simulate user think time
  sleep(Math.random() * 2 + 1);
}

// Teardown
export function teardown(data) {
  console.log('Load test completed');
}
