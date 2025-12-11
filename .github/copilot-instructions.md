# 🤖 Copilot Instructions - Yousoon Platform

> Contexte et décisions techniques pour le projet Yousoon  
> Dernière mise à jour : 9 décembre 2025

---

## ⚠️ RÈGLE IMPORTANTE

**Toujours mettre à jour cette documentation** lorsque des décisions techniques sont prises ou modifiées. Cette documentation est la source de vérité pour le projet.

---

## 📋 Résumé du Projet

**Yousoon** est une plateforme de sorties avec réductions qui met en relation clients et partenaires.

### Concept Business

1. **Apport de clients** : Yousoon apporte des clients aux partenaires (bars, restaurants, organismes de sorties) qui en échange offrent des réductions
2. **Abonnement utilisateurs** : Les clients paient un abonnement pour accéder aux sorties à prix réduit via les partenaires
3. **Intermédiaire** : Yousoon fait le lien entre clients (qui veulent sortir pas cher) et partenaires (qui veulent des clients)

### Architecture

- **App Mobile** : Flutter (iOS/Android)
- **Site Partenaires** : React TypeScript + Vite
- **Site Vitrine** : Next.js 14
- **Backend** : Go avec microservices DDD
- **Admin Backoffice** : React TypeScript (accès restreint)

### ⚠️ Règle d'Accès aux Données

**Tous les frontends (App Mobile, Site Partenaires, Admin Backoffice) communiquent UNIQUEMENT via l'API GraphQL.**

- ❌ **Jamais** d'accès direct à MongoDB depuis les frontends
- ✅ Toutes les données passent par l'API GraphQL (`api.yousoon.com`)
- ✅ Apollo Router fédère les requêtes vers les microservices

```
┌─────────────────────────────────────────────────────────────┐
│          FRONTENDS (Mobile, Partenaires, Admin)             │
└─────────────────────────────┬───────────────────────────────┘
                              │ GraphQL (HTTPS)
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    API GraphQL (Apollo Router)               │
│                      api.yousoon.com                         │
└─────────────────────────────┬───────────────────────────────┘
                              │ Federation
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Microservices (Go + gqlgen)                  │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    MongoDB / Redis / NATS                    │
└─────────────────────────────────────────────────────────────┘
```

### 🌐 URLs des Sites

