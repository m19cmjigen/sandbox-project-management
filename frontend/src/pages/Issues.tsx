import { useCallback, useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import {
  Box,
  Button,
  Checkbox,
  FormControl,
  FormControlLabel,
  InputLabel,
  MenuItem,
  Pagination,
  Paper,
  Select,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TableSortLabel,
  ToggleButton,
  ToggleButtonGroup,
  Tooltip,
  Typography,
} from '@mui/material'
import {
  FileDownload as DownloadIcon,
  OpenInNew as ExternalLinkIcon,
} from '@mui/icons-material'
import type { SelectChangeEvent } from '@mui/material'
import { getIssues } from '../api/issues'
import type { Issue, IssueSortKey, IssueListParams, SortOrder } from '../types/issue'
import type { DelayStatus, PaginationMeta } from '../types/project'
import { formatDate, isDueDateOverdue } from '../utils/dateUtils'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'
import StatusBadge from '../components/StatusBadge'

const JIRA_BASE_URL = import.meta.env.VITE_JIRA_BASE_URL || ''
const PER_PAGE = 25

function exportCSV(issues: Issue[]) {
  const headers = ['キー', 'サマリ', 'プロジェクト', 'ステータス', '期限', '担当者', '遅延ステータス', '優先度', 'タイプ']
  const rows = issues.map((i) => [
    i.jira_issue_key,
    `"${i.summary.replace(/"/g, '""')}"`,
    i.project_key,
    i.status,
    i.due_date ?? '',
    i.assignee_name ?? '',
    i.delay_status,
    i.priority ?? '',
    i.issue_type ?? '',
  ])
  const csv = [headers.join(','), ...rows.map((r) => r.join(','))].join('\n')
  const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `issues_${new Date().toISOString().slice(0, 10)}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

export default function Issues() {
  const [searchParams] = useSearchParams()
  const projectIdParam = searchParams.get('project_id')
  const initProjectId = projectIdParam ? parseInt(projectIdParam, 10) : undefined

  const [issues, setIssues] = useState<Issue[]>([])
  const [pagination, setPagination] = useState<PaginationMeta | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Filter state
  const [delayFilter, setDelayFilter] = useState<DelayStatus | 'ALL'>('ALL')
  const [noDueDate, setNoDueDate] = useState(false)
  const [statusCategory, setStatusCategory] = useState('')
  const [projectId, setProjectId] = useState<number | undefined>(initProjectId)
  const [assigneeName, setAssigneeName] = useState('')

  // Sort state
  const [sortKey, setSortKey] = useState<IssueSortKey>('due_date')
  const [sortOrder, setSortOrder] = useState<SortOrder>('asc')

  // Pagination
  const [page, setPage] = useState(1)

  // Options derived from loaded data
  const [projectOptions, setProjectOptions] = useState<{ id: number; key: string; name: string }[]>([])
  const [assigneeOptions, setAssigneeOptions] = useState<string[]>([])

  const buildParams = useCallback((): IssueListParams => ({
    page,
    per_page: PER_PAGE,
    sort: sortKey,
    order: sortOrder,
    delay_status: delayFilter,
    no_due_date: noDueDate || undefined,
    status_category: statusCategory || undefined,
    project_id: projectId,
    assignee_name: assigneeName || undefined,
  }), [page, sortKey, sortOrder, delayFilter, noDueDate, statusCategory, projectId, assigneeName])

  const fetchIssues = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await getIssues(buildParams())
      setIssues(res.data)
      setPagination(res.pagination)

      // Build option lists from first full load
      if (page === 1 && !delayFilter && !noDueDate && !statusCategory && !projectId && !assigneeName) {
        const uniqueProjects = Array.from(
          new Map(res.data.map((i) => [i.project_id, { id: i.project_id, key: i.project_key, name: i.project_name }])).values()
        )
        const uniqueAssignees = [...new Set(res.data.map((i) => i.assignee_name).filter(Boolean))] as string[]
        setProjectOptions(uniqueProjects)
        setAssigneeOptions(uniqueAssignees)
      }
    } catch {
      setError('チケットの取得に失敗しました。')
    } finally {
      setLoading(false)
    }
  }, [buildParams, page, delayFilter, noDueDate, statusCategory, projectId, assigneeName])

  // Load option lists on mount
  useEffect(() => {
    const loadOptions = async () => {
      try {
        const res = await getIssues({ per_page: 100 })
        const uniqueProjects = Array.from(
          new Map(res.data.map((i) => [i.project_id, { id: i.project_id, key: i.project_key, name: i.project_name }])).values()
        )
        const uniqueAssignees = [...new Set(res.data.map((i) => i.assignee_name).filter(Boolean))] as string[]
        setProjectOptions(uniqueProjects)
        setAssigneeOptions(uniqueAssignees)
      } catch { /* ignore */ }
    }
    loadOptions()
  }, [])

  useEffect(() => {
    fetchIssues()
  }, [fetchIssues])

  const handleSort = (key: IssueSortKey) => {
    if (key === sortKey) {
      setSortOrder((o) => (o === 'asc' ? 'desc' : 'asc'))
    } else {
      setSortKey(key)
      setSortOrder('asc')
    }
    setPage(1)
  }

  const handleFilterChange = () => setPage(1)

  const handleCSVExport = async () => {
    try {
      const res = await getIssues({ ...buildParams(), page: 1, per_page: 100 })
      exportCSV(res.data)
    } catch { /* ignore */ }
  }

  return (
    <Box>
      <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 2 }}>
        <Typography variant="h4">チケット一覧</Typography>
        <Button
          variant="outlined"
          size="small"
          startIcon={<DownloadIcon />}
          onClick={handleCSVExport}
          disabled={loading}
        >
          CSV出力
        </Button>
      </Stack>

      {/* Filter toolbar */}
      <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
        <Stack spacing={2}>
          {/* Row 1: delay status */}
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }}>
            <Typography variant="body2" color="text.secondary" sx={{ minWidth: 80 }}>
              遅延ステータス
            </Typography>
            <ToggleButtonGroup
              value={delayFilter}
              exclusive
              onChange={(_, v) => { if (v !== null) { setDelayFilter(v); handleFilterChange() } }}
              size="small"
            >
              <ToggleButton value="ALL">すべて</ToggleButton>
              <ToggleButton value="RED" sx={{ color: 'error.main', '&.Mui-selected': { bgcolor: 'error.light', color: 'error.contrastText' } }}>遅延</ToggleButton>
              <ToggleButton value="YELLOW" sx={{ color: 'warning.main', '&.Mui-selected': { bgcolor: 'warning.light', color: 'warning.contrastText' } }}>注意</ToggleButton>
              <ToggleButton value="GREEN" sx={{ color: 'success.main', '&.Mui-selected': { bgcolor: 'success.light', color: 'success.contrastText' } }}>正常</ToggleButton>
            </ToggleButtonGroup>
            <FormControlLabel
              control={
                <Checkbox
                  checked={noDueDate}
                  onChange={(e) => { setNoDueDate(e.target.checked); handleFilterChange() }}
                  size="small"
                />
              }
              label={<Typography variant="body2">期限未設定のみ</Typography>}
            />
          </Stack>

          {/* Row 2: selects */}
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <FormControl size="small" sx={{ minWidth: 180 }}>
              <InputLabel>プロジェクト</InputLabel>
              <Select
                value={projectId !== undefined ? String(projectId) : ''}
                label="プロジェクト"
                onChange={(e: SelectChangeEvent) => {
                  setProjectId(e.target.value ? Number(e.target.value) : undefined)
                  handleFilterChange()
                }}
              >
                <MenuItem value="">すべて</MenuItem>
                {projectOptions.map((p) => (
                  <MenuItem key={p.id} value={p.id}>{p.key} - {p.name}</MenuItem>
                ))}
              </Select>
            </FormControl>

            <FormControl size="small" sx={{ minWidth: 160 }}>
              <InputLabel>ステータス</InputLabel>
              <Select
                value={statusCategory}
                label="ステータス"
                onChange={(e: SelectChangeEvent) => { setStatusCategory(e.target.value); handleFilterChange() }}
              >
                <MenuItem value="">すべて</MenuItem>
                <MenuItem value="To Do">To Do</MenuItem>
                <MenuItem value="In Progress">In Progress</MenuItem>
                <MenuItem value="Done">Done</MenuItem>
              </Select>
            </FormControl>

            <FormControl size="small" sx={{ minWidth: 140 }}>
              <InputLabel>担当者</InputLabel>
              <Select
                value={assigneeName}
                label="担当者"
                onChange={(e: SelectChangeEvent) => { setAssigneeName(e.target.value); handleFilterChange() }}
              >
                <MenuItem value="">すべて</MenuItem>
                {assigneeOptions.map((a) => (
                  <MenuItem key={a} value={a}>{a}</MenuItem>
                ))}
              </Select>
            </FormControl>
          </Stack>
        </Stack>
      </Paper>

      {/* Summary */}
      {pagination && !loading && (
        <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
          {pagination.total}件中 {(page - 1) * PER_PAGE + 1}–{Math.min(page * PER_PAGE, pagination.total)}件を表示
        </Typography>
      )}

      {error && <ErrorMessage message={error} onRetry={fetchIssues} />}

      {loading ? (
        <LoadingSpinner minHeight={320} />
      ) : (
        <>
          <TableContainer component={Paper} variant="outlined">
            <Table size="small" stickyHeader>
              <TableHead>
                <TableRow>
                  <TableCell sx={{ fontWeight: 'bold', minWidth: 80 }}>ステータス</TableCell>
                  <TableCell sx={{ fontWeight: 'bold', minWidth: 120 }}>
                    <TableSortLabel
                      active={sortKey === 'jira_issue_key'}
                      direction={sortKey === 'jira_issue_key' ? sortOrder : 'asc'}
                      onClick={() => handleSort('jira_issue_key')}
                    >
                      キー
                    </TableSortLabel>
                  </TableCell>
                  <TableCell sx={{ fontWeight: 'bold', minWidth: 240 }}>サマリ</TableCell>
                  <TableCell sx={{ fontWeight: 'bold', minWidth: 120 }}>プロジェクト</TableCell>
                  <TableCell sx={{ fontWeight: 'bold', minWidth: 120 }}>ステータス</TableCell>
                  <TableCell sx={{ fontWeight: 'bold', minWidth: 110 }}>
                    <TableSortLabel
                      active={sortKey === 'due_date'}
                      direction={sortKey === 'due_date' ? sortOrder : 'asc'}
                      onClick={() => handleSort('due_date')}
                    >
                      期限
                    </TableSortLabel>
                  </TableCell>
                  <TableCell sx={{ fontWeight: 'bold', minWidth: 100 }}>担当者</TableCell>
                  <TableCell sx={{ fontWeight: 'bold', minWidth: 80 }}>優先度</TableCell>
                  <TableCell sx={{ fontWeight: 'bold', width: 48 }} />
                </TableRow>
              </TableHead>
              <TableBody>
                {issues.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={9} align="center" sx={{ py: 4, color: 'text.secondary' }}>
                      該当するチケットがありません
                    </TableCell>
                  </TableRow>
                ) : (
                  issues.map((issue) => {
                    const overdue = isDueDateOverdue(issue.due_date, issue.status_category)
                    const jiraUrl = JIRA_BASE_URL ? `${JIRA_BASE_URL}/browse/${issue.jira_issue_key}` : null
                    return (
                      <TableRow
                        key={issue.id}
                        hover
                        sx={{ '&:last-child td': { border: 0 } }}
                      >
                        <TableCell>
                          <StatusBadge status={issue.delay_status} />
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2" sx={{ fontFamily: 'monospace', fontWeight: 'bold' }}>
                            {issue.jira_issue_key}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2" sx={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', maxWidth: 320 }}>
                            {issue.summary}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                            {issue.project_key}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2">{issue.status}</Typography>
                        </TableCell>
                        <TableCell>
                          <Typography
                            variant="body2"
                            color={overdue ? 'error.main' : 'text.primary'}
                            fontWeight={overdue ? 'bold' : 'normal'}
                          >
                            {formatDate(issue.due_date)}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2">{issue.assignee_name ?? '—'}</Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2">{issue.priority ?? '—'}</Typography>
                        </TableCell>
                        <TableCell>
                          {jiraUrl && (
                            <Tooltip title="Jiraで開く">
                              <ExternalLinkIcon
                                fontSize="small"
                                component="a"
                                href={jiraUrl}
                                target="_blank"
                                rel="noopener noreferrer"
                                sx={{ color: 'primary.main', cursor: 'pointer', display: 'block' }}
                              />
                            </Tooltip>
                          )}
                        </TableCell>
                      </TableRow>
                    )
                  })
                )}
              </TableBody>
            </Table>
          </TableContainer>

          {pagination && pagination.total_pages > 1 && (
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
              <Pagination
                count={pagination.total_pages}
                page={page}
                onChange={(_, v) => { setPage(v); window.scrollTo({ top: 0, behavior: 'smooth' }) }}
                color="primary"
                showFirstButton
                showLastButton
              />
            </Box>
          )}
        </>
      )}
    </Box>
  )
}
