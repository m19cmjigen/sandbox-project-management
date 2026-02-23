import { useState } from 'react'
import {
  Box,
  Chip,
  Collapse,
  IconButton,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Tooltip,
  Typography,
} from '@mui/material'
import {
  ExpandLess,
  ExpandMore,
  ChevronRight,
  FiberManualRecord as DotIcon,
} from '@mui/icons-material'
import { useNavigate } from 'react-router-dom'
import type { OrganizationTreeNode } from '../types/organization'

interface OrganizationTreeProps {
  nodes: OrganizationTreeNode[]
  searchQuery: string
}

interface OrganizationTreeItemProps {
  node: OrganizationTreeNode
  depth: number
  searchQuery: string
  defaultExpanded: boolean
}

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

/** Returns true if the node or any descendant matches the search query. */
function nodeMatchesSearch(node: OrganizationTreeNode, query: string): boolean {
  const q = query.toLowerCase()
  if (node.name.toLowerCase().includes(q)) return true
  return node.children.some((child) => nodeMatchesSearch(child, q))
}

function OrganizationTreeItem({ node, depth, searchQuery, defaultExpanded }: OrganizationTreeItemProps) {
  const navigate = useNavigate()
  const hasChildren = node.children.length > 0
  const [expanded, setExpanded] = useState(defaultExpanded)

  const isMatch = searchQuery
    ? node.name.toLowerCase().includes(searchQuery.toLowerCase())
    : false

  // Filter children that match search
  const visibleChildren = searchQuery
    ? node.children.filter((c) => nodeMatchesSearch(c, searchQuery))
    : node.children

  // Auto-expand if search matches a descendant
  const shouldExpand = searchQuery
    ? node.children.some((c) => nodeMatchesSearch(c, searchQuery))
    : expanded

  const handleExpandToggle = (e: React.MouseEvent) => {
    e.stopPropagation()
    setExpanded((prev) => !prev)
  }

  const handleNavigate = () => {
    navigate(`/projects?organization_id=${node.id}`)
  }

  const status = node.subtree_status

  return (
    <>
      <ListItemButton
        onClick={handleNavigate}
        sx={{
          pl: depth * 3 + 1,
          pr: 1,
          borderRadius: 1,
          mb: 0.25,
          bgcolor: isMatch ? 'action.selected' : 'transparent',
          '&:hover': { bgcolor: 'action.hover' },
        }}
      >
        {/* Expand/collapse icon (takes up space even when no children for alignment) */}
        <ListItemIcon sx={{ minWidth: 32 }}>
          {hasChildren ? (
            <IconButton size="small" onClick={handleExpandToggle} sx={{ p: 0.25 }}>
              {shouldExpand ? <ExpandLess fontSize="small" /> : <ExpandMore fontSize="small" />}
            </IconButton>
          ) : (
            <ChevronRight fontSize="small" sx={{ color: 'transparent' }} />
          )}
        </ListItemIcon>

        {/* Status color dot */}
        <Tooltip title={STATUS_LABEL[status]}>
          <DotIcon sx={{ color: STATUS_COLOR[status], fontSize: 14, mr: 1, flexShrink: 0 }} />
        </Tooltip>

        {/* Organization name */}
        <ListItemText
          primary={
            <Typography
              variant="body2"
              fontWeight={depth === 0 ? 'bold' : 'normal'}
              sx={{ lineHeight: 1.4 }}
            >
              {node.name}
            </Typography>
          }
        />

        {/* Project count badges */}
        <Box sx={{ display: 'flex', gap: 0.5, flexShrink: 0, ml: 1 }}>
          {node.subtree_red > 0 && (
            <Chip
              label={node.subtree_red}
              size="small"
              color="error"
              sx={{ height: 20, fontSize: '0.65rem', '& .MuiChip-label': { px: 0.75 } }}
            />
          )}
          {node.subtree_yellow > 0 && (
            <Chip
              label={node.subtree_yellow}
              size="small"
              color="warning"
              sx={{ height: 20, fontSize: '0.65rem', '& .MuiChip-label': { px: 0.75 } }}
            />
          )}
          <Chip
            label={`${node.subtree_total}件`}
            size="small"
            variant="outlined"
            sx={{ height: 20, fontSize: '0.65rem', '& .MuiChip-label': { px: 0.75 } }}
          />
        </Box>
      </ListItemButton>

      {/* Children */}
      {hasChildren && (
        <Collapse in={shouldExpand} timeout="auto" unmountOnExit>
          <List disablePadding>
            {visibleChildren.map((child) => (
              <OrganizationTreeItem
                key={child.id}
                node={child}
                depth={depth + 1}
                searchQuery={searchQuery}
                defaultExpanded={true}
              />
            ))}
          </List>
        </Collapse>
      )}
    </>
  )
}

export default function OrganizationTree({ nodes, searchQuery }: OrganizationTreeProps) {
  const visibleNodes = searchQuery
    ? nodes.filter((n) => nodeMatchesSearch(n, searchQuery))
    : nodes

  if (visibleNodes.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary" sx={{ py: 2, textAlign: 'center' }}>
        該当する組織がありません
      </Typography>
    )
  }

  return (
    <List disablePadding>
      {visibleNodes.map((node) => (
        <OrganizationTreeItem
          key={node.id}
          node={node}
          depth={0}
          searchQuery={searchQuery}
          defaultExpanded={true}
        />
      ))}
    </List>
  )
}
