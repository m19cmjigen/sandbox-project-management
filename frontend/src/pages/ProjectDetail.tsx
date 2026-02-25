import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import {
  Box,
  Button,
  Chip,
  Divider,
  Grid,
  IconButton,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tooltip,
  Typography,
  alpha,
} from '@mui/material'
import {
  ArrowBack as BackIcon,
  ErrorOutline as RedIcon,
  WarningAmber as YellowIcon,
  CheckCircleOutline as GreenIcon,
  Assignment as IssuesIcon,
  OpenInNew as ExternalLinkIcon,
} from '@mui/icons-material'
import { getProjectSummary } from '../api/dashboard'
import type { ProjectSummaryResponse } from '../types/dashboard'
import type { Issue } from '../types/issue'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'
import StatusBadge from '../components/StatusBadge'
import { formatDate, isDueDateOverdue } from '../utils/dateUtils'

const JIRA_BASE_URL = import.meta.env.VITE_JIRA_BASE_URL || ''

interface SummaryChipProps {
  label: string
  value: number
  color: string
  icon: React.ReactNode
}

function SummaryChip({ label, value, color, icon }: SummaryChipProps) {
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: 0.5,
        px: 2.5,
        py: 1.5,
        borderRadius: 2,
        bgcolor: alpha(color, 0.08),
        border: '1px solid',
        borderColor: alpha(color, 0.2),
        minWidth: 80,
      }}
    >
      <Box sx={{ color, display: 'flex', alignItems: 'center' }}>{icon}</Box>
      <Typography variant="h5" fontWeight={700} color={color} lineHeight={1}>
        {value}
      </Typography>
      <Typography variant="caption" color="text.secondary" fontSize="0.6875rem">
        {label}
      </Typography>
    </Box>
  )
}

function DelayedIssueRow({ issue }: { issue: Issue }) {
  const overdue = isDueDateOverdue(issue.due_date, issue.status_category)
  const jiraUrl = JIRA_BASE_URL ? `${JIRA_BASE_URL}/browse/${issue.jira_issue_key}` : null

  return (
    <TableRow hover sx={{ '&:last-child td': { border: 0 } }}>
      <TableCell sx={{ width: 80 }}>
        <StatusBadge status={issue.delay_status} />
      </TableCell>
      <TableCell sx={{ width: 120 }}>
        <Typography variant="body2" sx={{ fontFamily: 'monospace', fontWeight: 600 }}>
          {issue.jira_issue_key}
        </Typography>
      </TableCell>
      <TableCell>
        <Typography
          variant="body2"
          sx={{
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            maxWidth: 400,
          }}
        >
          {issue.summary}
        </Typography>
      </TableCell>
      <TableCell sx={{ width: 110 }}>
        <Typography
          variant="body2"
          color={overdue ? 'error.main' : 'text.primary'}
          fontWeight={overdue ? 600 : 400}
        >
          {formatDate(issue.due_date)}
        </Typography>
      </TableCell>
      <TableCell sx={{ width: 120 }}>
        <Typography variant="body2" color="text.secondary">
          {issue.assignee_name ?? '—'}
        </Typography>
      </TableCell>
      <TableCell sx={{ width: 48 }}>
        {jiraUrl && (
          <Tooltip title="Jiraで開く">
            <IconButton
              size="small"
              component="a"
              href={jiraUrl}
              target="_blank"
              rel="noopener noreferrer"
              sx={{ p: 0.5, color: 'text.secondary', '&:hover': { color: 'primary.main' } }}
            >
              <ExternalLinkIcon sx={{ fontSize: 15 }} />
            </IconButton>
          </Tooltip>
        )}
      </TableCell>
    </TableRow>
  )
}

