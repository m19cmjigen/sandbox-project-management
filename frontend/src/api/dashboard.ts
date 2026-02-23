import apiClient from './apiClient'
import type { DashboardSummary, OrgSummaryResponse } from '../types/dashboard'

export const getDashboardSummary = async (): Promise<DashboardSummary> => {
  const response = await apiClient.get<DashboardSummary>('/dashboard/summary')
  return response.data
}

export const getOrgSummary = async (orgId: number): Promise<OrgSummaryResponse> => {
  const response = await apiClient.get<OrgSummaryResponse>(`/dashboard/organizations/${orgId}`)
  return response.data
}
