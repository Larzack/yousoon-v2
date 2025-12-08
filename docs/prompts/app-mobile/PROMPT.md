# 📱 App Mobile Flutter - Prompt Détaillé

> **Module** : Application Mobile Yousoon  
> **Technologie** : Flutter (Dart)  
> **Cibles** : iOS + Android (versions récentes uniquement)  
> **Figma** : [Yousoon-Test2](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2?node-id=121-114)

---

## 🎯 Objectifs

L'application mobile Yousoon doit être :
- **Ultra réactive** : Sensation native, animations fluides 60fps
- **Fidèle au design** : Respect pixel-perfect du Figma (Dark Mode natif)
- **Performante** : Temps de réponse API < 50ms, ressenti < 100ms
- **Offline-capable** : Cache local + QR code réservation disponible hors-ligne
- **Temps réel** : GraphQL Subscriptions pour nouvelles offres et statuts
- **Partageable** : Deep links pour partage d'offres
- **Sécurisée** : Biométrie (Face ID / Touch ID) pour reconnexion
- **Multi-langue** : FR, EN + architecture extensible
- **Maintenable** : Architecture propre, composants réutilisables

---

## 🛠️ Stack Technique

### Framework & Language
| Technologie | Version | Justification |
|-------------|---------|---------------|
| Flutter | 3.16+ | Cross-platform performant |
| Dart | 3.2+ | Null safety, moderne |

### State Management
| Technologie | Justification |
|-------------|---------------|
| **Riverpod 2.x** ✅ | Type-safe, testable, compile-time safety, code generation |

### Communication API
| Technologie | Usage |
|-------------|-------|
| **graphql_flutter** | Client GraphQL |
| **ferry** | Code generation type-safe |
| **websocket_channel** | GraphQL Subscriptions temps réel |
| **dio** | HTTP client (upload images, etc.) |

### Cache Local
| Technologie | Usage |
|-------------|-------|
| **Hive** ou **Isar** | Stockage local NoSQL rapide |
| **shared_preferences** | Préférences simples |
| **flutter_secure_storage** | Tokens, données sensibles |

### Navigation
| Technologie | Usage |
|-------------|-------|
| **go_router** | Navigation déclarative |
| **auto_route** | Alternative avec code gen |

### UI/UX
| Technologie | Usage |
|-------------|-------|
| **flutter_animate** | Animations déclaratives |
| **cached_network_image** | Images avec cache |
| **shimmer** | Loading states |

### Tests
| Type | Technologie |
|------|-------------|
| Unit | flutter_test, mockito |
| Widget | flutter_test |
| Integration | integration_test |
| E2E | patrol |

### Services Externes
| Technologie | Usage |
|-------------|-------|
| **OneSignal** | Push notifications |
| **Amplitude** | Analytics |
| **Sentry** (self-hosted) | Crash reporting |
| **Google Maps** | Cartes et géolocalisation |
| **S3 + CloudFront** | CDN images/assets |
| **in_app_purchase** | Paiements 100% in-app (Apple/Google) |
| **local_auth** | Biométrie (Face ID / Touch ID) |
| **share_plus** | Partage natif |
| **uni_links** | Deep links |
| **mobile_scanner** | Scan QR code (check-in) |

### Internationalisation
| Technologie | Usage |
|-------------|-------|
| **flutter_localizations** | i18n native |
| **intl** | Formatage dates/nombres |
| **slang** ou **easy_localization** | Gestion traductions |

**Langues V1** : Français (FR), Anglais (EN)  
**Architecture** : Extensible pour ajout de langues

---

## 📲 Notifications Push

### Types de Notifications Activées
| Type | Description | Activé |
|------|-------------|--------|
| `offer_nearby` | Nouvelles offres à proximité | ✅ Oui |
| `booking_reminder` | Rappel de réservation | ✅ Oui |
| `marketing` | Offres promotionnelles | ✅ Oui |
| `offer_expiring` | Offres qui expirent bientôt | ❌ Non |
| `new_partner` | Nouveau partenaire inscrit | ❌ Non |

### Configuration
- **Provider** : OneSignal
- **Permission** : Demandée après inscription
- **Paramètres** : Configurable par l'utilisateur dans Settings

---

## 🔐 Sécurité

### Biométrie
- **Face ID / Touch ID** pour reconnexion
- Package : `local_auth`
- Opt-in lors du premier login

### Check-in
- **QR Code uniquement** (pas de geofencing)
- Package : `mobile_scanner`
- QR Code disponible hors-ligne dans le cache local

---

## 🎨 Thème

- **Dark Mode par défaut** (selon design Figma)
- Pas de switch light/dark (thème unique)

---

## 🏗️ Architecture

### Clean Architecture Adaptée

