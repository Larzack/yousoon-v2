# 🚀 Plan de Génération - Application Mobile Flutter

> **Module** : App Mobile Yousoon (Flutter)  
> **Priorité** : 🟠 Haute (après Backend Phase 1-2)  
> **Dépendances** : Backend Gateway + Identity Service

---

## 📋 Vue d'Ensemble

```
┌─────────────────────────────────────────────────────────────────┐
│                    ORDRE DE GÉNÉRATION                          │
├─────────────────────────────────────────────────────────────────┤
│  Phase 1: Setup & Core (config, theme, navigation)             │
│  Phase 2: Design System (composants réutilisables)             │
│  Phase 3: Features Auth (inscription, login, CNI)              │
│  Phase 4: Features Core (offres, réservations, favoris)        │
│  Phase 5: Features Social (messagerie, avis)                   │
│  Phase 6: Tests & Polish                                        │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📁 Structure Cible

```
apps/mobile/
├── lib/
│   ├── main.dart
│   ├── app/
│   │   ├── app.dart
│   │   └── router.dart
│   ├── core/
│   │   ├── constants/
│   │   ├── errors/
│   │   ├── network/
│   │   ├── cache/
│   │   └── utils/
│   ├── features/
│   │   ├── auth/
│   │   ├── home/
│   │   ├── offers/
│   │   ├── bookings/
│   │   ├── favorites/
│   │   ├── profile/
│   │   ├── messaging/
│   │   ├── map/
│   │   └── settings/
│   ├── shared/
│   │   ├── widgets/
│   │   └── theme/
│   └── l10n/
├── assets/
├── test/
└── integration_test/
```

---

## 🔷 Phase 1 : Setup & Configuration

### Étape 1.1 : Initialisation Projet
**Fichiers à générer :**
```
apps/mobile/
├── pubspec.yaml                    # Dépendances
├── analysis_options.yaml           # Lint rules
├── lib/
│   └── main.dart                   # Entry point
```

**Dépendances clés :**
```yaml
dependencies:
  flutter_riverpod: ^2.4.0
  riverpod_annotation: ^2.3.0
  go_router: ^13.0.0
  graphql_flutter: ^5.1.0
  ferry: ^0.14.0
  hive_flutter: ^1.1.0
  flutter_secure_storage: ^9.0.0
  cached_network_image: ^3.3.0
  google_maps_flutter: ^2.5.0
  mobile_scanner: ^4.0.0
  local_auth: ^2.1.0
  in_app_purchase: ^3.1.0
  share_plus: ^7.2.0
  flutter_animate: ^4.3.0
  intl: ^0.18.0
```

### Étape 1.2 : Configuration Core
**Fichiers à générer :**
```
lib/core/
├── constants/
│   ├── app_constants.dart
│   ├── api_constants.dart
│   └── storage_keys.dart
├── config/
│   ├── env_config.dart
│   └── app_config.dart
├── errors/
│   ├── failures.dart
│   ├── exceptions.dart
│   └── error_handler.dart
└── utils/
    ├── extensions.dart
    ├── validators.dart
    └── formatters.dart
```

### Étape 1.3 : Network Layer
**Fichiers à générer :**
```
lib/core/network/
├── graphql_client.dart             # Client GraphQL configuré
├── network_info.dart               # Connectivity check
├── interceptors/
│   ├── auth_interceptor.dart
│   ├── logging_interceptor.dart
│   └── cache_interceptor.dart
└── api/
    └── api_client.dart
