import { useNavigate } from 'react-router-dom'
import { Box, Chip, Paper, Tooltip, Typography } from '@mui/material'
import type { DashboardOrgNode } from '../types/dashboard'
import type { DelayStatus } from '../types/project'

const STATUS_CONFIG: Record<DelayStatus, { label: string; bg: string; border: string; chipColor: 'error' | 'warning' | 'success' }> = {
  RED:    { label: '遅延あり', bg: '#fff5f5', border: '#f44336', chipColor: 'error' },
  YELLOW: { label: '注意',     bg: '#fffde7', border: '#ff9800', chipColor: 'warning' },
  GREEN:  { label: '正常',     bg: '#f1f8e9', border: '#4caf50', chipColor: 'success' },
}

interface HeatmapCardProps {
  node: DashboardOrgNode
  /** When true, renders as a child card (smaller) */
  isChild?: boolean
}

function DelayRateBar({ rate }: { rate: number }) {
  const pct = Math.round(rate * 100)
  const barColor = pct >= 50 ? '#f44336' : pct >= 20 ? '#ff9800' : '#4caf50'
  return (
    <Box sx={{ mt: 1 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.25 }}>
        <Typography variant="caption" color="text.secondary">遅延率</Typography>
        <Typography variant="caption" fontWeight="bold" color={barColor}>{pct}%</Typography>
      </Box>
      <Box sx={{ height: 6, bgcolor: 'grey.200', borderRadius: 3, overflow: 'hidden' }}>
        <Box sx={{ height: '100%', width: `${pct}%`, bgcolor: barColor, borderRadius: 3, transition: 'width 0.4s' }} />
      </Box>
    </Box>
  )
}

export default function HeatmapCard({ node, isChild = false }: HeatmapCardProps) {
  const navigate = useNavigate()
  const cfg = STATUS_CONFIG[node.delay_status]

  const tooltipContent = (
    <Box sx={{ p: 0.5 }}>
      <Typography variant="body2" fontWeight="bold">{node.name}</Typography>
      <Typography variant="caption" display="block">総プロジェクト: {node.total_projects}</Typography>
      <Typography variant="caption" display="block" color="error.light">遅延: {node.red_projects}</Typography>
      <Typography variant="caption" display="block" color="warning.light">注意: {node.yellow_projects}</Typography>
      <Typography variant="caption" display="block" color="success.light">正常: {node.green_projects}</Typography>
    </Box>
  )

  return (
    <Tooltip title={tooltipContent} arrow placement="top">
      <Paper
        variant="outlined"
        onClick={() => navigate(`/projects?organization_id=${node.id}`)}
        sx={{
          p: isChild ? 1.5 : 2,
          bgcolor: cfg.bg,
          borderLeft: `4px solid ${cfg.border}`,
          cursor: 'pointer',
          transition: 'box-shadow 0.2s, transform 0.15s',
          '&:hover': { boxShadow: 4, transform: 'translateY(-2px)' },
          height: '100%',
          boxSizing: 'border-box',
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'flex-start', gap: 1, mb: 0.5 }}>
          <Typography
            variant={isChild ? 'body2' : 'subtitle2'}
            fontWeight="bold"
            sx={{ flexGrow: 1, lineHeight: 1.3 }}
          >
            {node.name}
          </Typography>
          <Chip
            label={cfg.label}
            color={cfg.chipColor}
            size="small"
            variant="outlined"
            sx={{ flexShrink: 0 }}
          />
        </Box>

        {!isChild && (
          <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mb: 0.5 }}>
            {node.red_projects > 0 && (
              <Typography variant="caption" color="error.main" fontWeight="bold">
                遅延 {node.red_projects}
              </Typography>
            )}
            {node.yellow_projects > 0 && (
              <Typography variant="caption" color="warning.main" fontWeight="bold">
                注意 {node.yellow_projects}
              </Typography>
            )}
            <Typography variant="caption" color="text.secondary">
              全 {node.total_projects}件
            </Typography>
          </Box>
        )}

        {isChild && (
          <Typography variant="caption" color="text.secondary">
            {node.total_projects}件
            {node.red_projects > 0 && (
              <Box component="span" sx={{ color: 'error.main', fontWeight: 'bold', ml: 0.5 }}>
                遅延{node.red_projects}
              </Box>
            )}
          </Typography>
        )}

        {node.total_projects > 0 && <DelayRateBar rate={node.delay_rate} />}
      </Paper>
    </Tooltip>
  )
}
