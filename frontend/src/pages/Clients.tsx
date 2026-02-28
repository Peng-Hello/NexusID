import { useEffect, useState, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Plus, Copy, Key, Check } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '../components/ui/dialog'
import { Label } from '../components/ui/label'
import { Checkbox } from '../components/ui/checkbox'
import { Textarea } from '../components/ui/textarea'

interface Client {
  id: number
  client_id: string
  name: string
  redirect_uris: string[]
  scopes: string[]
  grant_types: string[]
  is_public: boolean
  is_active: boolean
  created_at: string
}

function Clients() {
  const { tenantId } = useParams()
  const { t } = useTranslation()
  const [clients, setClients] = useState<Client[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')

  // Copy feedback
  const [copiedId, setCopiedId] = useState<string | null>(null)

  // Create Client Form State
  const [showCreateDialog, setShowCreateDialog] = useState(false)
  const [newClient, setNewClient] = useState({
    name: '',
    redirect_uris: '',
    scopes: 'openid profile email',
    grant_types: 'authorization_code refresh_token',
    is_public: false
  })
  const [createError, setCreateError] = useState('')
  const [createdClientSecret, setCreatedClientSecret] = useState('')

  // Edit Client State
  const [editingClient, setEditingClient] = useState<Client | null>(null)
  const [editForm, setEditForm] = useState({
    name: '', redirect_uris: '', scopes: '', grant_types: '', is_public: false, is_active: true
  })
  const [editError, setEditError] = useState('')
  const [saving, setSaving] = useState(false)

  // Rotate Secret State
  const [rotatingClient, setRotatingClient] = useState<Client | null>(null)
  const [newSecret, setNewSecret] = useState('')
  const [rotating, setRotating] = useState(false)
  const [rotateError, setRotateError] = useState('')

  const fetchClients = useCallback(async () => {
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/clients?page=1&page_size=50`, {
        headers: { 'Authorization': `Bearer ${token}` },
      })
      if (response.ok) {
        const data = await response.json()
        setClients(data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch clients:', error)
    } finally {
      setLoading(false)
    }
  }, [tenantId])

  useEffect(() => {
    if (tenantId) fetchClients()
  }, [tenantId, fetchClients])

  const filteredClients = clients.filter(client =>
    client.name.toLowerCase().includes(search.toLowerCase()) ||
    client.client_id.toLowerCase().includes(search.toLowerCase())
  )

  // Copy with visual feedback
  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text)
    setCopiedId(id)
    setTimeout(() => setCopiedId(null), 2000)
  }

  // Create client
  const handleCreateClient = async () => {
    setCreateError('')
    setCreatedClientSecret('')
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/clients`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          name: newClient.name,
          redirect_uris: newClient.redirect_uris.split(',').map(s => s.trim()).filter(Boolean),
          scopes: newClient.scopes.split(',').map(s => s.trim()).filter(Boolean),
          grant_types: newClient.grant_types.split(',').map(s => s.trim()).filter(Boolean),
          is_public: newClient.is_public
        })
      })
      const data = await response.json()
      if (response.ok) {
        setCreatedClientSecret(data.secret)
        fetchClients()
      } else {
        setCreateError(data.error || t('clients.createFailed'))
      }
    } catch {
      setCreateError(t('clients.createFailed'))
    }
  }

  const resetCreateForm = () => {
    setShowCreateDialog(false)
    setCreatedClientSecret('')
    setCreateError('')
    setNewClient({
      name: '',
      redirect_uris: '',
      scopes: 'openid profile email',
      grant_types: 'authorization_code refresh_token',
      is_public: false
    })
  }

  // Edit client
  const openEditDialog = (client: Client) => {
    setEditingClient(client)
    setEditForm({
      name: client.name,
      redirect_uris: (client.redirect_uris ?? []).join('\n'),
      scopes: (client.scopes ?? []).join(' '),
      grant_types: (client.grant_types ?? []).join(' '),
      is_public: client.is_public,
      is_active: client.is_active,
    })
    setEditError('')
  }

  const handleEdit = async () => {
    if (!editingClient) return
    setEditError('')
    setSaving(true)
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/clients/${editingClient.id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: editForm.name,
          redirect_uris: editForm.redirect_uris.split('\n').map(s => s.trim()).filter(Boolean),
          scopes: editForm.scopes.split(/[\s,]+/).filter(Boolean),
          grant_types: editForm.grant_types.split(/[\s,]+/).filter(Boolean),
          is_public: editForm.is_public,
          is_active: editForm.is_active,
        }),
      })
      const data = await response.json()
      if (response.ok) {
        setEditingClient(null)
        fetchClients()
      } else {
        setEditError(data.error || t('clients.updateFailed'))
      }
    } catch {
      setEditError(t('clients.updateFailed'))
    } finally {
      setSaving(false)
    }
  }

  // Rotate secret
  const handleRotateSecret = async () => {
    if (!rotatingClient) return
    setRotating(true)
    setRotateError('')
    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch(`/api/v1/tenants/${tenantId}/clients/${rotatingClient.id}/rotate-secret`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` },
      })
      const data = await response.json()
      if (response.ok) {
        setNewSecret(data.secret || data.client_secret || '')
      } else {
        setRotateError(data.error || t('clients.rotateFailed'))
      }
    } catch {
      setRotateError(t('clients.rotateFailed'))
    } finally {
      setRotating(false)
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{t('clients.title')}</h1>
          <p className="text-gray-600">{t('clients.subtitle')}</p>
        </div>

        {/* Create Client Dialog */}
        <Dialog open={showCreateDialog} onOpenChange={(open) => {
          if (!open) resetCreateForm()
          else setShowCreateDialog(true)
        }}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              {t('clients.addClient')}
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>{t('clients.registerTitle')}</DialogTitle>
              <DialogDescription>{t('clients.registerDesc')}</DialogDescription>
            </DialogHeader>

            {!createdClientSecret ? (
              <div className="grid gap-4 py-4">
                {createError && (
                  <div className="p-3 bg-red-50 text-red-600 text-sm rounded">{createError}</div>
                )}
                <div className="grid gap-2">
                  <Label htmlFor="name">{t('clients.appName')} *</Label>
                  <Input
                    id="name"
                    value={newClient.name}
                    onChange={e => setNewClient({ ...newClient, name: e.target.value })}
                    placeholder="e.g. My Awesome App"
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="redirect_uris">{t('clients.redirectUris')} *</Label>
                  <Textarea
                    id="redirect_uris"
                    value={newClient.redirect_uris}
                    onChange={e => setNewClient({ ...newClient, redirect_uris: e.target.value })}
                    placeholder="http://localhost:3000/callback"
                  />
                  <p className="text-xs text-gray-500">{t('clients.redirectUrisHint')}</p>
                </div>
                <div className="grid gap-2">
                  <Label>{t('clients.scopes')}</Label>
                  <div className="flex flex-wrap gap-4 mt-1">
                    {['openid', 'profile', 'email'].map(scope => (
                      <div key={scope} className="flex items-center space-x-2">
                        <Checkbox
                          id={`create-scope-${scope}`}
                          checked={newClient.scopes.split(/[\s,]+/).includes(scope)}
                          onCheckedChange={(checked) => {
                            const current = new Set(newClient.scopes.split(/[\s,]+/).filter(Boolean))
                            if (checked) current.add(scope)
                            else current.delete(scope)
                            setNewClient({ ...newClient, scopes: Array.from(current).join(' ') })
                          }}
                        />
                        <Label htmlFor={`create-scope-${scope}`} className="font-normal cursor-pointer text-sm">
                          {t(`clients.scope_${scope}`)}
                        </Label>
                      </div>
                    ))}
                  </div>
                  <p className="text-xs text-gray-500 mt-1">{t('clients.scopesHint')}</p>
                </div>
                <div className="grid gap-2 mt-2">
                  <Label>{t('clients.grantTypes')}</Label>
                  <div className="flex flex-wrap gap-4 mt-1">
                    {['authorization_code', 'client_credentials', 'refresh_token'].map(grant => (
                      <div key={grant} className="flex items-center space-x-2">
                        <Checkbox
                          id={`create-grant-${grant}`}
                          checked={newClient.grant_types.split(/[\s,]+/).includes(grant)}
                          onCheckedChange={(checked) => {
                            const current = new Set(newClient.grant_types.split(/[\s,]+/).filter(Boolean))
                            if (checked) current.add(grant)
                            else current.delete(grant)
                            setNewClient({ ...newClient, grant_types: Array.from(current).join(' ') })
                          }}
                        />
                        <Label htmlFor={`create-grant-${grant}`} className="font-normal cursor-pointer text-sm">
                          {t(`clients.grant_${grant}`)}
                        </Label>
                      </div>
                    ))}
                  </div>
                </div>
                <div className="flex items-center space-x-2 mt-2">
                  <Checkbox
                    id="is_public"
                    checked={newClient.is_public}
                    onCheckedChange={(checked) => setNewClient({ ...newClient, is_public: checked === true })}
                  />
                  <Label htmlFor="is_public" className="font-normal">{t('clients.publicClient')}</Label>
                </div>
                <p className="text-xs text-gray-500 ml-6">{t('clients.publicClientHint')}</p>
              </div>
            ) : (
              <div className="py-6 space-y-4">
                <div className="p-4 bg-green-50 text-green-800 rounded-md border border-green-200">
                  <h4 className="font-semibold flex items-center gap-2 mb-2">{t('clients.createdSuccess')}</h4>
                  <p className="text-sm mb-4">{t('clients.createdSecretHint')}</p>
                  <div>
                    <Label className="text-xs uppercase text-green-700 font-bold mb-1 block">{t('clients.clientSecret')}</Label>
                    <div className="flex gap-2">
                      <code className="flex-1 p-2 bg-white rounded border border-green-200 font-mono text-sm break-all">
                        {createdClientSecret}
                      </code>
                      <Button variant="outline" size="sm" onClick={() => copyToClipboard(createdClientSecret, 'created-secret')}>
                        {copiedId === 'created-secret' ? <Check className="w-4 h-4 text-green-600" /> : <Copy className="w-4 h-4" />}
                      </Button>
                    </div>
                    {copiedId === 'created-secret' && (
                      <p className="text-xs text-green-600 mt-1">{t('clients.copied')}</p>
                    )}
                  </div>
                </div>
              </div>
            )}

            <DialogFooter>
              {!createdClientSecret ? (
                <>
                  <Button variant="outline" onClick={resetCreateForm}>{t('common.cancel')}</Button>
                  <Button onClick={handleCreateClient} disabled={!newClient.name || !newClient.redirect_uris}>{t('clients.createClient')}</Button>
                </>
              ) : (
                <Button onClick={resetCreateForm}>{t('common.close')}</Button>
              )}
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className="mb-6">
        <Input
          placeholder={t('clients.searchPlaceholder')}
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
        <div className="grid grid-cols-1 gap-6">
          {filteredClients.map((client) => (
            <Card key={client.id}>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle>{client.name}</CardTitle>
                    <CardDescription className="flex items-center gap-2 mt-1">
                      <Key className="w-4 h-4" />
                      <code className="text-xs bg-gray-100 px-2 py-1 rounded">{client.client_id}</code>
                      <button
                        onClick={() => copyToClipboard(client.client_id, `cid-${client.id}`)}
                        className="text-gray-500 hover:text-gray-700"
                        title={t('common.copy')}
                      >
                        {copiedId === `cid-${client.id}` ? <Check className="w-3 h-3 text-green-600" /> : <Copy className="w-3 h-3" />}
                      </button>
                      {copiedId === `cid-${client.id}` && (
                        <span className="text-xs text-green-600">{t('clients.copied')}</span>
                      )}
                    </CardDescription>
                  </div>
                  <div className="flex gap-2">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${client.is_active
                      ? 'bg-green-100 text-green-800'
                      : 'bg-gray-100 text-gray-800'
                      }`}>
                      {client.is_active ? t('common.active') : t('common.inactive')}
                    </span>
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${client.is_public
                      ? 'bg-blue-100 text-blue-800'
                      : 'bg-purple-100 text-purple-800'
                      }`}>
                      {client.is_public ? t('common.public') : t('clients.confidential')}
                    </span>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div>
                    <p className="text-sm font-medium text-gray-700">{t('clients.redirectUris')}</p>
                    <div className="mt-1 flex flex-wrap gap-1">
                      {(client.redirect_uris ?? []).map((uri, idx) => (
                        <code key={idx} className="text-xs bg-gray-100 px-2 py-1 rounded">{uri}</code>
                      ))}
                    </div>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-700">{t('clients.scopes')}</p>
                    <div className="mt-1 flex flex-wrap gap-1">
                      {(client.scopes ?? []).map((scope, idx) => (
                        <span key={idx} className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-blue-100 text-blue-800">
                          {scope}
                        </span>
                      ))}
                    </div>
                  </div>
                  <div className="flex gap-2 pt-2">
                    <Button variant="outline" size="sm" onClick={() => openEditDialog(client)}>
                      {t('common.edit')}
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => { setRotatingClient(client); setNewSecret('') }}>
                      {t('clients.rotateSecret')}
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {!loading && filteredClients.length === 0 && (
        <Card>
          <CardContent className="text-center py-12">
            <p className="text-gray-500">{t('clients.noClients')}</p>
          </CardContent>
        </Card>
      )}

      {/* Edit Client Dialog */}
      <Dialog open={!!editingClient} onOpenChange={(open) => { if (!open) setEditingClient(null) }}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>{t('clients.editTitle')}</DialogTitle>
            <DialogDescription>{editingClient?.client_id}</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            {editError && (
              <div className="p-3 bg-red-50 text-red-600 text-sm rounded">{editError}</div>
            )}
            <div className="grid gap-2">
              <Label>{t('clients.appName')} *</Label>
              <Input value={editForm.name} onChange={e => setEditForm({ ...editForm, name: e.target.value })} />
            </div>
            <div className="grid gap-2">
              <Label>{t('clients.redirectUris')}</Label>
              <Textarea
                value={editForm.redirect_uris}
                onChange={e => setEditForm({ ...editForm, redirect_uris: e.target.value })}
                rows={3}
              />
              <p className="text-xs text-gray-500">{t('clients.redirectUrisHint')}</p>
            </div>
            <div className="grid gap-2 mt-2">
              <Label>{t('clients.scopes')}</Label>
              <div className="flex flex-wrap gap-4 mt-1 border rounded p-3 bg-gray-50">
                {['openid', 'profile', 'email'].map(scope => (
                  <div key={scope} className="flex items-center space-x-2">
                    <Checkbox
                      id={`edit-scope-${scope}`}
                      checked={editForm.scopes.split(/[\s,]+/).includes(scope)}
                      onCheckedChange={(checked) => {
                        const current = new Set(editForm.scopes.split(/[\s,]+/).filter(Boolean))
                        if (checked) current.add(scope)
                        else current.delete(scope)
                        setEditForm({ ...editForm, scopes: Array.from(current).join(' ') })
                      }}
                    />
                    <Label htmlFor={`edit-scope-${scope}`} className="font-normal cursor-pointer text-sm">
                      {t(`clients.scope_${scope}`)}
                    </Label>
                  </div>
                ))}
              </div>
            </div>
            <div className="grid gap-2 mt-2">
              <Label>{t('clients.grantTypes')}</Label>
              <div className="flex flex-wrap gap-4 mt-1 border rounded p-3 bg-gray-50">
                {['authorization_code', 'client_credentials', 'refresh_token'].map(grant => (
                  <div key={grant} className="flex items-center space-x-2">
                    <Checkbox
                      id={`edit-grant-${grant}`}
                      checked={editForm.grant_types.split(/[\s,]+/).includes(grant)}
                      onCheckedChange={(checked) => {
                        const current = new Set(editForm.grant_types.split(/[\s,]+/).filter(Boolean))
                        if (checked) current.add(grant)
                        else current.delete(grant)
                        setEditForm({ ...editForm, grant_types: Array.from(current).join(' ') })
                      }}
                    />
                    <Label htmlFor={`edit-grant-${grant}`} className="font-normal cursor-pointer text-sm">
                      {t(`clients.grant_${grant}`)}
                    </Label>
                  </div>
                ))}
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="edit-public" checked={editForm.is_public} onCheckedChange={(c) => setEditForm({ ...editForm, is_public: c === true })} />
              <Label htmlFor="edit-public" className="font-normal">{t('clients.publicClient')}</Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="edit-active" checked={editForm.is_active} onCheckedChange={(c) => setEditForm({ ...editForm, is_active: c === true })} />
              <Label htmlFor="edit-active" className="font-normal">{t('clients.isActive')}</Label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditingClient(null)}>{t('common.cancel')}</Button>
            <Button onClick={handleEdit} disabled={!editForm.name || saving}>
              {saving ? t('common.loading') : t('common.save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Rotate Secret Dialog */}
      <Dialog open={!!rotatingClient} onOpenChange={(open) => { if (!open) { setRotatingClient(null); setNewSecret(''); setRotateError('') } }}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>{t('clients.rotateSecretTitle')}</DialogTitle>
            <DialogDescription>{rotatingClient?.name}</DialogDescription>
          </DialogHeader>

          {!newSecret ? (
            <div className="py-4">
              {rotateError && (
                <div className="p-3 mb-4 bg-red-50 text-red-600 text-sm rounded">{rotateError}</div>
              )}
              <p className="text-sm text-gray-600 mb-4">{t('clients.rotateSecretWarning')}</p>
              <DialogFooter>
                <Button variant="outline" onClick={() => setRotatingClient(null)}>{t('common.cancel')}</Button>
                <Button variant="destructive" onClick={handleRotateSecret} disabled={rotating}>
                  {rotating ? t('common.loading') : t('clients.rotateSecret')}
                </Button>
              </DialogFooter>
            </div>
          ) : (
            <div className="py-4 space-y-4">
              <div className="p-4 bg-green-50 text-green-800 rounded-md border border-green-200">
                <p className="text-sm mb-3">{t('clients.createdSecretHint')}</p>
                <Label className="text-xs uppercase text-green-700 font-bold mb-1 block">{t('clients.clientSecret')}</Label>
                <div className="flex gap-2">
                  <code className="flex-1 p-2 bg-white rounded border border-green-200 font-mono text-sm break-all">
                    {newSecret}
                  </code>
                  <Button variant="outline" size="sm" onClick={() => copyToClipboard(newSecret, 'rotated-secret')}>
                    {copiedId === 'rotated-secret' ? <Check className="w-4 h-4 text-green-600" /> : <Copy className="w-4 h-4" />}
                  </Button>
                </div>
                {copiedId === 'rotated-secret' && (
                  <p className="text-xs text-green-600 mt-1">{t('clients.copied')}</p>
                )}
              </div>
              <DialogFooter>
                <Button onClick={() => { setRotatingClient(null); setNewSecret('') }}>{t('common.close')}</Button>
              </DialogFooter>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default Clients
