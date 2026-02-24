import apiClient from './apiClient'
import type { JiraSettings } from '../types/settings'

export interface UpdateJiraSettingsRequest {
  jira_url: string
  email: string
  api_token: string
}

export const getJiraSettings = async (): Promise<JiraSettings> => {
  const res = await apiClient.get<JiraSettings>('/settings/jira')
  return res.data
}

export const updateJiraSettings = async (data: UpdateJiraSettingsRequest): Promise<void> => {
  await apiClient.put('/settings/jira', data)
}

export const testJiraConnection = async (data?: Partial<UpdateJiraSettingsRequest>): Promise<void> => {
  await apiClient.post('/settings/jira/test', data ?? {})
}

export const triggerSync = async (): Promise<{ log_id: number }> => {
  const res = await apiClient.post<{ log_id: number }>('/settings/jira/sync')
  return res.data
}
