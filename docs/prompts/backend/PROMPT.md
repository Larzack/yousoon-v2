# ⚙️ Backend Go + GraphQL - Prompt Détaillé

> **Module** : API Backend Yousoon  
> **Technologie** : Go + GraphQL (gqlgen)  
> **Architecture** : Microservices (regroupés par blocs fonctionnels)  
> **Infrastructure** : AKS (Azure Kubernetes Service)  
> **Figma** : [Yousoon-Test2](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2?node-id=121-114)

---

## 🎯 Objectifs

L'API backend doit :
- Servir l'App Flutter, le Site Partenaires et le Site Vitrine
- Garantir une latence < 50ms (P95) pour **TOUTES les requêtes**
- Être scalable horizontalement
- Supporter une architecture microservices par blocs fonctionnels
- Supporter GraphQL Subscriptions pour le temps réel
- Être déployable sur AKS (Azure Kubernetes Service)
- Respecter RGPD (données EU - Irlande)

---

## 🛠️ Stack Technique

### Core

| Technologie | Version | Justification |
|-------------|---------|---------------|
| **Go** | 1.21+ | Performance, concurrence native |
| **gqlgen** | 0.17+ | GraphQL type-safe, génération code |
| **Chi** ou **Fiber** | - | Router HTTP performant |

### Base de données

| Technologie | Usage |
|-------------|-------|
| **MongoDB** | Base principale (flexible, scale) |
| **Redis** | Cache distribué, sessions |

### Messagerie (optionnel)

| Technologie | Usage |
|-------------|-------|
| **NATS** ou **RabbitMQ** | Communication inter-services |
| **Kafka** | Event streaming (si besoin) |

### Observabilité

| Technologie | Usage |
|-------------|-------|
| **OpenTelemetry** | Tracing distribué |
| **Prometheus** | Métriques |
| **Grafana** | Dashboards |
| **Jaeger** | Trace visualization |
| **Loki** | Agrégation logs |
| **Amplitude** | Analytics produit (via SDK) |
| **Sentry** (self-hosted) | Error tracking (Go SDK) |

### Services Externes

| Technologie | Usage |
|-------------|-------|
| **S3 + CloudFront** | CDN images/assets |
| **OneSignal** | Push notifications |
| **Google Maps API** | Géocodage, distances |
| **Amazon SES** | Emails transactionnels |
| **Twilio** | SMS (OTP, notifications) |
| **Apple Pay / Google Pay** | Paiements in-app (100%)

### Tests

| Type | Technologie |
|------|-------------|
| Unit | testing + testify |
| Mock | gomock / mockery |
| Integration | testcontainers-go |
| Load | k6 |

### Sécurité

| Technologie | Usage |
|-------------|-------|
| **JWT** | Tokens authentification |
| **bcrypt** | Hash mots de passe |
| **go-playground/validator** | Validation entrées |

---

## 🏗️ Architecture Microservices

### Vue d'ensemble

```
┌─────────────────────────────────────────────────────────────────────┐
│                         API GATEWAY                                  │
│                    (Kong / Traefik / Custom)                        │
│  - Rate Limiting    - Auth Check     - Load Balancing               │
└────────────────────────────┬────────────────────────────────────────┘
                             │
┌────────────────────────────┴────────────────────────────────────────┐
│                     GRAPHQL GATEWAY                                  │
│                        (gqlgen)                                      │
│  - Schema Stitching/Federation                                       │
│  - Query Complexity Limiting                                         │
│  - Caching (DataLoader)                                              │
└────────────────────────────┬────────────────────────────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  AUTH SERVICE   │ │  USER SERVICE   │ │ PARTNER SERVICE │
│  - Login        │ │  - Profile      │ │  - CRUD Partner │
│  - Register     │ │  - Preferences  │ │  - Establishments│
│  - JWT          │ │  - Favorites    │ │  - Team mgmt    │
│  - Refresh      │ │  - History      │ │                 │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         └───────────────────┴───────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  OFFER SERVICE  │ │ OUTING SERVICE  │ │ NOTIF SERVICE   │
│  - CRUD Offers  │ │  - Bookings     │ │  - Push notif   │
│  - Categories   │ │  - Check-in     │ │  - Email        │
│  - Search       │ │  - QR Codes     │ │  - SMS          │
│  - Filters      │ │  - Calendar     │ │  - Templates    │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         └───────────────────┴───────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      IDENTITY SERVICE                                │
│             (Intégration Onfido/Jumio/Veriff)                       │
│  - CNI Upload      - Verification Status    - Webhook Handler       │
└─────────────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────┐                     ┌─────────────────┐
│     REDIS       │◄────────────────────►     MONGODB     │
│  - Sessions     │                     │  - Users        │
│  - Cache L1     │                     │  - Partners     │
│  - Rate Limit   │                     │  - Offers       │
│  - Pub/Sub      │                     │  - Bookings     │
└─────────────────┘                     └─────────────────┘
```

