import type { DelayStatus } from './project'
import type { Project } from './project'
import type { Issue } from './issue'

export interface DashboardOrg {
  id: number
  name: string
  parent_id: number | null
  level: number
  total_projects: number
  red_projects: number
  yellow_projects: number
  green_projects: number
  delay_status: DelayStatus
  delay_rate: number
}

export interface DashboardSummary {
  total_projects: number
  red_projects: number
  yellow_projects: number
  green_projects: number
  total_issues: number
  red_issues: number
  yellow_issues: number
  green_issues: number
  organizations: DashboardOrg[]
}

export interface OrgSummaryResponse {
  organization: DashboardOrg
  projects: Project[]
}

export interface ProjectIssueSummary {
  red_count: number
  yellow_count: number
  green_count: number
  open_count: number
  total_count: number
}

export interface ProjectSummaryResponse {
  project: Project
  delayed_issues: Issue[]
  summary: ProjectIssueSummary
}

export interface DashboardOrgNode extends DashboardOrg {
  children: DashboardOrgNode[]
}

/** Build a tree from a flat org list, propagating subtree stats upward. */
export function buildDashboardTree(orgs: DashboardOrg[]): DashboardOrgNode[] {
  const map = new Map<number, DashboardOrgNode>()
  orgs.forEach((o) => map.set(o.id, { ...o, children: [] }))

  const roots: DashboardOrgNode[] = []
  map.forEach((node) => {
    if (node.parent_id === null) {
      roots.push(node)
    } else {
      const parent = map.get(node.parent_id)
      if (parent) parent.children.push(node)
    }
  })

  // Propagate child stats up to parent so root orgs reflect subtree totals
  const propagate = (node: DashboardOrgNode) => {
    node.children.forEach(propagate)
    if (node.children.length > 0) {
      node.total_projects += node.children.reduce((s, c) => s + c.total_projects, 0)
      node.red_projects += node.children.reduce((s, c) => s + c.red_projects, 0)
      node.yellow_projects += node.children.reduce((s, c) => s + c.yellow_projects, 0)
      node.green_projects += node.children.reduce((s, c) => s + c.green_projects, 0)
      node.delay_status = node.red_projects > 0 ? 'RED' : node.yellow_projects > 0 ? 'YELLOW' : 'GREEN'
      node.delay_rate = node.total_projects > 0 ? node.red_projects / node.total_projects : 0
    }
  }
  roots.forEach(propagate)
  return roots
}
