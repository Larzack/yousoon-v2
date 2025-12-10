'use client'

import { useTranslations } from 'next-intl'
import { motion } from 'framer-motion'
import { Star, Quote } from 'lucide-react'
import { Card } from '@/components/ui'
import { TESTIMONIALS } from '@/lib/constants'

export function Testimonials() {
  const t = useTranslations('testimonials')

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

        {/* Testimonials grid */}
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8">
          {TESTIMONIALS.map((testimonial, index) => (
            <motion.div
              key={testimonial.id}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
            >
              <Card className="p-6 h-full border-dark-800 relative overflow-hidden">
                {/* Quote icon */}
                <Quote className="absolute top-4 right-4 w-8 h-8 text-primary/20" />

                {/* Stars */}
                <div className="flex gap-1 mb-4">
                  {Array.from({ length: 5 }).map((_, i) => (
                    <Star
                      key={i}
                      className={`w-5 h-5 ${
                        i < testimonial.rating
                          ? 'text-primary fill-primary'
                          : 'text-grey-jet'
                      }`}
                    />
                  ))}
                </div>

                {/* Content */}
                <p className="text-grey-eerie mb-6 italic">
                  "{testimonial.content}"
                </p>

                {/* Author */}
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 rounded-full bg-dark-800 flex items-center justify-center text-primary font-semibold">
                    {testimonial.name.charAt(0)}
                  </div>
                  <div>
                    <div className="font-semibold text-white">{testimonial.name}</div>
                    <div className="text-sm text-grey-jet">{testimonial.role}</div>
                  </div>
                </div>
              </Card>
            </motion.div>
          ))}
        </div>

        {/* App Store ratings */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mt-12 flex flex-col sm:flex-row items-center justify-center gap-8"
        >
          <div className="flex items-center gap-3 bg-dark-900 rounded-xl px-6 py-4">
            <img src="/images/app-store-icon.svg" alt="App Store" className="w-10 h-10" />
            <div>
              <div className="flex items-center gap-1">
                <span className="text-2xl font-bold text-white">4.8</span>
                <Star className="w-5 h-5 text-primary fill-primary" />
              </div>
              <div className="text-sm text-grey-jet">App Store</div>
            </div>
          </div>
          <div className="flex items-center gap-3 bg-dark-900 rounded-xl px-6 py-4">
            <img src="/images/play-store-icon.svg" alt="Play Store" className="w-10 h-10" />
            <div>
              <div className="flex items-center gap-1">
                <span className="text-2xl font-bold text-white">4.7</span>
                <Star className="w-5 h-5 text-primary fill-primary" />
              </div>
              <div className="text-sm text-grey-jet">Play Store</div>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  )
}
