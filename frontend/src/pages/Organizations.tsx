import { useEffect, useState } from 'react'
import {
  Box,
  Typography,
  Card,
  CardContent,
  List,
  ListItem,
  ListItemText,
  Collapse,
  IconButton,
  Chip,
  Stack,
} from '@mui/material'
import {
  ExpandMore as ExpandMoreIcon,
  ChevronRight as ChevronRightIcon,
  Business as BusinessIcon,
} from '@mui/icons-material'
import { OrganizationWithChildren } from '@/types'
import { organizationService } from '@/services/organizationService'
import Loading from '@/components/Loading'
import ErrorMessage from '@/components/ErrorMessage'

interface TreeNodeProps {
  node: OrganizationWithChildren
  level: number
}

function TreeNode({ node, level }: TreeNodeProps) {
  const [open, setOpen] = useState(level < 2) // 第2階層まで自動展開

  const hasChildren = node.children && node.children.length > 0

  return (
    <>
      <ListItem
        sx={{
          pl: level * 4,
          borderLeft: level > 0 ? '2px solid' : 'none',
          borderColor: 'divider',
          '&:hover': {
            bgcolor: 'action.hover',
          },
        }}
      >
        {hasChildren ? (
          <IconButton
            size="small"
            onClick={() => setOpen(!open)}
            sx={{ mr: 1 }}
          >
            {open ? <ExpandMoreIcon /> : <ChevronRightIcon />}
          </IconButton>
        ) : (
          <Box sx={{ width: 40 }} />
        )}
        <BusinessIcon sx={{ mr: 2, color: 'primary.main' }} />
        <ListItemText
          primary={
            <Stack direction="row" spacing={1} alignItems="center">
              <Typography variant="body1" fontWeight={level === 0 ? 'bold' : 'normal'}>
                {node.name}
              </Typography>
              <Chip
                label={`レベル ${node.level}`}
                size="small"
                variant="outlined"
              />
              {hasChildren && (
                <Chip
                  label={`${node.children!.length} 件の下位組織`}
                  size="small"
                  color="primary"
                />
              )}
            </Stack>
          }
          secondary={`組織ID: ${node.id} | パス: ${node.path}`}
        />
      </ListItem>
      {hasChildren && (
        <Collapse in={open} timeout="auto" unmountOnExit>
          <List component="div" disablePadding>
            {node.children!.map((child) => (
              <TreeNode key={child.id} node={child} level={level + 1} />
            ))}
          </List>
        </Collapse>
      )}
    </>
  )
}

export default function Organizations() {
  const [tree, setTree] = useState<OrganizationWithChildren[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchTree = async () => {
      try {
        setLoading(true)
        const data = await organizationService.getTree()
        setTree(data)
      } catch (err) {
        setError('組織データの取得に失敗しました')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    fetchTree()
  }, [])

  if (loading) return <Loading />
  if (error) return <ErrorMessage message={error} />

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        組織管理
      </Typography>

      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            組織階層
          </Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mb: 2 }}>
            全社の組織構造をツリー形式で表示しています
          </Typography>

          {tree.length > 0 ? (
            <List>
              {tree.map((rootNode) => (
                <TreeNode key={rootNode.id} node={rootNode} level={0} />
              ))}
            </List>
          ) : (
            <Typography variant="body2" color="textSecondary">
              組織データがありません
            </Typography>
          )}
        </CardContent>
      </Card>

      <Card sx={{ mt: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            組織管理について
          </Typography>
          <Typography variant="body2" color="textSecondary" paragraph>
            組織の追加・編集・削除機能は管理者権限が必要です。
          </Typography>
          <Typography variant="body2" color="textSecondary">
            各プロジェクトは組織に紐付けることで、組織単位での進捗管理が可能になります。
          </Typography>
        </CardContent>
      </Card>
    </Box>
  )
}
