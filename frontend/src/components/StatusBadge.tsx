import { Box, alpha } from '@mui/material'
import type { DelayStatus } from '../types/project'

const STATUS_CONFIG: Record<DelayStatus, { label: string; color: string }> = {
  RED:    { label: '遅延', color: '#ef4444' },
  YELLOW: { label: '注意', color: '#f59e0b' },
  GREEN:  { label: '正常', color: '#10b981' },
}

interface StatusBadgeProps {
  status: DelayStatus
  size?: 'small' | 'medium'
}

export default function StatusBadge({ status, size = 'small' }: StatusBadgeProps) {
  const cfg = STATUS_CONFIG[status]
  const isSmall = size === 'small'

  return (
    <Box
      sx={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: 0.5,
        px: isSmall ? 0.875 : 1.25,
        py: isSmall ? 0.25 : 0.5,
        borderRadius: 1,
        bgcolor: alpha(cfg.color, 0.1),
        fontSize: isSmall ? '0.6875rem' : '0.8125rem',
        fontWeight: 600,
        color: cfg.color,
        whiteSpace: 'nowrap',
        lineHeight: 1.4,
      }}
    >
      <Box
        sx={{
          width: isSmall ? 6 : 8,
          height: isSmall ? 6 : 8,
          borderRadius: '50%',
          bgcolor: cfg.color,
          flexShrink: 0,
        }}
      />
      {cfg.label}
    </Box>
  )
}
