# 🎨 Design System Yousoon

> Extrait du Figma : [Yousoon-Test2](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2?node-id=121-114)  
> **Date d'extraction** : 9 décembre 2025  
> **Fichier Figma** : `1GXJECHtsYzq46OYbSHiaj`

---

## 📋 Table des Matières

1. [Couleurs](#couleurs)
2. [Typographie](#typographie)
3. [Espacements](#espacements)
4. [Composants](#composants)
5. [Icônes](#icônes)
6. [Navigation](#navigation)

---

## 🎨 Couleurs

### Palette Principale

| Nom | Hex | RGB | Usage |
|-----|-----|-----|-------|
| **Dark Black** | `#000000` | 0, 0, 0 | Couleur principale, fond d'écran |
| **Indian Gold** | `#E99B27` | 233, 155, 39 | Couleur d'accent, CTA actifs, éléments interactifs |
| **Flash White** | `#FFFFFF` | 255, 255, 255 | Texte sur fond sombre, éléments secondaires |
| **Grey Jet** | `#6D6D6D` | 109, 109, 109 | Couleur inactive, placeholders |
| **Eerie Black** | `#CCCCCC` | 204, 204, 204 | Détails, texte secondaire |
| **Mantis Green** | `#5FC15C` | 95, 193, 92 | Validation, succès, profil vérifié |
| **Persian Red** | `#CC2936` | 204, 41, 54 | Erreurs, refus, actions critiques |

### Usage des Couleurs

```dart
// Flutter - app_colors.dart
abstract class AppColors {
  // Primaires
  static const Color primary = Color(0xFFE99B27);      // Indian Gold
  static const Color background = Color(0xFF000000);   // Dark Black
  static const Color surface = Color(0xFF000000);      // Dark Black
  
  // Texte
  static const Color textPrimary = Color(0xFFFFFFFF);  // Flash White
  static const Color textSecondary = Color(0xFFCCCCCC); // Eerie Black
  static const Color textDisabled = Color(0xFF6D6D6D); // Grey Jet
  
  // Feedback
  static const Color success = Color(0xFF5FC15C);      // Mantis Green
  static const Color error = Color(0xFFCC2936);        // Persian Red
  static const Color warning = Color(0xFFE99B27);      // Indian Gold
  
  // Éléments UI
  static const Color inactive = Color(0xFF6D6D6D);     // Grey Jet
  static const Color divider = Color(0xFF6D6D6D);      // Grey Jet
  static const Color cardBackground = Color(0xFF1A1A1A); // Noir légèrement plus clair
}
```

---

## 📝 Typographie

### Police Principale

| Police | Variantes | Usage |
|--------|-----------|-------|
| **Futura** | Medium, Bold | Titres, textes principaux |
| **Poppins** | Regular, SemiBold | Textes secondaires, corps |

### Échelle Typographique

| Style | Taille | Poids | Usage |
|-------|--------|-------|-------|
| **Titre 1** | 16px | Bold | Titres principaux |
| **Titre 2** | 16px | Medium | Sous-titres |
| **Titre 3** | 14px | Medium | Titres de section |
| **Corps texte** | 14px | Medium | Texte standard |
| **Texte secondaire** | 14px | Medium | Informations complémentaires |
| **CTA** | 16px | Medium | Boutons |
| **Instructions** | 14px | Medium Italic | Instructions d'inscription |
| **Grade Map** | 12px | Medium | Badges de niveau |

### Flutter Implementation

```dart
// Flutter - app_typography.dart
abstract class AppTypography {
  static const String fontFamilyPrimary = 'Futura';
  static const String fontFamilySecondary = 'Poppins';
  
  static TextStyle get headline1 => const TextStyle(
    fontFamily: fontFamilyPrimary,
    fontSize: 16,
    fontWeight: FontWeight.bold,
    color: Colors.white,
  );
  
  static TextStyle get headline2 => const TextStyle(
    fontFamily: fontFamilyPrimary,
    fontSize: 16,
    fontWeight: FontWeight.w500,
    color: Colors.white,
  );
  
  static TextStyle get headline3 => const TextStyle(
    fontFamily: fontFamilyPrimary,
    fontSize: 14,
    fontWeight: FontWeight.w500,
    color: Colors.white,
  );
  
  static TextStyle get bodyText => const TextStyle(
    fontFamily: fontFamilyPrimary,
    fontSize: 14,
    fontWeight: FontWeight.w500,
    color: Colors.white,
  );
  
  static TextStyle get button => const TextStyle(
    fontFamily: fontFamilyPrimary,
    fontSize: 16,
    fontWeight: FontWeight.w500,
    color: Colors.black,
  );
}
```

---

## 📐 Espacements

### Marges Standards

| Espacement | Valeur | Usage |
|------------|--------|-------|
| **xs** | 5px | Texte dans même bloc |
| **sm** | 10px | Entre éléments proches |
| **md** | 15px | Entre blocs de texte |
| **lg** | 20px | Entre blocs de même type |
| **xl** | 25px | Marge horizontale écran, avant CTA |
| **xxl** | 35px | Après titre principal |
| **xxxl** | 45px | Grande séparation, types différents |

### Flutter Implementation

```dart
// Flutter - app_spacing.dart
abstract class AppSpacing {
  static const double xs = 5.0;
  static const double sm = 10.0;
  static const double md = 15.0;
  static const double lg = 20.0;
  static const double xl = 25.0;
  static const double xxl = 35.0;
  static const double xxxl = 45.0;
  
  // Paddings écran
  static const EdgeInsets screenPadding = EdgeInsets.symmetric(horizontal: xl);
  static const EdgeInsets cardPadding = EdgeInsets.all(md);
}
```

---

## 🧩 Composants

### Boutons (CTA)

#### Primaire
- **Actif** : Fond Indian Gold (#E99B27), texte noir
- **Inactif** : Fond gris (#6D6D6D), texte gris
- **Tailles** : Grand (216x50px), Petit (150x50px)
- **Border radius** : 8px

#### Secondaire
- **Style** : Bordure blanche, fond transparent
- **Tailles** : Grand (216x50px), Petit (150x50px)

#### Tertiaire
- **Style** : Texte souligné uniquement
- **Couleur** : Blanc ou noir selon fond

### Inputs

| Type | Description |
|------|-------------|
| **Recherche** | Icône loupe + placeholder gris |
| **Description** | Zone de texte multi-lignes, max 30 caractères indiqué |
| **Infos perso** | Label + underline, requis avec * orange |
| **Message** | Input + icône envoi orange |

### Cards

- **Fond** : Noir ou image plein écran
- **Overlay** : Gradient du bas pour lisibilité texte
- **Contenu** : Titre, étoiles, date/lieu
- **Action** : Bouton cœur (favoris), chevron

### Onglets (Tabs)

- **Actif** : Texte blanc, underline orange
- **Inactif** : Texte gris

#### Types d'onglets
| Groupe | Options |
|--------|--------|
| **Principal** | ÉVÉNEMENTS / YOUSOONERS |
| **Mes events** | Passés / À venir / Créés |
| **Liste/Calendrier** | Vue liste / Vue calendrier |

### Pop-ups / Modales

| Type | Usage | Variantes Figma |
|------|-------|----------------|
| **Comportementale** | Information, changement abonnement | Défaut, Changement abonnement |
| **Interactionnelle** | Demande connexion, localisation | Connexion, Localisation |
| **Évaluation** | Notation étoiles + chips + commentaire | - |

> **Note Figma** : Les modales comportementales informent l'utilisateur pour continuer sur un parcours. Les modales interactionnelles invitent l'utilisateur à faire une action.

### Toasters

- Notification courte
- Avec image ou pictogramme
- Disparaît automatiquement

### Feedback Pleine Page

| État | Icône | Couleur |
|------|-------|---------|
| **En cours** | Loader orange | Orange |
| **Validé** | Check vert | Mantis Green |
| **Refusé** | Croix rouge | Persian Red |

---

## 🎯 Icônes

### Navigation (Tab Bar)

| Icône | Nom | États |
|-------|-----|-------|
| 📅 | Mes events | Actif (orange), Inactif (gris) |
| ❤️ | Favoris | Actif (orange), Inactif (gris) |
| 🃏 | Pour vous | Actif (orange), Inactif (gris) |
| 📍 | Carte | Actif (orange), Inactif (gris) |
| 💬 | Messagerie | Actif (orange), Inactif (gris) |

### Interaction

| Icône | Usage |
|-------|-------|
| 🔍 | Recherche |
| ✏️ | Édition |
| ❤️ | Favori (plein/vide) |
| 👁️ | Voir (ouvert/fermé) |
| 📤 | Partager |
| ➕ | Ajouter |
| ✅ | Validé |
| ❌ | Refusé/Annuler |
| 🔔 | Notifications |
| 📞 | Appeler |

### Réseaux Sociaux

- Instagram
- LinkedIn
- Facebook
- Google
- Apple

### Paiement

- Visa
- Mastercard
- American Express
- CB
- PayPal
- Apple Pay
- Prélèvement

### Grades Yousooner (sur carte)

| Grade | Emoji | Pin Map |
|-------|-------|--------|
| Explorateur | 🧭 | Pin personnalisé |
| Aventurier | 🎒 | Pin personnalisé |
| Grand voyageur | ✈️ | Pin personnalisé |
| Conquérant | 👑 | Pin personnalisé |

### Icônes Carte (Pins)

| Type | Description |
|------|-------------|
| **Position** | Position utilisateur actuelle |
| **Breakfast** | Établissement petit-déjeuner |
| **Restaurant** | Établissement restaurant |
| **Movie** | Cinéma/Divertissement |
| **Nature** | Activités nature |

### Catégories d'Intérêts

| Catégorie | Emoji |
|-----------|-------|
| Sport | 🎾 |
| Convivialité | 🥂 |
| Arts & Culture | 🎨 |
| Voyage & Escapade | 🌎 |
| Concert & Musique | 🎺 |
| Bien-être | 🧘 |
| Nature | 🌱 |
| Gastronomie & Dégustation | 🍴 |
| Évasion | ⛵ |
| Développement personnel | 💫 |

---

## 📱 Navigation

### Tab Bar (5 entrées)

```
┌────────┬────────┬────────┬────────┬────────┐
│  Mes   │        │  Pour  │        │        │
│ events │Favoris │  vous  │ Carte  │Message │
└────────┴────────┴────────┴────────┴────────┘
```

### Header Bar (2 entrées à droite)

- 🔔 Notifications
- 👤 Profil

### Page par défaut

**"Pour vous"** est la page d'accueil par défaut

---

## 🎬 Animations

### Swipe Events

- Swipe horizontal entre événements
- Image plein écran avec overlay gradient
- Bouton favori (cœur) en bas à droite
- Bouton recherche (loupe) en haut à gauche

### Transitions

- Fade in/out pour les modales
- Slide up pour les bottom sheets
- Scale pour les boutons au tap

---

## 📋 UX Writing

### Conventions

- **Vouvoiement** : Toujours utiliser "vous"
- **Proactif** : Verbes d'action (Se désinscrire, Se connecter, Créer)
- **Genre** : Adjectif féminin ou masculin selon le choix utilisateur

---

## 🎯 Guidelines de Design (extrait Figma)

### Choix des couleurs
- **Orange** : Contact humain, vitalité, communauté, convivialité, goût pour la nouveauté (couleur stimulante, couleur des Épicuriens)
- **Jaune/Or** : Action, productivité, luxe, pouvoir, puissance, chaleur, abondance
- **Nuances de gris** : Modernité et élégance

### Choix des images
- Photos (et non illustrations) pour la mise en avant de l'humain

### Choix des formes
- Droites et carrées pour le côté minimalisme et haut de gamme

### Choix des traits
- Longs et fins pour garder de la finesse et de l'élégance

---

## 🔗 Références

- [Figma Design](https://www.figma.com/design/1GXJECHtsYzq46OYbSHiaj/Yousoon-Test2?node-id=121-114)
- [App Mobile PROMPT](./app-mobile/PROMPT.md)
- [Composants Flutter](./app-mobile/COMPONENTS.md)
