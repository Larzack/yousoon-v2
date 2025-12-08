# 📊 Modèle de Données MongoDB - Yousoon

> Modèle de données organisé par Bounded Context (DDD)  
> Base de données : MongoDB (1 cluster, 1 database par context)  
> Région : Europe (Irlande) - Conformité RGPD  
> **Figma** : [Yousoon-Test2](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2?node-id=121-114)  
> **Dernière mise à jour** : 9 décembre 2025

---

## 📋 Table des Matières

1. [Vue d'Ensemble par Microservice](#vue-densemble-par-microservice)
2. [MCD Global](#mcd-global)
3. [Identity Service](#1-identity-service)
4. [Partner Service](#2-partner-service)
5. [Discovery Service](#3-discovery-service)
6. [Booking Service](#4-booking-service)
7. [Engagement Service](#5-engagement-service)
8. [Notification Service](#6-notification-service)
9. [Relations Cross-Context](#relations-cross-context)
10. [Conventions & Best Practices](#conventions--best-practices)

---

## Vue d'Ensemble par Microservice

### Architecture DDD - 6 Bounded Contexts

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           YOUSOON DATABASES                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │ identity_db     │  │  partner_db     │  │ discovery_db    │              │
│  │ (Core Domain)   │  │  (Core Domain)  │  │ (Core Domain)   │              │
│  │                 │  │                 │  │                 │              │
│  │ • users         │  │ • partners      │  │ • offers        │              │
│  │ • subscriptions │  │ • establishments│  │ • categories    │              │
│  │ • sub_plans     │  │ • team_members  │  │                 │              │
│  │ • id_verif      │  │ • invitations   │  │                 │              │
│  │ • user_grades   │  │ • partner_stats │  │                 │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │  booking_db     │  │ engagement_db   │  │ notification_db │              │
│  │ (Core Domain)   │  │ (Supporting)    │  │ (Generic)       │              │
│  │                 │  │                 │  │                 │              │
│  │ • outings       │  │ • favorites     │  │ • notifications │              │
│  │ • qr_codes      │  │ • reviews       │  │ • templates     │              │
│  │                 │  │ • conversations │  │ • push_tokens   │              │
│  │                 │  │ • messages      │  │ • admin_logs    │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Collections par Service

| Service | Database | Collections | Volume estimé |
|---------|----------|-------------|---------------|
| **Identity** | `identity_db` | users, subscriptions, subscription_plans, identity_verifications, user_grades | ~200k docs |
| **Partner** | `partner_db` | partners, establishments, team_members, invitations, partner_stats | ~10k docs |
| **Discovery** | `discovery_db` | offers, categories | ~15k docs |
| **Booking** | `booking_db` | outings, qr_codes | ~500k/an |
| **Engagement** | `engagement_db` | favorites, reviews, conversations, messages | ~1M docs |
| **Notification** | `notification_db` | notifications, templates, push_tokens, admin_logs | ~1M/an |

### Conventions Globales

- **_id** : ObjectId MongoDB
- **Timestamps** : `createdAt`, `updatedAt` (ISO 8601)
- **Soft delete** : `deletedAt` (null si actif)
- **Références cross-context** : Stockage de l'ID uniquement (pas de dénormalisation cross-service)
- **Références intra-context** : Dénormalisation autorisée pour performance

---

## MCD Global

### Diagramme Entité-Relation Complet

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    YOUSOON - MCD GLOBAL                                      │
└─────────────────────────────────────────────────────────────────────────────────────────────┘

                                    ┌─────────────────┐
                                    │  SUBSCRIPTION   │
                                    │     PLANS       │
                                    └────────┬────────┘
                                             │ 1
                                             │
                                             │ N
┌─────────────────┐                ┌─────────▼─────────┐                ┌─────────────────┐
│   IDENTITY      │                │                   │                │   USER          │
│ VERIFICATION    │ N          1   │   SUBSCRIPTIONS   │   1         N  │   GRADES        │
└────────┬────────┘◄───────────────┤                   ├───────────────►└─────────────────┘
         │                         └─────────┬─────────┘
         │ 1                                 │ N
         │                                   │
         │ N                                 │ 1
┌────────▼────────────────────────────────────▼────────────────────────────────────────────┐
│                                        USERS                                              │
│  (Aggregate Root - Identity Context)                                                     │
│  - email, passwordHash, phone                                                            │
│  - profile (firstName, lastName, avatar, birthDate, gender)                              │
│  - preferences (language, notifications, categories, maxDistance)                        │
│  - lastLocation (GeoJSON Point)                                                          │
│  - grade (explorateur, aventurier, grand_voyageur, conquerant)                          │
│  - status (active, suspended, deleted)                                                   │
└────────┬───────────────────────┬───────────────────────┬─────────────────────────────────┘
         │                       │                       │
         │ 1                     │ 1                     │ 1
         │                       │                       │
         │ N                     │ N                     │ N
┌────────▼────────┐    ┌─────────▼─────────┐   ┌────────▼────────┐
│    OUTINGS      │    │     FAVORITES     │   │    REVIEWS      │
│  (Réservations) │    │                   │   │                 │
└────────┬────────┘    └─────────┬─────────┘   └────────┬────────┘
         │                       │                       │
         │ N                     │ N                     │ N
         │                       │                       │
         │ 1                     │ 1                     │ 1
┌────────▼───────────────────────▼───────────────────────▼────────────────────────────────┐
│                                       OFFERS                                             │
│  (Aggregate Root - Discovery Context)                                                    │
│  - title, description, shortDescription                                                  │
│  - discount (type, value, originalPrice, formula)                                        │
│  - validity (startDate, endDate, timezone)                                               │
│  - schedule (allDay, slots[])                                                            │
│  - quota (total, perUser, perDay, used)                                                  │
│  - status (draft, pending, active, paused, expired, archived)                            │
└────────┬───────────────────────────────────────────────┬────────────────────────────────┘
         │                                               │
         │ N                                             │ N
         │                                               │
         │ 1                                             │ 1
┌────────▼────────┐                            ┌─────────▼─────────┐
│ ESTABLISHMENTS  │                            │    CATEGORIES     │
│                 │ N                     1    │                   │
└────────┬────────┘◄───────────────────────────┤ - name (fr, en)   │
         │                                     │ - slug, icon      │
         │ N                                   │ - parent (self)   │
         │                                     └───────────────────┘
         │ 1
┌────────▼────────────────────────────────────────────────────────────────────────────────┐
│                                      PARTNERS                                            │
│  (Aggregate Root - Partner Context)                                                      │
│  - company (name, tradeName, siret, vatNumber, legalForm)                               │
│  - branding (logo, coverImage, primaryColor, description)                                │
│  - contact (firstName, lastName, email, phone, role)                                     │
│  - category, subcategories[]                                                             │
│  - status (pending, active, suspended)                                                   │
└────────┬───────────────────────────────────────────────┬────────────────────────────────┘
         │                                               │
         │ 1                                             │ 1
         │                                               │
         │ N                                             │ N
┌────────▼────────┐                            ┌─────────▼─────────┐
│  TEAM_MEMBERS   │                            │    INVITATIONS    │
│  (Équipe)       │                            │  (En attente)     │
└─────────────────┘                            └───────────────────┘


                    ┌─────────────────────────────────────────────────┐
                    │              MESSAGING SYSTEM                    │
                    │                                                  │
                    │  ┌─────────────────┐    ┌─────────────────┐     │
                    │  │  CONVERSATIONS  │ 1  │    MESSAGES     │     │
                    │  │                 │◄───┤                 │     │
                    │  │ - participants  │  N │ - senderId      │     │
                    │  │ - type          │    │ - content       │     │
                    │  └────────┬────────┘    │ - readAt        │     │
                    │           │ N           └─────────────────┘     │
                    │           │                                      │
                    │           │ N                                    │
                    │     ┌─────▼─────┐                                │
                    │     │   USERS   │                                │
                    │     └───────────┘                                │
                    └─────────────────────────────────────────────────┘


                    ┌─────────────────────────────────────────────────┐
                    │            NOTIFICATION SYSTEM                   │
                    │                                                  │
                    │  ┌─────────────────┐    ┌─────────────────┐     │
                    │  │  PUSH_TOKENS    │    │   TEMPLATES     │     │
                    │  │                 │    │                 │     │
                    │  │ - userId        │    │ - code          │     │
                    │  │ - token         │    │ - channel       │     │
                    │  │ - platform      │    │ - content       │     │
                    │  └────────┬────────┘    └────────┬────────┘     │
                    │           │ N                    │ 1            │
                    │           │                      │              │
                    │           │ 1                    │ N            │
                    │     ┌─────▼──────────────────────▼─────┐        │
                    │     │          NOTIFICATIONS           │        │
                    │     │  - userId, type, channel         │        │
                    │     │  - content, status               │        │
                    │     │  - sentAt, readAt                │        │
                    │     └──────────────────────────────────┘        │
                    └─────────────────────────────────────────────────┘


                    ┌─────────────────────────────────────────────────┐
                    │              ADMIN / AUDIT SYSTEM                │
                    │                                                  │
                    │  ┌──────────────────────────────────────┐       │
                    │  │             ADMIN_LOGS               │       │
                    │  │  - adminId, adminEmail               │       │
                    │  │  - action, resource, resourceId      │       │
                    │  │  - details (before, after, reason)   │       │
                    │  │  - ip, userAgent                     │       │
                    │  └──────────────────────────────────────┘       │
                    └─────────────────────────────────────────────────┘
```

---

## Schémas Détaillés

### 1. Users

```javascript
{
  _id: ObjectId,
  
  // Authentification
  email: String,                      // unique, indexed
  passwordHash: String,               // bcrypt
  phone: String,                      // E.164 format, optional
  
  // Profil
  profile: {
    firstName: String,
    lastName: String,
    displayName: String,              // computed: firstName + lastName
    avatar: String,                   // URL Cloudinary/S3
    birthDate: Date,
    gender: String,                   // 'male', 'female', 'other', null
  },
  
  // Vérification identité
  identity: {
    status: String,                   // 'not_submitted', 'pending', 'verified', 'rejected'
    verificationId: String,           // Référence Onfido/Veriff
    verifiedAt: Date,
    documentType: String,             // 'cni', 'passport', 'permit'
  },
  
  // Préférences
  preferences: {
    language: String,                 // 'fr', 'en'
    notifications: {
      push: Boolean,
      email: Boolean,
      sms: Boolean,
      marketing: Boolean,
    },
    categories: [ObjectId],           // Catégories préférées
    maxDistance: Number,              // km, default 10
  },
  
  // Géolocalisation (dernière connue)
  lastLocation: {
    type: "Point",
    coordinates: [Number, Number],    // [longitude, latitude]
    updatedAt: Date,
  },
  
  // Favoris (dénormalisés pour performance)
  favorites: [{
    offerId: ObjectId,
    addedAt: Date,
  }],
  
  // Tokens
  fcmTokens: [{                       // Firebase Cloud Messaging
    token: String,
    platform: String,                 // 'ios', 'android'
    addedAt: Date,
  }],
  
  // Social login
  socialAccounts: [{
    provider: String,                 // 'google', 'apple', 'facebook'
    providerId: String,
    email: String,
  }],
  
  // Statut
  status: String,                     // 'active', 'suspended', 'deleted'
  emailVerified: Boolean,
  phoneVerified: Boolean,
  
  // Metadata
  createdAt: Date,
  updatedAt: Date,
  lastLoginAt: Date,
  deletedAt: Date,                    // Soft delete
}

// Indexes
db.users.createIndex({ email: 1 }, { unique: true })
db.users.createIndex({ phone: 1 }, { sparse: true })
db.users.createIndex({ "lastLocation": "2dsphere" })
db.users.createIndex({ "identity.status": 1 })
db.users.createIndex({ "favorites.offerId": 1 })
```

---

### 2. Partners

```javascript
{
  _id: ObjectId,
  
  // Informations entreprise
  company: {
    name: String,                     // Raison sociale
    tradeName: String,                // Nom commercial
    siret: String,                    // unique
    vatNumber: String,                // TVA intracommunautaire
    legalForm: String,                // SARL, SAS, etc.
  },
  
  // Branding
  branding: {
    logo: String,                     // URL
    coverImage: String,               // URL
    primaryColor: String,             // Hex
    description: String,              // Rich text
  },
  
  // Contact principal
  contact: {
    firstName: String,
    lastName: String,
    email: String,
    phone: String,
    role: String,
  },
  
  // Catégorie
  category: String,                   // 'bar', 'restaurant', 'club', etc.
  subcategories: [String],
  
  // Équipe
  team: [{
    userId: ObjectId,                 // Lien vers users (optionnel)
    email: String,
    firstName: String,
    lastName: String,
    role: String,                     // 'admin', 'manager', 'staff', 'viewer'
    invitedAt: Date,
    joinedAt: Date,
    status: String,                   // 'pending', 'active', 'inactive'
  }],
  
  // Abonnement (partenaires gratuits actuellement)
  subscription: {
    plan: String,                     // 'free', 'pro', 'enterprise'
    status: String,                   // 'active', 'past_due', 'cancelled'
    currentPeriodEnd: Date,
  },
  
  // Statistiques (dénormalisées, mises à jour périodiquement)
  stats: {
    totalOffers: Number,
    activeOffers: Number,
    totalBookings: Number,
    totalCheckins: Number,
    avgRating: Number,
    reviewCount: Number,
    lastUpdated: Date,
  },
  
  // Statut
  status: String,                     // 'pending', 'active', 'suspended'
  verifiedAt: Date,
  
  // Metadata
  createdAt: Date,
  updatedAt: Date,
  deletedAt: Date,
}

// Indexes
db.partners.createIndex({ "company.siret": 1 }, { unique: true })
db.partners.createIndex({ category: 1 })
db.partners.createIndex({ status: 1 })
db.partners.createIndex({ "team.email": 1 })
```

---

### 3. Establishments

```javascript
{
  _id: ObjectId,
  partnerId: ObjectId,                // Référence Partner
  
  // Informations
  name: String,
  description: String,
  
  // Adresse
  address: {
    street: String,
    streetNumber: String,
    complement: String,
    postalCode: String,
    city: String,
    country: String,                  // ISO 3166-1 alpha-2
    formatted: String,                // Adresse complète formatée
  },
  
  // Géolocalisation
  location: {
    type: "Point",
    coordinates: [Number, Number],    // [longitude, latitude]
  },
  
  // Contact
  contact: {
    phone: String,
    email: String,
    website: String,
  },
  
  // Horaires
  openingHours: [{
    dayOfWeek: Number,                // 0 = Dimanche, 1 = Lundi, ...
    open: String,                     // "09:00"
    close: String,                    // "23:00"
    isClosed: Boolean,
  }],
  
  // Jours fériés / Fermetures exceptionnelles
  closures: [{
    date: Date,
    reason: String,
  }],
  
  // Médias
  images: [{
    url: String,
    alt: String,
    isPrimary: Boolean,
    order: Number,
  }],
  
  // Caractéristiques
  features: [String],                 // ['terrasse', 'wifi', 'parking', 'handicap']
  
  // Catégorie spécifique
  type: String,                       // Plus précis que la catégorie partner
  priceRange: Number,                 // 1-4 (€ à €€€€)
  
  // Statut
  isActive: Boolean,
  
  // Metadata
  createdAt: Date,
  updatedAt: Date,
}

// Indexes
db.establishments.createIndex({ partnerId: 1 })
db.establishments.createIndex({ location: "2dsphere" })
db.establishments.createIndex({ "address.city": 1 })
db.establishments.createIndex({ isActive: 1, location: "2dsphere" })
```

---

### 4. Offers

```javascript
{
  _id: ObjectId,
  partnerId: ObjectId,                // Référence Partner
  establishmentId: ObjectId,          // Référence Establishment
  
  // Informations principales
  title: String,
  description: String,                // Rich text
  shortDescription: String,           // Max 100 chars
  
  // Catégorie
  categoryId: ObjectId,               // Référence Category
  tags: [String],
  
  // Réduction
  discount: {
    type: String,                     // 'percentage', 'fixed', 'formula'
    value: Number,                    // 20 pour 20% ou 5 pour 5€
    originalPrice: Number,            // Prix original (optionnel)
    formula: String,                  // "1 acheté = 1 offert"
  },
  
  // Conditions
  conditions: [{
    type: String,                     // 'min_purchase', 'min_people', 'first_visit'
    value: Mixed,
    label: String,
  }],
  termsAndConditions: String,
  
  // Validité temporelle
  validity: {
    startDate: Date,
    endDate: Date,
    timezone: String,                 // 'Europe/Paris'
  },
  
  // Planning hebdomadaire
  schedule: {
    allDay: Boolean,
    slots: [{
      dayOfWeek: Number,              // 0 = Dimanche
      startTime: String,              // "17:00"
      endTime: String,                // "20:00"
    }],
  },
  
  // Quotas
  quota: {
    total: Number,                    // Limite globale (null = illimité)
    perUser: Number,                  // Par utilisateur (null = illimité)
    perDay: Number,                   // Par jour (null = illimité)
    used: Number,                     // Compteur utilisations
  },
  
  // Médias
  images: [{
    url: String,
    alt: String,
    isPrimary: Boolean,
    order: Number,
  }],
  
  // Données dénormalisées (performance)
  _partner: {
    name: String,
    logo: String,
    category: String,
  },
  _establishment: {
    name: String,
    address: String,
    city: String,
    location: {
      type: "Point",
      coordinates: [Number, Number],
    },
  },
  
  // Statistiques
  stats: {
    views: Number,
    bookings: Number,
    checkins: Number,
    favorites: Number,
  },
  
  // Statut
  status: String,                     // 'draft', 'pending', 'active', 'paused', 'expired', 'archived'
  isActive: Boolean,                  // Computed: status === 'active' && now dans validity
  
  // Modération
  moderation: {
    status: String,                   // 'pending', 'approved', 'rejected'
    reviewedBy: ObjectId,
    reviewedAt: Date,
    comment: String,
  },
  
  // Metadata
  createdAt: Date,
  updatedAt: Date,
  publishedAt: Date,
  deletedAt: Date,
}

// Indexes
db.offers.createIndex({ partnerId: 1 })
db.offers.createIndex({ establishmentId: 1 })
db.offers.createIndex({ categoryId: 1 })
db.offers.createIndex({ status: 1, isActive: 1 })
db.offers.createIndex({ "_establishment.location": "2dsphere" })
db.offers.createIndex({ 
  "_establishment.location": "2dsphere",
  isActive: 1,
  categoryId: 1 
})
db.offers.createIndex({ title: "text", description: "text" }, { default_language: "french" })
db.offers.createIndex({ "validity.startDate": 1, "validity.endDate": 1 })
db.offers.createIndex({ "discount.type": 1, "discount.value": -1 })
```

---

### 5. Bookings

```javascript
{
  _id: ObjectId,
  
  // Références
  userId: ObjectId,
  offerId: ObjectId,
  partnerId: ObjectId,                // Dénormalisé pour queries partenaire
  establishmentId: ObjectId,
  
  // QR Code
  qrCode: {
    code: String,                     // UUID unique
    data: String,                     // Données encodées
    expiresAt: Date,
  },
  
  // Statut
  status: String,                     // 'pending', 'confirmed', 'checked_in', 'cancelled', 'expired', 'no_show'
  
  // Timeline
  timeline: [{
    status: String,
    timestamp: Date,
    actor: String,                    // 'user', 'partner', 'system'
    metadata: Object,
  }],
  
  // Check-in
  checkin: {
    checkedInAt: Date,
    checkedInBy: ObjectId,            // userId du staff partenaire
    method: String,                   // 'qr_scan', 'manual'
    location: {
      type: "Point",
      coordinates: [Number, Number],
    },
  },
  
  // Annulation
  cancellation: {
    cancelledAt: Date,
    cancelledBy: String,              // 'user', 'partner', 'system'
    reason: String,
  },
  
  // Données dénormalisées (snapshot au moment de la réservation)
  _offer: {
    title: String,
    discount: Object,
    images: [String],
  },
  _partner: {
    name: String,
    logo: String,
  },
  _establishment: {
    name: String,
    address: String,
  },
  _user: {
    firstName: String,
    lastName: String,
    email: String,
  },
  
  // Metadata
  createdAt: Date,
  updatedAt: Date,
  expiresAt: Date,                    // Auto-expiration
}

// Indexes
db.bookings.createIndex({ userId: 1, createdAt: -1 })
db.bookings.createIndex({ offerId: 1 })
db.bookings.createIndex({ partnerId: 1, createdAt: -1 })
db.bookings.createIndex({ establishmentId: 1 })
db.bookings.createIndex({ "qrCode.code": 1 }, { unique: true })
db.bookings.createIndex({ status: 1 })
db.bookings.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0 }) // TTL index
```

---

### 6. Categories

```javascript
{
  _id: ObjectId,
  
  name: {
    fr: String,
    en: String,
  },
  slug: String,                       // unique, URL-friendly
  description: {
    fr: String,
    en: String,
  },
  
  icon: String,                       // Nom icône ou URL
  color: String,                      // Hex
  image: String,                      // URL
  
  parent: ObjectId,                   // null si racine
  order: Number,                      // Ordre d'affichage
  
  isActive: Boolean,
  
  createdAt: Date,
  updatedAt: Date,
}

// Indexes
db.categories.createIndex({ slug: 1 }, { unique: true })
db.categories.createIndex({ parent: 1, order: 1 })
```

---

### 7. Identity Verifications

```javascript
{
  _id: ObjectId,
  userId: ObjectId,
  
  // Provider externe
  provider: String,                   // 'onfido', 'veriff', 'jumio'
  externalId: String,                 // ID chez le provider
  
  // Document
  document: {
    type: String,                     // 'cni', 'passport', 'driving_license'
    country: String,                  // ISO 3166-1 alpha-2
    frontImageUrl: String,            // URL sécurisée (temporaire)
    backImageUrl: String,
  },
  
  // Selfie
  selfie: {
    imageUrl: String,
    livenessScore: Number,
  },
  
  // Résultat
  result: {
    status: String,                   // 'pending', 'verified', 'rejected'
    confidence: Number,               // 0-100
    checks: [{
      name: String,
      status: String,
      details: Object,
    }],
    extractedData: {
      firstName: String,
      lastName: String,
      birthDate: Date,
      documentNumber: String,
      expiryDate: Date,
    },
    rejectionReasons: [String],
  },
  
  // Webhook
  webhookReceived: Boolean,
  webhookReceivedAt: Date,
  rawWebhookData: Object,
  
  // Metadata
  createdAt: Date,
  updatedAt: Date,
  completedAt: Date,
}

// Indexes
db.identity_verifications.createIndex({ userId: 1 })
db.identity_verifications.createIndex({ externalId: 1 })
db.identity_verifications.createIndex({ "result.status": 1 })
```

---

### 8. Notifications

```javascript
{
  _id: ObjectId,
  userId: ObjectId,
  
  type: String,                       // 'booking_confirmed', 'offer_nearby', 'reminder'
  channel: String,                    // 'push', 'email', 'sms'
  
  // Contenu
  content: {
    title: String,
    body: String,
    image: String,
    data: Object,                     // Deep link, metadata
  },
  
  // Envoi
  status: String,                     // 'pending', 'sent', 'delivered', 'failed', 'read'
  sentAt: Date,
  deliveredAt: Date,
  readAt: Date,
  error: String,
  
  // Références
  relatedTo: {
    type: String,                     // 'offer', 'booking', 'partner'
    id: ObjectId,
  },
  
  createdAt: Date,
}

// Indexes
db.notifications.createIndex({ userId: 1, createdAt: -1 })
db.notifications.createIndex({ status: 1 })
db.notifications.createIndex({ createdAt: 1 }, { expireAfterSeconds: 7776000 }) // TTL 90 jours
```

---

### 9. Subscription Plans (Plans d'abonnement)

```javascript
{
  _id: ObjectId,
  
  // Identifiant
  code: String,                       // 'free', 'monthly', 'yearly', 'premium'
  
  // Nom et description
  name: {
    fr: String,                       // "Mensuel"
    en: String,                       // "Monthly"
  },
  description: {
    fr: String,
    en: String,
  },
  
  // Tarification
  pricing: {
    amount: Number,                   // En centimes (990 = 9.90€)
    currency: String,                 // 'EUR'
    interval: String,                 // 'month', 'year', 'lifetime'
    intervalCount: Number,            // 1, 3, 12...
  },
  
  // Période d'essai
  trial: {
    enabled: Boolean,
    durationDays: Number,             // 30 jours par défaut
  },
  
  // Fonctionnalités incluses
  features: [{
    code: String,                     // 'unlimited_bookings', 'priority_support'
    name: { fr: String, en: String },
    included: Boolean,
    limit: Number,                    // null = illimité
  }],
  
  // Limites
  limits: {
    bookingsPerMonth: Number,         // null = illimité
    favoritesMax: Number,
  },
  
  // Affichage
  display: {
    order: Number,                    // Ordre d'affichage
    highlighted: Boolean,             // "Recommandé"
    badge: String,                    // "Populaire", "Meilleur rapport"
    color: String,                    // Hex
  },
  
  // Stripe
  stripeProductId: String,
  stripePriceId: String,
  
  // Statut
  isActive: Boolean,
  
  createdAt: Date,
  updatedAt: Date,
}

// Indexes
db.subscription_plans.createIndex({ code: 1 }, { unique: true })
db.subscription_plans.createIndex({ isActive: 1, "display.order": 1 })
```

---

### 10. Subscriptions (Abonnements utilisateurs)

```javascript
{
  _id: ObjectId,
  userId: ObjectId,
  planId: ObjectId,
  
  // In-App Purchase (Apple/Google)
  inAppPurchase: {
    platform: String,                 // 'apple', 'google'
    productId: String,                // ID produit in-app
    transactionId: String,            // ID transaction
    receipt: String,                  // Reçu pour validation
    receiptValidatedAt: Date,
  },
  
  // Statut
  status: String,                     // 'trialing', 'active', 'past_due', 'cancelled', 'expired'
  
  // Période d'essai
  trial: {
    startDate: Date,
    endDate: Date,                    // Date de fin essai (configurable, défaut 30j)
    converted: Boolean,               // A converti en payant ?
  },
  
  // Période courante
  currentPeriod: {
    startDate: Date,
    endDate: Date,
  },
  
  // Annulation
  cancellation: {
    requestedAt: Date,
    reason: String,
    effectiveAt: Date,                // Fin de la période payée
    feedback: String,
  },
  
  // Historique paiements (dénormalisé pour affichage rapide)
  lastPayment: {
    amount: Number,
    currency: String,
    date: Date,
    status: String,
  },
  
  // Données du plan au moment de la souscription (snapshot)
  _plan: {
    code: String,
    name: String,
    amount: Number,
    interval: String,
  },
  
  // Metadata
  createdAt: Date,
  updatedAt: Date,
}

// Indexes
db.subscriptions.createIndex({ userId: 1 })
db.subscriptions.createIndex({ "inAppPurchase.transactionId": 1 })
db.subscriptions.createIndex({ status: 1 })
db.subscriptions.createIndex({ "trial.endDate": 1 })
db.subscriptions.createIndex({ "currentPeriod.endDate": 1 })
```

---

### 11. Reviews (Avis)

```javascript
{
  _id: ObjectId,
  
  // Références
  userId: ObjectId,
  offerId: ObjectId,                  // Avis sur une offre
  partnerId: ObjectId,                // Avis sur un partenaire (dénormalisé)
  establishmentId: ObjectId,          // Établissement concerné
  bookingId: ObjectId,                // Réservation associée (optionnel)
  
  // Note
  rating: Number,                     // 1-5
  
  // Contenu
  title: String,                      // Optionnel
  content: String,                    // Texte de l'avis
  
  // Médias (optionnel)
  images: [String],                   // URLs photos
  
  // NOTE: Pas de réponse partenaire (désactivé)
  // response: { ... } - Non implémenté
  
  // Modération
  moderation: {
    status: String,                   // 'pending', 'approved', 'rejected', 'reported'
    reports: [{
      userId: ObjectId,
      reason: String,
      reportedAt: Date,
    }],
    reviewedBy: ObjectId,             // Admin
    reviewedAt: Date,
    rejectReason: String,
  },
  
  // Données dénormalisées
  _user: {
    firstName: String,
    avatar: String,
  },
  _offer: {
    title: String,
  },
  _partner: {
    name: String,
  },
  
  // Statistiques
  helpfulCount: Number,               // Nombre de "utile"
  
  // Metadata
  isVerifiedPurchase: Boolean,        // A réellement utilisé l'offre
  createdAt: Date,
  updatedAt: Date,
}

// Indexes
db.reviews.createIndex({ offerId: 1, createdAt: -1 })
db.reviews.createIndex({ partnerId: 1, createdAt: -1 })
db.reviews.createIndex({ userId: 1 })
db.reviews.createIndex({ "moderation.status": 1 })
db.reviews.createIndex({ rating: 1 })
```

---

### 12. Admin Logs (Audit)

```javascript
{
  _id: ObjectId,
  
  // Admin qui a fait l'action
  adminId: ObjectId,
  adminEmail: String,
  
  // Action
  action: String,                     // 'validate_partner', 'reject_identity', 'delete_review'
  resource: String,                   // 'partner', 'user', 'offer', 'review', 'identity'
  resourceId: ObjectId,
  
  // Détails
  details: {
    before: Object,                   // État avant (optionnel)
    after: Object,                    // État après (optionnel)
    reason: String,                   // Motif si rejet
  },
  
  // Contexte
  ip: String,
  userAgent: String,
  
  createdAt: Date,
}

// Indexes
db.admin_logs.createIndex({ adminId: 1, createdAt: -1 })
db.admin_logs.createIndex({ resource: 1, resourceId: 1 })
db.admin_logs.createIndex({ action: 1 })
db.admin_logs.createIndex({ createdAt: 1 }, { expireAfterSeconds: 31536000 }) // TTL 1 an
```

---

## Relations

### Résumé des Relations

| Collection A | Collection B | Type | Champ |
|--------------|--------------|------|-------|
| Partners | Establishments | 1:N | `establishments.partnerId` |
| Partners | Offers | 1:N | `offers.partnerId` |
| Partners | Users (team) | M:N | `partners.team[].userId` |
| Establishments | Offers | 1:N | `offers.establishmentId` |
| Users | Offers (favorites) | M:N | `users.favorites[]` |
| Users | Bookings | 1:N | `bookings.userId` |
| Offers | Bookings | 1:N | `bookings.offerId` |
| Users | Identity Verifications | 1:N | `identity_verifications.userId` |
| Users | Notifications | 1:N | `notifications.userId` |
| Categories | Offers | 1:N | `offers.categoryId` |
| Users | Subscriptions | 1:N | `subscriptions.userId` |
| Subscription Plans | Subscriptions | 1:N | `subscriptions.planId` |
| Users | Reviews | 1:N | `reviews.userId` |
| Offers | Reviews | 1:N | `reviews.offerId` |
| Partners | Reviews | 1:N | `reviews.partnerId` |

### Dénormalisation

Pour la performance (< 50ms), certaines données sont dénormalisées :

1. **Offers** : Contient `_partner`, `_establishment` pour éviter les lookups
2. **Bookings** : Snapshot de l'offre au moment de la réservation
3. **Users.favorites** : Liste des IDs offres pour filtrage rapide

### Mise à jour des données dénormalisées

```javascript
// Trigger lors de la mise à jour d'un partner
db.offers.updateMany(
  { partnerId: partnerId },
  { $set: { "_partner.name": newName, "_partner.logo": newLogo } }
)
```

---

## RGPD - Suppression de Compte

### Workflow de Suppression

1. **Demande de suppression** → Période de grâce de **30 jours**
2. **Pendant la période de grâce** :
   - Compte désactivé (pas d'accès)
   - L'utilisateur peut annuler et récupérer son compte
   - Email de confirmation envoyé
3. **Après 30 jours** → **Suppression totale** :
   - Toutes les données personnelles supprimées
   - Réservations archivées (anonymisées)
   - Favoris, préférences, historique : supprimés
   - Pas de conservation de données

### Champs de suppression (users collection)

```javascript
{
  // ... autres champs
  
  deletion: {
    requestedAt: Date,              // Date de la demande
    scheduledAt: Date,              // requestedAt + 30 jours
    reason: String,                 // Raison optionnelle
    cancelledAt: Date,              // Si annulé pendant période de grâce
    completedAt: Date,              // Date de suppression effective
  }
}
```

### Job CRON quotidien

```go
// Tous les jours à 2h du matin
// 1. Sélectionner users WHERE deletion.scheduledAt <= NOW() AND deletion.completedAt IS NULL
// 2. Pour chaque user :
//    - Supprimer définitivement le document
//    - Anonymiser les bookings (userId → "deleted_user")
//    - Supprimer reviews associées
//    - Supprimer notifications
//    - Log dans admin_logs
```

---

## Bonnes Pratiques

### 1. Validation avec JSON Schema

```javascript
db.createCollection("users", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["email", "passwordHash", "profile", "createdAt"],
      properties: {
        email: {
          bsonType: "string",
          pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
        },
        // ...
      }
    }
  }
})
```

### 2. Transactions

Pour les opérations multi-documents critiques (création booking + update quota) :

```go
session, _ := client.StartSession()
defer session.EndSession(ctx)

session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
    // 1. Vérifier quota
    // 2. Créer booking
    // 3. Incrémenter quota.used
    return nil, nil
})
```

### 3. Change Streams

Pour les mises à jour temps réel :

```go
pipeline := mongo.Pipeline{
    bson.D{{Key: "$match", Value: bson.D{
        {Key: "operationType", Value: "insert"},
        {Key: "fullDocument.partnerId", Value: partnerId},
    }}},
}

stream, _ := offersCollection.Watch(ctx, pipeline)
for stream.Next(ctx) {
    // Notifier les clients
}
```

---

## 🔗 Références

- [MASTER_PROMPT.md](./MASTER_PROMPT.md)
- [Backend Architecture](./backend/ARCHITECTURE.md)