---

## 📁 Structure Mono-Repo Services

```
services/
├── gateway/                          # GraphQL Gateway
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── server/
│   │   ├── middleware/
│   │   └── config/
│   ├── graph/
│   │   ├── schema.graphqls          # Schema principal
│   │   ├── schema.resolvers.go
│   │   ├── model/
│   │   └── generated.go
│   ├── Dockerfile
│   └── go.mod
│
├── auth-service/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── entity/
│   │   │   │   └── user.go
│   │   │   ├── repository/
│   │   │   │   └── user_repository.go
│   │   │   └── service/
│   │   │       └── auth_service.go
│   │   ├── infrastructure/
│   │   │   ├── mongodb/
│   │   │   ├── redis/
│   │   │   └── jwt/
│   │   ├── interface/
│   │   │   ├── grpc/               # Comm inter-services
│   │   │   └── http/
│   │   └── config/
│   ├── proto/                       # gRPC definitions
│   ├── Dockerfile
│   └── go.mod
│
├── user-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── domain/
│   │   ├── infrastructure/
│   │   └── interface/
│   ├── Dockerfile
│   └── go.mod
│
├── partner-service/
│   └── ... (même structure)
│
├── offer-service/
│   └── ...
│
├── outing-service/
│   └── ...
│
├── notification-service/
│   └── ...
│
├── identity-service/
│   └── ...
│
└── shared/                           # Code partagé
    ├── pkg/
    │   ├── logger/
    │   ├── errors/
    │   ├── validator/
    │   ├── mongodb/
    │   ├── redis/
    │   └── tracing/
    └── go.mod
```

---

## 📊 Schema GraphQL (Exemple)

