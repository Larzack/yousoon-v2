# 🚀 Plan de Génération - Site Partenaires

> **Module** : Portail Partenaires (business.yousoon.com)  
> **Priorité** : 🟡 Moyenne (après Backend + App Mobile)  
> **Dépendances** : Backend Gateway + Partner Service + Discovery Service

---

## 📋 Vue d'Ensemble

```
┌─────────────────────────────────────────────────────────────────┐
│                    ORDRE DE GÉNÉRATION                          │
├─────────────────────────────────────────────────────────────────┤
│  Phase 1: Setup & Configuration (Vite, TailwindCSS, shadcn)    │
│  Phase 2: Layout & Navigation                                   │
│  Phase 3: Auth & Onboarding                                     │
│  Phase 4: Dashboard & Analytics                                 │
│  Phase 5: Gestion Offres                                        │
│  Phase 6: Gestion Établissements & Équipe                       │
│  Phase 7: Réservations & Settings                               │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📁 Structure Cible

```
apps/partner-portal/
├── public/
│   └── assets/
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── routes/
│   ├── components/
│   │   ├── ui/              # shadcn/ui
│   │   ├── layout/
│   │   └── shared/
│   ├── lib/
│   │   ├── graphql/
│   │   ├── utils/
│   │   └── validators/
│   ├── hooks/
│   ├── stores/
│   ├── types/
│   └── styles/
├── tests/
├── package.json
├── vite.config.ts
├── tailwind.config.ts
└── tsconfig.json
```

---

## 🔷 Phase 1 : Setup & Configuration

### Étape 1.1 : Initialisation Projet
**Commandes :**
```bash
npm create vite@latest apps/partner-portal -- --template react-ts
cd apps/partner-portal
npm install
```

**Fichiers à générer :**
```
apps/partner-portal/
├── package.json
├── vite.config.ts
├── tsconfig.json
├── tailwind.config.ts
├── postcss.config.js
└── .eslintrc.cjs
```

### Étape 1.2 : Dépendances
**package.json :**
```json
{
  "dependencies": {
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "react-router-dom": "^6.20.0",
    "@tanstack/react-query": "^5.0.0",
    "zustand": "^4.4.0",
    "urql": "^4.0.0",
    "graphql": "^16.8.0",
    "react-hook-form": "^7.48.0",
    "@hookform/resolvers": "^3.3.0",
    "zod": "^3.22.0",
    "date-fns": "^2.30.0",
    "recharts": "^2.10.0",
    "@tanstack/react-table": "^8.10.0",
    "lucide-react": "^0.294.0",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.0.0",
    "tailwind-merge": "^2.0.0"
  }
}
```

### Étape 1.3 : Configuration shadcn/ui
**Fichiers à générer :**
```
src/components/ui/
├── button.tsx
├── input.tsx
├── label.tsx
├── card.tsx
├── dialog.tsx
├── dropdown-menu.tsx
├── select.tsx
├── table.tsx
├── tabs.tsx
├── toast.tsx
├── form.tsx
└── ...
```

### Étape 1.4 : Configuration GraphQL
**Fichiers à générer :**
```
src/lib/graphql/
├── client.ts                   # urql client
├── queries/
│   ├── partner.ts
│   ├── offers.ts
│   ├── establishments.ts
│   ├── bookings.ts
│   └── analytics.ts
├── mutations/
│   ├── auth.ts
│   ├── offers.ts
│   ├── establishments.ts
│   └── team.ts
└── types/
    └── generated.ts            # Types auto-générés
```

---

## 🔷 Phase 2 : Layout & Navigation

### Étape 2.1 : Layout Principal
**Fichiers à générer :**
```
src/components/layout/
├── RootLayout.tsx
├── DashboardLayout.tsx
│   ├── Sidebar.tsx
│   ├── Header.tsx
│   └── MobileNav.tsx
├── AuthLayout.tsx
└── Footer.tsx
```

### Étape 2.2 : Configuration Routes
**Fichiers à générer :**
```
src/routes/
├── index.tsx                   # Route config
├── auth/
│   ├── login.tsx
│   ├── register.tsx
│   └── forgot-password.tsx
├── dashboard/
│   └── index.tsx
├── offers/
│   ├── index.tsx
│   ├── create.tsx
│   └── [id].tsx
├── establishments/
│   ├── index.tsx
│   └── [id].tsx
├── analytics/
│   └── index.tsx
├── bookings/
│   └── index.tsx
└── settings/
    ├── profile.tsx
    ├── team.tsx
    └── billing.tsx
