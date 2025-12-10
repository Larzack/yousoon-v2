import 'package:freezed_annotation/freezed_annotation.dart';

part 'offer_model.freezed.dart';
part 'offer_model.g.dart';

/// Modèle d'une offre
@freezed
class OfferModel with _$OfferModel {
  const factory OfferModel({
    required String id,
    required String partnerId,
    required String establishmentId,
    required String title,
    required String description,
    String? shortDescription,
    required String categoryId,
    required String categoryName,
    required DiscountModel discount,
    required OfferValidityModel validity,
    required OfferScheduleModel schedule,
    OfferQuotaModel? quota,
    required List<String> images,
    required OfferLocationModel location,
    required OfferPartnerInfoModel partner,
    required OfferStatsModel stats,
    required String status,
    required DateTime createdAt,
    DateTime? publishedAt,
  }) = _OfferModel;

  factory OfferModel.fromJson(Map<String, dynamic> json) =>
      _$OfferModelFromJson(json);
}

/// Modèle de réduction
@freezed
class DiscountModel with _$DiscountModel {
  const factory DiscountModel({
    required String type, // percentage, fixed, formula
    required int value, // Pourcentage ou centimes
    int? originalPrice,
    String? formula, // ex: "1 acheté = 1 offert"
    String? displayText, // Texte à afficher
  }) = _DiscountModel;

  factory DiscountModel.fromJson(Map<String, dynamic> json) =>
      _$DiscountModelFromJson(json);
}

/// Période de validité de l'offre
@freezed
class OfferValidityModel with _$OfferValidityModel {
  const factory OfferValidityModel({
    required DateTime startDate,
    required DateTime endDate,
    @Default('Europe/Paris') String timezone,
  }) = _OfferValidityModel;

  factory OfferValidityModel.fromJson(Map<String, dynamic> json) =>
      _$OfferValidityModelFromJson(json);
}

/// Horaires de l'offre
@freezed
class OfferScheduleModel with _$OfferScheduleModel {
  const factory OfferScheduleModel({
    @Default(false) bool allDay,
    @Default([]) List<TimeSlotModel> slots,
  }) = _OfferScheduleModel;

  factory OfferScheduleModel.fromJson(Map<String, dynamic> json) =>
      _$OfferScheduleModelFromJson(json);
}

/// Créneau horaire
@freezed
class TimeSlotModel with _$TimeSlotModel {
  const factory TimeSlotModel({
    required int dayOfWeek, // 0 = Dimanche, 1 = Lundi, ...
    required String startTime, // "17:00"
    required String endTime, // "20:00"
  }) = _TimeSlotModel;

  factory TimeSlotModel.fromJson(Map<String, dynamic> json) =>
      _$TimeSlotModelFromJson(json);
}

/// Quotas de l'offre
@freezed
class OfferQuotaModel with _$OfferQuotaModel {
  const factory OfferQuotaModel({
    int? total,
    int? perUser,
    int? perDay,
    @Default(0) int used,
  }) = _OfferQuotaModel;

  factory OfferQuotaModel.fromJson(Map<String, dynamic> json) =>
      _$OfferQuotaModelFromJson(json);
}

/// Localisation de l'offre (dénormalisée depuis l'établissement)
@freezed
class OfferLocationModel with _$OfferLocationModel {
  const factory OfferLocationModel({
    required double latitude,
    required double longitude,
    required String address,
    required String city,
    String? postalCode,
    double? distance, // Distance depuis la position actuelle (calculée)
  }) = _OfferLocationModel;

  factory OfferLocationModel.fromJson(Map<String, dynamic> json) =>
      _$OfferLocationModelFromJson(json);
}

/// Infos du partenaire (dénormalisées)
@freezed
class OfferPartnerInfoModel with _$OfferPartnerInfoModel {
  const factory OfferPartnerInfoModel({
    required String id,
    required String name,
    String? logo,
    String? category,
  }) = _OfferPartnerInfoModel;

  factory OfferPartnerInfoModel.fromJson(Map<String, dynamic> json) =>
      _$OfferPartnerInfoModelFromJson(json);
}

/// Statistiques de l'offre
@freezed
class OfferStatsModel with _$OfferStatsModel {
  const factory OfferStatsModel({
    @Default(0) int views,
    @Default(0) int bookings,
    @Default(0) int checkins,
    @Default(0) int favorites,
    double? avgRating,
    @Default(0) int reviewCount,
  }) = _OfferStatsModel;

  factory OfferStatsModel.fromJson(Map<String, dynamic> json) =>
      _$OfferStatsModelFromJson(json);
}

/// Catégorie
@freezed
class CategoryModel with _$CategoryModel {
  const factory CategoryModel({
    required String id,
    required String name,
    required String slug,
    String? description,
    required String icon,
    required String color,
    String? parentId,
    @Default(true) bool isActive,
  }) = _CategoryModel;

  factory CategoryModel.fromJson(Map<String, dynamic> json) =>
      _$CategoryModelFromJson(json);
}

/// Paramètres de recherche d'offres
@freezed
class OfferSearchParams with _$OfferSearchParams {
  const factory OfferSearchParams({
    String? query,
    double? latitude,
    double? longitude,
    @Default(10) double radius, // km
    List<String>? categoryIds,
    String? discountType,
    int? minDiscount,
    @Default(OfferSortBy.distance) OfferSortBy sortBy,
    @Default(1) int page,
    @Default(20) int limit,
  }) = _OfferSearchParams;

  factory OfferSearchParams.fromJson(Map<String, dynamic> json) =>
      _$OfferSearchParamsFromJson(json);
}

/// Tri des offres
enum OfferSortBy {
  distance,
  discount,
  rating,
  newest,
  popularity,
}

/// Résultat de recherche paginé
@freezed
class PaginatedOffersResult with _$PaginatedOffersResult {
  const factory PaginatedOffersResult({
    required List<OfferModel> offers,
    required int totalCount,
    required int page,
    required int totalPages,
    required bool hasMore,
  }) = _PaginatedOffersResult;

  factory PaginatedOffersResult.fromJson(Map<String, dynamic> json) =>
      _$PaginatedOffersResultFromJson(json);
}
