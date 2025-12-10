'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { Check, X, Sparkles } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { cn } from '@/lib/utils';

const fadeInUp = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  transition: { duration: 0.5 },
};

interface PricingPlan {
  id: string;
  name: string;
  price: number;
  priceYearly: number;
  description: string;
  features: { name: string; included: boolean }[];
  popular?: boolean;
  cta: string;
}

const plans: PricingPlan[] = [
  {
    id: 'free',
    name: 'Découverte',
    price: 0,
    priceYearly: 0,
    description: 'Parfait pour découvrir Yousoon',
    features: [
      { name: 'Accès à toutes les offres', included: true },
      { name: '3 réservations par mois', included: true },
      { name: 'Favoris illimités', included: true },
      { name: 'Notifications basiques', included: true },
      { name: 'Offres exclusives', included: false },
      { name: 'Réservations illimitées', included: false },
      { name: 'Accès prioritaire', included: false },
      { name: 'Support prioritaire', included: false },
    ],
    cta: 'Commencer gratuitement',
  },
  {
    id: 'monthly',
    name: 'Premium',
    price: 9.99,
    priceYearly: 7.99,
    description: 'Pour profiter pleinement de Yousoon',
    features: [
      { name: 'Accès à toutes les offres', included: true },
      { name: 'Réservations illimitées', included: true },
      { name: 'Favoris illimités', included: true },
      { name: 'Notifications personnalisées', included: true },
      { name: 'Offres exclusives Premium', included: true },
      { name: 'Accès prioritaire aux nouveautés', included: true },
      { name: 'Support prioritaire', included: true },
      { name: 'Sans engagement', included: true },
    ],
    popular: true,
    cta: 'Essai gratuit 30 jours',
  },
  {
    id: 'yearly',
    name: 'Premium Annuel',
    price: 79.99,
    priceYearly: 79.99,
    description: '2 mois offerts',
    features: [
      { name: 'Tous les avantages Premium', included: true },
      { name: 'Réservations illimitées', included: true },
      { name: '2 mois gratuits', included: true },
      { name: 'Badge membre fidèle', included: true },
      { name: 'Offres VIP partenaires', included: true },
      { name: 'Avant-premières', included: true },
      { name: 'Support dédié', included: true },
      { name: 'Meilleur rapport qualité/prix', included: true },
    ],
    cta: 'Économiser 40€/an',
  },
];

const faqs = [
  {
    question: 'Comment fonctionne l\'essai gratuit ?',
    answer:
      'L\'essai gratuit Premium dure 30 jours. Vous avez accès à toutes les fonctionnalités Premium sans engagement. Vous pouvez annuler à tout moment avant la fin de la période d\'essai.',
  },
  {
    question: 'Puis-je changer de plan à tout moment ?',
    answer:
      'Oui, vous pouvez passer d\'un plan à un autre à tout moment. Si vous passez au plan annuel, vous bénéficierez immédiatement des avantages supplémentaires.',
  },
  {
    question: 'Comment fonctionne la facturation ?',
    answer:
      'La facturation se fait directement via l\'App Store (iOS) ou Google Play (Android). Tous les paiements sont sécurisés et gérés par ces plateformes.',
  },
  {
    question: 'Puis-je annuler mon abonnement ?',
    answer:
      'Oui, vous pouvez annuler votre abonnement à tout moment depuis les paramètres de votre compte. Vous conserverez l\'accès Premium jusqu\'à la fin de la période payée.',
  },
];

