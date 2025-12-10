import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';
import '../../../../shared/widgets/layouts/ys_scaffold.dart';
import '../../../../shared/widgets/buttons/ys_button.dart';

/// Écran principal des offres
/// Affiche la liste des offres avec filtres et recherche
class OffersScreen extends ConsumerStatefulWidget {
  const OffersScreen({super.key});

  @override
  ConsumerState<OffersScreen> createState() => _OffersScreenState();
}

class _OffersScreenState extends ConsumerState<OffersScreen> {
  final TextEditingController _searchController = TextEditingController();
  String? _selectedCategory;
  
  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return YsScaffold(
      title: 'Découvrir',
      showBackButton: false,
      actions: [
        IconButton(
          icon: const Icon(Icons.filter_list, color: AppColors.textPrimary),
          onPressed: _showFilters,
        ),
      ],
      body: Column(
        children: [
          // Search bar
          _buildSearchBar(),
          
          const SizedBox(height: AppSpacing.md),
          
          // Category chips
          _buildCategoryChips(),
          
          const SizedBox(height: AppSpacing.md),
          
          // Offers grid
          Expanded(
            child: _buildOffersGrid(),
          ),
        ],
      ),
    );
  }

  Widget _buildSearchBar() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: AppSpacing.xl),
      child: Container(
        decoration: BoxDecoration(
          color: AppColors.surface,
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
          border: Border.all(color: AppColors.border),
        ),
        child: TextField(
          controller: _searchController,
          style: AppTypography.bodyText,
          decoration: InputDecoration(
            hintText: 'Rechercher une offre...',
            hintStyle: AppTypography.bodyText.copyWith(color: AppColors.textDisabled),
            prefixIcon: const Icon(Icons.search, color: AppColors.textDisabled),
            border: InputBorder.none,
            contentPadding: const EdgeInsets.symmetric(
              horizontal: AppSpacing.md,
              vertical: AppSpacing.sm,
            ),
          ),
          onChanged: (value) {
            // TODO: Implement search
          },
        ),
      ),
    );
  }

  Widget _buildCategoryChips() {
    final categories = [
      {'id': null, 'name': 'Tout', 'emoji': '✨'},
      {'id': 'sport', 'name': 'Sport', 'emoji': '🎾'},
      {'id': 'convivialite', 'name': 'Convivialité', 'emoji': '🥂'},
      {'id': 'culture', 'name': 'Arts & Culture', 'emoji': '🎨'},
      {'id': 'voyage', 'name': 'Voyage', 'emoji': '🌎'},
      {'id': 'concert', 'name': 'Concert', 'emoji': '🎺'},
      {'id': 'bien-etre', 'name': 'Bien-être', 'emoji': '🧘'},
      {'id': 'nature', 'name': 'Nature', 'emoji': '🌱'},
      {'id': 'gastronomie', 'name': 'Gastronomie', 'emoji': '🍴'},
    ];

    return SizedBox(
      height: 40,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: AppSpacing.xl),
        itemCount: categories.length,
        separatorBuilder: (_, __) => const SizedBox(width: AppSpacing.sm),
        itemBuilder: (context, index) {
          final category = categories[index];
          final isSelected = _selectedCategory == category['id'];
          
          return FilterChip(
            selected: isSelected,
            label: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(category['emoji'] as String),
                const SizedBox(width: 4),
                Text(category['name'] as String),
              ],
            ),
            labelStyle: AppTypography.caption.copyWith(
              color: isSelected ? AppColors.onPrimary : AppColors.textPrimary,
            ),
            backgroundColor: AppColors.surface,
            selectedColor: AppColors.primary,
            side: BorderSide(
              color: isSelected ? AppColors.primary : AppColors.border,
            ),
            onSelected: (selected) {
              setState(() {
                _selectedCategory = selected ? category['id'] as String? : null;
              });
            },
          );
        },
      ),
    );
  }

  Widget _buildOffersGrid() {
    // TODO: Connect to actual data via Riverpod
    return GridView.builder(
      padding: const EdgeInsets.all(AppSpacing.xl),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        crossAxisSpacing: AppSpacing.md,
        mainAxisSpacing: AppSpacing.md,
        childAspectRatio: 0.75,
      ),
      itemCount: 10, // Placeholder
      itemBuilder: (context, index) {
        return _OfferGridItem(
          title: 'Offre ${index + 1}',
          partnerName: 'Partenaire ${index + 1}',
          discount: '-${(index + 1) * 5}%',
          imageUrl: 'https://picsum.photos/200/300?random=$index',
          onTap: () {
            context.push('/offers/offer_$index');
          },
        ).animate().fadeIn(delay: Duration(milliseconds: index * 100));
      },
    );
  }

  void _showFilters() {
    showModalBottomSheet(
      context: context,
      backgroundColor: AppColors.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(AppSpacing.radiusLg)),
      ),
      builder: (context) => const _FiltersBottomSheet(),
    );
  }
}

