import { describe, it, expect, vi, beforeEach } from 'vitest'

vi.mock('./apiClient', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

import apiClient from './apiClient'
import { getUsers, createUser, updateUser, deleteUser } from './users'

const mockGet = vi.mocked(apiClient.get)
const mockPost = vi.mocked(apiClient.post)
const mockPut = vi.mocked(apiClient.put)
const mockDelete = vi.mocked(apiClient.delete)

const mockUser = { id: 1, email: 'admin@example.com', role: 'admin' as const, is_active: true }

beforeEach(() => {
  mockGet.mockReset()
  mockPost.mockReset()
  mockPut.mockReset()
  mockDelete.mockReset()
})

describe('getUsers', () => {
  it('returns user list from API', async () => {
    mockGet.mockResolvedValueOnce({ data: { data: [mockUser] } })

    const result = await getUsers()

    expect(mockGet).toHaveBeenCalledWith('/users')
    expect(result).toEqual([mockUser])
  })

  it('throws on API error', async () => {
    mockGet.mockRejectedValueOnce(new Error('Unauthorized'))

    await expect(getUsers()).rejects.toThrow('Unauthorized')
  })
})

describe('createUser', () => {
  it('posts new user and returns created user', async () => {
    mockPost.mockResolvedValueOnce({ data: mockUser })

    const result = await createUser({
      email: 'admin@example.com',
      password: 'Password1!',
      role: 'admin',
    })

    expect(mockPost).toHaveBeenCalledWith('/users', {
      email: 'admin@example.com',
      password: 'Password1!',
      role: 'admin',
    })
    expect(result).toEqual(mockUser)
  })

  it('throws on conflict error', async () => {
    mockPost.mockRejectedValueOnce(new Error('Request failed with status code 409'))

    await expect(
      createUser({ email: 'dup@example.com', password: 'pass1234', role: 'viewer' }),
    ).rejects.toThrow('409')
  })
})

describe('updateUser', () => {
  it('puts role update and returns updated user', async () => {
    const updated = { ...mockUser, role: 'viewer' as const }
    mockPut.mockResolvedValueOnce({ data: updated })

    const result = await updateUser(1, { role: 'viewer' })

    expect(mockPut).toHaveBeenCalledWith('/users/1', { role: 'viewer' })
    expect(result.role).toBe('viewer')
  })

  it('puts is_active update', async () => {
    const updated = { ...mockUser, is_active: false }
    mockPut.mockResolvedValueOnce({ data: updated })

    const result = await updateUser(1, { is_active: false })

    expect(mockPut).toHaveBeenCalledWith('/users/1', { is_active: false })
    expect(result.is_active).toBe(false)
  })
})

describe('deleteUser', () => {
  it('sends DELETE request for the given user id', async () => {
    mockDelete.mockResolvedValueOnce({ data: undefined })

    await deleteUser(2)

    expect(mockDelete).toHaveBeenCalledWith('/users/2')
  })

  it('throws on forbidden error', async () => {
    mockDelete.mockRejectedValueOnce(new Error('Request failed with status code 403'))

    await expect(deleteUser(1)).rejects.toThrow('403')
  })
})
