import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface AdminUser {
  id: string
  email: string
  firstName: string
  lastName: string
  avatar?: string
  role: 'super_admin' | 'moderator' | 'support'
}

interface AuthState {
  user: AdminUser | null
  accessToken: string | null
  isAuthenticated: boolean
  
  setAuth: (user: AdminUser, token: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      isAuthenticated: false,

      setAuth: (user, token) => {
        localStorage.setItem('admin_access_token', token)
        set({
          user,
          accessToken: token,
          isAuthenticated: true,
        })
      },

      logout: () => {
        localStorage.removeItem('admin_access_token')
        set({
          user: null,
          accessToken: null,
          isAuthenticated: false,
        })
      },
    }),
    {
      name: 'admin-auth',
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
