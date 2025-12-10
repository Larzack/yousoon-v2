import { useQuery, useMutation } from '@tanstack/react-query';
import { gql } from '@urql/core';
import { graphqlClient as client } from '@/lib/graphql/client';
import type { TeamMember, TeamRole, TeamMemberFormData } from '@/types';

// ============================================
// GraphQL Queries & Mutations
// ============================================

const GET_TEAM_MEMBERS = gql`
  query GetTeamMembers($partnerId: ID!) {
    teamMembers(partnerId: $partnerId) {
      id
      userId
      email
      firstName
      lastName
      avatar
      role
      status
      invitedAt
      joinedAt
    }
  }
`;

const INVITE_TEAM_MEMBER = gql`
  mutation InviteTeamMember($partnerId: ID!, $input: InviteTeamMemberInput!) {
    inviteTeamMember(partnerId: $partnerId, input: $input) {
      id
      email
      role
      status
      invitedAt
    }
  }
`;

const UPDATE_TEAM_MEMBER_ROLE = gql`
  mutation UpdateTeamMemberRole($id: ID!, $role: TeamRole!) {
    updateTeamMemberRole(id: $id, role: $role) {
      id
      role
    }
  }
`;

const REMOVE_TEAM_MEMBER = gql`
  mutation RemoveTeamMember($id: ID!) {
    removeTeamMember(id: $id)
  }
`;

const RESEND_INVITATION = gql`
  mutation ResendInvitation($id: ID!) {
    resendInvitation(id: $id) {
      id
      invitedAt
    }
  }
`;

// ============================================
// Hooks
// ============================================

export function useTeamMembers(partnerId: string) {
  return useQuery({
    queryKey: ['teamMembers', partnerId],
    queryFn: async () => {
      const result = await client.query(GET_TEAM_MEMBERS, { partnerId }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.teamMembers as TeamMember[];
    },
    enabled: !!partnerId,
  });
}

export function useInviteTeamMember() {
  return useMutation({
    mutationFn: async ({ partnerId, input }: { partnerId: string; input: TeamMemberFormData }) => {
      const result = await client.mutation(INVITE_TEAM_MEMBER, {
        partnerId,
        input,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.inviteTeamMember as TeamMember;
    },
  });
}

export function useUpdateTeamMemberRole() {
  return useMutation({
    mutationFn: async ({ id, role }: { id: string; role: TeamRole }) => {
      const result = await client.mutation(UPDATE_TEAM_MEMBER_ROLE, { id, role }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.updateTeamMemberRole;
    },
  });
}

export function useRemoveTeamMember() {
  return useMutation({
    mutationFn: async (id: string) => {
      const result = await client.mutation(REMOVE_TEAM_MEMBER, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.removeTeamMember;
    },
  });
}

export function useResendInvitation() {
  return useMutation({
    mutationFn: async (id: string) => {
      const result = await client.mutation(RESEND_INVITATION, { id }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.resendInvitation;
    },
  });
}
