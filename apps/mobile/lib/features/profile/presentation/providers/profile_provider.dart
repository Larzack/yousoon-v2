import 'dart:io';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../data/models/profile_model.dart';
import '../../data/repositories/profile_repository.dart';

part 'profile_provider.freezed.dart';

/// État du profil utilisateur
@freezed
class ProfileState with _$ProfileState {
  const factory ProfileState({
    ProfileModel? profile,
    @Default(false) bool isLoading,
    @Default(false) bool isSaving,
    String? error,
    String? successMessage,
  }) = _ProfileState;
}

/// Notifier pour le profil
class ProfileNotifier extends StateNotifier<ProfileState> {
  final ProfileRepository _repository;

  ProfileNotifier(this._repository) : super(const ProfileState());

  /// Charger le profil
  Future<void> loadProfile() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final profile = await _repository.getMyProfile();
      state = state.copyWith(
        profile: profile,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Mettre à jour le profil
  Future<void> updateProfile(UpdateProfileParams params) async {
    state = state.copyWith(isSaving: true, error: null, successMessage: null);

    try {
      final updatedProfile = await _repository.updateProfile(params);
      state = state.copyWith(
        profile: updatedProfile,
        isSaving: false,
        successMessage: 'Profil mis à jour',
      );
    } catch (e) {
      state = state.copyWith(
        isSaving: false,
        error: e.toString(),
      );
    }
  }

  /// Mettre à jour l'avatar
  Future<void> updateAvatar(File imageFile) async {
    state = state.copyWith(isSaving: true, error: null);

    try {
      final avatarUrl = await _repository.updateAvatar(imageFile);
      
      if (state.profile != null) {
        state = state.copyWith(
          profile: state.profile!.copyWith(
            profile: state.profile!.profile.copyWith(avatar: avatarUrl),
          ),
          isSaving: false,
          successMessage: 'Photo de profil mise à jour',
        );
      }
    } catch (e) {
      state = state.copyWith(
        isSaving: false,
        error: e.toString(),
      );
    }
  }

  /// Supprimer l'avatar
  Future<void> deleteAvatar() async {
    state = state.copyWith(isSaving: true, error: null);

    try {
      await _repository.deleteAvatar();
      
      if (state.profile != null) {
        state = state.copyWith(
          profile: state.profile!.copyWith(
            profile: state.profile!.profile.copyWith(avatar: null),
          ),
          isSaving: false,
          successMessage: 'Photo de profil supprimée',
        );
      }
    } catch (e) {
      state = state.copyWith(
        isSaving: false,
        error: e.toString(),
      );
    }
  }

  /// Mettre à jour les préférences
  Future<void> updatePreferences(UpdatePreferencesParams params) async {
    state = state.copyWith(isSaving: true, error: null);

    try {
      final prefs = await _repository.updatePreferences(params);
      
      if (state.profile != null) {
        state = state.copyWith(
          profile: state.profile!.copyWith(preferences: prefs),
          isSaving: false,
          successMessage: 'Préférences mises à jour',
        );
      }
    } catch (e) {
      state = state.copyWith(
        isSaving: false,
        error: e.toString(),
      );
    }
  }

  /// Effacer le message de succès
  void clearSuccessMessage() {
    state = state.copyWith(successMessage: null);
  }

  /// Effacer l'erreur
  void clearError() {
    state = state.copyWith(error: null);
  }
}

/// Provider du profil
final profileProvider =
    StateNotifierProvider<ProfileNotifier, ProfileState>((ref) {
  final repository = ref.watch(profileRepositoryProvider);
  return ProfileNotifier(repository);
});

/// Provider pour le profil chargé (auto-load)
final myProfileProvider = FutureProvider<ProfileModel?>((ref) async {
  final repository = ref.watch(profileRepositoryProvider);
  try {
    return await repository.getMyProfile();
  } catch (e) {
    return null;
  }
});

/// État du changement de mot de passe
@freezed
class ChangePasswordState with _$ChangePasswordState {
  const factory ChangePasswordState.initial() = _CPInitial;
  const factory ChangePasswordState.loading() = _CPLoading;
  const factory ChangePasswordState.success() = _CPSuccess;
  const factory ChangePasswordState.error(String message) = _CPError;
}

/// Notifier pour le changement de mot de passe
class ChangePasswordNotifier extends StateNotifier<ChangePasswordState> {
  final ProfileRepository _repository;

  ChangePasswordNotifier(this._repository)
      : super(const ChangePasswordState.initial());

  Future<void> changePassword({
    required String currentPassword,
    required String newPassword,
  }) async {
    state = const ChangePasswordState.loading();

    try {
      await _repository.changePassword(
        currentPassword: currentPassword,
        newPassword: newPassword,
      );
      state = const ChangePasswordState.success();
    } catch (e) {
      state = ChangePasswordState.error(e.toString());
    }
  }

