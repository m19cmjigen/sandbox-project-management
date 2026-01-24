import apiClient from './api'
import { Issue, IssueFilter } from '@/types'

export const issueService = {
  // チケット一覧を取得
  async getAll(filter?: IssueFilter): Promise<{ issues: Issue[]; count: number }> {
    const response = await apiClient.get('/issues', { params: filter })
    return response.data
  },

  // チケット詳細を取得
  async getById(id: number): Promise<Issue> {
    const response = await apiClient.get(`/issues/${id}`)
    return response.data
  },

  // プロジェクトのチケット一覧を取得
  async getByProject(projectId: number): Promise<{ issues: Issue[]; count: number }> {
    const response = await apiClient.get(`/projects/${projectId}/issues`)
    return response.data
  },
}
