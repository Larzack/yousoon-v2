import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../features/auth/presentation/screens/login_screen.dart';
import '../features/auth/presentation/screens/register_screen.dart';
import '../features/auth/presentation/screens/onboarding_screen.dart';
import '../features/auth/presentation/screens/splash_screen.dart';
import '../features/auth/presentation/screens/identity_verification_screen.dart';
import '../features/home/presentation/screens/home_screen.dart';
import '../features/offers/presentation/screens/offer_detail_screen.dart';
import '../features/offers/presentation/screens/search_screen.dart';
import '../features/outings/presentation/screens/outing_detail_screen.dart';
import '../features/outings/presentation/screens/my_outings_screen.dart';
import '../features/favorites/presentation/screens/favorites_screen.dart';
import '../features/map/presentation/screens/map_screen.dart';
import '../features/messages/presentation/screens/messages_screen.dart';
import '../features/profile/presentation/screens/profile_screen.dart';
import '../features/settings/presentation/screens/settings_screen.dart';
import '../features/notifications/presentation/screens/notifications_screen.dart';
import '../shared/widgets/layouts/main_scaffold.dart';

final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/splash',
    debugLogDiagnostics: true,
    routes: [
      // Auth routes
      GoRoute(
        path: '/splash',
        name: 'splash',
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(
        path: '/onboarding',
        name: 'onboarding',
        builder: (context, state) => const OnboardingScreen(),
      ),
      GoRoute(
        path: '/login',
        name: 'login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/register',
        name: 'register',
        builder: (context, state) => const RegisterScreen(),
      ),
      GoRoute(
        path: '/identity-verification',
        name: 'identity-verification',
        builder: (context, state) => const IdentityVerificationScreen(),
      ),

      // Main app routes with bottom navigation
      ShellRoute(
        builder: (context, state, child) => MainScaffold(child: child),
        routes: [
          GoRoute(
            path: '/',
            name: 'home',
            builder: (context, state) => const HomeScreen(),
          ),
          GoRoute(
            path: '/outings',
            name: 'outings',
            builder: (context, state) => const MyOutingsScreen(),
          ),
          GoRoute(
            path: '/favorites',
            name: 'favorites',
            builder: (context, state) => const FavoritesScreen(),
          ),
          GoRoute(
            path: '/map',
            name: 'map',
            builder: (context, state) => const MapScreen(),
          ),
          GoRoute(
            path: '/messages',
            name: 'messages',
            builder: (context, state) => const MessagesScreen(),
          ),
        ],
      ),

      // Detail routes
      GoRoute(
        path: '/offer/:id',
        name: 'offer-detail',
        builder: (context, state) => OfferDetailScreen(
          offerId: state.pathParameters['id']!,
        ),
      ),
      GoRoute(
        path: '/outing/:id',
        name: 'outing-detail',
        builder: (context, state) => OutingDetailScreen(
          outingId: state.pathParameters['id']!,
        ),
      ),
      GoRoute(
        path: '/search',
        name: 'search',
        builder: (context, state) => const SearchScreen(),
      ),

      // Profile routes
      GoRoute(
        path: '/profile',
        name: 'profile',
        builder: (context, state) => const ProfileScreen(),
      ),
      GoRoute(
        path: '/settings',
        name: 'settings',
        builder: (context, state) => const SettingsScreen(),
      ),
      GoRoute(
        path: '/notifications',
        name: 'notifications',
        builder: (context, state) => const NotificationsScreen(),
      ),
    ],
    errorBuilder: (context, state) => Scaffold(
      backgroundColor: Colors.black,
      body: Center(
        child: Text(
          'Page non trouvée',
          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                color: Colors.white,
              ),
        ),
      ),
    ),
  );
});
