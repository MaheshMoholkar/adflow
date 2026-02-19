'use client';

import { useI18n } from './i18n-context';

export default function NotFound() {
  const { t } = useI18n();

  return (
    <main className="container">
      <div className="section not-found">
        <p className="label">{t('notFound.label')}</p>
        <h1>{t('notFound.title')}</h1>
        <p className="description">{t('notFound.text')}</p>
      </div>
    </main>
  );
}
