/**
 * k6 Load Test - Authenticated User Scenarios
 * 
 * Tests booking flow and user-specific endpoints
 * 
 * Run with: k6 run tests/load/authenticated.js -e BASE_URL=http://localhost:8080
 */

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const bookingSuccessRate = new Rate('booking_success');
const latencyTrend = new Trend('latency');

export const options = {
  scenarios: {
    // Normal load
    average_load: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 30 },
        { duration: '5m', target: 30 },
        { duration: '1m', target: 0 },
      ],
      gracefulRampDown: '30s',
    },
    // Peak hour simulation
    peak_hour: {
      executor: 'ramping-arrival-rate',
      startRate: 10,
      timeUnit: '1s',
      preAllocatedVUs: 100,
      maxVUs: 200,
      stages: [
        { duration: '1m', target: 50 },
        { duration: '3m', target: 100 },
        { duration: '1m', target: 10 },
      ],
      startTime: '8m', // Start after average_load
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.01'],
    booking_success: ['rate>0.95'],
    errors: ['rate<0.05'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// GraphQL mutations and queries
const mutations = {
  login: `
    mutation Login($email: String!, $password: String!) {
      login(email: $email, password: $password) {
        accessToken
        refreshToken
        user {
          id
          profile {
            firstName
          }
        }
      }
    }
  `,
  
  bookOuting: `
    mutation BookOuting($offerId: ID!) {
      createOuting(offerId: $offerId) {
        id
        status
        qrCode {
          code
        }
        offer {
          title
        }
      }
    }
  `,
  
  cancelOuting: `
    mutation CancelOuting($id: ID!, $reason: String) {
      cancelOuting(id: $id, reason: $reason) {
        id
        status
      }
    }
  `,
  
  addFavorite: `
    mutation AddFavorite($offerId: ID!) {
      addFavorite(offerId: $offerId) {
        id
        offerId
      }
    }
  `,
  
  removeFavorite: `
    mutation RemoveFavorite($offerId: ID!) {
      removeFavorite(offerId: $offerId)
    }
  `,
};

const queries = {
  myOutings: `
    query MyOutings($first: Int!, $status: OutingStatus) {
      myOutings(first: $first, status: $status) {
        edges {
          node {
            id
            status
            bookedAt
            offer {
              title
            }
          }
        }
      }
    }
  `,
  
  myFavorites: `
    query MyFavorites {
      myFavorites {
        id
        offer {
          id
          title
        }
      }
    }
  `,
  
  myProfile: `
    query MyProfile {
      me {
        id
        email
        profile {
          firstName
          lastName
          avatar
        }
        subscription {
          status
          plan {
            code
          }
        }
      }
    }
  `,
  
  nearbyOffers: `
    query NearbyOffers($location: GeoLocationInput!, $first: Int!) {
      offers(location: $location, first: $first) {
        edges {
          node {
            id
            title
            discount {
              type
              value
            }
          }
        }
      }
    }
  `,
};

// Helper functions
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
    'no GraphQL errors': (r) => {
      try {
        const body = JSON.parse(r.body);
        return !body.errors || body.errors.length === 0;
      } catch {
        return false;
      }
    },
  });
  
  errorRate.add(!success);
  latencyTrend.add(response.timings.duration);
  
  return response;
}

// Generate test user credentials
function getTestUser(vuId) {
  return {
    email: `loadtest+user${vuId}@yousoon.com`,
    password: 'LoadTest123!',
  };
}

// Setup: Create test users (run once before tests)
export function setup() {
  console.log('Setting up load test...');
  // In real scenario, you'd create test users here or use pre-created ones
  return {
    testOfferIds: ['test-offer-1', 'test-offer-2', 'test-offer-3'],
  };
}

// Main test scenario
export default function(data) {
  const user = getTestUser(__VU);
  let authToken = null;
  
  group('Authentication', () => {
    const loginResponse = graphqlRequest(mutations.login, {
      email: user.email,
      password: user.password,
    });
    
    if (loginResponse.status === 200) {
      try {
        const body = JSON.parse(loginResponse.body);
        if (body.data?.login?.accessToken) {
          authToken = body.data.login.accessToken;
        }
      } catch (e) {
        console.log('Login failed:', e);
      }
    }
    
    sleep(0.5);
  });
  
  if (!authToken) {
    // Skip authenticated tests if login failed
    console.log('Skipping authenticated tests - login failed');
    return;
  }
  
  group('User Profile', () => {
    graphqlRequest(queries.myProfile, {}, authToken);
    sleep(0.3);
  });
  
  group('Browse Offers', () => {
    const location = {
      latitude: 48.8566 + (Math.random() - 0.5) * 0.1,
      longitude: 2.3522 + (Math.random() - 0.5) * 0.1,
    };
    
    const offersResponse = graphqlRequest(queries.nearbyOffers, {
      location,
      first: 20,
    }, authToken);
    
    sleep(1);
    
    // Randomly like an offer (20% chance)
    if (Math.random() < 0.2 && offersResponse.status === 200) {
      try {
        const body = JSON.parse(offersResponse.body);
        if (body.data?.offers?.edges?.length > 0) {
          const randomOffer = body.data.offers.edges[Math.floor(Math.random() * body.data.offers.edges.length)];
          graphqlRequest(mutations.addFavorite, { offerId: randomOffer.node.id }, authToken);
        }
      } catch (e) {}
    }
    
    sleep(0.5);
  });
  
  group('Booking Flow', () => {
    // Get available offers
    const location = {
      latitude: 48.8566,
      longitude: 2.3522,
    };
    
    const offersResponse = graphqlRequest(queries.nearbyOffers, {
      location,
      first: 10,
    }, authToken);
    
    if (offersResponse.status === 200) {
      try {
        const body = JSON.parse(offersResponse.body);
        if (body.data?.offers?.edges?.length > 0) {
          // Book a random offer (10% chance to simulate realistic booking rate)
          if (Math.random() < 0.1) {
            const randomOffer = body.data.offers.edges[Math.floor(Math.random() * body.data.offers.edges.length)];
            
            const bookingResponse = graphqlRequest(mutations.bookOuting, {
              offerId: randomOffer.node.id,
            }, authToken);
            
            const bookingSuccess = bookingResponse.status === 200 && 
              JSON.parse(bookingResponse.body).data?.createOuting?.id;
            
            bookingSuccessRate.add(bookingSuccess);
            
            // Cancel booking (cleanup, 50% of the time)
            if (bookingSuccess && Math.random() < 0.5) {
              const bookingId = JSON.parse(bookingResponse.body).data.createOuting.id;
              sleep(0.5);
              graphqlRequest(mutations.cancelOuting, {
                id: bookingId,
                reason: 'Load test cleanup',
              }, authToken);
            }
          }
        }
      } catch (e) {
        console.log('Booking flow error:', e);
      }
    }
    
    sleep(1);
  });
  
  group('View My Outings', () => {
    graphqlRequest(queries.myOutings, { first: 10 }, authToken);
    sleep(0.5);
  });
  
  group('View My Favorites', () => {
    graphqlRequest(queries.myFavorites, {}, authToken);
    sleep(0.5);
  });
  
  // Simulate user think time
  sleep(Math.random() * 3 + 1);
}

export function teardown(data) {
  console.log('Load test completed');
  // Cleanup: Remove test data if needed
}
