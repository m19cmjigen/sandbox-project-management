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
import {
  getOrganizations,
  getOrganization,
  createOrganization,
  updateOrganization,
  deleteOrganization,
} from './organizations'

const mockGet = vi.mocked(apiClient.get)
const mockPost = vi.mocked(apiClient.post)
const mockPut = vi.mocked(apiClient.put)
const mockDelete = vi.mocked(apiClient.delete)

const mockOrg = {
  id: 1,
  name: '開発本部',
  parent_id: null,
  path: '/1/',
  level: 0,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

beforeEach(() => {
  mockGet.mockReset()
  mockPost.mockReset()
  mockPut.mockReset()
  mockDelete.mockReset()
})

describe('getOrganizations', () => {
  it('returns organization array', async () => {
    mockGet.mockResolvedValueOnce({ data: [mockOrg] })

    const result = await getOrganizations()

    expect(mockGet).toHaveBeenCalledWith('/organizations')
    expect(result).toEqual([mockOrg])
  })
})

describe('getOrganization', () => {
  it('fetches a single organization by id', async () => {
    mockGet.mockResolvedValueOnce({ data: mockOrg })

    const result = await getOrganization(1)

    expect(mockGet).toHaveBeenCalledWith('/organizations/1')
    expect(result.id).toBe(1)
  })

  it('throws on not-found error', async () => {
    mockGet.mockRejectedValueOnce(new Error('Request failed with status code 404'))

    await expect(getOrganization(999)).rejects.toThrow('404')
  })
})

describe('createOrganization', () => {
  it('posts new organization without parent', async () => {
    mockPost.mockResolvedValueOnce({ data: mockOrg })

    const result = await createOrganization({ name: '開発本部' })

    expect(mockPost).toHaveBeenCalledWith('/organizations', { name: '開発本部' })
    expect(result.name).toBe('開発本部')
  })

  it('posts new organization with parent_id', async () => {
    const child = { ...mockOrg, id: 2, name: '第一開発部', parent_id: 1 }
    mockPost.mockResolvedValueOnce({ data: child })

    const result = await createOrganization({ name: '第一開発部', parent_id: 1 })

    expect(mockPost).toHaveBeenCalledWith('/organizations', { name: '第一開発部', parent_id: 1 })
    expect(result.parent_id).toBe(1)
  })
})

describe('updateOrganization', () => {
  it('puts updated name', async () => {
    const updated = { ...mockOrg, name: '新開発本部' }
    mockPut.mockResolvedValueOnce({ data: updated })

    const result = await updateOrganization(1, { name: '新開発本部' })

    expect(mockPut).toHaveBeenCalledWith('/organizations/1', { name: '新開発本部' })
    expect(result.name).toBe('新開発本部')
  })
})

describe('deleteOrganization', () => {
  it('sends DELETE request for the given id', async () => {
    mockDelete.mockResolvedValueOnce({ data: undefined })

    await deleteOrganization(1)

    expect(mockDelete).toHaveBeenCalledWith('/organizations/1')
  })

  it('throws when delete is rejected', async () => {
    mockDelete.mockRejectedValueOnce(new Error('Request failed with status code 409'))

    await expect(deleteOrganization(1)).rejects.toThrow('409')
  })
})
