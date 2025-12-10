import 'package:flutter/material.dart';

import '../../../core/theme/app_colors.dart';

/// Widget d'affichage des étoiles de notation
class YsRating extends StatelessWidget {
  /// Note actuelle (0-5)
  final double rating;

  /// Nombre maximum d'étoiles
  final int maxRating;

  /// Taille des étoiles
  final double size;

  /// Couleur des étoiles pleines
  final Color? activeColor;

  /// Couleur des étoiles vides
  final Color? inactiveColor;

  /// Espacement entre les étoiles
  final double spacing;

  /// Afficher la valeur numérique
  final bool showValue;

  /// Nombre d'avis (optionnel)
  final int? reviewCount;

  /// Callback quand on tape sur une étoile (rend interactif)
  final ValueChanged<int>? onRatingChanged;

  const YsRating({
    super.key,
    required this.rating,
    this.maxRating = 5,
    this.size = 16,
    this.activeColor,
    this.inactiveColor,
    this.spacing = 2,
    this.showValue = false,
    this.reviewCount,
    this.onRatingChanged,
  });

  /// Créer un rating simple (lecture seule)
  const YsRating.small({
    super.key,
    required this.rating,
    this.reviewCount,
  })  : maxRating = 5,
        size = 12,
        activeColor = null,
        inactiveColor = null,
        spacing = 1,
        showValue = true,
        onRatingChanged = null;

  /// Créer un rating moyen
  const YsRating.medium({
    super.key,
    required this.rating,
    this.reviewCount,
  })  : maxRating = 5,
        size = 18,
        activeColor = null,
        inactiveColor = null,
        spacing = 2,
        showValue = true,
        onRatingChanged = null;

  /// Créer un rating large (pour les détails)
  const YsRating.large({
    super.key,
    required this.rating,
    this.reviewCount,
  })  : maxRating = 5,
        size = 24,
        activeColor = null,
        inactiveColor = null,
        spacing = 4,
        showValue = true,
        onRatingChanged = null;

  /// Créer un rating interactif
  const YsRating.interactive({
    super.key,
    required this.rating,
    required this.onRatingChanged,
  })  : maxRating = 5,
        size = 32,
        activeColor = null,
        inactiveColor = null,
        spacing = 8,
        showValue = false,
        reviewCount = null;

  @override
  Widget build(BuildContext context) {
    final active = activeColor ?? AppColors.primary;
    final inactive = inactiveColor ?? AppColors.textDisabled;

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        // Étoiles
        ...List.generate(maxRating, (index) {
          final starValue = index + 1;
          final isFull = rating >= starValue;
          final isHalf = rating >= starValue - 0.5 && rating < starValue;

          IconData icon;
          Color color;

          if (isFull) {
            icon = Icons.star_rounded;
            color = active;
          } else if (isHalf) {
            icon = Icons.star_half_rounded;
            color = active;
          } else {
            icon = Icons.star_outline_rounded;
            color = inactive;
          }

          Widget star = Icon(
            icon,
            size: size,
            color: color,
          );

          if (onRatingChanged != null) {
            star = GestureDetector(
              onTap: () => onRatingChanged!(starValue),
              child: Padding(
                padding: EdgeInsets.symmetric(horizontal: spacing / 2),
                child: star,
              ),
            );
          } else {
            star = Padding(
              padding: EdgeInsets.only(right: spacing),
              child: star,
            );
          }

          return star;
        }),

        // Valeur numérique
        if (showValue) ...[
          const SizedBox(width: 4),
          Text(
            rating.toStringAsFixed(1),
            style: TextStyle(
              fontSize: size * 0.75,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
        ],

        // Nombre d'avis
        if (reviewCount != null) ...[
          const SizedBox(width: 4),
          Text(
            '($reviewCount)',
            style: TextStyle(
              fontSize: size * 0.65,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ],
    );
  }
}

/// Widget de sélection de rating (pour les formulaires)
class YsRatingSelector extends StatefulWidget {
  final int initialRating;
  final ValueChanged<int> onChanged;
  final double size;

  const YsRatingSelector({
    super.key,
    this.initialRating = 0,
    required this.onChanged,
    this.size = 40,
  });

  @override
  State<YsRatingSelector> createState() => _YsRatingSelectorState();
}

class _YsRatingSelectorState extends State<YsRatingSelector> {
  late int _currentRating;

  @override
  void initState() {
    super.initState();
    _currentRating = widget.initialRating;
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: List.generate(5, (index) {
        final starValue = index + 1;
        final isSelected = _currentRating >= starValue;

        return GestureDetector(
          onTap: () {
            setState(() {
              _currentRating = starValue;
            });
            widget.onChanged(starValue);
          },
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: AnimatedScale(
              scale: isSelected ? 1.1 : 1.0,
              duration: const Duration(milliseconds: 150),
              child: Icon(
                isSelected ? Icons.star_rounded : Icons.star_outline_rounded,
                size: widget.size,
                color: isSelected ? AppColors.primary : AppColors.textDisabled,
              ),
            ),
          ),
        );
      }),
    );
  }
}
