# 🚀 Prompts de Génération - Yousoon Platform

> **Guide** : Copiez-collez ces prompts un par un pour générer la plateforme.  
> **Important** : Chaque prompt met à jour le fichier `GENERATION_STATUS.md`  
> **Reprise** : Si une génération échoue, relancez le même prompt, il reprendra là où il s'est arrêté.  
> **Parallélisation** : Certains prompts peuvent être lancés en parallèle (voir section dédiée)

---

## 📋 Instructions d'Utilisation

1. **Avant de commencer** : Vérifiez `GENERATION_STATUS.md` pour voir l'état actuel
2. **Pour chaque prompt** : Copiez le prompt complet et collez-le
3. **En cas d'erreur** : Relancez le même prompt, il vérifiera le statut et reprendra
4. **Validation** : Après chaque phase, vérifiez que les tests passent
5. **Parallélisation** : Utilisez plusieurs fenêtres/sessions pour les prompts parallèles

---

## 🔀 Guide de Parallélisation

### Quoi lancer en parallèle ?

```
TEMPS 0 ──────────────────────────────────────────────────────────────────────►

SESSION 1 (Backend)              SESSION 2 (Vitrine - Indépendant)
═══════════════════              ════════════════════════════════
                                 
[Prompt 1.1: Shared Domain]      [Prompt 7: Site Vitrine]
         │                                │
         ▼                                │
[Prompts 1.2+1.3+1.4 en //]              │
         │                                │
         ▼                                │
[Prompt 1.5: Federation]                  │
         │                                │
         ▼                                │
[Prompts 1.6+1.7 en //]                  │
         │                                │
         ▼                                │
[Prompts 2.1+2.2+2.3 en //]              │
         │                                │
         ▼                                │
[Prompts 3.1+3.2+3.3 en //]              ▼
         │                       ✅ Vitrine terminée
         ▼
   Backend prêt
         │
         ├──────────────────┬──────────────────┐
         ▼                  ▼                  ▼
SESSION 3 (Mobile)   SESSION 4 (Partner)  SESSION 5 (Admin)
═══════════════════  ═══════════════════  ═══════════════════
[Prompt 4.1]         [Prompt 5]           [Prompt 6]
[Prompt 4.2]              │                    │
[Prompt 4.3]              │                    │
[Prompt 4.4]              ▼                    ▼
     │               ✅ Terminé           ✅ Terminé
     ▼
✅ Terminé
```

### Règles de Parallélisation

| Groupe | Prompts Parallélisables | Condition |
|--------|-------------------------|-----------|
| **Groupe A** | 1.2, 1.3, 1.4 | Après 1.1 terminé |
| **Groupe B** | 1.6, 1.7 | Après 1.5 terminé |
| **Groupe C** | 2.1, 2.2, 2.3 | Après Phase 1 terminée |
| **Groupe D** | 3.1, 3.2, 3.3 | Après Phase 2 terminée |
| **Groupe E** | 4.x, 5, 6 | Après Phase 3 terminée |
| **Indépendant** | 7 (Vitrine) | Aucune dépendance |

---

## 🏗️ PHASE 1 : Backend Infrastructure

### Prompt 1.1 : Shared Domain (BLOQUANT)
```
Génère le Backend Phase 1 - Étape 1.1 : Package Shared Domain

Avant de commencer :
1. Lis le fichier docs/prompts/GENERATION_STATUS.md pour vérifier l'état actuel de l'étape 1.1
2. Si des fichiers sont déjà marqués ✅, ne les régénère pas
3. Génère uniquement les fichiers marqués ⬜ ou ❌

Fichiers à générer (si non complétés) :
- services/shared/domain/aggregate.go : AggregateRoot base avec gestion des domain events
- services/shared/domain/entity.go : Entity base avec ID
- services/shared/domain/valueobject.go : Value Objects communs (Email, Phone, Money, GeoLocation, Address)
- services/shared/domain/event.go : DomainEvent interface
- services/shared/domain/errors.go : Erreurs domain communes
- services/shared/domain/id.go : Types ID fortement typés (UserID, PartnerID, OfferID, etc.)

Référence : docs/prompts/backend/GENERATION_PLAN.md et docs/prompts/backend/ARCHITECTURE.md

Après génération de chaque fichier :
- Mets à jour GENERATION_STATUS.md avec le statut ✅ et la date
- Si erreur, marque ❌ avec la note d'erreur

À la fin, mets à jour le statut global de l'étape 1.1
```

