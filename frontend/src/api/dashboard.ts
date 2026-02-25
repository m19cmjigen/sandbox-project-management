import apiClient from './apiClient'
import type { DashboardSummary, OrgSummaryResponse, ProjectSummaryResponse } from '../types/dashboard'

export const getDashboardSummary = async (): Promise<DashboardSummary> => {
  const response = await apiClient.get<DashboardSummary>('/dashboard/summary')
  return response.data
}

export const getOrgSummary = async (orgId: number): Promise<OrgSummaryResponse> => {
  const response = await apiClient.get<OrgSummaryResponse>(`/dashboard/organizations/${orgId}`)
  return response.data
}

export const getProjectSummary = async (projectId: number): Promise<ProjectSummaryResponse> => {
  const response = await apiClient.get<ProjectSummaryResponse>(`/dashboard/projects/${projectId}`)
  return response.data
}
