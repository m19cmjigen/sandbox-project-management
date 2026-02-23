import { Chip } from '@mui/material'
import type { DelayStatus } from '../types/project'

const STATUS_COLOR: Record<DelayStatus, 'error' | 'warning' | 'success'> = {
  RED: 'error',
  YELLOW: 'warning',
  GREEN: 'success',
}

const STATUS_LABEL: Record<DelayStatus, string> = {
  RED: '遅延',
  YELLOW: '注意',
  GREEN: '正常',
}

interface StatusBadgeProps {
  status: DelayStatus
  size?: 'small' | 'medium'
}

export default function StatusBadge({ status, size = 'small' }: StatusBadgeProps) {
  return (
    <Chip
      label={STATUS_LABEL[status]}
      color={STATUS_COLOR[status]}
      size={size}
    />
  )
}
