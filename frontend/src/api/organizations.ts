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
