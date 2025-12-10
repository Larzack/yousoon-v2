import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React, { ReactNode } from 'react';
import { useBookings, useBooking, useCheckinBooking, useCancelBooking } from '../useBookings';
import { BookingStatus } from '@/types';

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

describe('useBookings', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = createQueryClient();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('useBookings hook', () => {
    it('should return bookings list when query succeeds', async () => {
      const mockBookings = {
        data: {
          bookings: {
            items: [
              {
                id: 'booking-1',
                status: 'CONFIRMED',
                user: {
                  id: 'user-1',
                  firstName: 'John',
                  lastName: 'Doe',
                  email: 'john@example.com',
                },
                offer: {
                  id: 'offer-1',
                  title: 'Happy Hour',
                },
                createdAt: '2025-12-01T10:00:00Z',
              },
              {
                id: 'booking-2',
                status: 'CHECKED_IN',
                user: {
                  id: 'user-2',
                  firstName: 'Jane',
                  lastName: 'Smith',
                  email: 'jane@example.com',
                },
                offer: {
                  id: 'offer-2',
                  title: 'Menu Déjeuner',
                },
                createdAt: '2025-12-01T12:00:00Z',
              },
            ],
            total: 2,
            page: 1,
            pageSize: 10,
            hasMore: false,
          },
        },
      };

      mockToPromise.mockResolvedValue(mockBookings);

      const { result } = renderHook(
        () => useBookings({ partnerId: 'partner-123' }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.items).toHaveLength(2);
      expect(result.current.data?.items[0].status).toBe('CONFIRMED');
    });

    it('should filter bookings by status', async () => {
      const mockBookings = {
        data: {
          bookings: {
            items: [
              { id: 'booking-1', status: 'CHECKED_IN' },
            ],
            total: 1,
            page: 1,
            pageSize: 10,
            hasMore: false,
          },
        },
      };

      mockToPromise.mockResolvedValue(mockBookings);

      const { result } = renderHook(
        () => useBookings({ partnerId: 'partner-123', status: BookingStatus.CHECKED_IN }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.items).toHaveLength(1);
    });

    it('should filter bookings by date range', async () => {
      const mockBookings = {
        data: {
          bookings: {
            items: [],
            total: 0,
            page: 1,
            pageSize: 10,
            hasMore: false,
          },
        },
      };

      mockToPromise.mockResolvedValue(mockBookings);

      const { result } = renderHook(
        () => useBookings({
          partnerId: 'partner-123',
          startDate: '2025-01-01',
          endDate: '2025-01-31',
        }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.total).toBe(0);
    });

    it('should handle error state', async () => {
      mockToPromise.mockResolvedValue({ error: { message: 'Network error' } });

      const { result } = renderHook(
        () => useBookings({ partnerId: 'partner-123' }),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });
  });

  describe('useBooking hook', () => {
    it('should return single booking details', async () => {
      const mockBooking = {
        data: {
          booking: {
            id: 'booking-1',
            status: 'CONFIRMED',
            user: {
              id: 'user-1',
              firstName: 'John',
              lastName: 'Doe',
              email: 'john@example.com',
            },
            offer: {
              id: 'offer-1',
              title: 'Happy Hour',
              discount: { type: 'percentage', value: 50 },
            },
            qrCode: {
              code: 'QR123456',
              expiresAt: '2025-12-01T12:00:00Z',
            },
            timeline: [
              { status: 'PENDING', timestamp: '2025-12-01T10:00:00Z', actor: 'user' },
              { status: 'CONFIRMED', timestamp: '2025-12-01T10:01:00Z', actor: 'system' },
            ],
            createdAt: '2025-12-01T10:00:00Z',
          },
        },
      };

      mockToPromise.mockResolvedValue(mockBooking);

      const { result } = renderHook(
        () => useBooking('booking-1'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data?.status).toBe('CONFIRMED');
      expect(result.current.data?.qrCode?.code).toBe('QR123456');
    });

    it('should not fetch when id is empty', () => {
      mockToPromise.mockResolvedValue({ data: { booking: null } });
      
      const { result } = renderHook(
        () => useBooking(''),
        { wrapper: createWrapper() }
      );

      expect(result.current.isFetching).toBe(false);
    });
  });

  describe('useCheckinBooking hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => useCheckinBooking(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });

    it('should checkin booking successfully', async () => {
      mockToPromise.mockResolvedValue({
        data: {
          checkinBooking: {
            id: 'booking-1',
            status: 'CHECKED_IN',
            checkin: {
              checkedInAt: new Date().toISOString(),
              method: 'qr_scan',
            },
          },
        },
      });

      const { result } = renderHook(
        () => useCheckinBooking(),
        { wrapper: createWrapper() }
      );

      await result.current.mutateAsync({ id: 'booking-1', method: 'qr_scan' });

      // Mutation was called
      expect(result.current.isSuccess).toBe(true);
    });
  });

  describe('useCancelBooking hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => useCancelBooking(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });

    it('should cancel booking successfully', async () => {
      mockToPromise.mockResolvedValue({
        data: {
          cancelBooking: {
            id: 'booking-1',
            status: 'CANCELLED',
            cancellation: {
              cancelledAt: new Date().toISOString(),
              cancelledBy: 'partner',
              reason: 'Customer request',
            },
          },
        },
      });

      const { result } = renderHook(
        () => useCancelBooking(),
        { wrapper: createWrapper() }
      );

      await result.current.mutateAsync({ id: 'booking-1', reason: 'Customer request' });

      expect(result.current.isSuccess).toBe(true);
    });
  });
});
