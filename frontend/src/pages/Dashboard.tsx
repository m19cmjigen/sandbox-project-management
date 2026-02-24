import { useEffect, useState } from 'react'
import {
  Box,
  Divider,
  Grid,
  Stack,
  Typography,
  alpha,
} from '@mui/material'
import {
  ErrorOutline as RedIcon,
  WarningAmber as YellowIcon,
  CheckCircleOutline as GreenIcon,
  FolderOpen as FolderIcon,
} from '@mui/icons-material'
import { getDashboardSummary } from '../api/dashboard'
import { buildDashboardTree } from '../types/dashboard'
import type { DashboardSummary, DashboardOrgNode } from '../types/dashboard'
import HeatmapCard from '../components/HeatmapCard'
import LoadingSpinner from '../components/LoadingSpinner'
import ErrorMessage from '../components/ErrorMessage'

interface SummaryCardProps {
  label: string
  value: number
  sublabel?: string
  subvalue?: number
  color: string
  icon: React.ReactNode
}

function SummaryCard({ label, value, sublabel, subvalue, color, icon }: SummaryCardProps) {
  return (
    <Box
      sx={{
        bgcolor: 'background.paper',
        borderRadius: 3,
        p: 3,
        border: '1px solid',
        borderColor: 'divider',
        boxShadow: '0 1px 3px 0 rgb(0 0 0 / 0.06)',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        gap: 2,
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <Typography variant="body2" color="text.secondary" fontWeight={500}>
          {label}
        </Typography>
        <Box
          sx={{
            width: 36,
            height: 36,
            borderRadius: 2,
            bgcolor: alpha(color, 0.12),
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color,
            flexShrink: 0,
          }}
        >
          {icon}
        </Box>
      </Box>
      <Box>
        <Typography variant="h3" fontWeight={700} color="text.primary" lineHeight={1}>
          {value}
        </Typography>
        {sublabel && subvalue !== undefined && (
          <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: 'block' }}>
            {sublabel} {subvalue.toLocaleString()}件
          </Typography>
        )}
      </Box>
      {/* Bottom accent bar */}
      <Box sx={{ height: 3, bgcolor: alpha(color, 0.25), borderRadius: 1, mt: 'auto' }}>
        <Box sx={{ height: '100%', width: '60%', bgcolor: color, borderRadius: 1 }} />
      </Box>
    </Box>
  )
}

function OrgHeatmap({ roots }: { roots: DashboardOrgNode[] }) {
  if (roots.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary" sx={{ py: 2 }}>
        組織データがありません
      </Typography>
    )
  }

  return (
    <Stack spacing={3}>
      {roots.map((root) => (
        <Box key={root.id}>
          <Grid container spacing={2} alignItems="stretch">
            <Grid item xs={12} sm={4} md={3}>
              <HeatmapCard node={root} />
            </Grid>
            {root.children.length > 0 && (
              <Grid item xs={12} sm={8} md={9}>
                <Grid container spacing={1.5} sx={{ height: '100%' }}>
                  {root.children.map((child) => (
                    <Grid item xs={6} sm={4} md={3} key={child.id}>
                      <HeatmapCard node={child} isChild />
                    </Grid>
                  ))}
                </Grid>
              </Grid>
            )}
          </Grid>
          <Divider sx={{ mt: 3, borderColor: 'divider' }} />
        </Box>
      ))}
    </Stack>
  )
}

export default function Dashboard() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const load = async () => {
      try {
        const data = await getDashboardSummary()
        setSummary(data)
      } catch {
        setError('ダッシュボードデータの取得に失敗しました。')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const handleRetry = () => {
    setError(null)
    setLoading(true)
    getDashboardSummary()
      .then(setSummary)
      .catch(() => setError('ダッシュボードデータの取得に失敗しました。'))
      .finally(() => setLoading(false))
  }

  const roots = summary ? buildDashboardTree(summary.organizations) : []

  return (
    <Box>
      {/* Page header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" fontWeight={700} color="text.primary">
          ダッシュボード
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
          全社プロジェクトの進捗状況
        </Typography>
      </Box>

      {error && <ErrorMessage message={error} onRetry={handleRetry} />}

      {loading ? (
        <LoadingSpinner minHeight={320} />
      ) : summary && (
        <>
          {/* Summary cards */}
          <Grid container spacing={2.5} sx={{ mb: 5 }}>
            <Grid item xs={6} sm={3}>
              <SummaryCard
                label="総プロジェクト"
                value={summary.total_projects}
                sublabel="チケット"
                subvalue={summary.total_issues}
                color="#6366f1"
                icon={<FolderIcon fontSize="small" />}
              />
            </Grid>
            <Grid item xs={6} sm={3}>
              <SummaryCard
                label="遅延プロジェクト"
                value={summary.red_projects}
                sublabel="遅延チケット"
                subvalue={summary.red_issues}
                color="#ef4444"
                icon={<RedIcon fontSize="small" />}
              />
            </Grid>
            <Grid item xs={6} sm={3}>
              <SummaryCard
                label="注意プロジェクト"
                value={summary.yellow_projects}
                sublabel="注意チケット"
                subvalue={summary.yellow_issues}
                color="#f59e0b"
                icon={<YellowIcon fontSize="small" />}
              />
            </Grid>
            <Grid item xs={6} sm={3}>
              <SummaryCard
                label="正常プロジェクト"
                value={summary.green_projects}
                sublabel="正常チケット"
                subvalue={summary.green_issues}
                color="#10b981"
                icon={<GreenIcon fontSize="small" />}
              />
            </Grid>
          </Grid>

          {/* Heatmap section */}
          <Box sx={{ mb: 2 }}>
            <Typography variant="h6" fontWeight={600} color="text.primary">
              組織別プロジェクト状況
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
              カードをクリックすると該当組織のプロジェクト一覧を表示します
            </Typography>
          </Box>
          <OrgHeatmap roots={roots} />
        </>
      )}
    </Box>
  )
}