export default function PricingPage() {
  const t = useTranslations('pricing');
  const [billingPeriod, setBillingPeriod] = useState<'monthly' | 'yearly'>(
    'monthly'
  );

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
              {t('title')}
            </h1>
            <p className="mt-6 text-lg text-muted-foreground md:text-xl">
              {t('subtitle')}
            </p>

            {/* Billing toggle */}
            <div className="mt-10 flex items-center justify-center gap-4">
              <span
                className={cn(
                  'text-sm font-medium',
                  billingPeriod === 'monthly'
                    ? 'text-foreground'
                    : 'text-muted-foreground'
                )}
              >
                {t('monthly')}
              </span>
              <button
                onClick={() =>
                  setBillingPeriod(
                    billingPeriod === 'monthly' ? 'yearly' : 'monthly'
                  )
                }
                className={cn(
                  'relative h-8 w-14 rounded-full transition-colors',
                  billingPeriod === 'yearly'
                    ? 'bg-primary'
                    : 'bg-white/20'
                )}
              >
                <span
                  className={cn(
                    'absolute top-1 h-6 w-6 rounded-full bg-white transition-all',
                    billingPeriod === 'yearly' ? 'left-7' : 'left-1'
                  )}
                />
              </button>
              <span
                className={cn(
                  'text-sm font-medium',
                  billingPeriod === 'yearly'
                    ? 'text-foreground'
                    : 'text-muted-foreground'
                )}
              >
                {t('yearly')}
                <span className="ml-2 rounded-full bg-primary/10 px-2 py-1 text-xs text-primary">
                  -20%
                </span>
              </span>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Pricing Cards */}
      <section className="pb-16 md:pb-24">
        <div className="container mx-auto px-4">
          <div className="grid gap-8 lg:grid-cols-3">
            {plans.map((plan, index) => (
              <motion.div
                key={plan.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                className={cn(
                  'relative rounded-2xl border p-8',
                  plan.popular
                    ? 'border-primary bg-primary/5'
                    : 'border-white/10 bg-white/5'
                )}
              >
                {plan.popular && (
                  <div className="absolute -top-4 left-1/2 -translate-x-1/2">
                    <span className="inline-flex items-center gap-1 rounded-full bg-primary px-4 py-1 text-sm font-medium text-black">
                      <Sparkles className="h-4 w-4" />
                      {t('popular')}
                    </span>
                  </div>
                )}

                <div className="text-center">
                  <h3 className="text-xl font-semibold">{plan.name}</h3>
                  <p className="mt-2 text-sm text-muted-foreground">
                    {plan.description}
                  </p>
                  <div className="mt-6">
                    <span className="text-5xl font-bold">
                      {billingPeriod === 'yearly'
                        ? plan.priceYearly
                        : plan.price}
                      €
                    </span>
                    {plan.price > 0 && (
                      <span className="text-muted-foreground">/mois</span>
                    )}
                  </div>
                </div>

                <ul className="mt-8 space-y-4">
                  {plan.features.map((feature) => (
                    <li key={feature.name} className="flex items-center gap-3">
                      {feature.included ? (
                        <Check className="h-5 w-5 flex-shrink-0 text-primary" />
                      ) : (
                        <X className="h-5 w-5 flex-shrink-0 text-muted-foreground/50" />
                      )}
                      <span
                        className={cn(
                          feature.included
                            ? 'text-foreground'
                            : 'text-muted-foreground/50'
                        )}
                      >
                        {feature.name}
                      </span>
                    </li>
                  ))}
                </ul>

                <Button
                  asChild
                  className="mt-8 w-full"
                  variant={plan.popular ? 'default' : 'outline'}
                  size="lg"
                >
                  <Link href="/telecharger">{plan.cta}</Link>
                </Button>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* FAQ Section */}
      <section className="border-t border-white/10 py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold md:text-4xl">
              {t('faq.title')}
            </h2>
            <p className="mt-4 text-muted-foreground">{t('faq.subtitle')}</p>
          </div>

          <div className="mx-auto mt-12 max-w-3xl space-y-4">
            {faqs.map((faq, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 10 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
                className="rounded-xl border border-white/10 bg-white/5 p-6"
              >
                <h3 className="font-semibold">{faq.question}</h3>
                <p className="mt-2 text-muted-foreground">{faq.answer}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="border-t border-white/10 bg-gradient-to-b from-primary/10 to-background py-16 md:py-24">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl font-bold md:text-4xl">{t('cta.title')}</h2>
          <p className="mx-auto mt-4 max-w-2xl text-muted-foreground">
            {t('cta.description')}
          </p>
          <Button asChild size="lg" className="mt-8">
            <Link href="/telecharger">{t('cta.button')}</Link>
          </Button>
        </div>
      </section>
    </>
  );
}
