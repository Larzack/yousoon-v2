'use client';

import { motion } from 'framer-motion';
import Image from 'next/image';
import Link from 'next/link';

interface AppStoreBadgesProps {
  size?: 'sm' | 'md' | 'lg';
  direction?: 'horizontal' | 'vertical';
  showLabels?: boolean;
  className?: string;
}

const sizeClasses = {
  sm: { width: 120, height: 40 },
  md: { width: 150, height: 50 },
  lg: { width: 180, height: 60 },
};

export function AppStoreBadges({
  size = 'md',
  direction = 'horizontal',
  showLabels = false,
  className = '',
}: AppStoreBadgesProps) {
  const { width, height } = sizeClasses[size];

  const stores = [
    {
      name: 'App Store',
      href: 'https://apps.apple.com/app/yousoon',
      image: '/images/app-store-badge.svg',
      alt: 'Télécharger sur l\'App Store',
    },
    {
      name: 'Google Play',
      href: 'https://play.google.com/store/apps/details?id=com.yousoon.yousoon',
      image: '/images/google-play-badge.svg',
      alt: 'Disponible sur Google Play',
    },
  ];

  return (
    <div
      className={`flex ${
        direction === 'horizontal' ? 'flex-row' : 'flex-col'
      } gap-4 ${className}`}
    >
      {stores.map((store, index) => (
        <motion.div
          key={store.name}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: index * 0.1 }}
          className="flex flex-col items-center gap-2"
        >
          <Link
            href={store.href}
            target="_blank"
            rel="noopener noreferrer"
            className="group relative block overflow-hidden rounded-lg transition-transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 focus:ring-offset-background"
          >
            {/* Glow effect on hover */}
            <div className="absolute inset-0 bg-primary/20 opacity-0 blur-xl transition-opacity group-hover:opacity-100" />
            
            <Image
              src={store.image}
              alt={store.alt}
              width={width}
              height={height}
              className="relative z-10"
            />
          </Link>
          {showLabels && (
            <span className="text-sm text-muted-foreground">{store.name}</span>
          )}
        </motion.div>
      ))}
    </div>
  );
}

// Composant simplifié pour les badges textuels (fallback sans images)
export function AppStoreBadgesText({
  className = '',
}: {
  className?: string;
}) {
  const stores = [
    {
      name: 'App Store',
      href: 'https://apps.apple.com/app/yousoon',
      icon: (
        <svg className="h-6 w-6" viewBox="0 0 24 24" fill="currentColor">
          <path d="M18.71 19.5C17.88 20.74 17 21.95 15.66 21.97C14.32 22 13.89 21.18 12.37 21.18C10.84 21.18 10.37 21.95 9.1 22C7.79 22.05 6.8 20.68 5.96 19.47C4.25 17 2.94 12.45 4.7 9.39C5.57 7.87 7.13 6.91 8.82 6.88C10.1 6.86 11.32 7.75 12.11 7.75C12.89 7.75 14.37 6.68 15.92 6.84C16.57 6.87 18.39 7.1 19.56 8.82C19.47 8.88 17.39 10.1 17.41 12.63C17.44 15.65 20.06 16.66 20.09 16.67C20.06 16.74 19.67 18.11 18.71 19.5ZM13 3.5C13.73 2.67 14.94 2.04 15.94 2C16.07 3.17 15.6 4.35 14.9 5.19C14.21 6.04 13.07 6.7 11.95 6.61C11.8 5.46 12.36 4.26 13 3.5Z" />
        </svg>
      ),
    },
    {
      name: 'Google Play',
      href: 'https://play.google.com/store/apps/details?id=com.yousoon.yousoon',
      icon: (
        <svg className="h-6 w-6" viewBox="0 0 24 24" fill="currentColor">
          <path d="M3,20.5V3.5C3,2.91 3.34,2.39 3.84,2.15L13.69,12L3.84,21.85C3.34,21.6 3,21.09 3,20.5M16.81,15.12L6.05,21.34L14.54,12.85L16.81,15.12M20.16,10.81C20.5,11.08 20.75,11.5 20.75,12C20.75,12.5 20.5,12.92 20.16,13.19L17.89,14.5L15.39,12L17.89,9.5L20.16,10.81M6.05,2.66L16.81,8.88L14.54,11.15L6.05,2.66Z" />
        </svg>
      ),
    },
  ];

  return (
    <div className={`flex flex-wrap gap-4 ${className}`}>
      {stores.map((store) => (
        <Link
          key={store.name}
          href={store.href}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-lg border border-white/20 bg-white/5 px-4 py-2 text-white backdrop-blur-sm transition-all hover:border-primary hover:bg-white/10"
        >
          {store.icon}
          <div className="flex flex-col">
            <span className="text-xs text-white/60">Télécharger sur</span>
            <span className="font-semibold">{store.name}</span>
          </div>
        </Link>
      ))}
    </div>
  );
}
