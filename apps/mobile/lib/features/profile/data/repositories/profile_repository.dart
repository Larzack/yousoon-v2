import 'dart:io';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../models/profile_model.dart';

/// Interface du repository du profil
abstract class ProfileRepository {
  /// Obtenir le profil de l'utilisateur courant
  Future<ProfileModel> getMyProfile();

  /// Mettre à jour le profil
  Future<ProfileModel> updateProfile(UpdateProfileParams params);

  /// Mettre à jour l'avatar
  Future<String> updateAvatar(File imageFile);

  /// Supprimer l'avatar
  Future<void> deleteAvatar();

  /// Mettre à jour les préférences
  Future<ProfilePreferencesModel> updatePreferences(UpdatePreferencesParams params);

  /// Changer le mot de passe
  Future<void> changePassword({
    required String currentPassword,
    required String newPassword,
  });

  /// Demander la suppression du compte (RGPD - 30 jours)
  Future<void> requestAccountDeletion(String? reason);

  /// Annuler la demande de suppression
  Future<void> cancelAccountDeletion();

  /// Soumettre une vérification d'identité
  Future<void> submitIdentityVerification({
    required String documentType,
    required File frontImage,
    File? backImage,
    File? selfieImage,
  });

  /// Gérer l'abonnement
  Future<void> cancelSubscription(String? reason);

  /// Restaurer l'abonnement (après annulation mais avant fin de période)
  Future<void> restoreSubscription();

  /// Exporter les données (RGPD)
  Future<String> exportMyData();
}

/// Implémentation du repository du profil
class ProfileRepositoryImpl implements ProfileRepository {
  // final Client _graphqlClient;

  ProfileRepositoryImpl({
    // required Client graphqlClient,
  });
  // : _graphqlClient = graphqlClient;

  @override
  Future<ProfileModel> getMyProfile() async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get my profile not implemented');
  }

  @override
  Future<ProfileModel> updateProfile(UpdateProfileParams params) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Update profile not implemented');
  }

  @override
  Future<String> updateAvatar(File imageFile) async {
    // TODO: Implémenter avec upload S3 puis GraphQL
    throw UnimplementedError('Update avatar not implemented');
  }

  @override
  Future<void> deleteAvatar() async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Delete avatar not implemented');
  }

  @override
  Future<ProfilePreferencesModel> updatePreferences(
      UpdatePreferencesParams params) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Update preferences not implemented');
  }

  @override
  Future<void> changePassword({
    required String currentPassword,
    required String newPassword,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Change password not implemented');
  }

  @override
  Future<void> requestAccountDeletion(String? reason) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Request account deletion not implemented');
  }

  @override
  Future<void> cancelAccountDeletion() async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Cancel account deletion not implemented');
  }

  @override
  Future<void> submitIdentityVerification({
    required String documentType,
    required File frontImage,
    File? backImage,
    File? selfieImage,
  }) async {
    // TODO: Implémenter avec upload S3 puis GraphQL
    throw UnimplementedError('Submit identity verification not implemented');
  }

  @override
  Future<void> cancelSubscription(String? reason) async {
    // TODO: Implémenter avec in-app purchase + GraphQL
    throw UnimplementedError('Cancel subscription not implemented');
  }

  @override
  Future<void> restoreSubscription() async {
    // TODO: Implémenter avec in-app purchase + GraphQL
    throw UnimplementedError('Restore subscription not implemented');
  }

  @override
  Future<String> exportMyData() async {
    // TODO: Implémenter avec GraphQL - retourne URL du fichier
    throw UnimplementedError('Export my data not implemented');
  }
}

/// Provider du repository du profil
final profileRepositoryProvider = Provider<ProfileRepository>((ref) {
  return ProfileRepositoryImpl();
});
