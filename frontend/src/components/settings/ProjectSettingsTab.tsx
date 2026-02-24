import { useCallback, useEffect, useState } from 'react'
import {
  Alert,
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  FormControl,
  IconButton,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Stack,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Tooltip,
  Typography,
} from '@mui/material'
import { Assignment as AssignIcon, Search as SearchIcon } from '@mui/icons-material'
import type { SelectChangeEvent } from '@mui/material'
import { getAllProjectsForAdmin, updateProject, assignProjectOrganization } from '../../api/projects'
import { getOrganizations } from '../../api/organizations'
import type { Project } from '../../types/project'
import type { Organization } from '../../types/organization'
import LoadingSpinner from '../LoadingSpinner'
import { useAuthStore } from '../../stores/authStore'
import { canAssignProjects } from '../../utils/permissions'

const delayColor: Record<string, 'error' | 'warning' | 'success'> = {
  RED: 'error',
  YELLOW: 'warning',
  GREEN: 'success',
}

export default function ProjectSettingsTab() {
  const currentUser = useAuthStore((s) => s.user)
  const assignEnabled = currentUser ? canAssignProjects(currentUser.role) : false

  const [projects, setProjects] = useState<Project[]>([])
  const [orgs, setOrgs] = useState<Organization[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [successMsg, setSuccessMsg] = useState<string | null>(null)
  const [errorMsg, setErrorMsg] = useState<string | null>(null)

  // 組織割り当てダイアログ
  const [assignDialogOpen, setAssignDialogOpen] = useState(false)
  const [assigningProject, setAssigningProject] = useState<Project | null>(null)
  const [selectedOrgId, setSelectedOrgId] = useState<string>('')
  const [assigning, setAssigning] = useState(false)

  const showSuccess = (msg: string) => {
    setSuccessMsg(msg)
    setTimeout(() => setSuccessMsg(null), 3000)
  }

  const loadData = useCallback(async () => {
    setLoading(true)
    try {
      const [projData, orgData] = await Promise.all([
        getAllProjectsForAdmin(),
        getOrganizations(),
      ])
      setProjects(projData.data)
      setOrgs(orgData)
    } catch {
      setErrorMsg('データの取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadData()
  }, [loadData])

  const handleToggleActive = async (project: Project) => {
    try {
      await updateProject(project.id, { is_active: !project.is_active })
      setProjects((prev) =>
        prev.map((p) => (p.id === project.id ? { ...p, is_active: !p.is_active } : p))
      )
      showSuccess(`「${project.name}」の表示状態を更新しました`)
    } catch {
      setErrorMsg('更新に失敗しました')
    }
  }

  const openAssignDialog = (project: Project) => {
    setAssigningProject(project)
    setSelectedOrgId(project.organization_id ? String(project.organization_id) : '')
    setAssignDialogOpen(true)
  }

  const handleAssign = async () => {
    if (!assigningProject) return
    setAssigning(true)
    try {
      const orgId = selectedOrgId ? Number(selectedOrgId) : null
      await assignProjectOrganization(assigningProject.id, orgId)
      showSuccess(`「${assigningProject.name}」の組織を更新しました`)
      setAssignDialogOpen(false)
      loadData()
    } catch {
      setErrorMsg('組織の割り当てに失敗しました')
      setAssignDialogOpen(false)
    } finally {
      setAssigning(false)
    }
  }

  const getOrgName = (id: number | null) =>
    id ? (orgs.find((o) => o.id === id)?.name ?? '-') : '未分類'

  const filtered = projects.filter(
    (p) =>
      p.name.toLowerCase().includes(search.toLowerCase()) ||
      p.key.toLowerCase().includes(search.toLowerCase())
  )

  if (loading) return <LoadingSpinner minHeight={320} />

  return (
    <Box>
      {successMsg && <Alert severity="success" sx={{ mb: 2 }}>{successMsg}</Alert>}
      {errorMsg && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setErrorMsg(null)}>
          {errorMsg}
        </Alert>
      )}

      <Paper variant="outlined" sx={{ p: 2 }}>
        <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 2 }}>
          <Typography variant="h6">
            プロジェクト一覧
            <Typography component="span" variant="body2" color="text.secondary" sx={{ ml: 1 }}>
              ({filtered.length}件)
            </Typography>
          </Typography>
          <TextField
            size="small"
            placeholder="名前・キーで検索"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            InputProps={{
              startAdornment: <SearchIcon fontSize="small" sx={{ mr: 0.5, color: 'text.disabled' }} />,
            }}
            sx={{ width: 220 }}
          />
        </Stack>
        <Divider sx={{ mb: 1 }} />

        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>キー</TableCell>
                <TableCell>プロジェクト名</TableCell>
                <TableCell>組織</TableCell>
                <TableCell>遅延状況</TableCell>
                <TableCell align="center">表示</TableCell>
                {assignEnabled && <TableCell align="center">組織割り当て</TableCell>}
              </TableRow>
            </TableHead>
            <TableBody>
              {filtered.map((project) => (
                <TableRow
                  key={project.id}
                  sx={{ opacity: project.is_active ? 1 : 0.5 }}
                >
                  <TableCell>
                    <Typography variant="body2" sx={{ fontFamily: 'monospace', fontWeight: 'bold' }}>
                      {project.key}
                    </Typography>
                  </TableCell>
                  <TableCell>{project.name}</TableCell>
                  <TableCell>{getOrgName(project.organization_id)}</TableCell>
                  <TableCell>
                    <Chip
                      label={project.delay_status}
                      size="small"
                      color={delayColor[project.delay_status] ?? 'default'}
                    />
                  </TableCell>
                  <TableCell align="center">
                    <Switch
                      size="small"
                      checked={project.is_active}
                      onChange={() => handleToggleActive(project)}
                    />
                  </TableCell>
                  {assignEnabled && (
                    <TableCell align="center">
                      <Tooltip title="組織を割り当て">
                        <IconButton size="small" onClick={() => openAssignDialog(project)}>
                          <AssignIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    </TableCell>
                  )}
                </TableRow>
              ))}
              {filtered.length === 0 && (
                <TableRow>
                  <TableCell colSpan={assignEnabled ? 6 : 5} align="center">
                    <Typography variant="body2" color="text.secondary" sx={{ py: 2 }}>
                      プロジェクトが見つかりません
                    </Typography>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* 組織割り当てダイアログ */}
      {assignEnabled && (
        <Dialog open={assignDialogOpen} onClose={() => setAssignDialogOpen(false)} maxWidth="xs" fullWidth>
          <DialogTitle>組織を割り当て</DialogTitle>
          <DialogContent>
            <Stack spacing={2} sx={{ pt: 1 }}>
              <Typography variant="body2" color="text.secondary">
                プロジェクト: <strong>{assigningProject?.key} - {assigningProject?.name}</strong>
              </Typography>
              <FormControl size="small" fullWidth>
                <InputLabel>割り当て先組織</InputLabel>
                <Select
                  value={selectedOrgId}
                  label="割り当て先組織"
                  onChange={(e: SelectChangeEvent) => setSelectedOrgId(e.target.value)}
                >
                  <MenuItem value="">未分類（割り当て解除）</MenuItem>
                  {orgs.map((o) => (
                    <MenuItem key={o.id} value={String(o.id)}>
                      {'　'.repeat(o.level)}{o.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Stack>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setAssignDialogOpen(false)}>キャンセル</Button>
            <Button onClick={handleAssign} variant="contained" disabled={assigning}>
              {assigning ? '更新中...' : '更新'}
            </Button>
          </DialogActions>
        </Dialog>
      )}
    </Box>
  )
}
