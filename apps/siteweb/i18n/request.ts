import { getRequestConfig } from 'next-intl/server';

export const locales = ['fr', 'en'] as const;
export const defaultLocale = 'fr' as const;

export type Locale = (typeof locales)[number];

export default getRequestConfig(async () => {
  // For now, use French as default
  // In production, this would be detected from the URL or Accept-Language header
  const locale = defaultLocale;

  return {
    locale,
    messages: (await import(`../messages/${locale}.json`)).default,
  };
});
