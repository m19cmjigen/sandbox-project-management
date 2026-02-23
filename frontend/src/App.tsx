import { Routes, Route } from 'react-router-dom'
import { Box } from '@mui/material'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Organizations from './pages/Organizations'
import Projects from './pages/Projects'
import Issues from './pages/Issues'
import OrganizationManagement from './pages/OrganizationManagement'

function App() {
  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/organizations" element={<Organizations />} />
          <Route path="/organizations/manage" element={<OrganizationManagement />} />
          <Route path="/projects" element={<Projects />} />
          <Route path="/issues" element={<Issues />} />
        </Routes>
      </Layout>
    </Box>
  )
}

export default App
