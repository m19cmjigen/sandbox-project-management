import {
  Box,
  Button,
  Divider,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Typography,
  useTheme,
  alpha,
} from '@mui/material'
import {
  CheckCircle as CheckCircleIcon,
  Error as ErrorIcon,
} from '@mui/icons-material'
import { useNotificationStore } from '../stores/notificationStore'
import type { Notification } from '../api/notifications'

function formatRelativeTime(isoString: string): string {
  const diff = Date.now() - new Date(isoString).getTime()
  const minutes = Math.floor(diff / 60_000)
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}

interface NotificationItemProps {
  notification: Notification
  onRead: (id: number) => void
}

function NotificationItem({ notification, onRead }: NotificationItemProps) {
  const theme = useTheme()
  const isFailed = notification.type === 'SYNC_FAILED'

  return (
    <ListItem disablePadding>
      <ListItemButton
        onClick={() => {
          if (!notification.is_read) {
            onRead(notification.id)
          }
        }}
        sx={{
          py: 1,
          px: 1.5,
          gap: 1,
          alignItems: 'flex-start',
          bgcolor: notification.is_read
            ? 'transparent'
            : alpha(theme.palette.primary.main, 0.04),
          '&:hover': {
            bgcolor: notification.is_read
              ? alpha(theme.palette.action.hover, 0.04)
              : alpha(theme.palette.primary.main, 0.08),
          },
        }}
      >
        <Box sx={{ mt: 0.25, flexShrink: 0 }}>
          {isFailed ? (
            <ErrorIcon sx={{ fontSize: 18, color: 'error.main' }} />
          ) : (
            <CheckCircleIcon sx={{ fontSize: 18, color: 'success.main' }} />
          )}
        </Box>
        <ListItemText
          primary={
            <Typography
              variant="body2"
              sx={{
                fontWeight: notification.is_read ? 400 : 600,
                fontSize: '0.8125rem',
                lineHeight: 1.4,
              }}
            >
              {notification.title}
            </Typography>
          }
          secondary={
            <Box component="span" sx={{ display: 'block' }}>
              <Typography
                variant="caption"
                sx={{ color: 'text.secondary', display: 'block', lineHeight: 1.4 }}
              >
                {notification.body}
              </Typography>
              <Typography
                variant="caption"
                sx={{ color: 'text.disabled', mt: 0.25, display: 'block' }}
              >
                {formatRelativeTime(notification.created_at)}
              </Typography>
            </Box>
          }
          disableTypography
        />
      </ListItemButton>
    </ListItem>
  )
}

export default function NotificationPanel() {
  const notifications = useNotificationStore((s) => s.notifications)
  const unreadCount = useNotificationStore((s) => s.unreadCount)
  const markRead = useNotificationStore((s) => s.markRead)
  const markAllRead = useNotificationStore((s) => s.markAllRead)

  return (
    <Box sx={{ width: 340 }}>
      {/* Header */}
      <Box
        sx={{
          px: 2,
          py: 1.5,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
        }}
      >
        <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
          通知
        </Typography>
        {unreadCount > 0 && (
          <Button
            size="small"
            onClick={markAllRead}
            sx={{ fontSize: '0.75rem', py: 0.25 }}
          >
            すべて既読にする
          </Button>
        )}
      </Box>
      <Divider />

      {/* List */}
      {notifications.length === 0 ? (
        <Box sx={{ px: 2, py: 4, textAlign: 'center' }}>
          <Typography variant="body2" color="text.secondary">
            通知はありません
          </Typography>
        </Box>
      ) : (
        <List
          disablePadding
          sx={{ maxHeight: 400, overflowY: 'auto' }}
        >
          {notifications.map((n, idx) => (
            <Box key={n.id}>
              <NotificationItem notification={n} onRead={markRead} />
              {idx < notifications.length - 1 && (
                <Divider component="li" sx={{ mx: 1.5 }} />
              )}
            </Box>
          ))}
        </List>
      )}
    </Box>
  )
}
