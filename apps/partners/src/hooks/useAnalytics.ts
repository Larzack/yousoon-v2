import { useQuery } from '@tanstack/react-query';
import { gql } from '@urql/core';
import { graphqlClient as client } from '@/lib/graphql/client';
import type { AnalyticsSummary, DailyStats, TopOffer } from '@/types';

// ============================================
// GraphQL Queries
// ============================================

const GET_ANALYTICS_SUMMARY = gql`
  query GetAnalyticsSummary($partnerId: ID!, $period: StatsPeriod!) {
    analyticsSummary(partnerId: $partnerId, period: $period) {
      period {
        start
        end
      }
      totalViews
      totalBookings
      totalCheckins
      conversionRate
      revenue
      trends {
        views
        bookings
        checkins
      }
    }
  }
`;

const GET_DAILY_STATS = gql`
  query GetDailyStats($partnerId: ID!, $startDate: DateTime!, $endDate: DateTime!) {
    dailyStats(partnerId: $partnerId, startDate: $startDate, endDate: $endDate) {
      date
      views
      bookings
      checkins
      revenue
    }
  }
`;

const GET_TOP_OFFERS = gql`
  query GetTopOffers($partnerId: ID!, $period: StatsPeriod!, $limit: Int) {
    topOffers(partnerId: $partnerId, period: $period, limit: $limit) {
      offerId
      title
      views
      bookings
      conversionRate
    }
  }
`;

const GET_HOURLY_DISTRIBUTION = gql`
  query GetHourlyDistribution($partnerId: ID!, $period: StatsPeriod!) {
    hourlyDistribution(partnerId: $partnerId, period: $period) {
      hour
      bookings
      checkins
    }
  }
`;

const GET_WEEKLY_DISTRIBUTION = gql`
  query GetWeeklyDistribution($partnerId: ID!, $period: StatsPeriod!) {
    weeklyDistribution(partnerId: $partnerId, period: $period) {
      dayOfWeek
      bookings
      checkins
    }
  }
`;

const GET_FUNNEL_DATA = gql`
  query GetFunnelData($partnerId: ID!, $period: StatsPeriod!) {
    funnelData(partnerId: $partnerId, period: $period) {
      views
      favorites
      bookings
      checkins
      reviews
    }
  }
`;

// ============================================
// Hooks
// ============================================

type Period = 'day' | 'week' | 'month' | 'quarter' | 'year';

interface UseAnalyticsParams {
  partnerId: string;
  period: Period;
}

export function useAnalyticsSummary({ partnerId, period }: UseAnalyticsParams) {
  return useQuery({
    queryKey: ['analyticsSummary', partnerId, period],
    queryFn: async () => {
      const result = await client.query(GET_ANALYTICS_SUMMARY, {
        partnerId,
        period,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.analyticsSummary as AnalyticsSummary;
    },
    enabled: !!partnerId,
  });
}

interface UseDailyStatsParams {
  partnerId: string;
  startDate: string;
  endDate: string;
}

export function useDailyStats({ partnerId, startDate, endDate }: UseDailyStatsParams) {
  return useQuery({
    queryKey: ['dailyStats', partnerId, startDate, endDate],
    queryFn: async () => {
      const result = await client.query(GET_DAILY_STATS, {
        partnerId,
        startDate,
        endDate,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.dailyStats as DailyStats[];
    },
    enabled: !!partnerId && !!startDate && !!endDate,
  });
}

interface UseTopOffersParams {
  partnerId: string;
  period: Period;
  limit?: number;
}

export function useTopOffers({ partnerId, period, limit = 10 }: UseTopOffersParams) {
  return useQuery({
    queryKey: ['topOffers', partnerId, period, limit],
    queryFn: async () => {
      const result = await client.query(GET_TOP_OFFERS, {
        partnerId,
        period,
        limit,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.topOffers as TopOffer[];
    },
    enabled: !!partnerId,
  });
}

export function useHourlyDistribution({ partnerId, period }: UseAnalyticsParams) {
  return useQuery({
    queryKey: ['hourlyDistribution', partnerId, period],
    queryFn: async () => {
      const result = await client.query(GET_HOURLY_DISTRIBUTION, {
        partnerId,
        period,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.hourlyDistribution as { hour: number; bookings: number; checkins: number }[];
    },
    enabled: !!partnerId,
  });
}

export function useWeeklyDistribution({ partnerId, period }: UseAnalyticsParams) {
  return useQuery({
    queryKey: ['weeklyDistribution', partnerId, period],
    queryFn: async () => {
      const result = await client.query(GET_WEEKLY_DISTRIBUTION, {
        partnerId,
        period,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.weeklyDistribution as { dayOfWeek: number; bookings: number; checkins: number }[];
    },
    enabled: !!partnerId,
  });
}

interface FunnelData {
  views: number;
  favorites: number;
  bookings: number;
  checkins: number;
  reviews: number;
}

export function useFunnelData({ partnerId, period }: UseAnalyticsParams) {
  return useQuery({
    queryKey: ['funnelData', partnerId, period],
    queryFn: async () => {
      const result = await client.query(GET_FUNNEL_DATA, {
        partnerId,
        period,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.funnelData as FunnelData;
    },
    enabled: !!partnerId,
  });
}