### Prompts 1.2 + 1.3 + 1.4 : Infrastructure (PARALLÉLISABLES après 1.1)

> **⚡ Ces 3 prompts peuvent être lancés en parallèle dans des sessions différentes**

#### Prompt 1.2 : Infrastructure MongoDB
```
Génère le Backend Phase 1 - Étape 1.2 : Infrastructure MongoDB

Avant de commencer :
1. Vérifie que l'étape 1.1 est ✅ COMPLETED dans GENERATION_STATUS.md
2. Lis le statut de l'étape 1.2 pour voir ce qui reste à faire
3. Ne régénère pas les fichiers déjà ✅

Fichiers à générer (si non complétés) :
- services/shared/infrastructure/mongodb/client.go : Client MongoDB avec connection pooling
- services/shared/infrastructure/mongodb/repository.go : Repository base générique
- services/shared/infrastructure/mongodb/transaction.go : Support transactions multi-documents
- services/shared/infrastructure/mongodb/mapper.go : Interface de mapping domain <-> mongo

Dépendances : Utilise les types de services/shared/domain/

Après génération, mets à jour GENERATION_STATUS.md
```

### Prompt 1.3 : Infrastructure Redis
```
Génère le Backend Phase 1 - Étape 1.3 : Infrastructure Redis

Prérequis : Étapes 1.1 et 1.2 doivent être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Fichiers à générer (si non complétés) :
- services/shared/infrastructure/redis/client.go : Client Redis avec connection pooling
- services/shared/infrastructure/redis/cache.go : Cache générique avec TTL
- services/shared/infrastructure/redis/distributed_lock.go : Locks distribués pour concurrence

Après génération, mets à jour GENERATION_STATUS.md
```

### Prompt 1.4 : Infrastructure NATS
```
Génère le Backend Phase 1 - Étape 1.4 : Infrastructure NATS JetStream

Prérequis : Étapes 1.1-1.3 doivent être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Fichiers à générer (si non complétés) :
- services/shared/infrastructure/nats/client.go : Client NATS JetStream
- services/shared/infrastructure/nats/publisher.go : Event Publisher avec retry
- services/shared/infrastructure/nats/subscriber.go : Event Subscriber avec consumer groups
- services/shared/infrastructure/nats/serializer.go : JSON serialization des events

Utilise les types DomainEvent de services/shared/domain/event.go

Après génération, mets à jour GENERATION_STATUS.md
```

### Prompt 1.5 : GraphQL Federation Shared
```
Génère le Backend Phase 1 - Étape 1.5 : GraphQL Federation Shared

Prérequis : Étapes 1.1-1.4 doivent être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Fichiers à générer (si non complétés) :

Registry Client :
- services/shared/federation/registry/client.go : Client pour s'enregistrer au registry
- services/shared/federation/registry/discovery.go : Service discovery (watch changes)
- services/shared/federation/registry/health.go : Health check pour subgraphs

Directives :
- services/shared/federation/directives/auth.go : @auth(requires: ADMIN) directive
- services/shared/federation/directives/validation.go : @constraint(min: 1, max: 100)
- services/shared/federation/directives/deprecated.go : @deprecated directive custom

Scalars :
- services/shared/federation/scalars/datetime.go : DateTime scalar (ISO 8601)
- services/shared/federation/scalars/money.go : Money scalar (centimes)
- services/shared/federation/scalars/geolocation.go : GeoLocation scalar
- services/shared/federation/scalars/objectid.go : MongoDB ObjectID scalar

Middleware :
- services/shared/federation/middleware/context.go : Context enrichment (user, claims)
- services/shared/federation/middleware/dataloader.go : DataLoader factory pour batching

Référence : Apollo Federation 2 + gqlgen

Après génération, mets à jour GENERATION_STATUS.md
```

### Prompt 1.6 : Apollo Router
```
Génère le Backend Phase 1 - Étape 1.6 : Apollo Router Configuration

Prérequis : Étapes 1.1-1.5 doivent être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Fichiers à générer (si non complétés) :
- services/router/config/router.yaml : Configuration Apollo Router (supergraph, plugins, telemetry)
- services/router/plugins/auth.rhai : Plugin auth custom (JWT validation)
- services/router/plugins/ratelimit.rhai : Rate limiting plugin
- services/router/plugins/logging.rhai : Custom logging
- services/router/scripts/compose.sh : Script de composition des subgraphs
- services/router/scripts/watch.sh : Watch mode pour dev
- services/router/Dockerfile : Docker image Apollo Router

Le router doit :
- Poll le Registry pour les schemas
- Composer le supergraph automatiquement
- Propager les headers Authorization et X-Request-ID
- Exporter les traces vers Jaeger

Après génération, mets à jour GENERATION_STATUS.md
```

