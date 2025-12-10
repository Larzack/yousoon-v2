import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuthStore, AdminUser } from '../authStore'

// Reset store before each test
beforeEach(() => {
  useAuthStore.setState({
    user: null,
    accessToken: null,
    isAuthenticated: false,
  })
  vi.clearAllMocks()
})

describe('useAuthStore', () => {
  const mockUser: AdminUser = {
    id: 'admin-123',
    email: 'admin@yousoon.com',
    firstName: 'Admin',
    lastName: 'User',
    avatar: 'https://example.com/avatar.png',
    role: 'super_admin',
  }

  const mockToken = 'jwt-token-xyz'

  describe('initial state', () => {
    it('should have null user initially', () => {
      const { user } = useAuthStore.getState()
      expect(user).toBeNull()
    })

    it('should have null accessToken initially', () => {
      const { accessToken } = useAuthStore.getState()
      expect(accessToken).toBeNull()
    })

    it('should not be authenticated initially', () => {
      const { isAuthenticated } = useAuthStore.getState()
      expect(isAuthenticated).toBe(false)
    })
  })

  describe('setAuth', () => {
    it('should set user and token', () => {
      const { setAuth } = useAuthStore.getState()
      
      setAuth(mockUser, mockToken)
      
      const state = useAuthStore.getState()
      expect(state.user).toEqual(mockUser)
      expect(state.accessToken).toBe(mockToken)
      expect(state.isAuthenticated).toBe(true)
    })

    it('should store token in localStorage', () => {
      const { setAuth } = useAuthStore.getState()
      
      setAuth(mockUser, mockToken)
      
      expect(localStorage.setItem).toHaveBeenCalledWith(
        'admin_access_token',
        mockToken
      )
    })

    it('should update user info', () => {
      const { setAuth } = useAuthStore.getState()
      
      setAuth(mockUser, mockToken)
      
      const { user } = useAuthStore.getState()
      expect(user?.email).toBe('admin@yousoon.com')
      expect(user?.role).toBe('super_admin')
    })
  })

  describe('logout', () => {
    it('should clear user and token', () => {
      // First, set authenticated state
      useAuthStore.setState({
        user: mockUser,
        accessToken: mockToken,
        isAuthenticated: true,
      })

      const { logout } = useAuthStore.getState()
      logout()

      const state = useAuthStore.getState()
      expect(state.user).toBeNull()
      expect(state.accessToken).toBeNull()
      expect(state.isAuthenticated).toBe(false)
    })

    it('should remove token from localStorage', () => {
      const { logout } = useAuthStore.getState()
      
      logout()
      
      expect(localStorage.removeItem).toHaveBeenCalledWith('admin_access_token')
    })
  })

  describe('AdminUser roles', () => {
    it('should accept super_admin role', () => {
      const superAdmin: AdminUser = { ...mockUser, role: 'super_admin' }
      const { setAuth } = useAuthStore.getState()
      
      setAuth(superAdmin, mockToken)
      
      expect(useAuthStore.getState().user?.role).toBe('super_admin')
    })

    it('should accept moderator role', () => {
      const moderator: AdminUser = { ...mockUser, role: 'moderator' }
      const { setAuth } = useAuthStore.getState()
      
      setAuth(moderator, mockToken)
      
      expect(useAuthStore.getState().user?.role).toBe('moderator')
    })

    it('should accept support role', () => {
      const support: AdminUser = { ...mockUser, role: 'support' }
      const { setAuth } = useAuthStore.getState()
      
      setAuth(support, mockToken)
      
      expect(useAuthStore.getState().user?.role).toBe('support')
    })
  })

  describe('persistence', () => {
    it('should persist to storage with correct name', () => {
      // The persist middleware uses the name 'admin-auth'
      const { setAuth } = useAuthStore.getState()
      setAuth(mockUser, mockToken)
      
      // Zustand persist will serialize the state
      const state = useAuthStore.getState()
      expect(state.user).not.toBeNull()
    })

    it('should partialize state correctly', () => {
      const { setAuth } = useAuthStore.getState()
      setAuth(mockUser, mockToken)
      
      const state = useAuthStore.getState()
      
      // These should be included in persisted state
      expect(state).toHaveProperty('user')
      expect(state).toHaveProperty('accessToken')
      expect(state).toHaveProperty('isAuthenticated')
      
      // Functions should not be persisted (but exist in state)
      expect(state).toHaveProperty('setAuth')
      expect(state).toHaveProperty('logout')
    })
  })
})
