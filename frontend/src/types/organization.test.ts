import { describe, it, expect } from 'vitest'
import { buildOrganizationTree, collectSubtreeIds } from './organization'
import type { Organization } from './organization'

const makeOrg = (partial: Partial<Organization> & Pick<Organization, 'id' | 'name'>): Organization => ({
  parent_id: null,
  path: `/${partial.id}/`,
  level: 0,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  total_projects: 0,
  red_projects: 0,
  yellow_projects: 0,
  green_projects: 0,
  delay_status: 'GREEN',
  ...partial,
})

describe('buildOrganizationTree', () => {
  it('returns empty array for empty input', () => {
    expect(buildOrganizationTree([])).toEqual([])
  })

  it('builds a flat list of root nodes', () => {
    const orgs = [makeOrg({ id: 1, name: 'A' }), makeOrg({ id: 2, name: 'B' })]
    const tree = buildOrganizationTree(orgs)
    expect(tree).toHaveLength(2)
    expect(tree[0].name).toBe('A')
    expect(tree[1].name).toBe('B')
    expect(tree[0].children).toEqual([])
  })

  it('nests children under their parent', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root', level: 0 }),
      makeOrg({ id: 2, name: 'Child', parent_id: 1, level: 1 }),
    ]
    const tree = buildOrganizationTree(orgs)
    expect(tree).toHaveLength(1)
    expect(tree[0].children).toHaveLength(1)
    expect(tree[0].children[0].name).toBe('Child')
  })

  it('propagates subtree project counts upward', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root', total_projects: 1, red_projects: 0, yellow_projects: 0, green_projects: 1 }),
      makeOrg({ id: 2, name: 'Child', parent_id: 1, level: 1, total_projects: 3, red_projects: 2, yellow_projects: 1, green_projects: 0 }),
    ]
    const tree = buildOrganizationTree(orgs)
    const root = tree[0]
    expect(root.subtree_total).toBe(4)  // 1 + 3
    expect(root.subtree_red).toBe(2)    // 0 + 2
    expect(root.subtree_yellow).toBe(1) // 0 + 1
    expect(root.subtree_green).toBe(1)  // 1 + 0
  })

  it('sets subtree_status to RED when any child has red projects', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root', total_projects: 0, red_projects: 0, yellow_projects: 0, green_projects: 0 }),
      makeOrg({ id: 2, name: 'Child', parent_id: 1, level: 1, total_projects: 1, red_projects: 1, yellow_projects: 0, green_projects: 0 }),
    ]
    const tree = buildOrganizationTree(orgs)
    expect(tree[0].subtree_status).toBe('RED')
  })

  it('sets subtree_status to YELLOW when only yellow projects exist', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root', total_projects: 0, red_projects: 0, yellow_projects: 0, green_projects: 0 }),
      makeOrg({ id: 2, name: 'Child', parent_id: 1, level: 1, total_projects: 1, red_projects: 0, yellow_projects: 1, green_projects: 0 }),
    ]
    const tree = buildOrganizationTree(orgs)
    expect(tree[0].subtree_status).toBe('YELLOW')
  })

  it('sets subtree_status to GREEN when all projects are green', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root', total_projects: 0, red_projects: 0, yellow_projects: 0, green_projects: 0 }),
      makeOrg({ id: 2, name: 'Child', parent_id: 1, level: 1, total_projects: 2, red_projects: 0, yellow_projects: 0, green_projects: 2 }),
    ]
    const tree = buildOrganizationTree(orgs)
    expect(tree[0].subtree_status).toBe('GREEN')
  })
})

describe('collectSubtreeIds', () => {
  it('returns just the root id for a leaf node', () => {
    const orgs = [makeOrg({ id: 5, name: 'Root' })]
    const tree = buildOrganizationTree(orgs)
    expect(collectSubtreeIds(tree[0])).toEqual([5])
  })

  it('returns root and all descendant ids', () => {
    const orgs = [
      makeOrg({ id: 1, name: 'Root' }),
      makeOrg({ id: 2, name: 'Child1', parent_id: 1, level: 1 }),
      makeOrg({ id: 3, name: 'Child2', parent_id: 1, level: 1 }),
    ]
    const tree = buildOrganizationTree(orgs)
    const ids = collectSubtreeIds(tree[0]).sort((a, b) => a - b)
    expect(ids).toEqual([1, 2, 3])
  })
})
