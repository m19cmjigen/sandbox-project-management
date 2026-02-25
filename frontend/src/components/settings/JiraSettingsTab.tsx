import { useCallback, useEffect, useState } from 'react'
import {
  Alert,
  Box,
  Button,
  Chip,
  CircularProgress,
  Divider,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from '@mui/material'
import {
  Refresh as RefreshIcon,
  Save as SaveIcon,
  CheckCircle as TestIcon,
  Sync as SyncIcon,
} from '@mui/icons-material'
import { getJiraSettings, updateJiraSettings, testJiraConnection, triggerSync } from '../../api/settings'
import { getSyncLogs } from '../../api/syncLogs'
import type { JiraSettings, SyncLog } from '../../types/settings'

const statusColor: Record<string, 'info' | 'success' | 'error' | 'default'> = {
  RUNNING: 'info',
  SUCCESS: 'success',
  FAILED: 'error',
}

export default function JiraSettingsTab() {
  const [settings, setSettings] = useState<JiraSettings | null>(null)
  const [logs, setLogs] = useState<SyncLog[]>([])
  const [loading, setLoading] = useState(true)
  const [formUrl, setFormUrl] = useState('')
  const [formEmail, setFormEmail] = useState('')
  const [formToken, setFormToken] = useState('')
  const [saving, setSaving] = useState(false)
  const [testing, setTesting] = useState(false)
  const [syncing, setSyncing] = useState(false)
  const [successMsg, setSuccessMsg] = useState<string | null>(null)
  const [errorMsg, setErrorMsg] = useState<string | null>(null)

  const showSuccess = (msg: string) => {
    setSuccessMsg(msg)
    setTimeout(() => setSuccessMsg(null), 4000)
  }

  const showError = (msg: string) => {
    setErrorMsg(msg)
  }

  const loadData = useCallback(async () => {
    setLoading(true)
    try {
      const [s, l] = await Promise.all([getJiraSettings(), getSyncLogs()])
      setSettings(s)
      setLogs(l)
      if (s.configured) {
        setFormUrl(s.jira_url)
        setFormEmail(s.email)
      }
    } catch {
      showError('設定の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadData()
  }, [loadData])

  const handleSave = async () => {
    if (!formUrl || !formEmail || !formToken) {
      showError('すべての項目を入力してください')
      return
    }
    setSaving(true)
    try {
      await updateJiraSettings({ jira_url: formUrl, email: formEmail, api_token: formToken })
      showSuccess('Jira設定を保存しました')
      setFormToken('')
      loadData()
    } catch (e: unknown) {
      const msg = (e as { response?: { data?: { error?: string } } })?.response?.data?.error
      showError(msg ?? '保存に失敗しました')
    } finally {
      setSaving(false)
    }
  }

  const handleTest = async () => {
    setTesting(true)
    try {
      // フォームにトークンが入力済みの場合はその値を使用、なければDB設定を使用
      const payload = formToken
        ? { jira_url: formUrl, email: formEmail, api_token: formToken }
        : {}
      await testJiraConnection(payload)
      showSuccess('Jira接続に成功しました')
    } catch (e: unknown) {
      const msg = (e as { response?: { data?: { error?: string } } })?.response?.data?.error
      showError(msg ?? 'Jira接続テストに失敗しました')
    } finally {
      setTesting(false)
    }
  }

  const handleSync = async () => {
    setSyncing(true)
    try {
      await triggerSync()
      showSuccess('同期を開始しました')
      // 少し待ってからログを再取得
      setTimeout(() => loadData(), 1500)
    } catch (e: unknown) {
      const msg = (e as { response?: { data?: { error?: string } } })?.response?.data?.error
      showError(msg ?? '同期の開始に失敗しました')
    } finally {
      setSyncing(false)
    }
  }

  if (loading) return <CircularProgress sx={{ display: 'block', mx: 'auto', mt: 4 }} />

  return (
    <Box>
      {successMsg && <Alert severity="success" sx={{ mb: 2 }}>{successMsg}</Alert>}
      {errorMsg && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setErrorMsg(null)}>
          {errorMsg}
        </Alert>
      )}

      {/* Jira接続設定フォーム */}
      <Paper variant="outlined" sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" sx={{ mb: 2 }}>Jira接続設定</Typography>
        <Divider sx={{ mb: 2 }} />
        <Stack spacing={2} sx={{ maxWidth: 480 }}>
          <TextField
            label="Jira URL"
            placeholder="https://yourcompany.atlassian.net"
            value={formUrl}
            onChange={(e) => setFormUrl(e.target.value)}
            size="small"
            fullWidth
            helperText={settings?.configured ? `現在: ${settings.jira_url}` : '未設定'}
          />
          <TextField
            label="メールアドレス"
            value={formEmail}
            onChange={(e) => setFormEmail(e.target.value)}
            size="small"
            fullWidth
          />
          <TextField
            label="APIトークン"
            type="password"
            value={formToken}
            onChange={(e) => setFormToken(e.target.value)}
            size="small"
            fullWidth
            helperText={settings?.configured ? `現在: ${settings.api_token_mask}` : '未設定'}
            placeholder={settings?.configured ? '変更する場合のみ入力' : ''}
          />
        </Stack>
        <Stack direction="row" spacing={2} sx={{ mt: 3 }}>
          <Button
            variant="contained"
            startIcon={saving ? <CircularProgress size={16} color="inherit" /> : <SaveIcon />}
            onClick={handleSave}
            disabled={saving}
          >
            保存
          </Button>
          <Button
            variant="outlined"
            startIcon={testing ? <CircularProgress size={16} color="inherit" /> : <TestIcon />}
            onClick={handleTest}
            disabled={testing || (!settings?.configured && !formToken)}
          >
            接続テスト
          </Button>
        </Stack>
      </Paper>

      {/* 同期管理 */}
      <Paper variant="outlined" sx={{ p: 3 }}>
        <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 2 }}>
          <Typography variant="h6">データ同期</Typography>
          <Stack direction="row" spacing={1}>
            <Button size="small" startIcon={<RefreshIcon />} onClick={loadData}>
              更新
            </Button>
            <Button
              variant="contained"
              size="small"
              startIcon={syncing ? <CircularProgress size={14} color="inherit" /> : <SyncIcon />}
              onClick={handleSync}
              disabled={syncing || !settings?.configured}
            >
              今すぐ同期
            </Button>
          </Stack>
        </Stack>
        <Divider sx={{ mb: 2 }} />

        {logs.length === 0 ? (
          <Typography variant="body2" color="text.secondary" sx={{ py: 2, textAlign: 'center' }}>
            同期ログがありません
          </Typography>
        ) : (
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>開始日時</TableCell>
                  <TableCell>ステータス</TableCell>
                  <TableCell align="right">PJ数</TableCell>
                  <TableCell align="right">チケット数</TableCell>
                  <TableCell align="right">所要時間</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {logs.map((log) => (
                  <TableRow key={log.id}>
                    <TableCell>
                      {new Date(log.executed_at).toLocaleString('ja-JP')}
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={log.status}
                        size="small"
                        color={statusColor[log.status] ?? 'default'}
                      />
                      {log.error_message && (
                        <Typography variant="caption" color="error" display="block">
                          {log.error_message}
                        </Typography>
                      )}
                    </TableCell>
                    <TableCell align="right">{log.projects_synced}</TableCell>
                    <TableCell align="right">{log.issues_synced}</TableCell>
                    <TableCell align="right">
                      {log.duration_seconds != null ? `${log.duration_seconds}秒` : '-'}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Paper>
    </Box>
  )
}
