import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import DashboardHome from './pages/DashboardHome'
import Tenants from './pages/Tenants'
import Users from './pages/Users'
import Roles from './pages/Roles'
import Clients from './pages/Clients'
import Consent from './pages/Consent'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/oauth/consent" element={<Consent />} />

        {/* Protected routes - wrapped in Dashboard layout */}
        <Route path="/" element={<Dashboard />}>
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<DashboardHome />} />
          <Route path="tenants" element={<Tenants />} />
          <Route path="tenants/:tenantId/users" element={<Users />} />
          <Route path="tenants/:tenantId/roles" element={<Roles />} />
          <Route path="tenants/:tenantId/clients" element={<Clients />} />
        </Route>

        {/* Catch all */}
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
