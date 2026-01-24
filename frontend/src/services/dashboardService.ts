import apiClient from './api'
import { DashboardSummary, OrganizationSummary, ProjectWithStats } from '@/types'

export const dashboardService = {
  // 全社サマリを取得
  async getSummary(): Promise<DashboardSummary> {
    const response = await apiClient.get('/dashboard/summary')
    return response.data
  },

  // 組織別サマリを取得
  async getOrganizationSummary(organizationId: number): Promise<OrganizationSummary> {
    const response = await apiClient.get(`/dashboard/organizations/${organizationId}`)
    return response.data
  },

  // プロジェクト別サマリを取得
  async getProjectSummary(projectId: number): Promise<ProjectWithStats> {
    const response = await apiClient.get(`/dashboard/projects/${projectId}`)
    return response.data
  },
}
