import { useQuery, useMutation } from '@tanstack/react-query';
import { gql } from '@urql/core';
import { graphqlClient as client } from '@/lib/graphql/client';
import type {
  Offer,
  OfferStatus,
  PaginatedResponse,
  OfferFormData,
} from '@/types';

// ============================================
// GraphQL Queries & Mutations
// ============================================

const GET_OFFERS = gql`
  query GetOffers($partnerId: ID!, $status: OfferStatus, $page: Int, $pageSize: Int) {
    offers(partnerId: $partnerId, status: $status, page: $page, pageSize: $pageSize) {
      items {
        id
        title
        shortDescription
        status
        isActive
        discount {
          type
          value
          originalPrice
        }
        validity {
          startDate
          endDate
        }
        stats {
          views
          bookings
          checkins
          favorites
        }
        images {
          url
          isPrimary
        }
        createdAt
        publishedAt
      }
      total
      page
      pageSize
      hasMore
    }
  }
`;

const GET_OFFER = gql`
  query GetOffer($id: ID!) {
    offer(id: $id) {
      id
      partnerId
      establishmentId
      title
      description
      shortDescription
      categoryId
      tags
      discount {
        type
        value
        originalPrice
        formula
      }
      conditions {
        type
        value
        label
      }
      termsAndConditions
      validity {
        startDate
        endDate
        timezone
      }
      schedule {
        allDay
        slots {
          dayOfWeek
          startTime
          endTime
        }
      }
      quota {
        total
        perUser
        perDay
        used
      }
      images {
        url
        alt
        isPrimary
        order
      }
      stats {
        views
        bookings
        checkins
        favorites
      }
      status
      isActive
      createdAt
      updatedAt
      publishedAt
    }
  }
`;

const CREATE_OFFER = gql`
  mutation CreateOffer($input: CreateOfferInput!) {
    createOffer(input: $input) {
      id
      title
      status
    }
  }
`;

const UPDATE_OFFER = gql`
  mutation UpdateOffer($id: ID!, $input: UpdateOfferInput!) {
    updateOffer(id: $id, input: $input) {
      id
      title
      status
    }
  }
`;

const PUBLISH_OFFER = gql`
  mutation PublishOffer($id: ID!) {
    publishOffer(id: $id) {
      id
      status
      publishedAt
    }
  }
`;

const PAUSE_OFFER = gql`
  mutation PauseOffer($id: ID!) {
    pauseOffer(id: $id) {
      id
      status
    }
  }
`;

const ARCHIVE_OFFER = gql`
  mutation ArchiveOffer($id: ID!) {
    archiveOffer(id: $id) {
      id
      status
    }
  }
`;

const DELETE_OFFER = gql`
  mutation DeleteOffer($id: ID!) {
    deleteOffer(id: $id)
  }
`;

// ============================================
// Hooks
// ============================================

interface UseOffersParams {
  partnerId: string;
  status?: OfferStatus;
  page?: number;
  pageSize?: number;
}

export function useOffers({ partnerId, status, page = 1, pageSize = 10 }: UseOffersParams) {
  return useQuery({
    queryKey: ['offers', partnerId, status, page, pageSize],
    queryFn: async () => {
      const result = await client.query(GET_OFFERS, {
        partnerId,
        status,
        page,
        pageSize,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.offers as PaginatedResponse<Offer>;
    },
    enabled: !!partnerId,
  });
}

export function useOffer(id: string) {
  return useQuery({
    queryKey: ['offer', id],
    queryFn: async () => {
      const result = await client.query(GET_OFFER, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.offer as Offer;
    },
    enabled: !!id,
  });
}

export function useCreateOffer() {
  return useMutation({
    mutationFn: async (input: OfferFormData) => {
      const result = await client.mutation(CREATE_OFFER, { input }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.createOffer;
    },
  });
}

export function useUpdateOffer() {
  return useMutation({
    mutationFn: async ({ id, input }: { id: string; input: Partial<OfferFormData> }) => {
      const result = await client.mutation(UPDATE_OFFER, { id, input }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.updateOffer;
    },
  });
}

export function usePublishOffer() {
  return useMutation({
    mutationFn: async (id: string) => {
      const result = await client.mutation(PUBLISH_OFFER, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.publishOffer;
    },
  });
}

export function usePauseOffer() {
  return useMutation({
    mutationFn: async (id: string) => {
      const result = await client.mutation(PAUSE_OFFER, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.pauseOffer;
    },
  });
}

export function useArchiveOffer() {
  return useMutation({
    mutationFn: async (id: string) => {
      const result = await client.mutation(ARCHIVE_OFFER, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.archiveOffer;
    },
  });
}

export function useDeleteOffer() {
  return useMutation({
    mutationFn: async (id: string) => {
      const result = await client.mutation(DELETE_OFFER, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.deleteOffer;
    },
  });
}
