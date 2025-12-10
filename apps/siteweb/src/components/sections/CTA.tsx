'use client'

import { useTranslations } from 'next-intl'
import { motion } from 'framer-motion'
import { Button } from '@/components/ui'
import { AppStoreBadges } from '@/components/shared/AppStoreBadges'

export function CTA() {
  const t = useTranslations('cta')

  return (
    <section className="section-padding bg-black relative overflow-hidden">
      {/* Background gradient */}
      <div className="absolute inset-0">
        <div className="absolute inset-0 bg-gradient-to-b from-dark-950 via-black to-dark-950" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[400px] bg-primary/10 rounded-full blur-3xl" />
      </div>

      <div className="container-narrow relative z-10">
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          className="text-center"
        >
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-white">
            {t('title')}
          </h2>
          <p className="mt-4 text-lg text-grey-eerie max-w-xl mx-auto">
            {t('subtitle')}
          </p>

          <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
            <AppStoreBadges size="lg" />
          </div>

          <p className="mt-6 text-sm text-grey-jet">
            {t('note')}
          </p>
        </motion.div>
      </div>
    </section>
  )
}
