import apiClient from './apiClient'
import type { ProjectListResponse, SortOption, DelayFilter } from '../types/project'

export interface ProjectListParams {
  page?: number
  per_page?: number
  sort?: SortOption
  organization_id?: number
  delay_status?: DelayFilter
}

export const getProjects = async (params?: ProjectListParams): Promise<ProjectListResponse> => {
  const queryParams: Record<string, string | number> = {}

  if (params?.page) queryParams.page = params.page
  if (params?.per_page) queryParams.per_page = params.per_page
  if (params?.sort) queryParams.sort = params.sort
  if (params?.organization_id) queryParams.organization_id = params.organization_id
  if (params?.delay_status && params.delay_status !== 'ALL') {
    queryParams.delay_status = params.delay_status
  }

  const response = await apiClient.get<ProjectListResponse>('/projects', { params: queryParams })
  return response.data
}
