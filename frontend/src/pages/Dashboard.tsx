import { useEffect, useState } from 'react'
import {
  Alert,
  Box,
  Card,
  CardContent,
  CircularProgress,
  Divider,
  Grid,
  Stack,
  Typography,
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
    <Card variant="outlined" sx={{ height: '100%', borderTop: `4px solid ${color}` }}>
      <CardContent>
        <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 1 }}>
          <Box sx={{ color }}>{icon}</Box>
          <Typography variant="body2" color="text.secondary">{label}</Typography>
        </Stack>
        <Typography variant="h4" fontWeight="bold" color={color}>{value}</Typography>
        {sublabel && subvalue !== undefined && (
          <Typography variant="caption" color="text.secondary">
            {sublabel}: {subvalue}件
          </Typography>
        )}
      </CardContent>
    </Card>
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
          {/* Root org (level 0) */}
          <Grid container spacing={2} alignItems="stretch">
            {/* Root card: takes up fixed width */}
            <Grid item xs={12} sm={4} md={3}>
              <HeatmapCard node={root} />
            </Grid>

            {/* Child orgs */}
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
          <Divider sx={{ mt: 3 }} />
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

  const roots = summary ? buildDashboardTree(summary.organizations) : []

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 3 }}>ダッシュボード</Typography>

      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
          <CircularProgress />
        </Box>
      ) : summary && (
        <>
          {/* Summary cards */}
          <Grid container spacing={2} sx={{ mb: 4 }}>
            <Grid item xs={6} sm={3}>
              <SummaryCard
                label="総プロジェクト"
                value={summary.total_projects}
                sublabel="チケット"
                subvalue={summary.total_issues}
                color="#1976d2"
                icon={<FolderIcon />}
              />
            </Grid>
            <Grid item xs={6} sm={3}>
              <SummaryCard
                label="遅延プロジェクト"
                value={summary.red_projects}
                sublabel="遅延チケット"
                subvalue={summary.red_issues}
                color="#d32f2f"
                icon={<RedIcon />}
              />
            </Grid>
            <Grid item xs={6} sm={3}>
              <SummaryCard
                label="注意プロジェクト"
                value={summary.yellow_projects}
                sublabel="注意チケット"
                subvalue={summary.yellow_issues}
                color="#ed6c02"
                icon={<YellowIcon />}
              />
            </Grid>
            <Grid item xs={6} sm={3}>
              <SummaryCard
                label="正常プロジェクト"
                value={summary.green_projects}
                sublabel="正常チケット"
                subvalue={summary.green_issues}
                color="#2e7d32"
                icon={<GreenIcon />}
              />
            </Grid>
          </Grid>

          {/* Heatmap section */}
          <Typography variant="h6" sx={{ mb: 2 }}>組織別プロジェクト状況</Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            カードをクリックすると該当組織のプロジェクト一覧を表示します
          </Typography>
          <OrgHeatmap roots={roots} />
        </>
      )}
    </Box>
  )
}
