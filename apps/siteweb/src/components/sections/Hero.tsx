'use client'

import { useTranslations } from 'next-intl'
import { motion } from 'framer-motion'
import { ArrowRight, ChevronDown } from 'lucide-react'
import { Button } from '@/components/ui'
import { STATS, SITE_CONFIG } from '@/lib/constants'
import { formatNumber } from '@/lib/utils'
import { AppStoreBadges } from '@/components/shared/AppStoreBadges'

export function Hero() {
  const t = useTranslations('hero')

  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden bg-black">
      {/* Background gradient */}
      <div className="absolute inset-0 bg-hero-gradient" />
      
      {/* Animated background circles */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-primary/10 rounded-full blur-3xl" />
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-primary/5 rounded-full blur-3xl" />
      </div>

      <div className="container-wide relative z-10 pt-32 pb-20">
        <div className="grid lg:grid-cols-2 gap-12 lg:gap-20 items-center">
          {/* Left content */}
          <motion.div
            initial={{ opacity: 0, x: -50 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.6 }}
          >
            <h1 className="text-4xl sm:text-5xl lg:text-6xl xl:text-7xl font-bold text-white leading-tight">
              {t('title')}
              <br />
              <span className="text-primary">{t('titleHighlight')}</span>
            </h1>
            
            <p className="mt-6 text-lg sm:text-xl text-grey-eerie max-w-xl">
              {t('subtitle')}
            </p>

            <div className="mt-8 flex flex-col sm:flex-row gap-4">
              <AppStoreBadges />
            </div>

            {/* Stats */}
            <div className="mt-12 grid grid-cols-3 gap-8">
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.3 }}
              >
                <div className="text-3xl sm:text-4xl font-bold text-primary">
                  {formatNumber(STATS.users)}+
                </div>
                <div className="text-grey-jet text-sm mt-1">{t('stats.users')}</div>
              </motion.div>
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.4 }}
              >
                <div className="text-3xl sm:text-4xl font-bold text-primary">
                  {STATS.partners}+
                </div>
                <div className="text-grey-jet text-sm mt-1">{t('stats.partners')}</div>
              </motion.div>
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.5 }}
              >
                <div className="text-3xl sm:text-4xl font-bold text-primary">
                  {STATS.savings}%
                </div>
                <div className="text-grey-jet text-sm mt-1">{t('stats.savings')}</div>
              </motion.div>
            </div>
          </motion.div>

          {/* Right - Phone mockup */}
          <motion.div
            initial={{ opacity: 0, x: 50 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="hidden lg:flex justify-center"
          >
            <div className="phone-frame animate-bounce-subtle">
              <div className="phone-screen">
                <img
                  src="/images/app-screenshot-home.png"
                  alt="Yousoon App"
                  className="w-full h-full object-cover"
                />
              </div>
            </div>
          </motion.div>
        </div>

        {/* Scroll indicator */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1 }}
          className="absolute bottom-10 left-1/2 -translate-x-1/2 flex flex-col items-center gap-2 text-grey-jet"
        >
          <span className="text-sm">{t('secondaryCta')}</span>
          <ChevronDown className="animate-bounce" size={24} />
        </motion.div>
      </div>
    </section>
  )
}
