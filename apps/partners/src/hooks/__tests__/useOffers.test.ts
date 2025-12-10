import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React, { ReactNode } from 'react';
import { useOffers, useOffer, useCreateOffer, useUpdateOffer, usePublishOffer } from '../useOffers';
import { OfferStatus, DiscountType } from '@/types';

// Mock URQL client
const mockToPromise = vi.fn();
const mockMutation = vi.fn(() => ({ toPromise: mockToPromise }));
const mockQueryFn = vi.fn(() => ({ toPromise: mockToPromise }));

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

describe('useOffers', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = createQueryClient();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('useOffers hook', () => {
    it('should return offers list when query succeeds', async () => {
      const mockOffers = {
        data: {
          offers: {
            items: [
              {
                id: 'offer-1',
                title: 'Happy Hour -50%',
                status: 'active',
                isActive: true,
                discount: { type: 'percentage', value: 50 },
                stats: { views: 100, bookings: 25 },
              },
              {
                id: 'offer-2',
                title: 'Menu Déjeuner',
                status: 'draft',
                isActive: false,
                discount: { type: 'fixed', value: 10 },
                stats: { views: 50, bookings: 10 },
              },
            ],
            total: 2,
            page: 1,
            pageSize: 10,
            hasMore: false,
          },
        },
      };

      mockToPromise.mockResolvedValue(mockOffers);

      const { result } = renderHook(
        () => useOffers({ partnerId: 'partner-123' }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.items).toHaveLength(2);
      expect(result.current.data?.items[0].title).toBe('Happy Hour -50%');
    });

    it('should filter offers by status', async () => {
      const mockOffers = {
        data: {
          offers: {
            items: [
              { id: 'offer-1', title: 'Active Offer', status: 'active' },
            ],
            total: 1,
            page: 1,
            pageSize: 10,
            hasMore: false,
          },
        },
      };

      mockToPromise.mockResolvedValue(mockOffers);

      const { result } = renderHook(
        () => useOffers({ partnerId: 'partner-123', status: OfferStatus.ACTIVE }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(mockQueryFn).toHaveBeenCalled();
    });

    it('should handle pagination', async () => {
      const mockOffers = {
        data: {
          offers: {
            items: [],
            total: 25,
            page: 2,
            pageSize: 10,
            hasMore: true,
          },
        },
      };

      mockToPromise.mockResolvedValue(mockOffers);

      const { result } = renderHook(
        () => useOffers({ partnerId: 'partner-123', page: 2, pageSize: 10 }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.total).toBe(25);
      expect(result.current.data?.hasMore).toBe(true);
    });

    it('should handle error state', async () => {
      mockToPromise.mockResolvedValue({ error: { message: 'Network error' } });

      const { result } = renderHook(
        () => useOffers({ partnerId: 'partner-123' }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });

  describe('useOffer hook', () => {
    it('should return single offer details', async () => {
      const mockOffer = {
        data: {
          offer: {
            id: 'offer-1',
            title: 'Happy Hour -50%',
            description: 'Profitez de -50% sur toutes les boissons',
            status: 'active',
            discount: { type: 'percentage', value: 50, originalPrice: 10 },
            validity: {
              startDate: '2025-01-01',
              endDate: '2025-12-31',
            },
            schedule: {
              allDay: false,
              slots: [{ dayOfWeek: 5, startTime: '17:00', endTime: '20:00' }],
            },
          },
        },
      };

      mockToPromise.mockResolvedValue(mockOffer);

      const { result } = renderHook(
        () => useOffer('offer-1'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.title).toBe('Happy Hour -50%');
      expect(result.current.data?.discount.value).toBe(50);
    });

    it('should not fetch when id is empty string', () => {
      mockToPromise.mockResolvedValue({ data: { offer: null } });
      
      const { result } = renderHook(
        () => useOffer(''),
        { wrapper: createWrapper() }
      );

      // Query is disabled when id is falsy
      expect(result.current.isFetching).toBe(false);
    });
  });

  describe('useCreateOffer hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => useCreateOffer(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });

    it('should provide mutateAsync function', () => {
      const { result } = renderHook(
        () => useCreateOffer(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutateAsync).toBeDefined();
      expect(typeof result.current.mutateAsync).toBe('function');
    });

    it('should create offer successfully', async () => {
      const newOffer = {
        data: {
          createOffer: {
            id: 'new-offer-1',
            title: 'New Offer',
            status: 'draft',
          },
        },
      };

      mockToPromise.mockResolvedValue(newOffer);

      const { result } = renderHook(
        () => useCreateOffer(),
        { wrapper: createWrapper() }
      );

      const offerData = {
        title: 'New Offer',
        description: 'Test description',
        categoryId: 'cat-1',
        establishmentId: 'est-1',
        tags: ['promo'],
        discountType: DiscountType.PERCENTAGE,
        discountValue: 20,
        conditions: [],
        startDate: '2025-01-01',
        endDate: '2025-12-31',
        allDay: true,
        slots: [],
        images: [],
      };

      await result.current.mutateAsync(offerData);

      expect(mockMutation).toHaveBeenCalled();
    });
  });

  describe('useUpdateOffer hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => useUpdateOffer(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });
  });

  describe('usePublishOffer hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => usePublishOffer(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });

    it('should call mutation with offer id', async () => {
      mockToPromise.mockResolvedValue({
        data: {
          publishOffer: {
            id: 'offer-1',
            status: 'active',
            publishedAt: new Date().toISOString(),
          },
        },
      });

      const { result } = renderHook(
        () => usePublishOffer(),
        { wrapper: createWrapper() }
      );

      await result.current.mutateAsync('offer-1');

      expect(mockMutation).toHaveBeenCalled();
    });
  });
});
