# 🌐 Site Vitrine - Prompt Détaillé

> **Module** : Site de présentation Yousoon  
> **URL** : www.yousoon.com  
> **Technologie** : Next.js 14  
> **Figma** : [Yousoon-Test2](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2?node-id=121-114)  
> **i18n** : FR, EN + extensible multi-langue  
> **Hébergement** : AKS (Azure)

---

## 🎯 Objectifs

Le site vitrine doit :
- Présenter l'application Yousoon aux utilisateurs potentiels
- Convaincre les partenaires de rejoindre la plateforme
- Optimiser le SEO pour l'acquisition organique
- Rediriger vers les stores (App Store, Play Store)
- Rediriger vers le portail partenaires

---

## 🛠️ Stack Technique

### Stack Confirmée : Next.js
| Technologie | Justification |
|-------------|---------------|
| **Next.js 14** | SSR/SSG, SEO optimal, performance, i18n natif |
| **React 18** | Composants |
| **TypeScript** | Type safety |
| **TailwindCSS** | Styling rapide |
| **Framer Motion** | Animations |
| **next-intl** | Internationalisation |

### Services Externes
| Technologie | Usage |
|-------------|-------|
| **S3 + CloudFront** | CDN images/assets |
| **Google Analytics 4** | Analytics (ou Amplitude) |

---

## 🏗️ Architecture

### Structure du Projet

```
apps/web-vitrine/
├── public/
│   ├── images/
│   ├── fonts/
│   └── favicon.ico
├── src/
│   ├── app/                      # App Router Next.js 14
│   │   ├── layout.tsx
│   │   ├── page.tsx              # Home
│   │   ├── fonctionnalites/
│   │   │   └── page.tsx
│   │   ├── partenaires/
│   │   │   └── page.tsx
│   │   ├── tarifs/
│   │   │   └── page.tsx
│   │   ├── a-propos/
│   │   │   └── page.tsx
│   │   ├── contact/
│   │   │   └── page.tsx
│   │   ├── blog/
│   │   │   ├── page.tsx
│   │   │   └── [slug]/
│   │   │       └── page.tsx
│   │   ├── mentions-legales/
│   │   │   └── page.tsx
│   │   ├── politique-confidentialite/
│   │   │   └── page.tsx
│   │   └── cgv/
│   │       └── page.tsx
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Header.tsx
│   │   │   ├── Footer.tsx
│   │   │   └── Navigation.tsx
│   │   ├── sections/
│   │   │   ├── Hero.tsx
│   │   │   ├── Features.tsx
│   │   │   ├── HowItWorks.tsx
│   │   │   ├── Testimonials.tsx
│   │   │   ├── Partners.tsx
│   │   │   ├── Pricing.tsx
│   │   │   ├── FAQ.tsx
│   │   │   ├── CTA.tsx
│   │   │   └── Newsletter.tsx
│   │   ├── ui/
│   │   │   ├── Button.tsx
│   │   │   ├── Card.tsx
│   │   │   └── Badge.tsx
│   │   └── shared/
│   │       ├── AppStoreBadges.tsx
│   │       ├── PhoneMockup.tsx
│   │       └── AnimatedCounter.tsx
│   ├── lib/
│   │   ├── utils.ts
│   │   └── constants.ts
│   ├── hooks/
│   ├── styles/
│   │   └── globals.css
│   └── types/
├── content/                      # Contenu Markdown (si pas de CMS)
│   ├── blog/
│   └── faq/
├── next.config.js
├── tailwind.config.ts
└── package.json
```

---

## 📄 Pages & Sections

### 1. Page d'Accueil (/)

#### Hero Section
- **Headline accrocheur** : Ex. "Sortez plus, payez moins"
- **Sous-titre** : Proposition de valeur claire
- **Mockup iPhone/Android** : Aperçu de l'app
- **CTA principal** : "Télécharger l'app"
- **Badges stores** : App Store + Play Store
- **Statistiques** : X utilisateurs, Y partenaires, Z% économisés

#### Comment ça marche
Processus en 3-4 étapes :
1. Téléchargez l'app
2. Découvrez les offres autour de vous
3. Réservez et profitez de la réduction
4. Partagez votre expérience

