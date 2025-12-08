# 🚀 Plan de Génération - Backend Microservices

> **Module** : Backend Go - Microservices DDD  
> **Priorité** : 🔴 CRITIQUE (doit être généré en premier)  
> **Dépendances** : Aucune (module racine)  
> **GraphQL** : Apollo Federation 2 avec gqlgen annotations

---

## 📋 Vue d'Ensemble

```
┌─────────────────────────────────────────────────────────────────┐
│                    ORDRE DE GÉNÉRATION                          │
├─────────────────────────────────────────────────────────────────┤
│  Phase 1: Infrastructure commune (shared, router, registry)     │
│  Phase 2: Core Services/Subgraphs (identity, partner, discovery)│
│  Phase 3: Business Services/Subgraphs (booking, engagement)     │
│  Phase 4: Generic Services/Subgraphs (notification)             │
│  Phase 5: Tests & Observabilité                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🏗️ Architecture GraphQL Federation 2

### Concept

```
┌─────────────────────────────────────────────────────────────────┐
│                     CLIENTS (App, Web)                          │
└──────────────────────────┬──────────────────────────────────────┘
                           │ GraphQL
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                   APOLLO ROUTER (Supergraph)                    │
│  • Compose les subgraphs automatiquement                        │
│  • Query planning & execution                                   │
│  • Auth middleware, rate limiting                               │
│  • Service Discovery via GraphQL Registry                       │
└──────────────────────────┬──────────────────────────────────────┘
                           │ Federation Subgraph Protocol
       ┌───────────────────┼───────────────────┬──────────────────┐
       ▼                   ▼                   ▼                  ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  IDENTITY    │   │   PARTNER    │   │  DISCOVERY   │   │   BOOKING    │
│  Subgraph    │   │   Subgraph   │   │  Subgraph    │   │   Subgraph   │
│              │   │              │   │              │   │              │
│ gqlgen +     │   │ gqlgen +     │   │ gqlgen +     │   │ gqlgen +     │
│ federation   │   │ federation   │   │ federation   │   │ federation   │
└──────────────┘   └──────────────┘   └──────────────┘   └──────────────┘
       │                   │                   │                  │
       └───────────────────┴───────────────────┴──────────────────┘
                           │
                    ┌──────▼──────┐
                    │  GRAPHQL    │
                    │  REGISTRY   │
                    │ (Schema +   │
                    │  Discovery) │
                    └─────────────┘
