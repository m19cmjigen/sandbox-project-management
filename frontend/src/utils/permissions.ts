import type { Role } from '../stores/authStore'

export const canManageUsers = (role: Role): boolean => role === 'admin'

export const canManageOrganizations = (role: Role): boolean => role === 'admin'

export const canAssignProjects = (role: Role): boolean =>
  role === 'admin' || role === 'project_manager'

export const canAccessSettings = (role: Role): boolean =>
  role === 'admin' || role === 'project_manager'
