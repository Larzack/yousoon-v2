'use client';

import { motion } from 'framer-motion';
import { useTranslations } from 'next-intl';
import Image from 'next/image';
import Link from 'next/link';
import { Heart, Users, Globe, Sparkles, ArrowRight } from 'lucide-react';
import { Button } from '@/components/ui/Button';

const fadeInUp = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  transition: { duration: 0.5 },
};

const values = [
  {
    icon: Heart,
    titleKey: 'passion',
    descriptionKey: 'passionDesc',
    color: 'text-red-400',
    bgColor: 'bg-red-400/10',
  },
  {
    icon: Users,
    titleKey: 'community',
    descriptionKey: 'communityDesc',
    color: 'text-blue-400',
    bgColor: 'bg-blue-400/10',
  },
  {
    icon: Globe,
    titleKey: 'accessibility',
    descriptionKey: 'accessibilityDesc',
    color: 'text-green-400',
    bgColor: 'bg-green-400/10',
  },
  {
    icon: Sparkles,
    titleKey: 'innovation',
    descriptionKey: 'innovationDesc',
    color: 'text-primary',
    bgColor: 'bg-primary/10',
  },
];

const milestones = [
  { year: '2023', titleKey: 'milestone1', descriptionKey: 'milestone1Desc' },
  { year: '2024', titleKey: 'milestone2', descriptionKey: 'milestone2Desc' },
  { year: '2024', titleKey: 'milestone3', descriptionKey: 'milestone3Desc' },
  { year: '2025', titleKey: 'milestone4', descriptionKey: 'milestone4Desc' },
];

export default function AboutPage() {
  const t = useTranslations('about');

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
            <span className="inline-block rounded-full bg-primary/10 px-4 py-2 text-sm font-medium text-primary">
              {t('badge')}
            </span>
            <h1 className="mt-6 text-4xl font-bold tracking-tight md:text-5xl lg:text-6xl">
              {t('title')}
            </h1>
            <p className="mt-6 text-lg text-muted-foreground md:text-xl">
              {t('subtitle')}
            </p>
          </motion.div>
        </div>
      </section>

      {/* Mission Section */}
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
                {t('mission.title')}
              </h2>
              <p className="mt-6 text-lg text-muted-foreground">
                {t('mission.description1')}
              </p>
              <p className="mt-4 text-lg text-muted-foreground">
                {t('mission.description2')}
              </p>
              <div className="mt-8 grid grid-cols-2 gap-6">
                <div>
                  <div className="text-4xl font-bold text-primary">50K+</div>
                  <p className="mt-1 text-muted-foreground">
                    {t('mission.stats.users')}
                  </p>
                </div>
                <div>
                  <div className="text-4xl font-bold text-primary">500+</div>
                  <p className="mt-1 text-muted-foreground">
                    {t('mission.stats.partners')}
                  </p>
                </div>
                <div>
                  <div className="text-4xl font-bold text-primary">100K+</div>
                  <p className="mt-1 text-muted-foreground">
                    {t('mission.stats.bookings')}
                  </p>
                </div>
                <div>
                  <div className="text-4xl font-bold text-primary">4.8/5</div>
                  <p className="mt-1 text-muted-foreground">
                    {t('mission.stats.rating')}
                  </p>
                </div>
              </div>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, x: 20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="relative"
            >
              <div className="aspect-square overflow-hidden rounded-2xl bg-gradient-to-br from-primary/20 to-primary/5">
                <div className="flex h-full items-center justify-center">
                  <span className="text-6xl">🎉</span>
                </div>
              </div>
              {/* Decorative elements */}
              <div className="absolute -right-4 -top-4 h-24 w-24 rounded-full bg-primary/20 blur-2xl" />
              <div className="absolute -bottom-4 -left-4 h-32 w-32 rounded-full bg-primary/10 blur-3xl" />
            </motion.div>
          </div>
        </div>
      </section>

      {/* Values Section */}
      <section className="border-y border-white/10 bg-white/5 py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold md:text-4xl">
              {t('values.title')}
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              {t('values.subtitle')}
            </p>
          </div>

          <div className="mt-16 grid gap-8 md:grid-cols-2 lg:grid-cols-4">
            {values.map((value, index) => (
              <motion.div
                key={value.titleKey}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
                className="text-center"
              >
                <div
                  className={`mx-auto flex h-16 w-16 items-center justify-center rounded-2xl ${value.bgColor}`}
                >
                  <value.icon className={`h-8 w-8 ${value.color}`} />
                </div>
                <h3 className="mt-6 text-xl font-semibold">
                  {t(`values.items.${value.titleKey}.title`)}
                </h3>
                <p className="mt-3 text-muted-foreground">
                  {t(`values.items.${value.titleKey}.description`)}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Timeline Section */}
      <section className="py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold md:text-4xl">
              {t('timeline.title')}
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              {t('timeline.subtitle')}
            </p>
          </div>

          <div className="relative mx-auto mt-16 max-w-3xl">
            {/* Timeline line */}
            <div className="absolute left-0 top-0 h-full w-px bg-gradient-to-b from-primary via-primary/50 to-transparent md:left-1/2 md:-translate-x-1/2" />

            {milestones.map((milestone, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, x: index % 2 === 0 ? -20 : 20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
                className={`relative mb-12 pl-8 md:w-1/2 md:pl-0 ${
                  index % 2 === 0
                    ? 'md:pr-12 md:text-right'
                    : 'md:ml-auto md:pl-12'
                }`}
              >
                {/* Dot */}
                <div
                  className={`absolute left-0 top-0 h-3 w-3 rounded-full bg-primary md:left-auto ${
                    index % 2 === 0
                      ? 'md:-right-1.5 md:left-auto'
                      : 'md:-left-1.5'
                  }`}
                />
                <span className="text-sm font-medium text-primary">
                  {milestone.year}
                </span>
                <h3 className="mt-2 text-xl font-semibold">
                  {t(`timeline.items.${milestone.titleKey}.title`)}
                </h3>
                <p className="mt-2 text-muted-foreground">
                  {t(`timeline.items.${milestone.titleKey}.description`)}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Team Section (placeholder) */}
      <section className="border-y border-white/10 bg-white/5 py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold md:text-4xl">{t('team.title')}</h2>
            <p className="mt-4 text-lg text-muted-foreground">
              {t('team.subtitle')}
            </p>
          </div>

          <div className="mt-16 grid gap-8 md:grid-cols-3">
            {[1, 2, 3].map((_, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
                className="text-center"
              >
                <div className="mx-auto h-32 w-32 rounded-full bg-gradient-to-br from-primary/30 to-primary/10" />
                <h3 className="mt-6 text-lg font-semibold">
                  {t(`team.members.member${index + 1}.name`)}
                </h3>
                <p className="text-sm text-primary">
                  {t(`team.members.member${index + 1}.role`)}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-16 md:py-24">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl font-bold md:text-4xl">{t('cta.title')}</h2>
          <p className="mx-auto mt-4 max-w-2xl text-muted-foreground">
            {t('cta.description')}
          </p>
          <div className="mt-8 flex flex-wrap justify-center gap-4">
            <Button asChild size="lg">
              <Link href="/telecharger">
                {t('cta.download')}
                <ArrowRight className="ml-2 h-5 w-5" />
              </Link>
            </Button>
            <Button asChild variant="outline" size="lg">
              <Link href="/contact">{t('cta.contact')}</Link>
            </Button>
          </div>
        </div>
      </section>
    </>
  );
}
