import { useQuery, useMutation } from '@tanstack/react-query';
import { gql } from '@urql/core';
import { graphqlClient as client } from '@/lib/graphql/client';
import type {
  Establishment,
  PaginatedResponse,
  EstablishmentFormData,
} from '@/types';

// ============================================
// GraphQL Queries & Mutations
// ============================================

const GET_ESTABLISHMENTS = gql`
  query GetEstablishments($partnerId: ID!, $page: Int, $pageSize: Int) {
    establishments(partnerId: $partnerId, page: $page, pageSize: $pageSize) {
      items {
        id
        name
        description
        address {
          street
          streetNumber
          postalCode
          city
          country
          formatted
        }
        location {
          type
          coordinates
        }
        contact {
          phone
          email
          website
        }
        images {
          url
          isPrimary
        }
        features
        type
        priceRange
        isActive
        createdAt
      }
      total
      page
      pageSize
      hasMore
    }
  }
`;

const GET_ESTABLISHMENT = gql`
  query GetEstablishment($id: ID!) {
    establishment(id: $id) {
      id
      partnerId
      name
      description
      address {
        street
        streetNumber
        complement
        postalCode
        city
        country
        formatted
      }
      location {
        type
        coordinates
      }
      contact {
        phone
        email
        website
      }
      openingHours {
        dayOfWeek
        open
        close
        isClosed
      }
      images {
        url
        alt
        isPrimary
        order
      }
      features
      type
      priceRange
      isActive
      createdAt
      updatedAt
    }
  }
`;

const CREATE_ESTABLISHMENT = gql`
  mutation CreateEstablishment($partnerId: ID!, $input: CreateEstablishmentInput!) {
    createEstablishment(partnerId: $partnerId, input: $input) {
      id
      name
    }
  }
`;

const UPDATE_ESTABLISHMENT = gql`
  mutation UpdateEstablishment($id: ID!, $input: UpdateEstablishmentInput!) {
    updateEstablishment(id: $id, input: $input) {
      id
      name
    }
  }
`;

const DELETE_ESTABLISHMENT = gql`
  mutation DeleteEstablishment($id: ID!) {
    deleteEstablishment(id: $id)
  }
`;

const TOGGLE_ESTABLISHMENT = gql`
  mutation ToggleEstablishment($id: ID!, $isActive: Boolean!) {
    toggleEstablishment(id: $id, isActive: $isActive) {
      id
      isActive
    }
  }
`;

// ============================================
// Hooks
// ============================================

interface UseEstablishmentsParams {
  partnerId: string;
  page?: number;
  pageSize?: number;
}

export function useEstablishments({ partnerId, page = 1, pageSize = 20 }: UseEstablishmentsParams) {
  return useQuery({
    queryKey: ['establishments', partnerId, page, pageSize],
    queryFn: async () => {
      const result = await client.query(GET_ESTABLISHMENTS, {
        partnerId,
        page,
        pageSize,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.establishments as PaginatedResponse<Establishment>;
    },
    enabled: !!partnerId,
  });
}

export function useEstablishment(id: string) {
  return useQuery({
    queryKey: ['establishment', id],
    queryFn: async () => {
      const result = await client.query(GET_ESTABLISHMENT, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.establishment as Establishment;
    },
    enabled: !!id,
  });
}

export function useCreateEstablishment() {
  return useMutation({
    mutationFn: async ({ partnerId, input }: { partnerId: string; input: EstablishmentFormData }) => {
      const result = await client.mutation(CREATE_ESTABLISHMENT, { partnerId, input }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.createEstablishment;
    },
  });
}

export function useUpdateEstablishment() {
  return useMutation({
    mutationFn: async ({ id, input }: { id: string; input: Partial<EstablishmentFormData> }) => {
      const result = await client.mutation(UPDATE_ESTABLISHMENT, { id, input }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.updateEstablishment;
    },
  });
}

export function useDeleteEstablishment() {
  return useMutation({
    mutationFn: async (id: string) => {
      const result = await client.mutation(DELETE_ESTABLISHMENT, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.deleteEstablishment;
    },
  });
}

export function useToggleEstablishment() {
  return useMutation({
    mutationFn: async ({ id, isActive }: { id: string; isActive: boolean }) => {
      const result = await client.mutation(TOGGLE_ESTABLISHMENT, { id, isActive }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.toggleEstablishment;
    },
  });
}