```

### Étape 2.3 : Navigation Sidebar
**Menu items :**
```typescript
const menuItems = [
  { icon: LayoutDashboard, label: 'Dashboard', href: '/' },
  { icon: Tag, label: 'Offres', href: '/offers' },
  { icon: Building2, label: 'Établissements', href: '/establishments' },
  { icon: Calendar, label: 'Réservations', href: '/bookings' },
  { icon: BarChart3, label: 'Analytics', href: '/analytics' },
  { icon: Settings, label: 'Paramètres', href: '/settings' },
];
```

---

## 🔷 Phase 3 : Auth & Onboarding

### Étape 3.1 : Pages Auth
**Fichiers à générer :**
```
src/routes/auth/
├── login.tsx
│   - Email + Password
│   - "Se souvenir de moi"
│   - Lien mot de passe oublié
│   - 2FA (TOTP)
├── register.tsx
│   - Step 1: Infos entreprise (Raison sociale, SIRET)
│   - Step 2: Adresse
│   - Step 3: Contact admin
│   - Step 4: Vérification email
└── forgot-password.tsx
```

### Étape 3.2 : Stores Auth
**Fichiers à générer :**
```
src/stores/
├── authStore.ts                # Zustand store
└── useAuth.ts                  # Hook
```

### Étape 3.3 : Composants Auth
**Fichiers à générer :**
```
src/components/auth/
├── LoginForm.tsx
├── RegisterForm.tsx
├── TwoFactorForm.tsx
└── PasswordResetForm.tsx
```

---

## 🔷 Phase 4 : Dashboard & Analytics

### Étape 4.1 : Dashboard Principal
**Fichiers à générer :**
```
src/routes/dashboard/
└── index.tsx
```

**KPIs affichés :**
- Offres actives
- Vues totales (30j)
- Réservations (30j)
- Taux de conversion
- Check-ins réalisés

### Étape 4.2 : Composants Dashboard
**Fichiers à générer :**
```
src/components/dashboard/
├── StatsCard.tsx
├── RecentBookings.tsx
├── TopOffers.tsx
├── QuickActions.tsx
└── AlertsWidget.tsx
```

### Étape 4.3 : Analytics Page
**Fichiers à générer :**
```
src/routes/analytics/
└── index.tsx

src/components/analytics/
├── ChartContainer.tsx
├── LineChart.tsx               # Évolution 365j
├── BarChart.tsx                # Comparaison offres
├── HeatmapCalendar.tsx         # Fréquentation
├── MetricsCards.tsx
├── DateRangePicker.tsx
├── FilterBar.tsx
└── ExportButton.tsx            # CSV, PDF
```

**Graphiques :**
- **Ligne** : Évolution des réservations/vues sur 365 jours
- **Heatmap** : Fréquentation par jour
- **Barres** : Comparaison entre offres
- **Prévisions** : Jours à venir (basé sur historique)

---

## 🔷 Phase 5 : Gestion Offres

### Étape 5.1 : Liste des Offres
**Fichiers à générer :**
```
src/routes/offers/
├── index.tsx                   # Liste avec DataTable
├── create.tsx                  # Création wizard
└── [id].tsx                    # Détail/Édition
```

### Étape 5.2 : Composants Offres
**Fichiers à générer :**
```
src/components/offers/
├── OfferTable.tsx              # TanStack Table
├── OfferCard.tsx
├── OfferForm/
│   ├── index.tsx               # Wizard multi-steps
│   ├── Step1GeneralInfo.tsx
│   ├── Step2Discount.tsx
│   ├── Step3Validity.tsx
│   ├── Step4Media.tsx
│   └── Step5Preview.tsx
├── OfferStatusBadge.tsx
├── OfferActions.tsx
└── OfferPreview.tsx            # Aperçu mobile
```

### Étape 5.3 : Formulaire Création (Wizard)
**Steps :**
1. **Informations générales** : Titre, description, catégorie, établissement
2. **Réduction** : Type (%, fixe, formule), valeur, conditions
3. **Validité** : Dates, jours, créneaux horaires, quota
4. **Médias** : Images (drag & drop), vidéo optionnelle
5. **Prévisualisation** : Aperçu mobile + publication

---

## 🔷 Phase 6 : Établissements & Équipe

### Étape 6.1 : Gestion Établissements
**Fichiers à générer :**
```
src/routes/establishments/
├── index.tsx
└── [id].tsx

