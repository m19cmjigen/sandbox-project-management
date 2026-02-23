import type { DelayStatus } from './project'
import type { PaginationMeta } from './project'

export interface Issue {
  id: number
  jira_issue_id: string
  jira_issue_key: string
  project_id: number
  project_key: string
  project_name: string
  summary: string
  status: string
  status_category: string
  due_date: string | null
  assignee_name: string | null
  assignee_account_id: string | null
  delay_status: DelayStatus
  priority: string | null
  issue_type: string | null
  last_updated_at: string
  created_at: string
  updated_at: string
}

export interface IssueListResponse {
  data: Issue[]
  pagination: PaginationMeta
}

export type IssueSortKey = 'due_date' | 'last_updated_at' | 'jira_issue_key' | 'delay_status'
export type SortOrder = 'asc' | 'desc'

export interface IssueListParams {
  page?: number
  per_page?: number
  sort?: IssueSortKey
  order?: SortOrder
  project_id?: number
  delay_status?: DelayStatus | 'ALL'
  no_due_date?: boolean
  status_category?: string
  assignee_name?: string
}
