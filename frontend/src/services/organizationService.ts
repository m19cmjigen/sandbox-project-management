import apiClient from './api'
import { Organization, OrganizationWithChildren } from '@/types'

export const organizationService = {
  // 組織一覧を取得
  async getAll(): Promise<Organization[]> {
    const response = await apiClient.get('/organizations')
    return response.data.organizations
  },

  // 組織ツリーを取得
  async getTree(): Promise<OrganizationWithChildren[]> {
    const response = await apiClient.get('/organizations/tree')
    return response.data.tree
  },

  // 組織詳細を取得
  async getById(id: number): Promise<Organization> {
    const response = await apiClient.get(`/organizations/${id}`)
    return response.data
  },

  // 子組織を取得
  async getChildren(id: number): Promise<Organization[]> {
    const response = await apiClient.get(`/organizations/${id}/children`)
    return response.data.children
  },

  // 組織を作成
  async create(name: string, parentId: number | null): Promise<Organization> {
    const response = await apiClient.post('/organizations', {
      name,
      parent_id: parentId,
    })
    return response.data
  },

  // 組織を更新
  async update(id: number, name: string, parentId: number | null): Promise<Organization> {
    const response = await apiClient.put(`/organizations/${id}`, {
      name,
      parent_id: parentId,
    })
    return response.data
  },

  // 組織を削除
  async delete(id: number): Promise<void> {
    await apiClient.delete(`/organizations/${id}`)
  },
}
