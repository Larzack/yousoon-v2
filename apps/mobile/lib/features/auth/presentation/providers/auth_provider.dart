import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../data/models/user_model.dart';
import '../../data/repositories/auth_repository.dart';

part 'auth_provider.freezed.dart';

/// État de l'authentification
@freezed
class AuthState with _$AuthState {
  const factory AuthState.initial() = _Initial;
  const factory AuthState.loading() = _Loading;
  const factory AuthState.authenticated(UserModel user) = _Authenticated;
  const factory AuthState.unauthenticated() = _Unauthenticated;
  const factory AuthState.error(String message) = _Error;
}

/// Notifier pour gérer l'état d'authentification
class AuthNotifier extends StateNotifier<AuthState> {
  final AuthRepository _repository;

  AuthNotifier(this._repository) : super(const AuthState.initial()) {
    _checkAuthStatus();
  }

  /// Vérifier le statut d'authentification au démarrage
  Future<void> _checkAuthStatus() async {
    try {
      final isLoggedIn = await _repository.isLoggedIn();
      if (isLoggedIn) {
        final user = await _repository.getCurrentUser();
        if (user != null) {
          state = AuthState.authenticated(user);
        } else {
          state = const AuthState.unauthenticated();
        }
      } else {
        state = const AuthState.unauthenticated();
      }
    } catch (e) {
      state = const AuthState.unauthenticated();
    }
  }

  /// Connexion avec email/mot de passe
  Future<void> login({
    required String email,
    required String password,
  }) async {
    state = const AuthState.loading();
    try {
      final result = await _repository.login(
        email: email,
        password: password,
      );
      state = AuthState.authenticated(result.user);
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  /// Inscription
  Future<void> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    String? phone,
  }) async {
    state = const AuthState.loading();
    try {
      final result = await _repository.register(
        email: email,
        password: password,
        firstName: firstName,
        lastName: lastName,
        phone: phone,
      );
      state = AuthState.authenticated(result.user);
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  /// Connexion sociale
  Future<void> socialLogin({
    required String provider,
    required String token,
  }) async {
    state = const AuthState.loading();
    try {
      final result = await _repository.socialLogin(
        provider: provider,
        token: token,
      );
      state = AuthState.authenticated(result.user);
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  /// Déconnexion
  Future<void> logout() async {
    await _repository.logout();
    state = const AuthState.unauthenticated();
  }

  /// Mettre à jour l'utilisateur dans le state
  void updateUser(UserModel user) {
    state = AuthState.authenticated(user);
  }

  /// Rafraîchir les informations utilisateur
  Future<void> refreshUser() async {
    final currentState = state;
    if (currentState is _Authenticated) {
      try {
        final user = await _repository.getCurrentUser();
        if (user != null) {
          state = AuthState.authenticated(user);
        }
      } catch (e) {
        // Ignorer l'erreur, garder l'état actuel
      }
    }
  }
}

/// Provider de l'état d'authentification
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final repository = ref.watch(authRepositoryProvider);
  return AuthNotifier(repository);
});

/// Provider pour savoir si l'utilisateur est authentifié
final isAuthenticatedProvider = Provider<bool>((ref) {
  final authState = ref.watch(authProvider);
  return authState is _Authenticated;
});

/// Provider pour l'utilisateur courant (peut être null)
final currentUserProvider = Provider<UserModel?>((ref) {
  final authState = ref.watch(authProvider);
  return authState.maybeWhen(
    authenticated: (user) => user,
    orElse: () => null,
  );
});

/// Provider pour vérifier si l'identité est vérifiée
final isIdentityVerifiedProvider = Provider<bool>((ref) {
  final user = ref.watch(currentUserProvider);
  return user?.identity.status == 'verified';
});

/// Provider pour vérifier si l'utilisateur a un abonnement actif
final hasActiveSubscriptionProvider = Provider<bool>((ref) {
  final user = ref.watch(currentUserProvider);
  final subscription = user?.subscription;
  if (subscription == null) return false;
  return subscription.status == 'active' || subscription.status == 'trialing';
});
