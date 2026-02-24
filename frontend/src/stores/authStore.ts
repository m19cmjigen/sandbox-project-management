import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type Role = 'admin' | 'project_manager' | 'viewer'

export interface AuthUser {
  id: number
  email: string
  role: Role
}

interface AuthState {
  token: string | null
  user: AuthUser | null
  login: (token: string, user: AuthUser) => void
  logout: () => void
  isAuthenticated: () => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,

      login: (token: string, user: AuthUser) => {
        set({ token, user })
      },

      logout: () => {
        set({ token: null, user: null })
        window.location.href = '/login'
      },

      isAuthenticated: () => {
        return get().token !== null && get().user !== null
      },
    }),
    {
      name: 'auth-storage',
    }
  )
)
