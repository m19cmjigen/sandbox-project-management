import { describe, it, expect } from 'vitest'
import { buildDashboardTree } from './dashboard'
import type { DashboardOrg } from './dashboard'

const makeOrg = (partial: Partial<DashboardOrg> & Pick<DashboardOrg, 'id' | 'name'>): DashboardOrg => ({
  parent_id: null,
  level: 0,
  total_projects: 0,
  red_projects: 0,
  yellow_projects: 0,
  green_projects: 0,
  delay_status: 'GREEN',
  delay_rate: 0,
  ...partial,
})

describe('buildDashboardTree', () => {
  it('returns empty array for empty input', () => {
    expect(buildDashboardTree([])).toEqual([])
  })

  it('places root orgs at the top level', () => {
    const orgs = [makeOrg({ id: 1, name: 'Root A' }), makeOrg({ id: 2, name: 'Root B' })]
    const tree = buildDashboardTree(orgs)
    expect(tree).toHaveLength(2)
    expect(tree.map((n) => n.name)).toEqual(['Root A', 'Root B'])
  })

  it('nests child orgs under their parent', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root', total_projects: 0 }),
      makeOrg({ id: 2, name: 'Child', parent_id: 1, level: 1, total_projects: 3, red_projects: 1 }),
    ]
    const tree = buildDashboardTree(orgs)
    expect(tree).toHaveLength(1)
    expect(tree[0].children).toHaveLength(1)
    expect(tree[0].children[0].name).toBe('Child')
  })

  it('propagates child stats upward to root', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root', total_projects: 1, red_projects: 0, yellow_projects: 0, green_projects: 1 }),
      makeOrg({ id: 2, name: 'Child', parent_id: 1, level: 1, total_projects: 4, red_projects: 2, yellow_projects: 1, green_projects: 1 }),
    ]
    const tree = buildDashboardTree(orgs)
    const root = tree[0]
    expect(root.total_projects).toBe(5)    // 1 + 4
    expect(root.red_projects).toBe(2)      // 0 + 2
    expect(root.yellow_projects).toBe(1)   // 0 + 1
    expect(root.green_projects).toBe(2)    // 1 + 1
  })

  it('computes delay_status as RED when red_projects > 0 after propagation', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root' }),
      makeOrg({ id: 2, name: 'Child', parent_id: 1, level: 1, total_projects: 1, red_projects: 1 }),
    ]
    const tree = buildDashboardTree(orgs)
    expect(tree[0].delay_status).toBe('RED')
  })

  it('computes delay_rate after propagation', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root' }),
      makeOrg({ id: 2, name: 'Child', parent_id: 1, level: 1, total_projects: 4, red_projects: 1 }),
    ]
    const tree = buildDashboardTree(orgs)
    expect(tree[0].delay_rate).toBeCloseTo(0.25) // 1 / 4
  })

  it('keeps delay_rate at 0 when total_projects is 0', () => {
    const orgs = [makeOrg({ id: 1, name: 'Root', total_projects: 0 })]
    const tree = buildDashboardTree(orgs)
    expect(tree[0].delay_rate).toBe(0)
  })
})
