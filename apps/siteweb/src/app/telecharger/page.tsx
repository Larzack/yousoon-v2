'use client';

import { useTranslations } from 'next-intl';
import { motion } from 'framer-motion';
import Image from 'next/image';
import { Header } from '@/components/layout/Header';
import { Footer } from '@/components/layout/Footer';
import { AppStoreBadges } from '@/components/shared/AppStoreBadges';
import { 
  Smartphone, 
  Shield, 
  Zap, 
  Download,
  CheckCircle,
  Star,
  MapPin,
  Heart,
  QrCode,
  Bell
} from 'lucide-react';

const fadeInUp = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  transition: { duration: 0.5 }
};

const staggerContainer = {
  animate: {
    transition: {
      staggerChildren: 0.1
    }
  }
};

export default function TelechargerPage() {
  const t = useTranslations('download');

  const features = [
    {
      icon: MapPin,
      title: t('features.discover.title'),
      description: t('features.discover.description')
    },
    {
      icon: Heart,
      title: t('features.favorites.title'),
      description: t('features.favorites.description')
    },
    {
      icon: QrCode,
      title: t('features.qrcode.title'),
      description: t('features.qrcode.description')
    },
    {
      icon: Bell,
      title: t('features.notifications.title'),
      description: t('features.notifications.description')
    }
  ];

  const requirements = {
    ios: {
      version: 'iOS 17.0',
      size: '~50 MB',
      languages: 'Français, English'
    },
    android: {
      version: 'Android 14 (API 34)',
      size: '~45 MB',
      languages: 'Français, English'
    }
  };

  const stats = [
    { value: '4.8', label: t('stats.rating'), icon: Star },
    { value: '10K+', label: t('stats.downloads'), icon: Download },
    { value: '500+', label: t('stats.partners'), icon: CheckCircle }
  ];

  return (
    <>
      <Header />
      <main className="min-h-screen bg-black">
        {/* Hero Section */}
        <section className="relative pt-32 pb-20 overflow-hidden">
          {/* Gradient Background */}
          <div className="absolute inset-0 bg-gradient-to-b from-primary/10 via-transparent to-transparent" />
          
          <div className="container mx-auto px-6 relative">
            <div className="grid lg:grid-cols-2 gap-12 items-center">
              {/* Content */}
              <motion.div
                initial={{ opacity: 0, x: -30 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.6 }}
              >
                <div className="inline-flex items-center gap-2 px-4 py-2 bg-primary/10 rounded-full mb-6">
                  <Smartphone className="w-4 h-4 text-primary" />
                  <span className="text-sm text-primary font-medium">
                    {t('badge')}
                  </span>
                </div>

                <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-white mb-6">
                  {t('title.line1')}
                  <span className="text-primary block">{t('title.line2')}</span>
                </h1>

                <p className="text-lg text-gray-400 mb-8 max-w-lg">
                  {t('description')}
                </p>

                {/* Stats */}
                <div className="flex gap-8 mb-8">
                  {stats.map((stat, index) => (
                    <motion.div
                      key={index}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: 0.2 + index * 0.1 }}
                      className="text-center"
                    >
                      <div className="flex items-center justify-center gap-1 mb-1">
                        <stat.icon className="w-4 h-4 text-primary" />
                        <span className="text-2xl font-bold text-white">{stat.value}</span>
                      </div>
                      <span className="text-sm text-gray-500">{stat.label}</span>
                    </motion.div>
                  ))}
                </div>

                {/* Download Buttons */}
                <AppStoreBadges size="lg" className="mb-8" />

                {/* QR Code */}
                <div className="flex items-center gap-4 p-4 bg-white/5 rounded-xl border border-white/10 max-w-fit">
                  <div className="w-20 h-20 bg-white rounded-lg flex items-center justify-center">
                    <QrCode className="w-16 h-16 text-black" />
                  </div>
                  <div>
                    <p className="text-white font-medium">{t('qrCode.title')}</p>
                    <p className="text-sm text-gray-400">{t('qrCode.description')}</p>
                  </div>
                </div>
              </motion.div>

              {/* Phone Mockup */}
              <motion.div
                initial={{ opacity: 0, x: 30 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.6, delay: 0.2 }}
                className="relative flex justify-center"
              >
                <div className="relative">
                  {/* Glow Effect */}
                  <div className="absolute inset-0 bg-primary/20 blur-3xl rounded-full scale-150" />
                  
                  {/* Phone Frame */}
                  <div className="relative w-72 h-[580px] bg-gradient-to-b from-gray-800 to-gray-900 rounded-[3rem] p-2 shadow-2xl border border-gray-700">
                    {/* Screen */}
                    <div className="w-full h-full bg-black rounded-[2.5rem] overflow-hidden relative">
                      {/* Notch */}
                      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-32 h-7 bg-black rounded-b-2xl z-10" />
                      
                      {/* App Screenshot Placeholder */}
                      <div className="w-full h-full bg-gradient-to-b from-gray-900 to-black flex flex-col">
                        {/* Status Bar */}
                        <div className="pt-8 px-6 flex justify-between items-center text-white text-xs">
                          <span>9:41</span>
                          <div className="flex gap-1">
                            <div className="w-4 h-2 bg-white rounded-sm" />
                            <div className="w-4 h-2 bg-white rounded-sm" />
                            <div className="w-6 h-3 bg-primary rounded-sm" />
                          </div>
                        </div>
                        
                        {/* App Content */}
                        <div className="flex-1 p-4">
                          <div className="text-white font-bold text-lg mb-4">Pour vous</div>
                          
                          {/* Card Placeholders */}
                          <div className="space-y-3">
                            <div className="h-32 bg-gradient-to-br from-primary/30 to-primary/10 rounded-xl border border-primary/20" />
                            <div className="h-24 bg-white/5 rounded-xl border border-white/10" />
                            <div className="h-24 bg-white/5 rounded-xl border border-white/10" />
                          </div>
                        </div>
                        
                        {/* Tab Bar */}
                        <div className="h-16 bg-gray-900/80 backdrop-blur border-t border-white/10 flex justify-around items-center px-4">
                          {['📅', '❤️', '🃏', '📍', '💬'].map((emoji, i) => (
                            <div
                              key={i}
                              className={`w-10 h-10 rounded-full flex items-center justify-center ${
                                i === 2 ? 'bg-primary' : 'bg-white/5'
                              }`}
                            >
                              <span className="text-lg">{emoji}</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </motion.div>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section className="py-20 bg-gradient-to-b from-black to-gray-900">
          <div className="container mx-auto px-6">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              className="text-center mb-16"
            >
              <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
                {t('featuresSection.title')}
              </h2>
              <p className="text-gray-400 max-w-2xl mx-auto">
                {t('featuresSection.description')}
              </p>
            </motion.div>

            <motion.div
              variants={staggerContainer}
              initial="initial"
              whileInView="animate"
              viewport={{ once: true }}
              className="grid md:grid-cols-2 lg:grid-cols-4 gap-6"
            >
              {features.map((feature, index) => (
                <motion.div
                  key={index}
                  variants={fadeInUp}
                  className="p-6 bg-white/5 rounded-2xl border border-white/10 hover:border-primary/50 transition-colors group"
                >
                  <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center mb-4 group-hover:bg-primary/20 transition-colors">
                    <feature.icon className="w-6 h-6 text-primary" />
                  </div>
                  <h3 className="text-lg font-semibold text-white mb-2">
                    {feature.title}
                  </h3>
                  <p className="text-gray-400 text-sm">
                    {feature.description}
                  </p>
                </motion.div>
              ))}
            </motion.div>
          </div>
        </section>

        {/* Requirements Section */}
        <section className="py-20 bg-gray-900">
          <div className="container mx-auto px-6">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              className="text-center mb-16"
            >
              <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
                {t('requirements.title')}
              </h2>
              <p className="text-gray-400">
                {t('requirements.description')}
              </p>
            </motion.div>

            <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
              {/* iOS */}
              <motion.div
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                className="p-8 bg-white/5 rounded-2xl border border-white/10"
              >
                <div className="flex items-center gap-4 mb-6">
                  <div className="w-14 h-14 bg-white rounded-xl flex items-center justify-center">
                    <svg className="w-8 h-8" viewBox="0 0 24 24" fill="black">
                      <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
                    </svg>
                  </div>
                  <div>
                    <h3 className="text-xl font-bold text-white">iOS</h3>
                    <p className="text-gray-400 text-sm">App Store</p>
                  </div>
                </div>
                
                <ul className="space-y-3">
                  <li className="flex items-center gap-3 text-gray-300">
                    <Shield className="w-5 h-5 text-primary" />
                    <span>{t('requirements.ios.version')}: {requirements.ios.version}</span>
                  </li>
                  <li className="flex items-center gap-3 text-gray-300">
                    <Download className="w-5 h-5 text-primary" />
                    <span>{t('requirements.ios.size')}: {requirements.ios.size}</span>
                  </li>
                  <li className="flex items-center gap-3 text-gray-300">
                    <Zap className="w-5 h-5 text-primary" />
                    <span>{t('requirements.ios.languages')}: {requirements.ios.languages}</span>
                  </li>
                </ul>
              </motion.div>

              {/* Android */}
              <motion.div
                initial={{ opacity: 0, x: 20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                className="p-8 bg-white/5 rounded-2xl border border-white/10"
              >
                <div className="flex items-center gap-4 mb-6">
                  <div className="w-14 h-14 bg-[#3DDC84] rounded-xl flex items-center justify-center">
                    <svg className="w-8 h-8" viewBox="0 0 24 24" fill="white">
                      <path d="M17.6 11.48V8.35l1.8-3.12a.5.5 0 00-.87-.5l-1.86 3.22a9.65 9.65 0 00-3.67-.69 9.65 9.65 0 00-3.67.69L7.47 4.73a.5.5 0 00-.87.5l1.8 3.12v3.13C4.5 12.25 2 14.83 2 18h20c0-3.17-2.5-5.75-4.4-6.52zm-10.1 2.02a1 1 0 110-2 1 1 0 010 2zm9 0a1 1 0 110-2 1 1 0 010 2z"/>
                    </svg>
                  </div>
                  <div>
                    <h3 className="text-xl font-bold text-white">Android</h3>
                    <p className="text-gray-400 text-sm">Google Play</p>
                  </div>
                </div>
                
                <ul className="space-y-3">
                  <li className="flex items-center gap-3 text-gray-300">
                    <Shield className="w-5 h-5 text-primary" />
                    <span>{t('requirements.android.version')}: {requirements.android.version}</span>
                  </li>
                  <li className="flex items-center gap-3 text-gray-300">
                    <Download className="w-5 h-5 text-primary" />
                    <span>{t('requirements.android.size')}: {requirements.android.size}</span>
                  </li>
                  <li className="flex items-center gap-3 text-gray-300">
                    <Zap className="w-5 h-5 text-primary" />
                    <span>{t('requirements.android.languages')}: {requirements.android.languages}</span>
                  </li>
                </ul>
              </motion.div>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="py-20 bg-black">
          <div className="container mx-auto px-6">
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              className="max-w-4xl mx-auto text-center p-12 bg-gradient-to-br from-primary/20 to-primary/5 rounded-3xl border border-primary/30"
            >
              <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
                {t('cta.title')}
              </h2>
              <p className="text-gray-400 mb-8 max-w-2xl mx-auto">
                {t('cta.description')}
              </p>
              
              <AppStoreBadges size="lg" className="justify-center" />
            </motion.div>
          </div>
        </section>
      </main>
      <Footer />
    </>
  );
}
