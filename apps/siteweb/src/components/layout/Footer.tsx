import Link from 'next/link'
import { useTranslations } from 'next-intl'
import { Instagram, Linkedin, Twitter } from 'lucide-react'
import { Button } from '@/components/ui'
import { SITE_CONFIG, FOOTER_LINKS } from '@/lib/constants'

export function Footer() {
  const t = useTranslations('footer')
  const currentYear = new Date().getFullYear()

  return (
    <footer className="bg-dark-950 border-t border-white/10">
      <div className="container-wide section-padding">
        {/* Main Footer Content */}
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-8 lg:gap-12">
          {/* Brand Column */}
          <div className="col-span-2 lg:col-span-2">
            <Link href="/" className="inline-block mb-4">
              <span className="text-2xl font-bold text-white">
                You<span className="text-primary">soon</span>
              </span>
            </Link>
            <p className="text-grey-jet text-sm mb-6 max-w-xs">
              {t('description')}
            </p>
            
            {/* Social Links */}
            <div className="flex items-center gap-4">
              <a
                href={SITE_CONFIG.links.instagram}
                target="_blank"
                rel="noopener noreferrer"
                className="text-grey-jet hover:text-primary transition-colors"
                aria-label="Instagram"
              >
                <Instagram size={20} />
              </a>
              <a
                href={SITE_CONFIG.links.linkedin}
                target="_blank"
                rel="noopener noreferrer"
                className="text-grey-jet hover:text-primary transition-colors"
                aria-label="LinkedIn"
              >
                <Linkedin size={20} />
              </a>
              <a
                href={SITE_CONFIG.links.twitter}
                target="_blank"
                rel="noopener noreferrer"
                className="text-grey-jet hover:text-primary transition-colors"
                aria-label="Twitter"
              >
                <Twitter size={20} />
              </a>
            </div>
          </div>

          {/* Product Links */}
          <div>
            <h4 className="text-white font-semibold mb-4">{t('product')}</h4>
            <ul className="space-y-3">
              {FOOTER_LINKS.product.map((link) => (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    className="text-grey-jet hover:text-white transition-colors text-sm"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Company Links */}
          <div>
            <h4 className="text-white font-semibold mb-4">{t('company')}</h4>
            <ul className="space-y-3">
              {FOOTER_LINKS.company.map((link) => (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    className="text-grey-jet hover:text-white transition-colors text-sm"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Partners Links */}
          <div>
            <h4 className="text-white font-semibold mb-4">{t('partners')}</h4>
            <ul className="space-y-3">
              {FOOTER_LINKS.partners.map((link) => (
                <li key={link.href}>
                  {'external' in link && link.external ? (
                    <a
                      href={link.href}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-grey-jet hover:text-white transition-colors text-sm"
                    >
                      {link.label}
                    </a>
                  ) : (
                    <Link
                      href={link.href}
                      className="text-grey-jet hover:text-white transition-colors text-sm"
                    >
                      {link.label}
                    </Link>
                  )}
                </li>
              ))}
            </ul>
          </div>

          {/* Legal Links */}
          <div>
            <h4 className="text-white font-semibold mb-4">{t('legal')}</h4>
            <ul className="space-y-3">
              {FOOTER_LINKS.legal.map((link) => (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    className="text-grey-jet hover:text-white transition-colors text-sm"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Newsletter Section */}
        <div className="mt-12 pt-8 border-t border-white/10">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-6">
            <div>
              <h4 className="text-white font-semibold mb-1">{t('newsletter.title')}</h4>
              <p className="text-grey-jet text-sm">{t('newsletter.description')}</p>
            </div>
            <form className="flex gap-2 w-full md:w-auto">
              <input
                type="email"
                placeholder={t('newsletter.placeholder')}
                className="flex-1 md:w-64 px-4 py-2 bg-dark-800 border border-dark-700 rounded-lg text-white placeholder:text-grey-jet focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              />
              <Button type="submit">{t('newsletter.button')}</Button>
            </form>
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="mt-12 pt-8 border-t border-white/10 flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <p className="text-grey-jet text-sm">
            {t('copyright', { year: currentYear })}
          </p>
          
          {/* App Store Badges */}
          <div className="flex items-center gap-4">
            <a
              href={SITE_CONFIG.links.appStore}
              target="_blank"
              rel="noopener noreferrer"
              className="opacity-80 hover:opacity-100 transition-opacity"
            >
              <img
                src="/images/app-store-badge.svg"
                alt="Download on App Store"
                className="h-10"
              />
            </a>
            <a
              href={SITE_CONFIG.links.playStore}
              target="_blank"
              rel="noopener noreferrer"
              className="opacity-80 hover:opacity-100 transition-opacity"
            >
              <img
                src="/images/play-store-badge.svg"
                alt="Get it on Google Play"
                className="h-10"
              />
            </a>
          </div>
        </div>
      </div>
    </footer>
  )
}
