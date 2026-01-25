import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface UserInfo {
  id: number;
  username: string;
  email: string;
  full_name?: string;
  role: 'admin' | 'manager' | 'viewer';
}

export interface LoginResponse {
  token: string;
  expires_at: string;
  user: UserInfo;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  full_name?: string;
  role: 'admin' | 'manager' | 'viewer';
}

export interface UpdateUserRequest {
  email?: string;
  full_name?: string;
  role?: 'admin' | 'manager' | 'viewer';
  is_active?: boolean;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

class AuthService {
  private readonly TOKEN_KEY = 'auth_token';
  private readonly USER_KEY = 'auth_user';

  // Login
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await axios.post<LoginResponse>(`${API_BASE_URL}/auth/login`, credentials);
    const { token, user } = response.data;

    // Store token and user info
    this.setToken(token);
    this.setUser(user);

    return response.data;
  }

  // Logout
  logout(): void {
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.USER_KEY);
  }

  // Get current user
  async getCurrentUser(): Promise<UserInfo> {
    const response = await axios.get<UserInfo>(`${API_BASE_URL}/auth/me`, {
      headers: this.getAuthHeaders(),
    });
    this.setUser(response.data);
    return response.data;
  }

  // Refresh token
  async refreshToken(): Promise<LoginResponse> {
    const token = this.getToken();
    if (!token) {
      throw new Error('No token to refresh');
    }

    const response = await axios.post<LoginResponse>(`${API_BASE_URL}/auth/refresh`, { token });
    const { token: newToken, user } = response.data;

    this.setToken(newToken);
    this.setUser(user);

    return response.data;
  }

  // Token management
  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  setToken(token: string): void {
    localStorage.setItem(this.TOKEN_KEY, token);
  }

  // User management
  getUser(): UserInfo | null {
    const userStr = localStorage.getItem(this.USER_KEY);
    return userStr ? JSON.parse(userStr) : null;
  }

  setUser(user: UserInfo): void {
    localStorage.setItem(this.USER_KEY, JSON.stringify(user));
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return this.getToken() !== null;
  }

  // Get auth headers
  getAuthHeaders(): Record<string, string> {
    const token = this.getToken();
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  // User CRUD operations (admin only)
  async createUser(data: CreateUserRequest): Promise<UserInfo> {
    const response = await axios.post<UserInfo>(`${API_BASE_URL}/users`, data, {
      headers: this.getAuthHeaders(),
    });
    return response.data;
  }

  async getUsers(filters?: { role?: string; is_active?: boolean; search?: string }): Promise<UserInfo[]> {
    const params = new URLSearchParams();
    if (filters?.role) params.append('role', filters.role);
    if (filters?.is_active !== undefined) params.append('is_active', String(filters.is_active));
    if (filters?.search) params.append('search', filters.search);

    const response = await axios.get<UserInfo[]>(`${API_BASE_URL}/users?${params}`, {
      headers: this.getAuthHeaders(),
    });
    return response.data;
  }

  async getUser(id: number): Promise<UserInfo> {
    const response = await axios.get<UserInfo>(`${API_BASE_URL}/users/${id}`, {
      headers: this.getAuthHeaders(),
    });
    return response.data;
  }

  async updateUser(id: number, data: UpdateUserRequest): Promise<UserInfo> {
    const response = await axios.put<UserInfo>(`${API_BASE_URL}/users/${id}`, data, {
      headers: this.getAuthHeaders(),
    });
    return response.data;
  }

  async deleteUser(id: number): Promise<void> {
    await axios.delete(`${API_BASE_URL}/users/${id}`, {
      headers: this.getAuthHeaders(),
    });
  }

  async changePassword(id: number, data: ChangePasswordRequest): Promise<void> {
    await axios.post(`${API_BASE_URL}/users/${id}/password`, data, {
      headers: this.getAuthHeaders(),
    });
  }

  // Role checks
  isAdmin(user?: UserInfo | null): boolean {
    const currentUser = user || this.getUser();
    return currentUser?.role === 'admin';
  }

  isManager(user?: UserInfo | null): boolean {
    const currentUser = user || this.getUser();
    return currentUser?.role === 'manager' || currentUser?.role === 'admin';
  }

  canManageUsers(user?: UserInfo | null): boolean {
    return this.isAdmin(user);
  }

  canManageProjects(user?: UserInfo | null): boolean {
    return this.isManager(user);
  }
}

export default new AuthService();
