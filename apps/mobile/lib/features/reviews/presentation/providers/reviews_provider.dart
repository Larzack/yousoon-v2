import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../data/models/review_model.dart';
import '../../data/repositories/reviews_repository.dart';

part 'reviews_provider.freezed.dart';

/// État des avis pour une offre
@freezed
class ReviewsState with _$ReviewsState {
  const factory ReviewsState({
    @Default([]) List<ReviewModel> reviews,
    @Default(false) bool isLoading,
    @Default(false) bool isLoadingMore,
    String? error,
    @Default(1) int currentPage,
    @Default(false) bool hasMore,
    @Default(0) int totalCount,
    double? averageRating,
    @Default({}) Map<int, int> ratingDistribution,
    @Default(ReviewSortBy.newest) ReviewSortBy sortBy,
  }) = _ReviewsState;
}

/// Notifier pour les avis d'une offre
class OfferReviewsNotifier extends StateNotifier<ReviewsState> {
  final ReviewsRepository _repository;
  final String offerId;

  OfferReviewsNotifier(this._repository, this.offerId) : super(const ReviewsState());

  /// Charger les avis
  Future<void> loadReviews({ReviewSortBy? sortBy}) async {
    state = state.copyWith(
      isLoading: true,
      error: null,
      sortBy: sortBy ?? state.sortBy,
    );

    try {
      final result = await _repository.getOfferReviews(
        offerId: offerId,
        sortBy: sortBy ?? state.sortBy,
        page: 1,
      );

      state = state.copyWith(
        reviews: result.reviews,
        isLoading: false,
        currentPage: result.page,
        hasMore: result.hasMore,
        totalCount: result.totalCount,
        averageRating: result.averageRating,
        ratingDistribution: result.ratingDistribution,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Charger plus d'avis
  Future<void> loadMore() async {
    if (state.isLoadingMore || !state.hasMore) return;

    state = state.copyWith(isLoadingMore: true);

    try {
      final nextPage = state.currentPage + 1;
      final result = await _repository.getOfferReviews(
        offerId: offerId,
        sortBy: state.sortBy,
        page: nextPage,
      );

      state = state.copyWith(
        reviews: [...state.reviews, ...result.reviews],
        isLoadingMore: false,
        currentPage: result.page,
        hasMore: result.hasMore,
      );
    } catch (e) {
      state = state.copyWith(
        isLoadingMore: false,
        error: e.toString(),
      );
    }
  }

  /// Changer le tri
  Future<void> changeSortBy(ReviewSortBy sortBy) async {
    await loadReviews(sortBy: sortBy);
  }

  /// Marquer un avis comme utile
  Future<void> markAsHelpful(String reviewId) async {
    try {
      await _repository.markAsHelpful(reviewId);
      
      // Mettre à jour localement
      state = state.copyWith(
        reviews: state.reviews.map((r) {
          if (r.id == reviewId) {
            return r.copyWith(helpfulCount: r.helpfulCount + 1);
          }
          return r;
        }).toList(),
      );
    } catch (e) {
      // Ignorer les erreurs pour cette action
    }
  }
}

/// Provider famille pour les avis d'une offre
final offerReviewsProvider = StateNotifierProvider.family<
    OfferReviewsNotifier, ReviewsState, String>((ref, offerId) {
  final repository = ref.watch(reviewsRepositoryProvider);
  return OfferReviewsNotifier(repository, offerId);
});

/// État de création d'un avis
@freezed
class CreateReviewState with _$CreateReviewState {
  const factory CreateReviewState.initial() = _Initial;
  const factory CreateReviewState.loading() = _Loading;
  const factory CreateReviewState.success(ReviewModel review) = _Success;
  const factory CreateReviewState.error(String message) = _Error;
}

/// Notifier pour la création d'un avis
class CreateReviewNotifier extends StateNotifier<CreateReviewState> {
  final ReviewsRepository _repository;

  CreateReviewNotifier(this._repository) : super(const CreateReviewState.initial());

  /// Créer un avis
  Future<void> createReview(CreateReviewParams params) async {
    state = const CreateReviewState.loading();

    try {
      final review = await _repository.createReview(params);
      state = CreateReviewState.success(review);
    } catch (e) {
      state = CreateReviewState.error(e.toString());
    }
  }

  /// Réinitialiser l'état
  void reset() {
    state = const CreateReviewState.initial();
  }
}

/// Provider pour la création d'un avis
final createReviewProvider =
    StateNotifierProvider<CreateReviewNotifier, CreateReviewState>((ref) {
  final repository = ref.watch(reviewsRepositoryProvider);
  return CreateReviewNotifier(repository);
});

/// Provider pour vérifier si on peut laisser un avis
final canReviewOfferProvider = FutureProvider.family<bool, String>((ref, offerId) async {
  final repository = ref.watch(reviewsRepositoryProvider);
  return repository.canReviewOffer(offerId);
});

/// Provider pour mes avis
final myReviewsProvider = FutureProvider<List<ReviewModel>>((ref) async {
  final repository = ref.watch(reviewsRepositoryProvider);
  final result = await repository.getMyReviews();
  return result.reviews;
});
