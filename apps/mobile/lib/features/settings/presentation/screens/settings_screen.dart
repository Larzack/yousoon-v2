import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';
import '../../../../shared/widgets/layouts/ys_scaffold.dart';

/// Écran de paramètres
class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return YsScaffold(
      title: 'Paramètres',
      body: ListView(
        padding: const EdgeInsets.all(AppSpacing.xl),
        children: [
          // Section Compte
          _buildSectionHeader('Compte'),
          _buildSettingItem(
            icon: Icons.person,
            title: 'Modifier le profil',
            onTap: () => context.push('/settings/profile'),
          ),
          _buildSettingItem(
            icon: Icons.lock,
            title: 'Sécurité et mot de passe',
            onTap: () => context.push('/settings/security'),
          ),
          _buildSettingItem(
            icon: Icons.credit_card,
            title: 'Méthodes de paiement',
            onTap: () => context.push('/settings/payments'),
          ),
          
          const SizedBox(height: AppSpacing.xl),
          
          // Section Notifications
          _buildSectionHeader('Notifications'),
          _buildSettingSwitch(
            icon: Icons.notifications,
            title: 'Notifications push',
            value: true,
            onChanged: (value) {
              // TODO: Update notification settings
            },
          ),
          _buildSettingSwitch(
            icon: Icons.email,
            title: 'Notifications email',
            value: true,
            onChanged: (value) {},
          ),
          _buildSettingSwitch(
            icon: Icons.campaign,
            title: 'Offres marketing',
            value: false,
            onChanged: (value) {},
          ),
          
          const SizedBox(height: AppSpacing.xl),
          
          // Section Préférences
          _buildSectionHeader('Préférences'),
          _buildSettingItem(
            icon: Icons.language,
            title: 'Langue',
            trailing: Text(
              'Français',
              style: AppTypography.bodyText.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            onTap: () => _showLanguageSelector(context),
          ),
          _buildSettingItem(
            icon: Icons.location_on,
            title: 'Distance maximale',
            trailing: Text(
              '10 km',
              style: AppTypography.bodyText.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            onTap: () => _showDistanceSelector(context),
          ),
          _buildSettingItem(
            icon: Icons.category,
            title: 'Centres d\'intérêt',
            onTap: () => context.push('/settings/interests'),
          ),
          
          const SizedBox(height: AppSpacing.xl),
          
          // Section Légal
          _buildSectionHeader('Légal'),
          _buildSettingItem(
            icon: Icons.description,
            title: 'Conditions générales',
            onTap: () {},
          ),
          _buildSettingItem(
            icon: Icons.privacy_tip,
            title: 'Politique de confidentialité',
            onTap: () {},
          ),
          _buildSettingItem(
            icon: Icons.cookie,
            title: 'Gestion des cookies',
            onTap: () {},
          ),
          
          const SizedBox(height: AppSpacing.xl),
          
          // Section Danger
          _buildSectionHeader('Zone de danger'),
          _buildSettingItem(
            icon: Icons.delete_forever,
            title: 'Supprimer mon compte',
            titleColor: AppColors.error,
            onTap: () => _showDeleteAccountDialog(context),
          ),
          
          const SizedBox(height: AppSpacing.xxl),
          
          // Version
          Center(
            child: Text(
              'Version 1.0.0',
              style: AppTypography.caption.copyWith(
                color: AppColors.textDisabled,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: Text(
        title,
        style: AppTypography.headline3.copyWith(color: AppColors.textSecondary),
      ),
    );
  }

  Widget _buildSettingItem({
    required IconData icon,
    required String title,
    Widget? trailing,
    Color? titleColor,
    required VoidCallback onTap,
  }) {
    return ListTile(
      onTap: onTap,
      contentPadding: const EdgeInsets.symmetric(
        horizontal: AppSpacing.md,
        vertical: AppSpacing.xs,
      ),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
      ),
      leading: Icon(icon, color: titleColor ?? AppColors.textPrimary),
      title: Text(
        title,
        style: AppTypography.bodyText.copyWith(color: titleColor),
      ),
      trailing: trailing ??
          const Icon(
            Icons.chevron_right,
            color: AppColors.textSecondary,
          ),
    );
  }

  Widget _buildSettingSwitch({
    required IconData icon,
    required String title,
    required bool value,
    required ValueChanged<bool> onChanged,
  }) {
    return ListTile(
      contentPadding: const EdgeInsets.symmetric(
        horizontal: AppSpacing.md,
        vertical: AppSpacing.xs,
      ),
      leading: Icon(icon, color: AppColors.textPrimary),
      title: Text(title, style: AppTypography.bodyText),
      trailing: Switch(
        value: value,
        activeColor: AppColors.primary,
        onChanged: onChanged,
      ),
    );
  }

  void _showLanguageSelector(BuildContext context) {
    showModalBottomSheet(
      context: context,
      backgroundColor: AppColors.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(AppSpacing.radiusLg)),
      ),
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.all(AppSpacing.lg),
              child: Text('Choisir la langue', style: AppTypography.headline2),
            ),
            ListTile(
              title: const Text('🇫🇷  Français'),
              trailing: const Icon(Icons.check, color: AppColors.primary),
              onTap: () => Navigator.pop(context),
            ),
            ListTile(
              title: const Text('🇬🇧  English'),
              onTap: () => Navigator.pop(context),
            ),
            const SizedBox(height: AppSpacing.lg),
          ],
        ),
      ),
    );
  }

  void _showDistanceSelector(BuildContext context) {
    double distance = 10;
    showModalBottomSheet(
      context: context,
      backgroundColor: AppColors.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(AppSpacing.radiusLg)),
      ),
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => SafeArea(
          child: Padding(
            padding: const EdgeInsets.all(AppSpacing.xl),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text('Distance maximale', style: AppTypography.headline2),
                const SizedBox(height: AppSpacing.xl),
                Text(
                  '${distance.round()} km',
                  style: AppTypography.headline1.copyWith(
                    color: AppColors.primary,
                    fontSize: 32,
                  ),
                ),
                Slider(
                  value: distance,
                  min: 1,
                  max: 50,
                  divisions: 49,
                  activeColor: AppColors.primary,
                  inactiveColor: AppColors.inactive,
                  onChanged: (value) => setState(() => distance = value),
                ),
                const SizedBox(height: AppSpacing.lg),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: () => Navigator.pop(context),
                        style: OutlinedButton.styleFrom(
                          foregroundColor: AppColors.textPrimary,
                          side: const BorderSide(color: AppColors.border),
                        ),
                        child: const Text('Annuler'),
                      ),
                    ),
                    const SizedBox(width: AppSpacing.md),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: () => Navigator.pop(context),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppColors.primary,
                          foregroundColor: AppColors.onPrimary,
                        ),
                        child: const Text('Appliquer'),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  void _showDeleteAccountDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: AppColors.surface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        ),
        title: Text(
          'Supprimer mon compte',
          style: AppTypography.headline2.copyWith(color: AppColors.error),
        ),
        content: Text(
          'Cette action est irréversible. Toutes vos données seront supprimées après un délai de 30 jours.',
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
              // TODO: Call delete account API
            },
            child: Text(
              'Supprimer',
              style: AppTypography.bodyText.copyWith(color: AppColors.error),
            ),
          ),
        ],
      ),
    );
  }
}
