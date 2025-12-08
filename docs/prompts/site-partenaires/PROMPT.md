# 💼 Site Partenaires - Prompt Détaillé

> **Module** : Portail Partenaires Yousoon  
> **URL** : business.yousoon.com  
> **Technologie** : React + TypeScript  
> **Figma** : [Yousoon-Test2](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2?node-id=121-114)  
> **i18n** : FR, EN + extensible multi-langue

---

## 🎯 Objectifs

Le portail partenaires permet aux fournisseurs (bars, restaurants, organismes de sorties) de :
- Créer et gérer leurs offres/réductions
- Consulter les statistiques d'utilisation
- Gérer leurs établissements
- Suivre les réservations/check-ins
- Gérer leur profil et paramètres

---

## 🛠️ Stack Technique

### Core
| Technologie | Version | Justification |
|-------------|---------|---------------|
| React | 19.x | Dernière version stable, Server Components ready |
| TypeScript | 5.x | Type safety, maintenabilité |
| Vite | 5.x | Build rapide, HMR instantané |

### State Management
| Technologie | Usage |
|-------------|-------|
| **Zustand** | État global léger |
| **TanStack Query** | Cache serveur, sync API |

### Communication API
| Technologie | Usage |
|-------------|-------|
| **urql** ou **Apollo Client** | Client GraphQL |
| **graphql-codegen** | Types auto-générés |

### UI/Styling
| Technologie | Usage |
|-------------|-------|
| **TailwindCSS** | Utility-first CSS |
| **shadcn/ui** | Composants accessibles |
| **Radix UI** | Primitives headless |
| **Lucide Icons** | Iconographie |

### Formulaires
| Technologie | Usage |
|-------------|-------|
| **React Hook Form** | Gestion formulaires |
| **Zod** | Validation schémas |

### Tests
| Type | Technologie |
|------|-------------|
| Unit | Vitest |
| Component | Testing Library |
| E2E | Playwright |

### Autres
| Technologie | Usage |
|-------------|-------|
| **React Router** | Routing |
| **date-fns** | Manipulation dates |
| **recharts** | Graphiques |
| **react-table** | Tableaux |

---

## 🏗️ Architecture

### Structure du Projet

```
apps/web-partner/
├── public/
│   └── assets/
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── routes/
│   │   ├── index.tsx
│   │   ├── auth/
│   │   │   ├── login.tsx
│   │   │   ├── register.tsx
│   │   │   └── forgot-password.tsx
│   │   ├── dashboard/
│   │   │   └── index.tsx
│   │   ├── offers/
│   │   │   ├── index.tsx
│   │   │   ├── [id].tsx
│   │   │   └── create.tsx
│   │   ├── establishments/
│   │   │   ├── index.tsx
│   │   │   └── [id].tsx
│   │   ├── analytics/
│   │   │   └── index.tsx
│   │   ├── bookings/
│   │   │   └── index.tsx
│   │   └── settings/
│   │       ├── profile.tsx
│   │       ├── billing.tsx
│   │       └── team.tsx
│   ├── components/
│   │   ├── ui/                    # Composants shadcn/ui
│   │   ├── layout/
│   │   │   ├── Sidebar.tsx
│   │   │   ├── Header.tsx
│   │   │   └── DashboardLayout.tsx
│   │   ├── offers/
│   │   │   ├── OfferForm.tsx
│   │   │   ├── OfferCard.tsx
│   │   │   └── OfferList.tsx
│   │   ├── analytics/
│   │   │   ├── StatsCard.tsx
│   │   │   └── Charts.tsx
│   │   └── shared/
│   │       ├── DataTable.tsx
│   │       └── FileUpload.tsx
│   ├── lib/
│   │   ├── graphql/
│   │   │   ├── client.ts
│   │   │   ├── queries/
│   │   │   └── mutations/
│   │   ├── utils/
│   │   └── validators/
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── useOffers.ts
│   │   └── useAnalytics.ts
│   ├── stores/
│   │   ├── authStore.ts
│   │   └── uiStore.ts
│   ├── types/
│   │   └── index.ts
│   └── styles/
│       └── globals.css
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── package.json
├── tsconfig.json
├── vite.config.ts
├── tailwind.config.ts
└── playwright.config.ts
```

---

## 📱 Fonctionnalités par Écran

### 1. Authentification

#### Login
- Email + mot de passe
- "Se souvenir de moi"
- Lien mot de passe oublié
- Redirection SSO (optionnel)

#### Register
- Informations entreprise
  - Raison sociale
  - SIRET
  - Adresse
  - Téléphone
- Informations administrateur
  - Nom, prénom
  - Email
  - Mot de passe
- Validation email

#### Forgot Password
- Saisie email
- Envoi lien reset
- Nouveau mot de passe

---

### 2. Dashboard

Tableau de bord avec KPIs :
- Nombre d'offres actives
- Vues totales
- Réservations du mois
- Taux de conversion
- Graphique évolution

