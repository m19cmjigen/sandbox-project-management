import { useEffect, useState } from 'react'
import {
  Alert,
  Box,
  Card,
  CardContent,
  CircularProgress,
  Divider,
  InputAdornment,
  Stack,
  TextField,
  Tooltip,
  Typography,
} from '@mui/material'
import {
  Search as SearchIcon,
  FiberManualRecord as DotIcon,
  OpenInNew as ExternalLinkIcon,
} from '@mui/icons-material'
import OrganizationTree from '../components/OrganizationTree'
import { getOrganizations } from '../api/organizations'
import { buildOrganizationTree } from '../types/organization'
import type { OrganizationTreeNode } from '../types/organization'

const STATUS_COLOR = {
  RED: '#f44336',
  YELLOW: '#ff9800',
  GREEN: '#4caf50',
} as const

const STATUS_LABEL = {
  RED: '遅延あり',
  YELLOW: '注意',
  GREEN: '正常',
} as const

function LegendItem({ status }: { status: 'RED' | 'YELLOW' | 'GREEN' }) {
  return (
    <Stack direction="row" alignItems="center" spacing={0.5}>
      <DotIcon sx={{ color: STATUS_COLOR[status], fontSize: 14 }} />
      <Typography variant="caption" color="text.secondary">
        {STATUS_LABEL[status]}
      </Typography>
    </Stack>
  )
}

function SummaryCard({ nodes }: { nodes: OrganizationTreeNode[] }) {
  const total = nodes.reduce((s, n) => s + n.subtree_total, 0)
  const red = nodes.reduce((s, n) => s + n.subtree_red, 0)
  const yellow = nodes.reduce((s, n) => s + n.subtree_yellow, 0)
  const green = nodes.reduce((s, n) => s + n.subtree_green, 0)
  const orgsCount = nodes.reduce((s, n) => s + 1 + n.children.length, 0)

  return (
    <Card variant="outlined" sx={{ mb: 2 }}>
      <CardContent sx={{ py: 1.5, '&:last-child': { pb: 1.5 } }}>
        <Stack direction="row" spacing={3} flexWrap="wrap">
          <Box>
            <Typography variant="caption" color="text.secondary">組織数</Typography>
            <Typography variant="h6" fontWeight="bold">{orgsCount}</Typography>
          </Box>
          <Divider orientation="vertical" flexItem />
          <Box>
            <Typography variant="caption" color="text.secondary">総プロジェクト</Typography>
            <Typography variant="h6" fontWeight="bold">{total}</Typography>
          </Box>
          <Divider orientation="vertical" flexItem />
          <Box>
            <Typography variant="caption" color="error.main">遅延プロジェクト</Typography>
            <Typography variant="h6" fontWeight="bold" color="error.main">{red}</Typography>
          </Box>
          <Divider orientation="vertical" flexItem />
          <Box>
            <Typography variant="caption" color="warning.main">注意プロジェクト</Typography>
            <Typography variant="h6" fontWeight="bold" color="warning.main">{yellow}</Typography>
          </Box>
          <Divider orientation="vertical" flexItem />
          <Box>
            <Typography variant="caption" color="success.main">正常プロジェクト</Typography>
            <Typography variant="h6" fontWeight="bold" color="success.main">{green}</Typography>
          </Box>
        </Stack>
      </CardContent>
    </Card>
  )
}

export default function Organizations() {
  const [treeNodes, setTreeNodes] = useState<OrganizationTreeNode[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')

  useEffect(() => {
    const fetchOrgs = async () => {
      setLoading(true)
      setError(null)
      try {
        const orgs = await getOrganizations()
        setTreeNodes(buildOrganizationTree(orgs))
      } catch {
        setError('組織情報の取得に失敗しました。')
      } finally {
        setLoading(false)
      }
    }
    fetchOrgs()
  }, [])

  return (
    <Box>
      {/* Page title */}
      <Typography variant="h4" gutterBottom>
        組織階層
      </Typography>

      {/* Summary cards */}
      {!loading && treeNodes.length > 0 && <SummaryCard nodes={treeNodes} />}

      {/* Search box + legend */}
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ xs: 'flex-start', sm: 'center' }} sx={{ mb: 2 }}>
        <TextField
          size="small"
          placeholder="組織名で検索..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon fontSize="small" />
              </InputAdornment>
            ),
          }}
          sx={{ width: { xs: '100%', sm: 280 } }}
        />
        <Box sx={{ flexGrow: 1 }} />
        <Stack direction="row" spacing={1.5} alignItems="center">
          <Typography variant="caption" color="text.secondary">
            色の意味:
          </Typography>
          <LegendItem status="RED" />
          <LegendItem status="YELLOW" />
          <LegendItem status="GREEN" />
          <Tooltip title="組織名をクリックするとプロジェクト一覧へ遷移します">
            <ExternalLinkIcon fontSize="small" sx={{ color: 'text.disabled', cursor: 'help' }} />
          </Tooltip>
        </Stack>
      </Stack>

      {/* Loading state */}
      {loading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
          <CircularProgress />
        </Box>
      )}

      {/* Error state */}
      {error && <Alert severity="error">{error}</Alert>}

      {/* Tree */}
      {!loading && !error && (
        <Card variant="outlined">
          <CardContent sx={{ p: 1, '&:last-child': { pb: 1 } }}>
            <OrganizationTree nodes={treeNodes} searchQuery={searchQuery} />
          </CardContent>
        </Card>
      )}
    </Box>
  )
}
