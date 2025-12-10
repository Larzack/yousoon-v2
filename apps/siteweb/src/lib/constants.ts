export const SITE_CONFIG = {
  name: 'Yousoon',
  description: 'Sortez plus, payez moins. Découvrez les meilleures offres de sorties près de chez vous.',
  url: 'https://www.yousoon.com',
  ogImage: 'https://www.yousoon.com/og-image.jpg',
  links: {
    appStore: 'https://apps.apple.com/app/yousoon',
    playStore: 'https://play.google.com/store/apps/details?id=com.yousoon.yousoon',
    partnerPortal: 'https://business.yousoon.com',
    instagram: 'https://instagram.com/yousoon',
    linkedin: 'https://linkedin.com/company/yousoon',
    twitter: 'https://twitter.com/yousoon',
  },
  contact: {
    email: 'contact@yousoon.com',
    support: 'support@yousoon.com',
  },
}

export const STATS = {
  users: 50000,
  partners: 500,
  savings: 30, // Pourcentage économisé
  cities: 15,
}

export const NAV_ITEMS = [
  { href: '/fonctionnalites', label: 'Fonctionnalités', key: 'features' },
  { href: '/partenaires', label: 'Devenir partenaire', key: 'partners' },
  { href: '/tarifs', label: 'Tarifs', key: 'pricing' },
  { href: '/a-propos', label: 'À propos', key: 'about' },
]

export const FOOTER_LINKS = {
  product: [
    { href: '/fonctionnalites', label: 'Fonctionnalités' },
    { href: '/tarifs', label: 'Tarifs' },
    { href: '/faq', label: 'FAQ' },
  ],
  company: [
    { href: '/a-propos', label: 'À propos' },
    { href: '/blog', label: 'Blog' },
    { href: '/contact', label: 'Contact' },
    { href: '/presse', label: 'Presse' },
  ],
  partners: [
    { href: '/partenaires', label: 'Devenir partenaire' },
    { href: 'https://business.yousoon.com', label: 'Espace partenaire', external: true },
  ],
  legal: [
    { href: '/mentions-legales', label: 'Mentions légales' },
    { href: '/politique-confidentialite', label: 'Confidentialité' },
    { href: '/cgv', label: 'CGV' },
    { href: '/cookies', label: 'Cookies' },
  ],
}

export const FEATURES = [
  {
    icon: 'MapPin',
    title: 'Offres géolocalisées',
    description: 'Découvrez les meilleures réductions autour de vous en temps réel.',
  },
  {
    icon: 'Percent',
    title: 'Réductions exclusives',
    description: 'Profitez d\'offres uniques négociées avec nos partenaires.',
  },
  {
    icon: 'Zap',
    title: 'Réservation instantanée',
    description: 'Réservez en un clic et validez sur place avec un QR code.',
  },
  {
    icon: 'Heart',
    title: 'Favoris personnalisés',
    description: 'Sauvegardez vos offres préférées et recevez des alertes.',
  },
  {
    icon: 'Users',
    title: 'Communauté active',
    description: 'Partagez vos expériences et découvrez les avis.',
  },
  {
    icon: 'Shield',
    title: 'Paiement sécurisé',
    description: 'Abonnement géré via Apple Pay et Google Pay.',
  },
]

export const HOW_IT_WORKS_STEPS = [
  {
    step: 1,
    title: 'Téléchargez l\'app',
    description: 'Disponible gratuitement sur iOS et Android.',
    icon: 'Download',
  },
  {
    step: 2,
    title: 'Explorez les offres',
    description: 'Parcourez les réductions près de vous ou dans une ville.',
    icon: 'Search',
  },
  {
    step: 3,
    title: 'Réservez',
    description: 'Sélectionnez une offre et réservez instantanément.',
    icon: 'CalendarCheck',
  },
  {
    step: 4,
    title: 'Profitez !',
    description: 'Présentez votre QR code sur place et profitez.',
    icon: 'PartyPopper',
  },
]

