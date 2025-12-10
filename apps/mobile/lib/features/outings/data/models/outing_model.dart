import 'package:freezed_annotation/freezed_annotation.dart';

part 'outing_model.freezed.dart';
part 'outing_model.g.dart';

/// Modèle d'une sortie (réservation)
@freezed
class OutingModel with _$OutingModel {
  const factory OutingModel({
    required String id,
    required String userId,
    required String offerId,
    required String partnerId,
    required String establishmentId,
    required QRCodeModel qrCode,
    required OutingOfferSnapshotModel offerSnapshot,
    required OutingPartnerSnapshotModel partnerSnapshot,
    required OutingEstablishmentSnapshotModel establishmentSnapshot,
    required String status, // pending, confirmed, checked_in, cancelled, expired, no_show
    required DateTime bookedAt,
    required DateTime expiresAt,
    DateTime? checkedInAt,
    DateTime? cancelledAt,
    String? cancellationReason,
    required DateTime createdAt,
    DateTime? updatedAt,
  }) = _OutingModel;

  factory OutingModel.fromJson(Map<String, dynamic> json) =>
      _$OutingModelFromJson(json);
}

/// QR Code pour le check-in
@freezed
class QRCodeModel with _$QRCodeModel {
  const factory QRCodeModel({
    required String code,
    required String data, // Données encodées dans le QR
    required DateTime expiresAt,
  }) = _QRCodeModel;

  factory QRCodeModel.fromJson(Map<String, dynamic> json) =>
      _$QRCodeModelFromJson(json);
}

/// Snapshot de l'offre au moment de la réservation
@freezed
class OutingOfferSnapshotModel with _$OutingOfferSnapshotModel {
  const factory OutingOfferSnapshotModel({
    required String id,
    required String title,
    String? shortDescription,
    required String discountType,
    required int discountValue,
    String? discountFormula,
    required List<String> images,
    required String category,
  }) = _OutingOfferSnapshotModel;

  factory OutingOfferSnapshotModel.fromJson(Map<String, dynamic> json) =>
      _$OutingOfferSnapshotModelFromJson(json);
}

/// Snapshot du partenaire au moment de la réservation
@freezed
class OutingPartnerSnapshotModel with _$OutingPartnerSnapshotModel {
  const factory OutingPartnerSnapshotModel({
    required String id,
    required String name,
    String? logo,
  }) = _OutingPartnerSnapshotModel;

  factory OutingPartnerSnapshotModel.fromJson(Map<String, dynamic> json) =>
      _$OutingPartnerSnapshotModelFromJson(json);
}

/// Snapshot de l'établissement au moment de la réservation
@freezed
class OutingEstablishmentSnapshotModel with _$OutingEstablishmentSnapshotModel {
  const factory OutingEstablishmentSnapshotModel({
    required String id,
    required String name,
    required String address,
    required String city,
    required double latitude,
    required double longitude,
    String? phone,
  }) = _OutingEstablishmentSnapshotModel;

  factory OutingEstablishmentSnapshotModel.fromJson(Map<String, dynamic> json) =>
      _$OutingEstablishmentSnapshotModelFromJson(json);
}

/// Statut d'une sortie
enum OutingStatus {
  pending,
  confirmed,
  checkedIn,
  cancelled,
  expired,
  noShow,
}

/// Extension pour le statut
extension OutingStatusExtension on OutingStatus {
  String get value {
    switch (this) {
      case OutingStatus.pending:
        return 'pending';
      case OutingStatus.confirmed:
        return 'confirmed';
      case OutingStatus.checkedIn:
        return 'checked_in';
      case OutingStatus.cancelled:
        return 'cancelled';
      case OutingStatus.expired:
        return 'expired';
      case OutingStatus.noShow:
        return 'no_show';
    }
  }

  static OutingStatus fromString(String value) {
    switch (value) {
      case 'pending':
        return OutingStatus.pending;
      case 'confirmed':
        return OutingStatus.confirmed;
      case 'checked_in':
        return OutingStatus.checkedIn;
      case 'cancelled':
        return OutingStatus.cancelled;
      case 'expired':
        return OutingStatus.expired;
      case 'no_show':
        return OutingStatus.noShow;
      default:
        return OutingStatus.pending;
    }
  }
}

/// Liste de sorties paginée
@freezed
class PaginatedOutingsResult with _$PaginatedOutingsResult {
  const factory PaginatedOutingsResult({
    required List<OutingModel> outings,
    required int totalCount,
    required int page,
    required int totalPages,
    required bool hasMore,
  }) = _PaginatedOutingsResult;

  factory PaginatedOutingsResult.fromJson(Map<String, dynamic> json) =>
      _$PaginatedOutingsResultFromJson(json);
}

/// Paramètres de filtre pour les sorties
@freezed
class OutingsFilterParams with _$OutingsFilterParams {
  const factory OutingsFilterParams({
    OutingFilter? filter, // upcoming, past, all
    @Default(1) int page,
    @Default(20) int limit,
  }) = _OutingsFilterParams;

  factory OutingsFilterParams.fromJson(Map<String, dynamic> json) =>
      _$OutingsFilterParamsFromJson(json);
}

/// Filtre de sorties
enum OutingFilter {
  upcoming, // À venir (confirmed)
  past, // Passées (checked_in, expired, no_show)
  cancelled, // Annulées
  all,
}
