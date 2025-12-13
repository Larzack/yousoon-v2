import 'package:flutter/material.dart';

import 'package:yousoon/core/theme/app_colors.dart';
import 'package:yousoon/core/theme/app_spacing.dart';

/// Bouton primaire Yousoon
/// Fond Indian Gold (#E99B27), texte noir
class YsButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
  final bool isLoading;
  final bool isFullWidth;
  final YsButtonSize size;
  final YsButtonVariant variant;
  final IconData? icon;

  const YsButton({
    super.key,
    required this.label,
    this.onPressed,
    this.isLoading = false,
    this.isFullWidth = true,
    this.size = YsButtonSize.large,
    this.variant = YsButtonVariant.primary,
    this.icon,
  });

  @override
  Widget build(BuildContext context) {
    final isDisabled = onPressed == null || isLoading;

    return SizedBox(
      width: isFullWidth ? double.infinity : null,
      height: size.height,
      child: switch (variant) {
        YsButtonVariant.primary => _buildPrimaryButton(isDisabled),
        YsButtonVariant.secondary => _buildSecondaryButton(isDisabled),
        YsButtonVariant.tertiary => _buildTertiaryButton(isDisabled),
      },
    );
  }

  Widget _buildPrimaryButton(bool isDisabled) {
    return ElevatedButton(
      onPressed: isDisabled ? null : onPressed,
      style: ElevatedButton.styleFrom(
        backgroundColor: isDisabled ? AppColors.inactive : AppColors.primary,
        foregroundColor: isDisabled ? AppColors.textDisabled : AppColors.onPrimary,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
        ),
        elevation: 0,
      ),
      child: _buildContent(
        isDisabled ? AppColors.textDisabled : AppColors.onPrimary,
      ),
    );
  }

  Widget _buildSecondaryButton(bool isDisabled) {
    return OutlinedButton(
      onPressed: isDisabled ? null : onPressed,
      style: OutlinedButton.styleFrom(
        foregroundColor: isDisabled ? AppColors.inactive : AppColors.textPrimary,
        side: BorderSide(
          color: isDisabled ? AppColors.inactive : AppColors.textPrimary,
          width: 1.5,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
        ),
      ),
      child: _buildContent(
        isDisabled ? AppColors.inactive : AppColors.textPrimary,
      ),
    );
  }

  Widget _buildTertiaryButton(bool isDisabled) {
    return TextButton(
      onPressed: isDisabled ? null : onPressed,
      style: TextButton.styleFrom(
        foregroundColor: isDisabled ? AppColors.inactive : AppColors.textPrimary,
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: size.fontSize,
          fontWeight: FontWeight.w500,
          decoration: TextDecoration.underline,
          color: isDisabled ? AppColors.inactive : AppColors.textPrimary,
        ),
      ),
    );
  }

  Widget _buildContent(Color color) {
    if (isLoading) {
      return SizedBox(
        width: 20,
        height: 20,
        child: CircularProgressIndicator(
          strokeWidth: 2,
          valueColor: AlwaysStoppedAnimation<Color>(color),
        ),
      );
    }

    if (icon != null) {
      return Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 20),
          const SizedBox(width: 8),
          Text(
            label,
            style: TextStyle(
              fontSize: size.fontSize,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      );
    }

    return Text(
      label,
      style: TextStyle(
        fontSize: size.fontSize,
        fontWeight: FontWeight.w500,
      ),
    );
  }
}

enum YsButtonSize {
  small(height: 40, fontSize: 14),
  medium(height: 44, fontSize: 14),
  large(height: 50, fontSize: 16);

  final double height;
  final double fontSize;

  const YsButtonSize({required this.height, required this.fontSize});
}

enum YsButtonVariant {
  primary,
  secondary,
  tertiary,
}
