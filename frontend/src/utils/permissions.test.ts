import { describe, it, expect } from 'vitest'
import { canManageUsers, canManageOrganizations, canAssignProjects, canAccessSettings } from './permissions'

describe('canManageUsers', () => {
  it('returns true for admin', () => {
    expect(canManageUsers('admin')).toBe(true)
  })

  it('returns false for project_manager', () => {
    expect(canManageUsers('project_manager')).toBe(false)
  })

  it('returns false for viewer', () => {
    expect(canManageUsers('viewer')).toBe(false)
  })
})

describe('canManageOrganizations', () => {
  it('returns true for admin', () => {
    expect(canManageOrganizations('admin')).toBe(true)
  })

  it('returns false for project_manager', () => {
    expect(canManageOrganizations('project_manager')).toBe(false)
  })

  it('returns false for viewer', () => {
    expect(canManageOrganizations('viewer')).toBe(false)
  })
})

describe('canAssignProjects', () => {
  it('returns true for admin', () => {
    expect(canAssignProjects('admin')).toBe(true)
  })

  it('returns true for project_manager', () => {
    expect(canAssignProjects('project_manager')).toBe(true)
  })

  it('returns false for viewer', () => {
    expect(canAssignProjects('viewer')).toBe(false)
  })
})

describe('canAccessSettings', () => {
  it('returns true for admin', () => {
    expect(canAccessSettings('admin')).toBe(true)
  })

  it('returns true for project_manager', () => {
    expect(canAccessSettings('project_manager')).toBe(true)
  })

  it('returns false for viewer', () => {
    expect(canAccessSettings('viewer')).toBe(false)
  })
})
