# 🚀 Plan de Génération - Site Vitrine

> **Module** : Site de présentation (www.yousoon.com)  
> **Priorité** : 🟢 Basse (peut être fait en parallèle)  
> **Dépendances** : Aucune (site statique)

---

## 📋 Vue d'Ensemble

```
┌─────────────────────────────────────────────────────────────────┐
│                    ORDRE DE GÉNÉRATION                          │
├─────────────────────────────────────────────────────────────────┤
│  Phase 1: Setup Next.js 14                                      │
│  Phase 2: Pages principales                                     │
│  Phase 3: Composants & Animations                               │
│  Phase 4: SEO & Performance                                     │
│  Phase 5: Internationalisation                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📁 Structure Cible

```
apps/vitrine/
├── app/
│   ├── layout.tsx
│   ├── page.tsx                    # Homepage
│   ├── [locale]/
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   ├── about/page.tsx
│   │   ├── features/page.tsx
│   │   ├── pricing/page.tsx
│   │   ├── partners/page.tsx
│   │   ├── contact/page.tsx
│   │   ├── legal/
│   │   │   ├── privacy/page.tsx
│   │   │   ├── terms/page.tsx
│   │   │   └── cookies/page.tsx
│   │   └── blog/
│   │       ├── page.tsx
│   │       └── [slug]/page.tsx
├── components/
│   ├── layout/
│   ├── sections/
│   └── ui/
├── lib/
├── public/
│   ├── images/
│   └── icons/
├── messages/
│   ├── fr.json
│   └── en.json
└── next.config.js
```

---

## 🔷 Phase 1 : Setup

### Étape 1.1 : Initialisation Next.js 14
**Commandes :**
```bash
npx create-next-app@latest apps/vitrine --typescript --tailwind --app
```

**Dépendances :**
```json
{
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.2.0",
    "next-intl": "^3.0.0",
    "framer-motion": "^10.16.0",
    "@vercel/analytics": "^1.1.0"
  }
}
```

---

## 🔷 Phase 2 : Pages Principales

### Étape 2.1 : Homepage
**Sections :**
1. **Hero** : Titre accrocheur + CTA téléchargement app
2. **Problème/Solution** : Pourquoi Yousoon
3. **Comment ça marche** : 3 étapes
4. **Features** : Fonctionnalités clés
5. **Témoignages** : Avis utilisateurs
6. **Partenaires** : Logos partenaires
7. **CTA Final** : Télécharger l'app

### Étape 2.2 : Pages Secondaires
```
app/[locale]/
├── about/page.tsx              # À propos
├── features/page.tsx           # Fonctionnalités détaillées
├── pricing/page.tsx            # Plans abonnement
├── partners/page.tsx           # Devenir partenaire (lien business.yousoon.com)
├── contact/page.tsx            # Formulaire contact
└── legal/
    ├── privacy/page.tsx        # Politique confidentialité
    ├── terms/page.tsx          # CGU
    └── cookies/page.tsx        # Politique cookies
```

---

## 🔷 Phase 3 : Composants

### Étape 3.1 : Layout
```
components/layout/
├── Header.tsx                  # Navigation
├── Footer.tsx
├── MobileMenu.tsx
└── LanguageSwitcher.tsx
```

### Étape 3.2 : Sections Homepage
```
components/sections/
├── Hero.tsx
├── ProblemSolution.tsx
├── HowItWorks.tsx
├── Features.tsx
├── Testimonials.tsx
├── Partners.tsx
├── CtaSection.tsx
└── AppShowcase.tsx             # Screenshots app
```

### Étape 3.3 : UI Components
```
components/ui/
├── Button.tsx
├── Card.tsx
├── Badge.tsx
├── AppStoreButtons.tsx         # iOS + Android
└── AnimatedCounter.tsx
```

---

## 🔷 Phase 4 : SEO & Performance

### Étape 4.1 : Métadonnées
```typescript
// app/layout.tsx
export const metadata: Metadata = {
  title: 'Yousoon - Sortez plus, payez moins',
  description: 'Découvrez des sorties avec réductions exclusives...',
  openGraph: { ... },
  twitter: { ... },
};
```

### Étape 4.2 : Sitemap & Robots
```
app/
├── sitemap.ts
└── robots.ts
```

---

## 🔷 Phase 5 : Internationalisation

### Étape 5.1 : Configuration next-intl
```
messages/
├── fr.json
└── en.json

middleware.ts                   # Locale detection
```

---

## ⏱️ Estimation des Temps

| Phase | Durée estimée |
|-------|---------------|
| Setup | 1h |
| Homepage | 4h |
| Pages secondaires | 3h |
| Composants | 2h |
| SEO | 1h |
| i18n | 1h |
| **Total** | **~12h** |

---

## ✅ Critères de Validation

- [ ] Lighthouse > 95
- [ ] Responsive parfait
- [ ] Animations fluides
- [ ] SEO optimisé
- [ ] FR/EN fonctionnel
- [ ] Links App Store/Play Store

---

## 🔗 Références

- [Prompt Site Vitrine](./PROMPT.md)
- [Design System](../DESIGN_SYSTEM.md)
