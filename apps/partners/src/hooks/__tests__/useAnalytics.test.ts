import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React, { ReactNode } from 'react';
import { useAnalyticsSummary, useDailyStats, useTopOffers, useFunnelData } from '../useAnalytics';

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

describe('useAnalytics', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = createQueryClient();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('useAnalyticsSummary hook', () => {
    it('should return analytics summary when query succeeds', async () => {
      const mockAnalytics = {
        data: {
          analyticsSummary: {
            period: { start: '2025-01-01', end: '2025-12-31' },
            totalViews: 15000,
            totalBookings: 500,
            totalCheckins: 450,
            conversionRate: 3.33,
            revenue: 25000,
            trends: {
              views: 12.5,
              bookings: 8.2,
              checkins: 10.0,
            },
          },
        },
      };

      mockToPromise.mockResolvedValue(mockAnalytics);

      const { result } = renderHook(
        () => useAnalyticsSummary({
          partnerId: 'partner-123',
          period: 'month',
        }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.totalViews).toBe(15000);
      expect(result.current.data?.conversionRate).toBe(3.33);
    });

    it('should handle error state', async () => {
      mockToPromise.mockResolvedValue({ error: { message: 'Network error' } });

      const { result } = renderHook(
        () => useAnalyticsSummary({
          partnerId: 'partner-123',
          period: 'week',
        }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });

  describe('useDailyStats hook', () => {
    it('should return daily statistics', async () => {
      const mockDailyStats = {
        data: {
          dailyStats: [
            { date: '2025-12-01', views: 100, bookings: 5, checkins: 4, revenue: 500 },
            { date: '2025-12-02', views: 120, bookings: 8, checkins: 7, revenue: 700 },
          ],
        },
      };

      mockToPromise.mockResolvedValue(mockDailyStats);

      const { result } = renderHook(
        () => useDailyStats({
          partnerId: 'partner-123',
          startDate: '2025-12-01',
          endDate: '2025-12-31',
        }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data).toHaveLength(2);
      expect(result.current.data?.[0].views).toBe(100);
    });

    it('should not fetch when dates are missing', () => {
      const { result } = renderHook(
        () => useDailyStats({
          partnerId: 'partner-123',
          startDate: '',
          endDate: '',
        }),
        { wrapper: createWrapper() }
      );

      expect(result.current.isFetching).toBe(false);
    });
  });

  describe('useTopOffers hook', () => {
    it('should return top offers', async () => {
      const mockTopOffers = {
        data: {
          topOffers: [
            { offerId: 'offer-1', title: 'Happy Hour', views: 1000, bookings: 50, conversionRate: 5.0 },
            { offerId: 'offer-2', title: 'Menu Déjeuner', views: 800, bookings: 40, conversionRate: 5.0 },
          ],
        },
      };

      mockToPromise.mockResolvedValue(mockTopOffers);

      const { result } = renderHook(
        () => useTopOffers({
          partnerId: 'partner-123',
          period: 'month',
          limit: 10,
        }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data).toHaveLength(2);
      expect(result.current.data?.[0].title).toBe('Happy Hour');
    });
  });

  describe('useFunnelData hook', () => {
    it('should return funnel data', async () => {
      const mockFunnelData = {
        data: {
          funnelData: {
            views: 10000,
            favorites: 500,
            bookings: 200,
            checkins: 180,
            reviews: 50,
          },
        },
      };

      mockToPromise.mockResolvedValue(mockFunnelData);

      const { result } = renderHook(
        () => useFunnelData({
          partnerId: 'partner-123',
          period: 'month',
        }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.views).toBe(10000);
      expect(result.current.data?.checkins).toBe(180);
    });
  });
});
