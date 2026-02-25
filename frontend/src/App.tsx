import { Routes, Route } from 'react-router-dom'
import { Box } from '@mui/material'
import Layout from './components/Layout'
import ProtectedRoute from './components/ProtectedRoute'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import Organizations from './pages/Organizations'
import Projects from './pages/Projects'
import Issues from './pages/Issues'
import OrganizationManagement from './pages/OrganizationManagement'
import UserManagement from './pages/UserManagement'
import Settings from './pages/Settings'
import ProjectDetail from './pages/ProjectDetail'

function App() {
  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route element={<ProtectedRoute />}>
          <Route
            element={
              <Layout>
                <Routes>
                  <Route path="/" element={<Dashboard />} />
                  <Route path="/organizations" element={<Organizations />} />
                  <Route path="/organizations/manage" element={<OrganizationManagement />} />
                  <Route path="/projects" element={<Projects />} />
                  <Route path="/projects/:id" element={<ProjectDetail />} />
                  <Route path="/issues" element={<Issues />} />
                  <Route path="/users" element={<UserManagement />} />
                  <Route path="/settings" element={<Settings />} />
                </Routes>
              </Layout>
            }
            path="/*"
          />
        </Route>
      </Routes>
    </Box>
  )
}

export default App
