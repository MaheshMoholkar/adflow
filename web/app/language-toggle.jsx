'use client';

import { useI18n } from './i18n-context';

export default function LanguageToggle() {
  const { language, setLanguage, t } = useI18n();

  return (
    <div className="lang-toggle" role="group" aria-label="Language selector">
      <button
        type="button"
        className={`lang-btn ${language === 'en' ? 'active' : ''}`}
        onClick={() => setLanguage('en')}
      >
        {t('toggle.en')}
      </button>
      <button
        type="button"
        className={`lang-btn ${language === 'mr' ? 'active' : ''}`}
        onClick={() => setLanguage('mr')}
      >
        {t('toggle.mr')}
      </button>
    </div>
  );
}