### Prompt 1.7 : Schema Registry
```
Génère le Backend Phase 1 - Étape 1.7 : Schema Registry Service

Prérequis : Étapes 1.1-1.6 doivent être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Fichiers à générer (si non complétés) :

Main :
- services/registry/cmd/main.go : Entrypoint
- services/registry/config/config.go : Configuration

Storage :
- services/registry/internal/storage/store.go : Interface storage
- services/registry/internal/storage/memory.go : In-memory (dev)
- services/registry/internal/storage/redis.go : Redis (prod)

Composer :
- services/registry/internal/composer/composer.go : Composition du supergraph
- services/registry/internal/composer/validator.go : Validation des schemas

Discovery :
- services/registry/internal/discovery/watcher.go : Watch Kubernetes services
- services/registry/internal/discovery/k8s.go : Kubernetes service discovery

API :
- services/registry/internal/api/handler.go : REST API (POST/DELETE /subgraphs, GET /supergraph)
- services/registry/internal/api/graphql.go : GraphQL API pour introspection

- services/registry/Dockerfile : Docker image

Le registry doit :
- Stocker les schemas SDL des subgraphs
- Composer automatiquement le supergraph quand un subgraph change
- Exposer une API REST pour registration/deregistration
- Watch les services Kubernetes avec label graphql.federation/subgraph=true

Après génération, mets à jour GENERATION_STATUS.md

Une fois terminé, mets à jour le statut global de la Phase 1 dans GENERATION_STATUS.md
```

---

## 🔷 PHASE 2 : Core Subgraphs

### Prompt 2.1 : Identity Service
```
Génère le Backend Phase 2 - Étape 2.1 : Identity Service (Subgraph)

Prérequis : Phase 1 complète (✅)
Vérifie GENERATION_STATUS.md avant de commencer.

Génère le service Identity complet en suivant docs/prompts/backend/GENERATION_PLAN.md

Structure à générer :
services/identity-service/
├── cmd/main.go                 # Entrypoint avec auto-registration au Registry
├── config/config.go
├── gqlgen.yml                  # Config gqlgen avec Federation 2
├── graph/
│   ├── schema.graphqls         # Schema avec @key(fields: "id")
│   ├── federation.graphqls
│   ├── resolver.go
│   ├── schema.resolvers.go
│   └── entity.resolvers.go
├── internal/
│   ├── domain/aggregate/user.go
│   ├── domain/entity/subscription.go, identity_verification.go
│   ├── domain/valueobject/profile.go, preferences.go, grade.go
│   ├── domain/event/user_events.go
│   ├── domain/repository/user_repository.go
│   ├── application/command/register_user.go, login_user.go, verify_identity.go, subscribe.go
│   ├── application/query/get_user.go, get_subscription.go
│   └── infrastructure/...
├── proto/identity.proto
└── Dockerfile

Le service doit :
- S'enregistrer automatiquement au Schema Registry au démarrage
- Exposer un schema GraphQL avec User @key(fields: "id")
- Gérer JWT (access 6h, refresh 30j stocké Redis)
- Gérer la vérification CNI via OCR interne
- Gérer les abonnements via In-App Purchase

Mets à jour GENERATION_STATUS.md après chaque composant généré.
```

### Prompt 2.2 : Partner Service
```
Génère le Backend Phase 2 - Étape 2.2 : Partner Service (Subgraph)

Prérequis : Étape 2.1 doit être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère le service Partner complet en suivant docs/prompts/backend/GENERATION_PLAN.md

Structure similaire à Identity Service avec :
- Partner @key(fields: "id")
- Establishment @key(fields: "id")
- Extension de User pour ajouter partners: [Partner!]!
- Gestion des équipes (invitations, rôles)
- 2FA obligatoire pour les partenaires

Mets à jour GENERATION_STATUS.md après génération.
```

