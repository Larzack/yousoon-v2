import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../models/user_model.dart';

/// Interface du repository d'authentification
abstract class AuthRepository {
  /// Connexion avec email et mot de passe
  Future<AuthResult> login({
    required String email,
    required String password,
  });

  /// Inscription
  Future<AuthResult> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    String? phone,
  });

  /// Connexion sociale (Google, Apple, Facebook)
  Future<AuthResult> socialLogin({
    required String provider,
    required String token,
  });

  /// Déconnexion
  Future<void> logout();

  /// Rafraîchir le token
  Future<AuthTokensModel?> refreshToken();

  /// Récupérer l'utilisateur courant
  Future<UserModel?> getCurrentUser();

  /// Vérifier si l'utilisateur est connecté
  Future<bool> isLoggedIn();

  /// Mot de passe oublié
  Future<void> forgotPassword(String email);

  /// Réinitialiser le mot de passe
  Future<void> resetPassword({
    required String token,
    required String newPassword,
  });

  /// Vérifier OTP email
  Future<bool> verifyEmailOTP(String code);

  /// Renvoyer OTP email
  Future<void> resendEmailOTP();

  /// Soumettre vérification identité (CNI)
  Future<void> submitIdentityVerification({
    required String documentType,
    required String frontImagePath,
    String? backImagePath,
  });
}

/// Résultat d'authentification
class AuthResult {
  final UserModel user;
  final AuthTokensModel tokens;

  const AuthResult({
    required this.user,
    required this.tokens,
  });
}

/// Implémentation du repository d'authentification
class AuthRepositoryImpl implements AuthRepository {
  // final Client _graphqlClient;
  final FlutterSecureStorage _secureStorage;

  static const _accessTokenKey = 'access_token';
  static const _refreshTokenKey = 'refresh_token';

  AuthRepositoryImpl({
    // required Client graphqlClient,
    FlutterSecureStorage? secureStorage,
  }) : // _graphqlClient = graphqlClient,
       _secureStorage = secureStorage ?? const FlutterSecureStorage();

  @override
  Future<AuthResult> login({
    required String email,
    required String password,
  }) async {
    // TODO: Implémenter avec GraphQL
    // final result = await _graphqlClient.execute(
    //   LoginMutation(
    //     variables: LoginArguments(
    //       email: email,
    //       password: password,
    //     ),
    //   ),
    // );
    
    // if (result.hasErrors) {
    //   throw AuthException(result.errors!.first.message);
    // }
    
    // final data = result.data!.login;
    // final tokens = AuthTokensModel(
    //   accessToken: data.accessToken,
    //   refreshToken: data.refreshToken,
    //   accessTokenExpiry: DateTime.now().add(const Duration(hours: 6)),
    //   refreshTokenExpiry: DateTime.now().add(const Duration(days: 30)),
    // );
    
    // await _saveTokens(tokens);
    // return AuthResult(user: UserModel.fromJson(data.user), tokens: tokens);
    
    throw UnimplementedError('Login not implemented');
  }

  @override
  Future<AuthResult> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    String? phone,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Register not implemented');
  }

  @override
  Future<AuthResult> socialLogin({
    required String provider,
    required String token,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Social login not implemented');
  }

  @override
  Future<void> logout() async {
    await _secureStorage.delete(key: _accessTokenKey);
    await _secureStorage.delete(key: _refreshTokenKey);
  }

  @override
  Future<AuthTokensModel?> refreshToken() async {
    final refreshToken = await _secureStorage.read(key: _refreshTokenKey);
    if (refreshToken == null) return null;

    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Refresh token not implemented');
  }

  @override
  Future<UserModel?> getCurrentUser() async {
    final accessToken = await _secureStorage.read(key: _accessTokenKey);
    if (accessToken == null) return null;

    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get current user not implemented');
  }

  @override
  Future<bool> isLoggedIn() async {
    final accessToken = await _secureStorage.read(key: _accessTokenKey);
    return accessToken != null;
  }

  @override
  Future<void> forgotPassword(String email) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Forgot password not implemented');
  }

  @override
  Future<void> resetPassword({
    required String token,
    required String newPassword,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Reset password not implemented');
  }

  @override
  Future<bool> verifyEmailOTP(String code) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Verify email OTP not implemented');
  }

  @override
  Future<void> resendEmailOTP() async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Resend email OTP not implemented');
  }

  @override
  Future<void> submitIdentityVerification({
    required String documentType,
    required String frontImagePath,
    String? backImagePath,
  }) async {
    // TODO: Implémenter avec GraphQL + upload multipart
    throw UnimplementedError('Submit identity verification not implemented');
  }

  Future<void> _saveTokens(AuthTokensModel tokens) async {
    await _secureStorage.write(key: _accessTokenKey, value: tokens.accessToken);
    await _secureStorage.write(key: _refreshTokenKey, value: tokens.refreshToken);
  }
}

/// Provider du repository d'authentification
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepositoryImpl();
});
