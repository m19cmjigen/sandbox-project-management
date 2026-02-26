import { ReactNode, useEffect, useState } from 'react'
import {
  Box,
  Chip,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Tooltip,
  Typography,
  useTheme,
  alpha,
} from '@mui/material'
import {
  Menu as MenuIcon,
  Dashboard as DashboardIcon,
  Business as BusinessIcon,
  Folder as FolderIcon,
  BugReport as BugReportIcon,
  Settings as SettingsIcon,
  BarChart as BarChartIcon,
  People as PeopleIcon,
  Logout as LogoutIcon,
  Tune as TuneIcon,
} from '@mui/icons-material'
import { useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import { useNotificationStore } from '../stores/notificationStore'
import { canManageUsers, canAccessSettings } from '../utils/permissions'
import NotificationBell from './NotificationBell'

const drawerWidth = 220

const baseMenuItems = [
  { text: 'ダッシュボード', icon: <DashboardIcon fontSize="small" />, path: '/' },
  { text: '組織ツリー',     icon: <BusinessIcon fontSize="small" />,  path: '/organizations' },
  { text: '組織管理',      icon: <SettingsIcon fontSize="small" />,   path: '/organizations/manage' },
  { text: 'プロジェクト',  icon: <FolderIcon fontSize="small" />,     path: '/projects' },
  { text: 'チケット',      icon: <BugReportIcon fontSize="small" />,  path: '/issues' },
]

const SIDEBAR_BG          = '#0f172a'  // Slate-900
const SIDEBAR_ITEM_ACTIVE = '#1e293b'  // Slate-800
const SIDEBAR_TEXT        = '#94a3b8'  // Slate-400
const SIDEBAR_TEXT_ACTIVE = '#f8fafc'  // Slate-50

const roleLabel: Record<string, string> = {
  admin:           '管理者',
  project_manager: 'PJ管理者',
  viewer:          '閲覧者',
}

const roleColor: Record<string, 'error' | 'warning' | 'default'> = {
  admin:           'error',
  project_manager: 'warning',
  viewer:          'default',
}

interface LayoutProps {
  children: ReactNode
}

function SidebarContent() {
  const navigate   = useNavigate()
  const location   = useLocation()
  const user       = useAuthStore((s) => s.user)
  const logout     = useAuthStore((s) => s.logout)

  const menuItems = [
    ...baseMenuItems,
    // ユーザー管理は admin のみ表示
    ...(user && canManageUsers(user.role)
      ? [{ text: 'ユーザー管理', icon: <PeopleIcon fontSize="small" />, path: '/users' }]
      : []),
    // 設定は admin / project_manager のみ表示
    ...(user && canAccessSettings(user.role)
      ? [{ text: '設定', icon: <TuneIcon fontSize="small" />, path: '/settings' }]
      : []),
  ]

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        height: '100%',
        bgcolor: SIDEBAR_BG,
      }}
    >
      {/* Logo / App title */}
      <Box sx={{ px: 2.5, py: 3, display: 'flex', alignItems: 'center', gap: 1.5 }}>
        <Box
          sx={{
            width: 32,
            height: 32,
            borderRadius: 2,
            bgcolor: '#6366f1',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            flexShrink: 0,
          }}
        >
          <BarChartIcon sx={{ fontSize: 18, color: '#fff' }} />
        </Box>
        <Box>
          <Typography
            variant="subtitle2"
            sx={{ color: '#f8fafc', fontWeight: 700, lineHeight: 1.2, fontSize: '0.8125rem' }}
          >
            ProjectViz
          </Typography>
          <Typography variant="caption" sx={{ color: SIDEBAR_TEXT, fontSize: '0.6875rem' }}>
            進捗可視化
          </Typography>
        </Box>
      </Box>

      {/* Nav items */}
      <Box sx={{ px: 1.5, flex: 1 }}>
        <Typography
          variant="caption"
          sx={{
            color: '#475569',
            fontWeight: 600,
            letterSpacing: '0.07em',
            textTransform: 'uppercase',
            px: 1,
            mb: 0.5,
            display: 'block',
            fontSize: '0.6875rem',
          }}
        >
          メニュー
        </Typography>
        <List disablePadding>
          {menuItems.map((item) => {
            const isActive = location.pathname === item.path
            return (
              <ListItem key={item.text} disablePadding sx={{ mb: 0.25 }}>
                <ListItemButton
                  selected={isActive}
                  onClick={() => navigate(item.path)}
                  sx={{
                    borderRadius: 1.5,
                    py: 0.875,
                    px: 1.25,
                    color: isActive ? SIDEBAR_TEXT_ACTIVE : SIDEBAR_TEXT,
                    bgcolor: isActive ? SIDEBAR_ITEM_ACTIVE : 'transparent',
                    '&:hover': {
                      bgcolor: isActive ? SIDEBAR_ITEM_ACTIVE : alpha('#94a3b8', 0.1),
                      color: SIDEBAR_TEXT_ACTIVE,
                    },
                    '&.Mui-selected': {
                      bgcolor: SIDEBAR_ITEM_ACTIVE,
                      '&:hover': { bgcolor: SIDEBAR_ITEM_ACTIVE },
                    },
                    transition: 'all 0.15s',
                  }}
                >
                  <ListItemIcon
                    sx={{
                      minWidth: 32,
                      color: isActive ? '#6366f1' : SIDEBAR_TEXT,
                    }}
                  >
                    {item.icon}
                  </ListItemIcon>
                  <ListItemText
                    primary={item.text}
                    primaryTypographyProps={{
                      fontSize: '0.8125rem',
                      fontWeight: isActive ? 600 : 400,
                      lineHeight: 1.4,
                    }}
                  />
                  {isActive && (
                    <Box
                      sx={{
                        width: 3,
                        height: 16,
                        borderRadius: 1,
                        bgcolor: '#6366f1',
                        flexShrink: 0,
                      }}
                    />
                  )}
                </ListItemButton>
              </ListItem>
            )
          })}
        </List>
      </Box>

      {/* Notification bell */}
      <Box sx={{ px: 1, py: 0.5, borderTop: '1px solid #1e293b' }}>
        <NotificationBell />
      </Box>

      {/* User info + Logout footer */}
      <Box sx={{ px: 2, py: 1.5, borderTop: '1px solid #1e293b' }}>
        {user && (
          <Box sx={{ mb: 1 }}>
            <Typography
              sx={{
                color: SIDEBAR_TEXT_ACTIVE,
                fontSize: '0.75rem',
                fontWeight: 500,
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
              }}
            >
              {user.email}
            </Typography>
            <Chip
              label={roleLabel[user.role] ?? user.role}
              size="small"
              color={roleColor[user.role] ?? 'default'}
              sx={{ mt: 0.5, fontSize: '0.625rem', height: 18 }}
            />
          </Box>
        )}
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Typography variant="caption" sx={{ color: '#334155', fontSize: '0.6875rem' }}>
            v1.0.0
          </Typography>
          <Tooltip title="ログアウト">
            <IconButton
              size="small"
              onClick={logout}
              sx={{
                color: SIDEBAR_TEXT,
                '&:hover': { color: '#f87171', bgcolor: alpha('#f87171', 0.1) },
                p: 0.5,
              }}
            >
              <LogoutIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>
    </Box>
  )
}

