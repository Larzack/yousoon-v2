# 🧩 Composants Flutter Réutilisables

> Liste des composants à développer pour l'App Mobile Yousoon  
> À affiner après analyse complète du Figma

---

## 📦 Structure des Composants

```
lib/shared/widgets/
├── buttons/
│   ├── ys_button.dart
│   ├── ys_icon_button.dart
│   ├── ys_text_button.dart
│   └── ys_social_button.dart
├── inputs/
│   ├── ys_text_field.dart
│   ├── ys_search_bar.dart
│   ├── ys_phone_field.dart
│   ├── ys_otp_field.dart
│   └── ys_dropdown.dart
├── cards/
│   ├── ys_card.dart
│   ├── offer_card.dart
│   ├── outing_card.dart
│   ├── partner_card.dart
│   └── booking_card.dart
├── badges/
│   ├── ys_badge.dart
│   ├── discount_badge.dart
│   ├── category_chip.dart
│   └── status_badge.dart
├── feedback/
│   ├── ys_loader.dart
│   ├── ys_shimmer.dart
│   ├── ys_snackbar.dart
│   ├── ys_dialog.dart
│   └── ys_toast.dart
├── states/
│   ├── ys_empty_state.dart
│   ├── ys_error_state.dart
│   └── ys_loading_state.dart
├── media/
│   ├── ys_avatar.dart
│   ├── ys_image.dart
│   ├── ys_carousel.dart
│   └── ys_gallery.dart
├── navigation/
│   ├── ys_app_bar.dart
│   ├── ys_bottom_nav.dart
│   ├── ys_tab_bar.dart
│   └── ys_bottom_sheet.dart
├── layout/
│   ├── ys_scaffold.dart
│   ├── ys_section.dart
│   ├── ys_divider.dart
│   └── ys_spacing.dart
└── specific/
    ├── ys_rating.dart
    ├── ys_map_marker.dart
    ├── ys_filter_sheet.dart
    └── ys_qr_code.dart
```

---

## 🔘 Boutons

### YsButton
Bouton principal de l'application.

```dart
/// Variantes: primary, secondary, outlined, text, danger
/// Tailles: small, medium, large
/// États: enabled, disabled, loading

YsButton(
  label: 'Réserver',
  variant: YsButtonVariant.primary,
  size: YsButtonSize.large,
  isLoading: false,
  isDisabled: false,
  icon: Icons.calendar,
  iconPosition: IconPosition.leading,
  onPressed: () {},
)
```

**Props**:
- `label` : String - Texte du bouton
- `variant` : YsButtonVariant - Style (primary, secondary, outlined, text, danger)
- `size` : YsButtonSize - Taille (small, medium, large)
- `isLoading` : bool - Affiche un loader
- `isDisabled` : bool - Désactive le bouton
- `icon` : IconData? - Icône optionnelle
- `iconPosition` : IconPosition - Position de l'icône
- `onPressed` : VoidCallback - Action au tap
- `fullWidth` : bool - Prend toute la largeur

---

### YsSocialButton
Bouton de connexion sociale.

```dart
YsSocialButton(
  provider: SocialProvider.google,
  onPressed: () {},
)
```

---

## 📝 Inputs

### YsTextField
Champ de texte standard.

```dart
YsTextField(
  label: 'Email',
  hint: 'votre@email.com',
  controller: _emailController,
  keyboardType: TextInputType.emailAddress,
  validator: Validators.email,
  prefixIcon: Icons.email,
  suffixIcon: Icons.clear,
  onSuffixTap: () => _emailController.clear(),
)
```

**Props**:
- `label` : String - Label du champ
- `hint` : String - Placeholder
- `controller` : TextEditingController
- `keyboardType` : TextInputType
- `validator` : FormFieldValidator?
- `prefixIcon` / `suffixIcon` : IconData?
- `obscureText` : bool - Pour mots de passe
- `maxLines` : int
- `enabled` : bool
- `errorText` : String?

---

### YsOtpField
Champ de saisie OTP (4-6 chiffres).

```dart
YsOtpField(
  length: 6,
  onCompleted: (code) => verifyOtp(code),
  onChanged: (value) {},
)
```

---

### YsSearchBar
Barre de recherche avec suggestions.

```dart
YsSearchBar(
  hint: 'Rechercher une sortie...',
  onSearch: (query) {},
  onFilterTap: () => showFilters(),
  suggestions: recentSearches,
)
```

---

## 🃏 Cards

### OfferCard
Carte d'affichage d'une offre.

```dart
OfferCard(
  offer: offer,
  variant: OfferCardVariant.compact, // ou .expanded
  onTap: () => navigateToOffer(offer.id),
  onFavorite: () => toggleFavorite(offer.id),
)
```

**Affiche**:
- Image principale
- Badge réduction (ex: -30%)
- Nom de l'offre
- Nom du partenaire
- Catégorie
- Distance (si géoloc)
- Bouton favori

---

### OutingCard
Carte de sortie/événement.

```dart
OutingCard(
  outing: outing,
  onTap: () => navigateToOuting(outing.id),
  showStatus: true,
)
```

**Affiche**:
- Image
- Date/heure
- Titre
- Lieu
- Nombre de participants
- Statut (à venir, en cours, passé)

---

### PartnerCard
Carte partenaire/établissement.

```dart
PartnerCard(
  partner: partner,
  onTap: () => navigateToPartner(partner.id),
)
```

