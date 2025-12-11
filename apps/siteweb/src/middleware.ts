import createMiddleware from 'next-intl/middleware';
import { locales, defaultLocale } from '../i18n/request';

export default createMiddleware({
  // Liste des locales supportées
  locales,

  // Locale par défaut
  defaultLocale,

  // Préfixer les URLs avec la locale (ex: /fr/about)
  localePrefix: 'as-needed',
});

export const config = {
  // Matcher pour les routes à internationaliser
  // Exclut les fichiers statiques, api, etc.
  matcher: [
    // Match all pathnames except for
    // - … if they start with `/api`, `/_next` or `/_vercel`
    // - … the ones containing a dot (e.g. `favicon.ico`)
    '/((?!api|_next|_vercel|.*\\..*).*)',
  ],
};