class _OfferGridItem extends StatelessWidget {
  final String title;
  final String partnerName;
  final String discount;
  final String imageUrl;
  final VoidCallback onTap;

  const _OfferGridItem({
    required this.title,
    required this.partnerName,
    required this.discount,
    required this.imageUrl,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
          border: Border.all(color: AppColors.border),
        ),
        clipBehavior: Clip.antiAlias,
        child: Stack(
          fit: StackFit.expand,
          children: [
            // Image
            CachedNetworkImage(
              imageUrl: imageUrl,
              fit: BoxFit.cover,
              placeholder: (context, url) => Container(
                color: AppColors.surface,
                child: const Center(
                  child: CircularProgressIndicator(color: AppColors.primary),
                ),
              ),
              errorWidget: (context, url, error) => Container(
                color: AppColors.surface,
                child: const Icon(Icons.error, color: AppColors.error),
              ),
            ),
            
            // Gradient overlay
            Container(
              decoration: const BoxDecoration(
                gradient: AppColors.cardOverlayGradient,
              ),
            ),
            
            // Discount badge
            Positioned(
              top: AppSpacing.sm,
              right: AppSpacing.sm,
              child: Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: AppSpacing.sm,
                  vertical: AppSpacing.xs,
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
                  ),
                ),
              ),
            ),
            
            // Content
            Positioned(
              left: AppSpacing.sm,
              right: AppSpacing.sm,
              bottom: AppSpacing.sm,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    title,
                    style: AppTypography.headline3,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 2),
                  Text(
                    partnerName,
                    style: AppTypography.caption.copyWith(
                      color: AppColors.textSecondary,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _FiltersBottomSheet extends StatelessWidget {
  const _FiltersBottomSheet();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.xl),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Filtres', style: AppTypography.headline1),
              IconButton(
                icon: const Icon(Icons.close, color: AppColors.textPrimary),
                onPressed: () => Navigator.pop(context),
              ),
            ],
          ),
          
          const SizedBox(height: AppSpacing.lg),
          
          // Distance slider
          Text('Distance maximale', style: AppTypography.headline3),
          const SizedBox(height: AppSpacing.sm),
          Slider(
            value: 10,
            min: 1,
            max: 50,
            divisions: 49,
            activeColor: AppColors.primary,
            inactiveColor: AppColors.inactive,
            label: '10 km',
            onChanged: (value) {
              // TODO: Update filter
            },
          ),
          
          const SizedBox(height: AppSpacing.lg),
          
          // Discount type
          Text('Type de réduction', style: AppTypography.headline3),
          const SizedBox(height: AppSpacing.sm),
          Wrap(
            spacing: AppSpacing.sm,
            children: [
              ChoiceChip(
                label: const Text('Tous'),
                selected: true,
                onSelected: (_) {},
              ),
              ChoiceChip(
                label: const Text('Pourcentage'),
                selected: false,
                onSelected: (_) {},
              ),
              ChoiceChip(
                label: const Text('Montant fixe'),
                selected: false,
                onSelected: (_) {},
              ),
            ],
          ),
          
          const SizedBox(height: AppSpacing.xl),
          
          // Apply button
          YsButton(
            label: 'Appliquer les filtres',
            onPressed: () => Navigator.pop(context),
          ),
          
          const SizedBox(height: AppSpacing.lg),
        ],
      ),
    );
  }
}
