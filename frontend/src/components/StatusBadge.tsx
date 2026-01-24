import { Chip } from '@mui/material'
import { DelayStatus } from '@/types'

interface StatusBadgeProps {
  status: DelayStatus
  size?: 'small' | 'medium'
}

export default function StatusBadge({ status, size = 'small' }: StatusBadgeProps) {
  const getColor = () => {
    switch (status) {
      case 'RED':
        return 'error'
      case 'YELLOW':
        return 'warning'
      case 'GREEN':
        return 'success'
      default:
        return 'default'
    }
  }

  const getLabel = () => {
    switch (status) {
      case 'RED':
        return '遅延'
      case 'YELLOW':
        return '注意'
      case 'GREEN':
        return '正常'
      default:
        return status
    }
  }

  return <Chip label={getLabel()} color={getColor()} size={size} />
}
