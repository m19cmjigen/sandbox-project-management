import { useEffect, useState } from 'react'
import {
  Card,
  CardContent,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  Box,
  Alert,
} from '@mui/material'
import { syncService, SyncLog } from '@/services/syncService'
import Loading from '@/components/Loading'

export default function SyncHistory() {
  const [logs, setLogs] = useState<SyncLog[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchLogs()
  }, [])

  const fetchLogs = async () => {
    try {
      setLoading(true)
      const data = await syncService.getSyncLogs()
      setLogs(data)
      setError(null)
    } catch (err) {
      setError('同期ログの取得に失敗しました')
      console.error(err)
    } finally {
      setLoading(false)
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

  const formatDateTime = (dateString: string) => {
    return new Date(dateString).toLocaleString('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  }

  const calculateDuration = (startedAt: string, completedAt: string | null) => {
    if (!completedAt) return '-'

    const start = new Date(startedAt).getTime()
    const end = new Date(completedAt).getTime()
    const durationMs = end - start

    const seconds = Math.floor(durationMs / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)

    if (hours > 0) {
      return `${hours}時間${minutes % 60}分`
    } else if (minutes > 0) {
      return `${minutes}分${seconds % 60}秒`
    } else {
      return `${seconds}秒`
    }
  }

  if (loading) return <Loading />

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          同期履歴
        </Typography>

        {error && (
          <Alert severity="warning" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {logs.length === 0 ? (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography variant="body2" color="textSecondary">
              同期履歴がありません
            </Typography>
          </Box>
        ) : (
          <TableContainer component={Paper}>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>開始時刻</TableCell>
                  <TableCell>ステータス</TableCell>
                  <TableCell align="right">プロジェクト</TableCell>
                  <TableCell align="right">Issue</TableCell>
                  <TableCell align="right">エラー</TableCell>
                  <TableCell>所要時間</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {logs.map((log) => (
                  <TableRow key={log.id} hover>
                    <TableCell>
                      <Typography variant="body2">
                        {formatDateTime(log.started_at)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={getStatusLabel(log.status)}
                        color={getStatusColor(log.status) as any}
                        size="small"
                      />
                    </TableCell>
                    <TableCell align="right">{log.projects_synced}</TableCell>
                    <TableCell align="right">{log.issues_synced}</TableCell>
                    <TableCell align="right">
                      {log.error_count > 0 ? (
                        <Chip
                          label={log.error_count}
                          color="error"
                          size="small"
                        />
                      ) : (
                        '-'
                      )}
                    </TableCell>
                    <TableCell>
                      {calculateDuration(log.started_at, log.completed_at)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </CardContent>
    </Card>
  )
}
