import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_animate/flutter_animate.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';
import '../../../../shared/widgets/layouts/ys_scaffold.dart';

/// Écran des notifications
class NotificationsScreen extends ConsumerStatefulWidget {
  const NotificationsScreen({super.key});

  @override
  ConsumerState<NotificationsScreen> createState() => _NotificationsScreenState();
}

class _NotificationsScreenState extends ConsumerState<NotificationsScreen> {
  // Mock data pour les notifications
  final List<Map<String, dynamic>> _notifications = [
    {
      'id': '1',
      'type': 'booking_confirmed',
      'title': 'Réservation confirmée',
      'message': 'Votre réservation au Petit Bistrot est confirmée pour ce soir à 20h.',
      'time': 'Il y a 5 min',
      'isRead': false,
      'icon': Icons.check_circle,
      'color': AppColors.success,
    },
    {
      'id': '2',
      'type': 'offer_nearby',
      'title': 'Nouvelle offre à proximité',
      'message': '-30% sur les cocktails au Bar L\'Éclipse, à seulement 500m de vous !',
      'time': 'Il y a 1h',
      'isRead': false,
      'icon': Icons.local_offer,
      'color': AppColors.primary,
    },
    {
      'id': '3',
      'type': 'reminder',
      'title': 'Rappel de réservation',
      'message': 'N\'oubliez pas votre sortie demain à 19h au Restaurant Le Gourmet.',
      'time': 'Il y a 3h',
      'isRead': true,
      'icon': Icons.alarm,
      'color': AppColors.warning,
    },
    {
      'id': '4',
      'type': 'grade_up',
      'title': 'Félicitations ! 🎉',
      'message': 'Vous êtes passé au grade Aventurier ! Continuez à explorer.',
      'time': 'Hier',
      'isRead': true,
      'icon': Icons.star,
      'color': AppColors.primary,
    },
    {
      'id': '5',
      'type': 'marketing',
      'title': 'Offre spéciale',
      'message': 'Ce weekend, profitez de -20% sur toutes les offres avec le code WEEKEND20.',
      'time': 'Il y a 2 jours',
      'isRead': true,
      'icon': Icons.campaign,
      'color': AppColors.info,
    },
  ];

  @override
  Widget build(BuildContext context) {
    return YsScaffold(
      title: 'Notifications',
      actions: [
        TextButton(
          onPressed: _markAllAsRead,
          child: Text(
            'Tout lire',
            style: AppTypography.bodyText.copyWith(color: AppColors.primary),
          ),
        ),
      ],
      body: _notifications.isEmpty
          ? _buildEmptyState()
          : ListView.separated(
              itemCount: _notifications.length,
              padding: const EdgeInsets.symmetric(vertical: AppSpacing.sm),
              separatorBuilder: (context, index) => Divider(
                color: AppColors.divider,
                height: 1,
                indent: AppSpacing.xl + 48,
                endIndent: AppSpacing.xl,
              ),
              itemBuilder: (context, index) {
                final notification = _notifications[index];
                return _buildNotificationTile(notification, index);
              },
            ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.xl),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.notifications_off_outlined,
              size: 80,
              color: AppColors.inactive,
            ),
            const SizedBox(height: AppSpacing.lg),
            Text(
              'Pas de notifications',
              style: AppTypography.headline2,
            ),
            const SizedBox(height: AppSpacing.sm),
            Text(
              'Vous recevrez ici les notifications\nde vos réservations et offres',
              style: AppTypography.bodyText.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildNotificationTile(Map<String, dynamic> notification, int index) {
    final isRead = notification['isRead'] as bool;
    
    return Dismissible(
      key: Key(notification['id']),
      direction: DismissDirection.endToStart,
      background: Container(
        alignment: Alignment.centerRight,
        padding: const EdgeInsets.only(right: AppSpacing.xl),
        color: AppColors.error,
        child: const Icon(Icons.delete, color: Colors.white),
      ),
      onDismissed: (_) {
        setState(() {
          _notifications.removeAt(index);
        });
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: const Text('Notification supprimée'),
            backgroundColor: AppColors.surface,
            action: SnackBarAction(
              label: 'Annuler',
              textColor: AppColors.primary,
              onPressed: () {
                setState(() {
                  _notifications.insert(index, notification);
                });
              },
            ),
          ),
        );
      },
      child: ListTile(
        onTap: () => _handleNotificationTap(notification),
        contentPadding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.xl,
          vertical: AppSpacing.sm,
        ),
        tileColor: isRead ? Colors.transparent : AppColors.surface,
        leading: Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: (notification['color'] as Color).withOpacity(0.1),
            shape: BoxShape.circle,
          ),
          child: Icon(
            notification['icon'] as IconData,
            color: notification['color'] as Color,
          ),
        ),
        title: Row(
          children: [
            Expanded(
              child: Text(
                notification['title'],
                style: AppTypography.headline3.copyWith(
                  fontWeight: isRead ? FontWeight.normal : FontWeight.bold,
                ),
              ),
            ),
            if (!isRead)
              Container(
                width: 8,
                height: 8,
                decoration: const BoxDecoration(
                  color: AppColors.primary,
                  shape: BoxShape.circle,
                ),
              ),
          ],
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 4),
            Text(
              notification['message'],
              style: AppTypography.bodyText.copyWith(
                color: AppColors.textSecondary,
              ),
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 4),
            Text(
              notification['time'],
              style: AppTypography.caption.copyWith(
                color: AppColors.textDisabled,
              ),
            ),
          ],
        ),
      ).animate().fadeIn(
        delay: Duration(milliseconds: index * 50),
      ),
    );
  }

  void _markAllAsRead() {
    setState(() {
      for (final notification in _notifications) {
        notification['isRead'] = true;
      }
    });
  }

  void _handleNotificationTap(Map<String, dynamic> notification) {
    // Mark as read
    setState(() {
      notification['isRead'] = true;
    });

    // Navigate based on type
    switch (notification['type']) {
      case 'booking_confirmed':
      case 'reminder':
        // TODO: Navigate to outing detail
        break;
      case 'offer_nearby':
        // TODO: Navigate to offer detail
        break;
      case 'grade_up':
        // TODO: Navigate to profile
        break;
      default:
        break;
    }
  }
}
