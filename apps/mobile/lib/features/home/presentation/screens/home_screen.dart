import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../shared/widgets/inputs/ys_text_field.dart';
import '../../../../shared/widgets/cards/offer_card.dart';

/// Écran d'accueil "Pour vous"
class HomeScreen extends ConsumerStatefulWidget {
  const HomeScreen({super.key});

  @override
  ConsumerState<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends ConsumerState<HomeScreen> {
  // Mock data pour démonstration
  final List<Map<String, dynamic>> _offers = [
    {
      'id': '1',
      'title': 'Happy Hour -50% sur les cocktails',
      'partnerName': 'Le Bar du Coin',
      'imageUrl': 'https://images.unsplash.com/photo-1514362545857-3bc16c4c7d1b?w=800',
      'discount': '-50%',
      'distance': '0.5 km',
      'rating': 4.5,
    },
    {
      'id': '2',
      'title': 'Menu découverte à prix réduit',
      'partnerName': 'Restaurant Gourmet',
      'imageUrl': 'https://images.unsplash.com/photo-1414235077428-338989a2e8c0?w=800',
      'discount': '-30%',
      'distance': '1.2 km',
      'rating': 4.8,
    },
    {
      'id': '3',
      'title': 'Séance de yoga en plein air',
      'partnerName': 'Zen Studio',
      'imageUrl': 'https://images.unsplash.com/photo-1544367567-0f2fcb009e0b?w=800',
      'discount': '-40%',
      'distance': '2.0 km',
      'rating': 4.9,
    },
    {
      'id': '4',
      'title': 'Escape Game pour 4 personnes',
      'partnerName': 'Escape Room Paris',
      'imageUrl': 'https://images.unsplash.com/photo-1587825140708-dfaf72ae4b04?w=800',
      'discount': '2 pour 1',
      'distance': '3.5 km',
      'rating': 4.7,
    },
  ];

  final Set<String> _favorites = {};

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: CustomScrollView(
          slivers: [
            // Header
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.all(AppSpacing.screenHorizontal),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Top bar
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Pour vous',
                              style: TextStyle(
                                fontSize: 28,
                                fontWeight: FontWeight.bold,
                                color: AppColors.textPrimary,
                              ),
                            ),
                            SizedBox(height: 4),
                            Row(
                              children: [
                                Icon(
                                  Icons.location_on,
                                  size: 14,
                                  color: AppColors.primary,
                                ),
                                SizedBox(width: 4),
                                Text(
                                  'Paris, France',
                                  style: TextStyle(
                                    fontSize: 14,
                                    color: AppColors.textSecondary,
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                        Row(
                          children: [
                            _HeaderButton(
                              icon: Icons.notifications_outlined,
                              onTap: () => context.push('/notifications'),
                            ),
                            const SizedBox(width: 8),
                            _HeaderButton(
                              icon: Icons.person_outline,
                              onTap: () => context.push('/profile'),
                            ),
                          ],
                        ),
                      ],
                    ),
                    const SizedBox(height: AppSpacing.lg),
                    // Search bar
                    GestureDetector(
                      onTap: () => context.push('/search'),
                      child: const YsSearchBar(
                        hint: 'Rechercher une sortie...',
                        readOnly: true,
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),
                    // Categories
                    SizedBox(
                      height: 40,
                      child: ListView(
                        scrollDirection: Axis.horizontal,
                        children: [
                          _CategoryChip(label: 'Tout', isSelected: true),
                          _CategoryChip(label: 'Bars'),
                          _CategoryChip(label: 'Restaurants'),
                          _CategoryChip(label: 'Activités'),
                          _CategoryChip(label: 'Bien-être'),
                          _CategoryChip(label: 'Culture'),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
            // Offers grid
            SliverPadding(
              padding: const EdgeInsets.symmetric(
                horizontal: AppSpacing.screenHorizontal,
              ),
              sliver: SliverGrid(
                gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisCount: 2,
                  mainAxisSpacing: AppSpacing.md,
                  crossAxisSpacing: AppSpacing.md,
                  childAspectRatio: 0.75,
                ),
                delegate: SliverChildBuilderDelegate(
                  (context, index) {
                    final offer = _offers[index % _offers.length];
                    return OfferCard(
                      id: offer['id'],
                      title: offer['title'],
                      partnerName: offer['partnerName'],
                      imageUrl: offer['imageUrl'],
                      discount: offer['discount'],
                      distance: offer['distance'],
                      rating: offer['rating'],
                      isFavorited: _favorites.contains(offer['id']),
                      onTap: () => context.push('/offer/${offer['id']}'),
                      onFavoriteTap: () {
                        setState(() {
                          if (_favorites.contains(offer['id'])) {
                            _favorites.remove(offer['id']);
                          } else {
                            _favorites.add(offer['id']);
                          }
                        });
                      },
                    );
                  },
                  childCount: 8,
                ),
              ),
            ),
            // Bottom padding
            const SliverToBoxAdapter(
              child: SizedBox(height: AppSpacing.xxl),
            ),
          ],
        ),
      ),
    );
  }
}

class _HeaderButton extends StatelessWidget {
  final IconData icon;
  final VoidCallback onTap;

  const _HeaderButton({
    required this.icon,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        width: 40,
        height: 40,
        decoration: BoxDecoration(
          color: AppColors.surface,
          shape: BoxShape.circle,
        ),
        child: Icon(
          icon,
          color: AppColors.textPrimary,
          size: 20,
        ),
      ),
    );
  }
}

class _CategoryChip extends StatelessWidget {
  final String label;
  final bool isSelected;

  const _CategoryChip({
    required this.label,
    this.isSelected = false,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(right: 8),
      child: Chip(
        label: Text(label),
        backgroundColor: isSelected ? AppColors.primary : AppColors.surface,
        labelStyle: TextStyle(
          color: isSelected ? AppColors.onPrimary : AppColors.textPrimary,
          fontSize: 12,
          fontWeight: FontWeight.w500,
        ),
        side: BorderSide.none,
        padding: const EdgeInsets.symmetric(horizontal: 8),
      ),
    );
  }
}
