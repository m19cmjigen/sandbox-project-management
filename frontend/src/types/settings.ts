export interface JiraSettings {
  id: number
  jira_url: string
  email: string
  api_token_mask: string
  configured: boolean
  created_at: string
  updated_at: string
}

export interface SyncLog {
  id: number
  sync_type: string
  executed_at: string
  completed_at: string | null
  status: 'RUNNING' | 'SUCCESS' | 'FAILED'
  projects_synced: number
  issues_synced: number
  error_message: string | null
  duration_seconds: number | null
}
