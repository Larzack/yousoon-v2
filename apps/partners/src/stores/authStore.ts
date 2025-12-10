import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  avatar?: string
  partnerId: string
  role: 'admin' | 'manager' | 'staff' | 'viewer'
}

export interface Partner {
  id: string
  name: string
  tradeName?: string
  logo?: string
  status: 'pending' | 'active' | 'suspended'
}

interface AuthState {
  user: User | null
  partner: Partner | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean
  
  // Actions
  setAuth: (data: {
    user: User
    partner: Partner
    accessToken: string
    refreshToken: string
  }) => void
  updateUser: (user: Partial<User>) => void
  updatePartner: (partner: Partial<Partner>) => void
  logout: () => void
  setLoading: (loading: boolean) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      partner: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,

      setAuth: (data) => {
        localStorage.setItem('access_token', data.accessToken)
        localStorage.setItem('refresh_token', data.refreshToken)
        set({
          user: data.user,
          partner: data.partner,
          accessToken: data.accessToken,
          refreshToken: data.refreshToken,
          isAuthenticated: true,
          isLoading: false,
        })
      },

      updateUser: (userData) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...userData } : null,
        })),

      updatePartner: (partnerData) =>
        set((state) => ({
          partner: state.partner ? { ...state.partner, ...partnerData } : null,
        })),

      logout: () => {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        set({
          user: null,
          partner: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
        })
      },

      setLoading: (loading) => set({ isLoading: loading }),
    }),
    {
      name: 'yousoon-auth',
      partialize: (state) => ({
        user: state.user,
        partner: state.partner,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