```

### Étape 1.4 : Cache Layer
**Fichiers à générer :**
```
lib/core/cache/
├── cache_manager.dart
├── cache_policy.dart
├── local_storage.dart              # Hive
└── secure_storage.dart             # flutter_secure_storage
```

### Étape 1.5 : Navigation
**Fichiers à générer :**
```
lib/app/
├── app.dart                        # MaterialApp
├── router.dart                     # GoRouter configuration
└── routes.dart                     # Route names constants
```

---

## 🔷 Phase 2 : Design System

### Étape 2.1 : Theme
**Fichiers à générer (basés sur Figma) :**
```
lib/shared/theme/
├── app_theme.dart                  # ThemeData
├── app_colors.dart                 # Palette Figma
│   - primary: #E99B27 (Indian Gold)
│   - background: #000000 (Dark Black)
│   - success: #5FC15C (Mantis Green)
│   - error: #CC2936 (Persian Red)
│   - inactive: #6D6D6D (Grey Jet)
├── app_typography.dart             # Futura font styles
└── app_spacing.dart                # Marges (5, 10, 15, 20, 25, 35, 45px)
```

### Étape 2.2 : Composants Boutons (CTAs)
**Basé sur Design System Figma :**
```
lib/shared/widgets/buttons/
├── ys_button.dart                  # Bouton générique
│   - Variants: primary, secondary, tertiary
│   - States: active, inactive
│   - Sizes: large (216x50), small (150x50)
├── ys_icon_button.dart
├── ys_text_button.dart
└── ys_floating_action_button.dart
```

### Étape 2.3 : Composants Inputs
**Fichiers à générer :**
```
lib/shared/widgets/inputs/
├── ys_text_field.dart              # Input avec underline
├── ys_search_field.dart            # Recherche avec loupe
├── ys_description_field.dart       # Multiline (30 chars)
├── ys_message_input.dart           # Input messagerie
└── ys_otp_field.dart               # Champs OTP
```

### Étape 2.4 : Composants Cards
**Fichiers à générer :**
```
lib/shared/widgets/cards/
├── ys_offer_card.dart              # Carte offre (swipe)
├── ys_outing_card.dart             # Carte réservation
├── ys_partner_card.dart            # Carte partenaire
├── ys_contact_card.dart            # Contact messagerie
└── ys_category_chip.dart           # Chips catégories
```

### Étape 2.5 : Composants Navigation
**Fichiers à générer :**
```
lib/shared/widgets/navigation/
├── ys_bottom_nav_bar.dart          # Tap Bar (5 entrées)
│   - Mes events, Favoris, Pour vous, Carte, Message
├── ys_app_bar.dart                 # Header avec notif + profil
├── ys_tab_bar.dart                 # Onglets (Événements/Yousooners)
└── ys_back_button.dart
```

### Étape 2.6 : Composants Feedback
**Fichiers à générer :**
```
lib/shared/widgets/feedback/
├── ys_loading.dart                 # Loader orange
├── ys_toast.dart                   # Toaster
├── ys_modal.dart                   # Pop-up comportementale/interactionnelle
├── ys_full_page_feedback.dart      # Validé/En cours/Refusé
└── ys_empty_state.dart
```

### Étape 2.7 : Composants Spécifiques
**Fichiers à générer :**
```
lib/shared/widgets/
├── ys_avatar.dart                  # Avatar avec grade
├── ys_rating.dart                  # Étoiles (avis)
├── ys_discount_badge.dart          # Badge réduction
├── ys_grade_badge.dart             # Explorateur, Aventurier, etc.
├── ys_qr_code.dart                 # Affichage QR
├── ys_image_slider.dart            # Carousel images
└── ys_map_marker.dart              # Pins carte
```

---

## 🔷 Phase 3 : Features Auth

### Étape 3.1 : Architecture Feature Auth
**Fichiers à générer :**
```
lib/features/auth/
├── data/
│   ├── datasources/
│   │   ├── auth_remote_datasource.dart
│   │   └── auth_local_datasource.dart
│   ├── models/
│   │   ├── user_model.dart
│   │   └── auth_token_model.dart
│   └── repositories/
│       └── auth_repository_impl.dart
├── domain/
│   ├── entities/
│   │   ├── user.dart
│   │   └── auth_token.dart
│   ├── repositories/
│   │   └── auth_repository.dart
│   └── usecases/
│       ├── register_usecase.dart
│       ├── login_usecase.dart
│       ├── logout_usecase.dart
│       └── verify_identity_usecase.dart
└── presentation/
    ├── providers/
    │   └── auth_provider.dart
    ├── screens/
    │   ├── splash_screen.dart
    │   ├── onboarding_screen.dart
    │   ├── login_screen.dart
    │   ├── register_screen.dart
    │   ├── otp_screen.dart
    │   ├── forgot_password_screen.dart
    │   └── identity_verification_screen.dart
    └── widgets/
        └── auth_form.dart
