import apiClient from './apiClient'
import type { Organization } from '../types/organization'

export const getOrganizations = async (): Promise<Organization[]> => {
  const response = await apiClient.get<Organization[]>('/organizations')
  return response.data
}

export const getOrganization = async (id: number): Promise<Organization> => {
  const response = await apiClient.get<Organization>(`/organizations/${id}`)
  return response.data
}

export const createOrganization = async (data: { name: string; parent_id?: number }): Promise<Organization> => {
  const response = await apiClient.post<Organization>('/organizations', data)
  return response.data
}

export const updateOrganization = async (id: number, data: { name: string }): Promise<Organization> => {
  const response = await apiClient.put<Organization>(`/organizations/${id}`, data)
  return response.data
}

export const deleteOrganization = async (id: number): Promise<void> => {
  await apiClient.delete(`/organizations/${id}`)
}
