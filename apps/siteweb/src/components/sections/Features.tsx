'use client'

import { useTranslations } from 'next-intl'
import { motion } from 'framer-motion'
import { MapPin, Percent, Zap, Heart, Users, Shield } from 'lucide-react'
import { Card } from '@/components/ui'

const icons = {
  MapPin,
  Percent,
  Zap,
  Heart,
  Users,
  Shield,
}

const FEATURES_DATA = [
  { key: 'geolocated', icon: 'MapPin' },
  { key: 'discounts', icon: 'Percent' },
  { key: 'instant', icon: 'Zap' },
  { key: 'favorites', icon: 'Heart' },
  { key: 'community', icon: 'Users' },
  { key: 'secure', icon: 'Shield' },
]

export function Features() {
  const t = useTranslations('features')

  return (
    <section className="section-padding bg-dark-950">
      <div className="container-wide">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center max-w-2xl mx-auto mb-16"
        >
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-white">
            {t('title')}
          </h2>
          <p className="mt-4 text-lg text-grey-eerie">
            {t('subtitle')}
          </p>
        </motion.div>

        {/* Features grid */}
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
          {FEATURES_DATA.map((feature, index) => {
            const Icon = icons[feature.icon as keyof typeof icons]
            return (
              <motion.div
                key={feature.key}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
              >
                <Card className="p-6 h-full card-hover border-dark-800 hover:border-primary/50">
                  <div className="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center mb-4">
                    <Icon className="w-6 h-6 text-primary" />
                  </div>
                  <h3 className="text-xl font-semibold text-white mb-2">
                    {t(`items.${feature.key}.title`)}
                  </h3>
                  <p className="text-grey-eerie">
                    {t(`items.${feature.key}.description`)}
                  </p>
                </Card>
              </motion.div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
