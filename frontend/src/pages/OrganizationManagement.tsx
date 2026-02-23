import { useCallback, useEffect, useState } from 'react'
import {
  Alert,
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Divider,
  FormControl,
  IconButton,
  InputLabel,
  List,
  ListItem,
  ListItemText,
  MenuItem,
  Paper,
  Select,
  Stack,
  TextField,
  Tooltip,
  Typography,
} from '@mui/material'
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  AccountTree as TreeIcon,
  Assignment as AssignIcon,
} from '@mui/icons-material'
import type { SelectChangeEvent } from '@mui/material'
import { getOrganizations, createOrganization, updateOrganization, deleteOrganization } from '../api/organizations'
import { getUnassignedProjects, assignProjectOrganization } from '../api/projects'
import type { Organization } from '../types/organization'
import type { Project } from '../types/project'
import LoadingSpinner from '../components/LoadingSpinner'

// ---- Organization tree helpers ----

interface OrgNode extends Organization {
  children: OrgNode[]
}

function buildTree(orgs: Organization[]): OrgNode[] {
  const map = new Map<number, OrgNode>()
  orgs.forEach((o) => map.set(o.id, { ...o, children: [] }))
  const roots: OrgNode[] = []
  map.forEach((node) => {
    if (node.parent_id === null) roots.push(node)
    else map.get(node.parent_id)?.children.push(node)
  })
  return roots
}

// ---- Sub-components ----

interface OrgTreeItemProps {
  node: OrgNode
  allOrgs: Organization[]
  onEdit: (org: Organization) => void
  onDelete: (org: Organization) => void
  onAddChild: (parent: Organization) => void
}

function OrgTreeItem({ node, allOrgs, onEdit, onDelete, onAddChild }: OrgTreeItemProps) {
  const canAddChild = node.level < 1 // max 2 levels (0 and 1)
  return (
    <Box>
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 1,
          py: 0.75,
          px: 1,
          borderRadius: 1,
          '&:hover': { bgcolor: 'action.hover' },
        }}
      >
        <TreeIcon fontSize="small" color="action" sx={{ flexShrink: 0 }} />
        <Typography variant="body2" sx={{ flexGrow: 1 }}>
          {node.name}
        </Typography>
        <Chip label={`Lv.${node.level}`} size="small" variant="outlined" sx={{ fontSize: 10 }} />
        <Chip
          label={`${node.total_projects}件`}
          size="small"
          variant="outlined"
          color={node.delay_status === 'RED' ? 'error' : node.delay_status === 'YELLOW' ? 'warning' : 'success'}
        />
        {canAddChild && (
          <Tooltip title="子組織を追加">
            <IconButton size="small" onClick={() => onAddChild(node)}>
              <AddIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        )}
        <Tooltip title="編集">
          <IconButton size="small" onClick={() => onEdit(node)}>
            <EditIcon fontSize="small" />
          </IconButton>
        </Tooltip>
        <Tooltip title="削除">
          <IconButton size="small" color="error" onClick={() => onDelete(node)}>
            <DeleteIcon fontSize="small" />
          </IconButton>
        </Tooltip>
      </Box>
      {node.children.length > 0 && (
        <Box sx={{ pl: 3, borderLeft: '2px solid', borderColor: 'divider', ml: 1.5, mt: 0.5 }}>
          {node.children.map((child) => (
            <OrgTreeItem
              key={child.id}
              node={child}
              allOrgs={allOrgs}
              onEdit={onEdit}
              onDelete={onDelete}
              onAddChild={onAddChild}
            />
          ))}
        </Box>
      )}
    </Box>
  )
}

// ---- Main page ----