### Prompt 2.3 : Discovery Service
```
Génère le Backend Phase 2 - Étape 2.3 : Discovery Service (Subgraph)

Prérequis : Étapes 2.1-2.2 doivent être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère le service Discovery complet avec :
- Offer @key(fields: "id")
- Category @key(fields: "id")
- Extension de Establishment et Partner pour ajouter offers
- Recherche Elasticsearch
- Queries géospatiales (getNearbyOffers)
- Système de recommandations

Mets à jour GENERATION_STATUS.md après génération.
Mets à jour le statut global Phase 2 une fois terminé.
```

---

## 🔷 PHASE 3 : Business Subgraphs

### Prompt 3.1 : Booking Service
```
Génère le Backend Phase 3 - Étape 3.1 : Booking Service (Subgraph)

Prérequis : Phase 2 complète (✅)
Vérifie GENERATION_STATUS.md avant de commencer.

Génère le service Booking complet avec :
- Outing @key(fields: "id") (réservation)
- OfferSnapshot (copie immutable au moment de la réservation)
- QRCode pour check-in
- Extension de User et Offer pour ajouter outings
- Subscriptions GraphQL pour outingStatusChanged
- Expiration automatique après 30min

Mets à jour GENERATION_STATUS.md après génération.
```

### Prompt 3.2 : Engagement Service
```
Génère le Backend Phase 3 - Étape 3.2 : Engagement Service (Subgraph)

Prérequis : Étape 3.1 doit être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère le service Engagement complet avec :
- Review @key(fields: "id")
- Favorite
- Conversation @key(fields: "id") avec Messages
- Extensions de User et Offer pour favorites, reviews
- Note : Les partenaires ne peuvent PAS répondre aux avis

Mets à jour GENERATION_STATUS.md après génération.
```

### Prompt 3.3 : Notification Service
```
Génère le Backend Phase 3 - Étape 3.3 : Notification Service (Subgraph)

Prérequis : Étape 3.2 doit être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère le service Notification complet avec :
- Notification @key(fields: "id")
- Extension de User pour notifications
- Intégration OneSignal (push)
- Intégration AWS SNS (email, SMS)
- Event handlers pour les events des autres services
- Subscriptions GraphQL pour newNotification

Types activés : offer_nearby, booking_reminder, marketing
Types désactivés : offer_expiring, new_partner

Mets à jour GENERATION_STATUS.md après génération.
Mets à jour le statut global Phase 3 une fois terminé.
```

---

## 📱 PHASE 4 : App Mobile Flutter

### Prompt 4.1 : Core & Design System
```
Génère l'App Mobile Phase 4 - Étape 4.1 : Core & Design System

Prérequis : Backend Phases 1-2 minimum (✅)
Vérifie GENERATION_STATUS.md avant de commencer.
Référence : docs/prompts/app-mobile/GENERATION_PLAN.md et docs/prompts/DESIGN_SYSTEM.md

Génère :
- Structure projet Flutter avec Clean Architecture
- Theme (Dark Mode natif) avec couleurs Figma :
  - Dark Black #000000, Indian Gold #E99B27, Flash White #FFFFFF
  - Grey Jet #6D6D6D, Mantis Green #5FC15C, Persian Red #CC2936
- Typography (Futura, Poppins)
- Spacings
- Shared Widgets (YsButton, YsTextField, YsCard, etc.)
- GraphQL Client (ferry) avec cache
- Riverpod setup

Mets à jour GENERATION_STATUS.md après génération.
```

### Prompt 4.2 : Features Auth
```
Génère l'App Mobile Phase 4 - Étape 4.2 : Features Auth

Prérequis : Étape 4.1 doit être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère les features d'authentification :
- Splash Screen
- Onboarding (slides)
- Login (email + password + biométrie optionnelle)
- Register
- Forgot Password
- OTP Verification
- Identity Verification (upload CNI + OCR)

Respecte scrupuleusement le design Figma.
Intègre avec Identity Service via GraphQL.

Mets à jour GENERATION_STATUS.md après génération.
```

### Prompt 4.3 : Features Core
```
Génère l'App Mobile Phase 4 - Étape 4.3 : Features Core

Prérequis : Étape 4.2 doit être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère les features principales :
- Home/Feed ("Pour vous" - page par défaut)
- Offers (liste, détail, recherche)
- Booking (réservation, QR code, check-in)
- Map (Google Maps avec pins par catégorie)
- Profile (informations, grade, historique)
- Tab Bar (5 entrées : Mes events, Favoris, Pour vous, Carte, Messages)

Respecte scrupuleusement le design Figma.

Mets à jour GENERATION_STATUS.md après génération.
```