export default function Layout({ children }: LayoutProps) {
  const [mobileOpen, setMobileOpen] = useState(false)
  const theme = useTheme()
  const startPolling = useNotificationStore((s) => s.startPolling)

  useEffect(() => {
    const stopPolling = startPolling()
    return stopPolling
  }, [startPolling])

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh', bgcolor: 'background.default' }}>
      {/* Mobile hamburger */}
      <Box
        sx={{
          position: 'fixed',
          top: 12,
          left: 12,
          zIndex: theme.zIndex.drawer + 2,
          display: { sm: 'none' },
        }}
      >
        <IconButton
          onClick={() => setMobileOpen(!mobileOpen)}
          sx={{
            bgcolor: SIDEBAR_BG,
            color: '#f8fafc',
            '&:hover': { bgcolor: SIDEBAR_ITEM_ACTIVE },
            borderRadius: 2,
            p: 1,
          }}
          size="small"
        >
          <MenuIcon fontSize="small" />
        </IconButton>
      </Box>

      {/* Mobile drawer */}
      <Drawer
        variant="temporary"
        open={mobileOpen}
        onClose={() => setMobileOpen(false)}
        ModalProps={{ keepMounted: true }}
        sx={{
          display: { xs: 'block', sm: 'none' },
          '& .MuiDrawer-paper': {
            width: drawerWidth,
            border: 'none',
            bgcolor: SIDEBAR_BG,
          },
        }}
      >
        <SidebarContent />
      </Drawer>

      {/* Desktop permanent sidebar */}
      <Drawer
        variant="permanent"
        sx={{
          display: { xs: 'none', sm: 'block' },
          width: drawerWidth,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: drawerWidth,
            boxSizing: 'border-box',
            border: 'none',
            bgcolor: SIDEBAR_BG,
          },
        }}
        open
      >
        <SidebarContent />
      </Drawer>

      {/* Main content */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          minWidth: 0,
          p: { xs: 2, sm: 3, md: 4 },
          pt: { xs: 7, sm: 4 },
        }}
      >
        {children}
      </Box>
    </Box>
  )
}
