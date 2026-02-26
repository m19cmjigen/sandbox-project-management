import apiClient from './apiClient'

export interface Notification {
  id: number
  type: 'SYNC_COMPLETED' | 'SYNC_FAILED'
  title: string
  body: string
  is_read: boolean
  related_log_id: number | null
  created_at: string
}

export interface NotificationsResponse {
  data: Notification[]
  unread_count: number
}

export const getNotifications = async (): Promise<NotificationsResponse> => {
  const res = await apiClient.get<NotificationsResponse>('/notifications')
  return res.data
}

export const markAsRead = async (id: number): Promise<void> => {
  await apiClient.put(`/notifications/${id}/read`)
}

export const markAllAsRead = async (): Promise<void> => {
  await apiClient.put('/notifications/read-all')
}
