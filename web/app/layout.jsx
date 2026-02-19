import './globals.css';
import { I18nProvider } from './i18n-context';
import LanguageToggle from './language-toggle';

export const metadata = {
  title: 'AdFlow Landing',
  description: 'AdFlow customer landing page',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <I18nProvider>
          <LanguageToggle />
          {children}
        </I18nProvider>
      </body>
    </html>
  );
}
