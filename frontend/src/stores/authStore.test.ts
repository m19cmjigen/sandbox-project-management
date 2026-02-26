import { describe, it, expect, beforeEach, vi } from 'vitest'

// Mock persist middleware before importing the store so it does not touch localStorage.
// This isolates the store logic from the storage layer in unit tests.
vi.mock('zustand/middleware', () => ({
  persist: (config: unknown) => config,
}))

import { useAuthStore } from './authStore'

const mockUser = { id: 1, email: 'admin@example.com', role: 'admin' as const }

// Mock window.location to capture href assignments from logout()
let locationHref = ''
Object.defineProperty(window, 'location', {
  value: {
    get href() {
      return locationHref
    },
    set href(val: string) {
      locationHref = val
    },
  },
  configurable: true,
})

beforeEach(() => {
  locationHref = ''
  useAuthStore.setState({ token: null, user: null })
})

describe('login', () => {
  it('sets token and user', () => {
    const { login } = useAuthStore.getState()
    login('my-token', mockUser)

    const state = useAuthStore.getState()
    expect(state.token).toBe('my-token')
    expect(state.user).toEqual(mockUser)
  })
})

describe('logout', () => {
  it('clears token and user', () => {
    useAuthStore.setState({ token: 'old-token', user: mockUser })

    const { logout } = useAuthStore.getState()
    logout()

    const state = useAuthStore.getState()
    expect(state.token).toBeNull()
    expect(state.user).toBeNull()
  })

  it('redirects to /login', () => {
    useAuthStore.setState({ token: 'old-token', user: mockUser })

    const { logout } = useAuthStore.getState()
    logout()

    expect(locationHref).toBe('/login')
  })
})

describe('isAuthenticated', () => {
  it('returns false when token is null', () => {
    const { isAuthenticated } = useAuthStore.getState()
    expect(isAuthenticated()).toBe(false)
  })

  it('returns false when user is null but token exists', () => {
    useAuthStore.setState({ token: 'my-token', user: null })
    const { isAuthenticated } = useAuthStore.getState()
    expect(isAuthenticated()).toBe(false)
  })

  it('returns true when token and user are set', () => {
    useAuthStore.setState({ token: 'my-token', user: mockUser })
    const { isAuthenticated } = useAuthStore.getState()
    expect(isAuthenticated()).toBe(true)
  })
})
