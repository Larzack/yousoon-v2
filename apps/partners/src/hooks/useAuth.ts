import { useCallback } from 'react';
import { useMutation } from '@tanstack/react-query';
import { gql } from '@urql/core';
import { graphqlClient as client } from '@/lib/graphql/client';
import { useAuthStore, User, Partner } from '@/stores/authStore';
import type { LoginFormData, RegisterFormData } from '@/types';

// ============================================
// GraphQL Mutations
// ============================================

const LOGIN = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      accessToken
      refreshToken
      user {
        id
        email
        firstName
        lastName
        avatar
        role
      }
      partner {
        id
        company {
          name
          tradeName
          siret
        }
        branding {
          logo
          primaryColor
        }
        status
      }
    }
  }
`;

const REGISTER = gql`
  mutation Register($input: RegisterPartnerInput!) {
    registerPartner(input: $input) {
      accessToken
      refreshToken
      user {
        id
        email
        firstName
        lastName
      }
      partner {
        id
        company {
          name
          siret
        }
        status
      }
    }
  }
`;

const LOGOUT = gql`
  mutation Logout {
    logout
  }
`;

const REFRESH_TOKEN = gql`
  mutation RefreshToken($refreshToken: String!) {
    refreshToken(refreshToken: $refreshToken) {
      accessToken
      refreshToken
    }
  }
`;

const FORGOT_PASSWORD = gql`
  mutation ForgotPassword($email: String!) {
    forgotPassword(email: $email)
  }
`;

const RESET_PASSWORD = gql`
  mutation ResetPassword($token: String!, $password: String!) {
    resetPassword(token: $token, password: $password)
  }
`;

const GET_ME = gql`
  query GetMe {
    me {
      user {
        id
        email
        firstName
        lastName
        avatar
        role
      }
      partner {
        id
        company {
          name
          tradeName
          siret
        }
        branding {
          logo
          primaryColor
        }
        status
        stats {
          totalOffers
          activeOffers
          totalBookings
          totalCheckins
          avgRating
          reviewCount
        }
      }
    }
  }
`;

// ============================================
// Hook
// ============================================

interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
  partner: Partner;
}

export function useAuth() {
  const { user, partner, setAuth, logout: clearAuth, isAuthenticated } = useAuthStore();

  // Login mutation
  const loginMutation = useMutation({
    mutationFn: async (data: LoginFormData): Promise<AuthResponse> => {
      const result = await client.mutation(LOGIN, {
        email: data.email,
        password: data.password,
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.login as AuthResponse;
    },
    onSuccess: (data) => {
      // Store tokens
      localStorage.setItem('accessToken', data.accessToken);
      localStorage.setItem('refreshToken', data.refreshToken);
      
      // Update auth store
      setAuth({
        user: data.user,
        partner: data.partner,
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
      });
    },
  });

  // Register mutation
  const registerMutation = useMutation({
    mutationFn: async (data: RegisterFormData): Promise<AuthResponse> => {
      const result = await client.mutation(REGISTER, {
        input: {
          firstName: data.firstName,
          lastName: data.lastName,
          email: data.email,
          phone: data.phone,
          password: data.password,
          company: {
            name: data.companyName,
            tradeName: data.tradeName,
            siret: data.siret,
          },
          category: data.category,
        },
      }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.registerPartner as AuthResponse;
    },
    onSuccess: (data) => {
      localStorage.setItem('accessToken', data.accessToken);
      localStorage.setItem('refreshToken', data.refreshToken);
      setAuth({
        user: data.user,
        partner: data.partner,
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
      });
    },
  });

  // Logout mutation
  const logoutMutation = useMutation({
    mutationFn: async () => {
      await client.mutation(LOGOUT, {}).toPromise();
    },
    onSettled: () => {
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
      clearAuth();
    },
  });

  // Forgot password mutation
  const forgotPasswordMutation = useMutation({
    mutationFn: async (email: string) => {
      const result = await client.mutation(FORGOT_PASSWORD, { email }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.forgotPassword;
    },
  });

  // Reset password mutation
  const resetPasswordMutation = useMutation({
    mutationFn: async ({ token, password }: { token: string; password: string }) => {
      const result = await client.mutation(RESET_PASSWORD, { token, password }).toPromise();

      if (result.error) {
        throw new Error(result.error.message);
      }

      return result.data?.resetPassword;
    },
  });

  // Refresh token
  const refreshToken = useCallback(async () => {
    const storedRefreshToken = localStorage.getItem('refreshToken');
    if (!storedRefreshToken) {
      throw new Error('No refresh token');
    }

    const result = await client.mutation(REFRESH_TOKEN, {
      refreshToken: storedRefreshToken,
    }).toPromise();

    if (result.error) {
      clearAuth();
      throw new Error(result.error.message);
    }

    const tokens = result.data?.refreshToken;
    localStorage.setItem('accessToken', tokens.accessToken);
    localStorage.setItem('refreshToken', tokens.refreshToken);

    return tokens.accessToken;
  }, [clearAuth]);

  // Fetch current user
  const fetchMe = useCallback(async () => {
    const result = await client.query(GET_ME, {}).toPromise();

    if (result.error) {
      throw new Error(result.error.message);
    }

    const data = result.data?.me;
    if (data) {
      setAuth({
        user: data.user,
        partner: data.partner,
        accessToken: localStorage.getItem('accessToken') || '',
        refreshToken: localStorage.getItem('refreshToken') || '',
      });
    }

    return data;
  }, [setAuth]);

  return {
    // State
    user,
    partner,
    isAuthenticated,
    
    // Actions
    login: loginMutation.mutateAsync,
    register: registerMutation.mutateAsync,
    logout: logoutMutation.mutate,
    forgotPassword: forgotPasswordMutation.mutateAsync,
    resetPassword: resetPasswordMutation.mutateAsync,
    refreshToken,
    fetchMe,
    
    // Loading states
    isLoggingIn: loginMutation.isPending,
    isRegistering: registerMutation.isPending,
    isLoggingOut: logoutMutation.isPending,
    
    // Errors
    loginError: loginMutation.error,
    registerError: registerMutation.error,
  };
}
