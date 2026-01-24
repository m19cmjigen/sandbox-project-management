import { useEffect, useState } from 'react'
import { Box, Card, CardContent, Grid, Typography } from '@mui/material'
import {
  TrendingUp as TrendingUpIcon,
  Warning as WarningIcon,
  CheckCircle as CheckCircleIcon,
} from '@mui/icons-material'
import { DashboardSummary } from '@/types'
import { dashboardService } from '@/services/dashboardService'
import Loading from '@/components/Loading'
import ErrorMessage from '@/components/ErrorMessage'
import ProjectCard from '@/components/ProjectCard'

export default function Dashboard() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchSummary = async () => {
      try {
        setLoading(true)
        const data = await dashboardService.getSummary()
        setSummary(data)
      } catch (err) {
        setError('ダッシュボードデータの取得に失敗しました')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    fetchSummary()
  }, [])

  if (loading) return <Loading />
  if (error) return <ErrorMessage message={error} />
  if (!summary) return <ErrorMessage message="データが見つかりません" />

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        ダッシュボード
      </Typography>

      {/* サマリカード */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                総プロジェクト数
              </Typography>
              <Typography variant="h4">{summary.total_projects}</Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={3}>
          <Card sx={{ bgcolor: 'error.light', color: 'white' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <WarningIcon />
                <Typography gutterBottom>遅延プロジェクト</Typography>
              </Box>
              <Typography variant="h4">{summary.delayed_projects}</Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={3}>
          <Card sx={{ bgcolor: 'warning.light', color: 'white' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <TrendingUpIcon />
                <Typography gutterBottom>注意プロジェクト</Typography>
              </Box>
              <Typography variant="h4">{summary.warning_projects}</Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={3}>
          <Card sx={{ bgcolor: 'success.light', color: 'white' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <CheckCircleIcon />
                <Typography gutterBottom>正常プロジェクト</Typography>
              </Box>
              <Typography variant="h4">{summary.normal_projects}</Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* チケット集計 */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                チケット状況
              </Typography>
              <Grid container spacing={2}>
                <Grid item xs={3}>
                  <Typography variant="body2" color="textSecondary">
                    総チケット数
                  </Typography>
                  <Typography variant="h5">{summary.total_issues}</Typography>
                </Grid>
                <Grid item xs={3}>
                  <Typography variant="body2" color="error">
                    遅延チケット
                  </Typography>
                  <Typography variant="h5" color="error.main">
                    {summary.red_issues}
                  </Typography>
                </Grid>
                <Grid item xs={3}>
                  <Typography variant="body2" color="warning.main">
                    注意チケット
                  </Typography>
                  <Typography variant="h5" color="warning.main">
                    {summary.yellow_issues}
                  </Typography>
                </Grid>
                <Grid item xs={3}>
                  <Typography variant="body2" color="success.main">
                    正常チケット
                  </Typography>
                  <Typography variant="h5" color="success.main">
                    {summary.green_issues}
                  </Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* プロジェクト一覧 */}
      <Typography variant="h6" gutterBottom>
        プロジェクト一覧
      </Typography>
      <Grid container spacing={2}>
        {summary.projects_by_status && summary.projects_by_status.length > 0 ? (
          summary.projects_by_status.map((project) => (
            <Grid item xs={12} md={6} lg={4} key={project.id}>
              <ProjectCard project={project} />
            </Grid>
          ))
        ) : (
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="body2" color="textSecondary">
                  プロジェクトがありません
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>
    </Box>
  )
}
