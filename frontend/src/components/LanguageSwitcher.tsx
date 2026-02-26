import { useTranslation } from 'react-i18next'
import { Globe } from 'lucide-react'
import { Button } from './ui/button'

function LanguageSwitcher() {
    const { i18n } = useTranslation()

    const toggleLanguage = () => {
        const nextLang = i18n.language?.startsWith('zh') ? 'en' : 'zh'
        i18n.changeLanguage(nextLang)
    }

    const currentLabel = i18n.language?.startsWith('zh') ? '中文' : 'EN'

    return (
        <Button
            variant="ghost"
            size="sm"
            onClick={toggleLanguage}
            className="gap-1.5 text-gray-600 hover:text-gray-900"
            title={i18n.language?.startsWith('zh') ? 'Switch to English' : '切换到中文'}
        >
            <Globe className="w-4 h-4" />
            {currentLabel}
        </Button>
    )
}

export default LanguageSwitcher
