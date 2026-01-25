import apiClient from './api'

export interface SyncLog {
  id: number
  started_at: string
  completed_at: string | null
  status: 'RUNNING' | 'COMPLETED' | 'COMPLETED_WITH_ERRORS' | 'FAILED'
  projects_synced: number
  issues_synced: number
  error_count: number
  error_message: string | null
  created_at: string
  updated_at: string
}

export interface TriggerSyncRequest {
  organization_id: number
}

export interface TriggerSyncResponse {
  message: string
  sync_log: SyncLog
}

export const syncService = {
  // 同期をトリガー
  async triggerSync(organizationId: number): Promise<TriggerSyncResponse> {
    const response = await apiClient.post('/sync/trigger', {
      organization_id: organizationId,
    })
    return response.data
  },

  // プロジェクト別同期をトリガー
  async syncProject(projectId: number): Promise<void> {
    await apiClient.post(`/sync/projects/${projectId}`)
  },

  // 同期ログ一覧を取得（仮実装 - バックエンドAPIが実装されたら使用）
  async getSyncLogs(): Promise<SyncLog[]> {
    try {
      const response = await apiClient.get('/sync/logs')
      return response.data.logs || []
    } catch (err) {
      // APIが未実装の場合は空配列を返す
      console.warn('Sync logs API not implemented yet')
      return []
    }
  },

  // 同期ログ詳細を取得（仮実装）
  async getSyncLog(id: number): Promise<SyncLog | null> {
    try {
      const response = await apiClient.get(`/sync/logs/${id}`)
      return response.data
    } catch (err) {
      console.warn('Sync log detail API not implemented yet')
      return null
    }
  },
}
