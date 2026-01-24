// 組織
export interface Organization {
  id: number
  name: string
  parent_id: number | null
  path: string
  level: number
  created_at: string
  updated_at: string
}

export interface OrganizationWithChildren extends Organization {
  children?: Organization[]
}

// プロジェクト
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
}

export interface ProjectWithStats extends Project {
  total_issues: number
  red_issues: number
  yellow_issues: number
  green_issues: number
  open_issues: number
  done_issues: number
}

// チケット
export type DelayStatus = 'RED' | 'YELLOW' | 'GREEN'
export type StatusCategory = 'To Do' | 'In Progress' | 'Done'

export interface Issue {
  id: number
  jira_issue_id: string
  jira_issue_key: string
  project_id: number
  summary: string
  status: string
  status_category: StatusCategory
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

export interface IssueFilter {
  project_id?: number
  delay_status?: DelayStatus
  status?: string
  assignee_id?: string
  due_date_from?: string
  due_date_to?: string
  limit?: number
  offset?: number
}

// ダッシュボード
export interface DashboardSummary {
  total_projects: number
  delayed_projects: number
  warning_projects: number
  normal_projects: number
  total_issues: number
  red_issues: number
  yellow_issues: number
  green_issues: number
  projects_by_status: ProjectWithStats[]
}

export interface OrganizationSummary {
  organization: Organization
  total_projects: number
  delayed_projects: number
  warning_projects: number
  projects: ProjectWithStats[]
}
