import type { DelayStatus } from './project'

export interface Organization {
  id: number
  name: string
  parent_id: number | null
  path: string
  level: number
  created_at: string
  updated_at: string
  total_projects: number
  red_projects: number
  yellow_projects: number
  green_projects: number
  delay_status: DelayStatus
}

export interface OrganizationTreeNode extends Organization {
  children: OrganizationTreeNode[]
  // Aggregated stats including all descendants
  subtree_total: number
  subtree_red: number
  subtree_yellow: number
  subtree_green: number
  subtree_status: DelayStatus
}

/** Build a tree from a flat list of organizations, computing subtree stats. */
export function buildOrganizationTree(orgs: Organization[]): OrganizationTreeNode[] {
  const map = new Map<number, OrganizationTreeNode>()

  orgs.forEach((org) => {
    map.set(org.id, {
      ...org,
      children: [],
      subtree_total: org.total_projects,
      subtree_red: org.red_projects,
      subtree_yellow: org.yellow_projects,
      subtree_green: org.green_projects,
      subtree_status: org.delay_status,
    })
  })

  const roots: OrganizationTreeNode[] = []

  map.forEach((node) => {
    if (node.parent_id === null) {
      roots.push(node)
    } else {
      const parent = map.get(node.parent_id)
      if (parent) parent.children.push(node)
    }
  })

  // Propagate subtree stats upward (bottom-up via path sort order, children come after parents)
  const propagate = (node: OrganizationTreeNode) => {
    node.children.forEach(propagate)
    if (node.children.length > 0) {
      node.subtree_total = node.total_projects + node.children.reduce((s, c) => s + c.subtree_total, 0)
      node.subtree_red = node.red_projects + node.children.reduce((s, c) => s + c.subtree_red, 0)
      node.subtree_yellow = node.yellow_projects + node.children.reduce((s, c) => s + c.subtree_yellow, 0)
      node.subtree_green = node.green_projects + node.children.reduce((s, c) => s + c.subtree_green, 0)
      node.subtree_status = node.subtree_red > 0 ? 'RED' : node.subtree_yellow > 0 ? 'YELLOW' : 'GREEN'
    }
  }

  roots.forEach(propagate)
  return roots
}

/** Return all IDs in the subtree rooted at node (including node itself). */
export function collectSubtreeIds(node: OrganizationTreeNode): number[] {
  return [node.id, ...node.children.flatMap(collectSubtreeIds)]
}
