import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Attach Bearer token from localStorage on every request
apiClient.interceptors.request.use((config) => {
  const stored = localStorage.getItem('auth-storage')
  if (stored) {
    try {
      const parsed = JSON.parse(stored) as { state?: { token?: string } }
      const token = parsed?.state?.token
      if (token) {
        config.headers['Authorization'] = `Bearer ${token}`
      }
    } catch {
      // JSON parse failure is non-critical; proceed without token
    }
  }
  return config
})

// Redirect to login on 401 (except when already on the login page)
apiClient.interceptors.response.use(
  (res) => res,
  (error) => {
    if (error.response?.status === 401 && !window.location.pathname.startsWith('/login')) {
      localStorage.removeItem('auth-storage')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default apiClient
