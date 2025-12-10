import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React, { ReactNode } from 'react';
import { useTeamMembers, useInviteTeamMember, useUpdateTeamMemberRole, useRemoveTeamMember } from '../useTeam';
import { TeamRole } from '@/types';

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

describe('useTeam', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    queryClient = createQueryClient();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('useTeamMembers hook', () => {
    it('should return team members list when query succeeds', async () => {
      const mockTeamMembers = {
        data: {
          teamMembers: [
            {
              id: 'member-1',
              email: 'admin@bistro.fr',
              firstName: 'Pierre',
              lastName: 'Martin',
              role: 'ADMIN',
              status: 'active',
              joinedAt: '2025-01-01T00:00:00Z',
            },
            {
              id: 'member-2',
              email: 'manager@bistro.fr',
              firstName: 'Marie',
              lastName: 'Dupont',
              role: 'MANAGER',
              status: 'active',
              joinedAt: '2025-02-01T00:00:00Z',
            },
            {
              id: 'member-3',
              email: 'staff@bistro.fr',
              firstName: null,
              lastName: null,
              role: 'STAFF',
              status: 'pending',
              invitedAt: '2025-12-01T00:00:00Z',
            },
          ],
        },
      };

      mockToPromise.mockResolvedValue(mockTeamMembers);

      const { result } = renderHook(
        () => useTeamMembers('partner-123'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      expect(result.current.data).toHaveLength(3);
      expect(result.current.data?.[0].role).toBe('ADMIN');
      expect(result.current.data?.[2].status).toBe('pending');
    });

    it('should handle error state', async () => {
      mockToPromise.mockResolvedValue({ error: { message: 'Network error' } });

      const { result } = renderHook(
        () => useTeamMembers('partner-123'),
        { wrapper: createWrapper() }
      );

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });
    });

    it('should not fetch when partnerId is empty', () => {
      const { result } = renderHook(
        () => useTeamMembers(''),
        { wrapper: createWrapper() }
      );

      expect(result.current.isFetching).toBe(false);
    });
  });

  describe('useInviteTeamMember hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => useInviteTeamMember(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });

    it('should invite team member successfully', async () => {
      mockToPromise.mockResolvedValue({
        data: {
          inviteTeamMember: {
            id: 'new-member-1',
            email: 'new@bistro.fr',
            role: 'STAFF',
            status: 'pending',
            invitedAt: new Date().toISOString(),
          },
        },
      });

      const { result } = renderHook(
        () => useInviteTeamMember(),
        { wrapper: createWrapper() }
      );

      await result.current.mutateAsync({
        partnerId: 'partner-123',
        input: {
          email: 'new@bistro.fr',
          role: TeamRole.STAFF,
        },
      });

      expect(result.current.isSuccess).toBe(true);
    });
  });

  describe('useUpdateTeamMemberRole hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => useUpdateTeamMemberRole(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });

    it('should update team member role successfully', async () => {
      mockToPromise.mockResolvedValue({
        data: {
          updateTeamMemberRole: {
            id: 'member-2',
            role: 'ADMIN',
          },
        },
      });

      const { result } = renderHook(
        () => useUpdateTeamMemberRole(),
        { wrapper: createWrapper() }
      );

      await result.current.mutateAsync({
        id: 'member-2',
        role: TeamRole.ADMIN,
      });

      expect(result.current.isSuccess).toBe(true);
    });
  });

  describe('useRemoveTeamMember hook', () => {
    it('should provide mutate function', () => {
      const { result } = renderHook(
        () => useRemoveTeamMember(),
        { wrapper: createWrapper() }
      );

      expect(result.current.mutate).toBeDefined();
      expect(typeof result.current.mutate).toBe('function');
    });

    it('should remove team member successfully', async () => {
      mockToPromise.mockResolvedValue({
        data: {
          removeTeamMember: true,
        },
      });

      const { result } = renderHook(
        () => useRemoveTeamMember(),
        { wrapper: createWrapper() }
      );

      await result.current.mutateAsync('member-3');

      expect(result.current.isSuccess).toBe(true);
    });
  });
});
