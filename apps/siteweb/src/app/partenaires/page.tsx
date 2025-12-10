'use client';

import { motion } from 'framer-motion';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import {
  Users,
  TrendingUp,
  Target,
  BarChart3,
  CheckCircle,
  ArrowRight,
  Building2,
  Utensils,
  Music,
  Dumbbell,
} from 'lucide-react';
import { Button } from '@/components/ui/Button';

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

const benefits = [
  {
    icon: Users,
    titleKey: 'newCustomers',
    descriptionKey: 'newCustomersDesc',
  },
  {
    icon: TrendingUp,
    titleKey: 'increaseRevenue',
    descriptionKey: 'increaseRevenueDesc',
  },
  {
    icon: Target,
    titleKey: 'targetedMarketing',
    descriptionKey: 'targetedMarketingDesc',
  },
  {
    icon: BarChart3,
    titleKey: 'analytics',
    descriptionKey: 'analyticsDesc',
  },
];

const partnerTypes = [
  {
    icon: Utensils,
    name: 'Restaurants',
    color: 'text-orange-400',
    bgColor: 'bg-orange-400/10',
  },
  {
    icon: Building2,
    name: 'Bars',
    color: 'text-blue-400',
    bgColor: 'bg-blue-400/10',
  },
  {
    icon: Music,
    name: 'Clubs & Events',
    color: 'text-purple-400',
    bgColor: 'bg-purple-400/10',
  },
  {
    icon: Dumbbell,
    name: 'Loisirs & Sport',
    color: 'text-green-400',
    bgColor: 'bg-green-400/10',
  },
];

const steps = [
  { number: '01', titleKey: 'step1', descriptionKey: 'step1Desc' },
  { number: '02', titleKey: 'step2', descriptionKey: 'step2Desc' },
  { number: '03', titleKey: 'step3', descriptionKey: 'step3Desc' },
  { number: '04', titleKey: 'step4', descriptionKey: 'step4Desc' },
];

