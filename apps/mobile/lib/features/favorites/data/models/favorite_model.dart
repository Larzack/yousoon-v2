import 'package:freezed_annotation/freezed_annotation.dart';

part 'favorite_model.freezed.dart';
part 'favorite_model.g.dart';

/// Modèle d'un favori
@freezed
class FavoriteModel with _$FavoriteModel {
  const factory FavoriteModel({
    required String id,
    required String userId,
    required String offerId,
    required FavoriteOfferModel offer,
    required DateTime addedAt,
  }) = _FavoriteModel;

  factory FavoriteModel.fromJson(Map<String, dynamic> json) =>
      _$FavoriteModelFromJson(json);
}

/// Offre dans un favori (version légère)
@freezed
class FavoriteOfferModel with _$FavoriteOfferModel {
  const factory FavoriteOfferModel({
    required String id,
    required String title,
    String? shortDescription,
    required String discountType,
    required int discountValue,
    String? discountFormula,
    required List<String> images,
    required String categoryName,
    required FavoritePartnerModel partner,
    required FavoriteLocationModel location,
    required String status, // active, paused, expired
    double? avgRating,
    int? reviewCount,
  }) = _FavoriteOfferModel;

  factory FavoriteOfferModel.fromJson(Map<String, dynamic> json) =>
      _$FavoriteOfferModelFromJson(json);
}

/// Partenaire dans un favori
@freezed
class FavoritePartnerModel with _$FavoritePartnerModel {
  const factory FavoritePartnerModel({
    required String id,
    required String name,
    String? logo,
  }) = _FavoritePartnerModel;

  factory FavoritePartnerModel.fromJson(Map<String, dynamic> json) =>
      _$FavoritePartnerModelFromJson(json);
}

/// Localisation dans un favori
@freezed
class FavoriteLocationModel with _$FavoriteLocationModel {
  const factory FavoriteLocationModel({
    required double latitude,
    required double longitude,
    required String address,
    required String city,
  }) = _FavoriteLocationModel;

  factory FavoriteLocationModel.fromJson(Map<String, dynamic> json) =>
      _$FavoriteLocationModelFromJson(json);
}

/// Liste paginée de favoris
@freezed
class PaginatedFavoritesResult with _$PaginatedFavoritesResult {
  const factory PaginatedFavoritesResult({
    required List<FavoriteModel> favorites,
    required int totalCount,
    required int page,
    required int totalPages,
    required bool hasMore,
  }) = _PaginatedFavoritesResult;

  factory PaginatedFavoritesResult.fromJson(Map<String, dynamic> json) =>
      _$PaginatedFavoritesResultFromJson(json);
}
