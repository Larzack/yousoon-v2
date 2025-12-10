import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';

import '../../../core/theme/app_colors.dart';

/// Avatar utilisateur ou partenaire
class YsAvatar extends StatelessWidget {
  /// URL de l'image
  final String? imageUrl;

  /// Initiales à afficher si pas d'image
  final String? initials;

  /// Taille de l'avatar
  final double size;

  /// Couleur de fond (si pas d'image)
  final Color? backgroundColor;

  /// Afficher une bordure
  final bool showBorder;

  /// Couleur de la bordure
  final Color? borderColor;

  /// Badge de statut (optionnel)
  final YsAvatarBadge? badge;

  const YsAvatar({
    super.key,
    this.imageUrl,
    this.initials,
    this.size = 40,
    this.backgroundColor,
    this.showBorder = false,
    this.borderColor,
    this.badge,
  });

  /// Avatar petit
  const YsAvatar.small({
    super.key,
    this.imageUrl,
    this.initials,
    this.backgroundColor,
    this.badge,
  })  : size = 32,
        showBorder = false,
        borderColor = null;

  /// Avatar moyen
  const YsAvatar.medium({
    super.key,
    this.imageUrl,
    this.initials,
    this.backgroundColor,
    this.badge,
  })  : size = 48,
        showBorder = false,
        borderColor = null;

  /// Avatar large
  const YsAvatar.large({
    super.key,
    this.imageUrl,
    this.initials,
    this.backgroundColor,
    this.badge,
  })  : size = 80,
        showBorder = true,
        borderColor = null;

  /// Avatar extra-large (profil)
  const YsAvatar.xl({
    super.key,
    this.imageUrl,
    this.initials,
    this.backgroundColor,
    this.badge,
  })  : size = 120,
        showBorder = true,
        borderColor = null;

  @override
  Widget build(BuildContext context) {
    Widget avatar = Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: backgroundColor ?? AppColors.primary.withOpacity(0.2),
        border: showBorder
            ? Border.all(
                color: borderColor ?? AppColors.primary,
                width: 2,
              )
            : null,
      ),
      child: ClipOval(
        child: _buildContent(),
      ),
    );

    // Ajouter le badge si présent
    if (badge != null) {
      avatar = Stack(
        clipBehavior: Clip.none,
        children: [
          avatar,
          Positioned(
            right: 0,
            bottom: 0,
            child: _buildBadge(),
          ),
        ],
      );
    }

    return avatar;
  }

  Widget _buildContent() {
    if (imageUrl != null && imageUrl!.isNotEmpty) {
      return CachedNetworkImage(
        imageUrl: imageUrl!,
        fit: BoxFit.cover,
        placeholder: (context, url) => _buildPlaceholder(),
        errorWidget: (context, url, error) => _buildPlaceholder(),
      );
    }

    return _buildPlaceholder();
  }

  Widget _buildPlaceholder() {
    if (initials != null && initials!.isNotEmpty) {
      return Center(
        child: Text(
          initials!.toUpperCase(),
          style: TextStyle(
            color: AppColors.primary,
            fontSize: size * 0.4,
            fontWeight: FontWeight.bold,
          ),
        ),
      );
    }

    return Center(
      child: Icon(
        Icons.person_rounded,
        size: size * 0.5,
        color: AppColors.primary,
      ),
    );
  }

  Widget _buildBadge() {
    final badgeSize = size * 0.3;

    Color badgeColor;
    IconData? badgeIcon;

    switch (badge!) {
      case YsAvatarBadge.verified:
        badgeColor = AppColors.success;
        badgeIcon = Icons.check;
      case YsAvatarBadge.premium:
        badgeColor = AppColors.primary;
        badgeIcon = Icons.star_rounded;
      case YsAvatarBadge.online:
        badgeColor = AppColors.success;
        badgeIcon = null;
      case YsAvatarBadge.offline:
        badgeColor = AppColors.textDisabled;
        badgeIcon = null;
    }

    return Container(
      width: badgeSize,
      height: badgeSize,
      decoration: BoxDecoration(
        color: badgeColor,
        shape: BoxShape.circle,
        border: Border.all(
          color: AppColors.background,
          width: 2,
        ),
      ),
      child: badgeIcon != null
          ? Icon(
              badgeIcon,
              size: badgeSize * 0.6,
              color: Colors.white,
            )
          : null,
    );
  }
}

/// Types de badge pour l'avatar
enum YsAvatarBadge {
  verified,
  premium,
  online,
  offline,
}

/// Groupe d'avatars (ex: participants)
class YsAvatarGroup extends StatelessWidget {
  final List<String?> imageUrls;
  final double size;
  final int maxVisible;
  final int? total;

  const YsAvatarGroup({
    super.key,
    required this.imageUrls,
    this.size = 32,
    this.maxVisible = 4,
    this.total,
  });

  @override
  Widget build(BuildContext context) {
    final visibleCount = imageUrls.length > maxVisible ? maxVisible : imageUrls.length;
    final remaining = (total ?? imageUrls.length) - visibleCount;

    return SizedBox(
      width: size + (visibleCount - 1) * (size * 0.6) + (remaining > 0 ? size * 0.6 : 0),
      height: size,
      child: Stack(
        children: [
          // Avatars visibles
          for (var i = 0; i < visibleCount; i++)
            Positioned(
              left: i * (size * 0.6),
              child: Container(
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  border: Border.all(
                    color: AppColors.background,
                    width: 2,
                  ),
                ),
                child: YsAvatar(
                  imageUrl: imageUrls[i],
                  size: size,
                ),
              ),
            ),

          // Compteur des restants
          if (remaining > 0)
            Positioned(
              left: visibleCount * (size * 0.6),
              child: Container(
                width: size,
                height: size,
                decoration: BoxDecoration(
                  color: AppColors.cardBackground,
                  shape: BoxShape.circle,
                  border: Border.all(
                    color: AppColors.background,
                    width: 2,
                  ),
                ),
                child: Center(
                  child: Text(
                    '+$remaining',
                    style: TextStyle(
                      color: AppColors.textPrimary,
                      fontSize: size * 0.35,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }
}
