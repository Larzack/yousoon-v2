import 'package:flutter/material.dart';

/// Classe de localisation pour Yousoon
/// Supporte FR et EN
class AppLocalizations {
  final Locale locale;

  AppLocalizations(this.locale);

  static AppLocalizations? of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations);
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  static final Map<String, Map<String, String>> _localizedValues = {
    'fr': {
      // Navigation
      'nav_for_you': 'Pour vous',
      'nav_favorites': 'Favoris',
      'nav_map': 'Carte',
      'nav_outings': 'Mes sorties',
      'nav_messages': 'Messages',
      
      // Auth
      'login': 'Connexion',
      'register': 'Inscription',
      'email': 'Email',
      'password': 'Mot de passe',
      'forgot_password': 'Mot de passe oublié ?',
      'no_account': 'Pas encore de compte ?',
      'already_account': 'Déjà un compte ?',
      'logout': 'Déconnexion',
      
      // Home
      'for_you': 'Pour vous',
      'discover': 'Découvrir',
      'nearby': 'À proximité',
      'popular': 'Populaires',
      'new_offers': 'Nouveautés',
      
      // Offers
      'offer_details': 'Détails de l\'offre',
      'book_now': 'Réserver',
      'add_to_favorites': 'Ajouter aux favoris',
      'remove_from_favorites': 'Retirer des favoris',
      'discount': 'Réduction',
      'valid_until': 'Valable jusqu\'au',
      'conditions': 'Conditions',
      
      // Bookings
      'my_outings': 'Mes sorties',
      'upcoming': 'À venir',
      'past': 'Passées',
      'cancelled': 'Annulées',
      'booking_confirmed': 'Réservation confirmée',
      'show_qr_code': 'Afficher le QR code',
      'check_in': 'Check-in',
      
      // Profile
      'profile': 'Profil',
      'edit_profile': 'Modifier le profil',
      'settings': 'Paramètres',
      'notifications': 'Notifications',
      'language': 'Langue',
      'help': 'Aide',
      'about': 'À propos',
      
      // Common
      'loading': 'Chargement...',
      'error': 'Erreur',
      'retry': 'Réessayer',
      'cancel': 'Annuler',
      'confirm': 'Confirmer',
      'save': 'Enregistrer',
      'search': 'Rechercher',
      'no_results': 'Aucun résultat',
      'see_all': 'Voir tout',
      'see_more': 'Voir plus',
    },
    'en': {
      // Navigation
      'nav_for_you': 'For you',
      'nav_favorites': 'Favorites',
      'nav_map': 'Map',
      'nav_outings': 'My outings',
      'nav_messages': 'Messages',
      
      // Auth
      'login': 'Login',
      'register': 'Sign up',
      'email': 'Email',
      'password': 'Password',
      'forgot_password': 'Forgot password?',
      'no_account': 'No account yet?',
      'already_account': 'Already have an account?',
      'logout': 'Logout',
      
      // Home
      'for_you': 'For you',
      'discover': 'Discover',
      'nearby': 'Nearby',
      'popular': 'Popular',
      'new_offers': 'New',
      
      // Offers
      'offer_details': 'Offer details',
      'book_now': 'Book now',
      'add_to_favorites': 'Add to favorites',
      'remove_from_favorites': 'Remove from favorites',
      'discount': 'Discount',
      'valid_until': 'Valid until',
      'conditions': 'Conditions',
      
      // Bookings
      'my_outings': 'My outings',
      'upcoming': 'Upcoming',
      'past': 'Past',
      'cancelled': 'Cancelled',
      'booking_confirmed': 'Booking confirmed',
      'show_qr_code': 'Show QR code',
      'check_in': 'Check-in',
      
      // Profile
      'profile': 'Profile',
      'edit_profile': 'Edit profile',
      'settings': 'Settings',
      'notifications': 'Notifications',
      'language': 'Language',
      'help': 'Help',
      'about': 'About',
      
      // Common
      'loading': 'Loading...',
      'error': 'Error',
      'retry': 'Retry',
      'cancel': 'Cancel',
      'confirm': 'Confirm',
      'save': 'Save',
      'search': 'Search',
      'no_results': 'No results',
      'see_all': 'See all',
      'see_more': 'See more',
    },
  };

  String translate(String key) {
    return _localizedValues[locale.languageCode]?[key] ?? key;
  }

  // Getters pour les traductions courantes
  String get navForYou => translate('nav_for_you');
  String get navFavorites => translate('nav_favorites');
  String get navMap => translate('nav_map');
  String get navOutings => translate('nav_outings');
  String get navMessages => translate('nav_messages');
  String get login => translate('login');
  String get register => translate('register');
  String get email => translate('email');
  String get password => translate('password');
  String get bookNow => translate('book_now');
  String get loading => translate('loading');
  String get error => translate('error');
  String get retry => translate('retry');
  String get cancel => translate('cancel');
  String get confirm => translate('confirm');
  String get search => translate('search');
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) {
    return ['fr', 'en'].contains(locale.languageCode);
  }

  @override
  Future<AppLocalizations> load(Locale locale) async {
    return AppLocalizations(locale);
  }

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}
