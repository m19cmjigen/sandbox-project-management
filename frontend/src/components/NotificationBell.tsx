import { useState } from 'react'
import { Badge, IconButton, Popover, Tooltip } from '@mui/material'
import { Notifications as NotificationsIcon } from '@mui/icons-material'
import { useNotificationStore } from '../stores/notificationStore'
import NotificationPanel from './NotificationPanel'

export default function NotificationBell() {
  const unreadCount = useNotificationStore((s) => s.unreadCount)
  const [anchor, setAnchor] = useState<HTMLElement | null>(null)

  const handleOpen = (e: React.MouseEvent<HTMLElement>) => {
    setAnchor(e.currentTarget)
  }

  const handleClose = () => {
    setAnchor(null)
  }

  const open = Boolean(anchor)

  return (
    <>
      <Tooltip title="通知">
        <IconButton
          onClick={handleOpen}
          size="small"
          sx={{
            color: '#94a3b8',
            '&:hover': { color: '#f8fafc', bgcolor: 'rgba(148,163,184,0.1)' },
            p: 0.75,
          }}
        >
          <Badge badgeContent={unreadCount} color="error" max={99}>
            <NotificationsIcon fontSize="small" />
          </Badge>
        </IconButton>
      </Tooltip>

      <Popover
        open={open}
        anchorEl={anchor}
        onClose={handleClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
        transformOrigin={{ vertical: 'top', horizontal: 'left' }}
        slotProps={{
          paper: { elevation: 4, sx: { mt: 0.5 } },
        }}
      >
        <NotificationPanel />
      </Popover>
    </>
  )
}
