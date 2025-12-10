import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_animate/flutter_animate.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';
import '../../../../shared/widgets/layouts/ys_scaffold.dart';
import '../../../../shared/widgets/buttons/ys_button.dart';

/// Écran de profil utilisateur
class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return YsScaffold(
      title: 'Mon profil',
      showBackButton: false,
      actions: [
        IconButton(
          icon: const Icon(Icons.settings, color: AppColors.textPrimary),
          onPressed: () => context.push('/settings'),
        ),
      ],
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(AppSpacing.xl),
        child: Column(
          children: [
            // Profile header
            _buildProfileHeader().animate().fadeIn(),
            
            const SizedBox(height: AppSpacing.xl),
            
            // Stats
            _buildStats().animate().fadeIn(delay: 100.ms),
            
            const SizedBox(height: AppSpacing.xl),
            
            // Menu items
            _buildMenuItem(
              icon: Icons.receipt_long,
              title: 'Mes réservations',
              subtitle: '12 sorties effectuées',
              onTap: () => context.push('/bookings'),
            ).animate().fadeIn(delay: 200.ms),
            
            _buildMenuItem(
              icon: Icons.favorite,
              title: 'Mes favoris',
              subtitle: '5 offres sauvegardées',
              onTap: () => context.push('/favorites'),
            ).animate().fadeIn(delay: 300.ms),
            
            _buildMenuItem(
              icon: Icons.star,
              title: 'Mes avis',
              subtitle: '8 avis publiés',
              onTap: () => context.push('/reviews'),
            ).animate().fadeIn(delay: 400.ms),
            
            _buildMenuItem(
              icon: Icons.card_membership,
              title: 'Mon abonnement',
              subtitle: 'Premium - Actif',
              trailing: Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                decoration: BoxDecoration(
                  color: AppColors.success.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(AppSpacing.radiusXs),
                ),
                child: Text(
                  'Actif',
                  style: AppTypography.caption.copyWith(
                    color: AppColors.success,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              onTap: () => context.push('/subscription'),
            ).animate().fadeIn(delay: 500.ms),
            
            _buildMenuItem(
              icon: Icons.verified_user,
              title: 'Vérification d\'identité',
              subtitle: 'Compte vérifié',
              trailing: const Icon(Icons.check_circle, color: AppColors.success, size: 20),
              onTap: () {},
            ).animate().fadeIn(delay: 600.ms),
            
            const SizedBox(height: AppSpacing.xxl),
            
            // Logout button
            YsButton(
              label: 'Se déconnecter',
              variant: YsButtonVariant.secondary,
              onPressed: () => _showLogoutDialog(context),
            ).animate().fadeIn(delay: 700.ms),
          ],
        ),
      ),
    );
  }

  Widget _buildProfileHeader() {
    return Column(
      children: [
        // Avatar
        Stack(
          children: [
            Container(
              width: 100,
              height: 100,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                border: Border.all(color: AppColors.primary, width: 3),
              ),
              child: ClipOval(
                child: CachedNetworkImage(
                  imageUrl: 'https://i.pravatar.cc/200',
                  fit: BoxFit.cover,
                  placeholder: (_, __) => Container(
                    color: AppColors.surface,
                    child: const Icon(Icons.person, color: AppColors.inactive, size: 50),
                  ),
                ),
              ),
            ),
            Positioned(
              bottom: 0,
              right: 0,
              child: Container(
                width: 32,
                height: 32,
                decoration: BoxDecoration(
                  color: AppColors.primary,
                  shape: BoxShape.circle,
                  border: Border.all(color: AppColors.background, width: 2),
                ),
                child: const Icon(Icons.edit, color: AppColors.onPrimary, size: 16),
              ),
            ),
          ],
        ),
        
        const SizedBox(height: AppSpacing.md),
        
        // Name
        Text(
          'Jean Dupont',
          style: AppTypography.headline1.copyWith(fontSize: 22),
        ),
        
        const SizedBox(height: 4),
        
        // Email
        Text(
          'jean.dupont@example.com',
          style: AppTypography.bodyText.copyWith(color: AppColors.textSecondary),
        ),
        
        const SizedBox(height: AppSpacing.sm),
        
        // Grade badge
        Container(
          padding: const EdgeInsets.symmetric(
            horizontal: AppSpacing.md,
            vertical: AppSpacing.xs,
          ),
          decoration: BoxDecoration(
            gradient: AppColors.primaryGradient,
            borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Text('👑', style: TextStyle(fontSize: 14)),
              const SizedBox(width: 4),
              Text(
                'Conquérant',
                style: AppTypography.caption.copyWith(
                  color: AppColors.onPrimary,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildStats() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: [
          _buildStatItem('12', 'Sorties'),
          Container(width: 1, height: 40, color: AppColors.border),
          _buildStatItem('€156', 'Économisés'),
          Container(width: 1, height: 40, color: AppColors.border),
          _buildStatItem('5', 'Favoris'),
        ],
      ),
    );
  }

  Widget _buildStatItem(String value, String label) {
    return Column(
      children: [
        Text(
          value,
          style: AppTypography.headline1.copyWith(
            color: AppColors.primary,
            fontSize: 20,
          ),
        ),
        const SizedBox(height: 2),
        Text(
          label,
          style: AppTypography.caption.copyWith(
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  Widget _buildMenuItem({
    required IconData icon,
    required String title,
    required String subtitle,
    Widget? trailing,
    required VoidCallback onTap,
  }) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: ListTile(
        onTap: onTap,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.md,
          vertical: AppSpacing.xs,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
          side: BorderSide(color: AppColors.border),
        ),
        tileColor: AppColors.surface,
        leading: Container(
          width: 44,
          height: 44,
          decoration: BoxDecoration(
            color: AppColors.primary.withOpacity(0.1),
            borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
          ),
          child: Icon(icon, color: AppColors.primary),
        ),
        title: Text(title, style: AppTypography.headline3),
        subtitle: Text(
          subtitle,
          style: AppTypography.caption.copyWith(color: AppColors.textSecondary),
        ),
        trailing: trailing ?? const Icon(
          Icons.chevron_right,
          color: AppColors.textSecondary,
        ),
      ),
    );
  }

  void _showLogoutDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: AppColors.surface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        ),
        title: Text('Se déconnecter', style: AppTypography.headline2),
        content: Text(
          'Êtes-vous sûr de vouloir vous déconnecter ?',
          style: AppTypography.bodyText.copyWith(color: AppColors.textSecondary),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(
              'Annuler',
              style: AppTypography.bodyText.copyWith(color: AppColors.textSecondary),
            ),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              context.go('/login');
            },
            child: Text(
              'Déconnexion',
              style: AppTypography.bodyText.copyWith(color: AppColors.error),
            ),
          ),
        ],
      ),
    );
  }
}
