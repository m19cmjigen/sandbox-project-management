import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock apiClient before importing auth so HTTP calls are intercepted.
vi.mock('./apiClient', () => ({
  default: {
    post: vi.fn(),
  },
}))

import apiClient from './apiClient'
import { login } from './auth'

const mockPost = vi.mocked(apiClient.post)

beforeEach(() => {
  mockPost.mockReset()
})

describe('login', () => {
  it('returns token and user on success', async () => {
    mockPost.mockResolvedValueOnce({
      data: {
        access_token: 'jwt-token-abc',
        token_type: 'Bearer',
        expires_in: 3600,
        user: { id: 1, email: 'admin@example.com', role: 'admin' },
      },
    })

    const result = await login({ email: 'admin@example.com', password: 'secret' })

    expect(mockPost).toHaveBeenCalledWith('/auth/login', {
      email: 'admin@example.com',
      password: 'secret',
    })
    expect(result.token).toBe('jwt-token-abc')
    expect(result.user).toEqual({ id: 1, email: 'admin@example.com', role: 'admin' })
  })

  it('throws when the API returns an error', async () => {
    const apiError = new Error('Request failed with status code 401')
    mockPost.mockRejectedValueOnce(apiError)

    await expect(login({ email: 'wrong@example.com', password: 'bad' })).rejects.toThrow(
      'Request failed with status code 401'
    )
  })
})
