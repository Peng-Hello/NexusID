import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Plus, Copy, Key } from 'lucide-react'
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

  useEffect(() => {
    if (tenantId) fetchClients()
  }, [tenantId])

  const fetchClients = async () => {
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
  }

  const filteredClients = clients.filter(client =>
    client.name.toLowerCase().includes(search.toLowerCase()) ||
    client.client_id.toLowerCase().includes(search.toLowerCase())
  )

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
  }

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
        setCreateError(data.error || 'Failed to create client')
      }
    } catch {
      setCreateError('Network error occurred')
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

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{t('clients.title')}</h1>
          <p className="text-gray-600">{t('clients.subtitle')}</p>
        </div>

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
              <DialogDescription>
                {t('clients.registerDesc')}
              </DialogDescription>
            </DialogHeader>

            {!createdClientSecret ? (
              <div className="grid gap-4 py-4">
                {createError && (
                  <div className="p-3 bg-red-50 text-red-600 text-sm rounded">
                    {createError}
                  </div>
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
                  <Label htmlFor="scopes">{t('clients.scopes')}</Label>
                  <Input
                    id="scopes"
                    value={newClient.scopes}
                    onChange={e => setNewClient({ ...newClient, scopes: e.target.value })}
                  />
                  <p className="text-xs text-gray-500">{t('clients.scopesHint')}</p>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="grant_types">{t('clients.grantTypes')}</Label>
                  <Input
                    id="grant_types"
                    value={newClient.grant_types}
                    onChange={e => setNewClient({ ...newClient, grant_types: e.target.value })}
                  />
                </div>
                <div className="flex items-center space-x-2 mt-2">
                  <Checkbox
                    id="is_public"
                    checked={newClient.is_public}
                    onCheckedChange={(checked) => setNewClient({ ...newClient, is_public: checked === true })}
                  />
                  <Label htmlFor="is_public" className="font-normal">
                    {t('clients.publicClient')}
                  </Label>
                </div>
                <p className="text-xs text-gray-500 ml-6">{t('clients.publicClientHint')}</p>
              </div>
            ) : (
              <div className="py-6 space-y-4">
                <div className="p-4 bg-green-50 text-green-800 rounded-md border border-green-200">
                  <h4 className="font-semibold flex items-center gap-2 mb-2">
                    {t('clients.createdSuccess')}
                  </h4>
                  <p className="text-sm mb-4">{t('clients.createdSecretHint')}</p>

                  <div className="space-y-4">
                    <div>
                      <Label className="text-xs uppercase text-green-700 font-bold mb-1 block">{t('clients.clientSecret')}</Label>
                      <div className="flex gap-2">
                        <code className="flex-1 p-2 bg-white rounded border border-green-200 font-mono text-sm break-all">
                          {createdClientSecret}
                        </code>
                        <Button variant="outline" size="sm" onClick={() => copyToClipboard(createdClientSecret)}>
                          <Copy className="w-4 h-4" />
                        </Button>
                      </div>
                    </div>
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
                      <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                        {client.client_id}
                      </code>
                      <button
                        onClick={() => copyToClipboard(client.client_id)}
                        className="text-gray-500 hover:text-gray-700"
                      >
                        <Copy className="w-3 h-3" />
                      </button>
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
                      {client.is_public ? t('common.public') : 'Confidential'}
                    </span>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div>
                    <p className="text-sm font-medium text-gray-700">{t('clients.redirectUris')}</p>
                    <div className="mt-1 flex flex-wrap gap-1">
                      {client.redirect_uris.map((uri, idx) => (
                        <code key={idx} className="text-xs bg-gray-100 px-2 py-1 rounded">
                          {uri}
                        </code>
                      ))}
                    </div>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-700">{t('clients.scopes')}</p>
                    <div className="mt-1 flex flex-wrap gap-1">
                      {client.scopes.map((scope, idx) => (
                        <span key={idx} className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-blue-100 text-blue-800">
                          {scope}
                        </span>
                      ))}
                    </div>
                  </div>
                  <div className="flex gap-2 pt-2">
                    <Button variant="outline" size="sm">
                      {t('common.edit')}
                    </Button>
                    <Button variant="outline" size="sm">
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
    </div>
  )
}

export default Clients
