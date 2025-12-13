import { test, expect } from '@playwright/test';

/**
 * E2E Tests for Yousoon Backend API
 * 
 * These tests validate the GraphQL API is functional after deployment.
 */

const API_URL = process.env.API_URL || 'http://api.yousoon.com';

test.describe('API Health Checks', () => {
  test('GraphQL endpoint responds', async ({ request }) => {
    const response = await request.post(`${API_URL}/graphql`, {
      data: {
        query: '{ __typename }',
      },
    });

    expect(response.ok()).toBeTruthy();
    expect(response.status()).toBe(200);
  });

  test('GraphQL introspection works', async ({ request }) => {
    const response = await request.post(`${API_URL}/graphql`, {
      data: {
        query: `
          query IntrospectionQuery {
            __schema {
              queryType {
                name
              }
            }
          }
        `,
      },
    });

    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.data.__schema.queryType.name).toBe('Query');
  });
});

test.describe('Public Queries', () => {
  test('can fetch offers', async ({ request }) => {
    const response = await request.post(`${API_URL}/graphql`, {
      data: {
        query: `
          query GetOffers {
            offers(first: 5) {
              edges {
                node {
                  id
                  title
                }
              }
              pageInfo {
                hasNextPage
              }
            }
          }
        `,
      },
    });

    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.errors).toBeUndefined();
    expect(body.data.offers).toBeDefined();
  });

  test('can fetch categories', async ({ request }) => {
    const response = await request.post(`${API_URL}/graphql`, {
      data: {
        query: `
          query GetCategories {
            categories {
              id
              name
              slug
            }
          }
        `,
      },
    });

    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.errors).toBeUndefined();
    expect(body.data.categories).toBeDefined();
  });
});

test.describe('Error Handling', () => {
  test('returns error for invalid query', async ({ request }) => {
    const response = await request.post(`${API_URL}/graphql`, {
      data: {
        query: '{ invalidField }',
      },
    });

    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.errors).toBeDefined();
    expect(body.errors.length).toBeGreaterThan(0);
  });

  test('returns error for unauthenticated protected query', async ({ request }) => {
    const response = await request.post(`${API_URL}/graphql`, {
      data: {
        query: `
          query GetMe {
            me {
              id
              email
            }
          }
        `,
      },
    });

    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    // Should either have errors or null data for unauthenticated request
    expect(body.errors || body.data.me === null).toBeTruthy();
  });
});
