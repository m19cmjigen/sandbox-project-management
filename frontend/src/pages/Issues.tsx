import { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  Grid,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Card,
  CardContent,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
} from '@mui/material'
import { Issue, Organization, ProjectWithStats, DelayStatus } from '@/types'
import { issueService } from '@/services/issueService'
import { organizationService } from '@/services/organizationService'
import { projectService } from '@/services/projectService'
import Loading from '@/components/Loading'
import ErrorMessage from '@/components/ErrorMessage'
import StatusBadge from '@/components/StatusBadge'

export default function Issues() {
  const [issues, setIssues] = useState<Issue[]>([])
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [projects, setProjects] = useState<ProjectWithStats[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // フィルター状態
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedOrgId, setSelectedOrgId] = useState<number | ''>('')
  const [selectedProjectId, setSelectedProjectId] = useState<number | ''>('')
  const [selectedStatus, setSelectedStatus] = useState<DelayStatus | 'all'>('all')

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)
        const [issuesData, orgsData, projectsData] = await Promise.all([
          issueService.getAll(),
          organizationService.getAll(),
          projectService.getAll(),
        ])
        setIssues(issuesData)
        setOrganizations(orgsData)
        setProjects(projectsData)
      } catch (err) {
        setError('データの取得に失敗しました')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  const filteredIssues = issues.filter((issue) => {
    // 検索条件
    const matchesSearch =
      searchTerm === '' ||
      issue.summary.toLowerCase().includes(searchTerm.toLowerCase()) ||
      issue.jira_key.toLowerCase().includes(searchTerm.toLowerCase())

    // プロジェクトフィルター
    const matchesProject =
      selectedProjectId === '' || issue.project_id === selectedProjectId

    // 組織フィルター (プロジェクト経由)
    let matchesOrg = true
    if (selectedOrgId !== '') {
      const project = projects.find((p) => p.id === issue.project_id)
      matchesOrg = project?.organization_id === selectedOrgId
    }

    // ステータスフィルター
    const matchesStatus =
      selectedStatus === 'all' || issue.delay_status === selectedStatus

    return matchesSearch && matchesProject && matchesOrg && matchesStatus
  })

  const formatDate = (dateString: string | null) => {
    if (!dateString) return '-'
    return new Date(dateString).toLocaleDateString('ja-JP')
  }

  const getProjectName = (projectId: number) => {
    const project = projects.find((p) => p.id === projectId)
    return project?.name || '-'
  }

  if (loading) return <Loading />
  if (error) return <ErrorMessage message={error} />

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        チケット一覧
      </Typography>

      {/* フィルター */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2}>
            <Grid item xs={12} md={3}>
              <TextField
                fullWidth
                label="検索"
                placeholder="チケット名またはキー"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </Grid>
            <Grid item xs={12} md={3}>
              <FormControl fullWidth>
                <InputLabel>組織</InputLabel>
                <Select
                  value={selectedOrgId}
                  label="組織"
                  onChange={(e) => setSelectedOrgId(e.target.value as number | '')}
                >
                  <MenuItem value="">すべて</MenuItem>
                  {organizations.map((org) => (
                    <MenuItem key={org.id} value={org.id}>
                      {org.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={3}>
              <FormControl fullWidth>
                <InputLabel>プロジェクト</InputLabel>
                <Select
                  value={selectedProjectId}
                  label="プロジェクト"
                  onChange={(e) =>
                    setSelectedProjectId(e.target.value as number | '')
                  }
                >
                  <MenuItem value="">すべて</MenuItem>
                  {projects.map((project) => (
                    <MenuItem key={project.id} value={project.id}>
                      {project.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={3}>
              <FormControl fullWidth>
                <InputLabel>ステータス</InputLabel>
                <Select
                  value={selectedStatus}
                  label="ステータス"
                  onChange={(e) =>
                    setSelectedStatus(e.target.value as DelayStatus | 'all')
                  }
                >
                  <MenuItem value="all">すべて</MenuItem>
                  <MenuItem value="RED">遅延</MenuItem>
                  <MenuItem value="YELLOW">注意</MenuItem>
                  <MenuItem value="GREEN">正常</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* チケット一覧テーブル */}
      <Typography variant="h6" gutterBottom>
        {filteredIssues.length} 件のチケット
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>キー</TableCell>
              <TableCell>プロジェクト</TableCell>
              <TableCell>概要</TableCell>
              <TableCell>担当者</TableCell>
              <TableCell>期日</TableCell>
              <TableCell>ステータス</TableCell>
              <TableCell>遅延状況</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredIssues.length > 0 ? (
              filteredIssues.map((issue) => (
                <TableRow
                  key={issue.id}
                  hover
                  sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                >
                  <TableCell>
                    <Typography variant="body2" fontWeight="bold">
                      {issue.jira_key}
                    </Typography>
                  </TableCell>
                  <TableCell>{getProjectName(issue.project_id)}</TableCell>
                  <TableCell>{issue.summary}</TableCell>
                  <TableCell>{issue.assignee || '未割り当て'}</TableCell>
                  <TableCell>{formatDate(issue.due_date)}</TableCell>
                  <TableCell>
                    <Chip label={issue.status} size="small" />
                  </TableCell>
                  <TableCell>
                    <StatusBadge status={issue.delay_status} />
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  <Typography variant="body2" color="textSecondary" sx={{ py: 3 }}>
                    条件に一致するチケットがありません
                  </Typography>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  )
}
