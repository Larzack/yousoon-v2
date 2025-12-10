import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_animate/flutter_animate.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';
import '../../../../shared/widgets/layouts/ys_scaffold.dart';

/// Écran des favoris
class FavoritesScreen extends ConsumerWidget {
  const FavoritesScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // TODO: Connect to actual data via Riverpod
    final favorites = List.generate(5, (i) => i);
    
    return YsScaffold(
      title: 'Mes favoris',
      showBackButton: false,
      body: favorites.isEmpty
          ? _buildEmptyState()
          : ListView.separated(
              padding: const EdgeInsets.all(AppSpacing.xl),
              itemCount: favorites.length,
              separatorBuilder: (_, __) => const SizedBox(height: AppSpacing.md),
              itemBuilder: (context, index) {
                return _FavoriteItem(
                  title: 'Offre ${index + 1}',
                  partnerName: 'Partenaire ${index + 1}',
                  discount: '-${(index + 1) * 10}%',
                  imageUrl: 'https://picsum.photos/300/200?random=$index',
                  validUntil: 'Valide jusqu\'au 31 déc. 2025',
                  onTap: () => context.push('/offers/offer_$index'),
                  onRemove: () {
                    // TODO: Remove from favorites
                  },
                ).animate().fadeIn(delay: Duration(milliseconds: index * 100));
              },
            ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.xl),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: AppColors.surface,
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.favorite_border,
                color: AppColors.inactive,
                size: 40,
              ),
            ),
            const SizedBox(height: AppSpacing.lg),
            Text(
              'Aucun favori',
              style: AppTypography.headline2,
            ),
            const SizedBox(height: AppSpacing.sm),
            Text(
              'Ajoutez des offres à vos favoris pour les retrouver facilement.',
              style: AppTypography.bodyText.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }
}

class _FavoriteItem extends StatelessWidget {
  final String title;
  final String partnerName;
  final String discount;
  final String imageUrl;
  final String validUntil;
  final VoidCallback onTap;
  final VoidCallback onRemove;

  const _FavoriteItem({
    required this.title,
    required this.partnerName,
    required this.discount,
    required this.imageUrl,
    required this.validUntil,
    required this.onTap,
    required this.onRemove,
  });

  @override
  Widget build(BuildContext context) {
    return Dismissible(
      key: Key(title),
      direction: DismissDirection.endToStart,
      background: Container(
        alignment: Alignment.centerRight,
        padding: const EdgeInsets.only(right: AppSpacing.lg),
        decoration: BoxDecoration(
          color: AppColors.error,
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        ),
        child: const Icon(Icons.delete, color: Colors.white),
      ),
      onDismissed: (_) => onRemove(),
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          decoration: BoxDecoration(
            color: AppColors.surface,
            borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
            border: Border.all(color: AppColors.border),
          ),
          clipBehavior: Clip.antiAlias,
          child: Row(
            children: [
              // Image
              SizedBox(
                width: 100,
                height: 100,
                child: Stack(
                  fit: StackFit.expand,
                  children: [
                    CachedNetworkImage(
                      imageUrl: imageUrl,
                      fit: BoxFit.cover,
                    ),
                    Positioned(
                      top: AppSpacing.xs,
                      left: AppSpacing.xs,
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: AppSpacing.xs,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color: AppColors.primary,
                          borderRadius: BorderRadius.circular(AppSpacing.radiusXs),
                        ),
                        child: Text(
                          discount,
                          style: AppTypography.caption.copyWith(
                            color: AppColors.onPrimary,
                            fontWeight: FontWeight.bold,
                            fontSize: 10,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              
              // Content
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        title,
                        style: AppTypography.headline3,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 2),
                      Text(
                        partnerName,
                        style: AppTypography.caption.copyWith(
                          color: AppColors.textSecondary,
                        ),
                      ),
                      const SizedBox(height: AppSpacing.sm),
                      Row(
                        children: [
                          const Icon(
                            Icons.schedule,
                            size: 12,
                            color: AppColors.textSecondary,
                          ),
                          const SizedBox(width: 4),
                          Expanded(
                            child: Text(
                              validUntil,
                              style: AppTypography.caption.copyWith(
                                color: AppColors.textSecondary,
                                fontSize: 11,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
              
              // Favorite button
              IconButton(
                icon: const Icon(Icons.favorite, color: AppColors.error),
                onPressed: onRemove,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