```

### Étape 3.2 : Écrans Onboarding (Slider)
**Basé sur Figma - Écrans slide :**
- Slide 1 : Bienvenue
- Slide 2 : Découvrez les offres
- Slide 3 : Réservez facilement
- CTA : S'inscrire / Se connecter

### Étape 3.3 : Écran Inscription
**Champs (selon Figma) :**
- Email
- Mot de passe
- Confirmation mot de passe
- Nom, Prénom
- Date de naissance
- Genre (optionnel)
- Acceptation CGU

### Étape 3.4 : Vérification CNI (OCR)
**Écran multi-étapes :**
1. Choix document (CNI, Passeport, Permis)
2. Photo recto
3. Photo verso (si CNI)
4. Selfie (optionnel)
5. Validation en cours
6. Résultat (Validé/Refusé)

---

## 🔷 Phase 4 : Features Core

### Étape 4.1 : Feature Home / Pour Vous
**Fichiers à générer :**
```
lib/features/home/
├── data/...
├── domain/...
└── presentation/
    ├── providers/
    │   └── home_provider.dart
    ├── screens/
    │   └── home_screen.dart        # "Pour vous" - swipe cards
    └── widgets/
        ├── offer_swipe_card.dart   # Carte plein écran avec swipe
        └── event_tabs.dart         # Onglets Événements/Yousooners
```

### Étape 4.2 : Feature Offers / Discovery
**Fichiers à générer :**
```
lib/features/offers/
├── data/
│   ├── datasources/
│   ├── models/
│   └── repositories/
├── domain/
│   ├── entities/
│   │   ├── offer.dart
│   │   ├── category.dart
│   │   └── discount.dart
│   └── usecases/
│       ├── get_offers_usecase.dart
│       ├── search_offers_usecase.dart
│       └── get_nearby_offers_usecase.dart
└── presentation/
    ├── providers/
    ├── screens/
    │   ├── offers_list_screen.dart
    │   ├── offer_detail_screen.dart
    │   └── search_screen.dart
    └── widgets/
```

### Étape 4.3 : Feature Bookings / Outings
**Fichiers à générer :**
```
lib/features/bookings/
├── data/...
├── domain/
│   ├── entities/
│   │   └── outing.dart
│   └── usecases/
│       ├── book_outing_usecase.dart
│       ├── get_my_outings_usecase.dart
│       └── checkin_usecase.dart
└── presentation/
    ├── screens/
    │   ├── my_outings_screen.dart  # Onglets: Passés/À venir/Créés
    │   ├── outing_detail_screen.dart
    │   ├── booking_flow_screen.dart
    │   └── checkin_screen.dart     # QR Code scanner
    └── widgets/
        └── qr_display.dart
```

### Étape 4.4 : Feature Favorites
**Fichiers à générer :**
```
lib/features/favorites/
├── data/...
├── domain/...
└── presentation/
    ├── screens/
    │   └── favorites_screen.dart
    └── widgets/
```

### Étape 4.5 : Feature Map
**Fichiers à générer :**
```
lib/features/map/
├── data/...
├── domain/...
└── presentation/
    ├── providers/
    │   └── map_provider.dart
    ├── screens/
    │   └── map_screen.dart         # Google Maps
    └── widgets/
        ├── offer_map_marker.dart
        ├── yousooner_marker.dart   # Pins par grade
        └── map_bottom_sheet.dart
