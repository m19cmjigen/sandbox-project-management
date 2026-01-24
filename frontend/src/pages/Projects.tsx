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
} from '@mui/material'
import { ProjectWithStats, Organization } from '@/types'
import { projectService } from '@/services/projectService'
import { organizationService } from '@/services/organizationService'
import Loading from '@/components/Loading'
import ErrorMessage from '@/components/ErrorMessage'
import ProjectCard from '@/components/ProjectCard'

export default function Projects() {
  const [projects, setProjects] = useState<ProjectWithStats[]>([])
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // フィルター状態
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedOrgId, setSelectedOrgId] = useState<number | ''>('')
  const [selectedStatus, setSelectedStatus] = useState<string>('all')

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)
        const [projectsData, orgsData] = await Promise.all([
          projectService.getAll(),
          organizationService.getAll(),
        ])
        setProjects(projectsData)
        setOrganizations(orgsData)
      } catch (err) {
        setError('データの取得に失敗しました')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  const getProjectStatus = (project: ProjectWithStats) => {
    if (project.red_issues > 0) return 'RED'
    if (project.yellow_issues > 0) return 'YELLOW'
    return 'GREEN'
  }

  const filteredProjects = projects.filter((project) => {
    // 検索条件
    const matchesSearch =
      searchTerm === '' ||
      project.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      project.key.toLowerCase().includes(searchTerm.toLowerCase())

    // 組織フィルター
    const matchesOrg =
      selectedOrgId === '' || project.organization_id === selectedOrgId

    // ステータスフィルター
    const projectStatus = getProjectStatus(project)
    const matchesStatus =
      selectedStatus === 'all' || projectStatus === selectedStatus

    return matchesSearch && matchesOrg && matchesStatus
  })

  if (loading) return <Loading />
  if (error) return <ErrorMessage message={error} />

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        プロジェクト一覧
      </Typography>

      {/* フィルター */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2}>
            <Grid item xs={12} md={4}>
              <TextField
                fullWidth
                label="検索"
                placeholder="プロジェクト名またはキー"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </Grid>
            <Grid item xs={12} md={4}>
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
            <Grid item xs={12} md={4}>
              <FormControl fullWidth>
                <InputLabel>ステータス</InputLabel>
                <Select
                  value={selectedStatus}
                  label="ステータス"
                  onChange={(e) => setSelectedStatus(e.target.value)}
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

      {/* プロジェクト一覧 */}
      <Typography variant="h6" gutterBottom>
        {filteredProjects.length} 件のプロジェクト
      </Typography>
      <Grid container spacing={2}>
        {filteredProjects.length > 0 ? (
          filteredProjects.map((project) => (
            <Grid item xs={12} md={6} lg={4} key={project.id}>
              <ProjectCard
                project={project}
                onClick={() => {
                  // TODO: プロジェクト詳細ページへ遷移
                  console.log('Navigate to project:', project.id)
                }}
              />
            </Grid>
          ))
        ) : (
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="body2" color="textSecondary">
                  条件に一致するプロジェクトがありません
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>
    </Box>
  )
}
