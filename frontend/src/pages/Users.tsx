import { useEffect, useState, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Card, CardContent } from '../components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '../components/ui/dialog'
import { Label } from '../components/ui/label'
import { Checkbox } from '../components/ui/checkbox'

interface User {
  id: number
  email: string
  full_name: string
  is_active: boolean
  is_email_verified: boolean
  roles: string[]
  created_at: string
}

interface Role {
  id: number
  name: string
  description: string
  is_system: boolean
}

function Users() {
  const { tenantId } = useParams()
  const { t } = useTranslation()
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')

  // Create user state
  const [showCreateDialog, setShowCreateDialog] = useState(false)
  const [createForm, setCreateForm] = useState({ email: '', password: '', full_name: '' })
  const [createError, setCreateError] = useState('')
  const [creating, setCreating] = useState(false)
  const [createUserRoleIds, setCreateUserRoleIds] = useState<Set<number>>(new Set())

  // Edit user state
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [editForm, setEditForm] = useState({ full_name: '', is_active: true })
  const [editError, setEditError] = useState('')
  const [saving, setSaving] = useState(false)

  // Delete user state
  const [deletingUser, setDeletingUser] = useState<User | null>(null)
  const [deleting, setDeleting] = useState(false)

  // Role assignment state
  const [allRoles, setAllRoles] = useState<Role[]>([])
  const [userRoleIds, setUserRoleIds] = useState<Set<number>>(new Set())
  const [roleLoading, setRoleLoading] = useState(false)

  const fetchUsers = useCallback(async () => {
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/users?page=1&page_size=50`, {
        headers: { 'Authorization': `Bearer ${token}` },
      })
      if (response.ok) {
        const data = await response.json()
        setUsers(data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch users:', error)
    } finally {
      setLoading(false)
    }
  }, [tenantId])

  useEffect(() => {
    if (tenantId) fetchUsers()
  }, [tenantId, fetchUsers])

  const filteredUsers = users.filter(user =>
    user.email.toLowerCase().includes(search.toLowerCase()) ||
    (user.full_name && user.full_name.toLowerCase().includes(search.toLowerCase()))
  )

  // Create user
  const openCreateDialog = async () => {
    setShowCreateDialog(true)
    setCreateForm({ email: '', password: '', full_name: '' })
    setCreateError('')
    setCreateUserRoleIds(new Set())

    setRoleLoading(true)
    try {
      const token = localStorage.getItem('accessToken')
      const rolesRes = await fetch(`/api/v1/tenants/${tenantId}/roles?page=1&page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } })
      const rolesData = await rolesRes.json()
      setAllRoles(rolesData.data || [])
    } catch {
      setAllRoles([])
    } finally {
      setRoleLoading(false)
    }
  }

  const handleCreate = async () => {
    setCreateError('')
    setCreating(true)
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/users`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(createForm),
      })
      const data = await response.json()
      if (response.ok) {
        // Assign roles if any were checked
        if (createUserRoleIds.size > 0) {
          const assignPromises = Array.from(createUserRoleIds).map(roleId =>
            fetch(`/api/v1/tenants/${tenantId}/roles/assign`, {
              method: 'POST',
              headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
              body: JSON.stringify({ user_id: data.id, role_id: roleId }),
            })
          )
          await Promise.allSettled(assignPromises)
        }

        setShowCreateDialog(false)
        setCreateForm({ email: '', password: '', full_name: '' })
        setCreateUserRoleIds(new Set())
        fetchUsers()
      } else {
        setCreateError(data.error || t('users.createFailed'))
      }
    } catch {
      setCreateError(t('users.createFailed'))
    } finally {
      setCreating(false)
    }
  }

  const openEditDialog = async (user: User) => {
    setEditingUser(user)
    setEditForm({ full_name: user.full_name || '', is_active: user.is_active })
    setEditError('')
    setRoleLoading(true)
    try {
      const token = localStorage.getItem('accessToken')
      const [rolesRes, userRolesRes] = await Promise.all([
        fetch(`/api/v1/tenants/${tenantId}/roles?page=1&page_size=100`, { headers: { 'Authorization': `Bearer ${token}` } }),
        fetch(`/api/v1/tenants/${tenantId}/users/${user.id}/roles`, { headers: { 'Authorization': `Bearer ${token}` } }),
      ])
      const [rolesData, userRolesData] = await Promise.all([rolesRes.json(), userRolesRes.json()])
      setAllRoles(rolesData.data || [])
      const assigned: Role[] = userRolesData.data || []
      setUserRoleIds(new Set(assigned.map((r: Role) => r.id)))
    } catch {
      setAllRoles([])
    } finally {
      setRoleLoading(false)
    }
  }

  const handleEdit = async () => {
    if (!editingUser) return
    setEditError('')
    setSaving(true)
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/users/${editingUser.id}`, {
        method: 'PUT',
        headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
        body: JSON.stringify(editForm),
      })
      const data = await response.json()
      if (response.ok) {
        setEditingUser(null)
        fetchUsers()
      } else {
        setEditError(data.error || t('users.updateFailed'))
      }
    } catch {
      setEditError(t('users.updateFailed'))
    } finally {
      setSaving(false)
    }
  }

  const handleToggleRole = async (role: Role, checked: boolean) => {
    if (!editingUser) return
    setEditError('')
    const token = localStorage.getItem('accessToken')
    const endpoint = checked ? 'assign' : 'revoke'
    try {
      const response = await fetch(`/api/v1/tenants/${tenantId}/roles/${endpoint}`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: editingUser.id, role_id: role.id }),
      })
      if (!response.ok) {
        const data = await response.json()
        setEditError(data.error || t('users.updateFailed'))
        return
      }
      setUserRoleIds(prev => {
        const next = new Set(prev)
        if (checked) {
          next.add(role.id)
        } else {
          next.delete(role.id)
        }
        return next
      })
    } catch {
      setEditError(t('users.updateFailed'))
    }
  }

  // Delete user
  const handleDelete = async () => {
    if (!deletingUser) return
    setDeleting(true)
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/users/${deletingUser.id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` },
      })
      if (response.ok) {
        setDeletingUser(null)
        fetchUsers()
      }
    } catch (error) {
      console.error('Failed to delete user:', error)
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{t('users.title')}</h1>
          <p className="text-gray-600">{t('users.subtitle')}</p>
        </div>

        {/* Create User Dialog */}
        <Dialog open={showCreateDialog} onOpenChange={(open) => {
          if (!open) {
            setShowCreateDialog(false)
            setCreateForm({ email: '', password: '', full_name: '' })
            setCreateError('')
            setCreateUserRoleIds(new Set())
          }
        }}>
          <DialogTrigger asChild>
            <Button onClick={(e) => { e.preventDefault(); openCreateDialog(); }}>
              <Plus className="w-4 h-4 mr-2" />
              {t('users.addUser')}
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[450px]">
            <DialogHeader>
              <DialogTitle>{t('users.createTitle')}</DialogTitle>
              <DialogDescription>{t('users.createDesc')}</DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              {createError && (
                <div className="p-3 bg-red-50 text-red-600 text-sm rounded">{createError}</div>
              )}
              <div className="grid gap-2">
                <Label htmlFor="create-name">{t('users.fullNameLabel')}</Label>
                <Input
                  id="create-name"
                  value={createForm.full_name}
                  onChange={e => setCreateForm({ ...createForm, full_name: e.target.value })}
                  placeholder={t('users.fullNamePlaceholder')}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="create-email">{t('users.emailLabel')} *</Label>
                <Input
                  id="create-email"
                  type="email"
                  value={createForm.email}
                  onChange={e => setCreateForm({ ...createForm, email: e.target.value })}
                  placeholder={t('users.emailPlaceholder')}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="create-password">{t('users.passwordLabel')} *</Label>
                <Input
                  id="create-password"
                  type="password"
                  value={createForm.password}
                  onChange={e => setCreateForm({ ...createForm, password: e.target.value })}
                  placeholder={t('users.passwordPlaceholder')}
                />
                <p className="text-xs text-gray-500">{t('users.passwordHint')}</p>
              </div>

              {/* Role Assignment During Creation */}
              <div className="mt-2">
                <Label className="text-sm font-medium">{t('users.rolesLabel')}</Label>
                {roleLoading ? (
                  <p className="text-xs text-gray-500 mt-1">{t('common.loading')}</p>
                ) : allRoles.length === 0 ? (
                  <p className="text-xs text-gray-400 mt-1">{t('roles.noRoles')}</p>
                ) : (
                  <div className="mt-2 border rounded p-3 max-h-40 overflow-y-auto space-y-2">
                    {allRoles.map(role => (
                      <div key={role.id} className="flex items-center space-x-2">
                        <Checkbox
                          id={`create-role-${role.id}`}
                          checked={createUserRoleIds.has(role.id)}
                          onCheckedChange={(checked) => {
                            setCreateUserRoleIds(prev => {
                              const next = new Set(prev)
                              if (checked) next.add(role.id)
                              else next.delete(role.id)
                              return next
                            })
                          }}
                        />
                        <Label htmlFor={`create-role-${role.id}`} className="font-normal cursor-pointer">
                          {role.name}
                          {role.is_system && <span className="ml-1 text-xs text-blue-500">(system)</span>}
                        </Label>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setShowCreateDialog(false)}>{t('common.cancel')}</Button>
              <Button onClick={handleCreate} disabled={!createForm.email || !createForm.password || creating}>
                {creating ? t('common.loading') : t('common.create')}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className="mb-6">
        <Input
          placeholder={t('users.searchPlaceholder')}
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
        <Card>
          <CardContent className="p-0">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{t('nav.users')}</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{t('nav.roles')}</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{t('users.statusColumn')}</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{t('users.actionsColumn')}</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredUsers.map((user) => (
                  <tr key={user.id}>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>
                        <div className="text-sm font-medium text-gray-900">{user.full_name || 'N/A'}</div>
                        <div className="text-sm text-gray-500">{user.email}</div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex flex-wrap gap-1">
                        {(user.roles ?? []).map((role) => (
                          <span key={role} className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-blue-100 text-blue-800">
                            {role}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${user.is_active
                        ? 'bg-green-100 text-green-800'
                        : 'bg-gray-100 text-gray-800'
                        }`}>
                        {user.is_active ? t('common.active') : t('common.inactive')}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button
                        className="text-blue-600 hover:text-blue-900 mr-3"
                        onClick={() => openEditDialog(user)}
                        title={t('common.edit')}
                      >
                        <Pencil className="w-4 h-4" />
                      </button>
                      <button
                        className="text-red-600 hover:text-red-900"
                        onClick={() => setDeletingUser(user)}
                        title={t('common.delete')}
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      )}

      {!loading && filteredUsers.length === 0 && (
        <Card>
          <CardContent className="text-center py-12">
            <p className="text-gray-500">{t('users.noUsers')}</p>
          </CardContent>
        </Card>
      )}

      {/* Edit User Dialog */}
      <Dialog open={!!editingUser} onOpenChange={(open) => { if (!open) { setEditingUser(null); fetchUsers(); } }}>
        <DialogContent className="sm:max-w-[480px]">
          <DialogHeader>
            <DialogTitle>{t('users.editTitle')}</DialogTitle>
            <DialogDescription>{editingUser?.email}</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            {editError && (
              <div className="p-3 bg-red-50 text-red-600 text-sm rounded">{editError}</div>
            )}
            <div className="grid gap-2">
              <Label htmlFor="edit-name">{t('users.fullNameLabel')}</Label>
              <Input
                id="edit-name"
                value={editForm.full_name}
                onChange={e => setEditForm({ ...editForm, full_name: e.target.value })}
              />
            </div>
            <div className="flex items-center space-x-2 mt-2">
              <Checkbox
                id="edit-active"
                checked={editForm.is_active}
                onCheckedChange={(checked) => setEditForm({ ...editForm, is_active: checked === true })}
              />
              <Label htmlFor="edit-active" className="font-normal">{t('users.isActive')}</Label>
            </div>

            {/* Role Assignment */}
            <div className="mt-2">
              <Label className="text-sm font-medium">{t('users.rolesLabel')}</Label>
              {roleLoading ? (
                <p className="text-xs text-gray-500 mt-1">{t('common.loading')}</p>
              ) : allRoles.length === 0 ? (
                <p className="text-xs text-gray-400 mt-1">{t('roles.noRoles')}</p>
              ) : (
                <div className="mt-2 border rounded p-3 max-h-40 overflow-y-auto space-y-2">
                  {allRoles.map(role => (
                    <div key={role.id} className="flex items-center space-x-2">
                      <Checkbox
                        id={`role-${role.id}`}
                        checked={userRoleIds.has(role.id)}
                        onCheckedChange={(checked) => handleToggleRole(role, checked === true)}
                      />
                      <Label htmlFor={`role-${role.id}`} className="font-normal cursor-pointer">
                        {role.name}
                        {role.is_system && <span className="ml-1 text-xs text-blue-500">(system)</span>}
                      </Label>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditingUser(null)}>{t('common.cancel')}</Button>
            <Button onClick={handleEdit} disabled={saving}>
              {saving ? t('common.loading') : t('common.save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={!!deletingUser} onOpenChange={(open) => { if (!open) setDeletingUser(null) }}>
        <DialogContent className="sm:max-w-[400px]">
          <DialogHeader>
            <DialogTitle>{t('users.deleteTitle')}</DialogTitle>
            <DialogDescription>
              {t('users.deleteConfirm', { email: deletingUser?.email })}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeletingUser(null)}>{t('common.cancel')}</Button>
            <Button variant="destructive" onClick={handleDelete} disabled={deleting}>
              {deleting ? t('common.loading') : t('common.delete')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default Users
