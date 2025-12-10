import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React, { ReactNode } from 'react';
import { useAuth } from '../useAuth';
import { useAuthStore } from '@/stores/authStore';

// Mock GraphQL client
vi.mock('@/lib/graphql/client', () => ({
  graphqlClient: {
    mutation: vi.fn(),
    query: vi.fn(),
  },
}));

// Mock auth store
vi.mock('@/stores/authStore', () => ({
  useAuthStore: vi.fn(),
}));

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: false },
    mutations: { retry: false },
  },
});

const wrapper = ({ children }: { children: ReactNode }) => 
  React.createElement(QueryClientProvider, { client: queryClient }, children);

describe('useAuth', () => {
  const mockSetAuth = vi.fn();
  const mockClearAuth = vi.fn();
  const mockSetLoading = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    queryClient.clear();
    
    (useAuthStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      user: null,
      partner: null,
      accessToken: null,
      isAuthenticated: false,
      isLoading: false,
      setAuth: mockSetAuth,
      clearAuth: mockClearAuth,
      setLoading: mockSetLoading,
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('initial state', () => {
    it('should return unauthenticated state initially', () => {
      const { result } = renderHook(() => useAuth(), { wrapper });

      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.user).toBeNull();
      expect(result.current.partner).toBeNull();
    });
  });

  describe('login', () => {
    it('should provide login function', () => {
      const { result } = renderHook(() => useAuth(), { wrapper });

      expect(result.current.login).toBeDefined();
      expect(typeof result.current.login).toBe('function');
    });
  });

  describe('register', () => {
    it('should provide register function', () => {
      const { result } = renderHook(() => useAuth(), { wrapper });

      expect(result.current.register).toBeDefined();
      expect(typeof result.current.register).toBe('function');
    });
  });

  describe('logout', () => {
    it('should provide logout function', () => {
      const { result } = renderHook(() => useAuth(), { wrapper });

      expect(result.current.logout).toBeDefined();
      expect(typeof result.current.logout).toBe('function');
    });
  });

  describe('forgotPassword', () => {
    it('should provide forgotPassword function', () => {
      const { result } = renderHook(() => useAuth(), { wrapper });

      expect(result.current.forgotPassword).toBeDefined();
      expect(typeof result.current.forgotPassword).toBe('function');
    });
  });

  describe('resetPassword', () => {
    it('should provide resetPassword function', () => {
      const { result } = renderHook(() => useAuth(), { wrapper });

      expect(result.current.resetPassword).toBeDefined();
      expect(typeof result.current.resetPassword).toBe('function');
    });
  });

  describe('fetchMe', () => {
    it('should provide fetchMe function', () => {
      const { result } = renderHook(() => useAuth(), { wrapper });

      expect(result.current.fetchMe).toBeDefined();
      expect(typeof result.current.fetchMe).toBe('function');
    });
  });
});

describe('useAuth store integration', () => {
  const mockUser = {
    id: 'user-123',
    email: 'test@example.com',
    firstName: 'John',
    lastName: 'Doe',
    avatar: null,
    role: 'admin',
  };

  const mockPartner = {
    id: 'partner-123',
    company: {
      name: 'Test Company',
      tradeName: 'Test',
      siret: '12345678901234',
    },
    branding: {
      logo: null,
      primaryColor: '#E99B27',
    },
    status: 'active',
  };

  const wrapper = ({ children }: { children: ReactNode }) => 
    React.createElement(QueryClientProvider, { client: queryClient }, children);

  beforeEach(() => {
    vi.clearAllMocks();
    queryClient.clear();
  });

  it('should return authenticated state when user is logged in', () => {
    (useAuthStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      user: mockUser,
      partner: mockPartner,
      accessToken: 'valid-token',
      isAuthenticated: true,
      isLoading: false,
      setAuth: vi.fn(),
      clearAuth: vi.fn(),
      setLoading: vi.fn(),
    });

    const { result } = renderHook(() => useAuth(), { wrapper });

    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.user).toEqual(mockUser);
    expect(result.current.partner).toEqual(mockPartner);
  });

  it('should expose loading states for mutations', () => {
    (useAuthStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      user: null,
      partner: null,
      accessToken: null,
      isAuthenticated: false,
      isLoading: false,
      setAuth: vi.fn(),
      clearAuth: vi.fn(),
      setLoading: vi.fn(),
    });

    const { result } = renderHook(() => useAuth(), { wrapper });

    expect(result.current.isLoggingIn).toBe(false);
    expect(result.current.isRegistering).toBe(false);
    expect(result.current.loginError).toBeNull();
    expect(result.current.registerError).toBeNull();
  });
});
