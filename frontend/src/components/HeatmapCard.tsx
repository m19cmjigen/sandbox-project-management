import { useNavigate } from 'react-router-dom'
import { Box, Tooltip, Typography, alpha } from '@mui/material'
import type { DashboardOrgNode } from '../types/dashboard'
import type { DelayStatus } from '../types/project'

const STATUS_CONFIG: Record<DelayStatus, { label: string; color: string; bg: string }> = {
  RED:    { label: '遅延あり', color: '#ef4444', bg: alpha('#ef4444', 0.08) },
  YELLOW: { label: '注意',    color: '#f59e0b', bg: alpha('#f59e0b', 0.08) },
  GREEN:  { label: '正常',    color: '#10b981', bg: alpha('#10b981', 0.08) },
}

interface HeatmapCardProps {
  node: DashboardOrgNode
  isChild?: boolean
}

function DelayRateBar({ rate }: { rate: number }) {
  const pct = Math.round(rate * 100)
  const color = pct >= 50 ? '#ef4444' : pct >= 20 ? '#f59e0b' : '#10b981'
  return (
    <Box sx={{ mt: 1.5 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
        <Typography variant="caption" color="text.secondary" sx={{ fontSize: '0.6875rem' }}>
          遅延率
        </Typography>
        <Typography
          variant="caption"
          fontWeight={700}
          sx={{ color, fontSize: '0.6875rem' }}
        >
          {pct}%
        </Typography>
      </Box>
      <Box
        sx={{
          height: 5,
          bgcolor: 'grey.100',
          borderRadius: 99,
          overflow: 'hidden',
        }}
      >
        <Box
          sx={{
            height: '100%',
            width: `${pct}%`,
            bgcolor: color,
            borderRadius: 99,
            transition: 'width 0.5s cubic-bezier(0.4, 0, 0.2, 1)',
          }}
        />
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
      <Typography variant="caption" display="block" sx={{ color: '#fca5a5' }}>遅延: {node.red_projects}</Typography>
      <Typography variant="caption" display="block" sx={{ color: '#fcd34d' }}>注意: {node.yellow_projects}</Typography>
      <Typography variant="caption" display="block" sx={{ color: '#6ee7b7' }}>正常: {node.green_projects}</Typography>
    </Box>
  )

  return (
    <Tooltip title={tooltipContent} arrow placement="top">
      <Box
        onClick={() => navigate(`/projects?organization_id=${node.id}`)}
        sx={{
          bgcolor: 'background.paper',
          border: '1px solid',
          borderColor: 'divider',
          borderRadius: 2.5,
          p: isChild ? 1.75 : 2.5,
          cursor: 'pointer',
          height: '100%',
          boxSizing: 'border-box',
          boxShadow: '0 1px 3px 0 rgb(0 0 0 / 0.06)',
          transition: 'box-shadow 0.2s, transform 0.15s, border-color 0.2s',
          position: 'relative',
          overflow: 'hidden',
          '&:hover': {
            boxShadow: '0 4px 12px 0 rgb(0 0 0 / 0.1)',
            transform: 'translateY(-2px)',
            borderColor: cfg.color,
          },
          // Top accent strip
          '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            height: 3,
            bgcolor: cfg.color,
          },
        }}
      >
        {/* Header */}
        <Box sx={{ display: 'flex', alignItems: 'flex-start', gap: 1, mb: 0.5 }}>
          <Typography
            variant={isChild ? 'body2' : 'subtitle2'}
            fontWeight={600}
            color="text.primary"
            sx={{ flexGrow: 1, lineHeight: 1.4 }}
          >
            {node.name}
          </Typography>

          {/* Status dot badge */}
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 0.5,
              px: 0.75,
              py: 0.25,
              borderRadius: 1,
              bgcolor: alpha(cfg.color, 0.1),
              flexShrink: 0,
            }}
          >
            <Box
              sx={{
                width: 6,
                height: 6,
                borderRadius: '50%',
                bgcolor: cfg.color,
              }}
            />
            <Typography
              sx={{
                fontSize: '0.6875rem',
                fontWeight: 600,
                color: cfg.color,
                lineHeight: 1,
              }}
            >
              {cfg.label}
            </Typography>
          </Box>
        </Box>

        {/* Stats */}
        {!isChild && (
          <Box sx={{ display: 'flex', gap: 1.5, mt: 1 }}>
            {node.red_projects > 0 && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                <Box sx={{ width: 8, height: 8, borderRadius: '50%', bgcolor: '#ef4444' }} />
                <Typography variant="caption" fontWeight={700} sx={{ color: '#ef4444', fontSize: '0.75rem' }}>
                  {node.red_projects}
                </Typography>
              </Box>
            )}
            {node.yellow_projects > 0 && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                <Box sx={{ width: 8, height: 8, borderRadius: '50%', bgcolor: '#f59e0b' }} />
                <Typography variant="caption" fontWeight={700} sx={{ color: '#f59e0b', fontSize: '0.75rem' }}>
                  {node.yellow_projects}
                </Typography>
              </Box>
            )}
            <Typography variant="caption" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
              全{node.total_projects}件
            </Typography>
          </Box>
        )}

        {isChild && (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.75, mt: 0.5 }}>
            <Typography variant="caption" color="text.secondary" sx={{ fontSize: '0.6875rem' }}>
              {node.total_projects}件
            </Typography>
            {node.red_projects > 0 && (
              <Typography
                variant="caption"
                sx={{ color: '#ef4444', fontWeight: 700, fontSize: '0.6875rem' }}
              >
                遅延{node.red_projects}
              </Typography>
            )}
          </Box>
        )}

        {node.total_projects > 0 && <DelayRateBar rate={node.delay_rate} />}
      </Box>
    </Tooltip>
  )
}