export default function ProjectDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [data, setData] = useState<ProjectSummaryResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    const projectId = parseInt(id, 10)
    if (isNaN(projectId)) {
      setError('不正なプロジェクトIDです。')
      setLoading(false)
      return
    }
    setLoading(true)
    setError(null)
    getProjectSummary(projectId)
      .then(setData)
      .catch(() => setError('プロジェクト情報の取得に失敗しました。'))
      .finally(() => setLoading(false))
  }, [id])

  const handleRetry = () => {
    if (!id) return
    const projectId = parseInt(id, 10)
    setError(null)
    setLoading(true)
    getProjectSummary(projectId)
      .then(setData)
      .catch(() => setError('プロジェクト情報の取得に失敗しました。'))
      .finally(() => setLoading(false))
  }

  const jiraProjectUrl = data && JIRA_BASE_URL
    ? `${JIRA_BASE_URL}/browse/${data.project.key}`
    : null

  const delayStatusColor: Record<string, string> = {
    RED: '#ef4444',
    YELLOW: '#f59e0b',
    GREEN: '#10b981',
  }

  return (
    <Box>
      {/* Back navigation */}
      <Button
        size="small"
        startIcon={<BackIcon />}
        onClick={() => navigate(-1)}
        sx={{ mb: 2 }}
      >
        戻る
      </Button>

      {loading && <LoadingSpinner minHeight={320} />}
      {error && <ErrorMessage message={error} onRetry={handleRetry} />}

      {!loading && !error && data && (
        <>
          {/* Header */}
          <Box sx={{ mb: 3 }}>
            <Stack direction="row" alignItems="flex-start" justifyContent="space-between" flexWrap="wrap" gap={1}>
              <Box>
                <Stack direction="row" alignItems="center" spacing={1.5} flexWrap="wrap">
                  <Typography
                    variant="caption"
                    sx={{
                      fontFamily: '"JetBrains Mono", "Fira Code", monospace',
                      fontSize: '0.75rem',
                      fontWeight: 700,
                      color: 'text.secondary',
                      bgcolor: 'grey.100',
                      px: 1,
                      py: 0.375,
                      borderRadius: 1,
                    }}
                  >
                    {data.project.key}
                  </Typography>
                  <Chip
                    label={data.project.delay_status === 'RED' ? '遅延あり' : data.project.delay_status === 'YELLOW' ? '注意' : '正常'}
                    size="small"
                    sx={{
                      bgcolor: alpha(delayStatusColor[data.project.delay_status] ?? '#10b981', 0.1),
                      color: delayStatusColor[data.project.delay_status] ?? '#10b981',
                      fontWeight: 700,
                      fontSize: '0.75rem',
                    }}
                  />
                </Stack>
                <Typography variant="h4" fontWeight={700} color="text.primary" sx={{ mt: 0.75 }}>
                  {data.project.name}
                </Typography>
                {data.project.lead_email && (
                  <Typography variant="body2" color="text.secondary" sx={{ mt: 0.25 }}>
                    リード: {data.project.lead_email}
                  </Typography>
                )}
              </Box>

              <Stack direction="row" spacing={1}>
                <Button
                  variant="outlined"
                  size="small"
                  startIcon={<IssuesIcon />}
                  onClick={() => navigate(`/issues?project_id=${data.project.id}`)}
                >
                  全チケット一覧
                </Button>
                {jiraProjectUrl && (
                  <Button
                    variant="outlined"
                    size="small"
                    startIcon={<ExternalLinkIcon />}
                    component="a"
                    href={jiraProjectUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    Jiraで開く
                  </Button>
                )}
              </Stack>
            </Stack>
          </Box>

          {/* Summary chips */}
          <Grid container spacing={1.5} sx={{ mb: 4 }}>
            <Grid item>
              <SummaryChip
                label="遅延"
                value={data.summary.red_count}
                color="#ef4444"
                icon={<RedIcon fontSize="small" />}
              />
            </Grid>
            <Grid item>
              <SummaryChip
                label="注意"
                value={data.summary.yellow_count}
                color="#f59e0b"
                icon={<YellowIcon fontSize="small" />}
              />
            </Grid>
            <Grid item>
              <SummaryChip
                label="正常"
                value={data.summary.green_count}
                color="#10b981"
                icon={<GreenIcon fontSize="small" />}
              />
            </Grid>
            <Grid item>
              <Box
                sx={{
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  gap: 0.5,
                  px: 2.5,
                  py: 1.5,
                  borderRadius: 2,
                  bgcolor: 'grey.50',
                  border: '1px solid',
                  borderColor: 'divider',
                  minWidth: 80,
                }}
              >
                <Typography variant="h5" fontWeight={700} color="text.primary" lineHeight={1}>
                  {data.summary.open_count}
                </Typography>
                <Typography variant="caption" color="text.secondary" fontSize="0.6875rem">
                  未完了
                </Typography>
              </Box>
            </Grid>
            <Grid item>
              <Box
                sx={{
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  gap: 0.5,
                  px: 2.5,
                  py: 1.5,
                  borderRadius: 2,
                  bgcolor: 'grey.50',
                  border: '1px solid',
                  borderColor: 'divider',
                  minWidth: 80,
                }}
              >
                <Typography variant="h5" fontWeight={700} color="text.secondary" lineHeight={1}>
                  {data.summary.total_count}
                </Typography>
                <Typography variant="caption" color="text.secondary" fontSize="0.6875rem">
                  合計
                </Typography>
              </Box>
            </Grid>
          </Grid>

          {/* Delayed issues section */}
          <Box sx={{ mb: 2 }}>
            <Stack direction="row" alignItems="center" justifyContent="space-between">
              <Box>
                <Typography variant="h6" fontWeight={600}>
                  遅延・注意チケット
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mt: 0.25 }}>
                  期限切れ・期限間近のチケット（最大20件）
                </Typography>
              </Box>
              {data.summary.red_count + data.summary.yellow_count > 20 && (
                <Button
                  size="small"
                  variant="text"
                  onClick={() => navigate(`/issues?project_id=${data.project.id}&delay_status=RED`)}
                >
                  全件表示
                </Button>
              )}
            </Stack>
          </Box>

          <Divider sx={{ mb: 2 }} />

          {data.delayed_issues.length === 0 ? (
            <Paper
              variant="outlined"
              sx={{ p: 4, textAlign: 'center' }}
            >
              <GreenIcon sx={{ fontSize: 40, color: '#10b981', mb: 1 }} />
              <Typography variant="body1" color="text.secondary">
                遅延・注意チケットはありません
              </Typography>
            </Paper>
          ) : (
            <TableContainer component={Paper} variant="outlined">
              <Table size="small" stickyHeader>
                <TableHead>
                  <TableRow>
                    <TableCell sx={{ fontWeight: 700 }}>ステータス</TableCell>
                    <TableCell sx={{ fontWeight: 700 }}>キー</TableCell>
                    <TableCell sx={{ fontWeight: 700 }}>サマリ</TableCell>
                    <TableCell sx={{ fontWeight: 700 }}>期限</TableCell>
                    <TableCell sx={{ fontWeight: 700 }}>担当者</TableCell>
                    <TableCell />
                  </TableRow>
                </TableHead>
                <TableBody>
                  {data.delayed_issues.map((issue) => (
                    <DelayedIssueRow key={issue.id} issue={issue} />
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}

          {/* Footer link */}
          <Box sx={{ mt: 3, textAlign: 'center' }}>
            <Button
              variant="contained"
              startIcon={<IssuesIcon />}
              onClick={() => navigate(`/issues?project_id=${data.project.id}`)}
            >
              全チケット一覧を表示（{data.summary.total_count}件）
            </Button>
          </Box>
        </>
      )}
    </Box>
  )
}
