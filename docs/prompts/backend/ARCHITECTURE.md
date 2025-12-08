# 🏛️ Architecture DDD - Microservices Yousoon

> Architecture Domain-Driven Design pour la plateforme Yousoon  
> **Pattern** : Hexagonal Architecture (Ports & Adapters)  
> **Communication** : gRPC (sync) + NATS JetStream (async events)  
> **API Gateway** : Apollo Federation 2 (GraphQL)  
> **Infrastructure** : AWS EKS (Elastic Kubernetes Service) - Région Irlande (RGPD)

---

## 🛠️ Stack Technique Validée

| Composant | Technologie |
|-----------|-------------|
| **Cloud** | AWS EKS (Kubernetes) |
| **API Gateway** | Apollo Router (Federation 2) |
| **GraphQL** | gqlgen avec annotations + federation |
| **Service Discovery** | Schema Registry custom + Kubernetes labels |
| **API Sync** | gRPC + protobuf (inter-service) |
| **API Async** | NATS JetStream |
| **Database** | MongoDB (par context) |
| **Cache** | Redis |
| **Search** | Elasticsearch |
| **Storage** | AWS S3 + CloudFront |
| **Notifications** | OneSignal (Push) + AWS SNS (Email/SMS) |
| **Analytics** | Amplitude |
| **Vérification CNI** | OCR interne (Tesseract/OpenCV) |
| **Observability** | OpenTelemetry + Jaeger + Prometheus + Loki + Grafana |
| **Langues** | FR + EN (traduction auto) |
| **Mode Offline** | Oui (app mobile) |

---

## 🌐 Architecture GraphQL Federation 2

### Vue d'Ensemble

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CLIENTS (App Mobile, Web)                          │
└──────────────────────────────────┬──────────────────────────────────────────┘
                                   │ GraphQL (HTTPS)
                                   ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         APOLLO ROUTER (Supergraph)                           │
│                                                                              │
│  • Compose automatiquement les subgraphs                                     │
│  • Query planning & execution distribuée                                     │
│  • Auth middleware (JWT validation)                                          │
│  • Rate limiting & caching                                                   │
│  • Tracing OpenTelemetry                                                     │
│                                                                              │
│  Plugins: auth.rhai, ratelimit.rhai, logging.rhai                           │
└──────────────────────────────────┬──────────────────────────────────────────┘
                                   │ Federation Protocol
       ┌───────────────────────────┼───────────────────────────┐
       ▼                           ▼                           ▼
┌──────────────┐           ┌──────────────┐           ┌──────────────┐
│   IDENTITY   │           │   PARTNER    │           │  DISCOVERY   │
│   Subgraph   │           │   Subgraph   │           │   Subgraph   │
│              │           │              │           │              │
│ gqlgen +     │◄─────────►│ gqlgen +     │◄─────────►│ gqlgen +     │
│ federation   │   gRPC    │ federation   │   gRPC    │ federation   │
└──────┬───────┘           └──────┬───────┘           └──────┬───────┘
       │                          │                          │
       ▼                          ▼                          ▼
┌──────────────┐           ┌──────────────┐           ┌──────────────┐
│   BOOKING    │           │  ENGAGEMENT  │           │ NOTIFICATION │
│   Subgraph   │           │   Subgraph   │           │   Subgraph   │
└──────────────┘           └──────────────┘           └──────────────┘
       │                          │                          │
       └──────────────────────────┴──────────────────────────┘
                                  │
                           ┌──────▼──────┐
                           │   SCHEMA    │
                           │  REGISTRY   │
                           │             │
                           │ • Stockage  │
                           │ • Discovery │
                           │ • Compose   │
                           └─────────────┘
```

### Concepts Clés Federation 2

| Directive | Usage |
|-----------|-------|
| `@key` | Définit l'identifiant unique pour référencer une entité cross-subgraph |
| `@external` | Champ défini dans un autre subgraph |
| `@requires` | Champs externes requis pour résoudre un champ |
| `@provides` | Champs fournis par ce subgraph pour une entité externe |
| `@shareable` | Champ pouvant être résolu par plusieurs subgraphs |

### Exemple de Type Partagé

```graphql
# Dans Identity Subgraph
type User @key(fields: "id") {
  id: ID!
  email: String!
  profile: Profile!
}

# Dans Booking Subgraph (extension)
extend type User @key(fields: "id") {
  id: ID! @external
  outings: [Outing!]!        # Nouveau champ ajouté par Booking
}

# Dans Engagement Subgraph (extension)
extend type User @key(fields: "id") {
  id: ID! @external
  favorites: [Favorite!]!    # Nouveau champ ajouté par Engagement
  reviews: [Review!]!
}

