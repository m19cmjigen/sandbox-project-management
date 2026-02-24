import apiClient from './apiClient'
import type { SyncLog } from '../types/settings'

export const getSyncLogs = async (): Promise<SyncLog[]> => {
  const res = await apiClient.get<{ data: SyncLog[] }>('/sync-logs')
  return res.data.data
}
