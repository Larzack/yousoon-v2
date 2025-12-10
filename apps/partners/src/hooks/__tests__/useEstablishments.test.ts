import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React, { ReactNode } from 'react';
import { useEstablishments, useEstablishment, useCreateEstablishment, useUpdateEstablishment } from '../useEstablishments';

// Mock URQL client
const mockToPromise = vi.fn();

vi.mock('@/lib/graphql/client', () => ({
  graphqlClient: {
    mutation: () => ({ toPromise: mockToPromise }),
    query: () => ({ toPromise: mockToPromise }),
  },
}));

const createQueryClient = () => new QueryClient({
  defaultOptions: {
    queries: { retry: false },
    mutations: { retry: false },
  },
});

let queryClient: QueryClient;

const createWrapper = () => {
  return ({ children }: { children: ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);
};

describe('useEstablishments', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = createQueryClient();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('useEstablishments hook', () => {
    it('should return establishments list when query succeeds', async () => {
      const mockEstablishments = {
        data: {
          establishments: {
            items: [
              {
                id: 'est-1',
                name: 'Le Petit Bistro',
                address: {
                  street: '10 Rue de Paris',
                  city: 'Paris',
                  postalCode: '75001',
                },
                isActive: true,
              },
              {
                id: 'est-2',
                name: 'Café des Arts',
                address: {
                  street: '25 Boulevard Saint-Michel',
                  city: 'Paris',
                  postalCode: '75006',
                },
                isActive: true,
              },
            ],
            total: 2,
            page: 1,
            pageSize: 20,
            hasMore: false,
          },
        },
      };

      mockToPromise.mockResolvedValue(mockEstablishments);

      const { result } = renderHook(
        () => useEstablishments({ partnerId: 'partner-123' }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.items).toHaveLength(2);
      expect(result.current.data?.items[0].name).toBe('Le Petit Bistro');
    });

    it('should handle error state', async () => {
      mockToPromise.mockResolvedValue({ error: { message: 'Network error' } });

      const { result } = renderHook(
        () => useEstablishments({ partnerId: 'partner-123' }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });

    it('should not fetch when partnerId is empty', () => {
      const { result } = renderHook(
        () => useEstablishments({ partnerId: '' }),
        { wrapper: createWrapper() }
      );

      expect(result.current.isFetching).toBe(false);
    });
  });

  describe('useEstablishment hook', () => {
    it('should return single establishment details', async () => {
      const mockEstablishment = {
        data: {
          establishment: {
            id: 'est-1',
            name: 'Le Petit Bistro',
            description: 'Un bistro typiquement parisien',
            address: {
              street: '10 Rue de Paris',
              streetNumber: '10',
              city: 'Paris',
              postalCode: '75001',
              country: 'FR',
              formatted: '10 Rue de Paris, 75001 Paris',
            },
            location: {
              type: 'Point',
              coordinates: [2.3522, 48.8566],
            },
            contact: {
              phone: '+33123456789',
              email: 'contact@petitbistro.fr',
              website: 'https://petitbistro.fr',
            },
            openingHours: [
              { dayOfWeek: 1, open: '09:00', close: '23:00', isClosed: false },
              { dayOfWeek: 2, open: '09:00', close: '23:00', isClosed: false },
            ],
            features: ['terrasse', 'wifi', 'parking'],
            priceRange: 2,
            isActive: true,
          },
        },
      };

      mockToPromise.mockResolvedValue(mockEstablishment);

      const { result } = renderHook(
        () => useEstablishment('est-1'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.name).toBe('Le Petit Bistro');
      expect(result.current.data?.features).toContain('wifi');
    });

    it('should not fetch when id is empty', () => {
      const { result } = renderHook(
        () => useEstablishment(''),
        { wrapper: createWrapper() }
      );

      expect(result.current.isFetching).toBe(false);
    });
  });

  describe('useCreateEstablishment hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => useCreateEstablishment(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });

    it('should create establishment successfully', async () => {
      mockToPromise.mockResolvedValue({
        data: {
          createEstablishment: {
            id: 'new-est-1',
            name: 'Nouveau Restaurant',
            isActive: true,
          },
        },
      });

      const { result } = renderHook(
        () => useCreateEstablishment(),
        { wrapper: createWrapper() }
      );

      const establishmentData = {
        name: 'Nouveau Restaurant',
        street: '5 Avenue des Champs-Élysées',
        postalCode: '75008',
        city: 'Paris',
        country: 'FR',
        openingHours: [],
        features: [],
        images: [],
      };

      await result.current.mutateAsync({
        partnerId: 'partner-123',
        input: establishmentData,
      });

      expect(result.current.isSuccess).toBe(true);
    });
  });

  describe('useUpdateEstablishment hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => useUpdateEstablishment(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });

    it('should update establishment successfully', async () => {
      mockToPromise.mockResolvedValue({
        data: {
          updateEstablishment: {
            id: 'est-1',
            name: 'Le Petit Bistro Updated',
            isActive: true,
          },
        },
      });

      const { result } = renderHook(
        () => useUpdateEstablishment(),
        { wrapper: createWrapper() }
      );

      await result.current.mutateAsync({
        id: 'est-1',
        input: { name: 'Le Petit Bistro Updated' },
      });

      expect(result.current.isSuccess).toBe(true);
    });
  });
});