# Le Router compose automatiquement:
type User {
  id: ID!
  email: String!
  profile: Profile!
  outings: [Outing!]!        # Vient de Booking
  favorites: [Favorite!]!    # Vient de Engagement
  reviews: [Review!]!        # Vient de Engagement
}
```

### Service Discovery

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          SERVICE DISCOVERY FLOW                              │
└─────────────────────────────────────────────────────────────────────────────┘

1. Subgraph démarre
   ├── Génère son schema SDL (gqlgen)
   └── S'enregistre au Schema Registry
   
2. Schema Registry
   ├── Stocke le schema (Redis)
   ├── Valide la compatibilité
   └── Re-compose le supergraph
   
3. Apollo Router
   ├── Poll le Registry (interval: 10s)
   ├── Détecte le nouveau supergraph
   └── Hot-reload la configuration

4. Kubernetes (backup discovery)
   ├── Watch services avec label: graphql.federation/subgraph=true
   └── Fallback si Registry indisponible
```

---

## 📋 Table des Matières

1. [Bounded Contexts](#bounded-contexts)
2. [Context Map](#context-map)
3. [Ubiquitous Language](#ubiquitous-language)
4. [Architecture Hexagonale](#architecture-hexagonale)
5. [Aggregates & Entities](#aggregates--entities)
6. [Value Objects](#value-objects)
7. [Domain Events](#domain-events)
8. [Structure des Services](#structure-des-services)
9. [Anti-Corruption Layers](#anti-corruption-layers)

---

## Bounded Contexts

### Vue d'Ensemble Stratégique

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              YOUSOON PLATFORM                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │    IDENTITY     │  │    PARTNER      │  │   DISCOVERY     │              │
│  │    CONTEXT      │  │    CONTEXT      │  │    CONTEXT      │              │
│  │                 │  │                 │  │                 │              │
│  │ • Authentication│  │ • Partner Mgmt  │  │ • Offer Catalog │              │
│  │ • User Profile  │  │ • Establishment │  │ • Search        │              │
│  │ • Verification  │  │ • Team          │  │ • Recommendations│             │
│  │ • Subscription  │  │ • Analytics     │  │ • Categories    │              │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘              │
│           │                    │                    │                       │
│           └────────────────────┼────────────────────┘                       │
│                                │                                            │
│  ┌─────────────────┐  ┌────────┴────────┐  ┌─────────────────┐              │
│  │   ENGAGEMENT    │  │    BOOKING      │  │  NOTIFICATION   │              │
│  │    CONTEXT      │  │    CONTEXT      │  │    CONTEXT      │              │
│  │                 │  │                 │  │                 │              │
│  │ • Favorites     │  │ • Reservations  │  │ • Push          │              │
│  │ • Reviews       │  │ • Check-in (QR) │  │ • Email         │              │
│  │ • Social        │  │ • Outing History│  │ • SMS           │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Définition des Bounded Contexts

| Context | Responsabilité | Core/Support/Generic |
|---------|---------------|----------------------|
| **Identity** | Authentification, profils, vérification CNI, abonnements | **Core Domain** |
| **Partner** | Gestion partenaires, établissements, équipes | **Core Domain** |
| **Discovery** | Catalogue offres, recherche, recommandations | **Core Domain** |
| **Booking** | Réservations, check-in QR, historique sorties | **Core Domain** |
| **Engagement** | Favoris, avis, interactions sociales | **Supporting Domain** |
| **Notification** | Envoi push/email/SMS | **Generic Domain** |

---

## Context Map

### Relations entre Bounded Contexts

```
                    ┌──────────────┐
                    │   IDENTITY   │
                    │   (Core)     │
                    └──────┬───────┘
                           │
           ┌───────────────┼───────────────┐
           │ U/D           │ U/D           │ U/D
           ▼               ▼               ▼
    ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
    │   PARTNER    │ │  DISCOVERY   │ │   BOOKING    │
    │   (Core)     │ │   (Core)     │ │   (Core)     │
    └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
           │                │                │
           │ D/U            │ D/U            │ U/D
           ▼                ▼                ▼
    ┌──────────────────────────────────────────────┐
    │              ENGAGEMENT (Supporting)          │
    └──────────────────────────────────────────────┘
                           │
                           │ U/D (ACL)
                           ▼
                    ┌──────────────┐
                    │ NOTIFICATION │
                    │  (Generic)   │
                    └──────────────┘

Légende:
  U = Upstream (fournit des données)
  D = Downstream (consomme des données)
  ACL = Anti-Corruption Layer
```

### Types de Relations

| Relation | Upstream | Downstream | Pattern |
|----------|----------|------------|---------|
| Identity → Partner | Identity | Partner | **Customer/Supplier** |
| Identity → Booking | Identity | Booking | **Customer/Supplier** |
| Partner → Discovery | Partner | Discovery | **Conformist** |
| Discovery → Booking | Discovery | Booking | **Shared Kernel** (OfferSnapshot) |
| * → Notification | Tous | Notification | **ACL** (Anti-Corruption Layer) |

---

## Ubiquitous Language

### Glossaire Métier

| Terme | Définition | Context |
|-------|------------|---------|
| **Yousooner** | Utilisateur vérifié de l'application | Identity |
| **Partner** | Entreprise proposant des offres | Partner |
| **Establishment** | Lieu physique d'un partenaire | Partner |
| **Offer** | Réduction ou sortie proposée | Discovery |
| **Discount** | Pourcentage ou montant de réduction | Discovery |
| **Outing** | Réservation d'une offre par un utilisateur | Booking |
| **Check-in** | Validation de présence via QR code | Booking |
| **Grade** | Niveau de l'utilisateur (Explorateur → Conquérant) | Identity |
| **Subscription** | Abonnement payant (via In-App Purchase) | Identity |

### Règles Métier Clés

```yaml
Identity:
  - Un Yousooner doit avoir une CNI vérifiée pour réserver
  - Un compte peut être utilisateur ET partenaire
  - L'abonnement est géré 100% via Apple/Google Pay

Partner:
  - Un partenaire peut avoir plusieurs établissements
  - Chaque établissement a sa propre géolocalisation
  - Les équipes ont des rôles (admin, manager, staff)

Discovery:
  - Une offre est toujours rattachée à un établissement
  - Les offres ont une période de validité
  - Le rayon de recherche par défaut est de 10km

Booking:
  - Une réservation n'est valide qu'une seule fois
  - Le check-in se fait uniquement par QR code
  - L'utilisateur a 30 min après réservation pour check-in

Engagement:
  - Un utilisateur ne peut laisser qu'un seul avis par offre
  - Les partenaires ne peuvent PAS répondre aux avis
```

---

## Architecture Hexagonale

### Pattern Ports & Adapters

```
                         ┌─────────────────────────────────────┐
                         │          DRIVING ADAPTERS           │
                         │     (Primary/Input Adapters)        │
                         ├─────────────────────────────────────┤
                         │  ┌─────────┐  ┌─────────┐  ┌─────┐  │
                         │  │  gRPC   │  │  HTTP   │  │ CLI │  │
                         │  │ Handler │  │ Handler │  │     │  │
                         │  └────┬────┘  └────┬────┘  └──┬──┘  │
                         └───────┼────────────┼─────────┼─────┘
                                 │            │         │
                         ┌───────▼────────────▼─────────▼─────┐
                         │           INPUT PORTS               │
                         │     (Use Cases / Application)       │
                         ├─────────────────────────────────────┤
                         │                                     │
                         │  ┌─────────────────────────────┐    │
                         │  │       APPLICATION LAYER      │    │
                         │  │                             │    │
                         │  │  • Command Handlers         │    │
                         │  │  • Query Handlers           │    │
                         │  │  • Use Cases                │    │
                         │  └─────────────┬───────────────┘    │
                         │                │                    │
                         │  ┌─────────────▼───────────────┐    │
                         │  │         DOMAIN LAYER         │    │
                         │  │                             │    │
                         │  │  • Aggregates               │    │
                         │  │  • Entities                 │    │
                         │  │  • Value Objects            │    │
                         │  │  • Domain Events            │    │
                         │  │  • Domain Services          │    │
                         │  │  • Repository Interfaces    │    │
                         │  └─────────────────────────────┘    │
                         │                                     │
                         └───────────────┬─────────────────────┘
                                         │
                         ┌───────────────▼─────────────────────┐
                         │          OUTPUT PORTS               │
                         │    (Repository Interfaces)          │
                         └───────────────┬─────────────────────┘
                                         │
                         ┌───────────────▼─────────────────────┐
                         │         DRIVEN ADAPTERS             │
                         │    (Secondary/Output Adapters)      │
                         ├─────────────────────────────────────┤
                         │  ┌─────────┐  ┌─────────┐  ┌─────┐  │
                         │  │ MongoDB │  │  Redis  │  │NATS │  │
                         │  │ Adapter │  │ Adapter │  │Pub  │  │
                         │  └─────────┘  └─────────┘  └─────┘  │
                         └─────────────────────────────────────┘
```

---

## Aggregates & Entities

### 1. Identity Context

```go
// ============================================
// AGGREGATE ROOT: User
// ============================================

package identity

import (
    "time"
    "github.com/yousoon/shared/domain"
)

// User est l'Aggregate Root du contexte Identity
type User struct {
    domain.AggregateRoot
    
    // Identity
    id           UserID
    email        Email
    passwordHash PasswordHash
    phone        *Phone
    
    // Profile (Value Object)
    profile      Profile
    
    // Verification (Entity)
    identity     *IdentityVerification
    
    // Subscription (Entity)
    subscription *Subscription
    
    // Preferences (Value Object)
    preferences  Preferences
    
    // State
    status       UserStatus
    grade        UserGrade
    
    // Metadata
    createdAt    time.Time
    updatedAt    time.Time
    deletedAt    *time.Time
}

// Invariants (Business Rules)
func (u *User) CanBook() error {
    if u.status != UserStatusActive {
        return ErrUserNotActive
    }
    if u.identity == nil || u.identity.Status != VerificationStatusVerified {
        return ErrIdentityNotVerified
    }
    return nil
}

func (u *User) CanCreateOffer() error {
    if !u.HasPartnerRole() {
        return ErrNotAPartner
    }
    return nil
}

// Commands
func (u *User) VerifyIdentity(verification IdentityVerification) error {
    if u.identity != nil && u.identity.Status == VerificationStatusVerified {
        return ErrAlreadyVerified
    }
    u.identity = &verification
    u.AddDomainEvent(UserIdentityVerified{
        UserID:    u.id,
        Method:    verification.Method,
        Timestamp: time.Now(),
    })
    return nil
}

func (u *User) Subscribe(plan SubscriptionPlan, receipt InAppReceipt) error {
    if u.subscription != nil && u.subscription.IsActive() {
        return ErrAlreadySubscribed
    }
    
    sub, err := NewSubscription(u.id, plan, receipt)
    if err != nil {
        return err
    }
    
    u.subscription = sub
    u.AddDomainEvent(UserSubscribed{
        UserID:   u.id,
        PlanID:   plan.ID,
        Platform: receipt.Platform,
    })
    return nil
}

// ============================================
// ENTITY: IdentityVerification
// ============================================

type IdentityVerification struct {
    id           VerificationID
    status       VerificationStatus
    documentType DocumentType
    method       VerificationMethod  // internal_ocr, external_provider
    submittedAt  time.Time
    verifiedAt   *time.Time
    rejectedAt   *time.Time
    reason       *string
}

// ============================================
// ENTITY: Subscription
// ============================================

type Subscription struct {
    id            SubscriptionID
    planID        PlanID
    platform      Platform  // ios, android
    transactionID string
    startDate     time.Time
    endDate       time.Time
    autoRenew     bool
    cancelledAt   *time.Time
}

func (s *Subscription) IsActive() bool {
    return time.Now().Before(s.endDate) && s.cancelledAt == nil
}
```

### 2. Partner Context

```go
// ============================================
// AGGREGATE ROOT: Partner
// ============================================

package partner

type Partner struct {
    domain.AggregateRoot
    
    id              PartnerID
    ownerUserID     UserID  // Référence cross-context
    
    // Company (Value Object)
    company         Company
    
    // Establishments (Entities - partie de l'Aggregate)
    establishments  []Establishment
    
    // Team (Entities)
    teamMembers     []TeamMember
    
    // State
    status          PartnerStatus
    verifiedAt      *time.Time
    
    createdAt       time.Time
    updatedAt       time.Time
}

// Invariants
func (p *Partner) CanPublishOffer() error {
    if p.status != PartnerStatusVerified {
        return ErrPartnerNotVerified
    }
    if len(p.establishments) == 0 {
        return ErrNoEstablishment
    }
    return nil
}

// Commands
func (p *Partner) AddEstablishment(est Establishment) error {
    if p.HasEstablishment(est.Address) {
        return ErrEstablishmentAlreadyExists
    }
    
    p.establishments = append(p.establishments, est)
    p.AddDomainEvent(EstablishmentAdded{
        PartnerID:       p.id,
        EstablishmentID: est.ID,
        Name:            est.Name,
    })
    return nil
}

func (p *Partner) AddTeamMember(member TeamMember) error {
    if p.HasTeamMember(member.Email) {
        return ErrTeamMemberExists
    }
    
    p.teamMembers = append(p.teamMembers, member)
    p.AddDomainEvent(TeamMemberInvited{
        PartnerID: p.id,
        Email:     member.Email,
        Role:      member.Role,
    })
    return nil
}

// ============================================
// ENTITY: Establishment
// ============================================

type Establishment struct {
    id          EstablishmentID
    name        string
    description string
    address     Address        // Value Object
    location    GeoLocation    // Value Object
    photos      []Photo
    openingHours OpeningHours  // Value Object
    isActive    bool
}

// ============================================
// ENTITY: TeamMember
// ============================================

type TeamMember struct {
    id       TeamMemberID
    userID   *UserID  // nil si invitation en attente
    email    Email
    role     TeamRole  // admin, manager, staff
    invitedAt time.Time
    joinedAt  *time.Time
}
```

### 3. Discovery Context

```go
// ============================================
// AGGREGATE ROOT: Offer
// ============================================

package discovery

type Offer struct {
    domain.AggregateRoot
    
    id              OfferID
    partnerID       PartnerID        // Cross-context reference
    establishmentID EstablishmentID  // Cross-context reference
    
    // Core
    title           string
    description     string
    category        CategoryID
    
    // Discount (Value Object)
    discount        Discount
    
    // Schedule (Value Object)
    schedule        Schedule
    
    // Media
    images          []Image
    
    // Constraints
    maxParticipants *int
    currentBookings int
    
    // State
    status          OfferStatus
    
    createdAt       time.Time
    updatedAt       time.Time
    publishedAt     *time.Time
}

// Invariants
func (o *Offer) CanBeBooked() error {
    if o.status != OfferStatusPublished {
        return ErrOfferNotPublished
    }
    if o.schedule.IsExpired() {
        return ErrOfferExpired
    }
    if o.maxParticipants != nil && o.currentBookings >= *o.maxParticipants {
        return ErrOfferFullyBooked
    }
    return nil
}

// Commands
func (o *Offer) Publish() error {
    if o.status == OfferStatusPublished {
        return ErrAlreadyPublished
    }
    
    now := time.Now()
    o.status = OfferStatusPublished
    o.publishedAt = &now
    
    o.AddDomainEvent(OfferPublished{
        OfferID:         o.id,
        PartnerID:       o.partnerID,
        EstablishmentID: o.establishmentID,
        Category:        o.category,
        Location:        o.GetLocation(), // Via Establishment
    })
    return nil
}

func (o *Offer) IncrementBookings() error {
    if err := o.CanBeBooked(); err != nil {
        return err
    }
    o.currentBookings++
    return nil
}

// ============================================
// AGGREGATE ROOT: Category
// ============================================

type Category struct {
    domain.AggregateRoot
    
    id          CategoryID
    name        string
    slug        string
    icon        string  // Emoji ou URL icône
    color       string  // Hex color
    parentID    *CategoryID
    sortOrder   int
    isActive    bool
}
```

### 4. Booking Context

```go
// ============================================
// AGGREGATE ROOT: Outing (Réservation)
// ============================================

package booking

type Outing struct {
    domain.AggregateRoot
    
    id        OutingID
    userID    UserID
    
    // Snapshot de l'offre au moment de la réservation (Value Object)
    offer     OfferSnapshot
    
    // QR Code (Value Object)
    qrCode    QRCode
    
    // Timeline
    status    OutingStatus
    bookedAt  time.Time
    expiresAt time.Time  // 30 min après booking
    checkedInAt *time.Time
    cancelledAt *time.Time
    
    // Metadata
    createdAt time.Time
    updatedAt time.Time
}

// Invariants
func (o *Outing) CanCheckIn() error {
    if o.status != OutingStatusBooked {
        return ErrInvalidStatus
    }
    if time.Now().After(o.expiresAt) {
        return ErrOutingExpired
    }
    return nil
}

// Commands
func NewOuting(userID UserID, offer OfferSnapshot) (*Outing, error) {
    now := time.Now()
    
    outing := &Outing{
        id:        NewOutingID(),
        userID:    userID,
        offer:     offer,
        qrCode:    GenerateQRCode(),
        status:    OutingStatusBooked,
        bookedAt:  now,
        expiresAt: now.Add(30 * time.Minute),
        createdAt: now,
        updatedAt: now,
    }
    
    outing.AddDomainEvent(OutingBooked{
        OutingID:  outing.id,
        UserID:    userID,
        OfferID:   offer.ID,
        PartnerID: offer.PartnerID,
    })
    
    return outing, nil
}

func (o *Outing) CheckIn(scannedQR string) error {
    if err := o.CanCheckIn(); err != nil {
        return err
    }
    
    if !o.qrCode.Matches(scannedQR) {
        return ErrInvalidQRCode
    }
    
    now := time.Now()
    o.status = OutingStatusCheckedIn
    o.checkedInAt = &now
    o.updatedAt = now
    
    o.AddDomainEvent(OutingCheckedIn{
        OutingID:  o.id,
        UserID:    o.userID,
        OfferID:   o.offer.ID,
        PartnerID: o.offer.PartnerID,
        Timestamp: now,
    })
    
    return nil
}

func (o *Outing) Cancel(reason string) error {
    if o.status == OutingStatusCheckedIn {
        return ErrCannotCancelCheckedIn
    }
    
    now := time.Now()
    o.status = OutingStatusCancelled
    o.cancelledAt = &now
    o.updatedAt = now
    
    o.AddDomainEvent(OutingCancelled{
        OutingID: o.id,
        UserID:   o.userID,
        OfferID:  o.offer.ID,
        Reason:   reason,
    })
    
    return nil
}
```

---

## Value Objects

```go
// ============================================
// VALUE OBJECTS
// ============================================

package domain

// Email - Value Object
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return Email{}, ErrInvalidEmail
    }
    return Email{value: strings.ToLower(email)}, nil
}

func (e Email) String() string { return e.value }
func (e Email) Equals(other Email) bool { return e.value == other.value }

// Money - Value Object
type Money struct {
    amount   int64   // En centimes
    currency string  // ISO 4217
}

func NewMoney(amount int64, currency string) Money {
    return Money{amount: amount, currency: currency}
}

func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, ErrCurrencyMismatch
    }
    return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

// Discount - Value Object
type Discount struct {
    discountType DiscountType  // percentage, fixed
    value        int           // % ou centimes
    minPurchase  *Money
    maxDiscount  *Money
}

func (d Discount) Apply(original Money) Money {
    switch d.discountType {
    case DiscountTypePercentage:
        reduction := original.amount * int64(d.value) / 100
        if d.maxDiscount != nil && reduction > d.maxDiscount.amount {
            reduction = d.maxDiscount.amount
        }
        return Money{amount: original.amount - reduction, currency: original.currency}
    case DiscountTypeFixed:
        return Money{amount: original.amount - int64(d.value), currency: original.currency}
    }
    return original
}

// GeoLocation - Value Object
type GeoLocation struct {
    longitude float64
    latitude  float64
}

func NewGeoLocation(lng, lat float64) (GeoLocation, error) {
    if lng < -180 || lng > 180 || lat < -90 || lat > 90 {
        return GeoLocation{}, ErrInvalidCoordinates
    }
    return GeoLocation{longitude: lng, latitude: lat}, nil
}

func (g GeoLocation) DistanceTo(other GeoLocation) float64 {
    // Haversine formula
    return haversine(g.latitude, g.longitude, other.latitude, other.longitude)
}

// Address - Value Object
type Address struct {
    street     string
    city       string
    postalCode string
    country    string
}

// Schedule - Value Object
type Schedule struct {
    startDate time.Time
    endDate   time.Time
    timeSlots []TimeSlot
    recurring RecurringPattern
}

func (s Schedule) IsExpired() bool {
    return time.Now().After(s.endDate)
}

func (s Schedule) IsActiveNow() bool {
    now := time.Now()
    return now.After(s.startDate) && now.Before(s.endDate)
}

// QRCode - Value Object
type QRCode struct {
    code      string
    signature string
    createdAt time.Time
}

func GenerateQRCode() QRCode {
    code := uuid.New().String()
    signature := hmacSign(code, secretKey)
    return QRCode{
        code:      code,
        signature: signature,
        createdAt: time.Now(),
    }
}

func (q QRCode) Matches(scanned string) bool {
    return q.code == scanned || q.FullCode() == scanned
}

func (q QRCode) FullCode() string {
    return fmt.Sprintf("%s.%s", q.code, q.signature)
}

// OfferSnapshot - Value Object (immutable copy for Booking)
type OfferSnapshot struct {
    ID              OfferID
    PartnerID       PartnerID
    EstablishmentID EstablishmentID
    Title           string
    Description     string
    Discount        Discount
    Category        string
    Location        GeoLocation
    CapturedAt      time.Time
}

// Profile - Value Object
type Profile struct {
    firstName   string
    lastName    string
    displayName string
    avatar      *string
    birthDate   *time.Time
    gender      *Gender
}

func (p Profile) FullName() string {
    return fmt.Sprintf("%s %s", p.firstName, p.lastName)
}

func (p Profile) Age() int {
    if p.birthDate == nil {
        return 0
    }
    return int(time.Since(*p.birthDate).Hours() / 24 / 365)
}
```

---

## Domain Events

### Catalogue des Events

```go
// ============================================
// DOMAIN EVENTS
// ============================================

package events

// Base event
type DomainEvent interface {
    EventName() string
    OccurredAt() time.Time
    AggregateID() string
}

// ========================
// Identity Context Events
// ========================

type UserRegistered struct {
    UserID    UserID
    Email     Email
    Platform  string  // ios, android, web
    Timestamp time.Time
}

type UserIdentityVerified struct {
    UserID    UserID
    Method    VerificationMethod
    Timestamp time.Time
}

type UserSubscribed struct {
    UserID       UserID
    PlanID       PlanID
    Platform     Platform
    TransactionID string
    Timestamp    time.Time
}

type UserSubscriptionCancelled struct {
    UserID    UserID
    PlanID    PlanID
    Reason    string
    Timestamp time.Time
}

type UserDeleted struct {
    UserID    UserID
    Reason    string  // gdpr_request, self_delete
    Timestamp time.Time
}

// ========================
// Partner Context Events
// ========================

type PartnerRegistered struct {
    PartnerID PartnerID
    OwnerID   UserID
    Company   string
    Timestamp time.Time
}

type PartnerVerified struct {
    PartnerID PartnerID
    VerifiedBy UserID  // Admin
    Timestamp time.Time
}

type EstablishmentAdded struct {
    PartnerID       PartnerID
    EstablishmentID EstablishmentID
    Name            string
    Location        GeoLocation
    Timestamp       time.Time
}

type TeamMemberInvited struct {
    PartnerID PartnerID
    Email     Email
    Role      TeamRole
    Timestamp time.Time
}

// ========================
// Discovery Context Events
// ========================

type OfferCreated struct {
    OfferID         OfferID
    PartnerID       PartnerID
    EstablishmentID EstablishmentID
    Title           string
    Category        CategoryID
    Timestamp       time.Time
}

type OfferPublished struct {
    OfferID         OfferID
    PartnerID       PartnerID
    EstablishmentID EstablishmentID
    Category        CategoryID
    Location        GeoLocation
    Schedule        Schedule
    Timestamp       time.Time
}

type OfferExpired struct {
    OfferID   OfferID
    PartnerID PartnerID
    Timestamp time.Time
}

// ========================
// Booking Context Events
// ========================

type OutingBooked struct {
    OutingID  OutingID
    UserID    UserID
    OfferID   OfferID
    PartnerID PartnerID
    Timestamp time.Time
}

type OutingCheckedIn struct {
    OutingID  OutingID
    UserID    UserID
    OfferID   OfferID
    PartnerID PartnerID
    Timestamp time.Time
}

type OutingCancelled struct {
    OutingID  OutingID
    UserID    UserID
    OfferID   OfferID
    Reason    string
    Timestamp time.Time
}

type OutingExpired struct {
    OutingID  OutingID
    UserID    UserID
    OfferID   OfferID
    Timestamp time.Time
}

// ========================
// Engagement Context Events
// ========================

type OfferAddedToFavorites struct {
    UserID    UserID
    OfferID   OfferID
    Timestamp time.Time
}

type ReviewSubmitted struct {
    ReviewID  ReviewID
    UserID    UserID
    OfferID   OfferID
    Rating    int
    Timestamp time.Time
}
```

### Event Bus (NATS)

```go
// ============================================
// EVENT PUBLISHER
// ============================================

package infrastructure

type EventPublisher interface {
    Publish(ctx context.Context, event events.DomainEvent) error
    PublishAll(ctx context.Context, events []events.DomainEvent) error
}

type NATSEventPublisher struct {
    conn *nats.Conn
    js   nats.JetStreamContext
}

func (p *NATSEventPublisher) Publish(ctx context.Context, event events.DomainEvent) error {
    subject := fmt.Sprintf("yousoon.events.%s", event.EventName())
    
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    _, err = p.js.Publish(subject, data, nats.Context(ctx))
    return err
}

// ============================================
// EVENT HANDLERS (Subscribers)
// ============================================

// Dans Notification Service
type NotificationEventHandler struct {
    pushService  PushService
    emailService EmailService
}

func (h *NotificationEventHandler) HandleOutingBooked(event OutingBooked) error {
    // Notifier l'utilisateur de sa réservation
    return h.pushService.Send(event.UserID, PushNotification{
        Title: "Réservation confirmée",
        Body:  "Votre sortie a été réservée. Présentez le QR code.",
    })
}

func (h *NotificationEventHandler) HandleOutingCheckedIn(event OutingCheckedIn) error {
    // Notifier le partenaire
    return h.pushService.SendToPartner(event.PartnerID, PushNotification{
        Title: "Check-in effectué",
        Body:  "Un client vient de valider sa réservation.",
    })
}
```

---

## Structure des Services

### Template de Service DDD

```
services/{context}-service/
├── cmd/
│   └── main.go                          # Entrypoint
│
├── internal/
│   ├── domain/                          # 🎯 DOMAIN LAYER (Pure, no dependencies)
│   │   ├── aggregate/
│   │   │   └── {aggregate}.go           # Aggregate Root
│   │   ├── entity/
│   │   │   └── {entity}.go              # Entities
│   │   ├── valueobject/
│   │   │   └── {vo}.go                  # Value Objects
│   │   ├── event/
│   │   │   └── events.go                # Domain Events
│   │   ├── repository/
│   │   │   └── {aggregate}_repository.go # Repository Interface (Port)
│   │   ├── service/
│   │   │   └── domain_service.go        # Domain Services
│   │   └── error/
│   │       └── errors.go                # Domain Errors
│   │
│   ├── application/                     # 📦 APPLICATION LAYER
│   │   ├── command/
│   │   │   ├── handler.go               # Command Handlers
│   │   │   └── commands.go              # Command DTOs
│   │   ├── query/
│   │   │   ├── handler.go               # Query Handlers
│   │   │   └── queries.go               # Query DTOs
│   │   ├── service/
│   │   │   └── application_service.go   # Orchestration
│   │   └── dto/
│   │       └── responses.go             # Response DTOs
│   │
│   ├── infrastructure/                  # 🔧 INFRASTRUCTURE LAYER
│   │   ├── persistence/
│   │   │   └── mongodb/
│   │   │       ├── repository_impl.go   # Repository Implementation (Adapter)
│   │   │       └── mapper.go            # Domain <-> MongoDB mapping
│   │   ├── messaging/
│   │   │   └── nats/
│   │   │       ├── publisher.go         # Event Publisher
│   │   │       └── subscriber.go        # Event Subscriber
│   │   ├── cache/
│   │   │   └── redis/
│   │   │       └── cache.go             # Cache Adapter
│   │   └── external/
│   │       └── {service}_client.go      # External API Clients
│   │
│   └── interface/                       # 🌐 INTERFACE LAYER
│       ├── grpc/
│       │   ├── server.go                # gRPC Server
│       │   ├── handler.go               # gRPC Handlers
│       │   └── mapper.go                # Proto <-> Domain mapping
│       └── http/                        # (si nécessaire)
│           └── handler.go
│
├── proto/
│   └── {service}.proto                  # gRPC definitions
│
├── config/
│   └── config.go
│
├── Dockerfile
├── go.mod
└── go.sum
```

### Exemple Concret: Booking Service

```
services/booking-service/
├── cmd/
│   └── main.go
│
├── internal/
│   ├── domain/
│   │   ├── aggregate/
│   │   │   └── outing.go               # Outing Aggregate Root
│   │   ├── valueobject/
│   │   │   ├── qrcode.go
│   │   │   └── offer_snapshot.go
│   │   ├── event/
│   │   │   └── outing_events.go
│   │   ├── repository/
│   │   │   └── outing_repository.go    # interface
│   │   └── error/
│   │       └── errors.go
│   │
│   ├── application/
│   │   ├── command/
│   │   │   ├── book_outing.go
│   │   │   ├── checkin_outing.go
│   │   │   └── cancel_outing.go
│   │   ├── query/
│   │   │   ├── get_outing.go
│   │   │   ├── list_user_outings.go
│   │   │   └── get_outing_by_qr.go
│   │   └── service/
│   │       └── booking_service.go
│   │
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   └── mongodb/
│   │   │       └── outing_repository.go
│   │   ├── messaging/
│   │   │   └── nats/
│   │   │       ├── publisher.go
│   │   │       └── offer_subscriber.go  # Écoute OfferExpired
│   │   └── grpc/
│   │       └── discovery_client.go      # Pour récupérer OfferSnapshot
│   │
│   └── interface/
│       └── grpc/
│           ├── server.go
│           └── booking_handler.go
│
├── proto/
│   └── booking.proto
│
└── Dockerfile
```

---

## Anti-Corruption Layers

### Cross-Context Communication

```go
// ============================================
// ACL: Booking -> Discovery
// ============================================

package acl

// Dans Booking Service, on ne dépend PAS du domain Discovery
// On utilise un ACL pour traduire

type DiscoveryACL interface {
    GetOfferSnapshot(ctx context.Context, offerID string) (OfferSnapshot, error)
}

type DiscoveryACLImpl struct {
    client discoverygrpc.DiscoveryServiceClient
}

func (a *DiscoveryACLImpl) GetOfferSnapshot(ctx context.Context, offerID string) (OfferSnapshot, error) {
    // Appel gRPC vers Discovery Service
    resp, err := a.client.GetOffer(ctx, &discoverygrpc.GetOfferRequest{
        OfferId: offerID,
    })
    if err != nil {
        return OfferSnapshot{}, err
    }
    
    // TRANSLATION: Proto -> Domain Value Object
    return OfferSnapshot{
        ID:              OfferID(resp.Id),
        PartnerID:       PartnerID(resp.PartnerId),
        EstablishmentID: EstablishmentID(resp.EstablishmentId),
        Title:           resp.Title,
        Description:     resp.Description,
        Discount:        mapDiscount(resp.Discount),
        Category:        resp.Category,
        Location:        GeoLocation{
            Longitude: resp.Location.Longitude,
            Latitude:  resp.Location.Latitude,
        },
        CapturedAt:      time.Now(),
    }, nil
}

// ============================================
// ACL: Notification (Generic) <- All Contexts
// ============================================

// Le Notification Service définit son propre modèle
// Il ne connaît PAS les domain objects des autres contexts

package notification

type NotificationRequest struct {
    RecipientType string   // user, partner, admin
    RecipientID   string
    Channel       string   // push, email, sms
    Template      string
    Data          map[string]interface{}
}

// Event Handler avec ACL
type EventHandler struct {
    service NotificationService
}

func (h *EventHandler) HandleOutingBooked(data []byte) error {
    // Parse event générique
    var event struct {
        OutingID  string    `json:"outing_id"`
        UserID    string    `json:"user_id"`
        OfferID   string    `json:"offer_id"`
        Timestamp time.Time `json:"timestamp"`
    }
    json.Unmarshal(data, &event)
    
    // TRANSLATION vers le domain Notification
    return h.service.Send(NotificationRequest{
        RecipientType: "user",
        RecipientID:   event.UserID,
        Channel:       "push",
        Template:      "booking_confirmed",
        Data: map[string]interface{}{
            "outing_id": event.OutingID,
            "offer_id":  event.OfferID,
        },
    })
}
```

---

## Résumé

### Principes DDD Appliqués

| Principe | Application |
|----------|-------------|
| **Bounded Contexts** | 6 contextes indépendants |
| **Ubiquitous Language** | Glossaire partagé par équipe |
| **Aggregates** | User, Partner, Offer, Outing |
| **Value Objects** | Email, Money, GeoLocation, QRCode |
| **Domain Events** | Communication asynchrone NATS |
| **Repository Pattern** | Interface dans domain, impl dans infra |
| **ACL** | Protection des frontières de contexte |

### Stack Technique

| Composant | Technologie |
|-----------|-------------|
| Langage | Go 1.21+ |
| API Sync | gRPC + protobuf |
| API Async | NATS JetStream |
| Database | MongoDB (par context) |
| Cache | Redis |
| Observability | OpenTelemetry + Jaeger + Prometheus + Loki + Grafana |

---

*Document généré pour Yousoon - Architecture DDD v1.0*