```graphql
# schema.graphqls

type Query {
  # Auth
  me: User!
  
  # Users
  user(id: ID!): User
  
  # Offers
  offers(input: OffersInput!): OfferConnection!
  offer(id: ID!): Offer
  
  # Partners
  partner(id: ID!): Partner
  partners(input: PartnersInput!): PartnerConnection!
  
  # Outings
  myBookings(status: BookingStatus): [Booking!]!
  booking(id: ID!): Booking
}

type Mutation {
  # Auth
  register(input: RegisterInput!): AuthPayload!
  login(input: LoginInput!): AuthPayload!
  refreshToken(token: String!): AuthPayload!
  logout: Boolean!
  
  # Profile
  updateProfile(input: UpdateProfileInput!): User!
  
  # Identity
  submitIdentityVerification(input: IdentityInput!): IdentityVerification!
  
  # Offers (Partner)
  createOffer(input: CreateOfferInput!): Offer!
  updateOffer(id: ID!, input: UpdateOfferInput!): Offer!
  deleteOffer(id: ID!): Boolean!
  
  # Bookings
  createBooking(offerId: ID!): Booking!
  cancelBooking(id: ID!): Booking!
  checkIn(bookingId: ID!, qrCode: String!): Booking!
  
  # Favorites
  addFavorite(offerId: ID!): Boolean!
  removeFavorite(offerId: ID!): Boolean!
}

type Subscription {
  bookingStatusChanged(bookingId: ID!): Booking!
  newOfferNearby(location: LocationInput!): Offer!
}

# Types

type User {
  id: ID!
  email: String!
  firstName: String!
  lastName: String!
  phone: String
  avatar: String
  isVerified: Boolean!
  identityStatus: IdentityStatus!
  createdAt: DateTime!
  favorites: [Offer!]!
  bookings: [Booking!]!
}

type Partner {
  id: ID!
  name: String!
  description: String
  logo: String
  category: PartnerCategory!
  establishments: [Establishment!]!
  offers: [Offer!]!
  rating: Float
  reviewCount: Int!
  createdAt: DateTime!
}

type Establishment {
  id: ID!
  name: String!
  address: Address!
  location: Location!
  phone: String
  email: String
  website: String
  openingHours: [OpeningHour!]!
  images: [String!]!
}

type Offer {
  id: ID!
  title: String!
  description: String!
  partner: Partner!
  establishment: Establishment!
  category: OfferCategory!
  discount: Discount!
  conditions: [String!]!
  validFrom: DateTime!
  validUntil: DateTime!
  schedule: OfferSchedule
  quota: Quota
  images: [String!]!
  isActive: Boolean!
  isFavorite: Boolean!
  distance: Float
  createdAt: DateTime!
}

type Booking {
  id: ID!
  user: User!
  offer: Offer!
  status: BookingStatus!
  qrCode: String!
  checkedInAt: DateTime
  createdAt: DateTime!
}

# Enums

enum BookingStatus {
  PENDING
  CONFIRMED
  CHECKED_IN
  CANCELLED
  EXPIRED
  NO_SHOW
}

enum IdentityStatus {
  NOT_SUBMITTED
  PENDING
  VERIFIED
  REJECTED
}

enum PartnerCategory {
  BAR
  RESTAURANT
  CLUB
  CINEMA
  SPORT
  LEISURE
  EVENT
  OTHER
}

# Inputs

input OffersInput {
  location: LocationInput
  categories: [OfferCategory!]
  minDiscount: Int
  maxDistance: Float
  search: String
  first: Int
  after: String
}

input RegisterInput {
  email: String!
  password: String!
  firstName: String!
  lastName: String!
  phone: String
}

input LoginInput {
  email: String!
  password: String!
}
```

---

## 🚀 Stratégie de Cache

### Architecture Cache Multi-Niveaux

```
┌─────────────────────────────────────────────────────────────────────┐
│                      CLIENT (Flutter/React)                          │
│                    Cache L0 : Local Storage                          │
└────────────────────────────┬────────────────────────────────────────┘
                             │
┌────────────────────────────┴────────────────────────────────────────┐
│                      GRAPHQL GATEWAY                                 │
│                 Cache L1 : DataLoader (Request-level)                │
│                 - Batch queries                                      │
│                 - Dedupe requests                                    │
└────────────────────────────┬────────────────────────────────────────┘
                             │
┌────────────────────────────┴────────────────────────────────────────┐
│                          REDIS                                       │
│                 Cache L2 : Distributed Cache                         │
│                 - Offers list (TTL: 5min)                           │
│                 - Offer details (TTL: 15min)                        │
│                 - User sessions (TTL: 24h)                          │
│                 - Rate limiting counters                            │
└────────────────────────────┬────────────────────────────────────────┘
                             │
┌────────────────────────────┴────────────────────────────────────────┐
│                         MONGODB                                      │
│                 Cache L3 : Indexes + WiredTiger Cache                │
│                 - Compound indexes                                   │
│                 - Geospatial indexes (2dsphere)                     │
│                 - Text indexes (search)                             │
└─────────────────────────────────────────────────────────────────────┘
```

### Stratégies par Type de Données

| Donnée | Stratégie | TTL Redis | Invalidation |
|--------|-----------|-----------|--------------|
| Offers list | Cache-aside + Stale-while-revalidate | 5 min | TTL + Event |
| Offer detail | Cache-aside | 15 min | On update |
| User profile | Cache-aside | 1h | On update |
| Categories | Cache-aside | 24h | Manual |
| Partner info | Cache-aside | 1h | On update |
| Bookings | Write-through | - | Real-time |
| Sessions | Redis primary | 24h | On logout |

### DataLoader Pattern

```go
// Batch loading pour éviter N+1
type OfferLoader struct {
    wait     time.Duration
    maxBatch int
    cache    map[string]*Offer
    batch    func([]string) ([]*Offer, []error)
}

// Usage dans resolver
func (r *queryResolver) Offers(ctx context.Context, input OffersInput) (*OfferConnection, error) {
    loader := dataloader.For(ctx)
    offers, err := loader.Offer.LoadAll(ctx, offerIDs)
    // ...
}
```

