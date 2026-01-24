import apiClient from './api'
import { Project, ProjectWithStats } from '@/types'

export const projectService = {
  // プロジェクト一覧を取得（デフォルトで統計情報付き）
  async getAll(withStats = true): Promise<ProjectWithStats[]> {
    const response = await apiClient.get('/projects', {
      params: { with_stats: withStats },
    })
    return response.data.projects
  },

  // 組織別プロジェクト一覧を取得
  async getByOrganization(organizationId: number): Promise<Project[]> {
    const response = await apiClient.get('/projects', {
      params: { organization_id: organizationId },
    })
    return response.data.projects
  },

  // 未分類プロジェクトを取得
  async getUnassigned(): Promise<Project[]> {
    const response = await apiClient.get('/projects', {
      params: { unassigned: true },
    })
    return response.data.projects
  },

  // プロジェクト詳細を取得
  async getById(id: number, withStats = false): Promise<Project | ProjectWithStats> {
    const response = await apiClient.get(`/projects/${id}`, {
      params: { with_stats: withStats },
    })
    return response.data
  },

  // プロジェクトを組織に紐付け
  async assignToOrganization(projectId: number, organizationId: number | null): Promise<void> {
    await apiClient.put(`/projects/${projectId}/organization`, {
      organization_id: organizationId,
    })
  },
}
