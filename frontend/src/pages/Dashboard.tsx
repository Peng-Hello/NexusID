import { useEffect, useState } from 'react'
import { Link, Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Button } from '../components/ui/button'
import LanguageSwitcher from '../components/LanguageSwitcher'
import {
  LayoutDashboard,
  Users,
  Building,
  Shield,
  Key,
  LogOut,
  Menu,
  X,
} from 'lucide-react'

interface DashboardProps {
  children?: React.ReactNode
}

function Dashboard({ children }: DashboardProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const { t } = useTranslation()
  const [user, setUser] = useState<Record<string, string> | null>(null)
  const [sidebarOpen, setSidebarOpen] = useState(false)

  useEffect(() => {
    const userData = localStorage.getItem('user')
    if (!userData) {
      navigate('/login')
      return
    }
    setUser(JSON.parse(userData))
  }, [navigate])

  const handleLogout = async () => {
    try {
      const token = localStorage.getItem('token')
      const refreshToken = localStorage.getItem('refresh_token')

      if (refreshToken && token) {
        await fetch('/api/v1/oauth/revoke', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
          body: JSON.stringify({
            token: refreshToken,
            token_type_hint: 'refresh_token',
          }),
        })
      }
    } catch (err) {
      console.warn('Token revocation failed:', err)
    } finally {
      localStorage.clear()
      navigate('/login')
    }
  }

  const navigation = [
    { name: t('nav.dashboard'), href: '/dashboard', icon: LayoutDashboard, exact: true },
    { name: t('nav.tenants'), href: '/tenants', icon: Building, exact: true },
    { name: t('nav.users'), href: '/tenants/1/users', icon: Users },
    { name: t('nav.roles'), href: '/tenants/1/roles', icon: Shield },
    { name: t('nav.clients'), href: '/tenants/1/clients', icon: Key },
  ]

  if (!user) return null

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Mobile sidebar backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 lg:hidden bg-gray-600 bg-opacity-75"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div
        className={`fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-lg transform transition-transform duration-300 ease-in-out lg:translate-x-0 ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'
          }`}
      >
        <div className="flex flex-col h-full">
          <div className="flex items-center justify-between h-16 px-4 border-b">
            <h1 className="text-xl font-bold text-blue-600">NexusID</h1>
            <button
              onClick={() => setSidebarOpen(false)}
              className="lg:hidden"
            >
              <X className="w-6 h-6" />
            </button>
          </div>

          <nav className="flex-1 px-4 py-6 space-y-2 overflow-y-auto">
            {navigation.map((item) => {
              const Icon = item.icon
              const isActive = item.exact ? location.pathname === item.href : location.pathname.startsWith(item.href)
              return (
                <Link
                  key={item.href}
                  to={item.href}
                  onClick={() => setSidebarOpen(false)}
                  className={`flex items-center px-4 py-3 text-sm font-medium rounded-lg transition-colors ${isActive
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-700 hover:bg-gray-100'
                    }`}
                >
                  <Icon className="w-5 h-5 mr-3" />
                  {item.name}
                </Link>
              )
            })}
          </nav>

          <div className="p-4 border-t">
            <div className="flex items-center mb-4">
              <div className="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold">
                {user.full_name?.charAt(0).toUpperCase() || user.email.charAt(0).toUpperCase()}
              </div>
              <div className="ml-3 flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 truncate">
                  {user.full_name || 'User'}
                </p>
                <p className="text-xs text-gray-500 truncate">{user.email}</p>
              </div>
            </div>
            <Button
              onClick={handleLogout}
              variant="outline"
              className="w-full"
            >
              <LogOut className="w-4 h-4 mr-2" />
              {t('common.logout')}
            </Button>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="lg:pl-64">
        {/* Top bar */}
        <div className="sticky top-0 z-30 flex items-center justify-between h-16 px-4 bg-white border-b shadow-sm">
          <div className="flex items-center">
            <button
              onClick={() => setSidebarOpen(true)}
              className="lg:hidden mr-4"
            >
              <Menu className="w-6 h-6" />
            </button>
            <h2 className="text-lg font-semibold text-gray-900">
              {navigation.find((item) => item.exact ? location.pathname === item.href : location.pathname.startsWith(item.href))?.name || t('nav.dashboard')}
            </h2>
          </div>
          <LanguageSwitcher />
        </div>

        {/* Page content */}
        <main className="p-6">
          {children || <Outlet />}
        </main>
      </div>
    </div>
  )
}

export default Dashboard
