import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../models/review_model.dart';

/// Interface du repository des avis
abstract class ReviewsRepository {
  /// Obtenir les avis d'une offre
  Future<PaginatedReviewsResult> getOfferReviews({
    required String offerId,
    ReviewSortBy sortBy = ReviewSortBy.newest,
    int page = 1,
    int limit = 20,
  });

  /// Obtenir les avis d'un partenaire
  Future<PaginatedReviewsResult> getPartnerReviews({
    required String partnerId,
    ReviewSortBy sortBy = ReviewSortBy.newest,
    int page = 1,
    int limit = 20,
  });

  /// Obtenir les avis de l'utilisateur courant
  Future<PaginatedReviewsResult> getMyReviews({
    int page = 1,
    int limit = 20,
  });

  /// Obtenir un avis par ID
  Future<ReviewModel> getReviewById(String id);

  /// Créer un avis
  Future<ReviewModel> createReview(CreateReviewParams params);

  /// Modifier un avis
  Future<ReviewModel> updateReview({
    required String reviewId,
    int? rating,
    String? title,
    String? content,
  });

  /// Supprimer un avis
  Future<void> deleteReview(String reviewId);

  /// Marquer un avis comme utile
  Future<void> markAsHelpful(String reviewId);

  /// Signaler un avis
  Future<void> reportReview(String reviewId, String reason);

  /// Vérifier si l'utilisateur peut laisser un avis sur une offre
  Future<bool> canReviewOffer(String offerId);
}

/// Implémentation du repository des avis
class ReviewsRepositoryImpl implements ReviewsRepository {
  // final Client _graphqlClient;

  ReviewsRepositoryImpl({
    // required Client graphqlClient,
  });
  // : _graphqlClient = graphqlClient;

  @override
  Future<PaginatedReviewsResult> getOfferReviews({
    required String offerId,
    ReviewSortBy sortBy = ReviewSortBy.newest,
    int page = 1,
    int limit = 20,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get offer reviews not implemented');
  }

  @override
  Future<PaginatedReviewsResult> getPartnerReviews({
    required String partnerId,
    ReviewSortBy sortBy = ReviewSortBy.newest,
    int page = 1,
    int limit = 20,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get partner reviews not implemented');
  }

  @override
  Future<PaginatedReviewsResult> getMyReviews({
    int page = 1,
    int limit = 20,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get my reviews not implemented');
  }

  @override
  Future<ReviewModel> getReviewById(String id) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get review by id not implemented');
  }

  @override
  Future<ReviewModel> createReview(CreateReviewParams params) async {
    // TODO: Implémenter avec GraphQL + upload images
    throw UnimplementedError('Create review not implemented');
  }

  @override
  Future<ReviewModel> updateReview({
    required String reviewId,
    int? rating,
    String? title,
    String? content,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Update review not implemented');
  }

  @override
  Future<void> deleteReview(String reviewId) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Delete review not implemented');
  }

  @override
  Future<void> markAsHelpful(String reviewId) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Mark as helpful not implemented');
  }

  @override
  Future<void> reportReview(String reviewId, String reason) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Report review not implemented');
  }

  @override
  Future<bool> canReviewOffer(String offerId) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Can review offer not implemented');
  }
}

/// Provider du repository des avis
final reviewsRepositoryProvider = Provider<ReviewsRepository>((ref) {
  return ReviewsRepositoryImpl();
});
