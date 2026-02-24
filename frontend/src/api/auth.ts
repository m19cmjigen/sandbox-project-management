import apiClient from './apiClient'
import type { AuthUser } from '../stores/authStore'

interface LoginRequest {
  email: string
  password: string
}

interface LoginResponse {
  access_token: string
  token_type: string
  expires_in: number
  user: {
    id: number
    email: string
    role: string
  }
}

export async function login(data: LoginRequest): Promise<{ token: string; user: AuthUser }> {
  const res = await apiClient.post<LoginResponse>('/auth/login', data)
  return {
    token: res.data.access_token,
    user: {
      id: res.data.user.id,
      email: res.data.user.email,
      role: res.data.user.role as AuthUser['role'],
    },
  }
}
