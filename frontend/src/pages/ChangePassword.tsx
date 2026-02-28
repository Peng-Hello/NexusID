import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Card } from '../components/ui/card'

function ChangePassword() {
  const { t } = useTranslation()
  const [formData, setFormData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess(false)

    // Client-side validation
    if (formData.newPassword.length < 8) {
      setError(t('changePassword.passwordTooShort'))
      return
    }

    if (formData.newPassword !== formData.confirmPassword) {
      setError(t('changePassword.passwordsNotMatch'))
      return
    }

    setIsLoading(true)

    try {
      const token = localStorage.getItem('accessToken')
      const response = await fetch('/api/v1/user/change-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          old_password: formData.currentPassword,
          new_password: formData.newPassword,
        }),
      })

      const data = await response.json()

      if (!response.ok) {
        if (data.error === 'incorrect old password') {
          throw new Error(t('changePassword.incorrectPassword'))
        }
        throw new Error(data.error || t('changePassword.error'))
      }

      setSuccess(true)
      // Clear form
      setFormData({
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : t('changePassword.error'))
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto">
      <Card className="p-6">
        <h2 className="text-2xl font-bold text-gray-900 mb-2">
          {t('changePassword.title')}
        </h2>
        <p className="text-sm text-gray-600 mb-6">
          {t('changePassword.description')}
        </p>

        {success && (
          <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded mb-4">
            {t('changePassword.success')}
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="currentPassword" className="block text-sm font-medium text-gray-700 mb-1">
              {t('changePassword.currentPasswordLabel')}
            </label>
            <Input
              id="currentPassword"
              type="password"
              autoComplete="current-password"
              required
              value={formData.currentPassword}
              onChange={(e) => setFormData({ ...formData, currentPassword: e.target.value })}
              placeholder={t('changePassword.currentPasswordPlaceholder')}
            />
          </div>

          <div>
            <label htmlFor="newPassword" className="block text-sm font-medium text-gray-700 mb-1">
              {t('changePassword.newPasswordLabel')}
            </label>
            <Input
              id="newPassword"
              type="password"
              autoComplete="new-password"
              required
              value={formData.newPassword}
              onChange={(e) => setFormData({ ...formData, newPassword: e.target.value })}
              placeholder={t('changePassword.newPasswordPlaceholder')}
            />
            <p className="text-xs text-gray-500 mt-1">{t('changePassword.passwordHint')}</p>
          </div>

          <div>
            <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-1">
              {t('changePassword.confirmPasswordLabel')}
            </label>
            <Input
              id="confirmPassword"
              type="password"
              autoComplete="new-password"
              required
              value={formData.confirmPassword}
              onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
              placeholder={t('changePassword.confirmPasswordPlaceholder')}
            />
          </div>

          <Button
            type="submit"
            disabled={isLoading}
            className="w-full"
          >
            {isLoading ? t('changePassword.submitting') : t('changePassword.submit')}
          </Button>
        </form>
      </Card>
    </div>
  )
}

export default ChangePassword
