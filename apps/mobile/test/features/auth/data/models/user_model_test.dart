import 'package:flutter_test/flutter_test.dart';
import 'package:yousoon/features/auth/data/models/user_model.dart';

void main() {
  group('UserModel', () {
    test('should create a valid UserModel', () {
      final user = UserModel(
        id: 'user-123',
        email: 'test@example.com',
        phone: '+33612345678',
        profile: const ProfileModel(
          firstName: 'Jean',
          lastName: 'Dupont',
          displayName: 'Jean Dupont',
          avatar: 'https://example.com/avatar.jpg',
          birthDate: null,
          gender: 'male',
        ),
        identity: const IdentityStatusModel(
          status: 'verified',
          verifiedAt: null,
          documentType: 'cni',
          attemptsRemaining: 10,
        ),
        subscription: null,
        preferences: const PreferencesModel(
          language: 'fr',
          notifications: NotificationPreferencesModel(
            push: true,
            email: true,
            sms: false,
            marketing: true,
          ),
          favoriteCategories: ['sport', 'culture'],
          maxDistance: 15,
        ),
        status: 'active',
        grade: 'explorateur',
        createdAt: DateTime(2024, 1, 1),
        lastLoginAt: DateTime(2024, 12, 10),
      );

      expect(user.id, 'user-123');
      expect(user.email, 'test@example.com');
      expect(user.phone, '+33612345678');
      expect(user.profile.firstName, 'Jean');
      expect(user.profile.lastName, 'Dupont');
      expect(user.identity.status, 'verified');
      expect(user.status, 'active');
      expect(user.grade, 'explorateur');
    });

    test('should create UserModel from JSON', () {
      final json = {
        'id': 'user-456',
        'email': 'user@test.com',
        'phone': null,
        'profile': {
          'firstName': 'Marie',
          'lastName': 'Martin',
          'displayName': 'Marie Martin',
          'avatar': null,
          'birthDate': null,
          'gender': 'female',
        },
        'identity': {
          'status': 'pending',
          'verifiedAt': null,
          'documentType': null,
          'attemptsRemaining': 9,
        },
        'subscription': null,
        'preferences': {
          'language': 'en',
          'notifications': {
            'push': false,
            'email': true,
            'sms': false,
            'marketing': false,
          },
          'favoriteCategories': [],
          'maxDistance': 10,
        },
        'status': 'active',
        'grade': 'aventurier',
        'createdAt': '2024-06-15T10:30:00.000Z',
        'lastLoginAt': null,
      };

      final user = UserModel.fromJson(json);

      expect(user.id, 'user-456');
      expect(user.email, 'user@test.com');
      expect(user.phone, isNull);
      expect(user.profile.firstName, 'Marie');
      expect(user.identity.status, 'pending');
      expect(user.preferences.language, 'en');
    });

    test('should handle optional fields', () {
      final user = UserModel(
        id: 'user-789',
        email: 'minimal@test.com',
        profile: const ProfileModel(
          firstName: 'Test',
          lastName: 'User',
        ),
        identity: const IdentityStatusModel(
          status: 'not_submitted',
        ),
        preferences: const PreferencesModel(
          notifications: NotificationPreferencesModel(),
        ),
        status: 'active',
        grade: 'explorateur',
        createdAt: DateTime.now(),
      );

      expect(user.phone, isNull);
      expect(user.subscription, isNull);
      expect(user.profile.avatar, isNull);
      expect(user.profile.displayName, isNull);
      expect(user.lastLoginAt, isNull);
    });
  });

  group('ProfileModel', () {
    test('should create a valid ProfileModel', () {
      const profile = ProfileModel(
        firstName: 'Jean',
        lastName: 'Dupont',
        displayName: 'JD',
        avatar: 'https://cdn.example.com/avatar.png',
        birthDate: null,
        gender: 'male',
      );

      expect(profile.firstName, 'Jean');
      expect(profile.lastName, 'Dupont');
      expect(profile.displayName, 'JD');
      expect(profile.gender, 'male');
    });

    test('should create ProfileModel with minimal data', () {
      const profile = ProfileModel(
        firstName: 'Test',
        lastName: 'User',
      );

      expect(profile.firstName, 'Test');
      expect(profile.lastName, 'User');
      expect(profile.displayName, isNull);
      expect(profile.avatar, isNull);
      expect(profile.birthDate, isNull);
      expect(profile.gender, isNull);
    });
  });

  group('IdentityStatusModel', () {
    test('should create verified identity status', () {
      final identity = IdentityStatusModel(
        status: 'verified',
        verifiedAt: DateTime(2024, 11, 1),
        documentType: 'cni',
        attemptsRemaining: 10,
      );

      expect(identity.status, 'verified');
      expect(identity.documentType, 'cni');
      expect(identity.attemptsRemaining, 10);
    });

    test('should handle not_submitted status', () {
      const identity = IdentityStatusModel(
        status: 'not_submitted',
      );

      expect(identity.status, 'not_submitted');
      expect(identity.verifiedAt, isNull);
      expect(identity.documentType, isNull);
    });

    test('should handle rejected status', () {
      const identity = IdentityStatusModel(
        status: 'rejected',
        attemptsRemaining: 8,
      );

      expect(identity.status, 'rejected');
      expect(identity.attemptsRemaining, 8);
    });
  });

  group('SubscriptionModel', () {
    test('should create active subscription', () {
      final subscription = SubscriptionModel(
        id: 'sub-123',
        planId: 'monthly',
        planName: 'Mensuel',
        status: 'active',
        platform: 'ios',
        trialEndDate: null,
        currentPeriodStart: DateTime(2024, 12, 1),
        currentPeriodEnd: DateTime(2025, 1, 1),
        autoRenew: true,
      );

      expect(subscription.id, 'sub-123');
      expect(subscription.status, 'active');
      expect(subscription.platform, 'ios');
      expect(subscription.autoRenew, true);
    });

    test('should create trialing subscription', () {
      final subscription = SubscriptionModel(
        id: 'sub-456',
        planId: 'monthly',
        planName: 'Mensuel',
        status: 'trialing',
        platform: 'android',
        trialEndDate: DateTime(2025, 1, 10),
        currentPeriodStart: DateTime(2024, 12, 10),
        currentPeriodEnd: DateTime(2025, 1, 10),
        autoRenew: true,
      );

      expect(subscription.status, 'trialing');
      expect(subscription.trialEndDate, isNotNull);
    });
  });

  group('PreferencesModel', () {
    test('should create preferences with defaults', () {
      const prefs = PreferencesModel(
        notifications: NotificationPreferencesModel(),
      );

      expect(prefs.language, 'fr');
      expect(prefs.maxDistance, 10);
      expect(prefs.favoriteCategories, isEmpty);
    });

    test('should create preferences with custom values', () {
      const prefs = PreferencesModel(
        language: 'en',
        notifications: NotificationPreferencesModel(
          push: false,
          email: true,
          sms: true,
          marketing: false,
        ),
        favoriteCategories: ['sport', 'music', 'food'],
        maxDistance: 25,
      );

      expect(prefs.language, 'en');
      expect(prefs.maxDistance, 25);
      expect(prefs.favoriteCategories.length, 3);
      expect(prefs.notifications.push, false);
      expect(prefs.notifications.sms, true);
    });
  });

  group('NotificationPreferencesModel', () {
    test('should create with defaults', () {
      const notifs = NotificationPreferencesModel();

      expect(notifs.push, true);
      expect(notifs.email, true);
      expect(notifs.sms, false);
      expect(notifs.marketing, true);
    });

    test('should allow all disabled', () {
      const notifs = NotificationPreferencesModel(
        push: false,
        email: false,
        sms: false,
        marketing: false,
      );

      expect(notifs.push, false);
      expect(notifs.email, false);
      expect(notifs.sms, false);
      expect(notifs.marketing, false);
    });
  });

  group('AuthTokensModel', () {
    test('should create valid auth tokens', () {
      final tokens = AuthTokensModel(
        accessToken: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
        refreshToken: 'refresh-token-xyz',
        accessTokenExpiry: DateTime(2024, 12, 10, 18, 0),
        refreshTokenExpiry: DateTime(2025, 1, 10),
      );

      expect(tokens.accessToken, startsWith('eyJ'));
      expect(tokens.refreshToken, 'refresh-token-xyz');
      expect(tokens.accessTokenExpiry.hour, 18);
    });
  });
}
