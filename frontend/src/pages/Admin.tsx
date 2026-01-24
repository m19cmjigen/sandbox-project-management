import { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  Tabs,
  Tab,
  Button,
  Card,
  CardContent,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Chip,
} from '@mui/material'
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Business as BusinessIcon,
} from '@mui/icons-material'
import { Organization, ProjectWithStats } from '@/types'
import { organizationService } from '@/services/organizationService'
import { projectService } from '@/services/projectService'
import Loading from '@/components/Loading'
import ErrorMessage from '@/components/ErrorMessage'
import OrganizationForm, {
  OrganizationFormData,
} from '@/components/OrganizationForm'
import ProjectAssignmentForm from '@/components/ProjectAssignmentForm'

export default function Admin() {
  const [activeTab, setActiveTab] = useState(0)
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [projects, setProjects] = useState<ProjectWithStats[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Form states
  const [orgFormOpen, setOrgFormOpen] = useState(false)
  const [editingOrg, setEditingOrg] = useState<Organization | null>(null)
  const [projectAssignOpen, setProjectAssignOpen] = useState(false)
  const [assigningProject, setAssigningProject] = useState<ProjectWithStats | null>(null)

  // Delete confirmation
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [deletingOrg, setDeletingOrg] = useState<Organization | null>(null)

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      setLoading(true)
      const [orgsData, projectsData] = await Promise.all([
        organizationService.getAll(),
        projectService.getAll(),
      ])
      setOrganizations(orgsData)
      setProjects(projectsData)
    } catch (err) {
      setError('データの取得に失敗しました')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateOrg = () => {
    setEditingOrg(null)
    setOrgFormOpen(true)
  }

  const handleEditOrg = (org: Organization) => {
    setEditingOrg(org)
    setOrgFormOpen(true)
  }

  const handleOrgFormSubmit = async (data: OrganizationFormData) => {
    if (editingOrg) {
      // Update
      await organizationService.update(editingOrg.id, data.name, data.parent_id)
    } else {
      // Create
      await organizationService.create(data.name, data.parent_id)
    }
    await fetchData()
  }

  const handleDeleteClick = (org: Organization) => {
    setDeletingOrg(org)
    setDeleteConfirmOpen(true)
  }

  const handleDeleteConfirm = async () => {
    if (!deletingOrg) return

    try {
      await organizationService.delete(deletingOrg.id)
      await fetchData()
      setDeleteConfirmOpen(false)
      setDeletingOrg(null)
    } catch (err) {
      setError('組織の削除に失敗しました')
      console.error(err)
    }
  }

  const handleAssignProject = (project: ProjectWithStats) => {
    setAssigningProject(project)
    setProjectAssignOpen(true)
  }

  const handleProjectAssignSubmit = async (
    projectId: number,
    organizationId: number | null
  ) => {
    await projectService.assignToOrganization(projectId, organizationId)
    await fetchData()
  }

  const getOrgName = (orgId: number | null) => {
    if (!orgId) return '未割り当て'
    const org = organizations.find((o) => o.id === orgId)
    return org?.name || '不明'
  }

  if (loading) return <Loading />
  if (error) return <ErrorMessage message={error} />

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        管理画面
      </Typography>

      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={activeTab} onChange={(_, v) => setActiveTab(v)}>
          <Tab label="組織管理" />
          <Tab label="プロジェクト割り当て" />
        </Tabs>
      </Box>

      {/* 組織管理タブ */}
      {activeTab === 0 && (
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
              <Typography variant="h6">組織一覧</Typography>
              <Button
                variant="contained"
                startIcon={<AddIcon />}
                onClick={handleCreateOrg}
              >
                新しい組織
              </Button>
            </Box>

            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>組織名</TableCell>
                    <TableCell>レベル</TableCell>
                    <TableCell>パス</TableCell>
                    <TableCell>親組織</TableCell>
                    <TableCell align="right">操作</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {organizations.map((org) => (
                    <TableRow key={org.id} hover>
                      <TableCell>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <BusinessIcon color="primary" />
                          <Typography fontWeight="medium">{org.name}</Typography>
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Chip label={`レベル ${org.level}`} size="small" />
                      </TableCell>
                      <TableCell>
                        <Typography variant="caption" color="textSecondary">
                          {org.path}
                        </Typography>
                      </TableCell>
                      <TableCell>{getOrgName(org.parent_id)}</TableCell>
                      <TableCell align="right">
                        <IconButton
                          size="small"
                          onClick={() => handleEditOrg(org)}
                          color="primary"
                        >
                          <EditIcon />
                        </IconButton>
                        <IconButton
                          size="small"
                          onClick={() => handleDeleteClick(org)}
                          color="error"
                        >
                          <DeleteIcon />
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </CardContent>
        </Card>
      )}

      {/* プロジェクト割り当てタブ */}
      {activeTab === 1 && (
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              プロジェクト一覧
            </Typography>
            <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
              各プロジェクトを組織に割り当てることができます
            </Typography>

            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>プロジェクト名</TableCell>
                    <TableCell>キー</TableCell>
                    <TableCell>現在の組織</TableCell>
                    <TableCell>チケット数</TableCell>
                    <TableCell align="right">操作</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {projects.map((project) => (
                    <TableRow key={project.id} hover>
                      <TableCell>
                        <Typography fontWeight="medium">{project.name}</Typography>
                      </TableCell>
                      <TableCell>
                        <Chip label={project.key} size="small" variant="outlined" />
                      </TableCell>
                      <TableCell>{getOrgName(project.organization_id)}</TableCell>
                      <TableCell>{project.total_issues}</TableCell>
                      <TableCell align="right">
                        <Button
                          size="small"
                          variant="outlined"
                          onClick={() => handleAssignProject(project)}
                        >
                          割り当て
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </CardContent>
        </Card>
      )}

      {/* 組織フォーム */}
      <OrganizationForm
        open={orgFormOpen}
        onClose={() => setOrgFormOpen(false)}
        onSubmit={handleOrgFormSubmit}
        organizations={organizations}
        editingOrganization={editingOrg}
      />

      {/* プロジェクト割り当てフォーム */}
      <ProjectAssignmentForm
        open={projectAssignOpen}
        onClose={() => setProjectAssignOpen(false)}
        onSubmit={handleProjectAssignSubmit}
        project={assigningProject}
        organizations={organizations}
      />

      {/* 削除確認ダイアログ */}
      <Dialog open={deleteConfirmOpen} onClose={() => setDeleteConfirmOpen(false)}>
        <DialogTitle>組織を削除</DialogTitle>
        <DialogContent>
          <DialogContentText>
            本当に「{deletingOrg?.name}」を削除しますか？
            <br />
            この操作は取り消せません。
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteConfirmOpen(false)}>キャンセル</Button>
          <Button onClick={handleDeleteConfirm} color="error" variant="contained">
            削除
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