---

## 📏 Contrainte Performance < 50ms

### Optimisations Requises

1. **Indexes MongoDB**
```javascript
// Offres par localisation et catégorie
db.offers.createIndex({ "location": "2dsphere", "category": 1, "isActive": 1 })

// Recherche textuelle
db.offers.createIndex({ "title": "text", "description": "text" })

// Pagination efficace
db.offers.createIndex({ "createdAt": -1, "_id": -1 })
```

2. **Connection Pooling**
```go
// MongoDB connection pool
clientOptions := options.Client().
    SetMaxPoolSize(100).
    SetMinPoolSize(10).
    SetMaxConnIdleTime(30 * time.Second)
```

3. **Projection MongoDB**
```go
// Ne récupérer que les champs nécessaires
projection := bson.M{
    "title": 1,
    "discount": 1,
    "images": bson.M{"$slice": 1},
}
```

4. **Pagination Cursor-based**
```go
// Plus performant que offset
filter := bson.M{
    "_id": bson.M{"$gt": lastID},
}
opts := options.Find().SetLimit(20)
```

---

## 🧪 Tests de Charge (k6)

### Script de Test

```javascript
// tests/load/offers_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 },  // Ramp up
    { duration: '1m', target: 100 },   // Stay at 100 users
    { duration: '30s', target: 200 },  // Ramp to 200
    { duration: '1m', target: 200 },   // Stay at 200
    { duration: '30s', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<50'],   // 95% des requêtes < 50ms
    http_req_failed: ['rate<0.01'],    // < 1% d'erreurs
  },
};

const GRAPHQL_ENDPOINT = 'https://api.yousoon.com/graphql';

const OFFERS_QUERY = `
  query GetOffers($input: OffersInput!) {
    offers(input: $input) {
      edges {
        node {
          id
          title
          discount { percentage }
        }
      }
      pageInfo {
        hasNextPage
      }
    }
  }
`;

export default function () {
  const payload = JSON.stringify({
    query: OFFERS_QUERY,
    variables: {
      input: {
        location: { lat: 48.8566, lng: 2.3522 },
        first: 20,
      },
    },
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${__ENV.AUTH_TOKEN}`,
    },
  };

  const res = http.post(GRAPHQL_ENDPOINT, payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 50ms': (r) => r.timings.duration < 50,
    'no errors': (r) => !JSON.parse(r.body).errors,
  });

  sleep(1);
}
```

### Exécution

```bash
# Local
k6 run tests/load/offers_test.js

# Avec Grafana Cloud
k6 cloud run tests/load/offers_test.js
```

---

## 🔐 Sécurité

### JWT Structure

```go
type Claims struct {
    UserID    string `json:"uid"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    IsPartner bool   `json:"is_partner"`
    jwt.RegisteredClaims
}

// Access token: 15min
// Refresh token: 7 days (stored in Redis)
```

### Rate Limiting

```go
// Redis-based rate limiter
type RateLimiter struct {
    redis   *redis.Client
    limit   int
    window  time.Duration
}

// 100 requests per minute per IP
// 1000 requests per minute per authenticated user
```

### Validation Inputs

```go
type RegisterInput struct {
    Email     string `validate:"required,email"`
    Password  string `validate:"required,min=8,max=72"`
    FirstName string `validate:"required,min=2,max=50"`
    LastName  string `validate:"required,min=2,max=50"`
    Phone     string `validate:"omitempty,e164"`
}
```

---

## 📋 Checklist

- [ ] GraphQL schema complet
- [ ] DataLoader implémenté
- [ ] Cache Redis configuré
- [ ] Indexes MongoDB créés
- [ ] JWT + Refresh tokens
- [ ] Rate limiting
- [ ] Tracing OpenTelemetry
- [ ] Métriques Prometheus
- [ ] Tests unitaires > 80%
- [ ] Tests de charge k6
- [ ] Dockerfiles optimisés
- [ ] Manifests Kubernetes

---

## 🔗 Références

- [Questions à clarifier](./QUESTIONS.md)
- [Architecture détaillée](./ARCHITECTURE.md)
- [Modèle de données](../DATA_MODEL.md)
