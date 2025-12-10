import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_animate/flutter_animate.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';

/// Écran de recherche d'offres
class SearchScreen extends ConsumerStatefulWidget {
  const SearchScreen({super.key});

  @override
  ConsumerState<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends ConsumerState<SearchScreen> {
  final TextEditingController _searchController = TextEditingController();
  final FocusNode _searchFocusNode = FocusNode();
  
  String _searchQuery = '';
  List<String> _recentSearches = [
    'Cocktails',
    'Restaurant italien',
    'Brunch Paris',
    'Happy hour',
  ];

  // Mock categories
  final List<Map<String, dynamic>> _categories = [
    {'name': 'Restaurants', 'icon': '🍽️', 'count': 42},
    {'name': 'Bars', 'icon': '🍸', 'count': 28},
    {'name': 'Cafés', 'icon': '☕', 'count': 15},
    {'name': 'Loisirs', 'icon': '🎯', 'count': 12},
    {'name': 'Bien-être', 'icon': '🧘', 'count': 8},
    {'name': 'Culture', 'icon': '🎨', 'count': 6},
  ];

  // Mock search results
  final List<Map<String, dynamic>> _searchResults = [
    {
      'id': '1',
      'title': '-30% sur les cocktails',
      'partner': 'Le Petit Bistrot',
      'category': 'Bars',
      'distance': '0.5 km',
      'image': 'https://picsum.photos/seed/search1/400/200',
      'discount': '-30%',
    },
    {
      'id': '2',
      'title': 'Menu dégustation -25%',
      'partner': 'Restaurant Le Gourmet',
      'category': 'Restaurants',
      'distance': '1.2 km',
      'image': 'https://picsum.photos/seed/search2/400/200',
      'discount': '-25%',
    },
    {
      'id': '3',
      'title': 'Cocktail offert',
      'partner': 'Bar L\'Éclipse',
      'category': 'Bars',
      'distance': '0.8 km',
      'image': 'https://picsum.photos/seed/search3/400/200',
      'discount': 'OFFERT',
    },
  ];

  @override
  void initState() {
    super.initState();
    // Auto-focus search field
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _searchFocusNode.requestFocus();
    });
  }

  @override
  void dispose() {
    _searchController.dispose();
    _searchFocusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        backgroundColor: AppColors.background,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: AppColors.textPrimary),
          onPressed: () => context.pop(),
        ),
        title: _buildSearchField(),
      ),
      body: _searchQuery.isEmpty
          ? _buildSuggestions()
          : _buildSearchResults(),
    );
  }

  Widget _buildSearchField() {
    return Container(
      height: 40,
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      ),
      child: TextField(
        controller: _searchController,
        focusNode: _searchFocusNode,
        style: AppTypography.bodyText,
        decoration: InputDecoration(
          hintText: 'Rechercher une offre, un lieu...',
          hintStyle: AppTypography.bodyText.copyWith(
            color: AppColors.textDisabled,
          ),
          prefixIcon: const Icon(Icons.search, color: AppColors.textDisabled),
          suffixIcon: _searchQuery.isNotEmpty
              ? IconButton(
                  icon: const Icon(Icons.close, color: AppColors.textSecondary),
                  onPressed: () {
                    _searchController.clear();
                    setState(() => _searchQuery = '');
                  },
                )
              : null,
          border: InputBorder.none,
          contentPadding: const EdgeInsets.symmetric(
            horizontal: AppSpacing.md,
            vertical: AppSpacing.sm,
          ),
        ),
        onChanged: (value) {
          setState(() => _searchQuery = value);
        },
        onSubmitted: (value) {
          if (value.isNotEmpty && !_recentSearches.contains(value)) {
            setState(() {
              _recentSearches.insert(0, value);
              if (_recentSearches.length > 10) {
                _recentSearches.removeLast();
              }
            });
          }
        },
      ),
    );
  }

  Widget _buildSuggestions() {
    return ListView(
      padding: const EdgeInsets.all(AppSpacing.xl),
      children: [
        // Recent searches
        if (_recentSearches.isNotEmpty) ...[
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Recherches récentes', style: AppTypography.headline3),
              TextButton(
                onPressed: () {
                  setState(() => _recentSearches.clear());
                },
                child: Text(
                  'Effacer',
                  style: AppTypography.bodyText.copyWith(
                    color: AppColors.textSecondary,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSpacing.sm),
          Wrap(
            spacing: AppSpacing.sm,
            runSpacing: AppSpacing.sm,
            children: _recentSearches.map((search) {
              return GestureDetector(
                onTap: () {
                  _searchController.text = search;
                  setState(() => _searchQuery = search);
                },
                child: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: AppSpacing.md,
                    vertical: AppSpacing.sm,
                  ),
                  decoration: BoxDecoration(
                    color: AppColors.surface,
                    borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      const Icon(
                        Icons.history,
                        size: 16,
                        color: AppColors.textSecondary,
                      ),
                      const SizedBox(width: AppSpacing.xs),
                      Text(
                        search,
                        style: AppTypography.bodyText.copyWith(
                          color: AppColors.textSecondary,
                        ),
                      ),
                    ],
                  ),
                ),
              );
            }).toList(),
          ),
          const SizedBox(height: AppSpacing.xxl),
        ],

        // Categories
        Text('Catégories', style: AppTypography.headline3),
        const SizedBox(height: AppSpacing.md),
        GridView.builder(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: 2,
            crossAxisSpacing: AppSpacing.md,
            mainAxisSpacing: AppSpacing.md,
            childAspectRatio: 2.5,
          ),
          itemCount: _categories.length,
          itemBuilder: (context, index) {
            final category = _categories[index];
            return _buildCategoryTile(category, index);
          },
        ),
        
        const SizedBox(height: AppSpacing.xxl),

        // Popular searches
        Text('Recherches populaires', style: AppTypography.headline3),
        const SizedBox(height: AppSpacing.md),
        _buildPopularSearchItem('Happy hour près de moi', Icons.local_bar),
        _buildPopularSearchItem('Brunch le dimanche', Icons.brunch_dining),
        _buildPopularSearchItem('Restaurant romantique', Icons.favorite),
        _buildPopularSearchItem('Sortie entre amis', Icons.groups),
      ],
    );
  }

  Widget _buildCategoryTile(Map<String, dynamic> category, int index) {
    return GestureDetector(
      onTap: () {
        _searchController.text = category['name'];
        setState(() => _searchQuery = category['name']);
      },
      child: Container(
        padding: const EdgeInsets.all(AppSpacing.md),
        decoration: BoxDecoration(
          color: AppColors.surface,
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        ),
        child: Row(
          children: [
            Text(
              category['icon'],
              style: const TextStyle(fontSize: 24),
            ),
            const SizedBox(width: AppSpacing.sm),
            Expanded(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    category['name'],
                    style: AppTypography.headline3,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  Text(
                    '${category['count']} offres',
                    style: AppTypography.caption.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ).animate().fadeIn(
        delay: Duration(milliseconds: index * 50),
      ),
    );
  }

  Widget _buildPopularSearchItem(String query, IconData icon) {
    return ListTile(
      onTap: () {
        _searchController.text = query;
        setState(() => _searchQuery = query);
      },
      contentPadding: EdgeInsets.zero,
      leading: Container(
        width: 40,
        height: 40,
        decoration: BoxDecoration(
          color: AppColors.primary.withOpacity(0.1),
          borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
        ),
        child: Icon(icon, color: AppColors.primary),
      ),
      title: Text(query, style: AppTypography.bodyText),
      trailing: const Icon(
        Icons.arrow_forward_ios,
        size: 16,
        color: AppColors.textSecondary,
      ),
    );
  }

  Widget _buildSearchResults() {
    // Filter results based on search query
    final filteredResults = _searchResults.where((result) {
      final query = _searchQuery.toLowerCase();
      return result['title'].toString().toLowerCase().contains(query) ||
          result['partner'].toString().toLowerCase().contains(query) ||
          result['category'].toString().toLowerCase().contains(query);
    }).toList();

    if (filteredResults.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.xl),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.search_off,
                size: 80,
                color: AppColors.inactive,
              ),
              const SizedBox(height: AppSpacing.lg),
              Text(
                'Aucun résultat',
                style: AppTypography.headline2,
              ),
              const SizedBox(height: AppSpacing.sm),
              Text(
                'Essayez avec d\'autres termes\nou explorez les catégories',
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

    return ListView.builder(
      padding: const EdgeInsets.all(AppSpacing.xl),
      itemCount: filteredResults.length + 1,
      itemBuilder: (context, index) {
        if (index == 0) {
          return Padding(
            padding: const EdgeInsets.only(bottom: AppSpacing.md),
            child: Text(
              '${filteredResults.length} résultat${filteredResults.length > 1 ? 's' : ''}',
              style: AppTypography.caption.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
          );
        }

        final result = filteredResults[index - 1];
        return _buildResultCard(result, index - 1);
      },
    );
  }

  Widget _buildResultCard(Map<String, dynamic> result, int index) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSpacing.md),
      child: GestureDetector(
        onTap: () => context.push('/offer/${result['id']}'),
        child: Container(
          decoration: BoxDecoration(
            color: AppColors.surface,
            borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
          ),
          clipBehavior: Clip.antiAlias,
          child: Row(
            children: [
              // Image
              Stack(
                children: [
                  CachedNetworkImage(
                    imageUrl: result['image'],
                    width: 100,
                    height: 100,
                    fit: BoxFit.cover,
                  ),
                  Positioned(
                    top: AppSpacing.xs,
                    left: AppSpacing.xs,
                    child: Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 6,
                        vertical: 2,
                      ),
                      decoration: BoxDecoration(
                        color: AppColors.primary,
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Text(
                        result['discount'],
                        style: AppTypography.caption.copyWith(
                          color: AppColors.onPrimary,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  ),
                ],
              ),

              // Content
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        result['title'],
                        style: AppTypography.headline3,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 4),
                      Text(
                        result['partner'],
                        style: AppTypography.bodyText.copyWith(
                          color: AppColors.primary,
                        ),
                      ),
                      const SizedBox(height: AppSpacing.xs),
                      Row(
                        children: [
                          Icon(
                            Icons.location_on,
                            size: 14,
                            color: AppColors.textSecondary,
                          ),
                          const SizedBox(width: 4),
                          Text(
                            result['distance'],
                            style: AppTypography.caption.copyWith(
                              color: AppColors.textSecondary,
                            ),
                          ),
                          const SizedBox(width: AppSpacing.md),
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6,
                              vertical: 2,
                            ),
                            decoration: BoxDecoration(
                              color: AppColors.surface,
                              borderRadius: BorderRadius.circular(4),
                              border: Border.all(color: AppColors.border),
                            ),
                            child: Text(
                              result['category'],
                              style: AppTypography.caption.copyWith(
                                color: AppColors.textSecondary,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),

              // Arrow
              Padding(
                padding: const EdgeInsets.only(right: AppSpacing.md),
                child: Icon(
                  Icons.arrow_forward_ios,
                  size: 16,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
        ),
      ).animate().fadeIn(
        delay: Duration(milliseconds: index * 50),
      ).slideX(begin: 0.1, end: 0),
    );
  }
}
