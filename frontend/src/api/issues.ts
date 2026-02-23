import apiClient from './apiClient'
import type { IssueListResponse, IssueListParams } from '../types/issue'

export const getIssues = async (params?: IssueListParams): Promise<IssueListResponse> => {
  const query: Record<string, string | number | boolean> = {}

  if (params?.page) query.page = params.page
  if (params?.per_page) query.per_page = params.per_page
  if (params?.sort) query.sort = params.sort
  if (params?.order) query.order = params.order
  if (params?.project_id) query.project_id = params.project_id
  if (params?.delay_status && params.delay_status !== 'ALL') query.delay_status = params.delay_status
  if (params?.no_due_date) query.no_due_date = true
  if (params?.status_category) query.status_category = params.status_category
  if (params?.assignee_name) query.assignee_name = params.assignee_name

  const response = await apiClient.get<IssueListResponse>('/issues', { params: query })
  return response.data
}
