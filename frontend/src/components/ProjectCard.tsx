import {
  Card,
  CardContent,
  Typography,
  Box,
  LinearProgress,
  Stack,
  Chip,
} from '@mui/material'
import { ProjectWithStats } from '@/types'
import StatusBadge from './StatusBadge'

interface ProjectCardProps {
  project: ProjectWithStats
  onClick?: () => void
}

export default function ProjectCard({ project, onClick }: ProjectCardProps) {
  const getProjectStatus = () => {
    if (project.red_issues > 0) return 'RED'
    if (project.yellow_issues > 0) return 'YELLOW'
    return 'GREEN'
  }

  const completionRate =
    project.total_issues > 0
      ? Math.round((project.done_issues / project.total_issues) * 100)
      : 0

  return (
    <Card
      sx={{
        cursor: onClick ? 'pointer' : 'default',
        '&:hover': onClick
          ? {
              boxShadow: 3,
              transform: 'translateY(-2px)',
              transition: 'all 0.2s',
            }
          : {},
      }}
      onClick={onClick}
    >
      <CardContent>
        <Box sx={{ mb: 2 }}>
          <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1 }}>
            <Typography variant="h6" component="div">
              {project.name}
            </Typography>
            <StatusBadge status={getProjectStatus()} />
          </Stack>
          <Typography variant="body2" color="text.secondary">
            {project.key}
          </Typography>
        </Box>

        <Box sx={{ mb: 2 }}>
          <Typography variant="body2" color="text.secondary" gutterBottom>
            進捗率: {completionRate}%
          </Typography>
          <LinearProgress variant="determinate" value={completionRate} />
        </Box>

        <Stack direction="row" spacing={2}>
          <Box>
            <Typography variant="caption" color="text.secondary">
              総チケット数
            </Typography>
            <Typography variant="h6">{project.total_issues}</Typography>
          </Box>
          <Box>
            <Typography variant="caption" color="error">
              遅延
            </Typography>
            <Typography variant="h6" color="error.main">
              {project.red_issues}
            </Typography>
          </Box>
          <Box>
            <Typography variant="caption" color="warning.main">
              注意
            </Typography>
            <Typography variant="h6" color="warning.main">
              {project.yellow_issues}
            </Typography>
          </Box>
          <Box>
            <Typography variant="caption" color="success.main">
              正常
            </Typography>
            <Typography variant="h6" color="success.main">
              {project.green_issues}
            </Typography>
          </Box>
        </Stack>
      </CardContent>
    </Card>
  )
}
