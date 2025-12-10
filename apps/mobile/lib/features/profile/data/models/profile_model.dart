import 'package:freezed_annotation/freezed_annotation.dart';

part 'profile_model.freezed.dart';
part 'profile_model.g.dart';

/// Modèle du profil utilisateur complet
@freezed
class ProfileModel with _$ProfileModel {
  const factory ProfileModel({
    required String id,
    required String email,
    String? phone,
    required ProfileInfoModel profile,
    required ProfileIdentityModel identity,
    ProfileSubscriptionModel? subscription,
    required ProfilePreferencesModel preferences,
    required ProfileStatsModel stats,
    required String status, // active, suspended, deleted
    required String grade, // explorateur, aventurier, grand_voyageur, conquerant
    @Default(false) bool emailVerified,
    @Default(false) bool phoneVerified,
    required DateTime createdAt,
    DateTime? lastLoginAt,
  }) = _ProfileModel;

  factory ProfileModel.fromJson(Map<String, dynamic> json) =>
      _$ProfileModelFromJson(json);
}

/// Informations du profil
@freezed
class ProfileInfoModel with _$ProfileInfoModel {
  const factory ProfileInfoModel({
    required String firstName,
    required String lastName,
    String? displayName,
    String? avatar,
    DateTime? birthDate,
    String? gender, // male, female, other
    String? bio,
  }) = _ProfileInfoModel;

  factory ProfileInfoModel.fromJson(Map<String, dynamic> json) =>
      _$ProfileInfoModelFromJson(json);
}

/// Statut de vérification d'identité
@freezed
class ProfileIdentityModel with _$ProfileIdentityModel {
  const factory ProfileIdentityModel({
    required String status, // not_submitted, pending, verified, rejected
    DateTime? submittedAt,
    DateTime? verifiedAt,
    String? documentType,
    String? rejectionReason,
    @Default(0) int attemptCount,
    @Default(10) int maxAttempts,
  }) = _ProfileIdentityModel;

  factory ProfileIdentityModel.fromJson(Map<String, dynamic> json) =>
      _$ProfileIdentityModelFromJson(json);
}

/// Abonnement de l'utilisateur
@freezed
class ProfileSubscriptionModel with _$ProfileSubscriptionModel {
  const factory ProfileSubscriptionModel({
    required String id,
    required String planId,
    required String planCode, // free, monthly, yearly
    required String planName,
    required String status, // trialing, active, past_due, cancelled, expired
    required String platform, // apple, google
    DateTime? trialStartDate,
    DateTime? trialEndDate,
    required DateTime currentPeriodStart,
    required DateTime currentPeriodEnd,
    @Default(true) bool autoRenew,
    DateTime? cancelledAt,
  }) = _ProfileSubscriptionModel;

  factory ProfileSubscriptionModel.fromJson(Map<String, dynamic> json) =>
      _$ProfileSubscriptionModelFromJson(json);
}

/// Préférences utilisateur
@freezed
class ProfilePreferencesModel with _$ProfilePreferencesModel {
  const factory ProfilePreferencesModel({
    @Default('fr') String language,
    @Default(true) bool pushNotifications,
    @Default(true) bool emailNotifications,
    @Default(false) bool smsNotifications,
    @Default(true) bool marketingNotifications,
    @Default([]) List<String> favoriteCategories,
    @Default(10) int maxDistance, // km
    @Default(true) bool biometricEnabled,
  }) = _ProfilePreferencesModel;

  factory ProfilePreferencesModel.fromJson(Map<String, dynamic> json) =>
      _$ProfilePreferencesModelFromJson(json);
}

/// Statistiques utilisateur
@freezed
class ProfileStatsModel with _$ProfileStatsModel {
  const factory ProfileStatsModel({
    @Default(0) int totalOutings,
    @Default(0) int completedOutings,
    @Default(0) int cancelledOutings,
    @Default(0) int totalFavorites,
    @Default(0) int totalReviews,
    @Default(0) int pointsEarned,
    DateTime? lastOutingAt,
  }) = _ProfileStatsModel;

  factory ProfileStatsModel.fromJson(Map<String, dynamic> json) =>
      _$ProfileStatsModelFromJson(json);
}

/// Paramètres de mise à jour du profil
@freezed
class UpdateProfileParams with _$UpdateProfileParams {
  const factory UpdateProfileParams({
    String? firstName,
    String? lastName,
    String? displayName,
    DateTime? birthDate,
    String? gender,
    String? bio,
  }) = _UpdateProfileParams;

  factory UpdateProfileParams.fromJson(Map<String, dynamic> json) =>
      _$UpdateProfileParamsFromJson(json);
}

/// Paramètres de mise à jour des préférences
@freezed
class UpdatePreferencesParams with _$UpdatePreferencesParams {
  const factory UpdatePreferencesParams({
    String? language,
    bool? pushNotifications,
    bool? emailNotifications,
    bool? smsNotifications,
    bool? marketingNotifications,
    List<String>? favoriteCategories,
    int? maxDistance,
    bool? biometricEnabled,
  }) = _UpdatePreferencesParams;

  factory UpdatePreferencesParams.fromJson(Map<String, dynamic> json) =>
      _$UpdatePreferencesParamsFromJson(json);
}

/// Grades utilisateur avec leurs seuils
enum UserGrade {
  explorateur(0, '🧭', 'Explorateur'),
  aventurier(10, '🎒', 'Aventurier'),
  grandVoyageur(50, '✈️', 'Grand Voyageur'),
  conquerant(100, '👑', 'Conquérant');

  final int requiredOutings;
  final String emoji;
  final String displayName;

  const UserGrade(this.requiredOutings, this.emoji, this.displayName);

  static UserGrade fromOutings(int outings) {
    if (outings >= 100) return UserGrade.conquerant;
    if (outings >= 50) return UserGrade.grandVoyageur;
    if (outings >= 10) return UserGrade.aventurier;
    return UserGrade.explorateur;
  }

  static UserGrade fromString(String value) {
    switch (value.toLowerCase()) {
      case 'conquerant':
      case 'conquérant':
        return UserGrade.conquerant;
      case 'grand_voyageur':
      case 'grandvoyageur':
        return UserGrade.grandVoyageur;
      case 'aventurier':
        return UserGrade.aventurier;
      default:
        return UserGrade.explorateur;
    }
  }

  /// Prochain grade et nombre de sorties restantes
  (UserGrade? nextGrade, int outingsRemaining) progressTo(int currentOutings) {
    final values = UserGrade.values;
    final currentIndex = values.indexOf(this);
    
    if (currentIndex >= values.length - 1) {
      return (null, 0); // Déjà au max
    }
    
    final next = values[currentIndex + 1];
    return (next, next.requiredOutings - currentOutings);
  }
}