**Affiche**:
- Logo/Image
- Nom
- Catégorie
- Note moyenne
- Nombre d'offres actives
- Distance

---

## 🏷️ Badges

### DiscountBadge
Badge affichant une réduction.

```dart
DiscountBadge(
  percentage: 30,
  size: BadgeSize.medium,
)
```

---

### CategoryChip
Chip de catégorie cliquable.

```dart
CategoryChip(
  category: category,
  isSelected: true,
  onTap: () => selectCategory(category),
)
```

---

### StatusBadge
Badge de statut.

```dart
StatusBadge(
  status: BookingStatus.confirmed,
)
```

---

## 💬 Feedback

### YsLoader
Indicateur de chargement.

```dart
YsLoader(
  size: LoaderSize.medium,
  color: AppColors.primary,
)

// Variante overlay
YsLoader.overlay(
  message: 'Chargement...',
)
```

---

### YsShimmer
Skeleton loading.

```dart
YsShimmer(
  child: OfferCard.skeleton(),
)

// Ou liste
YsShimmer.list(
  itemCount: 5,
  itemBuilder: () => OfferCard.skeleton(),
)
```

---

### YsSnackbar
Notification en bas d'écran.

```dart
YsSnackbar.show(
  context,
  message: 'Offre ajoutée aux favoris',
  type: SnackbarType.success,
  action: SnackbarAction(
    label: 'Annuler',
    onPressed: () => undoFavorite(),
  ),
)
```

---

## 📭 États

### YsEmptyState
État vide personnalisable.

```dart
YsEmptyState(
  icon: Icons.search_off,
  title: 'Aucun résultat',
  description: 'Essayez avec d\'autres critères',
  action: YsButton(
    label: 'Réinitialiser',
    onPressed: () => resetFilters(),
  ),
)
```

---

### YsErrorState
État d'erreur.

```dart
YsErrorState(
  error: error,
  onRetry: () => refetch(),
)
```

---

## 🖼️ Media

### YsCarousel
Carousel d'images.

```dart
YsCarousel(
  images: offer.images,
  height: 200,
  autoPlay: true,
  showIndicators: true,
  onPageChanged: (index) {},
)
```

---

### YsAvatar
Avatar utilisateur ou partenaire.

```dart
YsAvatar(
  imageUrl: user.avatarUrl,
  name: user.displayName, // Pour initiales si pas d'image
  size: AvatarSize.medium,
  badge: AvatarBadge.verified,
)
```

---

## 🧭 Navigation

### YsBottomNav
Barre de navigation principale.

```dart
YsBottomNav(
  currentIndex: selectedIndex,
  onTap: (index) => navigateTo(index),
  items: [
    YsNavItem(icon: Icons.home, label: 'Accueil'),
    YsNavItem(icon: Icons.search, label: 'Explorer'),
    YsNavItem(icon: Icons.map, label: 'Carte'),
    YsNavItem(icon: Icons.favorite, label: 'Favoris'),
    YsNavItem(icon: Icons.person, label: 'Profil'),
  ],
)
```

---

### YsBottomSheet
Bottom sheet réutilisable.

```dart
YsBottomSheet.show(
  context,
  title: 'Filtres',
  child: FilterContent(),
  actions: [
    YsButton(label: 'Appliquer', onPressed: applyFilters),
  ],
)
```

---

## 🗺️ Spécifiques

### YsRating
Affichage/saisie de note.

```dart
// Lecture seule
YsRating(
  value: 4.5,
  size: RatingSize.small,
)

// Éditable
YsRating.interactive(
  value: rating,
  onChanged: (value) => setRating(value),
)
```

---

### YsQrCode
Affichage QR code pour check-in.

```dart
YsQrCode(
  data: booking.qrCodeData,
  size: 200,
  logo: 'assets/logo.png',
)
```

---

### YsFilterSheet
Sheet de filtres.

```dart
YsFilterSheet(
  filters: activeFilters,
  onApply: (filters) => applyFilters(filters),
  onReset: () => resetFilters(),
  sections: [
    FilterSection(
      title: 'Catégories',
      type: FilterType.multiSelect,
      options: categories,
    ),
    FilterSection(
      title: 'Distance',
      type: FilterType.slider,
      range: Range(0, 50),
    ),
    FilterSection(
      title: 'Réduction minimum',
      type: FilterType.slider,
      range: Range(0, 100),
    ),
  ],
)
```

---

## 📐 Design Tokens

### Espacements

```dart
abstract class YsSpacing {
  static const double xs = 4;
  static const double sm = 8;
  static const double md = 16;
  static const double lg = 24;
  static const double xl = 32;
  static const double xxl = 48;
}
```

### Rayons

```dart
abstract class YsRadius {
  static const double xs = 4;
  static const double sm = 8;
  static const double md = 12;
  static const double lg = 16;
  static const double xl = 24;
  static const double full = 999;
}
```

### Ombres

```dart
abstract class YsShadows {
  static const BoxShadow sm = BoxShadow(...);
  static const BoxShadow md = BoxShadow(...);
  static const BoxShadow lg = BoxShadow(...);
}
```

---

## ✅ Checklist Composants

- [ ] Props bien typées
- [ ] États gérés (enabled, disabled, loading, error)
- [ ] Thème light/dark supporté
- [ ] Accessibilité (Semantics)
- [ ] Animations fluides
- [ ] Tests unitaires
- [ ] Documentation Dart
- [ ] Exemple d'utilisation
