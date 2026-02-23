import { Box, CircularProgress } from '@mui/material'

interface LoadingSpinnerProps {
  size?: number
  minHeight?: string | number
}

export default function LoadingSpinner({ size = 40, minHeight = 200 }: LoadingSpinnerProps) {
  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      minHeight={minHeight}
    >
      <CircularProgress size={size} />
    </Box>
  )
}