Widgets :
- Dernières réservations
- Offres les plus vues
- Alertes/notifications

---

### 3. Gestion des Offres

#### Liste des offres
- Tableau avec filtres (statut, établissement, catégorie)
- Actions : voir, éditer, dupliquer, archiver, supprimer
- Statuts : brouillon, en attente, active, expirée, archivée

#### Création/Édition offre
Formulaire multi-étapes :

**Étape 1 - Informations générales**
- Titre de l'offre
- Description (rich text)
- Catégorie
- Établissement(s) concerné(s)

**Étape 2 - Réduction**
- Type de réduction (pourcentage, montant fixe, formule)
- Valeur de la réduction
- Conditions (minimum d'achat, etc.)

**Étape 3 - Validité**
- Date de début
- Date de fin
- Jours de la semaine
- Créneaux horaires
- Quota (nombre max d'utilisation)

**Étape 4 - Médias**
- Image principale
- Images additionnelles
- Vidéo (optionnel)

**Étape 5 - Prévisualisation**
- Aperçu mobile
- Publication

---

### 4. Gestion des Établissements

#### Liste des établissements
- Cartes ou tableau
- Ajout nouvel établissement

#### Fiche établissement
- Informations générales
  - Nom
  - Adresse
  - Téléphone
  - Email
  - Site web
- Catégorie (bar, restaurant, loisirs, etc.)
- Horaires d'ouverture
- Photos
- Description
- Position GPS (carte)

---

### 5. Analytics

**Métriques principales** :
- Fréquentation journalière (365 derniers jours + prévisions)
- Vues par offre
- Réservations par période
- Taux de conversion
- Check-ins réalisés

**Graphiques** :
- Ligne : évolution temporelle sur 365 jours
- Calendrier heatmap : fréquentation par jour
- Barres : comparaison offres
- Prévisions : jours à venir (basé sur historique)

**Filtres** :
- Par établissement
- Par offre
- Par période

**Export** :
- CSV
- PDF (rapport)

> Note : Pas de validation requise, consultation libre par le partenaire

---

### 6. Réservations/Check-ins

Liste des réservations :
- Utilisateur
- Offre
- Date/heure
- Statut (réservé, check-in, annulé, no-show)

Actions :
- Valider check-in manuel
- Annuler réservation
- Contacter utilisateur

---

### 7. Paramètres

#### Profil entreprise
- Informations légales
- Logo
- Description

#### Facturation
- Plan actuel
- Historique factures
- Moyens de paiement

#### Équipe
- Gestion utilisateurs
- Rôles et permissions
- Invitations

#### Notifications
- Préférences email
- Alertes

---

## 🎨 Design System

### Palette de couleurs
```css
--primary: #[à définir depuis Figma]
--primary-foreground: #ffffff
--secondary: #[à définir]
--accent: #[à définir]
--background: #ffffff
--foreground: #1a1a1a
--muted: #f5f5f5
--border: #e5e5e5
--error: #ef4444
--success: #22c55e
--warning: #f59e0b
```

### Typographie
- **Titres** : Inter (ou font Figma)
- **Corps** : Inter
- **Monospace** : JetBrains Mono

---

## 🔐 Sécurité

### Authentification
- JWT avec refresh token
- Session timeout configurable
- Blocage après X tentatives
- **2FA obligatoire** pour tous les partenaires (TOTP)
- **Social Login** : Google, Apple, Facebook (autant que possible)

### Autorisations
Rôles :
- **Admin** : Accès complet
- **Manager** : Gestion offres et stats
- **Viewer** : Consultation uniquement

### Validation
- SIRET vérifié via API INSEE
- Email vérifié
- Modération offres (optionnel)

### Multi-Établissements
- Un partenaire peut gérer plusieurs établissements
- Chaque offre est liée à un établissement spécifique

---

## 📊 Intégrations

- **API INSEE** : Vérification SIRET
- **Google Maps** : Géolocalisation établissements
- **S3 + CloudFront** : Upload médias

> Note : Pas de facturation Stripe pour les partenaires (modèle gratuit ou futur changement)

---

## 🧪 Tests Requis

### Tests Unitaires
- Hooks personnalisés
- Fonctions utilitaires
- Validators

### Tests d'Intégration
- Formulaires complets
- Flows navigation

### Tests E2E
- Inscription partenaire
- Création offre complète
- Consultation analytics

---

## 📋 Checklist Qualité

- [ ] Responsive (desktop-first, support tablet)
- [ ] Accessibilité WCAG 2.1 AA
- [ ] Internationalisation (FR/EN)
- [ ] SEO (meta, sitemap)
- [ ] Performance (Lighthouse > 90)
- [ ] Sécurité (OWASP)
- [ ] Documentation composants

---

## 🔗 Références

- [Questions à clarifier](./QUESTIONS.md)
- [Modèle de données](../DATA_MODEL.md)
- [Backend API](../backend/PROMPT.md)
