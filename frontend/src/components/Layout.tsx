import { ReactNode, useState } from 'react'
import {
  Box,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
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
} from '@mui/icons-material'
import { useNavigate, useLocation } from 'react-router-dom'

const drawerWidth = 220

const menuItems = [
  { text: 'ダッシュボード', icon: <DashboardIcon fontSize="small" />, path: '/' },
  { text: '組織ツリー',     icon: <BusinessIcon fontSize="small" />,  path: '/organizations' },
  { text: '組織管理',      icon: <SettingsIcon fontSize="small" />,   path: '/organizations/manage' },
  { text: 'プロジェクト',  icon: <FolderIcon fontSize="small" />,     path: '/projects' },
  { text: 'チケット',      icon: <BugReportIcon fontSize="small" />,  path: '/issues' },
]

const SIDEBAR_BG   = '#0f172a'  // Slate-900
const SIDEBAR_ITEM_ACTIVE = '#1e293b' // Slate-800
const SIDEBAR_TEXT = '#94a3b8'  // Slate-400
const SIDEBAR_TEXT_ACTIVE = '#f8fafc' // Slate-50

interface LayoutProps {
  children: ReactNode
}

function SidebarContent() {
  const navigate = useNavigate()
  const location = useLocation()

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

      {/* Footer */}
      <Box sx={{ px: 2.5, py: 2, borderTop: '1px solid #1e293b' }}>
        <Typography variant="caption" sx={{ color: '#334155', fontSize: '0.6875rem' }}>
          v1.0.0
        </Typography>
      </Box>
    </Box>
  )
}

export default function Layout({ children }: LayoutProps) {
  const [mobileOpen, setMobileOpen] = useState(false)
  const theme = useTheme()

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