#### Fonctionnalités clés
- Offres exclusives
- Géolocalisation
- Réservation instantanée
- Validation facile
- Favoris et recommandations

#### Témoignages
- Carrousel avis utilisateurs
- Notes App Store/Play Store

#### Partenaires
- Logos partenaires (avec accord)
- "Ils nous font confiance"
- Nombre de partenaires

#### CTA Final
- Répétition du call-to-action téléchargement

---

### 2. Page Fonctionnalités (/fonctionnalites)

Détail des fonctionnalités avec visuels :
- Découverte des offres
- Filtres et recherche
- Carte interactive
- Réservation
- QR Code check-in
- Historique et favoris
- Notifications personnalisées

---

### 3. Section Partenaires (/partenaires)

Page dédiée à l'acquisition B2B :

#### Hero B2B
- Headline : "Attirez de nouveaux clients"
- Proposition de valeur partenaire

#### Avantages partenaires
- Visibilité locale
- Acquisition de nouveaux clients
- Fidélisation
- Analytics détaillés
- Interface simple

#### Comment devenir partenaire
1. Inscription gratuite
2. Création de votre profil
3. Publication de vos offres
4. Accueil des clients

#### Tarification (si applicable)
- Plans et prix
- Comparatif fonctionnalités

#### Témoignages partenaires
- Success stories
- Chiffres clés

#### CTA
- "Devenir partenaire" → business.yousoon.com

---

### 4. Page Tarifs (/tarifs) - Optionnel

Si modèle freemium :
- Comparatif des plans
- FAQ tarification
- CTA inscription

---

### 5. Page À propos (/a-propos)

- Histoire de Yousoon
- Mission et vision
- Équipe (photos + bios)
- Valeurs
- Chiffres clés

---

### 6. Page Contact (/contact)

- Formulaire de contact
- Email support
- Réseaux sociaux
- FAQ rapide

---

### 7. Blog (/blog) - Optionnel

- Articles SEO (sorties, bons plans, guides ville)
- Actualités Yousoon
- Nouveaux partenaires

---

### 8. Pages légales

- Mentions légales
- Politique de confidentialité
- CGU/CGV
- Gestion des cookies

---

## 🎨 Design & Animations

### Principes
- Design moderne et épuré
- Couleurs vives (cohérent avec l'app)
- Illustrations ou photos lifestyle
- Mockups iPhone/Android réalistes

### Animations suggérées
- Hero : Fade-in progressif
- Scroll : Reveal animations (Framer Motion)
- Compteurs : Animation des chiffres
- Carrousel : Smooth sliding
- Hover : Micro-interactions boutons/cartes

### Responsive
- Mobile-first
- Breakpoints : sm(640), md(768), lg(1024), xl(1280)

---

## 🔍 SEO

### Optimisations

- **Meta tags** : Title, description par page
- **Open Graph** : Preview réseaux sociaux
- **Sitemap XML** : Auto-généré
- **Robots.txt** : Configuré
- **Schema.org** : Structured data (Organization, App)
- **Performance** : Core Web Vitals optimisés
- **Images** : WebP, lazy loading, alt tags

### Mots-clés cibles
- "réductions sorties [ville]"
- "bons plans bars restaurants"
- "application réductions sorties"
- "offres happy hour [ville]"

---

## 📊 Intégrations

- **Google Analytics 4** : Tracking
- **Google Tag Manager** : Gestion tags
- **Hotjar/Clarity** : Heatmaps (optionnel)
- **Newsletter** : Mailchimp/SendGrid
- **Formulaire contact** : Email ou CRM

---

## 🧪 Tests

### Tests E2E (Playwright)
- Navigation complète
- Formulaire contact
- Liens stores
- Responsive check

### Tests Performance
- Lighthouse > 95
- Core Web Vitals verts

---

## 📋 Checklist

- [ ] Responsive parfait
- [ ] Lighthouse > 95
- [ ] SEO optimisé
- [ ] Accessibilité WCAG AA
- [ ] HTTPS obligatoire
- [ ] Cookies consent (RGPD)
- [ ] Analytics configuré
- [ ] Liens stores fonctionnels
- [ ] Redirection business.yousoon.com

---

## 🔗 Références

- [Questions à clarifier](./QUESTIONS.md)
- [Design Figma](TODO)
