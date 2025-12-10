import type { Metadata, Viewport } from 'next';
import { Inter, Poppins } from 'next/font/google';
import { NextIntlClientProvider } from 'next-intl';
import { getMessages, getLocale } from 'next-intl/server';
import { Header } from '@/components/layout/Header';
import { Footer } from '@/components/layout/Footer';
import '@/styles/globals.css';

// Fonts
const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
  display: 'swap',
});

// Using Poppins as a fallback for Futura (similar geometric sans-serif)
const poppins = Poppins({
  subsets: ['latin'],
  weight: ['500', '700'],
  variable: '--font-futura',
  display: 'swap',
});

export const metadata: Metadata = {
  metadataBase: new URL('https://www.yousoon.com'),
  title: {
    default: 'Yousoon - Sorties avec réductions',
    template: '%s | Yousoon',
  },
  description:
    'Découvrez des sorties uniques à prix réduits. Bars, restaurants, événements... Profitez de réductions exclusives avec Yousoon.',
  keywords: [
    'sorties',
    'réductions',
    'bars',
    'restaurants',
    'événements',
    'bons plans',
    'lifestyle',
    'partenaires',
  ],
  authors: [{ name: 'Yousoon' }],
  creator: 'Yousoon',
  publisher: 'Yousoon',
  formatDetection: {
    email: false,
    address: false,
    telephone: false,
  },
  openGraph: {
    type: 'website',
    locale: 'fr_FR',
    alternateLocale: 'en_US',
    url: 'https://www.yousoon.com',
    siteName: 'Yousoon',
    title: 'Yousoon - Sorties avec réductions',
    description:
      'Découvrez des sorties uniques à prix réduits. Bars, restaurants, événements... Profitez de réductions exclusives.',
    images: [
      {
        url: '/images/og-image.jpg',
        width: 1200,
        height: 630,
        alt: 'Yousoon - Sorties avec réductions',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Yousoon - Sorties avec réductions',
    description:
      'Découvrez des sorties uniques à prix réduits. Bars, restaurants, événements...',
    images: ['/images/twitter-image.jpg'],
    creator: '@yousoon',
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-video-preview': -1,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
  icons: {
    icon: '/favicon.ico',
    shortcut: '/favicon-16x16.png',
    apple: '/apple-touch-icon.png',
  },
  manifest: '/manifest.json',
};

export const viewport: Viewport = {
  themeColor: '#000000',
  width: 'device-width',
  initialScale: 1,
  maximumScale: 5,
};

interface RootLayoutProps {
  children: React.ReactNode;
}

export default async function RootLayout({ children }: RootLayoutProps) {
  const locale = await getLocale();
  const messages = await getMessages();

  return (
    <html lang={locale} className="dark scroll-smooth">
      <body
        className={`${inter.variable} ${futura.variable} font-sans antialiased`}
      >
        <NextIntlClientProvider messages={messages}>
          <div className="relative flex min-h-screen flex-col bg-background text-foreground">
            {/* Background gradient */}
            <div className="pointer-events-none fixed inset-0 bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(233,155,39,0.1),rgba(0,0,0,0))]" />
            
            <Header />
            <main className="relative z-10 flex-1">{children}</main>
            <Footer />
          </div>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
