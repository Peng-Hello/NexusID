import { Outlet } from 'react-router-dom'
import Dashboard from './Dashboard'

function ProtectedLayout() {
  return (
    <Dashboard>
      <Outlet />
    </Dashboard>
  )
}

export default ProtectedLayout
