import { useEffect } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { Box, Tab, Tabs, Typography } from '@mui/material'
import {
  IntegrationInstructions as JiraIcon,
  FolderOpen as ProjectIcon,
  Business as OrgIcon,
} from '@mui/icons-material'
import { useAuthStore } from '../stores/authStore'
import { canAccessSettings } from '../utils/permissions'
import JiraSettingsTab from '../components/settings/JiraSettingsTab'
import ProjectSettingsTab from '../components/settings/ProjectSettingsTab'
import OrganizationSettingsTab from '../components/settings/OrganizationSettingsTab'

type TabId = 'jira' | 'projects' | 'organizations'

const allTabs: { id: TabId; label: string; icon: React.ReactElement; adminOnly?: boolean }[] = [
  { id: 'jira', label: 'Jira連携', icon: <JiraIcon />, adminOnly: true },
  { id: 'projects', label: 'PJ管理', icon: <ProjectIcon /> },
  { id: 'organizations', label: '組織管理', icon: <OrgIcon /> },
]

export default function Settings() {
  const [searchParams, setSearchParams] = useSearchParams()
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const isAdmin = user?.role === 'admin'

  // 設定ページへのアクセス権がないユーザーをリダイレクト
  useEffect(() => {
    if (user && !canAccessSettings(user.role)) {
      navigate('/', { replace: true })
    }
  }, [user, navigate])

  const currentTab = (searchParams.get('tab') ?? (isAdmin ? 'jira' : 'projects')) as TabId

  // project_manager が admin 専用タブにアクセスした場合はリダイレクト
  useEffect(() => {
    const tab = allTabs.find((t) => t.id === currentTab)
    if (tab?.adminOnly && !isAdmin) {
      setSearchParams({ tab: 'projects' }, { replace: true })
    }
  }, [isAdmin, currentTab, setSearchParams])

  const visibleTabs = allTabs.filter((t) => !t.adminOnly || isAdmin)
  const activeTab = visibleTabs.some((t) => t.id === currentTab)
    ? currentTab
    : (visibleTabs[0]?.id ?? 'projects')

  const handleTabChange = (_: React.SyntheticEvent, value: TabId) => {
    setSearchParams({ tab: value })
  }

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 3 }}>設定</Typography>

      <Tabs
        value={activeTab}
        onChange={handleTabChange}
        sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}
      >
        {visibleTabs.map((tab) => (
          <Tab
            key={tab.id}
            value={tab.id}
            label={tab.label}
            icon={tab.icon}
            iconPosition="start"
          />
        ))}
      </Tabs>

      {activeTab === 'jira' && isAdmin && <JiraSettingsTab />}
      {activeTab === 'projects' && <ProjectSettingsTab />}
      {activeTab === 'organizations' && <OrganizationSettingsTab />}
    </Box>
  )
}