export default function PartnersPage() {
  const t = useTranslations('partners');

  return (
    <>
      {/* Hero Section */}
      <section className="relative overflow-hidden py-24 md:py-32">
        <div className="container mx-auto px-4">
          <div className="grid items-center gap-12 lg:grid-cols-2">
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6 }}
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
              <div className="mt-8 flex flex-wrap gap-4">
                <Button asChild size="lg">
                  <Link href="https://business.yousoon.com/register">
                    {t('cta.primary')}
                    <ArrowRight className="ml-2 h-5 w-5" />
                  </Link>
                </Button>
                <Button asChild variant="outline" size="lg">
                  <Link href="/contact">{t('cta.secondary')}</Link>
                </Button>
              </div>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="relative"
            >
              {/* Stats cards */}
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="rounded-2xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
                  <div className="text-4xl font-bold text-primary">+45%</div>
                  <p className="mt-2 text-muted-foreground">
                    {t('stats.traffic')}
                  </p>
                </div>
                <div className="rounded-2xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
                  <div className="text-4xl font-bold text-primary">50K+</div>
                  <p className="mt-2 text-muted-foreground">
                    {t('stats.users')}
                  </p>
                </div>
                <div className="rounded-2xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
                  <div className="text-4xl font-bold text-primary">500+</div>
                  <p className="mt-2 text-muted-foreground">
                    {t('stats.partners')}
                  </p>
                </div>
                <div className="rounded-2xl border border-white/10 bg-white/5 p-6 backdrop-blur-sm">
                  <div className="text-4xl font-bold text-primary">4.8/5</div>
                  <p className="mt-2 text-muted-foreground">
                    {t('stats.satisfaction')}
                  </p>
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Partner Types */}
      <section className="border-y border-white/10 bg-white/5 py-16">
        <div className="container mx-auto px-4">
          <div className="text-center">
            <h2 className="text-2xl font-bold md:text-3xl">
              {t('partnerTypes.title')}
            </h2>
          </div>
          <div className="mt-12 flex flex-wrap items-center justify-center gap-8">
            {partnerTypes.map((type) => (
              <div key={type.name} className="flex items-center gap-3">
                <div className={`rounded-xl ${type.bgColor} p-3`}>
                  <type.icon className={`h-6 w-6 ${type.color}`} />
                </div>
                <span className="font-medium">{type.name}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Benefits */}
      <section className="py-16 md:py-24">
        <div className="container mx-auto px-4">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="mx-auto max-w-2xl text-center"
          >
            <h2 className="text-3xl font-bold md:text-4xl">
              {t('benefits.title')}
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              {t('benefits.subtitle')}
            </p>
          </motion.div>

          <motion.div
            variants={staggerContainer}
            initial="initial"
            whileInView="animate"
            viewport={{ once: true }}
            className="mt-16 grid gap-8 md:grid-cols-2 lg:grid-cols-4"
          >
            {benefits.map((benefit) => (
              <motion.div
                key={benefit.titleKey}
                variants={fadeInUp}
                className="text-center"
              >
                <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl bg-primary/10">
                  <benefit.icon className="h-8 w-8 text-primary" />
                </div>
                <h3 className="mt-6 text-xl font-semibold">
                  {t(`benefits.items.${benefit.titleKey}.title`)}
                </h3>
                <p className="mt-3 text-muted-foreground">
                  {t(`benefits.items.${benefit.titleKey}.description`)}
                </p>
              </motion.div>
            ))}
          </motion.div>
        </div>
      </section>

      {/* How it works */}
      <section className="border-y border-white/10 bg-white/5 py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold md:text-4xl">
              {t('howItWorks.title')}
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              {t('howItWorks.subtitle')}
            </p>
          </div>

          <div className="mt-16 grid gap-8 md:grid-cols-2 lg:grid-cols-4">
            {steps.map((step, index) => (
              <motion.div
                key={step.number}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
                className="relative"
              >
                <div className="text-6xl font-bold text-primary/20">
                  {step.number}
                </div>
                <h3 className="mt-4 text-xl font-semibold">
                  {t(`howItWorks.steps.${step.titleKey}.title`)}
                </h3>
                <p className="mt-2 text-muted-foreground">
                  {t(`howItWorks.steps.${step.titleKey}.description`)}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Features checklist */}
      <section className="py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="grid items-center gap-12 lg:grid-cols-2">
            <div>
              <h2 className="text-3xl font-bold md:text-4xl">
                {t('features.title')}
              </h2>
              <p className="mt-4 text-lg text-muted-foreground">
                {t('features.subtitle')}
              </p>
              <ul className="mt-8 space-y-4">
                {[
                  'feature1',
                  'feature2',
                  'feature3',
                  'feature4',
                  'feature5',
                  'feature6',
                ].map((feature) => (
                  <li key={feature} className="flex items-center gap-3">
                    <CheckCircle className="h-5 w-5 flex-shrink-0 text-primary" />
                    <span>{t(`features.items.${feature}`)}</span>
                  </li>
                ))}
              </ul>
            </div>

            <div className="rounded-2xl border border-white/10 bg-white/5 p-8">
              <h3 className="text-2xl font-bold">{t('pricing.title')}</h3>
              <div className="mt-6">
                <span className="text-5xl font-bold">{t('pricing.price')}</span>
                <span className="text-muted-foreground">
                  {t('pricing.period')}
                </span>
              </div>
              <p className="mt-4 text-muted-foreground">
                {t('pricing.description')}
              </p>
              <Button asChild className="mt-8 w-full" size="lg">
                <Link href="https://business.yousoon.com/register">
                  {t('pricing.cta')}
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="border-t border-white/10 bg-gradient-to-b from-primary/10 to-background py-16 md:py-24">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl font-bold md:text-4xl">{t('cta.title')}</h2>
          <p className="mx-auto mt-4 max-w-2xl text-lg text-muted-foreground">
            {t('cta.description')}
          </p>
          <div className="mt-8 flex flex-wrap justify-center gap-4">
            <Button asChild size="lg">
              <Link href="https://business.yousoon.com/register">
                {t('cta.button')}
                <ArrowRight className="ml-2 h-5 w-5" />
              </Link>
            </Button>
          </div>
        </div>
      </section>
    </>
  );
}
