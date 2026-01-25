import { useState } from 'react'
import {
  Card,
  CardContent,
  Typography,
  Button,
  Box,
  Alert,
  CircularProgress,
  LinearProgress,
  Chip,
} from '@mui/material'
import { Sync as SyncIcon } from '@mui/icons-material'
import { syncService, SyncLog } from '@/services/syncService'

export default function SyncTrigger() {
  const [syncing, setSyncing] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [lastSync, setLastSync] = useState<SyncLog | null>(null)

  const handleSync = async () => {
    setError(null)
    setSuccess(null)
    setSyncing(true)

    try {
      const result = await syncService.triggerSync(1) // デフォルト組織ID: 1
      setLastSync(result.sync_log)

      if (result.sync_log.status === 'COMPLETED') {
        setSuccess(
          `同期が完了しました。プロジェクト: ${result.sync_log.projects_synced}件、Issue: ${result.sync_log.issues_synced}件`
        )
      } else if (result.sync_log.status === 'COMPLETED_WITH_ERRORS') {
        setSuccess(
          `同期が完了しました（エラー: ${result.sync_log.error_count}件）。プロジェクト: ${result.sync_log.projects_synced}件、Issue: ${result.sync_log.issues_synced}件`
        )
      } else if (result.sync_log.status === 'FAILED') {
        setError(
          result.sync_log.error_message || '同期に失敗しました'
        )
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Jira同期の実行に失敗しました'
      )
    } finally {
      setSyncing(false)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'COMPLETED':
        return 'success'
      case 'COMPLETED_WITH_ERRORS':
        return 'warning'
      case 'FAILED':
        return 'error'
      case 'RUNNING':
        return 'info'
      default:
        return 'default'
    }
  }

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'COMPLETED':
        return '完了'
      case 'COMPLETED_WITH_ERRORS':
        return '完了（エラーあり）'
      case 'FAILED':
        return '失敗'
      case 'RUNNING':
        return '実行中'
      default:
        return status
    }
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Jira同期
        </Typography>
        <Typography variant="body2" color="textSecondary" sx={{ mb: 3 }}>
          Jira Cloudから最新のプロジェクトとIssueデータを取得します
        </Typography>

        {syncing && <LinearProgress sx={{ mb: 2 }} />}

        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        {success && (
          <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess(null)}>
            {success}
          </Alert>
        )}

        {lastSync && (
          <Box sx={{ mb: 2, p: 2, bgcolor: 'background.default', borderRadius: 1 }}>
            <Typography variant="subtitle2" gutterBottom>
              前回の同期結果
            </Typography>
            <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mt: 1 }}>
              <Chip
                label={getStatusLabel(lastSync.status)}
                color={getStatusColor(lastSync.status) as any}
                size="small"
              />
              <Chip
                label={`プロジェクト: ${lastSync.projects_synced}件`}
                size="small"
                variant="outlined"
              />
              <Chip
                label={`Issue: ${lastSync.issues_synced}件`}
                size="small"
                variant="outlined"
              />
              {lastSync.error_count > 0 && (
                <Chip
                  label={`エラー: ${lastSync.error_count}件`}
                  color="error"
                  size="small"
                  variant="outlined"
                />
              )}
            </Box>
            <Typography variant="caption" color="textSecondary" sx={{ mt: 1, display: 'block' }}>
              開始: {new Date(lastSync.started_at).toLocaleString('ja-JP')}
            </Typography>
          </Box>
        )}

        <Button
          variant="contained"
          startIcon={syncing ? <CircularProgress size={20} /> : <SyncIcon />}
          onClick={handleSync}
          disabled={syncing}
          fullWidth
        >
          {syncing ? '同期中...' : '今すぐ同期'}
        </Button>

        <Typography variant="caption" color="textSecondary" sx={{ mt: 2, display: 'block' }}>
          注意: 大量のプロジェクトがある場合、同期に数分かかることがあります
        </Typography>
      </CardContent>
    </Card>
  )
}
