import 'package:freezed_annotation/freezed_annotation.dart';

part 'user_model.freezed.dart';
part 'user_model.g.dart';

/// Modèle utilisateur
@freezed
class UserModel with _$UserModel {
  const factory UserModel({
    required String id,
    required String email,
    String? phone,
    required ProfileModel profile,
    required IdentityStatusModel identity,
    SubscriptionModel? subscription,
    required PreferencesModel preferences,
    required String status,
    required String grade,
    required DateTime createdAt,
    DateTime? lastLoginAt,
  }) = _UserModel;

  factory UserModel.fromJson(Map<String, dynamic> json) =>
      _$UserModelFromJson(json);
}

/// Modèle profil utilisateur
@freezed
class ProfileModel with _$ProfileModel {
  const factory ProfileModel({
    required String firstName,
    required String lastName,
    String? displayName,
    String? avatar,
    DateTime? birthDate,
    String? gender,
  }) = _ProfileModel;

  factory ProfileModel.fromJson(Map<String, dynamic> json) =>
      _$ProfileModelFromJson(json);
}

/// Statut de vérification d'identité
@freezed
class IdentityStatusModel with _$IdentityStatusModel {
  const factory IdentityStatusModel({
    required String status, // not_submitted, pending, verified, rejected
    DateTime? verifiedAt,
    String? documentType,
    int? attemptsRemaining,
  }) = _IdentityStatusModel;

  factory IdentityStatusModel.fromJson(Map<String, dynamic> json) =>
      _$IdentityStatusModelFromJson(json);
}

/// Modèle abonnement
@freezed
class SubscriptionModel with _$SubscriptionModel {
  const factory SubscriptionModel({
    required String id,
    required String planId,
    required String planName,
    required String status, // trialing, active, past_due, cancelled, expired
    required String platform, // ios, android
    DateTime? trialEndDate,
    required DateTime currentPeriodStart,
    required DateTime currentPeriodEnd,
    required bool autoRenew,
  }) = _SubscriptionModel;

  factory SubscriptionModel.fromJson(Map<String, dynamic> json) =>
      _$SubscriptionModelFromJson(json);
}

/// Préférences utilisateur
@freezed
class PreferencesModel with _$PreferencesModel {
  const factory PreferencesModel({
    @Default('fr') String language,
    required NotificationPreferencesModel notifications,
    @Default([]) List<String> favoriteCategories,
    @Default(10) int maxDistance,
  }) = _PreferencesModel;

  factory PreferencesModel.fromJson(Map<String, dynamic> json) =>
      _$PreferencesModelFromJson(json);
}

/// Préférences de notifications
@freezed
class NotificationPreferencesModel with _$NotificationPreferencesModel {
  const factory NotificationPreferencesModel({
    @Default(true) bool push,
    @Default(true) bool email,
    @Default(false) bool sms,
    @Default(true) bool marketing,
  }) = _NotificationPreferencesModel;

  factory NotificationPreferencesModel.fromJson(Map<String, dynamic> json) =>
      _$NotificationPreferencesModelFromJson(json);
}

/// Token d'authentification
@freezed
class AuthTokensModel with _$AuthTokensModel {
  const factory AuthTokensModel({
    required String accessToken,
    required String refreshToken,
    required DateTime accessTokenExpiry,
    required DateTime refreshTokenExpiry,
  }) = _AuthTokensModel;

  factory AuthTokensModel.fromJson(Map<String, dynamic> json) =>
      _$AuthTokensModelFromJson(json);
}