  void reset() {
    state = const ChangePasswordState.initial();
  }
}

/// Provider pour le changement de mot de passe
final changePasswordProvider =
    StateNotifierProvider<ChangePasswordNotifier, ChangePasswordState>((ref) {
  final repository = ref.watch(profileRepositoryProvider);
  return ChangePasswordNotifier(repository);
});

/// État de la vérification d'identité
@freezed
class IdentityVerificationState with _$IdentityVerificationState {
  const factory IdentityVerificationState.initial() = _IVInitial;
  const factory IdentityVerificationState.uploading(double progress) = _IVUploading;
  const factory IdentityVerificationState.processing() = _IVProcessing;
  const factory IdentityVerificationState.success() = _IVSuccess;
  const factory IdentityVerificationState.error(String message) = _IVError;
}

/// Notifier pour la vérification d'identité
class IdentityVerificationNotifier
    extends StateNotifier<IdentityVerificationState> {
  final ProfileRepository _repository;

  IdentityVerificationNotifier(this._repository)
      : super(const IdentityVerificationState.initial());

  Future<void> submitVerification({
    required String documentType,
    required File frontImage,
    File? backImage,
    File? selfieImage,
  }) async {
    state = const IdentityVerificationState.uploading(0.0);

    try {
      // Simuler progression upload
      state = const IdentityVerificationState.uploading(0.5);
      
      state = const IdentityVerificationState.processing();
      
      await _repository.submitIdentityVerification(
        documentType: documentType,
        frontImage: frontImage,
        backImage: backImage,
        selfieImage: selfieImage,
      );
      
      state = const IdentityVerificationState.success();
    } catch (e) {
      state = IdentityVerificationState.error(e.toString());
    }
  }

  void reset() {
    state = const IdentityVerificationState.initial();
  }
}

/// Provider pour la vérification d'identité
final identityVerificationProvider = StateNotifierProvider<
    IdentityVerificationNotifier, IdentityVerificationState>((ref) {
  final repository = ref.watch(profileRepositoryProvider);
  return IdentityVerificationNotifier(repository);
});

/// État de la suppression de compte
@freezed
class AccountDeletionState with _$AccountDeletionState {
  const factory AccountDeletionState.initial() = _ADInitial;
  const factory AccountDeletionState.loading() = _ADLoading;
  const factory AccountDeletionState.requested(DateTime scheduledAt) = _ADRequested;
  const factory AccountDeletionState.cancelled() = _ADCancelled;
  const factory AccountDeletionState.error(String message) = _ADError;
}

/// Notifier pour la suppression de compte
class AccountDeletionNotifier extends StateNotifier<AccountDeletionState> {
  final ProfileRepository _repository;

  AccountDeletionNotifier(this._repository)
      : super(const AccountDeletionState.initial());

  Future<void> requestDeletion(String? reason) async {
    state = const AccountDeletionState.loading();

    try {
      await _repository.requestAccountDeletion(reason);
      // Suppression programmée dans 30 jours
      final scheduledAt = DateTime.now().add(const Duration(days: 30));
      state = AccountDeletionState.requested(scheduledAt);
    } catch (e) {
      state = AccountDeletionState.error(e.toString());
    }
  }

  Future<void> cancelDeletion() async {
    state = const AccountDeletionState.loading();

    try {
      await _repository.cancelAccountDeletion();
      state = const AccountDeletionState.cancelled();
    } catch (e) {
      state = AccountDeletionState.error(e.toString());
    }
  }

  void reset() {
    state = const AccountDeletionState.initial();
  }
}

/// Provider pour la suppression de compte
final accountDeletionProvider =
    StateNotifierProvider<AccountDeletionNotifier, AccountDeletionState>((ref) {
  final repository = ref.watch(profileRepositoryProvider);
  return AccountDeletionNotifier(repository);
});

/// Provider pour les statistiques du profil
final profileStatsProvider = Provider<ProfileStatsModel?>((ref) {
  final profileState = ref.watch(profileProvider);
  return profileState.profile?.stats;
});

/// Provider pour le grade actuel
final currentGradeProvider = Provider<UserGrade>((ref) {
  final profileState = ref.watch(profileProvider);
  final grade = profileState.profile?.grade ?? 'explorateur';
  return UserGrade.fromString(grade);
});

/// Provider pour la progression vers le prochain grade
final gradeProgressProvider = Provider<(UserGrade? nextGrade, int outingsRemaining, double progress)>((ref) {
  final stats = ref.watch(profileStatsProvider);
  final currentGrade = ref.watch(currentGradeProvider);
  
  final currentOutings = stats?.completedOutings ?? 0;
  final (nextGrade, remaining) = currentGrade.progressTo(currentOutings);
  
  if (nextGrade == null) {
    return (null, 0, 1.0); // Maximum atteint
  }
  
  final totalRequired = nextGrade.requiredOutings - currentGrade.requiredOutings;
  final done = currentOutings - currentGrade.requiredOutings;
  final progress = totalRequired > 0 ? done / totalRequired : 0.0;
  
  return (nextGrade, remaining, progress.clamp(0.0, 1.0));
});