```
lib/
├── main.dart
├── app/
│   ├── app.dart
│   └── router.dart
├── core/
│   ├── constants/
│   ├── errors/
│   ├── network/
│   │   ├── graphql_client.dart
│   │   └── network_info.dart
│   ├── cache/
│   │   ├── cache_manager.dart
│   │   └── cache_policy.dart
│   └── utils/
├── features/
│   ├── auth/
│   │   ├── data/
│   │   │   ├── datasources/
│   │   │   ├── models/
│   │   │   └── repositories/
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   ├── repositories/
│   │   │   └── usecases/
│   │   └── presentation/
│   │       ├── providers/
│   │       ├── screens/
│   │       └── widgets/
│   ├── outings/
│   ├── offers/
│   ├── profile/
│   └── settings/
├── shared/
│   ├── widgets/                 # Composants réutilisables
│   │   ├── buttons/
│   │   ├── cards/
│   │   ├── inputs/
│   │   └── layouts/
│   └── theme/
│       ├── app_theme.dart
│       ├── app_colors.dart
│       └── app_typography.dart
└── l10n/                        # Internationalisation
```

---

## 🎨 Composants Réutilisables Identifiés

> À affiner après analyse du Figma

### Composants UI Génériques
- [ ] `YsButton` - Boutons (primary, secondary, outlined, text)
- [ ] `YsTextField` - Champs de saisie
- [ ] `YsCard` - Cartes génériques
- [ ] `YsBottomSheet` - Bottom sheets
- [ ] `YsChip` - Tags/chips
- [ ] `YsAvatar` - Avatars utilisateur/partenaire
- [ ] `YsRating` - Notation étoiles
- [ ] `YsLoader` - Indicateurs de chargement
- [ ] `YsEmptyState` - États vides
- [ ] `YsErrorState` - États d'erreur

### Composants Métier
- [ ] `OfferCard` - Carte d'offre
- [ ] `OutingCard` - Carte de sortie
- [ ] `PartnerCard` - Carte partenaire
- [ ] `CategoryChip` - Chip catégorie
- [ ] `DiscountBadge` - Badge réduction
- [ ] `BookingStatus` - Statut réservation
- [ ] `SearchBar` - Barre de recherche
- [ ] `FilterSheet` - Filtres
- [ ] `MapView` - Vue carte (Google Maps)

### Layouts
- [ ] `YsScaffold` - Scaffold personnalisé
- [ ] `YsAppBar` - AppBar personnalisée
- [ ] `YsBottomNav` - Navigation bottom
- [ ] `YsSlider` - Carousel/slider

---

## 📱 Écrans Principaux

> À compléter avec le Figma

### Authentification
1. Splash Screen
2. Onboarding (slides)
3. Login
4. Register
5. Forgot Password
6. OTP Verification
7. Identity Verification (CNI)

### Navigation Principale
1. Home (Feed)
2. Search/Explore
3. Map View
4. Favorites
5. Profile

### Détails & Actions
1. Offer Detail
2. Partner Detail
3. Outing Detail
4. Booking Flow
5. Check-in

### Profil & Settings
1. My Profile
2. My Bookings
3. My Favorites
4. Settings
5. Notifications

---

## 🔄 Stratégie de Cache

### Niveaux de Cache

```
┌─────────────────────────────────────┐
│          CACHE L1 (Memory)          │
│  - Données session courante         │
│  - Provider state                   │
│  TTL: Session                       │
└─────────────────┬───────────────────┘
                  ▼
┌─────────────────────────────────────┐
│          CACHE L2 (Hive/Isar)       │
│  - Offres récentes                  │
│  - Profil utilisateur               │
│  - Historique recherches            │
│  TTL: 1h - 24h selon type           │
└─────────────────┬───────────────────┘
                  ▼
┌─────────────────────────────────────┐
│          CACHE L3 (Network)         │
│  - Images (cached_network_image)    │
│  - Assets CDN                       │
│  TTL: 7 jours                       │
└─────────────────────────────────────┘
```

### Politiques de Cache

| Donnée | Stratégie | TTL |
|--------|-----------|-----|
| Profil utilisateur | Cache-first | 1h |
| Offres liste | Stale-while-revalidate | 5min |
| Offre détail | Cache-first | 15min |
| Catégories | Cache-first | 24h |
| Images | Cache permanent | 7j |

---

## 🧪 Tests Requis

### Tests Unitaires
- [ ] UseCases (domain layer)
- [ ] Repositories
- [ ] Providers/State

### Tests Widget
- [ ] Composants réutilisables
- [ ] Écrans principaux

### Tests E2E - Parcours Inscription
1. Lancement app
2. Skip onboarding
3. Tap "S'inscrire"
4. Saisie email + mot de passe
5. Validation OTP
6. Upload CNI
7. Validation identité
8. Accès home

---

## 📋 Checklist Qualité

- [ ] Animations 60fps
- [ ] Support dark mode
- [ ] Accessibilité (semantics)
- [ ] Internationalisation (FR/EN minimum)
- [ ] Deep linking
- [ ] Push notifications
- [ ] Analytics (Amplitude)
- [ ] Crash reporting (Sentry self-hosted)
- [ ] Performance monitoring

---

## 🔗 Références

- [Questions à clarifier](./QUESTIONS.md)
- [Composants identifiés](./COMPONENTS.md)
- [Design Figma](TODO)
