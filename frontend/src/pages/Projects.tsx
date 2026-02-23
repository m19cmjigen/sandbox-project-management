import { useCallback, useEffect, useState } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Pagination,
  Select,
  Stack,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from '@mui/material'
import { ArrowBack as BackIcon } from '@mui/icons-material'
import type { SelectChangeEvent } from '@mui/material'
import ProjectCard from '../components/ProjectCard'
import { getProjects } from '../api/projects'
import type { DelayFilter, PaginationMeta, Project, SortOption } from '../types/project'

const PER_PAGE = 12

export default function Projects() {
  const [searchParams, setSearchParams] = useSearchParams()
  const navigate = useNavigate()

  // Read organization_id from URL params (set when drilling down from org tree)
  const orgIdParam = searchParams.get('organization_id')
  const orgId = orgIdParam ? parseInt(orgIdParam, 10) : undefined

  const [projects, setProjects] = useState<Project[]>([])
  const [pagination, setPagination] = useState<PaginationMeta | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [page, setPage] = useState(1)
  const [sort, setSort] = useState<SortOption>('name')
  const [delayFilter, setDelayFilter] = useState<DelayFilter>('ALL')

  const fetchProjects = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await getProjects({
        page,
        per_page: PER_PAGE,
        sort,
        delay_status: delayFilter,
        organization_id: orgId,
      })
      setProjects(res.data)
      setPagination(res.pagination)
    } catch {
      setError('プロジェクトの取得に失敗しました。バックエンドの接続を確認してください。')
    } finally {
      setLoading(false)
    }
  }, [page, sort, delayFilter, orgId])

  useEffect(() => {
    fetchProjects()
  }, [fetchProjects])

  // Reset to page 1 when filter/sort changes
  const handleFilterChange = (_: React.MouseEvent<HTMLElement>, value: DelayFilter | null) => {
    if (value !== null) {
      setDelayFilter(value)
      setPage(1)
    }
  }

  const handleSortChange = (e: SelectChangeEvent) => {
    setSort(e.target.value as SortOption)
    setPage(1)
  }

  const handlePageChange = (_: React.ChangeEvent<unknown>, value: number) => {
    setPage(value)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const handleClearOrgFilter = () => {
    setSearchParams({})
    setPage(1)
  }

  return (
    <Box>
      {/* Page title */}
      <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 1 }}>
        {orgId && (
          <Button
            size="small"
            startIcon={<BackIcon />}
            onClick={() => navigate('/organizations')}
            sx={{ mr: 1 }}
          >
            組織一覧へ戻る
          </Button>
        )}
        <Typography variant="h4">
          プロジェクト一覧
        </Typography>
      </Stack>

      {/* Organization filter badge */}
      {orgId && (
        <Alert
          severity="info"
          sx={{ mb: 2 }}
          action={
            <Button color="inherit" size="small" onClick={handleClearOrgFilter}>
              絞り込み解除
            </Button>
          }
        >
          組織ID {orgId} でフィルタ中
        </Alert>
      )}

      {/* Toolbar: filter + sort */}
      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        spacing={2}
        alignItems={{ xs: 'flex-start', sm: 'center' }}
        sx={{ mb: 3 }}
      >
        <ToggleButtonGroup
          value={delayFilter}
          exclusive
          onChange={handleFilterChange}
          size="small"
          aria-label="ステータスフィルタ"
        >
          <ToggleButton value="ALL">すべて</ToggleButton>
          <ToggleButton value="RED" sx={{ color: 'error.main', '&.Mui-selected': { bgcolor: 'error.light', color: 'error.contrastText' } }}>
            遅延あり
          </ToggleButton>
          <ToggleButton value="YELLOW" sx={{ color: 'warning.main', '&.Mui-selected': { bgcolor: 'warning.light', color: 'warning.contrastText' } }}>
            注意
          </ToggleButton>
          <ToggleButton value="GREEN" sx={{ color: 'success.main', '&.Mui-selected': { bgcolor: 'success.light', color: 'success.contrastText' } }}>
            正常
          </ToggleButton>
        </ToggleButtonGroup>

        <Box sx={{ flexGrow: 1 }} />

        <FormControl size="small" sx={{ minWidth: 160 }}>
          <InputLabel>並び順</InputLabel>
          <Select value={sort} label="並び順" onChange={handleSortChange}>
            <MenuItem value="name">名前順 (A→Z)</MenuItem>
            <MenuItem value="name_desc">名前順 (Z→A)</MenuItem>
            <MenuItem value="delay_count">遅延数の多い順</MenuItem>
          </Select>
        </FormControl>
      </Stack>

      {/* Summary bar */}
      {pagination && !loading && (
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          {pagination.total}件中 {(page - 1) * PER_PAGE + 1}–{Math.min(page * PER_PAGE, pagination.total)}件を表示
        </Typography>
      )}

      {/* Error state */}
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {/* Loading state */}
      {loading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
          <CircularProgress />
        </Box>
      )}

      {/* Empty state */}
      {!loading && !error && projects.length === 0 && (
        <Alert severity="info">
          該当するプロジェクトがありません。
        </Alert>
      )}

      {/* Project grid */}
      {!loading && projects.length > 0 && (
        <Grid container spacing={2}>
          {projects.map((project) => (
            <Grid item key={project.id} xs={12} sm={6} md={4} lg={3}>
              <ProjectCard project={project} />
            </Grid>
          ))}
        </Grid>
      )}

      {/* Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
          <Pagination
            count={pagination.total_pages}
            page={page}
            onChange={handlePageChange}
            color="primary"
            showFirstButton
            showLastButton
          />
        </Box>
      )}
    </Box>
  )
}