export default function OrganizationManagement() {
  const [orgs, setOrgs] = useState<Organization[]>([])
  const [unassignedProjects, setUnassignedProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [successMsg, setSuccessMsg] = useState<string | null>(null)

  // Org dialog
  const [orgDialogOpen, setOrgDialogOpen] = useState(false)
  const [editingOrg, setEditingOrg] = useState<Organization | null>(null)
  const [orgName, setOrgName] = useState('')
  const [orgParentId, setOrgParentId] = useState<number | undefined>(undefined)
  const [orgFormError, setOrgFormError] = useState('')
  const [orgSaving, setOrgSaving] = useState(false)

  // Delete dialog
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [deletingOrg, setDeletingOrg] = useState<Organization | null>(null)
  const [deleting, setDeleting] = useState(false)

  // Assign dialog
  const [assignDialogOpen, setAssignDialogOpen] = useState(false)
  const [assigningProject, setAssigningProject] = useState<Project | null>(null)
  const [selectedOrgId, setSelectedOrgId] = useState<string>('')
  const [assigning, setAssigning] = useState(false)

  const loadData = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const [orgData, projectData] = await Promise.all([
        getOrganizations(),
        getUnassignedProjects(),
      ])
      setOrgs(orgData)
      setUnassignedProjects(projectData.data)
    } catch {
      setError('データの取得に失敗しました。')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadData()
  }, [loadData])

  const showSuccess = (msg: string) => {
    setSuccessMsg(msg)
    setTimeout(() => setSuccessMsg(null), 3000)
  }

  // ---- Org CRUD handlers ----

  const openCreateDialog = (parent?: Organization) => {
    setEditingOrg(null)
    setOrgName('')
    setOrgParentId(parent?.id)
    setOrgFormError('')
    setOrgDialogOpen(true)
  }

  const openEditDialog = (org: Organization) => {
    setEditingOrg(org)
    setOrgName(org.name)
    setOrgParentId(undefined)
    setOrgFormError('')
    setOrgDialogOpen(true)
  }

  const handleOrgSave = async () => {
    if (!orgName.trim()) {
      setOrgFormError('組織名を入力してください')
      return
    }
    setOrgSaving(true)
    setOrgFormError('')
    try {
      if (editingOrg) {
        await updateOrganization(editingOrg.id, { name: orgName.trim() })
        showSuccess(`「${orgName.trim()}」を更新しました`)
      } else {
        await createOrganization({ name: orgName.trim(), parent_id: orgParentId })
        showSuccess(`「${orgName.trim()}」を作成しました`)
      }
      setOrgDialogOpen(false)
      loadData()
    } catch (e: unknown) {
      const msg = (e as { response?: { data?: { error?: string } } })?.response?.data?.error
      setOrgFormError(msg ?? '保存に失敗しました')
    } finally {
      setOrgSaving(false)
    }
  }

  const openDeleteDialog = (org: Organization) => {
    setDeletingOrg(org)
    setDeleteDialogOpen(true)
  }

  const handleDelete = async () => {
    if (!deletingOrg) return
    setDeleting(true)
    try {
      await deleteOrganization(deletingOrg.id)
      showSuccess(`「${deletingOrg.name}」を削除しました`)
      setDeleteDialogOpen(false)
      setDeletingOrg(null)
      loadData()
    } catch (e: unknown) {
      const msg = (e as { response?: { data?: { error?: string } } })?.response?.data?.error
      setError(msg ?? '削除に失敗しました')
      setDeleteDialogOpen(false)
    } finally {
      setDeleting(false)
    }
  }

  // ---- Project assignment handlers ----

  const openAssignDialog = (project: Project) => {
    setAssigningProject(project)
    setSelectedOrgId('')
    setAssignDialogOpen(true)
  }

  const handleAssign = async () => {
    if (!assigningProject || !selectedOrgId) return
    setAssigning(true)
    try {
      await assignProjectOrganization(assigningProject.id, Number(selectedOrgId))
      showSuccess(`「${assigningProject.name}」を組織に紐付けました`)
      setAssignDialogOpen(false)
      setAssigningProject(null)
      loadData()
    } catch {
      setError('プロジェクトの紐付けに失敗しました')
      setAssignDialogOpen(false)
    } finally {
      setAssigning(false)
    }
  }

  const roots = buildTree(orgs)
  // Flat list of orgs that can be parents (level 0 only for now)
  const parentCandidates = orgs.filter((o) => o.level === 0)

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 3 }}>組織管理</Typography>

      {successMsg && <Alert severity="success" sx={{ mb: 2 }}>{successMsg}</Alert>}
      {error && <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>{error}</Alert>}

      {loading ? (
        <LoadingSpinner minHeight={320} />
      ) : (
        <Stack direction={{ xs: 'column', md: 'row' }} spacing={3} alignItems="flex-start">
          {/* Left: Organization tree */}
          <Paper variant="outlined" sx={{ flex: 1, p: 2, minWidth: 0 }}>
            <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 2 }}>
              <Typography variant="h6">組織ツリー</Typography>
              <Button
                variant="contained"
                size="small"
                startIcon={<AddIcon />}
                onClick={() => openCreateDialog()}
              >
                本部を追加
              </Button>
            </Stack>
            <Divider sx={{ mb: 1.5 }} />
            {roots.length === 0 ? (
              <Typography variant="body2" color="text.secondary" sx={{ py: 2, textAlign: 'center' }}>
                組織が登録されていません
              </Typography>
            ) : (
              roots.map((root) => (
                <OrgTreeItem
                  key={root.id}
                  node={root}
                  allOrgs={orgs}
                  onEdit={openEditDialog}
                  onDelete={openDeleteDialog}
                  onAddChild={(parent) => openCreateDialog(parent)}
                />
              ))
            )}
          </Paper>

          {/* Right: Unclassified projects */}
          <Paper variant="outlined" sx={{ flex: 1, p: 2, minWidth: 0 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              未分類プロジェクト
              {unassignedProjects.length > 0 && (
                <Chip label={unassignedProjects.length} size="small" color="warning" sx={{ ml: 1 }} />
              )}
            </Typography>
            <Divider sx={{ mb: 1.5 }} />
            {unassignedProjects.length === 0 ? (
              <Typography variant="body2" color="text.secondary" sx={{ py: 2, textAlign: 'center' }}>
                未分類のプロジェクトはありません
              </Typography>
            ) : (
              <List dense disablePadding>
                {unassignedProjects.map((project) => (
                  <ListItem
                    key={project.id}
                    divider
                    secondaryAction={
                      <Tooltip title="組織に紐付け">
                        <IconButton
                          size="small"
                          color="primary"
                          onClick={() => openAssignDialog(project)}
                          disabled={orgs.length === 0}
                        >
                          <AssignIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    }
                  >
                    <ListItemText
                      primary={
                        <Stack direction="row" spacing={1} alignItems="center">
                          <Typography variant="body2" sx={{ fontFamily: 'monospace', fontWeight: 'bold' }}>
                            {project.key}
                          </Typography>
                          <Typography variant="body2">{project.name}</Typography>
                        </Stack>
                      }
                      secondary={project.lead_email ?? undefined}
                    />
                  </ListItem>
                ))}
              </List>
            )}
          </Paper>
        </Stack>
      )}

      {/* Create / Edit org dialog */}
      <Dialog open={orgDialogOpen} onClose={() => setOrgDialogOpen(false)} maxWidth="xs" fullWidth>
        <DialogTitle>{editingOrg ? '組織を編集' : '組織を追加'}</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ pt: 1 }}>
            <TextField
              label="組織名"
              value={orgName}
              onChange={(e) => setOrgName(e.target.value)}
              error={!!orgFormError}
              helperText={orgFormError || undefined}
              autoFocus
              fullWidth
              size="small"
              onKeyDown={(e) => { if (e.key === 'Enter') handleOrgSave() }}
            />
            {!editingOrg && (
              <FormControl size="small" fullWidth>
                <InputLabel>親組織（省略可）</InputLabel>
                <Select
                  value={orgParentId !== undefined ? String(orgParentId) : ''}
                  label="親組織（省略可）"
                  onChange={(e: SelectChangeEvent) =>
                    setOrgParentId(e.target.value ? Number(e.target.value) : undefined)
                  }
                >
                  <MenuItem value="">なし（本部として追加）</MenuItem>
                  {parentCandidates.map((o) => (
                    <MenuItem key={o.id} value={o.id}>{o.name}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            )}
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOrgDialogOpen(false)}>キャンセル</Button>
          <Button onClick={handleOrgSave} variant="contained" disabled={orgSaving}>
            {orgSaving ? '保存中...' : '保存'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete confirm dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)} maxWidth="xs" fullWidth>
        <DialogTitle>組織の削除</DialogTitle>
        <DialogContent>
          <DialogContentText>
            「{deletingOrg?.name}」を削除しますか？この操作は取り消せません。
            子組織やプロジェクトが紐付いている場合は削除できません。
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>キャンセル</Button>
          <Button onClick={handleDelete} color="error" variant="contained" disabled={deleting}>
            {deleting ? '削除中...' : '削除'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Assign project dialog */}
      <Dialog open={assignDialogOpen} onClose={() => setAssignDialogOpen(false)} maxWidth="xs" fullWidth>
        <DialogTitle>組織に紐付け</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ pt: 1 }}>
            <Typography variant="body2" color="text.secondary">
              プロジェクト: <strong>{assigningProject?.key} - {assigningProject?.name}</strong>
            </Typography>
            <FormControl size="small" fullWidth>
              <InputLabel>紐付け先組織</InputLabel>
              <Select
                value={selectedOrgId}
                label="紐付け先組織"
                onChange={(e: SelectChangeEvent) => setSelectedOrgId(e.target.value)}
              >
                <MenuItem value="" disabled>組織を選択してください</MenuItem>
                {orgs.map((o) => (
                  <MenuItem key={o.id} value={o.id}>
                    {'　'.repeat(o.level)}{o.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAssignDialogOpen(false)}>キャンセル</Button>
          <Button
            onClick={handleAssign}
            variant="contained"
            disabled={assigning || !selectedOrgId}
          >
            {assigning ? '紐付け中...' : '紐付け'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
