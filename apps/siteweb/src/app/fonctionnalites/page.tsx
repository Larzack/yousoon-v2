'use client';

import { motion } from 'framer-motion';
import { useTranslations } from 'next-intl';
import {
  MapPin,
  Percent,
  Zap,
  Heart,
  Users,
  Shield,
  Smartphone,
  Bell,
  QrCode,
  Search,
  Star,
  Clock,
} from 'lucide-react';
import { CTA } from '@/components/sections/CTA';

const fadeInUp = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  transition: { duration: 0.5 },
};

const staggerContainer = {
  animate: {
    transition: {
      staggerChildren: 0.1,
    },
  },
};

const features = [
  {
    icon: MapPin,
    titleKey: 'geolocatedOffers',
    descriptionKey: 'geolocatedOffersDesc',
    color: 'text-blue-400',
    bgColor: 'bg-blue-400/10',
  },
  {
    icon: Percent,
    titleKey: 'exclusiveDiscounts',
    descriptionKey: 'exclusiveDiscountsDesc',
    color: 'text-primary',
    bgColor: 'bg-primary/10',
  },
  {
    icon: Zap,
    titleKey: 'instantBooking',
    descriptionKey: 'instantBookingDesc',
    color: 'text-yellow-400',
    bgColor: 'bg-yellow-400/10',
  },
  {
    icon: QrCode,
    titleKey: 'qrCheckin',
    descriptionKey: 'qrCheckinDesc',
    color: 'text-green-400',
    bgColor: 'bg-green-400/10',
  },
  {
    icon: Heart,
    titleKey: 'favorites',
    descriptionKey: 'favoritesDesc',
    color: 'text-red-400',
    bgColor: 'bg-red-400/10',
  },
  {
    icon: Search,
    titleKey: 'smartSearch',
    descriptionKey: 'smartSearchDesc',
    color: 'text-purple-400',
    bgColor: 'bg-purple-400/10',
  },
  {
    icon: Bell,
    titleKey: 'notifications',
    descriptionKey: 'notificationsDesc',
    color: 'text-orange-400',
    bgColor: 'bg-orange-400/10',
  },
  {
    icon: Star,
    titleKey: 'reviews',
    descriptionKey: 'reviewsDesc',
    color: 'text-amber-400',
    bgColor: 'bg-amber-400/10',
  },
  {
    icon: Users,
    titleKey: 'community',
    descriptionKey: 'communityDesc',
    color: 'text-cyan-400',
    bgColor: 'bg-cyan-400/10',
  },
  {
    icon: Clock,
    titleKey: 'history',
    descriptionKey: 'historyDesc',
    color: 'text-indigo-400',
    bgColor: 'bg-indigo-400/10',
  },
  {
    icon: Shield,
    titleKey: 'security',
    descriptionKey: 'securityDesc',
    color: 'text-emerald-400',
    bgColor: 'bg-emerald-400/10',
  },
  {
    icon: Smartphone,
    titleKey: 'offlineMode',
    descriptionKey: 'offlineModeDesc',
    color: 'text-rose-400',
    bgColor: 'bg-rose-400/10',
  },
];

export default function FeaturesPage() {
  const t = useTranslations('features');

  return (
    <>
      {/* Hero Section */}
      <section className="relative overflow-hidden py-24 md:py-32">
        <div className="container mx-auto px-4">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="mx-auto max-w-3xl text-center"
          >
            <h1 className="text-4xl font-bold tracking-tight md:text-5xl lg:text-6xl">
              {t('pageTitle')}
            </h1>
            <p className="mt-6 text-lg text-muted-foreground md:text-xl">
              {t('pageDescription')}
            </p>
          </motion.div>
        </div>
      </section>

      {/* Features Grid */}
      <section className="py-16 md:py-24">
        <div className="container mx-auto px-4">
          <motion.div
            variants={staggerContainer}
            initial="initial"
            whileInView="animate"
            viewport={{ once: true }}
            className="grid gap-6 md:grid-cols-2 lg:grid-cols-3"
          >
            {features.map((feature, index) => (
              <motion.div
                key={feature.titleKey}
                variants={fadeInUp}
                className="group relative rounded-2xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm transition-all duration-300 hover:border-primary/50 hover:bg-white/10"
              >
                {/* Icon */}
                <div
                  className={`mb-4 inline-flex rounded-xl ${feature.bgColor} p-3`}
                >
                  <feature.icon className={`h-6 w-6 ${feature.color}`} />
                </div>

                {/* Content */}
                <h3 className="mb-2 text-xl font-semibold">
                  {t(`items.${feature.titleKey}.title`)}
                </h3>
                <p className="text-muted-foreground">
                  {t(`items.${feature.titleKey}.description`)}
                </p>

                {/* Hover glow */}
                <div className="absolute -inset-px -z-10 rounded-2xl bg-gradient-to-r from-primary/0 via-primary/10 to-primary/0 opacity-0 blur-xl transition-opacity duration-500 group-hover:opacity-100" />
              </motion.div>
            ))}
          </motion.div>
        </div>
      </section>

      {/* App Preview Section */}
      <section className="py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="grid items-center gap-12 lg:grid-cols-2">
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.6 }}
            >
              <h2 className="text-3xl font-bold md:text-4xl">
                {t('appPreview.title')}
              </h2>
              <p className="mt-4 text-lg text-muted-foreground">
                {t('appPreview.description')}
              </p>
              <ul className="mt-8 space-y-4">
                {['point1', 'point2', 'point3', 'point4'].map((point) => (
                  <li key={point} className="flex items-start gap-3">
                    <div className="mt-1 rounded-full bg-primary/20 p-1">
                      <svg
                        className="h-4 w-4 text-primary"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M5 13l4 4L19 7"
                        />
                      </svg>
                    </div>
                    <span>{t(`appPreview.${point}`)}</span>
                  </li>
                ))}
              </ul>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, x: 20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="relative"
            >
              {/* Phone mockup placeholder */}
              <div className="relative mx-auto max-w-xs">
                <div className="phone-frame aspect-[9/19] rounded-[3rem] border-4 border-white/20 bg-gradient-to-b from-white/10 to-white/5 p-4 shadow-2xl">
                  <div className="h-full w-full rounded-[2.5rem] bg-background">
                    {/* App screenshot placeholder */}
                    <div className="flex h-full items-center justify-center text-muted-foreground">
                      <Smartphone className="h-16 w-16" />
                    </div>
                  </div>
                </div>
                {/* Glow effect */}
                <div className="absolute -inset-10 -z-10 bg-primary/20 blur-3xl" />
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      <CTA />
    </>
  );
}
