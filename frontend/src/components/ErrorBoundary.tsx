import { Component, ErrorInfo, ReactNode } from 'react'
import { Box, Typography, Button, Card, CardContent, Alert } from '@mui/material'
import { ErrorOutline as ErrorIcon, Refresh as RefreshIcon } from '@mui/icons-material'

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
  errorInfo: ErrorInfo | null
}

class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
    }
  }

  static getDerivedStateFromError(error: Error): Partial<State> {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ErrorBoundary caught an error:', error, errorInfo)
    this.setState({
      error,
      errorInfo,
    })
  }

  handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
    })
    window.location.reload()
  }

  render() {
    if (this.state.hasError) {
      return (
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: '100vh',
            bgcolor: 'background.default',
            p: 3,
          }}
        >
          <Card sx={{ maxWidth: 600 }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <ErrorIcon color="error" sx={{ fontSize: 48, mr: 2 }} />
                <Typography variant="h4">エラーが発生しました</Typography>
              </Box>

              <Alert severity="error" sx={{ mb: 3 }}>
                アプリケーションで予期しないエラーが発生しました。
                ページをリロードして再試行してください。
              </Alert>

              {this.state.error && (
                <Box sx={{ mb: 3 }}>
                  <Typography variant="subtitle2" color="error" gutterBottom>
                    エラー詳細:
                  </Typography>
                  <Typography
                    variant="body2"
                    component="pre"
                    sx={{
                      bgcolor: 'grey.100',
                      p: 2,
                      borderRadius: 1,
                      overflow: 'auto',
                      fontSize: '0.875rem',
                      fontFamily: 'monospace',
                    }}
                  >
                    {this.state.error.toString()}
                  </Typography>
                </Box>
              )}

              {process.env.NODE_ENV === 'development' && this.state.errorInfo && (
                <Box sx={{ mb: 3 }}>
                  <Typography variant="subtitle2" color="textSecondary" gutterBottom>
                    スタックトレース:
                  </Typography>
                  <Typography
                    variant="body2"
                    component="pre"
                    sx={{
                      bgcolor: 'grey.100',
                      p: 2,
                      borderRadius: 1,
                      overflow: 'auto',
                      fontSize: '0.75rem',
                      fontFamily: 'monospace',
                      maxHeight: 200,
                    }}
                  >
                    {this.state.errorInfo.componentStack}
                  </Typography>
                </Box>
              )}

              <Box sx={{ display: 'flex', gap: 2 }}>
                <Button
                  variant="contained"
                  startIcon={<RefreshIcon />}
                  onClick={this.handleReset}
                  fullWidth
                >
                  ページをリロード
                </Button>
                <Button
                  variant="outlined"
                  onClick={() => window.history.back()}
                  fullWidth
                >
                  前のページに戻る
                </Button>
              </Box>

              <Typography
                variant="caption"
                color="textSecondary"
                sx={{ mt: 2, display: 'block', textAlign: 'center' }}
              >
                問題が解決しない場合は、ブラウザのキャッシュをクリアするか、
                管理者にお問い合わせください。
              </Typography>
            </CardContent>
          </Card>
        </Box>
      )
    }

    return this.props.children
  }
}

export default ErrorBoundary