| Site | URL | Description |
|------|-----|-------------|
| **Site Vitrine** | [www.yousoon.com](https://www.yousoon.com) | Landing page publique |
| **Portail Partenaires** | [business.yousoon.com](https://business.yousoon.com) | Gestion partenaires |
| **API GraphQL** | [api.yousoon.com](https://api.yousoon.com) | Apollo Router Federation |
| **Admin Backoffice** | `kubectl port-forward` | Accès interne uniquement |

---

## 🏛️ Architecture Backend

### Domain-Driven Design (DDD)

**6 Bounded Contexts** :
| Context | Type | Responsabilité |
|---------|------|----------------|
| **Identity** | Core | Auth, profils, vérification CNI, abonnements |
| **Partner** | Core | Partenaires, établissements, équipes |
| **Discovery** | Core | Catalogue offres, recherche, recommandations |
| **Booking** | Core | Réservations, check-in QR, historique |
| **Engagement** | Supporting | Favoris, avis |
| **Notification** | Generic | Push, email, SMS |

### Architecture Hexagonale

```
┌─────────────────────────────────────┐
│      INTERFACE (gRPC/HTTP)          │  ← Driving Adapters
├─────────────────────────────────────┤
│      APPLICATION (Commands/Queries)  │  ← Use Cases
├─────────────────────────────────────┤
│      DOMAIN (Aggregates/Events)      │  ← Business Rules (PURE)
├─────────────────────────────────────┤
│      INFRASTRUCTURE (MongoDB/NATS)   │  ← Driven Adapters
└─────────────────────────────────────┘
```

### Communication Inter-Services

| Type | Technologie | Usage |
|------|-------------|-------|
| **Synchrone** | **gRPC** | Appels requête/réponse entre services |
| **Asynchrone** | **NATS JetStream** | Domain Events (pub/sub) |

**Pourquoi NATS JetStream ?**
- Latence < 1ms (requis < 50ms)
- Ultra-léger (~15MB)
- Kubernetes-native
- Simplicité d'opération
- Persistance des events

---

## 🛠️ Stack Technique

### Backend (Go)
```yaml
Langage: Go 1.21+
GraphQL: gqlgen + Apollo Federation 2
API Gateway: Apollo Router (Federation)
Service Discovery: Schema Registry custom + Kubernetes labels
gRPC: google.golang.org/grpc (inter-service)
Events: NATS JetStream
Database: MongoDB (Europe/Irlande - RGPD)
Cache: Redis
ORM: go.mongodb.org/mongo-driver
```

### App Mobile (Flutter)
```yaml
Framework: Flutter 3.x
State: Riverpod 2.x
GraphQL: ferry
Cache local: Hive/Isar
Paiements: in_app_purchase (100% Apple/Google Pay)
Auth biométrique: local_auth
QR Scanner: mobile_scanner
Theme: Dark Mode natif
```

### Sites Web (React/Next.js)
```yaml
Partner Portal: React 18 + TypeScript + Vite
Vitrine: Next.js 14 + next-intl + MySQL
UI: TailwindCSS + shadcn/ui
GraphQL: urql (Partner Portal uniquement)
```

> **Note** : Le Site Vitrine utilise **MySQL** pour son propre contenu (blog, FAQ, pages). Il n'utilise PAS l'API GraphQL backend.

### Infrastructure
```yaml
Cloud: AWS (EKS)
Région: Europe (Irlande) - RGPD
CI/CD: GitHub Actions
IaC: Helm + Helmfile
Storage: AWS S3 + CloudFront
Search: Elasticsearch
Observability: OpenTelemetry + Jaeger + Prometheus + Loki + Grafana
Crash Reporting: Sentry (self-hosted)
Analytics: Amplitude
Notifications: OneSignal (Push) + AWS SNS (Email/SMS)
```

### Déploiement Infrastructure (Helmfile)

L'infrastructure est déployée via **Helmfile** :

```
deploy/helm/
├── helmfile.yaml                    # Orchestration principale
├── secrets-README.md                # Instructions secrets
└── values/
    ├── mongodb.yaml
    ├── redis.yaml
    ├── nats.yaml
    ├── elasticsearch.yaml
    ├── prometheus-stack.yaml
    ├── loki.yaml
    └── jaeger.yaml
```

**Composants déployés** :
| Composant | Chart Helm | Usage |
|-----------|-----------|-------|
| MongoDB | bitnami/mongodb | Base de données principale |
| Redis | bitnami/redis | Cache et sessions |
| NATS | nats/nats | Messaging (events) |
| Elasticsearch | elastic/elasticsearch | Recherche full-text |
| Prometheus + Grafana | prometheus-community/kube-prometheus-stack | Monitoring |
| Loki | grafana/loki-stack | Agrégation de logs |
| Jaeger | jaegertracing/jaeger | Tracing distribué |

**Workflow CI/CD** : `.github/workflows/helmfile-deploy.yml`
- Branche `staging` → Namespace `yousoon-staging` → Mode `sidecar` (4 pods)
- Branche `prod` → Namespace `yousoon-prod` → Mode `classic` (~18 pods)
- Déploiement automatique sur push dans `deploy/helm/`

---

## 📊 Modèle de Données

### Collections MongoDB

| Collection | Description |
|------------|-------------|
| `users` | Utilisateurs (Yousooners) |
| `partners` | Partenaires/Fournisseurs |
| `establishments` | Établissements physiques |
| `offers` | Offres/Réductions |
| `bookings` | Réservations (Outings) |
| `categories` | Catégories d'offres |
| `subscriptions` | Abonnements utilisateurs |
| `reviews` | Avis et notes |

### Aggregates DDD

```
Identity:    User (+ IdentityVerification, Subscription)
Partner:     Partner (+ Establishment, TeamMember)
Discovery:   Offer, Category
Booking:     Outing (+ QRCode, OfferSnapshot)
Engagement:  Favorite, Review
```

---

## 🎨 Design System (Figma)

**Fichier Figma** : `1GXJECHtsYzq46OYbSHiaj`

### Palette de Couleurs
| Nom | Hex | Usage |
|-----|-----|-------|
| Dark Black | `#000000` | Background principal |
| Indian Gold | `#E99B27` | Accent, CTAs |
| Flash White | `#FFFFFF` | Texte sur fond noir |
| Grey Jet | `#6D6D6D` | Éléments inactifs |
| Eerie Black | `#CCCCCC` | Texte secondaire |
| Mantis Green | `#5FC15C` | Validation |
| Persian Red | `#CC2936` | Erreurs |

### Typographie
- **Titres** : Futura Bold/Medium
- **Corps** : Futura Medium (14-16pt)

### Navigation
- Page par défaut : **"Pour vous"**
- 5 entrées Tap Bar + 2 en haut à droite

---

## ✅ Décisions Techniques Validées

| Sujet | Décision |
|-------|----------|
| **Cloud** | AWS EKS (Kubernetes) - Région Irlande |
| **Architecture** | Microservices DDD (ou monolithe modulaire si trop complexe) |
| **MongoDB** | 1 cluster avec 1 database par context (self-hosted EKS) |
| **MongoDB HA** | Non pour commencer (Standalone) |
| **JWT** | Identity génère, Gateway valide (Access: 6h, Refresh: 30j) |
| **Refresh Token** | Stocké dans Redis |
| **Paiements** | 100% in-app (Apple Pay / Google Pay) |
| **2FA** | Obligatoire Admin + Partenaires uniquement |
| **Check-in** | QR Code uniquement (pas de geofencing) |
| **Comptes** | Unifiés (user peut être partenaire) |
| **Réponses avis** | Les partenaires ne peuvent PAS répondre |
| **RGPD** | Suppression sous 30 jours (grace period) |
| **Biométrie** | Pour reconnexion utilisateur (optionnel) |
| **Theme** | Dark Mode natif |
| **Vérification CNI** | OCR interne - tous documents - 10 tentatives max |
| **Notifications** | OneSignal (push) + AWS SNS (Email, SMS) |
| **Recherche** | Elasticsearch |
| **Stockage média** | AWS S3 + CloudFront |
| **Analytics** | Amplitude |
| **Cartographie** | Google Maps |
| **Observabilité** | OpenTelemetry + Jaeger + Prometheus + Loki + Grafana |
| **Crash Reporting** | Sentry (self-hosted) |
| **GraphQL Subscriptions** | Oui (temps réel WebSocket) |
| **Persisted Queries** | Oui |
| **Ingress** | Nginx Ingress |
| **Secrets** | Kubernetes Secrets |
| **DNS** | Route53 |
| **SSL** | Let's Encrypt (cert-manager) |
| **Rate Limiting** | Par user, détection abus réservations |
| **Géo-restriction** | Aucune (monde entier) |
| **Langues** | FR + EN, traduction automatique |
| **Mode Offline** | Oui (favoris, historique) |

### App Mobile
| Sujet | Décision |
|-------|----------|
| **iOS minimum** | Dernière version (iOS 17+) |
| **Android minimum** | Dernière version (API 34+) |
| **Bundle ID** | com.yousoon.yousoon |
| **Catégorie stores** | Lifestyle |
| **CI/CD** | GitHub Actions |
| **Beta iOS** | TestFlight |
| **Beta Android** | Google Play Internal Testing |

### Sites Web (React)
| Sujet | Décision |
|-------|----------|
| **React** | 19.x (dernière version) |
| **TypeScript** | 5.x |
| **Build** | Vite 5.x |

### Performance
| Sujet | Décision |
|-------|----------|
| **Objectif** | 5000 utilisateurs/heure minimum |

---

## 📁 Structure des Fichiers

```
yousoon-v2/
├── .github/
│   ├── copilot-instructions.md     # CE FICHIER
│   └── workflows/
│       └── helmfile-deploy.yml     # CI/CD Infrastructure
├── docs/
│   └── prompts/
│       ├── DATA_MODEL.md           # Schémas MongoDB
│       ├── DESIGN_SYSTEM.md        # Extrait Figma
│       ├── app-mobile/
│       │   └── PROMPT.md           # Specs Flutter
│       ├── site-partenaires/
│       │   └── PROMPT.md           # Specs React
│       ├── site-vitrine/
│       │   └── PROMPT.md           # Specs Next.js
│       ├── admin/
│       │   └── PROMPT.md           # Specs Admin
│       └── backend/
│           └── ARCHITECTURE.md     # Architecture DDD détaillée
├── deploy/
│   └── helm/                       # Helmfile + values
│       ├── helmfile.yaml
│       ├── secrets-README.md
│       └── values/
└── apps/
    ├── mobile/                     # Flutter App
    ├── partners/                   # React Partner Site (business.yousoon.com)
    ├── siteweb/                    # Next.js Landing (www.yousoon.com)
    ├── admin/                      # React Admin (accès interne)
    └── services/                   # Backend Microservices
        ├── router/                 # Apollo Router (Federation Gateway)
        ├── registry/               # Schema Registry (Service Discovery)
        ├── shared/                 # Shared Go modules
        ├── identity-service/       # Auth, Users, Subscriptions (Subgraph)
        ├── partner-service/        # Partners, Establishments (Subgraph)
        ├── discovery-service/      # Offers, Search (Subgraph)
        ├── booking-service/        # Outings, Check-in (Subgraph)
        ├── engagement-service/     # Favorites, Reviews (Subgraph)
        └── notification-service/   # Push, Email, SMS (Subgraph)
```

---

## 🔧 Conventions de Code

### Go (Backend)
```go
// Package naming: lowercase, single word
package booking

// Interface naming: verb + "er"
type OutingRepository interface {}

// Aggregate methods: verb
func (o *Outing) CheckIn(qr string) error {}

// Value Objects: immutable, no setters
type Email struct { value string }
```

### Flutter (Mobile)
```dart
// Feature-first structure
// lib/features/{feature}/

// Riverpod providers
final userProvider = StateNotifierProvider<UserNotifier, UserState>

// Repository pattern
abstract class OfferRepository {}
```

### React (Web)
```typescript
// Component naming: PascalCase
export function OfferCard({ offer }: Props) {}

// Hooks: use prefix
export function useOffers() {}

// Types: suffix with Type or interface
interface OfferType {}
```

---

## ⚠️ Questions En Suspens

1. **Abonnements** : Détails des plans (noms, prix, limites)
2. **Rayon recherche** : Valeur par défaut (actuellement 10km)
3. **Catégories** : Liste définitive des catégories d'intérêts

---

## 📞 Contexte MCP

Les MCPs disponibles dans ce projet :
- **Figma MCP** : Analyse des designs (`mcp_figma_*`)
- **GitKraken MCP** : Git, PRs, issues (`mcp_gitkraken_*`)
- **Container MCP** : Docker (`mcp_copilot_conta_*`)

Les MCPs s'exécutent **localement** et communiquent avec les APIs cloud respectives.

---

## 🔗 Références Détaillées

Pour plus de détails, voir :
- **Architecture DDD** : [docs/prompts/backend/ARCHITECTURE.md](../docs/prompts/backend/ARCHITECTURE.md)
- **Modèle de données** : [docs/prompts/DATA_MODEL.md](../docs/prompts/DATA_MODEL.md)
- **Design System** : [docs/prompts/DESIGN_SYSTEM.md](../docs/prompts/DESIGN_SYSTEM.md)
- **Specs Flutter** : [docs/prompts/app-mobile/PROMPT.md](../docs/prompts/app-mobile/PROMPT.md)
- **Specs Partenaires** : [docs/prompts/site-partenaires/PROMPT.md](../docs/prompts/site-partenaires/PROMPT.md)
- **Specs Admin** : [docs/prompts/admin/PROMPT.md](../docs/prompts/admin/PROMPT.md)
- **Specs Site Vitrine** : [docs/prompts/site-vitrine/PROMPT.md](../docs/prompts/site-vitrine/PROMPT.md)
- **Design Figma** : [Figma Yousoon-Test2](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2?node-id=121-114)

---

*Généré automatiquement - Yousoon v2*
