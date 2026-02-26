import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Plus } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'

interface Role {
  id: number
  name: string
  description: string
  is_system_role: boolean
  created_at: string
}

function Roles() {
  const { tenantId } = useParams()
  const { t } = useTranslation()
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')

  useEffect(() => {
    if (tenantId) fetchRoles()
  }, [tenantId])

  const fetchRoles = async () => {
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/roles?page=1&page_size=50`, {
        headers: { 'Authorization': `Bearer ${token}` },
      })

      if (response.ok) {
        const data = await response.json()
        setRoles(data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch roles:', error)
    } finally {
      setLoading(false)
    }
  }

  const filteredRoles = roles.filter(role =>
    role.name.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{t('roles.title')}</h1>
          <p className="text-gray-600">{t('roles.subtitle')}</p>
        </div>
        <Button>
          <Plus className="w-4 h-4 mr-2" />
          {t('roles.addRole')}
        </Button>
      </div>

      <div className="mb-6">
        <Input
          placeholder={t('roles.searchPlaceholder')}
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
          {filteredRoles.map((role) => (
            <Card key={role.id}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle>{role.name}</CardTitle>
                  {role.is_system_role && (
                    <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-gray-100 text-gray-800">
                      System
                    </span>
                  )}
                </div>
                <CardDescription>{role.description || t('common.noDescription')}</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex gap-2">
                  <Button variant="outline" size="sm" className="flex-1">
                    {t('common.edit')}
                  </Button>
                  {!role.is_system_role && (
                    <Button variant="outline" size="sm" className="flex-1 text-red-600 hover:text-red-700">
                      {t('common.delete')}
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {!loading && filteredRoles.length === 0 && (
        <Card>
          <CardContent className="text-center py-12">
            <p className="text-gray-500">{t('roles.noRoles')}</p>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

export default Roles