export const TESTIMONIALS = [
  {
    id: 1,
    name: 'Marie D.',
    role: 'Utilisatrice depuis 6 mois',
    avatar: '/avatars/user-1.jpg',
    content: 'J\'ai économisé plus de 200€ en 3 mois ! Les offres sont vraiment intéressantes et variées.',
    rating: 5,
  },
  {
    id: 2,
    name: 'Thomas L.',
    role: 'Yousooner fidèle',
    avatar: '/avatars/user-2.jpg',
    content: 'Super app pour découvrir de nouveaux endroits. L\'interface est top et les réductions valent le coup.',
    rating: 5,
  },
  {
    id: 3,
    name: 'Sophie M.',
    role: 'Grande voyageuse',
    avatar: '/avatars/user-3.jpg',
    content: 'Parfait quand je voyage ! Je trouve toujours des bons plans dans chaque ville.',
    rating: 4,
  },
]

export const PARTNER_BENEFITS = [
  {
    icon: 'TrendingUp',
    title: 'Augmentez votre visibilité',
    description: 'Touchez des milliers de nouveaux clients potentiels dans votre zone.',
  },
  {
    icon: 'Users',
    title: 'Attirez de nouveaux clients',
    description: 'Nos utilisateurs recherchent activement des expériences uniques.',
  },
  {
    icon: 'BarChart3',
    title: 'Analysez vos performances',
    description: 'Accédez à des statistiques détaillées sur vos offres.',
  },
  {
    icon: 'Settings',
    title: 'Gérez facilement',
    description: 'Interface intuitive pour créer et gérer vos offres.',
  },
]

export const FAQ_ITEMS = [
  {
    question: 'Comment fonctionne Yousoon ?',
    answer: 'Yousoon est une application mobile qui vous permet de découvrir des offres et réductions exclusives chez nos partenaires (bars, restaurants, loisirs). Téléchargez l\'app, créez votre compte, et réservez des offres près de chez vous.',
  },
  {
    question: 'L\'application est-elle gratuite ?',
    answer: 'L\'application est gratuite au téléchargement. Nous proposons un abonnement premium pour accéder à toutes les offres et bénéficier d\'avantages exclusifs. Une période d\'essai gratuite de 30 jours est offerte.',
  },
  {
    question: 'Comment utiliser une offre ?',
    answer: 'Une fois une offre réservée, rendez-vous chez le partenaire et présentez votre QR code depuis l\'application. Le personnel validera votre réservation et vous pourrez profiter de la réduction.',
  },
  {
    question: 'Puis-je annuler une réservation ?',
    answer: 'Oui, vous pouvez annuler une réservation jusqu\'à 30 minutes après la réservation, tant que vous n\'avez pas encore fait le check-in.',
  },
  {
    question: 'Comment devenir partenaire ?',
    answer: 'Inscrivez-vous gratuitement sur notre portail partenaires (business.yousoon.com). Notre équipe validera votre demande et vous pourrez commencer à publier vos offres.',
  },
]

export const PRICING_PLANS = [
  {
    name: 'Gratuit',
    price: 0,
    period: '',
    description: 'Pour découvrir Yousoon',
    features: [
      '5 réservations / mois',
      'Accès aux offres standard',
      'Historique des sorties',
    ],
    cta: 'Commencer gratuitement',
    highlighted: false,
  },
  {
    name: 'Explorer',
    price: 9.99,
    period: '/mois',
    description: 'Pour les sorteurs réguliers',
    features: [
      'Réservations illimitées',
      'Accès à toutes les offres',
      'Offres exclusives',
      'Support prioritaire',
      '30 jours d\'essai gratuit',
    ],
    cta: 'Essayer gratuitement',
    highlighted: true,
  },
  {
    name: 'Voyager',
    price: 79.99,
    period: '/an',
    description: 'Meilleur rapport qualité/prix',
    features: [
      'Tout Explorer inclus',
      '2 mois offerts',
      'Accès anticipé nouvelles offres',
      'Événements VIP',
    ],
    cta: 'Économiser 40%',
    highlighted: false,
  },
]
