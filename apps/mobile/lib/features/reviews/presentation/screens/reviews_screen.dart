import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../shared/widgets/feedback/ys_empty_state.dart';
import '../../../../shared/widgets/feedback/ys_loader.dart';
import '../../../../shared/widgets/feedback/ys_rating.dart';
import '../../../../shared/widgets/images/ys_avatar.dart';
import '../../../../shared/widgets/layouts/ys_scaffold.dart';
import '../../data/models/review_model.dart';
import '../providers/reviews_provider.dart';

/// Écran des avis d'une offre
class ReviewsScreen extends ConsumerStatefulWidget {
  final String offerId;
  final String offerTitle;

  const ReviewsScreen({
    super.key,
    required this.offerId,
    required this.offerTitle,
  });

  @override
  ConsumerState<ReviewsScreen> createState() => _ReviewsScreenState();
}

class _ReviewsScreenState extends ConsumerState<ReviewsScreen> {
  @override
  void initState() {
    super.initState();
    // Charger les avis au démarrage
    Future.microtask(() {
      ref.read(offerReviewsProvider(widget.offerId).notifier).loadReviews();
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(offerReviewsProvider(widget.offerId));

    return YsScaffold(
      title: 'Avis',
      body: state.isLoading
          ? const YsLoadingCenter()
          : state.error != null
              ? YsErrorState(
                  message: state.error!,
                  onRetry: () => ref
                      .read(offerReviewsProvider(widget.offerId).notifier)
                      .loadReviews(),
                )
              : state.reviews.isEmpty
                  ? const YsEmptyState(
                      icon: Icons.rate_review_outlined,
                      title: 'Aucun avis',
                      message:
                          'Soyez le premier à donner votre avis sur cette offre !',
                    )
                  : _buildContent(state),
    );
  }

  Widget _buildContent(ReviewsState state) {
    return CustomScrollView(
      slivers: [
        // En-tête avec résumé des notes
        SliverToBoxAdapter(
          child: _buildRatingSummary(state),
        ),

        // Tri
        SliverToBoxAdapter(
          child: _buildSortSelector(state),
        ),

        // Liste des avis
        SliverPadding(
          padding: const EdgeInsets.all(AppSpacing.lg),
          sliver: SliverList(
            delegate: SliverChildBuilderDelegate(
              (context, index) {
                if (index >= state.reviews.length) {
                  // Loader de pagination
                  if (state.hasMore) {
                    return const Padding(
                      padding: EdgeInsets.all(AppSpacing.lg),
                      child: YsLoadingCenter(),
                    );
                  }
                  return null;
                }

                return _ReviewCard(review: state.reviews[index]);
              },
              childCount: state.reviews.length + (state.hasMore ? 1 : 0),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildRatingSummary(ReviewsState state) {
    final avgRating = state.averageRating ?? 0.0;
    final totalReviews = state.totalCount;
    final distribution = state.ratingDistribution;

    return Container(
      padding: const EdgeInsets.all(AppSpacing.lg),
      color: AppColors.cardBackground,
      child: Row(
        children: [
          // Note moyenne
          Expanded(
            flex: 2,
            child: Column(
              children: [
                Text(
                  avgRating.toStringAsFixed(1),
                  style: const TextStyle(
                    fontSize: 48,
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
                ),
                const SizedBox(height: AppSpacing.sm),
                YsRating.medium(rating: avgRating),
                const SizedBox(height: AppSpacing.xs),
                Text(
                  '$totalReviews avis',
                  style: const TextStyle(
                    fontSize: 14,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),

          // Distribution
          Expanded(
            flex: 3,
            child: Column(
              children: List.generate(5, (index) {
                final rating = 5 - index;
                final count = distribution[rating] ?? 0;
                final percentage = totalReviews > 0 ? count / totalReviews : 0.0;

                return Padding(
                  padding: const EdgeInsets.symmetric(vertical: 2),
                  child: Row(
                    children: [
                      Text(
                        '$rating',
                        style: const TextStyle(
                          fontSize: 12,
                          color: AppColors.textSecondary,
                        ),
                      ),
                      const SizedBox(width: 4),
                      const Icon(
                        Icons.star_rounded,
                        size: 12,
                        color: AppColors.primary,
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: ClipRRect(
                          borderRadius: BorderRadius.circular(4),
                          child: LinearProgressIndicator(
                            value: percentage,
                            backgroundColor: AppColors.textDisabled.withOpacity(0.2),
                            valueColor: const AlwaysStoppedAnimation<Color>(
                              AppColors.primary,
                            ),
                            minHeight: 8,
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      SizedBox(
                        width: 30,
                        child: Text(
                          '$count',
                          style: const TextStyle(
                            fontSize: 12,
                            color: AppColors.textSecondary,
                          ),
                          textAlign: TextAlign.end,
                        ),
                      ),
                    ],
                  ),
                );
              }),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSortSelector(ReviewsState state) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSpacing.lg,
        vertical: AppSpacing.md,
      ),
      child: Row(
        children: [
          const Text(
            'Trier par',
            style: TextStyle(
              color: AppColors.textSecondary,
              fontSize: 14,
            ),
          ),
          const SizedBox(width: AppSpacing.md),
          DropdownButton<ReviewSortBy>(
            value: state.sortBy,
            dropdownColor: AppColors.cardBackground,
            style: const TextStyle(color: AppColors.textPrimary),
            underline: const SizedBox.shrink(),
            items: const [
              DropdownMenuItem(
                value: ReviewSortBy.newest,
                child: Text('Plus récents'),
              ),
              DropdownMenuItem(
                value: ReviewSortBy.oldest,
                child: Text('Plus anciens'),
              ),
              DropdownMenuItem(
                value: ReviewSortBy.highestRating,
                child: Text('Meilleures notes'),
              ),
              DropdownMenuItem(
                value: ReviewSortBy.lowestRating,
                child: Text('Notes les plus basses'),
              ),
              DropdownMenuItem(
                value: ReviewSortBy.mostHelpful,
                child: Text('Plus utiles'),
              ),
            ],
            onChanged: (value) {
              if (value != null) {
                ref
                    .read(offerReviewsProvider(widget.offerId).notifier)
                    .changeSortBy(value);
              }
            },
          ),
        ],
      ),
    );
  }
}

/// Card d'un avis
class _ReviewCard extends ConsumerWidget {
  final ReviewModel review;

  const _ReviewCard({required this.review});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppSpacing.lg),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.cardBackground,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // En-tête
          Row(
            children: [
              YsAvatar.small(
                imageUrl: review.user.avatar,
                initials: review.user.firstName.isNotEmpty
                    ? review.user.firstName[0]
                    : '?',
                badge: review.isVerifiedPurchase ? YsAvatarBadge.verified : null,
              ),
              const SizedBox(width: AppSpacing.md),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      review.user.firstName,
                      style: const TextStyle(
                        fontWeight: FontWeight.w600,
                        color: AppColors.textPrimary,
                      ),
                    ),
                    Text(
                      _formatDate(review.createdAt),
                      style: const TextStyle(
                        fontSize: 12,
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
              YsRating.small(rating: review.rating.toDouble()),
            ],
          ),

          const SizedBox(height: AppSpacing.md),

          // Badge achat vérifié
          if (review.isVerifiedPurchase)
            Container(
              margin: const EdgeInsets.only(bottom: AppSpacing.sm),
              padding: const EdgeInsets.symmetric(
                horizontal: AppSpacing.sm,
                vertical: 2,
              ),
              decoration: BoxDecoration(
                color: AppColors.success.withOpacity(0.1),
                borderRadius: BorderRadius.circular(4),
              ),
              child: const Text(
                '✓ Achat vérifié',
                style: TextStyle(
                  fontSize: 11,
                  color: AppColors.success,
                ),
              ),
            ),

          // Titre
          if (review.title != null && review.title!.isNotEmpty) ...[
            Text(
              review.title!,
              style: const TextStyle(
                fontWeight: FontWeight.w600,
                fontSize: 16,
                color: AppColors.textPrimary,
              ),
            ),
            const SizedBox(height: AppSpacing.sm),
          ],

          // Contenu
          if (review.content != null && review.content!.isNotEmpty)
            Text(
              review.content!,
              style: const TextStyle(
                color: AppColors.textSecondary,
                height: 1.5,
              ),
            ),

          // Images
          if (review.images.isNotEmpty) ...[
            const SizedBox(height: AppSpacing.md),
            SizedBox(
              height: 80,
              child: ListView.separated(
                scrollDirection: Axis.horizontal,
                itemCount: review.images.length,
                separatorBuilder: (_, __) => const SizedBox(width: AppSpacing.sm),
                itemBuilder: (context, index) {
                  return ClipRRect(
                    borderRadius: BorderRadius.circular(8),
                    child: Image.network(
                      review.images[index],
                      width: 80,
                      height: 80,
                      fit: BoxFit.cover,
                    ),
                  );
                },
              ),
            ),
          ],

          const SizedBox(height: AppSpacing.md),

          // Actions
          Row(
            children: [
              // Bouton utile
              TextButton.icon(
                onPressed: () {
                  // TODO: Appeler l'API pour marquer comme utile
                },
                icon: const Icon(
                  Icons.thumb_up_outlined,
                  size: 18,
                  color: AppColors.textSecondary,
                ),
                label: Text(
                  'Utile (${review.helpfulCount})',
                  style: const TextStyle(
                    color: AppColors.textSecondary,
                    fontSize: 13,
                  ),
                ),
              ),
              const Spacer(),
              // Bouton signaler
              IconButton(
                onPressed: () {
                  _showReportDialog(context, ref);
                },
                icon: const Icon(
                  Icons.flag_outlined,
                  size: 18,
                  color: AppColors.textDisabled,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);

    if (diff.inDays == 0) {
      return 'Aujourd\'hui';
    } else if (diff.inDays == 1) {
      return 'Hier';
    } else if (diff.inDays < 7) {
      return 'Il y a ${diff.inDays} jours';
    } else if (diff.inDays < 30) {
      return 'Il y a ${(diff.inDays / 7).floor()} semaines';
    } else if (diff.inDays < 365) {
      return 'Il y a ${(diff.inDays / 30).floor()} mois';
    } else {
      return 'Il y a ${(diff.inDays / 365).floor()} ans';
    }
  }

  void _showReportDialog(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: AppColors.cardBackground,
        title: const Text(
          'Signaler cet avis',
          style: TextStyle(color: AppColors.textPrimary),
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            _ReportOption(
              label: 'Contenu inapproprié',
              onTap: () {
                Navigator.pop(context);
                // TODO: Appeler l'API
              },
            ),
            _ReportOption(
              label: 'Faux avis',
              onTap: () {
                Navigator.pop(context);
              },
            ),
            _ReportOption(
              label: 'Spam',
              onTap: () {
                Navigator.pop(context);
              },
            ),
            _ReportOption(
              label: 'Autre',
              onTap: () {
                Navigator.pop(context);
              },
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text(
              'Annuler',
              style: TextStyle(color: AppColors.textSecondary),
            ),
          ),
        ],
      ),
    );
  }
}

class _ReportOption extends StatelessWidget {
  final String label;
  final VoidCallback onTap;

  const _ReportOption({
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      title: Text(
        label,
        style: const TextStyle(color: AppColors.textPrimary),
      ),
      onTap: onTap,
      contentPadding: EdgeInsets.zero,
    );
  }
}