### Prompt 4.4 : Features Social
```
Génère l'App Mobile Phase 4 - Étape 4.4 : Features Social

Prérequis : Étape 4.3 doit être ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère les features sociales :
- Favorites (ajout/suppression, liste)
- Reviews (notation étoiles, commentaire, photos)
- Messaging (conversations, temps réel via Subscriptions)
- Notifications (liste, préférences)
- Partage (deep links)

Mets à jour GENERATION_STATUS.md après génération.
Mets à jour le statut global Phase 4 une fois terminé.
```

---

## 💼 PHASE 5 : Site Partenaires

### Prompt 5 : Site Partenaires Complet
```
Génère le Site Partenaires - Phase 5

Prérequis : Backend complet (Phases 1-3 ✅)
Vérifie GENERATION_STATUS.md avant de commencer.
Référence : docs/prompts/site-partenaires/GENERATION_PLAN.md

Génère le portail partenaires (business.yousoon.com) :
- Setup projet (React 19, TypeScript, Vite, TailwindCSS, shadcn/ui)
- Auth (login, register, 2FA obligatoire, social login)
- Layout (sidebar, header)
- Dashboard (KPIs, graphiques)
- Gestion Offres (CRUD, multi-étapes, médias)
- Gestion Établissements
- Analytics (365 jours + prévisions, export CSV/PDF)
- Réservations/Check-ins
- Settings (profil, équipe, notifications)

Stack : urql pour GraphQL, Zustand, TanStack Query, React Hook Form + Zod

Mets à jour GENERATION_STATUS.md après chaque section.
Mets à jour le statut global Phase 5 une fois terminé.
```

---

## 🔐 PHASE 6 : Admin Backoffice

### Prompt 6 : Admin Backoffice Complet
```
Génère l'Admin Backoffice - Phase 6

Prérequis : Backend complet (Phases 1-3 ✅)
Vérifie GENERATION_STATUS.md avant de commencer.
Référence : docs/prompts/admin/GENERATION_PLAN.md

Génère le backoffice admin (admin.yousoon.com - accès interne) :
- Setup projet (React 19, TypeScript, Vite)
- Auth admin avec rôles (super_admin, moderator, support)
- Dashboard (KPIs globaux)
- Gestion Users (liste, détail, suspendre)
- Gestion Partners (validation, blocage)
- Validation CNI (affichage images, valider/rejeter)
- Modération Avis (signalés, suppression)
- Gestion Abonnements (plans, historique)
- Configuration (catégories, paramètres)
- Audit Logs

Note : Pas d'Ingress public, accès via kubectl port-forward

Mets à jour GENERATION_STATUS.md après génération.
Mets à jour le statut global Phase 6 une fois terminé.
```

---

## 🌐 PHASE 7 : Site Vitrine

### Prompt 7 : Site Vitrine Complet
```
Génère le Site Vitrine - Phase 7

Cette phase peut être faite en parallèle des autres.
Vérifie GENERATION_STATUS.md avant de commencer.
Référence : docs/prompts/site-vitrine/GENERATION_PLAN.md

Génère le site vitrine (www.yousoon.com) :
- Setup Next.js 14 avec App Router
- Pages : Accueil, Fonctionnalités, Tarifs, FAQ, Contact
- Section Partenaires (CTA vers business.yousoon.com)
- Téléchargement App (liens App Store / Play Store)
- SEO optimisé
- i18n (FR, EN) avec next-intl
- Responsive (mobile-first)
- Dark mode (cohérent avec l'app)

Mets à jour GENERATION_STATUS.md après génération.
Mets à jour le statut global Phase 7 une fois terminé.
```

---

## 🚀 PHASE 8 : Déploiement & Tests

### Prompt 8.1 : Kubernetes Manifests
```
Génère le Déploiement Phase 8 - Étape 8.1 : Kubernetes Manifests

Prérequis : Toutes les phases précédentes ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère les manifests Kubernetes :
- Namespace, ConfigMaps, Secrets
- Deployments pour tous les services avec labels Federation
- Services (ClusterIP pour subgraphs, LoadBalancer pour Router)
- Ingress (Nginx) pour Router, Partner Portal, Vitrine
- Infrastructure : MongoDB, Redis, NATS, Elasticsearch
- Monitoring : Prometheus, Grafana, Jaeger, Loki

Structure : deploy/kubernetes/

Mets à jour GENERATION_STATUS.md après génération.
```