```

### Stack Federation

| Composant | Technologie | Rôle |
|-----------|-------------|------|
| **Supergraph Router** | Apollo Router | Composition des subgraphs, query planning |
| **Subgraphs** | gqlgen + federation | Chaque microservice expose son schema |
| **Schema Registry** | Apollo GraphOS ou Self-hosted | Stockage schemas + service discovery |
| **Code Generation** | gqlgen annotations | Schema généré depuis le code Go |

### Avantages

1. **Schema-first depuis le code** : Annotations gqlgen → schema auto-généré
2. **Service Discovery auto** : Les subgraphs s'enregistrent au démarrage
3. **Composition automatique** : Le Router compose le supergraph à la volée
4. **Type sharing** : Entités partagées via `@key` directive

---

## 📁 Structure Cible

```
services/
├── shared/                          # Phase 1
│   ├── domain/
│   │   ├── aggregate.go
│   │   ├── entity.go
│   │   ├── valueobject.go
│   │   └── event.go
│   ├── infrastructure/
│   │   ├── mongodb/
│   │   ├── redis/
│   │   ├── nats/
│   │   └── grpc/
│   ├── federation/                  # 🆕 GraphQL Federation shared
│   │   ├── registry/
│   │   │   ├── client.go           # Client pour s'enregistrer au registry
│   │   │   └── discovery.go        # Service discovery
│   │   ├── directives/
│   │   │   ├── auth.go             # @auth directive
│   │   │   └── validation.go       # @constraint directive
│   │   └── scalars/
│   │       ├── datetime.go
│   │       ├── money.go
│   │       └── geolocation.go
│   └── utils/
│
├── router/                          # Phase 1 - Apollo Router
│   ├── config/
│   │   ├── router.yaml             # Config Apollo Router
│   │   └── supergraph.graphql      # Schema composé (généré)
│   ├── plugins/                    # Custom plugins Rust/Rhai
│   │   ├── auth.rhai
│   │   └── ratelimit.rhai
│   └── Dockerfile
│
├── registry/                        # Phase 1 - Schema Registry
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── storage/                # Stockage des schemas
│   │   ├── composer/               # Composition supergraph
│   │   └── api/                    # API REST/GraphQL pour registration
│   └── Dockerfile
│
├── identity-service/                # Phase 2 - Subgraph
│   ├── graph/                      # 🆕 gqlgen folder
│   │   ├── schema.graphqls         # Generated from annotations
│   │   ├── model/                  # Generated models
│   │   ├── resolver.go
│   │   └── generated.go
│   └── ...
│
├── partner-service/                 # Phase 2 - Subgraph
├── discovery-service/               # Phase 2 - Subgraph
├── booking-service/                 # Phase 3 - Subgraph
├── engagement-service/              # Phase 3 - Subgraph
└── notification-service/            # Phase 4 - Subgraph
```

---

## 🔷 Phase 1 : Infrastructure Commune

### Étape 1.1 : Package Shared Domain
**Fichiers à générer :**
```
services/shared/domain/
├── aggregate.go          # AggregateRoot base avec events
├── entity.go             # Entity base avec ID
├── valueobject.go        # Value Objects communs (Email, Money, etc.)
├── event.go              # DomainEvent interface
├── errors.go             # Erreurs domain communes
└── id.go                 # Types ID (UserID, PartnerID, etc.)
```

**Contenu clé :**
- `AggregateRoot` avec gestion des domain events
- Value Objects : `Email`, `Phone`, `Money`, `GeoLocation`, `Address`
- Types ID fortement typés pour chaque aggregate

### Étape 1.2 : Infrastructure MongoDB
**Fichiers à générer :**
```
services/shared/infrastructure/mongodb/
├── client.go             # Client MongoDB avec connection pooling
├── repository.go         # Repository base générique
├── transaction.go        # Support transactions multi-documents
└── mapper.go             # Interface de mapping domain <-> mongo
```

### Étape 1.3 : Infrastructure Redis
**Fichiers à générer :**
```
services/shared/infrastructure/redis/
├── client.go             # Client Redis
├── cache.go              # Cache générique avec TTL
└── distributed_lock.go   # Locks distribués
```

### Étape 1.4 : Infrastructure NATS
**Fichiers à générer :**
```
services/shared/infrastructure/nats/
├── client.go             # Client NATS JetStream
├── publisher.go          # Event Publisher
├── subscriber.go         # Event Subscriber
└── serializer.go         # JSON serialization
```

### Étape 1.5 : GraphQL Federation Shared
**🆕 Fichiers à générer :**
```
services/shared/federation/
├── registry/
│   ├── client.go         # Client pour s'enregistrer au registry
│   ├── discovery.go      # Service discovery (watch for changes)
│   └── health.go         # Health check pour subgraphs
├── directives/
│   ├── auth.go           # @auth(requires: ADMIN) directive
│   ├── validation.go     # @constraint(min: 1, max: 100)
│   └── deprecated.go     # @deprecated directive custom
├── scalars/
│   ├── datetime.go       # DateTime scalar (ISO 8601)
│   ├── money.go          # Money scalar (centimes)
│   ├── geolocation.go    # GeoLocation scalar
│   └── objectid.go       # MongoDB ObjectID scalar
└── middleware/
    ├── context.go        # Context enrichment (user, claims)
    └── dataloader.go     # DataLoader factory pour batching
```

**Contenu clé :**
```go
// registry/client.go
type RegistryClient interface {
    Register(ctx context.Context, subgraph SubgraphInfo) error
    Deregister(ctx context.Context, name string) error
    Heartbeat(ctx context.Context, name string) error
}

type SubgraphInfo struct {
    Name      string    // "identity", "partner", etc.
    URL       string    // "http://identity-service:4000/graphql"
    SchemaSDL string    // Schema SDL généré par gqlgen
    Version   string    // Pour rolling updates
}
```

### Étape 1.6 : Apollo Router (Supergraph)
**Fichiers à générer :**
```
services/router/
├── config/
│   ├── router.yaml           # Configuration Apollo Router
│   └── supergraph.graphql    # Schema composé (auto-généré)
├── plugins/
│   ├── auth.rhai             # Plugin auth custom (Rhai script)
│   ├── ratelimit.rhai        # Rate limiting plugin
│   └── logging.rhai          # Custom logging
├── scripts/
│   ├── compose.sh            # Script de composition des subgraphs
│   └── watch.sh              # Watch mode pour dev
└── Dockerfile
```

**router.yaml :**
```yaml
supergraph:
  introspection: true
  listen: 0.0.0.0:4000
  
