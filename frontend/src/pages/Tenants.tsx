import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Plus, Pencil, Users, Shield, Key } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '../components/ui/dialog'
import { Label } from '../components/ui/label'
import { Textarea } from '../components/ui/textarea'
import { Checkbox } from '../components/ui/checkbox'

interface Tenant {
  id: number
  name: string
  slug: string
  description: string
  is_active: boolean
  created_at: string
}

function Tenants() {
  const navigate = useNavigate()
  const { t } = useTranslation()
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')

  // Create dialog state
  const [showCreateDialog, setShowCreateDialog] = useState(false)
  const [createForm, setCreateForm] = useState({ name: '', slug: '', description: '' })
  const [createError, setCreateError] = useState('')
  const [creating, setCreating] = useState(false)

  // Edit dialog state
  const [editingTenant, setEditingTenant] = useState<Tenant | null>(null)
  const [editForm, setEditForm] = useState({ name: '', slug: '', description: '', is_active: true })
  const [editError, setEditError] = useState('')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    fetchTenants()
  }, [])

  const fetchTenants = async () => {
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch('/api/v1/tenants?page=1&page_size=50', {
        headers: { 'Authorization': `Bearer ${token}` },
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

  // Auto-generate slug from name
  const generateSlug = (name: string) => {
    return name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')
  }

  // Create tenant
  const handleCreate = async () => {
    setCreateError('')
    setCreating(true)
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch('/api/v1/tenants', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(createForm),
      })
      const data = await response.json()
      if (response.ok) {
        setShowCreateDialog(false)
        setCreateForm({ name: '', slug: '', description: '' })
        fetchTenants()
      } else {
        setCreateError(data.error || t('tenants.createFailed'))
      }
    } catch {
      setCreateError(t('tenants.createFailed'))
    } finally {
      setCreating(false)
    }
  }

  // Edit tenant
  const openEditDialog = (tenant: Tenant) => {
    setEditingTenant(tenant)
    setEditForm({
      name: tenant.name,
      slug: tenant.slug,
      description: tenant.description || '',
      is_active: tenant.is_active,
    })
    setEditError('')
  }

  const handleEdit = async () => {
    if (!editingTenant) return
    setEditError('')
    setSaving(true)
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${editingTenant.id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(editForm),
      })
      const data = await response.json()
      if (response.ok) {
        setEditingTenant(null)
        fetchTenants()
      } else {
        setEditError(data.error || t('tenants.updateFailed'))
      }
    } catch {
      setEditError(t('tenants.updateFailed'))
    } finally {
      setSaving(false)
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{t('tenants.title')}</h1>
          <p className="text-gray-600">{t('tenants.subtitle')}</p>
        </div>

        {/* Create Tenant Dialog */}
        <Dialog open={showCreateDialog} onOpenChange={(open) => {
          setShowCreateDialog(open)
          if (!open) {
            setCreateForm({ name: '', slug: '', description: '' })
            setCreateError('')
          }
        }}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              {t('tenants.addTenant')}
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[450px]">
            <DialogHeader>
              <DialogTitle>{t('tenants.createTitle')}</DialogTitle>
              <DialogDescription>{t('tenants.createDesc')}</DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              {createError && (
                <div className="p-3 bg-red-50 text-red-600 text-sm rounded">{createError}</div>
              )}
              <div className="grid gap-2">
                <Label htmlFor="create-name">{t('tenants.nameLabel')} *</Label>
                <Input
                  id="create-name"
                  value={createForm.name}
                  onChange={e => {
                    const name = e.target.value
                    setCreateForm({ ...createForm, name, slug: generateSlug(name) })
                  }}
                  placeholder={t('tenants.namePlaceholder')}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="create-slug">{t('tenants.slugLabel')} *</Label>
                <Input
                  id="create-slug"
                  value={createForm.slug}
                  onChange={e => setCreateForm({ ...createForm, slug: e.target.value })}
                  placeholder={t('tenants.slugPlaceholder')}
                />
                <p className="text-xs text-gray-500">{t('tenants.slugHint')}</p>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="create-desc">{t('tenants.descLabel')}</Label>
                <Textarea
                  id="create-desc"
                  value={createForm.description}
                  onChange={e => setCreateForm({ ...createForm, description: e.target.value })}
                  placeholder={t('tenants.descPlaceholder')}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setShowCreateDialog(false)}>{t('common.cancel')}</Button>
              <Button onClick={handleCreate} disabled={!createForm.name || !createForm.slug || creating}>
                {creating ? t('common.loading') : t('common.create')}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
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
            <Card key={tenant.id} className="hover:shadow-md transition-shadow">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle className="text-lg">{tenant.name}</CardTitle>
                    <CardDescription className="mt-0.5 font-mono text-xs">{tenant.slug}</CardDescription>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${tenant.is_active
                      ? 'bg-green-100 text-green-800'
                      : 'bg-gray-100 text-gray-800'
                      }`}>
                      {tenant.is_active ? t('common.active') : t('common.inactive')}
                    </span>
                    <Button variant="ghost" size="sm" onClick={() => openEditDialog(tenant)} title={t('common.edit')}>
                      <Pencil className="w-3.5 h-3.5" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-gray-500 mb-4 min-h-[36px]">
                  {tenant.description || t('common.noDescription')}
                </p>
                {/* Management entry buttons */}
                <div className="grid grid-cols-3 gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    className="flex flex-col h-auto py-2 gap-1"
                    onClick={() => navigate(`/tenants/${tenant.id}/users`)}
                  >
                    <Users className="w-4 h-4" />
                    <span className="text-xs">{t('nav.users')}</span>
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="flex flex-col h-auto py-2 gap-1"
                    onClick={() => navigate(`/tenants/${tenant.id}/roles`)}
                  >
                    <Shield className="w-4 h-4" />
                    <span className="text-xs">{t('nav.roles')}</span>
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="flex flex-col h-auto py-2 gap-1"
                    onClick={() => navigate(`/tenants/${tenant.id}/clients`)}
                  >
                    <Key className="w-4 h-4" />
                    <span className="text-xs">{t('nav.clients')}</span>
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

      {/* Edit Tenant Dialog */}
      <Dialog open={!!editingTenant} onOpenChange={(open) => { if (!open) setEditingTenant(null) }}>
        <DialogContent className="sm:max-w-[450px]">
          <DialogHeader>
            <DialogTitle>{t('tenants.editTitle')}</DialogTitle>
            <DialogDescription>{t('tenants.editDesc')}</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            {editError && (
              <div className="p-3 bg-red-50 text-red-600 text-sm rounded">{editError}</div>
            )}
            <div className="grid gap-2">
              <Label htmlFor="edit-name">{t('tenants.nameLabel')} *</Label>
              <Input
                id="edit-name"
                value={editForm.name}
                onChange={e => setEditForm({ ...editForm, name: e.target.value })}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-slug">{t('tenants.slugLabel')} *</Label>
              <Input
                id="edit-slug"
                value={editForm.slug}
                onChange={e => setEditForm({ ...editForm, slug: e.target.value })}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-desc">{t('tenants.descLabel')}</Label>
              <Textarea
                id="edit-desc"
                value={editForm.description}
                onChange={e => setEditForm({ ...editForm, description: e.target.value })}
              />
            </div>
            <div className="flex items-center space-x-2 mt-2">
              <Checkbox
                id="edit-active"
                checked={editForm.is_active}
                onCheckedChange={(checked) => setEditForm({ ...editForm, is_active: checked === true })}
              />
              <Label htmlFor="edit-active" className="font-normal">{t('tenants.isActive')}</Label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditingTenant(null)}>{t('common.cancel')}</Button>
            <Button onClick={handleEdit} disabled={!editForm.name || !editForm.slug || saving}>
              {saving ? t('common.loading') : t('common.save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default Tenants
