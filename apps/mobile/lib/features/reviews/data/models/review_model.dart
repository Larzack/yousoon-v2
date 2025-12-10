import 'package:freezed_annotation/freezed_annotation.dart';

part 'review_model.freezed.dart';
part 'review_model.g.dart';

/// Modèle d'un avis
@freezed
class ReviewModel with _$ReviewModel {
  const factory ReviewModel({
    required String id,
    required String userId,
    required String offerId,
    required String partnerId,
    String? establishmentId,
    String? outingId,
    required int rating, // 1-5
    String? title,
    String? content,
    @Default([]) List<String> images,
    required ReviewUserModel user,
    required ReviewOfferModel offer,
    required ReviewModerationModel moderation,
    @Default(0) int helpfulCount,
    @Default(false) bool isVerifiedPurchase,
    required DateTime createdAt,
    DateTime? updatedAt,
  }) = _ReviewModel;

  factory ReviewModel.fromJson(Map<String, dynamic> json) =>
      _$ReviewModelFromJson(json);
}

/// Utilisateur de l'avis
@freezed
class ReviewUserModel with _$ReviewUserModel {
  const factory ReviewUserModel({
    required String id,
    required String firstName,
    String? avatar,
    String? grade, // explorateur, aventurier, etc.
  }) = _ReviewUserModel;

  factory ReviewUserModel.fromJson(Map<String, dynamic> json) =>
      _$ReviewUserModelFromJson(json);
}

/// Offre liée à l'avis
@freezed
class ReviewOfferModel with _$ReviewOfferModel {
  const factory ReviewOfferModel({
    required String id,
    required String title,
    String? partnerName,
  }) = _ReviewOfferModel;

  factory ReviewOfferModel.fromJson(Map<String, dynamic> json) =>
      _$ReviewOfferModelFromJson(json);
}

/// Modération de l'avis
@freezed
class ReviewModerationModel with _$ReviewModerationModel {
  const factory ReviewModerationModel({
    required String status, // pending, approved, rejected, reported
    @Default(0) int reportCount,
  }) = _ReviewModerationModel;

  factory ReviewModerationModel.fromJson(Map<String, dynamic> json) =>
      _$ReviewModerationModelFromJson(json);
}

/// Liste paginée d'avis
@freezed
class PaginatedReviewsResult with _$PaginatedReviewsResult {
  const factory PaginatedReviewsResult({
    required List<ReviewModel> reviews,
    required int totalCount,
    required int page,
    required int totalPages,
    required bool hasMore,
    double? averageRating,
    @Default({}) Map<int, int> ratingDistribution, // {5: 120, 4: 80, ...}
  }) = _PaginatedReviewsResult;

  factory PaginatedReviewsResult.fromJson(Map<String, dynamic> json) =>
      _$PaginatedReviewsResultFromJson(json);
}

/// Paramètres de création d'un avis
@freezed
class CreateReviewParams with _$CreateReviewParams {
  const factory CreateReviewParams({
    required String offerId,
    String? outingId,
    required int rating,
    String? title,
    String? content,
    @Default([]) List<String> imagePaths,
  }) = _CreateReviewParams;

  factory CreateReviewParams.fromJson(Map<String, dynamic> json) =>
      _$CreateReviewParamsFromJson(json);
}

/// Filtre de tri des avis
enum ReviewSortBy {
  newest,
  oldest,
  highestRating,
  lowestRating,
  mostHelpful,
}