# Service Discovery dynamique
subgraphs:
  registry:
    url: http://registry:8080/graphql
    poll_interval: 10s

# Plugins
plugins:
  - path: plugins/auth.rhai
  - path: plugins/ratelimit.rhai

# Headers propagation
headers:
  all:
    - propagate:
        named: Authorization
        rename: Authorization
    - propagate:
        named: X-Request-ID

# Telemetry
telemetry:
  instrumentation:
    spans:
      mode: spec_compliant
  exporters:
    tracing:
      otlp:
        endpoint: http://jaeger:4317
```

### Étape 1.7 : Schema Registry
**Fichiers à générer :**
```
services/registry/
├── cmd/main.go
├── internal/
│   ├── storage/
│   │   ├── store.go          # Interface storage
│   │   ├── memory.go         # In-memory (dev)
│   │   └── redis.go          # Redis (prod)
│   ├── composer/
│   │   ├── composer.go       # Composition du supergraph
│   │   └── validator.go      # Validation des schemas
│   ├── discovery/
│   │   ├── watcher.go        # Watch Kubernetes services
│   │   └── k8s.go            # Kubernetes service discovery
│   └── api/
│       ├── handler.go        # REST API
│       └── graphql.go        # GraphQL API pour introspection
├── config/config.go
└── Dockerfile
```

**API du Registry :**
```go
// POST /subgraphs/:name - Enregistrer un subgraph
// DELETE /subgraphs/:name - Désenregistrer
// GET /subgraphs - Lister tous les subgraphs
// GET /supergraph - Obtenir le schema composé
// GET /health - Health check
```

**Service Discovery Kubernetes :**
```go
// discovery/k8s.go
// Watch les services avec le label: graphql.federation/subgraph=true
// Détecte automatiquement les nouveaux pods et récupère leur schema
```

---

## 🔷 Phase 2 : Core Services (Subgraphs)

### Étape 2.1 : Identity Service (Subgraph)
**Priorité** : 🔴 Critique (auth requise pour tout)

**Fichiers à générer :**
```
services/identity-service/
├── cmd/main.go
├── config/config.go
├── graph/                            # 🆕 gqlgen federation
│   ├── schema.graphqls              # Schema généré depuis annotations
│   ├── federation.graphqls          # Directives federation (@key, etc.)
│   ├── model/
│   │   └── models_gen.go            # Modèles générés
│   ├── resolver.go                  # Resolver principal
│   ├── schema.resolvers.go          # Resolvers générés
│   ├── entity.resolvers.go          # Entity resolvers pour @key
│   └── generated/
│       └── generated.go             # Code gqlgen généré
├── internal/
│   ├── domain/
│   │   ├── aggregate/
│   │   │   └── user.go              # User Aggregate Root avec annotations gqlgen
│   │   ├── entity/
│   │   │   ├── subscription.go
│   │   │   └── identity_verification.go
│   │   ├── valueobject/
│   │   │   ├── profile.go
│   │   │   ├── preferences.go
│   │   │   └── grade.go             # Explorateur, Aventurier, etc.
│   │   ├── event/
│   │   │   └── user_events.go
│   │   ├── repository/
│   │   │   └── user_repository.go
│   │   └── service/
│   │       └── auth_service.go
│   ├── application/
│   │   ├── command/
│   │   │   ├── register_user.go
│   │   │   ├── login_user.go
│   │   │   ├── verify_identity.go
│   │   │   └── subscribe.go
│   │   ├── query/
│   │   │   ├── get_user.go
│   │   │   └── get_subscription.go
│   │   └── service/
│   │       └── identity_service.go
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   └── mongodb/
│   │   │       ├── user_repository.go
│   │   │       └── mapper.go
│   │   ├── auth/
│   │   │   ├── jwt.go
│   │   │   └── biometric.go
│   │   └── external/
│   │       └── ocr_service.go       # Vérification CNI
│   └── interface/
│       └── grpc/
│           ├── server.go
│           └── handler.go
├── proto/
│   └── identity.proto
├── gqlgen.yml                       # 🆕 Configuration gqlgen
└── Dockerfile
```

**gqlgen.yml avec Federation :**
```yaml
schema:
  - graph/*.graphqls

exec:
  filename: graph/generated/generated.go
  package: generated

federation:
  filename: graph/federation.go
  package: graph
  version: 2                         # Apollo Federation 2

model:
  filename: graph/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: graph
  package: graph

# Autobind: génère le schema depuis les annotations Go
autobind:
  - github.com/yousoon/services/identity-service/internal/domain/aggregate
  - github.com/yousoon/services/identity-service/internal/domain/valueobject

directives:
  auth:
    skip_runtime: true
  constraint:
    skip_runtime: true
```

**Annotations gqlgen dans le code Go :**
```go
// internal/domain/aggregate/user.go

// User est l'Aggregate Root du contexte Identity
// @key directive indique que User peut être référencé depuis d'autres subgraphs
type User struct {
    ID        UserID    `json:"id" gqlgen:"id"`
    Email     Email     `json:"email"`
    Profile   Profile   `json:"profile"`
    Grade     UserGrade `json:"grade"`
    Status    UserStatus `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
}

// Génère automatiquement dans le schema:
// type User @key(fields: "id") {
//   id: ID!
//   email: String!
//   profile: Profile!
//   grade: UserGrade!
//   status: UserStatus!
//   createdAt: DateTime!
// }
```

**Schema généré (graph/schema.graphqls) :**
```graphql
extend schema @link(
  url: "https://specs.apollo.dev/federation/v2.3"
  import: ["@key", "@external", "@requires", "@provides", "@shareable"]
)

type User @key(fields: "id") {
  id: ID!
  email: String!
  profile: Profile!
  grade: UserGrade!
  status: UserStatus!
  subscription: Subscription
  createdAt: DateTime!
}

type Profile {
  firstName: String!
  lastName: String!
  displayName: String!
  avatar: String
  birthDate: Date
}

enum UserGrade {
  EXPLORATEUR
  AVENTURIER
  GRAND_VOYAGEUR
  CONQUERANT
}

type Query {
  me: User! @auth
  getSubscriptionPlans: [SubscriptionPlan!]!
}

type Mutation {
  registerUser(input: RegisterInput!): AuthPayload!
  loginUser(email: String!, password: String!): AuthPayload!
  refreshToken(token: String!): AuthPayload!
  verifyIdentity(input: VerifyIdentityInput!): VerificationResult! @auth
  subscribe(planId: ID!): Subscription! @auth
}
```

**Auto-registration au démarrage :**
```go
// cmd/main.go
func main() {
    // ... init
    
    // Générer le schema SDL depuis gqlgen
    schemaSDL := generated.GetSchemaSDL()
    
    // S'enregistrer au registry
    registryClient := federation.NewRegistryClient(cfg.RegistryURL)
    err := registryClient.Register(ctx, federation.SubgraphInfo{
        Name:      "identity",
        URL:       fmt.Sprintf("http://%s:%d/graphql", hostname, port),
        SchemaSDL: schemaSDL,
        Version:   version,
    })
    
    // Heartbeat en background
    go registryClient.StartHeartbeat(ctx, "identity", 10*time.Second)
    
    // Graceful shutdown -> deregister
    defer registryClient.Deregister(ctx, "identity")
    
    // Start server
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
        Resolvers: &graph.Resolver{},
    }))
    http.Handle("/graphql", srv)
}
```

### Étape 2.2 : Partner Service (Subgraph)
**Fichiers à générer :**
```
services/partner-service/
├── cmd/main.go
├── graph/                            # 🆕 gqlgen federation
│   ├── schema.graphqls
│   ├── federation.graphqls
│   ├── model/
│   ├── resolver.go
│   ├── schema.resolvers.go
│   ├── entity.resolvers.go          # Resolve User @external
│   └── generated/
├── internal/
│   ├── domain/
│   │   ├── aggregate/
│   │   │   └── partner.go
│   │   ├── entity/
│   │   │   ├── establishment.go
│   │   │   └── team_member.go
│   │   ├── valueobject/
│   │   │   ├── company.go
│   │   │   ├── branding.go
│   │   │   └── opening_hours.go
│   │   ├── event/
│   │   │   └── partner_events.go
│   │   └── repository/
│   │       └── partner_repository.go
│   ├── application/
│   │   ├── command/
│   │   │   ├── register_partner.go
│   │   │   ├── add_establishment.go
│   │   │   └── invite_team_member.go
│   │   └── query/
│   │       ├── get_partner.go
│   │       └── get_establishments.go
│   └── infrastructure/
│       └── persistence/mongodb/
├── proto/
│   └── partner.proto
├── gqlgen.yml
└── Dockerfile
```

**Schema avec référence cross-subgraph :**
```graphql
extend schema @link(
  url: "https://specs.apollo.dev/federation/v2.3"
  import: ["@key", "@external", "@requires", "@provides", "@shareable"]
)

type Partner @key(fields: "id") {
  id: ID!
  company: Company!
  branding: Branding!
  establishments: [Establishment!]!
  owner: User!                        # Référence vers Identity subgraph
  status: PartnerStatus!
  createdAt: DateTime!
}

# Extension de User depuis Identity subgraph
extend type User @key(fields: "id") {
  id: ID! @external
  partners: [Partner!]!              # Ajoute le champ partners à User
}

type Establishment @key(fields: "id") {
  id: ID!
  name: String!
  address: Address!
  location: GeoLocation!
  openingHours: [OpeningHours!]!
  isActive: Boolean!
}
```

### Étape 2.3 : Discovery Service (Subgraph)
**Fichiers à générer :**
```
services/discovery-service/
├── cmd/main.go
├── graph/
│   ├── schema.graphqls
│   ├── federation.graphqls
│   ├── model/
│   ├── resolver.go
│   ├── schema.resolvers.go
│   ├── entity.resolvers.go
│   └── generated/
├── internal/
│   ├── domain/
│   │   ├── aggregate/
│   │   │   ├── offer.go
│   │   │   └── category.go
│   │   ├── valueobject/
│   │   │   ├── discount.go
│   │   │   ├── schedule.go
│   │   │   └── quota.go
│   │   ├── event/
│   │   │   └── offer_events.go
│   │   └── repository/
│   │       ├── offer_repository.go
│   │       └── category_repository.go
│   ├── application/
│   │   ├── command/
│   │   │   ├── create_offer.go
│   │   │   ├── publish_offer.go
│   │   │   └── update_offer.go
│   │   └── query/
│   │       ├── search_offers.go      # Elasticsearch
│   │       ├── get_nearby_offers.go  # Géospatial
│   │       └── get_recommendations.go
│   └── infrastructure/
│       ├── persistence/mongodb/
│       └── search/
│           └── elasticsearch.go
├── proto/
│   └── discovery.proto
├── gqlgen.yml
└── Dockerfile
```

**Schema avec références multiples :**
```graphql
extend schema @link(
  url: "https://specs.apollo.dev/federation/v2.3"
  import: ["@key", "@external", "@requires", "@provides", "@shareable"]
)

type Offer @key(fields: "id") {
  id: ID!
  title: String!
  description: String!
  discount: Discount!
  validity: Validity!
  schedule: Schedule!
  quota: Quota
  images: [Image!]!
  category: Category!
  establishment: Establishment!       # Référence vers Partner subgraph
  partner: Partner!                   # Référence vers Partner subgraph
  status: OfferStatus!
  createdAt: DateTime!
}

# Extension depuis Partner subgraph
extend type Establishment @key(fields: "id") {
  id: ID! @external
  offers: [Offer!]!                   # Ajoute offers à Establishment
}

extend type Partner @key(fields: "id") {
  id: ID! @external
  offers: [Offer!]!                   # Ajoute offers à Partner
}

type Category @key(fields: "id") {
  id: ID!
  name: LocalizedString!
  slug: String!
  icon: String
  color: String
  parent: Category
}

type Query {
  searchOffers(input: SearchOffersInput!): OfferConnection!
  getNearbyOffers(location: GeoLocationInput!, radius: Float): [Offer!]!
  getRecommendations(limit: Int): [Offer!]! @auth
  getCategories: [Category!]!
}
```

---

## 🔷 Phase 3 : Business Services (Subgraphs)

### Étape 3.1 : Booking Service (Subgraph)
**Fichiers à générer :**
```
services/booking-service/
├── cmd/main.go
├── graph/
│   ├── schema.graphqls
│   ├── federation.graphqls
│   ├── model/
│   ├── resolver.go
│   ├── schema.resolvers.go
│   ├── entity.resolvers.go
│   └── generated/
├── internal/
│   ├── domain/
│   │   ├── aggregate/
│   │   │   └── outing.go         # Réservation
│   │   ├── valueobject/
│   │   │   ├── qrcode.go
│   │   │   └── offer_snapshot.go
│   │   ├── event/
│   │   │   └── outing_events.go
│   │   └── repository/
│   │       └── outing_repository.go
│   ├── application/
│   │   ├── command/
│   │   │   ├── book_outing.go
│   │   │   ├── checkin_outing.go
│   │   │   └── cancel_outing.go
│   │   └── query/
│   │       ├── get_user_outings.go
│   │       └── get_partner_outings.go
│   └── infrastructure/
│       ├── persistence/mongodb/
│       └── acl/
│           └── discovery_acl.go  # Anti-Corruption Layer
├── proto/
│   └── booking.proto
├── gqlgen.yml
└── Dockerfile
```

**Schema Federation :**
```graphql
extend schema @link(
  url: "https://specs.apollo.dev/federation/v2.3"
  import: ["@key", "@external", "@requires", "@provides", "@shareable"]
)

type Outing @key(fields: "id") {
  id: ID!
  user: User!                         # Référence Identity
  offer: OfferSnapshot!               # Snapshot immutable
  qrCode: QRCode!
  status: OutingStatus!
  bookedAt: DateTime!
  checkedInAt: DateTime
  expiresAt: DateTime!
}

type OfferSnapshot @shareable {
  offerId: ID!
  title: String!
  discount: Discount!
  establishmentName: String!
  address: String!
}

type QRCode {
  code: String!
  expiresAt: DateTime!
}

# Extension: Ajouter outings à User
extend type User @key(fields: "id") {
  id: ID! @external
  outings(status: OutingStatus, first: Int): OutingConnection!
}

# Extension: Ajouter outings à Offer
extend type Offer @key(fields: "id") {
  id: ID! @external
  outings(first: Int): [Outing!]!     # Pour les partenaires
}

type Query {
  getOuting(id: ID!): Outing @auth
  getOutingByQR(code: String!): Outing
}

type Mutation {
  bookOuting(offerId: ID!): Outing! @auth
  checkInOuting(outingId: ID!, qrCode: String!): Outing!
  cancelOuting(outingId: ID!, reason: String): Outing! @auth
}

type Subscription {
  outingStatusChanged(outingId: ID!): Outing!
}
```

### Étape 3.2 : Engagement Service (Subgraph)
**Fichiers à générer :**
```
services/engagement-service/
├── cmd/main.go
├── graph/
│   ├── schema.graphqls
│   ├── federation.graphqls
│   ├── model/
│   ├── resolver.go
│   ├── schema.resolvers.go
│   ├── entity.resolvers.go
│   └── generated/
├── internal/
│   ├── domain/
│   │   ├── aggregate/
│   │   │   ├── favorite.go
│   │   │   ├── review.go
│   │   │   └── conversation.go
│   │   ├── entity/
│   │   │   └── message.go
│   │   └── repository/
│   │       ├── favorite_repository.go
│   │       ├── review_repository.go
│   │       └── conversation_repository.go
│   ├── application/
│   │   ├── command/
│   │   │   ├── add_favorite.go
│   │   │   ├── submit_review.go
│   │   │   └── send_message.go
│   │   └── query/
│   │       ├── get_favorites.go
│   │       ├── get_reviews.go
│   │       └── get_conversations.go
│   └── infrastructure/
│       └── persistence/mongodb/
├── proto/
│   └── engagement.proto
├── gqlgen.yml
└── Dockerfile
```

**Schema Federation :**
```graphql
extend schema @link(
  url: "https://specs.apollo.dev/federation/v2.3"
  import: ["@key", "@external", "@requires", "@provides", "@shareable"]
)

type Review @key(fields: "id") {
  id: ID!
  user: User!
  offer: Offer!
  rating: Int!
  content: String
  images: [String!]
  isVerifiedPurchase: Boolean!
  createdAt: DateTime!
}

type Favorite {
  id: ID!
  user: User!
  offer: Offer!
  addedAt: DateTime!
}

type Conversation @key(fields: "id") {
  id: ID!
  participants: [User!]!
  messages: [Message!]!
  lastMessageAt: DateTime
  createdAt: DateTime!
}

type Message {
  id: ID!
  sender: User!
  content: String!
  readAt: DateTime
  createdAt: DateTime!
}

# Extensions
extend type User @key(fields: "id") {
  id: ID! @external
  favorites: [Favorite!]!
  reviews: [Review!]!
  conversations: [Conversation!]!
}

extend type Offer @key(fields: "id") {
  id: ID! @external
  reviews: ReviewConnection!
  averageRating: Float
  reviewCount: Int!
  isFavorited: Boolean! @auth         # Requiert user context
}

type Mutation {
  addFavorite(offerId: ID!): Favorite! @auth
  removeFavorite(offerId: ID!): Boolean! @auth
  submitReview(input: SubmitReviewInput!): Review! @auth
  sendMessage(conversationId: ID!, content: String!): Message! @auth
}

type Subscription {
  newMessage(conversationId: ID!): Message!
}
```
---

## 🔷 Phase 4 : Generic Services (Subgraphs)

### Étape 4.1 : Notification Service (Subgraph)
**Fichiers à générer :**
```
services/notification-service/
├── cmd/main.go
├── graph/
│   ├── schema.graphqls
│   ├── federation.graphqls
│   ├── model/
│   ├── resolver.go
│   ├── schema.resolvers.go
│   └── generated/
├── internal/
│   ├── domain/
│   │   ├── aggregate/
│   │   │   └── notification.go
│   │   ├── entity/
│   │   │   └── template.go
│   │   └── repository/
│   │       └── notification_repository.go
│   ├── application/
│   │   ├── command/
│   │   │   ├── send_push.go
│   │   │   ├── send_email.go
│   │   │   └── send_sms.go
│   │   └── handler/
│   │       └── event_handler.go  # Écoute events des autres services
│   └── infrastructure/
│       ├── persistence/mongodb/
│       └── external/
│           ├── onesignal.go      # Push
│           └── aws_sns.go        # Email/SMS
├── proto/
│   └── notification.proto
├── gqlgen.yml
└── Dockerfile
```

**Schema Federation :**
```graphql
extend schema @link(
  url: "https://specs.apollo.dev/federation/v2.3"
  import: ["@key", "@external"]
)

type Notification @key(fields: "id") {
  id: ID!
  type: NotificationType!
  channel: NotificationChannel!
  content: NotificationContent!
  status: NotificationStatus!
  sentAt: DateTime
  readAt: DateTime
  createdAt: DateTime!
}

type NotificationContent {
  title: String!
  body: String!
  image: String
  data: JSON
}

enum NotificationType {
  BOOKING_CONFIRMED
  BOOKING_REMINDER
  OFFER_NEARBY
  NEW_MESSAGE
  MARKETING
}

# Extension: Ajouter notifications à User
extend type User @key(fields: "id") {
  id: ID! @external
  notifications(first: Int, unreadOnly: Boolean): NotificationConnection!
  unreadNotificationCount: Int!
}

type Query {
  getNotification(id: ID!): Notification @auth
}

type Mutation {
  markNotificationAsRead(id: ID!): Notification! @auth
  markAllNotificationsAsRead: Int! @auth
  updateNotificationPreferences(input: NotificationPreferencesInput!): Boolean! @auth
}

type Subscription {
  newNotification: Notification! @auth
}
```

---

## 🔷 Phase 5 : Tests & Observabilité

### Étape 5.1 : Tests
**Fichiers à générer par service :**
```
services/{service}/
├── tests/
│   ├── unit/
│   │   ├── domain/
│   │   └── application/
│   ├── integration/
│   │   ├── repository_test.go
│   │   └── graphql_test.go          # 🆕 Tests GraphQL
│   └── e2e/
│       └── api_test.go
```

### Étape 5.2 : Observabilité
**Fichiers à générer :**
```
services/shared/observability/
├── tracing/
│   └── opentelemetry.go
├── metrics/
│   └── prometheus.go
└── logging/
    └── structured.go
```

### Étape 5.3 : Kubernetes Manifests
**Fichiers à générer :**
```
deploy/kubernetes/
├── base/
│   ├── namespace.yaml
│   ├── configmap.yaml
│   └── secrets.yaml
├── services/
│   ├── router/                      # 🆕 Apollo Router
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── ingress.yaml
│   │   └── configmap.yaml           # router.yaml
│   ├── registry/                    # 🆕 Schema Registry
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   ├── identity-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── labels.yaml              # graphql.federation/subgraph=true
│   ├── partner-service/
│   ├── discovery-service/
│   ├── booking-service/
│   ├── engagement-service/
│   └── notification-service/
├── infrastructure/
│   ├── mongodb/
│   ├── redis/
│   ├── nats/
│   └── elasticsearch/
└── monitoring/
    ├── prometheus/
    ├── grafana/
    └── jaeger/
```

**Labels pour Service Discovery Kubernetes :**
```yaml
# services/identity-service/deployment.yaml
metadata:
  labels:
    app: identity-service
    graphql.federation/subgraph: "true"
    graphql.federation/name: "identity"
    graphql.federation/port: "4000"
```

---

## 🔄 Workflow de Développement

### Génération de Schema depuis le Code

```bash
# 1. Écrire les modèles Go avec annotations
# internal/domain/aggregate/user.go

# 2. Générer le schema GraphQL et le code
cd services/identity-service
go generate ./...

# Ou utiliser gqlgen directement
go run github.com/99designs/gqlgen generate

# 3. Le schema est généré dans graph/schema.graphqls

# 4. Au démarrage, le service s'enregistre au registry avec son schema
```

### Composition du Supergraph

```bash
# Le Registry compose automatiquement le supergraph quand:
# - Un nouveau subgraph s'enregistre
# - Un subgraph met à jour son schema
# - Un subgraph se désenregistre

# Composition manuelle (debug)
cd services/registry
go run cmd/compose/main.go > config/supergraph.graphql

# Validation
rover supergraph compose --config supergraph.yaml
```

### Watch Mode en Développement

```bash
# Terminal 1: Registry
cd services/registry && go run cmd/main.go

# Terminal 2: Router (watch le registry)
cd services/router && ./router --config config/router.yaml

# Terminal 3-N: Subgraphs
cd services/identity-service && go run cmd/main.go
cd services/partner-service && go run cmd/main.go
# etc.

# Les changements sont propagés automatiquement via le registry
```

---

## ⏱️ Estimation des Temps

| Phase | Étape | Durée estimée |
|-------|-------|---------------|
| **Phase 1** | Shared Domain | 2h |
| | Infrastructure (Mongo, Redis, NATS) | 3h |
| | Federation Shared | 3h |
| | Apollo Router | 2h |
| | Schema Registry | 3h |
| **Phase 2** | Identity Service (Subgraph) | 5h |
| | Partner Service (Subgraph) | 4h |
| | Discovery Service (Subgraph) | 5h |
| **Phase 3** | Booking Service (Subgraph) | 4h |
| | Engagement Service (Subgraph) | 4h |
| **Phase 4** | Notification Service (Subgraph) | 3h |
| **Phase 5** | Tests | 5h |
| | Observabilité | 2h |
| | Kubernetes | 4h |
| **Total** | | **~49h** |

---

## ✅ Critères de Validation

### Par Subgraph
- [ ] Compilation sans erreur
- [ ] Tests unitaires passent (>80% coverage domain)
- [ ] Tests d'intégration passent
- [ ] gRPC server démarre
- [ ] GraphQL server démarre (subgraph)
- [ ] Connexion MongoDB OK
- [ ] Events NATS publiés/consommés
- [ ] Métriques Prometheus exposées
- [ ] **Auto-registration au Registry OK**
- [ ] **Schema SDL généré correctement**

### Apollo Router
- [ ] Composition du supergraph réussie
- [ ] Query planning fonctionne
- [ ] Cross-subgraph queries OK
- [ ] Auth plugin actif
- [ ] Rate limiting actif
- [ ] Tracing propagé

### Schema Registry
- [ ] API REST fonctionnelle
- [ ] Composition automatique
- [ ] Service Discovery Kubernetes OK
- [ ] Persistance Redis OK
- [ ] Health checks OK

### Global
- [ ] `docker-compose up` démarre tous les services
- [ ] Latence < 50ms (tests de charge)
- [ ] Tracing Jaeger visible
- [ ] **Supergraph Schema correct**
- [ ] **Federation queries cross-subgraph OK**

---

## 🔗 Références

- [Architecture DDD](./ARCHITECTURE.md)
- [Data Model](../DATA_MODEL.md)
- [Copilot Instructions](../../.github/copilot-instructions.md)
- [Apollo Federation 2 Docs](https://www.apollographql.com/docs/federation/)
- [gqlgen Federation](https://gqlgen.com/recipes/federation/)
- [Apollo Router](https://www.apollographql.com/docs/router/)