### Prompt 8.2 : CI/CD Pipelines
```
Génère le Déploiement Phase 8 - Étape 8.2 : CI/CD Pipelines

Prérequis : Étape 8.1 ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère les pipelines GitHub Actions :
- .github/workflows/backend.yml : Build, test, push images, deploy
- .github/workflows/mobile.yml : Build, test, deploy TestFlight/Play Store
- .github/workflows/partner-portal.yml
- .github/workflows/vitrine.yml
- .github/workflows/admin.yml

Inclure : Tests, linting, build Docker, push ECR, deploy EKS

Mets à jour GENERATION_STATUS.md après génération.
```

### Prompt 8.3 : Tests E2E
```
Génère le Déploiement Phase 8 - Étape 8.3 : Tests E2E

Prérequis : Étape 8.2 ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Génère les tests E2E :
- Backend : Tests GraphQL Federation cross-subgraph
- App Mobile : Parcours inscription complet (Patrol)
- Site Partenaires : Création offre (Playwright)
- Admin : Validation partenaire (Playwright)

Mets à jour GENERATION_STATUS.md après génération.
Mets à jour le statut global Phase 8 une fois terminé.

🎉 Si Phase 8 complète, la génération est terminée !
```

---

## ⚡ PROMPTS COMBINÉS (Parallélisation Maximale)

> Ces prompts regroupent plusieurs étapes parallélisables en un seul prompt.  
> Idéal si vous utilisez une seule session mais voulez profiter de la parallélisation.

### Prompt Combiné : Phase 1 Infrastructure Complète
```
Génère le Backend Phase 1 COMPLÈTE : Infrastructure

Vérifie GENERATION_STATUS.md et génère uniquement les étapes non complétées.

Ordre de génération :
1. ÉTAPE 1.1 (Shared Domain) - BLOQUANT, faire en premier
2. ÉTAPES 1.2, 1.3, 1.4 (MongoDB, Redis, NATS) - Parallélisables après 1.1
3. ÉTAPE 1.5 (Federation Shared) - Après 1.2-1.4
4. ÉTAPES 1.6, 1.7 (Router, Registry) - Parallélisables après 1.5

Pour chaque étape :
- Vérifie le statut dans GENERATION_STATUS.md
- Si ⬜ ou ❌, génère les fichiers
- Mets à jour le statut après chaque fichier

Référence : docs/prompts/backend/GENERATION_PLAN.md

À la fin, marque la Phase 1 comme ✅ COMPLETED dans GENERATION_STATUS.md
```

### Prompt Combiné : Phase 2 Core Subgraphs (Parallèle)
```
Génère le Backend Phase 2 COMPLÈTE : Core Subgraphs

Prérequis : Phase 1 ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Les 3 subgraphs sont indépendants et peuvent être générés en parallèle :
1. Identity Service (Subgraph) - User @key, Auth, Subscriptions
2. Partner Service (Subgraph) - Partner @key, Establishments, Teams
3. Discovery Service (Subgraph) - Offer @key, Categories, Search

Pour chaque service, génère la structure complète :
- cmd/main.go avec auto-registration
- graph/ (gqlgen federation)
- internal/domain/, application/, infrastructure/
- proto/, gqlgen.yml, Dockerfile

Mets à jour GENERATION_STATUS.md après chaque service.
Marque la Phase 2 comme ✅ à la fin.
```

### Prompt Combiné : Phase 3 Business Subgraphs (Parallèle)
```
Génère le Backend Phase 3 COMPLÈTE : Business Subgraphs

Prérequis : Phase 2 ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Les 3 subgraphs sont indépendants :
1. Booking Service - Outing @key, QR Code, Check-in
2. Engagement Service - Reviews, Favorites, Messaging
3. Notification Service - Push (OneSignal), Email/SMS (AWS SNS)

Génère la structure complète pour chaque service.
Mets à jour GENERATION_STATUS.md après chaque service.
Marque la Phase 3 comme ✅ à la fin.

🎉 Après cette phase, le Backend est prêt pour les frontends !
```

