import { Box, Paper, Typography, Tooltip, Grid } from '@mui/material'
import { ProjectWithStats, DelayStatus } from '@/types'

interface ProjectHeatmapProps {
  projects: ProjectWithStats[]
  title?: string
}

interface HeatmapCell {
  project: ProjectWithStats
  status: DelayStatus
}

export default function ProjectHeatmap({
  projects,
  title = 'プロジェクトヒートマップ',
}: ProjectHeatmapProps) {
  const getProjectStatus = (project: ProjectWithStats): DelayStatus => {
    if (project.red_issues > 0) return 'RED'
    if (project.yellow_issues > 0) return 'YELLOW'
    return 'GREEN'
  }

  const getStatusColor = (status: DelayStatus) => {
    switch (status) {
      case 'RED':
        return '#ef5350' // error.main
      case 'YELLOW':
        return '#ff9800' // warning.main
      case 'GREEN':
        return '#66bb6a' // success.main
      default:
        return '#9e9e9e' // grey
    }
  }

  const getStatusLabel = (status: DelayStatus) => {
    switch (status) {
      case 'RED':
        return '遅延'
      case 'YELLOW':
        return '注意'
      case 'GREEN':
        return '正常'
      default:
        return '不明'
    }
  }

  const heatmapCells: HeatmapCell[] = projects.map((project) => ({
    project,
    status: getProjectStatus(project),
  }))

  // プロジェクトを最大20個まで表示（グリッドが大きくなりすぎないように）
  const displayCells = heatmapCells.slice(0, 20)

  return (
    <Paper sx={{ p: 3 }}>
      <Typography variant="h6" gutterBottom>
        {title}
      </Typography>
      <Typography variant="body2" color="textSecondary" sx={{ mb: 3 }}>
        各プロジェクトの遅延状況を色で表示しています
      </Typography>

      <Grid container spacing={1}>
        {displayCells.map((cell) => (
          <Grid item xs={6} sm={4} md={3} lg={2.4} key={cell.project.id}>
            <Tooltip
              title={
                <Box>
                  <Typography variant="body2" fontWeight="bold">
                    {cell.project.name}
                  </Typography>
                  <Typography variant="caption" display="block">
                    キー: {cell.project.key}
                  </Typography>
                  <Typography variant="caption" display="block" sx={{ mt: 1 }}>
                    ステータス: {getStatusLabel(cell.status)}
                  </Typography>
                  <Typography variant="caption" display="block">
                    総チケット: {cell.project.total_issues}
                  </Typography>
                  <Typography variant="caption" display="block">
                    遅延: {cell.project.red_issues} / 注意:{' '}
                    {cell.project.yellow_issues} / 正常: {cell.project.green_issues}
                  </Typography>
                  <Typography variant="caption" display="block">
                    完了: {cell.project.done_issues}
                  </Typography>
                </Box>
              }
              arrow
              placement="top"
            >
              <Box
                sx={{
                  aspectRatio: '1',
                  bgcolor: getStatusColor(cell.status),
                  borderRadius: 1,
                  cursor: 'pointer',
                  transition: 'all 0.2s',
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'center',
                  p: 1,
                  '&:hover': {
                    transform: 'scale(1.05)',
                    boxShadow: 3,
                  },
                }}
              >
                <Typography
                  variant="caption"
                  sx={{
                    color: 'white',
                    fontWeight: 'bold',
                    textAlign: 'center',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    display: '-webkit-box',
                    WebkitLineClamp: 2,
                    WebkitBoxOrient: 'vertical',
                    lineHeight: 1.2,
                    fontSize: '0.7rem',
                  }}
                >
                  {cell.project.key}
                </Typography>
              </Box>
            </Tooltip>
          </Grid>
        ))}
      </Grid>

      {/* 凡例 */}
      <Box sx={{ mt: 3, display: 'flex', gap: 3, justifyContent: 'center' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Box
            sx={{
              width: 20,
              height: 20,
              bgcolor: getStatusColor('RED'),
              borderRadius: 1,
            }}
          />
          <Typography variant="caption">遅延</Typography>
        </Box>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Box
            sx={{
              width: 20,
              height: 20,
              bgcolor: getStatusColor('YELLOW'),
              borderRadius: 1,
            }}
          />
          <Typography variant="caption">注意</Typography>
        </Box>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Box
            sx={{
              width: 20,
              height: 20,
              bgcolor: getStatusColor('GREEN'),
              borderRadius: 1,
            }}
          />
          <Typography variant="caption">正常</Typography>
        </Box>
      </Box>

      {projects.length > 20 && (
        <Typography
          variant="caption"
          color="textSecondary"
          sx={{ mt: 2, display: 'block', textAlign: 'center' }}
        >
          {projects.length - 20} 件のプロジェクトが非表示です
        </Typography>
      )}
    </Paper>
  )
}
