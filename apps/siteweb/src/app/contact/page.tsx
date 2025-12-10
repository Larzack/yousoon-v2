'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';
import { useTranslations } from 'next-intl';
import {
  Mail,
  Phone,
  MapPin,
  Send,
  MessageSquare,
  Building2,
  Users,
} from 'lucide-react';
import { Button } from '@/components/ui/Button';

const contactMethods = [
  {
    icon: Mail,
    titleKey: 'email',
    value: 'contact@yousoon.com',
    href: 'mailto:contact@yousoon.com',
  },
  {
    icon: Phone,
    titleKey: 'phone',
    value: '+33 1 23 45 67 89',
    href: 'tel:+33123456789',
  },
  {
    icon: MapPin,
    titleKey: 'address',
    value: 'Paris, France',
    href: null,
  },
];

const subjects = [
  { value: 'general', icon: MessageSquare },
  { value: 'partnership', icon: Building2 },
  { value: 'support', icon: Users },
  { value: 'press', icon: Mail },
];

export default function ContactPage() {
  const t = useTranslations('contact');
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    subject: 'general',
    message: '',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    // Simulate form submission
    await new Promise((resolve) => setTimeout(resolve, 1500));

    setIsSubmitting(false);
    setIsSubmitted(true);
  };

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

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
          </motion.div>
        </div>
      </section>

      {/* Contact Methods */}
      <section className="border-y border-white/10 bg-white/5 py-12">
        <div className="container mx-auto px-4">
          <div className="grid gap-8 md:grid-cols-3">
            {contactMethods.map((method, index) => (
              <motion.div
                key={method.titleKey}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                className="flex items-center gap-4"
              >
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10">
                  <method.icon className="h-6 w-6 text-primary" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">
                    {t(`methods.${method.titleKey}`)}
                  </p>
                  {method.href ? (
                    <a
                      href={method.href}
                      className="font-medium hover:text-primary"
                    >
                      {method.value}
                    </a>
                  ) : (
                    <p className="font-medium">{method.value}</p>
                  )}
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Contact Form */}
      <section className="py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-2xl">
            {isSubmitted ? (
              <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="rounded-2xl border border-primary/50 bg-primary/10 p-8 text-center"
              >
                <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-primary/20">
                  <Send className="h-8 w-8 text-primary" />
                </div>
                <h2 className="mt-6 text-2xl font-bold">
                  {t('form.success.title')}
                </h2>
                <p className="mt-4 text-muted-foreground">
                  {t('form.success.description')}
                </p>
                <Button
                  onClick={() => {
                    setIsSubmitted(false);
                    setFormData({
                      name: '',
                      email: '',
                      subject: 'general',
                      message: '',
                    });
                  }}
                  variant="outline"
                  className="mt-6"
                >
                  {t('form.success.button')}
                </Button>
              </motion.div>
            ) : (
              <motion.form
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
                onSubmit={handleSubmit}
                className="space-y-6"
              >
                {/* Subject Selection */}
                <div>
                  <label className="mb-3 block text-sm font-medium">
                    {t('form.subject.label')}
                  </label>
                  <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
                    {subjects.map((subject) => (
                      <button
                        key={subject.value}
                        type="button"
                        onClick={() =>
                          setFormData((prev) => ({
                            ...prev,
                            subject: subject.value,
                          }))
                        }
                        className={`flex flex-col items-center gap-2 rounded-xl border p-4 transition-all ${
                          formData.subject === subject.value
                            ? 'border-primary bg-primary/10 text-primary'
                            : 'border-white/10 bg-white/5 hover:border-white/20'
                        }`}
                      >
                        <subject.icon className="h-6 w-6" />
                        <span className="text-sm">
                          {t(`form.subject.options.${subject.value}`)}
                        </span>
                      </button>
                    ))}
                  </div>
                </div>

                {/* Name & Email */}
                <div className="grid gap-6 md:grid-cols-2">
                  <div>
                    <label
                      htmlFor="name"
                      className="mb-2 block text-sm font-medium"
                    >
                      {t('form.name.label')}
                    </label>
                    <input
                      type="text"
                      id="name"
                      name="name"
                      value={formData.name}
                      onChange={handleChange}
                      required
                      placeholder={t('form.name.placeholder')}
                      className="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 outline-none transition-colors placeholder:text-muted-foreground focus:border-primary focus:ring-1 focus:ring-primary"
                    />
                  </div>
                  <div>
                    <label
                      htmlFor="email"
                      className="mb-2 block text-sm font-medium"
                    >
                      {t('form.email.label')}
                    </label>
                    <input
                      type="email"
                      id="email"
                      name="email"
                      value={formData.email}
                      onChange={handleChange}
                      required
                      placeholder={t('form.email.placeholder')}
                      className="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 outline-none transition-colors placeholder:text-muted-foreground focus:border-primary focus:ring-1 focus:ring-primary"
                    />
                  </div>
                </div>

                {/* Message */}
                <div>
                  <label
                    htmlFor="message"
                    className="mb-2 block text-sm font-medium"
                  >
                    {t('form.message.label')}
                  </label>
                  <textarea
                    id="message"
                    name="message"
                    value={formData.message}
                    onChange={handleChange}
                    required
                    rows={6}
                    placeholder={t('form.message.placeholder')}
                    className="w-full resize-none rounded-xl border border-white/10 bg-white/5 px-4 py-3 outline-none transition-colors placeholder:text-muted-foreground focus:border-primary focus:ring-1 focus:ring-primary"
                  />
                </div>

                {/* Submit */}
                <Button
                  type="submit"
                  size="lg"
                  className="w-full"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? (
                    <>
                      <svg
                        className="mr-2 h-5 w-5 animate-spin"
                        viewBox="0 0 24 24"
                      >
                        <circle
                          className="opacity-25"
                          cx="12"
                          cy="12"
                          r="10"
                          stroke="currentColor"
                          strokeWidth="4"
                          fill="none"
                        />
                        <path
                          className="opacity-75"
                          fill="currentColor"
                          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        />
                      </svg>
                      {t('form.submitting')}
                    </>
                  ) : (
                    <>
                      {t('form.submit')}
                      <Send className="ml-2 h-5 w-5" />
                    </>
                  )}
                </Button>
              </motion.form>
            )}
          </div>
        </div>
      </section>

      {/* FAQ CTA */}
      <section className="border-t border-white/10 py-16">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-2xl font-bold">{t('faq.title')}</h2>
          <p className="mt-2 text-muted-foreground">{t('faq.description')}</p>
          <Button asChild variant="outline" className="mt-6">
            <a href="/#faq">{t('faq.button')}</a>
          </Button>
        </div>
      </section>
    </>
  );
}