src/components/establishments/
├── EstablishmentCard.tsx
├── EstablishmentForm.tsx
├── OpeningHoursEditor.tsx
├── LocationPicker.tsx          # Google Maps
└── PhotoGallery.tsx
```

### Étape 6.2 : Gestion Équipe
**Fichiers à générer :**
```
src/routes/settings/
└── team.tsx

src/components/team/
├── TeamTable.tsx
├── InviteMemberDialog.tsx
├── RoleSelect.tsx
└── MemberActions.tsx
```

**Rôles :**
- Admin : Accès complet
- Manager : Gestion offres et stats
- Viewer : Consultation uniquement

---

## 🔷 Phase 7 : Réservations & Settings

### Étape 7.1 : Liste Réservations
**Fichiers à générer :**
```
src/routes/bookings/
└── index.tsx

src/components/bookings/
├── BookingTable.tsx
├── BookingFilters.tsx
├── BookingStatusBadge.tsx
├── BookingActions.tsx
└── ManualCheckinDialog.tsx
```

### Étape 7.2 : Settings
**Fichiers à générer :**
```
src/routes/settings/
├── profile.tsx                 # Infos entreprise
├── team.tsx                    # Gestion équipe
├── billing.tsx                 # Facturation (future)
└── notifications.tsx           # Préférences

src/components/settings/
├── ProfileForm.tsx
├── BrandingUpload.tsx
├── NotificationPreferences.tsx
└── DangerZone.tsx              # Suppression compte
```

---

## ⏱️ Estimation des Temps

| Phase | Étape | Durée estimée |
|-------|-------|---------------|
| **Phase 1** | Setup projet | 1h |
| | shadcn/ui | 1h |
| | GraphQL config | 1h |
| **Phase 2** | Layout | 2h |
| | Navigation | 1h |
| **Phase 3** | Auth pages | 3h |
| | 2FA | 1h |
| **Phase 4** | Dashboard | 2h |
| | Analytics | 4h |
| **Phase 5** | Liste offres | 2h |
| | Formulaire wizard | 4h |
| **Phase 6** | Établissements | 3h |
| | Équipe | 2h |
| **Phase 7** | Réservations | 2h |
| | Settings | 2h |
| **Total** | | **~31h** |

---

## ✅ Critères de Validation

### UI/UX
- [ ] Responsive (desktop-first, tablet support)
- [ ] Accessibilité WCAG 2.1 AA
- [ ] Dark mode optionnel
- [ ] Loading states partout
- [ ] Error handling gracieux

### Features
- [ ] Inscription partenaire complète
- [ ] 2FA fonctionnel
- [ ] CRUD offres complet
- [ ] CRUD établissements
- [ ] Analytics avec graphiques
- [ ] Export CSV/PDF

### Performance
- [ ] Lighthouse > 90
- [ ] First contentful paint < 1.5s
- [ ] Time to interactive < 3s

### Tests
- [ ] Tests unitaires hooks
- [ ] Tests intégration formulaires
- [ ] E2E création offre

---

## 🔗 Références

- [Prompt Site Partenaires](./PROMPT.md)
- [Data Model](../DATA_MODEL.md)
- [Backend API](../backend/ARCHITECTURE.md)
