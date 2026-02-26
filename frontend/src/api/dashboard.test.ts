import { describe, it, expect, vi, beforeEach } from 'vitest'

vi.mock('./apiClient', () => ({
  default: {
    get: vi.fn(),
  },
}))

import apiClient from './apiClient'
import { getDashboardSummary, getOrgSummary, getProjectSummary } from './dashboard'

const mockGet = vi.mocked(apiClient.get)

beforeEach(() => {
  mockGet.mockReset()
})

describe('getDashboardSummary', () => {
  it('returns summary data from API', async () => {
    const summary = {
      total_projects: 10,
      red_projects: 3,
      yellow_projects: 2,
      green_projects: 5,
      total_issues: 80,
      red_issues: 15,
      yellow_issues: 10,
      green_issues: 55,
      organizations: [],
    }
    mockGet.mockResolvedValueOnce({ data: summary })

    const result = await getDashboardSummary()

    expect(mockGet).toHaveBeenCalledWith('/dashboard/summary')
    expect(result.total_projects).toBe(10)
    expect(result.red_projects).toBe(3)
  })

  it('throws on API error', async () => {
    mockGet.mockRejectedValueOnce(new Error('Internal Server Error'))

    await expect(getDashboardSummary()).rejects.toThrow('Internal Server Error')
  })
})

describe('getOrgSummary', () => {
  it('fetches org summary by id', async () => {
    const orgSummary = {
      organization: { id: 1, name: '開発本部' },
      summary: { red_count: 2, yellow_count: 1, green_count: 3 },
    }
    mockGet.mockResolvedValueOnce({ data: orgSummary })

    const result = await getOrgSummary(1)

    expect(mockGet).toHaveBeenCalledWith('/dashboard/organizations/1')
    expect(result).toEqual(orgSummary)
  })
})

describe('getProjectSummary', () => {
  it('fetches project summary by id', async () => {
    const projectSummary = {
      project: { id: 1, name: 'プロジェクトA', delay_status: 'RED' },
      summary: { red_count: 5, yellow_count: 2, green_count: 3, open_count: 10 },
      delayed_issues: [],
    }
    mockGet.mockResolvedValueOnce({ data: projectSummary })

    const result = await getProjectSummary(1)

    expect(mockGet).toHaveBeenCalledWith('/dashboard/projects/1')
    expect(result).toEqual(projectSummary)
  })
})
