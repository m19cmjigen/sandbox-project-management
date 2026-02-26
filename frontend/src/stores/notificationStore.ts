import { create } from 'zustand'
import {
  getNotifications,
  markAsRead,
  markAllAsRead,
  type Notification,
} from '../api/notifications'

const POLLING_INTERVAL_MS = 30_000

interface NotificationState {
  notifications: Notification[]
  unreadCount: number
  fetch: () => Promise<void>
  markRead: (id: number) => Promise<void>
  markAllRead: () => Promise<void>
  startPolling: () => () => void
}

export const useNotificationStore = create<NotificationState>((set) => ({
  notifications: [],
  unreadCount: 0,

  fetch: async () => {
    try {
      const res = await getNotifications()
      set({ notifications: res.data, unreadCount: res.unread_count })
    } catch {
      // ネットワークエラー時はサイレントに無視し、表示状態を維持する
    }
  },

  markRead: async (id: number) => {
    try {
      await markAsRead(id)
      set((state) => ({
        notifications: state.notifications.map((n) =>
          n.id === id ? { ...n, is_read: true } : n,
        ),
        unreadCount: Math.max(0, state.unreadCount - 1),
      }))
    } catch {
      // エラー時はサイレントに無視する
    }
  },

  markAllRead: async () => {
    try {
      await markAllAsRead()
      set((state) => ({
        notifications: state.notifications.map((n) => ({ ...n, is_read: true })),
        unreadCount: 0,
      }))
    } catch {
      // エラー時はサイレントに無視する
    }
  },

  startPolling: () => {
    const store = useNotificationStore.getState()
    // 初回即時フェッチ
    void store.fetch()

    const timer = setInterval(() => {
      void useNotificationStore.getState().fetch()
    }, POLLING_INTERVAL_MS)

    return () => clearInterval(timer)
  },
}))
