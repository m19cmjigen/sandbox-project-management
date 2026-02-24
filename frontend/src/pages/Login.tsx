import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Paper,
  TextField,
  Typography,
} from '@mui/material'
import { BarChart as BarChartIcon } from '@mui/icons-material'
import { login } from '../api/auth'
import { useAuthStore } from '../stores/authStore'

export default function Login() {
  const navigate = useNavigate()
  const authLogin = useAuthStore((s) => s.login)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const { token, user } = await login({ email, password })
      authLogin(token, user)
      navigate('/', { replace: true })
    } catch (err: unknown) {
      if (
        err &&
        typeof err === 'object' &&
        'response' in err &&
        ((err as { response?: { status?: number } }).response?.status === 401 ||
         (err as { response?: { status?: number } }).response?.status === 400)
      ) {
        setError('メールアドレスまたはパスワードが正しくありません。')
      } else {
        setError('ログインに失敗しました。バックエンドの接続を確認してください。')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: '#0f172a',
      }}
    >
      <Paper
        elevation={0}
        sx={{
          width: '100%',
          maxWidth: 400,
          p: 4,
          bgcolor: '#1e293b',
          borderRadius: 3,
          border: '1px solid #334155',
        }}
      >
        {/* Logo */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 4 }}>
          <Box
            sx={{
              width: 40,
              height: 40,
              borderRadius: 2,
              bgcolor: '#6366f1',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <BarChartIcon sx={{ fontSize: 22, color: '#fff' }} />
          </Box>
          <Box>
            <Typography sx={{ color: '#f8fafc', fontWeight: 700, fontSize: '1rem', lineHeight: 1.2 }}>
              ProjectViz
            </Typography>
            <Typography sx={{ color: '#94a3b8', fontSize: '0.75rem' }}>進捗可視化プラットフォーム</Typography>
          </Box>
        </Box>

        <Typography variant="h6" sx={{ color: '#f8fafc', fontWeight: 600, mb: 0.5 }}>
          ログイン
        </Typography>
        <Typography sx={{ color: '#94a3b8', fontSize: '0.875rem', mb: 3 }}>
          アカウント情報を入力してください
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2, bgcolor: '#450a0a', color: '#fca5a5', '& .MuiAlert-icon': { color: '#f87171' } }}>
            {error}
          </Alert>
        )}

        <Box component="form" onSubmit={handleSubmit}>
          <TextField
            label="メールアドレス"
            type="email"
            fullWidth
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            disabled={loading}
            sx={{
              mb: 2,
              '& .MuiOutlinedInput-root': {
                color: '#f8fafc',
                '& fieldset': { borderColor: '#334155' },
                '&:hover fieldset': { borderColor: '#6366f1' },
                '&.Mui-focused fieldset': { borderColor: '#6366f1' },
              },
              '& .MuiInputLabel-root': { color: '#94a3b8' },
              '& .MuiInputLabel-root.Mui-focused': { color: '#6366f1' },
            }}
          />
          <TextField
            label="パスワード"
            type="password"
            fullWidth
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={loading}
            sx={{
              mb: 3,
              '& .MuiOutlinedInput-root': {
                color: '#f8fafc',
                '& fieldset': { borderColor: '#334155' },
                '&:hover fieldset': { borderColor: '#6366f1' },
                '&.Mui-focused fieldset': { borderColor: '#6366f1' },
              },
              '& .MuiInputLabel-root': { color: '#94a3b8' },
              '& .MuiInputLabel-root.Mui-focused': { color: '#6366f1' },
            }}
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            disabled={loading}
            sx={{
              bgcolor: '#6366f1',
              '&:hover': { bgcolor: '#4f46e5' },
              py: 1.25,
              fontWeight: 600,
              fontSize: '0.9375rem',
              textTransform: 'none',
              borderRadius: 2,
            }}
          >
            {loading ? <CircularProgress size={20} sx={{ color: '#fff' }} /> : 'ログイン'}
          </Button>
        </Box>
      </Paper>
    </Box>
  )
}
