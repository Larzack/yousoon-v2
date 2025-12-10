import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:yousoon/features/auth/data/models/user_model.dart';
import 'package:yousoon/features/auth/data/repositories/auth_repository.dart';
import 'package:yousoon/features/auth/presentation/providers/auth_provider.dart';

// Mocks
class MockAuthRepository extends Mock implements AuthRepository {}

void main() {
  late MockAuthRepository mockRepository;
  late ProviderContainer container;

  // Fixtures
  final testUser = UserModel(
    id: 'test-user-id',
    email: 'test@example.com',
    phone: '+33612345678',
    profile: const ProfileModel(
      firstName: 'Test',
      lastName: 'User',
      displayName: 'Test User',
    ),
    identity: const IdentityStatusModel(
      status: 'verified',
      documentType: 'cni',
    ),
    subscription: const SubscriptionModel(
      id: 'sub-123',
      planId: 'monthly',
      planName: 'Mensuel',
      status: 'active',
      platform: 'ios',
      currentPeriodStart: null,
      currentPeriodEnd: null,
      autoRenew: true,
    ),
    preferences: const PreferencesModel(
      language: 'fr',
      notifications: NotificationPreferencesModel(),
      favoriteCategories: ['sport'],
      maxDistance: 15,
    ),
    status: 'active',
    grade: 'aventurier',
    createdAt: DateTime(2024, 1, 1),
    lastLoginAt: DateTime(2024, 12, 10),
  );

  final testTokens = AuthTokensModel(
    accessToken: 'test-access-token',
    refreshToken: 'test-refresh-token',
    accessTokenExpiry: DateTime.now().add(const Duration(hours: 6)),
    refreshTokenExpiry: DateTime.now().add(const Duration(days: 30)),
  );

  setUp(() {
    mockRepository = MockAuthRepository();

    // Default stubs
    when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => false);
    when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => null);
  });

  tearDown(() {
    container.dispose();
  });

  ProviderContainer createContainer() {
    return ProviderContainer(
      overrides: [
        authRepositoryProvider.overrideWithValue(mockRepository),
      ],
    );
  }

  group('AuthNotifier', () {
    test('initial state should be unauthenticated when not logged in', () async {
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => false);
      
      container = createContainer();
      
      // Wait for initial check
      await Future.delayed(const Duration(milliseconds: 100));
      
      final state = container.read(authProvider);
      expect(state, const AuthState.unauthenticated());
    });

    test('initial state should be authenticated when logged in', () async {
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => testUser);
      
      container = createContainer();
      
      // Wait for initial check
      await Future.delayed(const Duration(milliseconds: 100));
      
      final state = container.read(authProvider);
      state.maybeWhen(
        authenticated: (user) {
          expect(user.id, testUser.id);
          expect(user.email, testUser.email);
        },
        orElse: () => fail('Expected authenticated state'),
      );
    });

    test('login should update state to authenticated on success', () async {
      when(() => mockRepository.login(
        email: any(named: 'email'),
        password: any(named: 'password'),
      )).thenAnswer((_) async => AuthResult(user: testUser, tokens: testTokens));
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      await container.read(authProvider.notifier).login(
        email: 'test@example.com',
        password: 'password123',
      );
      
      final state = container.read(authProvider);
      state.maybeWhen(
        authenticated: (user) => expect(user.email, 'test@example.com'),
        orElse: () => fail('Expected authenticated state'),
      );
    });

    test('login should update state to error on failure', () async {
      when(() => mockRepository.login(
        email: any(named: 'email'),
        password: any(named: 'password'),
      )).thenThrow(Exception('Invalid credentials'));
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      await container.read(authProvider.notifier).login(
        email: 'test@example.com',
        password: 'wrong-password',
      );
      
      final state = container.read(authProvider);
      state.maybeWhen(
        error: (message) => expect(message, contains('Invalid credentials')),
        orElse: () => fail('Expected error state'),
      );
    });

    test('register should update state to authenticated on success', () async {
      when(() => mockRepository.register(
        email: any(named: 'email'),
        password: any(named: 'password'),
        firstName: any(named: 'firstName'),
        lastName: any(named: 'lastName'),
        phone: any(named: 'phone'),
      )).thenAnswer((_) async => AuthResult(user: testUser, tokens: testTokens));
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      await container.read(authProvider.notifier).register(
        email: 'new@example.com',
        password: 'password123',
        firstName: 'New',
        lastName: 'User',
      );
      
      final state = container.read(authProvider);
      state.maybeWhen(
        authenticated: (user) => expect(user, isNotNull),
        orElse: () => fail('Expected authenticated state'),
      );
    });

    test('socialLogin should work with Google', () async {
      when(() => mockRepository.socialLogin(
        provider: any(named: 'provider'),
        token: any(named: 'token'),
      )).thenAnswer((_) async => AuthResult(user: testUser, tokens: testTokens));
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      await container.read(authProvider.notifier).socialLogin(
        provider: 'google',
        token: 'google-oauth-token',
      );
      
      verify(() => mockRepository.socialLogin(
        provider: 'google',
        token: 'google-oauth-token',
      )).called(1);
    });

    test('logout should update state to unauthenticated', () async {
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => testUser);
      when(() => mockRepository.logout()).thenAnswer((_) async {});
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      // Verify initially authenticated
      expect(container.read(authProvider), isA<AuthState>());
      
      await container.read(authProvider.notifier).logout();
      
      final state = container.read(authProvider);
      expect(state, const AuthState.unauthenticated());
      verify(() => mockRepository.logout()).called(1);
    });

    test('updateUser should update the current user', () async {
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => testUser);
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      final updatedUser = testUser.copyWith(grade: 'conquerant');
      container.read(authProvider.notifier).updateUser(updatedUser);
      
      final state = container.read(authProvider);
      state.maybeWhen(
        authenticated: (user) => expect(user.grade, 'conquerant'),
        orElse: () => fail('Expected authenticated state'),
      );
    });
  });

  group('Derived Providers', () {
    test('isAuthenticatedProvider should return true when authenticated', () async {
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => testUser);
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      expect(container.read(isAuthenticatedProvider), true);
    });

    test('isAuthenticatedProvider should return false when unauthenticated', () async {
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      expect(container.read(isAuthenticatedProvider), false);
    });

    test('currentUserProvider should return user when authenticated', () async {
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => testUser);
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      final user = container.read(currentUserProvider);
      expect(user?.id, testUser.id);
    });

    test('currentUserProvider should return null when unauthenticated', () async {
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      expect(container.read(currentUserProvider), isNull);
    });

    test('isIdentityVerifiedProvider should return true when verified', () async {
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => testUser);
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      expect(container.read(isIdentityVerifiedProvider), true);
    });

    test('isIdentityVerifiedProvider should return false when not verified', () async {
      final unverifiedUser = testUser.copyWith(
        identity: const IdentityStatusModel(status: 'pending'),
      );
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => unverifiedUser);
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      expect(container.read(isIdentityVerifiedProvider), false);
    });

    test('hasActiveSubscriptionProvider should return true for active subscription', () async {
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => testUser);
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      expect(container.read(hasActiveSubscriptionProvider), true);
    });

    test('hasActiveSubscriptionProvider should return true for trialing', () async {
      final trialUser = testUser.copyWith(
        subscription: testUser.subscription?.copyWith(status: 'trialing'),
      );
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => trialUser);
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      expect(container.read(hasActiveSubscriptionProvider), true);
    });

    test('hasActiveSubscriptionProvider should return false without subscription', () async {
      final noSubUser = testUser.copyWith(subscription: null);
      when(() => mockRepository.isLoggedIn()).thenAnswer((_) async => true);
      when(() => mockRepository.getCurrentUser()).thenAnswer((_) async => noSubUser);
      
      container = createContainer();
      await Future.delayed(const Duration(milliseconds: 100));
      
      expect(container.read(hasActiveSubscriptionProvider), false);
    });
  });

  group('AuthState', () {
    test('initial state should match', () {
      const state = AuthState.initial();
      state.when(
        initial: () => expect(true, true),
        loading: () => fail('Should be initial'),
        authenticated: (_) => fail('Should be initial'),
        unauthenticated: () => fail('Should be initial'),
        error: (_) => fail('Should be initial'),
      );
    });

    test('loading state should match', () {
      const state = AuthState.loading();
      state.maybeWhen(
        loading: () => expect(true, true),
        orElse: () => fail('Should be loading'),
      );
    });

    test('authenticated state should contain user', () {
      final state = AuthState.authenticated(testUser);
      state.maybeWhen(
        authenticated: (user) => expect(user.id, testUser.id),
        orElse: () => fail('Should be authenticated'),
      );
    });

    test('error state should contain message', () {
      const state = AuthState.error('Something went wrong');
      state.maybeWhen(
        error: (msg) => expect(msg, 'Something went wrong'),
        orElse: () => fail('Should be error'),
      );
    });
  });
}
