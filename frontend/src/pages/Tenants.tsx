import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Plus } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'

interface Tenant {
  id: number
  name: string
  slug: string
  description: string
  is_active: boolean
  created_at: string
}

function Tenants() {
  const { t } = useTranslation()
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')

  useEffect(() => {
    fetchTenants()
  }, [])

  const fetchTenants = async () => {
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch('/api/v1/tenants?page=1&page_size=50', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (response.ok) {
        const data = await response.json()
        setTenants(data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch tenants:', error)
    } finally {
      setLoading(false)
    }
  }

  const filteredTenants = tenants.filter(tenant =>
    tenant.name.toLowerCase().includes(search.toLowerCase()) ||
    tenant.slug.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{t('tenants.title')}</h1>
          <p className="text-gray-600">{t('tenants.subtitle')}</p>
        </div>
        <Button>
          <Plus className="w-4 h-4 mr-2" />
          {t('tenants.addTenant')}
        </Button>
      </div>

      <div className="mb-6">
        <Input
          placeholder={t('tenants.searchPlaceholder')}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredTenants.map((tenant) => (
            <Card key={tenant.id}>
              <CardHeader>
                <CardTitle>{tenant.name}</CardTitle>
                <CardDescription>{tenant.slug}</CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-gray-600 mb-4">
                  {tenant.description || t('common.noDescription')}
                </p>
                <div className="flex items-center justify-between">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${tenant.is_active
                    ? 'bg-green-100 text-green-800'
                    : 'bg-gray-100 text-gray-800'
                    }`}>
                    {tenant.is_active ? t('common.active') : t('common.inactive')}
                  </span>
                  <Button variant="outline" size="sm">
                    {t('common.manage')}
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {!loading && filteredTenants.length === 0 && (
        <Card>
          <CardContent className="text-center py-12">
            <p className="text-gray-500">{t('tenants.noTenants')}</p>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

export default Tenants
