'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useTranslations } from 'next-intl'
import { motion, AnimatePresence } from 'framer-motion'
import { Menu, X, ChevronDown } from 'lucide-react'
import { Button } from '@/components/ui'
import { cn } from '@/lib/utils'
import { SITE_CONFIG, NAV_ITEMS } from '@/lib/constants'

export function Header() {
  const t = useTranslations('nav')
  const [isScrolled, setIsScrolled] = useState(false)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 20)
    }
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  return (
    <header
      className={cn(
        'fixed top-0 left-0 right-0 z-50 transition-all duration-300',
        isScrolled
          ? 'bg-black/90 backdrop-blur-md border-b border-white/10 py-3'
          : 'bg-transparent py-5'
      )}
    >
      <div className="container-wide">
        <nav className="flex items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2">
            <span className="text-2xl font-bold text-white">
              You<span className="text-primary">soon</span>
            </span>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden lg:flex items-center gap-8">
            {NAV_ITEMS.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="text-grey-eerie hover:text-white transition-colors text-sm font-medium"
              >
                {t(item.label.toLowerCase().replace(/[^a-z]/g, ''))}
              </Link>
            ))}
          </div>

          {/* Desktop CTA */}
          <div className="hidden lg:flex items-center gap-4">
            <Link
              href={SITE_CONFIG.links.partnerPortal}
              className="text-grey-eerie hover:text-white transition-colors text-sm font-medium"
            >
              {t('login')}
            </Link>
            <Button size="sm">{t('downloadApp')}</Button>
          </div>

          {/* Mobile Menu Button */}
          <button
            className="lg:hidden p-2 text-white"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            aria-label="Toggle menu"
          >
            {isMobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
          </button>
        </nav>
      </div>

      {/* Mobile Menu */}
      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="lg:hidden bg-black/95 backdrop-blur-lg border-t border-white/10"
          >
            <div className="container-wide py-6 flex flex-col gap-4">
              {NAV_ITEMS.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="text-grey-eerie hover:text-white transition-colors py-2 text-lg"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  {t(item.label.toLowerCase().replace(/[^a-z]/g, ''))}
                </Link>
              ))}
              <div className="flex flex-col gap-3 pt-4 border-t border-white/10">
                <Link
                  href={SITE_CONFIG.links.partnerPortal}
                  className="text-grey-eerie hover:text-white transition-colors py-2"
                >
                  {t('login')}
                </Link>
                <Button className="w-full">{t('downloadApp')}</Button>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  )
}
