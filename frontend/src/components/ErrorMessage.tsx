import { Alert, Box, Button } from '@mui/material'

interface ErrorMessageProps {
  message: string
  onRetry?: () => void
}

export default function ErrorMessage({ message, onRetry }: ErrorMessageProps) {
  return (
    <Box sx={{ py: 2 }}>
      <Alert
        severity="error"
        action={
          onRetry ? (
            <Button color="inherit" size="small" onClick={onRetry}>
              再試行
            </Button>
          ) : undefined
        }
      >
        {message}
      </Alert>
    </Box>
  )
}
