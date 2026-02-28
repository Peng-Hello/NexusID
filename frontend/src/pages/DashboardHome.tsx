import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Button } from '../components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'
import { Building, Users, Key, Shield, ArrowRight } from 'lucide-react'

function DashboardHome() {
  const { t } = useTranslation()

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">{t('dashboard.welcome')}</h1>
        <p className="text-gray-600 mt-2">{t('dashboard.subtitle')}</p>
      </div>

      {/* Stats cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">{t('nav.tenants')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold">1</div>
            <p className="text-xs text-gray-500 mt-1">{t('dashboard.activeOrganizations')}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">{t('nav.users')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold">0</div>
            <p className="text-xs text-gray-500 mt-1">{t('dashboard.totalUsers')}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">{t('nav.roles')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold">3</div>
            <p className="text-xs text-gray-500 mt-1">{t('dashboard.systemRoles')}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">{t('nav.clients')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold">1</div>
            <p className="text-xs text-gray-500 mt-1">{t('dashboard.oidcApplications')}</p>
          </CardContent>
        </Card>
      </div>

      {/* Quick actions */}
      <h2 className="text-xl font-semibold text-gray-900 mb-4">{t('dashboard.quickActions')}</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Link to="/tenants">
          <Card className="hover:shadow-md transition-shadow cursor-pointer">
            <CardHeader>
              <Building className="w-8 h-8 text-blue-600 mb-2" />
              <CardTitle>{t('dashboard.manageTenants')}</CardTitle>
              <CardDescription>{t('dashboard.manageTenantsDesc')}</CardDescription>
            </CardHeader>
            <CardContent>
              <Button variant="ghost" className="w-full justify-start">
                {t('dashboard.goToTenants')} <ArrowRight className="w-4 h-4 ml-auto" />
              </Button>
            </CardContent>
          </Card>
        </Link>

        <Link to="/tenants/1/users">
          <Card className="hover:shadow-md transition-shadow cursor-pointer">
            <CardHeader>
              <Users className="w-8 h-8 text-green-600 mb-2" />
              <CardTitle>{t('dashboard.manageUsers')}</CardTitle>
              <CardDescription>{t('dashboard.manageUsersDesc')}</CardDescription>
            </CardHeader>
            <CardContent>
              <Button variant="ghost" className="w-full justify-start">
                {t('dashboard.goToUsers')} <ArrowRight className="w-4 h-4 ml-auto" />
              </Button>
            </CardContent>
          </Card>
        </Link>

        <Link to="/tenants/1/roles">
          <Card className="hover:shadow-md transition-shadow cursor-pointer">
            <CardHeader>
              <Shield className="w-8 h-8 text-purple-600 mb-2" />
              <CardTitle>{t('dashboard.manageRoles')}</CardTitle>
              <CardDescription>{t('dashboard.manageRolesDesc')}</CardDescription>
            </CardHeader>
            <CardContent>
              <Button variant="ghost" className="w-full justify-start">
                {t('dashboard.goToRoles')} <ArrowRight className="w-4 h-4 ml-auto" />
              </Button>
            </CardContent>
          </Card>
        </Link>

        <Link to="/tenants/1/clients">
          <Card className="hover:shadow-md transition-shadow cursor-pointer">
            <CardHeader>
              <Key className="w-8 h-8 text-orange-600 mb-2" />
              <CardTitle>{t('dashboard.manageClients')}</CardTitle>
              <CardDescription>{t('dashboard.manageClientsDesc')}</CardDescription>
            </CardHeader>
            <CardContent>
              <Button variant="ghost" className="w-full justify-start">
                {t('dashboard.goToClients')} <ArrowRight className="w-4 h-4 ml-auto" />
              </Button>
            </CardContent>
          </Card>
        </Link>
      </div>
    </div>
  )
}

export default DashboardHome
