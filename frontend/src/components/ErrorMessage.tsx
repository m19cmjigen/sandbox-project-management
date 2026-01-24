import { Alert, AlertTitle } from '@mui/material'

interface ErrorMessageProps {
  title?: string
  message: string
}

export default function ErrorMessage({
  title = 'エラー',
  message,
}: ErrorMessageProps) {
  return (
    <Alert severity="error">
      <AlertTitle>{title}</AlertTitle>
      {message}
    </Alert>
  )
}
