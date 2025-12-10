import { useQuery, useMutation } from '@tanstack/react-query';
import { gql } from '@urql/core';
import { graphqlClient as client } from '@/lib/graphql/client';
import type {
  Booking,
  BookingStatus,
  PaginatedResponse,
} from '@/types';

// ============================================
// GraphQL Queries & Mutations
// ============================================

const GET_BOOKINGS = gql`
  query GetBookings(
    $partnerId: ID!
    $status: BookingStatus
    $establishmentId: ID
    $offerId: ID
    $startDate: DateTime
    $endDate: DateTime
    $page: Int
    $pageSize: Int
  ) {
    bookings(
      partnerId: $partnerId
      status: $status
      establishmentId: $establishmentId
      offerId: $offerId
      startDate: $startDate
      endDate: $endDate
      page: $page
      pageSize: $pageSize
    ) {
      items {
        id
        status
        user {
          id
          firstName
          lastName
          email
          avatar
        }
        offer {
          id
          title
          discount {
            type
            value
          }
        }
        establishment {
          id
          name
          address
        }
        qrCode {
          code
          expiresAt
        }
        checkin {
          checkedInAt
          method
        }
        createdAt
        expiresAt
      }
      total
      page
      pageSize
      hasMore
    }
  }
`;

const GET_BOOKING = gql`
  query GetBooking($id: ID!) {
    booking(id: $id) {
      id
      userId
      offerId
      partnerId
      establishmentId
      status
      user {
        id
        firstName
        lastName
        email
        avatar
      }
      offer {
        id
        title
        discount {
          type
          value
          originalPrice
        }
        images
      }
      establishment {
        id
        name
        address
      }
      qrCode {
        code
        data
        expiresAt
      }
      timeline {
        status
        timestamp
        actor
        metadata
      }
      checkin {
        checkedInAt
        checkedInBy
        method
      }
      cancellation {
        cancelledAt
        cancelledBy
        reason
      }
      createdAt
      updatedAt
      expiresAt
    }
  }
`;

const CHECKIN_BOOKING = gql`
  mutation CheckinBooking($id: ID!, $method: CheckinMethod!) {
    checkinBooking(id: $id, method: $method) {
      id
      status
      checkin {
        checkedInAt
        method
      }
    }
  }
`;

const CANCEL_BOOKING = gql`
  mutation CancelBooking($id: ID!, $reason: String) {
    cancelBooking(id: $id, reason: $reason) {
      id
      status
      cancellation {
        cancelledAt
        reason
      }
    }
  }
`;

const GET_BOOKING_STATS = gql`
  query GetBookingStats($partnerId: ID!, $period: StatsPeriod!) {
    bookingStats(partnerId: $partnerId, period: $period) {
      total
      confirmed
      checkedIn
      cancelled
      noShow
      conversionRate
      byDay {
        date
        count
      }
      byHour {
        hour
        count
      }
    }
  }
`;

// ============================================
// Hooks
// ============================================

interface UseBookingsParams {
  partnerId: string;
  status?: BookingStatus;
  establishmentId?: string;
  offerId?: string;
  startDate?: string;
  endDate?: string;
  page?: number;
  pageSize?: number;
}

export function useBookings({
  partnerId,
  status,
  establishmentId,
  offerId,
  startDate,
  endDate,
  page = 1,
  pageSize = 20,
}: UseBookingsParams) {
  return useQuery({
    queryKey: ['bookings', partnerId, status, establishmentId, offerId, startDate, endDate, page, pageSize],
    queryFn: async () => {
      const result = await client.query(GET_BOOKINGS, {
        partnerId,
        status,
        establishmentId,
        offerId,
        startDate,
        endDate,
        page,
        pageSize,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.bookings as PaginatedResponse<Booking>;
    },
    enabled: !!partnerId,
  });
}

export function useBooking(id: string) {
  return useQuery({
    queryKey: ['booking', id],
    queryFn: async () => {
      const result = await client.query(GET_BOOKING, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.booking as Booking;
    },
    enabled: !!id,
  });
}

export function useCheckinBooking() {
  return useMutation({
    mutationFn: async ({ id, method }: { id: string; method: 'qr_scan' | 'manual' }) => {
      const result = await client.mutation(CHECKIN_BOOKING, { id, method }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.checkinBooking;
    },
  });
}

export function useCancelBooking() {
  return useMutation({
    mutationFn: async ({ id, reason }: { id: string; reason?: string }) => {
      const result = await client.mutation(CANCEL_BOOKING, { id, reason }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.cancelBooking;
    },
  });
}

interface UseBookingStatsParams {
  partnerId: string;
  period: 'day' | 'week' | 'month' | 'year';
}

export function useBookingStats({ partnerId, period }: UseBookingStatsParams) {
  return useQuery({
    queryKey: ['bookingStats', partnerId, period],
    queryFn: async () => {
      const result = await client.query(GET_BOOKING_STATS, {
        partnerId,
        period,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.bookingStats;
    },
    enabled: !!partnerId,
  });
}
