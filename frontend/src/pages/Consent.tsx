import { useEffect, useState } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { Shield, ShieldAlert } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../components/ui/card'

function Consent() {
    const [searchParams] = useSearchParams()
    const navigate = useNavigate()
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')
    const [clientName, setClientName] = useState('')

    const clientId = searchParams.get('client_id')
    const redirectUri = searchParams.get('redirect_uri')
    const scope = searchParams.get('scope')
    const state = searchParams.get('state')
    const responseType = searchParams.get('response_type')

    useEffect(() => {
        // Validate required parameters
        if (!clientId || !redirectUri || responseType !== 'code') {
            setError('Invalid authorization request. Missing or invalid parameters.')
            return
        }

        // Check if user is logged in
        const token = localStorage.getItem('accessToken')
        if (!token) {
            const returnUrl = encodeURIComponent(window.location.pathname + window.location.search)
            navigate(`/login?returnUrl=${returnUrl}`)
            return
        }

        // Fetch real client name
        const fetchClientName = async () => {
            try {
                // The user's tenant_id is in their stored user object
                const userStr = localStorage.getItem('user')
                const user = userStr ? JSON.parse(userStr) : null
                const tenantId = user?.tenant_id
                if (!tenantId) {
                    setClientName(clientId)
                    return
                }
                const res = await fetch(`/api/v1/tenants/${tenantId}/clients?page=1&page_size=100`, {
                    headers: { 'Authorization': `Bearer ${token}` },
                })
                if (res.ok) {
                    const data = await res.json()
                    const clients: { client_id: string; name: string }[] = data.data || []
                    const found = clients.find(c => c.client_id === clientId)
                    setClientName(found?.name || clientId || 'Unknown Application')
                } else {
                    setClientName(clientId || 'Unknown Application')
                }
            } catch {
                setClientName(clientId || 'Unknown Application')
            } finally {
                // done loading
            }
        }

        fetchClientName()
    }, [clientId, redirectUri, responseType, navigate])

    const scopesList = scope ? scope.split(' ') : []

    const handleDecision = async (approved: boolean) => {
        setLoading(true)
        setError('')

        try {
            const token = localStorage.getItem('accessToken')

            if (!approved) {
                // Redirect back with access_denied error
                const url = new URL(redirectUri!)
                url.searchParams.append('error', 'access_denied')
                url.searchParams.append('error_description', 'The user denied the request')
                if (state) url.searchParams.append('state', state)
                window.location.href = url.toString()
                return
            }

            // Submit consent to backend
            const response = await fetch('/api/v1/oauth/consent', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    client_id: clientId,
                    redirect_uri: redirectUri,
                    scope: scope,
                    state: state,
                    response_type: responseType
                })
            })

            const data = await response.json()

            if (response.ok && data.redirect_url) {
                // Backend should return the full redirect URL with the authorization code
                window.location.href = data.redirect_url
            } else {
                setError(data.error_description || data.error || 'Failed to authorize application')
            }
        } catch {
            setError('A network error occurred while processing your request.')
        } finally {
            setLoading(false)
        }
    }

    if (error && !clientId) {
        return (
            <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
                <div className="sm:mx-auto sm:w-full sm:max-w-md">
                    <Card className="border-red-200">
                        <CardHeader className="text-center">
                            <ShieldAlert className="w-12 h-12 text-red-500 mx-auto mb-4" />
                            <CardTitle className="text-red-700">Authorization Error</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-center text-gray-600">{error}</p>
                        </CardContent>
                    </Card>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
            <div className="sm:mx-auto sm:w-full sm:max-w-md">
                <div className="flex justify-center mb-6">
                    <div className="w-16 h-16 bg-blue-600 rounded-xl flex items-center justify-center shadow-lg">
                        <Shield className="w-8 h-8 text-white" />
                    </div>
                </div>
                <h2 className="mt-2 text-center text-3xl font-extrabold text-gray-900">
                    NexusID Auth
                </h2>
            </div>

            <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
                <Card>
                    <CardHeader className="text-center pb-2">
                        <CardTitle className="text-xl">Request for Permission</CardTitle>
                        <CardDescription className="mt-2 text-base">
                            <strong className="text-gray-900">{clientName}</strong> is requesting access to your NexusID account.
                        </CardDescription>
                    </CardHeader>

                    <CardContent className="pt-4">
                        {error && (
                            <div className="mb-4 p-3 bg-red-50 text-red-600 text-sm rounded border border-red-200">
                                {error}
                            </div>
                        )}

                        <h3 className="text-sm font-semibold text-gray-900 mb-3 uppercase tracking-wider">
                            This application will be able to:
                        </h3>

                        <ul className="space-y-3 mb-6">
                            {scopesList.includes('openid') && (
                                <li className="flex items-start">
                                    <span className="flex-shrink-0 w-5 h-5 flex items-center justify-center rounded-full bg-blue-100 text-blue-600 mt-0.5 mr-3">✓</span>
                                    <span className="text-sm text-gray-600">Authenticate your identity using NexusID</span>
                                </li>
                            )}
                            {scopesList.includes('profile') && (
                                <li className="flex items-start">
                                    <span className="flex-shrink-0 w-5 h-5 flex items-center justify-center rounded-full bg-blue-100 text-blue-600 mt-0.5 mr-3">✓</span>
                                    <span className="text-sm text-gray-600">Access your basic profile information (name)</span>
                                </li>
                            )}
                            {scopesList.includes('email') && (
                                <li className="flex items-start">
                                    <span className="flex-shrink-0 w-5 h-5 flex items-center justify-center rounded-full bg-blue-100 text-blue-600 mt-0.5 mr-3">✓</span>
                                    <span className="text-sm text-gray-600">Access your email address</span>
                                </li>
                            )}
                            {!scopesList.length && (
                                <li className="text-sm text-gray-500 italic">No specific permissions requested.</li>
                            )}
                        </ul>

                        <div className="text-xs text-gray-500 bg-gray-50 p-3 rounded">
                            <p>By clicking <strong>Allow</strong>, you allow this app to use your information in accordance with their terms of service and privacy policies.</p>
                        </div>
                    </CardContent>

                    <CardFooter className="flex flex-col gap-3 sm:flex-row-reverse sm:justify-start">
                        <Button
                            className="w-full sm:w-auto"
                            onClick={() => handleDecision(true)}
                            disabled={loading || !!error}
                        >
                            Allow Access
                        </Button>
                        <Button
                            variant="outline"
                            className="w-full sm:w-auto"
                            onClick={() => handleDecision(false)}
                            disabled={loading}
                        >
                            Cancel
                        </Button>
                    </CardFooter>
                </Card>
            </div>
        </div>
    )
}

export default Consent
