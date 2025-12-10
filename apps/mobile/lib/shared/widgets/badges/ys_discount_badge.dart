import 'package:flutter/material.dart';

import '../../../core/theme/app_colors.dart';

/// Badge affichant une réduction
class YsDiscountBadge extends StatelessWidget {
  /// Type de réduction (percentage, fixed, formula)
  final String type;

  /// Valeur de la réduction
  final int value;

  /// Formule personnalisée (ex: "1 acheté = 1 offert")
  final String? formula;

  /// Taille du badge
  final YsDiscountBadgeSize size;

  const YsDiscountBadge({
    super.key,
    required this.type,
    required this.value,
    this.formula,
    this.size = YsDiscountBadgeSize.medium,
  });

  /// Badge pourcentage
  factory YsDiscountBadge.percentage(int value, {YsDiscountBadgeSize size = YsDiscountBadgeSize.medium}) {
    return YsDiscountBadge(type: 'percentage', value: value, size: size);
  }

  /// Badge montant fixe
  factory YsDiscountBadge.fixed(int value, {YsDiscountBadgeSize size = YsDiscountBadgeSize.medium}) {
    return YsDiscountBadge(type: 'fixed', value: value, size: size);
  }

  /// Badge formule
  factory YsDiscountBadge.formula(String formula, {YsDiscountBadgeSize size = YsDiscountBadgeSize.medium}) {
    return YsDiscountBadge(type: 'formula', value: 0, formula: formula, size: size);
  }

  String get _displayText {
    switch (type) {
      case 'percentage':
        return '-$value%';
      case 'fixed':
        // Convertir les centimes en euros
        final euros = value / 100;
        if (euros == euros.toInt()) {
          return '-${euros.toInt()}€';
        }
        return '-${euros.toStringAsFixed(2)}€';
      case 'formula':
        return formula ?? '';
      default:
        return '-$value%';
    }
  }

  @override
  Widget build(BuildContext context) {
    final config = size._config;

    return Container(
      padding: EdgeInsets.symmetric(
        horizontal: config.horizontalPadding,
        vertical: config.verticalPadding,
      ),
      decoration: BoxDecoration(
        color: AppColors.primary,
        borderRadius: BorderRadius.circular(config.borderRadius),
      ),
      child: Text(
        _displayText,
        style: TextStyle(
          color: Colors.black,
          fontSize: config.fontSize,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}

/// Tailles du badge
enum YsDiscountBadgeSize {
  small,
  medium,
  large,
}

extension on YsDiscountBadgeSize {
  _BadgeConfig get _config {
    switch (this) {
      case YsDiscountBadgeSize.small:
        return const _BadgeConfig(
          fontSize: 10,
          horizontalPadding: 6,
          verticalPadding: 2,
          borderRadius: 4,
        );
      case YsDiscountBadgeSize.medium:
        return const _BadgeConfig(
          fontSize: 12,
          horizontalPadding: 8,
          verticalPadding: 4,
          borderRadius: 6,
        );
      case YsDiscountBadgeSize.large:
        return const _BadgeConfig(
          fontSize: 16,
          horizontalPadding: 12,
          verticalPadding: 6,
          borderRadius: 8,
        );
    }
  }
}

class _BadgeConfig {
  final double fontSize;
  final double horizontalPadding;
  final double verticalPadding;
  final double borderRadius;

  const _BadgeConfig({
    required this.fontSize,
    required this.horizontalPadding,
    required this.verticalPadding,
    required this.borderRadius,
  });
}
