import apiClient from './apiClient'
import type { Role } from '../stores/authStore'

export interface User {
  id: number
  email: string
  role: Role
  is_active: boolean
}

interface CreateUserRequest {
  email: string
  password: string
  role: Role
}

interface UpdateUserRequest {
  role?: Role
  is_active?: boolean
}

export async function getUsers(): Promise<User[]> {
  const res = await apiClient.get<{ data: User[] }>('/users')
  return res.data.data
}

export async function createUser(data: CreateUserRequest): Promise<User> {
  const res = await apiClient.post<User>('/users', data)
  return res.data
}

export async function updateUser(id: number, data: UpdateUserRequest): Promise<User> {
  const res = await apiClient.put<User>(`/users/${id}`, data)
  return res.data
}

export async function deleteUser(id: number): Promise<void> {
  await apiClient.delete(`/users/${id}`)
}

export async function changePassword(id: number, newPassword: string): Promise<void> {
  await apiClient.put(`/users/${id}/password`, { new_password: newPassword })
}
