import { useNavigate } from 'react-router-dom'
import {
  Box,
  Card,
  CardContent,
  Chip,
  IconButton,
  Tooltip,
  Typography,
} from '@mui/material'
import {
  ErrorOutline as RedIcon,
  WarningAmber as YellowIcon,
  CheckCircleOutline as GreenIcon,
  OpenInNew as ExternalLinkIcon,
  Assignment as IssuesIcon,
} from '@mui/icons-material'
import type { Project } from '../types/project'

const JIRA_BASE_URL = import.meta.env.VITE_JIRA_BASE_URL || ''

interface ProjectCardProps {
  project: Project
}

const statusConfig = {
  RED: {
    label: '遅延あり',
    color: 'error' as const,
    icon: <RedIcon fontSize="small" />,
    bgColor: '#fff5f5',
    borderColor: '#f44336',
  },
  YELLOW: {
    label: '注意',
    color: 'warning' as const,
    icon: <YellowIcon fontSize="small" />,
    bgColor: '#fffde7',
    borderColor: '#ff9800',
  },
  GREEN: {
    label: '正常',
    color: 'success' as const,
    icon: <GreenIcon fontSize="small" />,
    bgColor: '#f1f8e9',
    borderColor: '#4caf50',
  },
}

export default function ProjectCard({ project }: ProjectCardProps) {
  const navigate = useNavigate()
  const config = statusConfig[project.delay_status]
  const jiraUrl = JIRA_BASE_URL
    ? `${JIRA_BASE_URL}/browse/${project.key}`
    : null

  return (
    <Card
      variant="outlined"
      sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        borderLeft: `4px solid ${config.borderColor}`,
        bgcolor: config.bgColor,
        transition: 'box-shadow 0.2s',
        '&:hover': {
          boxShadow: 3,
        },
      }}
    >
      <CardContent sx={{ flexGrow: 1, pb: 1 }}>
        {/* Header: key + status chip + Jira link */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
          <Typography
            variant="caption"
            sx={{ fontFamily: 'monospace', color: 'text.secondary', fontWeight: 'bold' }}
          >
            {project.key}
          </Typography>
          <Chip
            icon={config.icon}
            label={config.label}
            color={config.color}
            size="small"
            variant="outlined"
          />
          <Box sx={{ flexGrow: 1 }} />
          <Tooltip title="チケット一覧">
            <IconButton
              size="small"
              color="primary"
              onClick={() => navigate(`/issues?project_id=${project.id}`)}
            >
              <IssuesIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          {jiraUrl && (
            <Tooltip title="Jiraで開く">
              <IconButton
                size="small"
                component="a"
                href={jiraUrl}
                target="_blank"
                rel="noopener noreferrer"
                color="primary"
              >
                <ExternalLinkIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          )}
        </Box>

        {/* Project name */}
        <Typography variant="subtitle1" fontWeight="bold" sx={{ mb: 2, lineHeight: 1.4 }}>
          {project.name}
        </Typography>

        {/* Issue counts */}
        <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
          {project.red_count > 0 && (
            <Tooltip title="期限切れチケット数">
              <Chip
                label={`遅延 ${project.red_count}`}
                size="small"
                color="error"
                sx={{ fontWeight: 'bold' }}
              />
            </Tooltip>
          )}
          {project.yellow_count > 0 && (
            <Tooltip title="期限切れ間近チケット数（3日以内）">
              <Chip
                label={`注意 ${project.yellow_count}`}
                size="small"
                color="warning"
                sx={{ fontWeight: 'bold' }}
              />
            </Tooltip>
          )}
          <Tooltip title="未完了チケット数">
            <Chip
              label={`未完了 ${project.open_count}`}
              size="small"
              variant="outlined"
            />
          </Tooltip>
          <Tooltip title="総チケット数">
            <Chip
              label={`全 ${project.total_count}`}
              size="small"
              variant="outlined"
              color="default"
            />
          </Tooltip>
        </Box>
      </CardContent>
    </Card>
  )
}
