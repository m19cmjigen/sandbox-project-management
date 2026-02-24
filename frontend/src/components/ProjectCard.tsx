import { useNavigate } from 'react-router-dom'
import {
  Box,
  Chip,
  IconButton,
  Tooltip,
  Typography,
  alpha,
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
    icon: <RedIcon sx={{ fontSize: 13 }} />,
    color: '#ef4444',
    bg: alpha('#ef4444', 0.08),
    chipBg: alpha('#ef4444', 0.1),
    chipColor: '#dc2626',
    accentBg: alpha('#ef4444', 0.06),
  },
  YELLOW: {
    label: '注意',
    icon: <YellowIcon sx={{ fontSize: 13 }} />,
    color: '#f59e0b',
    bg: alpha('#f59e0b', 0.08),
    chipBg: alpha('#f59e0b', 0.1),
    chipColor: '#d97706',
    accentBg: alpha('#f59e0b', 0.06),
  },
  GREEN: {
    label: '正常',
    icon: <GreenIcon sx={{ fontSize: 13 }} />,
    color: '#10b981',
    bg: alpha('#10b981', 0.08),
    chipBg: alpha('#10b981', 0.1),
    chipColor: '#059669',
    accentBg: alpha('#10b981', 0.06),
  },
}

export default function ProjectCard({ project }: ProjectCardProps) {
  const navigate = useNavigate()
  const cfg = statusConfig[project.delay_status]
  const jiraUrl = JIRA_BASE_URL ? `${JIRA_BASE_URL}/browse/${project.key}` : null

  return (
    <Box
      sx={{
        bgcolor: 'background.paper',
        borderRadius: 3,
        border: '1px solid',
        borderColor: 'divider',
        boxShadow: '0 1px 3px 0 rgb(0 0 0 / 0.06)',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden',
        transition: 'box-shadow 0.2s, transform 0.15s',
        '&:hover': {
          boxShadow: '0 4px 12px 0 rgb(0 0 0 / 0.1)',
          transform: 'translateY(-1px)',
        },
      }}
    >
      {/* Top accent strip */}
      <Box sx={{ height: 3, bgcolor: cfg.color, flexShrink: 0 }} />

      <Box sx={{ p: 2.5, flexGrow: 1, display: 'flex', flexDirection: 'column', gap: 2 }}>
        {/* Header row */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Typography
            variant="caption"
            sx={{
              fontFamily: '"JetBrains Mono", "Fira Code", monospace',
              fontSize: '0.7rem',
              fontWeight: 600,
              color: 'text.secondary',
              bgcolor: 'grey.100',
              px: 0.75,
              py: 0.25,
              borderRadius: 1,
              letterSpacing: '0.02em',
            }}
          >
            {project.key}
          </Typography>

          {/* Status badge */}
          <Box
            sx={{
              display: 'inline-flex',
              alignItems: 'center',
              gap: 0.5,
              px: 0.875,
              py: 0.25,
              borderRadius: 1,
              bgcolor: cfg.chipBg,
              color: cfg.chipColor,
              fontSize: '0.6875rem',
              fontWeight: 600,
            }}
          >
            {cfg.icon}
            {cfg.label}
          </Box>

          <Box sx={{ flexGrow: 1 }} />

          <Tooltip title="チケット一覧">
            <IconButton
              size="small"
              onClick={() => navigate(`/issues?project_id=${project.id}`)}
              sx={{
                p: 0.5,
                color: 'text.secondary',
                '&:hover': { color: 'primary.main', bgcolor: alpha('#6366f1', 0.08) },
                borderRadius: 1.5,
              }}
            >
              <IssuesIcon sx={{ fontSize: 16 }} />
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
                sx={{
                  p: 0.5,
                  color: 'text.secondary',
                  '&:hover': { color: 'primary.main', bgcolor: alpha('#6366f1', 0.08) },
                  borderRadius: 1.5,
                }}
              >
                <ExternalLinkIcon sx={{ fontSize: 16 }} />
              </IconButton>
            </Tooltip>
          )}
        </Box>

        {/* Project name */}
        <Typography
          variant="subtitle2"
          fontWeight={600}
          sx={{ lineHeight: 1.5, color: 'text.primary', flexGrow: 1 }}
        >
          {project.name}
        </Typography>

        {/* Issue count chips */}
        <Box sx={{ display: 'flex', gap: 0.75, flexWrap: 'wrap' }}>
          {project.red_count > 0 && (
            <Tooltip title="期限切れチケット数">
              <Chip
                label={`遅延 ${project.red_count}`}
                size="small"
                sx={{
                  height: 20,
                  fontSize: '0.6875rem',
                  fontWeight: 700,
                  bgcolor: alpha('#ef4444', 0.1),
                  color: '#dc2626',
                  border: 'none',
                }}
              />
            </Tooltip>
          )}
          {project.yellow_count > 0 && (
            <Tooltip title="期限切れ間近（3日以内）">
              <Chip
                label={`注意 ${project.yellow_count}`}
                size="small"
                sx={{
                  height: 20,
                  fontSize: '0.6875rem',
                  fontWeight: 700,
                  bgcolor: alpha('#f59e0b', 0.1),
                  color: '#d97706',
                  border: 'none',
                }}
              />
            </Tooltip>
          )}
          <Tooltip title="未完了チケット数">
            <Chip
              label={`未完了 ${project.open_count}`}
              size="small"
              sx={{
                height: 20,
                fontSize: '0.6875rem',
                fontWeight: 500,
                bgcolor: 'grey.100',
                color: 'text.secondary',
                border: 'none',
              }}
            />
          </Tooltip>
          <Tooltip title="総チケット数">
            <Chip
              label={`全 ${project.total_count}`}
              size="small"
              sx={{
                height: 20,
                fontSize: '0.6875rem',
                fontWeight: 500,
                bgcolor: 'grey.100',
                color: 'text.secondary',
                border: 'none',
              }}
            />
          </Tooltip>
        </Box>
      </Box>
    </Box>
  )
}