```

---

## 🔷 Phase 5 : Features Social

### Étape 5.1 : Feature Messaging
**Fichiers à générer :**
```
lib/features/messaging/
├── data/...
├── domain/
│   ├── entities/
│   │   ├── conversation.dart
│   │   └── message.dart
│   └── usecases/
└── presentation/
    ├── screens/
    │   ├── conversations_list_screen.dart
    │   └── chat_screen.dart
    └── widgets/
        ├── message_bubble.dart
        └── contact_list.dart       # Ordre alphabétique
```

### Étape 5.2 : Feature Profile
**Fichiers à générer :**
```
lib/features/profile/
├── data/...
├── domain/...
└── presentation/
    ├── screens/
    │   ├── profile_screen.dart
    │   ├── edit_profile_screen.dart
    │   ├── subscription_screen.dart
    │   └── yousooner_profile_screen.dart  # Profil autre user
    └── widgets/
        ├── grade_progress.dart
        └── stats_card.dart
```

### Étape 5.3 : Feature Settings
**Fichiers à générer :**
```
lib/features/settings/
└── presentation/
    ├── screens/
    │   ├── settings_screen.dart
    │   ├── notifications_settings_screen.dart
    │   ├── privacy_screen.dart
    │   └── language_screen.dart
    └── widgets/
```

---

## 🔷 Phase 6 : Tests & Polish

### Étape 6.1 : Tests Unitaires
```
test/
├── core/
├── features/
│   ├── auth/
│   │   ├── domain/usecases/
│   │   └── data/repositories/
│   └── ...
└── shared/
```

### Étape 6.2 : Tests Widget
```
test/
└── widgets/
    ├── buttons/
    ├── cards/
    └── inputs/
```

### Étape 6.3 : Tests E2E - Parcours Inscription
```
integration_test/
├── app_test.dart
└── auth/
    └── registration_test.dart
```

**Scénario E2E :**
1. Launch app
2. Skip onboarding
3. Tap "S'inscrire"
4. Fill registration form
5. Verify OTP
6. Upload CNI
7. Wait validation
8. Access home

### Étape 6.4 : Internationalisation
```
lib/l10n/
├── app_fr.arb
├── app_en.arb
└── l10n.dart
```

---

## ⏱️ Estimation des Temps

| Phase | Étape | Durée estimée |
|-------|-------|---------------|
| **Phase 1** | Setup & Config | 2h |
| | Network & Cache | 2h |
| | Navigation | 1h |
| **Phase 2** | Theme | 1h |
| | Composants UI | 4h |
| | Navigation Components | 2h |
| **Phase 3** | Auth Architecture | 2h |
| | Écrans Auth | 4h |
| | Vérification CNI | 3h |
| **Phase 4** | Home/Pour Vous | 3h |
| | Offers/Discovery | 4h |
| | Bookings | 3h |
| | Favorites | 1h |
| | Map | 3h |
| **Phase 5** | Messaging | 3h |
| | Profile | 2h |
| | Settings | 1h |
| **Phase 6** | Tests | 4h |
| | i18n | 1h |
| **Total** | | **~46h** |

---

## ✅ Critères de Validation

### Composants UI
- [ ] Tous les composants matchent le Figma
- [ ] Dark mode appliqué partout
- [ ] Animations fluides 60fps
- [ ] Responsif (différentes tailles écran)

### Features
- [ ] Inscription complète fonctionne
- [ ] Connexion avec biométrie
- [ ] Offres chargent et s'affichent
- [ ] Réservation + QR code OK
- [ ] Carte Google Maps fonctionnelle
- [ ] Messagerie temps réel

### Performance
- [ ] Cold start < 2s
- [ ] Navigation instantanée
- [ ] Images cached
- [ ] Mode offline (favoris, historique)

### Tests
- [ ] Tests unitaires >80% coverage
- [ ] Tests widget pour composants clés
- [ ] E2E inscription passe

---

## 🔗 Références

- [Design System](../DESIGN_SYSTEM.md)
- [Prompt App Mobile](./PROMPT.md)
- [Composants identifiés](./COMPONENTS.md)
- [Figma](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2)