### Prompt Combiné : Frontends en Parallèle (Session 1/3 - Mobile)
```
Génère l'App Mobile Flutter - Phase 4 COMPLÈTE

Prérequis : Backend Phases 1-3 ✅ (au minimum Phase 2)
Vérifie GENERATION_STATUS.md avant de commencer.

Cette génération peut être lancée en PARALLÈLE avec les Phases 5 et 6.

Génère dans l'ordre :
1. Étape 4.1 : Core & Design System
2. Étape 4.2 : Features Auth
3. Étape 4.3 : Features Core
4. Étape 4.4 : Features Social

Référence : docs/prompts/app-mobile/GENERATION_PLAN.md et DESIGN_SYSTEM.md
Respecte scrupuleusement le design Figma.

Mets à jour GENERATION_STATUS.md après chaque étape.
Marque la Phase 4 comme ✅ à la fin.
```

### Prompt Combiné : Frontends en Parallèle (Session 2/3 - Partner)
```
Génère le Site Partenaires - Phase 5 COMPLÈTE

Prérequis : Backend Phases 1-3 ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Cette génération peut être lancée en PARALLÈLE avec les Phases 4 et 6.

Génère le portail business.yousoon.com complet :
- Setup React 19 + TypeScript + Vite + TailwindCSS + shadcn/ui
- Auth (login, register, 2FA, social)
- Dashboard, Offres, Établissements, Analytics, Settings

Référence : docs/prompts/site-partenaires/GENERATION_PLAN.md

Mets à jour GENERATION_STATUS.md progressivement.
Marque la Phase 5 comme ✅ à la fin.
```

### Prompt Combiné : Frontends en Parallèle (Session 3/3 - Admin)
```
Génère l'Admin Backoffice - Phase 6 COMPLÈTE

Prérequis : Backend Phases 1-3 ✅
Vérifie GENERATION_STATUS.md avant de commencer.

Cette génération peut être lancée en PARALLÈLE avec les Phases 4 et 5.

Génère le backoffice admin.yousoon.com complet :
- Setup React 19 + TypeScript + Vite
- Auth admin avec rôles
- Dashboard, Users, Partners, CNI Validation, Moderation, Config

Note : Accès interne uniquement (kubectl port-forward)

Référence : docs/prompts/admin/GENERATION_PLAN.md

Mets à jour GENERATION_STATUS.md progressivement.
Marque la Phase 6 comme ✅ à la fin.
```

---

## ✅ Validation Finale

```
Vérifie la génération complète de Yousoon.

1. Lis GENERATION_STATUS.md et vérifie que toutes les phases sont ✅
2. Liste les phases/étapes encore incomplètes
3. Pour chaque élément incomplet, indique le prompt à relancer
4. Si tout est ✅, génère un rapport final avec :
   - Nombre total de fichiers générés
   - Temps estimé total
   - Prochaines étapes recommandées (review, tests manuels, etc.)
```

---

## 🔄 Reprise après Erreur

Si une génération échoue :

```
Reprends la génération de [PHASE X - ÉTAPE Y].

1. Lis GENERATION_STATUS.md pour voir l'état actuel
2. Identifie les fichiers marqués ❌ ou ⬜
3. Régénère uniquement ces fichiers
4. Mets à jour les statuts au fur et à mesure
5. Si l'erreur persiste, note-la dans GENERATION_STATUS.md
```

---

## 📊 Récapitulatif Temps avec Parallélisation

| Scénario | Temps Total | Sessions Requises |
|----------|-------------|-------------------|
| **Séquentiel (1 session)** | ~167h | 1 |
| **Parallélisé (2 sessions)** | ~95h | 2 (Backend + Vitrine, puis Frontend) |
| **Parallélisé (3 sessions)** | ~85h | 3 (+ Mobile/Partner/Admin en //) |
| **Parallélisé (5 sessions)** | ~75h | 5 (Maximum de parallélisme) |

### Stratégie Optimale (3 sessions)

```
T=0h   Session A: Phase 1.1 Shared Domain
       Session B: Phase 7 Site Vitrine (indépendant)

T=2h   Session A: Phases 1.2+1.3+1.4 (parallèle interne)
       Session B: Continue Vitrine

T=4h   Session A: Phases 1.5, 1.6+1.7
       Session B: Vitrine terminée ✅

T=8h   Session A: Phase 2 (3 subgraphs //)

T=13h  Session A: Phase 3 (3 subgraphs //)

T=17h  Session A: Phase 4 Mobile
       Session B: Phase 5 Partner Portal
       Session C: Phase 6 Admin

T=63h  Tous frontends terminés ✅

T=63h  Phase 8 Déploiement

T=78h  🎉 TERMINÉ
```
