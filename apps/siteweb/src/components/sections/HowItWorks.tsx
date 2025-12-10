'use client'

import { useTranslations } from 'next-intl'
import { motion } from 'framer-motion'
import { Download, Search, CalendarCheck, PartyPopper } from 'lucide-react'

const icons = {
  Download,
  Search,
  CalendarCheck,
  PartyPopper,
}

const STEPS_DATA = [
  { key: 'download', icon: 'Download' },
  { key: 'explore', icon: 'Search' },
  { key: 'book', icon: 'CalendarCheck' },
  { key: 'enjoy', icon: 'PartyPopper' },
]

export function HowItWorks() {
  const t = useTranslations('howItWorks')

  return (
    <section className="section-padding bg-black relative overflow-hidden">
      {/* Background */}
      <div className="absolute inset-0">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-primary/5 rounded-full blur-3xl" />
      </div>

      <div className="container-wide relative z-10">
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

        {/* Steps */}
        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-8 lg:gap-4">
          {STEPS_DATA.map((step, index) => {
            const Icon = icons[step.icon as keyof typeof icons]
            return (
              <motion.div
                key={step.key}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.15 }}
                className="relative"
              >
                {/* Connector line */}
                {index < STEPS_DATA.length - 1 && (
                  <div className="hidden lg:block absolute top-12 left-1/2 w-full h-0.5 bg-gradient-to-r from-primary/50 to-transparent" />
                )}

                <div className="flex flex-col items-center text-center">
                  {/* Step number */}
                  <div className="relative">
                    <div className="w-24 h-24 rounded-full bg-dark-900 border-2 border-primary flex items-center justify-center mb-6 glow-primary">
                      <Icon className="w-10 h-10 text-primary" />
                    </div>
                    <div className="absolute -top-2 -right-2 w-8 h-8 rounded-full bg-primary text-black font-bold flex items-center justify-center text-sm">
                      {index + 1}
                    </div>
                  </div>

                  <h3 className="text-xl font-semibold text-white mb-2">
                    {t(`steps.${step.key}.title`)}
                  </h3>
                  <p className="text-grey-eerie text-sm max-w-[200px]">
                    {t(`steps.${step.key}.description`)}
                  </p>
                </div>
              </motion.div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
