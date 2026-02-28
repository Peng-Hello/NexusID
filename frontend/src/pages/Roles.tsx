import { useEffect, useState, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Plus } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '../components/ui/dialog'
import { Label } from '../components/ui/label'
import { Textarea } from '../components/ui/textarea'

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

  // Create dialog state
  const [showCreateDialog, setShowCreateDialog] = useState(false)
  const [createForm, setCreateForm] = useState({ name: '', description: '' })
  const [createError, setCreateError] = useState('')
  const [creating, setCreating] = useState(false)

  // Edit dialog state
  const [editingRole, setEditingRole] = useState<Role | null>(null)
  const [editForm, setEditForm] = useState({ name: '', description: '' })
  const [editError, setEditError] = useState('')
  const [saving, setSaving] = useState(false)

  // Delete dialog state
  const [deletingRole, setDeletingRole] = useState<Role | null>(null)
  const [deleting, setDeleting] = useState(false)

  const fetchRoles = useCallback(async () => {
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
  }, [tenantId])

  useEffect(() => {
    if (tenantId) fetchRoles()
  }, [tenantId, fetchRoles])

  const filteredRoles = roles.filter(role =>
    role.name.toLowerCase().includes(search.toLowerCase())
  )

  // Create role
  const handleCreate = async () => {
    setCreateError('')
    setCreating(true)
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/roles`, {
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
        setCreateForm({ name: '', description: '' })
        fetchRoles()
      } else {
        setCreateError(data.error || t('roles.createFailed'))
      }
    } catch {
      setCreateError(t('roles.createFailed'))
    } finally {
      setCreating(false)
    }
  }

  // Edit role
  const openEditDialog = (role: Role) => {
    setEditingRole(role)
    setEditForm({ name: role.name, description: role.description || '' })
    setEditError('')
  }

  const handleEdit = async () => {
    if (!editingRole) return
    setEditError('')
    setSaving(true)
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/roles/${editingRole.id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(editForm),
      })
      const data = await response.json()
      if (response.ok) {
        setEditingRole(null)
        fetchRoles()
      } else {
        setEditError(data.error || t('roles.updateFailed'))
      }
    } catch {
      setEditError(t('roles.updateFailed'))
    } finally {
      setSaving(false)
    }
  }

  // Delete role
  const handleDelete = async () => {
    if (!deletingRole) return
    setDeleting(true)
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/roles/${deletingRole.id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` },
      })
      if (response.ok) {
        setDeletingRole(null)
        fetchRoles()
      }
    } catch (error) {
      console.error('Failed to delete role:', error)
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{t('roles.title')}</h1>
          <p className="text-gray-600">{t('roles.subtitle')}</p>
        </div>

        {/* Create Role Dialog */}
        <Dialog open={showCreateDialog} onOpenChange={(open) => {
          setShowCreateDialog(open)
          if (!open) { setCreateForm({ name: '', description: '' }); setCreateError('') }
        }}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              {t('roles.addRole')}
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[450px]">
            <DialogHeader>
              <DialogTitle>{t('roles.createTitle')}</DialogTitle>
              <DialogDescription>{t('roles.createDesc')}</DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              {createError && (
                <div className="p-3 bg-red-50 text-red-600 text-sm rounded">{createError}</div>
              )}
              <div className="grid gap-2">
                <Label htmlFor="create-name">{t('roles.nameLabel')} *</Label>
                <Input
                  id="create-name"
                  value={createForm.name}
                  onChange={e => setCreateForm({ ...createForm, name: e.target.value })}
                  placeholder={t('roles.namePlaceholder')}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="create-desc">{t('roles.descLabel')}</Label>
                <Textarea
                  id="create-desc"
                  value={createForm.description}
                  onChange={e => setCreateForm({ ...createForm, description: e.target.value })}
                  placeholder={t('roles.descPlaceholder')}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setShowCreateDialog(false)}>{t('common.cancel')}</Button>
              <Button onClick={handleCreate} disabled={!createForm.name || creating}>
                {creating ? t('common.loading') : t('common.create')}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
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
                      {t('roles.systemRole')}
                    </span>
                  )}
                </div>
                <CardDescription>{role.description || t('common.noDescription')}</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex gap-2">
                  <Button variant="outline" size="sm" className="flex-1" onClick={() => openEditDialog(role)}>
                    {t('common.edit')}
                  </Button>
                  {!role.is_system_role && (
                    <Button variant="outline" size="sm" className="flex-1 text-red-600 hover:text-red-700" onClick={() => setDeletingRole(role)}>
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

      {/* Edit Role Dialog */}
      <Dialog open={!!editingRole} onOpenChange={(open) => { if (!open) setEditingRole(null) }}>
        <DialogContent className="sm:max-w-[450px]">
          <DialogHeader>
            <DialogTitle>{t('roles.editTitle')}</DialogTitle>
            <DialogDescription>{editingRole?.name}</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            {editError && (
              <div className="p-3 bg-red-50 text-red-600 text-sm rounded">{editError}</div>
            )}
            <div className="grid gap-2">
              <Label htmlFor="edit-name">{t('roles.nameLabel')} *</Label>
              <Input
                id="edit-name"
                value={editForm.name}
                onChange={e => setEditForm({ ...editForm, name: e.target.value })}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-desc">{t('roles.descLabel')}</Label>
              <Textarea
                id="edit-desc"
                value={editForm.description}
                onChange={e => setEditForm({ ...editForm, description: e.target.value })}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditingRole(null)}>{t('common.cancel')}</Button>
            <Button onClick={handleEdit} disabled={!editForm.name || saving}>
              {saving ? t('common.loading') : t('common.save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={!!deletingRole} onOpenChange={(open) => { if (!open) setDeletingRole(null) }}>
        <DialogContent className="sm:max-w-[400px]">
          <DialogHeader>
            <DialogTitle>{t('roles.deleteTitle')}</DialogTitle>
            <DialogDescription>
              {t('roles.deleteConfirm', { name: deletingRole?.name })}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeletingRole(null)}>{t('common.cancel')}</Button>
            <Button variant="destructive" onClick={handleDelete} disabled={deleting}>
              {deleting ? t('common.loading') : t('common.delete')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default Roles
