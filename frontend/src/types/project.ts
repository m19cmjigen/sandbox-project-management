export type DelayStatus = 'RED' | 'YELLOW' | 'GREEN'

export interface Project {
  id: number
  jira_project_id: string
  key: string
  name: string
  lead_account_id: string | null
  lead_email: string | null
  organization_id: number | null
  created_at: string
  updated_at: string
  red_count: number
  yellow_count: number
  green_count: number
  open_count: number
  total_count: number
  delay_status: DelayStatus
}

export interface PaginationMeta {
  page: number
  per_page: number
  total: number
  total_pages: number
}

export interface ProjectListResponse {
  data: Project[]
  pagination: PaginationMeta
}

export type SortOption = 'name' | 'name_desc' | 'delay_count'

export type DelayFilter = 'ALL' | DelayStatus
