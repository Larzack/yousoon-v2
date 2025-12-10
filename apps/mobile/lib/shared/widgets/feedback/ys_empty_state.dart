import 'package:flutter/material.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../buttons/ys_button.dart';

/// Widget pour afficher un état vide
class YsEmptyState extends StatelessWidget {
  /// Icône à afficher
  final IconData icon;

  /// Titre principal
  final String title;

  /// Message descriptif
  final String? message;

  /// Texte du bouton d'action
  final String? actionText;

  /// Callback du bouton d'action
  final VoidCallback? onAction;

  /// Taille de l'icône
  final double iconSize;

  const YsEmptyState({
    super.key,
    required this.icon,
    required this.title,
    this.message,
    this.actionText,
    this.onAction,
    this.iconSize = 80,
  });

  /// État vide pour les favoris
  factory YsEmptyState.favorites({VoidCallback? onExplore}) {
    return YsEmptyState(
      icon: Icons.favorite_outline_rounded,
      title: 'Aucun favori',
      message: 'Explorez les offres et ajoutez vos préférées à vos favoris.',
      actionText: 'Explorer',
      onAction: onExplore,
    );
  }

  /// État vide pour les sorties
  factory YsEmptyState.outings({VoidCallback? onExplore}) {
    return YsEmptyState(
      icon: Icons.event_available_outlined,
      title: 'Aucune sortie',
      message: 'Réservez votre première offre pour commencer l\'aventure !',
      actionText: 'Découvrir les offres',
      onAction: onExplore,
    );
  }

  /// État vide pour les recherches
  factory YsEmptyState.searchResults() {
    return const YsEmptyState(
      icon: Icons.search_off_rounded,
      title: 'Aucun résultat',
      message: 'Essayez avec d\'autres mots-clés ou modifiez vos filtres.',
    );
  }

  /// État vide pour les notifications
  factory YsEmptyState.notifications() {
    return const YsEmptyState(
      icon: Icons.notifications_none_rounded,
      title: 'Aucune notification',
      message: 'Vous êtes à jour ! Les nouvelles notifications apparaîtront ici.',
    );
  }

  /// État vide pour les messages
  factory YsEmptyState.messages() {
    return const YsEmptyState(
      icon: Icons.chat_bubble_outline_rounded,
      title: 'Aucun message',
      message: 'Vos conversations apparaîtront ici.',
    );
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.xl),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          // Icône
          Container(
            padding: const EdgeInsets.all(AppSpacing.lg),
            decoration: BoxDecoration(
              color: AppColors.primary.withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: Icon(
              icon,
              size: iconSize,
              color: AppColors.primary,
            ),
          ),

          const SizedBox(height: AppSpacing.xl),

          // Titre
          Text(
            title,
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
            textAlign: TextAlign.center,
          ),

          // Message
          if (message != null) ...[
            const SizedBox(height: AppSpacing.md),
            Text(
              message!,
              style: const TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],

          // Bouton d'action
          if (actionText != null && onAction != null) ...[
            const SizedBox(height: AppSpacing.xl),
            YsButton(
              text: actionText!,
              onPressed: onAction,
            ),
          ],
        ],
      ),
    );
  }
}

/// Widget pour afficher un état d'erreur
class YsErrorState extends StatelessWidget {
  /// Message d'erreur
  final String message;

  /// Texte du bouton de retry
  final String retryText;

  /// Callback du retry
  final VoidCallback? onRetry;

  /// Icône
  final IconData icon;

  const YsErrorState({
    super.key,
    required this.message,
    this.retryText = 'Réessayer',
    this.onRetry,
    this.icon = Icons.error_outline_rounded,
  });

  /// Erreur réseau
  factory YsErrorState.network({VoidCallback? onRetry}) {
    return YsErrorState(
      icon: Icons.wifi_off_rounded,
      message: 'Impossible de se connecter. Vérifiez votre connexion internet.',
      onRetry: onRetry,
    );
  }

  /// Erreur générique
  factory YsErrorState.generic({VoidCallback? onRetry}) {
    return YsErrorState(
      message: 'Une erreur est survenue. Veuillez réessayer.',
      onRetry: onRetry,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.xl),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          // Icône
          Container(
            padding: const EdgeInsets.all(AppSpacing.lg),
            decoration: BoxDecoration(
              color: AppColors.error.withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: Icon(
              icon,
              size: 60,
              color: AppColors.error,
            ),
          ),

          const SizedBox(height: AppSpacing.xl),

          // Message
          Text(
            message,
            style: const TextStyle(
              fontSize: 14,
              color: AppColors.textSecondary,
            ),
            textAlign: TextAlign.center,
          ),

          // Bouton retry
          if (onRetry != null) ...[
            const SizedBox(height: AppSpacing.xl),
            YsButton.secondary(
              text: retryText,
              onPressed: onRetry,
              icon: Icons.refresh_rounded,
            ),
          ],
        ],
      ),
    );
  }
}

/// Widget pour afficher un état de succès
class YsSuccessState extends StatelessWidget {
  /// Titre
  final String title;

  /// Message
  final String? message;

  /// Texte du bouton
  final String? actionText;

  /// Callback du bouton
  final VoidCallback? onAction;

  const YsSuccessState({
    super.key,
    required this.title,
    this.message,
    this.actionText,
    this.onAction,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.xl),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          // Icône
          Container(
            padding: const EdgeInsets.all(AppSpacing.lg),
            decoration: BoxDecoration(
              color: AppColors.success.withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.check_circle_outline_rounded,
              size: 80,
              color: AppColors.success,
            ),
          ),

          const SizedBox(height: AppSpacing.xl),

          // Titre
          Text(
            title,
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: AppColors.textPrimary,
            ),
            textAlign: TextAlign.center,
          ),

          // Message
          if (message != null) ...[
            const SizedBox(height: AppSpacing.md),
            Text(
              message!,
              style: const TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],

          // Bouton
          if (actionText != null && onAction != null) ...[
            const SizedBox(height: AppSpacing.xl),
            YsButton(
              text: actionText!,
              onPressed: onAction,
            ),
          ],
        ],
      ),
    );
  }
}
