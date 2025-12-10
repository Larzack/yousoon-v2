import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_animate/flutter_animate.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';

/// Écran Mes Sorties (Outings)
class MyOutingsScreen extends ConsumerStatefulWidget {
  const MyOutingsScreen({super.key});

  @override
  ConsumerState<MyOutingsScreen> createState() => _MyOutingsScreenState();
}

class _MyOutingsScreenState extends ConsumerState<MyOutingsScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  // Mock data
  final List<Map<String, dynamic>> _upcomingOutings = [
    {
      'id': '1',
      'title': '-30% sur les cocktails',
      'partner': 'Le Petit Bistrot',
      'address': '12 Rue de la Paix, Paris',
      'date': 'Aujourd\'hui',
      'time': '20:00',
      'image': 'https://picsum.photos/seed/outing1/400/200',
      'status': 'confirmed',
      'qrCode': 'ABC123',
    },
    {
      'id': '2',
      'title': 'Menu dégustation -25%',
      'partner': 'Restaurant Le Gourmet',
      'address': '45 Avenue des Champs',
      'date': 'Demain',
      'time': '19:30',
      'image': 'https://picsum.photos/seed/outing2/400/200',
      'status': 'confirmed',
      'qrCode': 'DEF456',
    },
  ];

  final List<Map<String, dynamic>> _pastOutings = [
    {
      'id': '3',
      'title': '-20% sur la carte',
      'partner': 'Brasserie du Coin',
      'address': '8 Place du Marché',
      'date': '15 nov.',
      'time': '12:30',
      'image': 'https://picsum.photos/seed/outing3/400/200',
      'status': 'checked_in',
      'hasReview': false,
    },
    {
      'id': '4',
      'title': 'Happy Hour -50%',
      'partner': 'Bar L\'Éclipse',
      'address': '22 Rue du Temple',
      'date': '10 nov.',
      'time': '18:00',
      'image': 'https://picsum.photos/seed/outing4/400/200',
      'status': 'checked_in',
      'hasReview': true,
    },
    {
      'id': '5',
      'title': 'Brunch illimité',
      'partner': 'Café Central',
      'address': '5 Boulevard Saint-Michel',
      'date': '5 nov.',
      'time': '11:00',
      'image': 'https://picsum.photos/seed/outing5/400/200',
      'status': 'no_show',
      'hasReview': false,
    },
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        backgroundColor: AppColors.background,
        elevation: 0,
        centerTitle: true,
        title: Text('Mes sorties', style: AppTypography.headline2),
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: AppColors.primary,
          labelColor: AppColors.textPrimary,
          unselectedLabelColor: AppColors.textSecondary,
          labelStyle: AppTypography.headline3,
          tabs: const [
            Tab(text: 'À venir'),
            Tab(text: 'Passées'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildOutingsList(_upcomingOutings, isUpcoming: true),
          _buildOutingsList(_pastOutings, isUpcoming: false),
        ],
      ),
    );
  }

  Widget _buildOutingsList(List<Map<String, dynamic>> outings, {required bool isUpcoming}) {
    if (outings.isEmpty) {
      return _buildEmptyState(isUpcoming);
    }

    return ListView.builder(
      padding: const EdgeInsets.all(AppSpacing.xl),
      itemCount: outings.length,
      itemBuilder: (context, index) {
        final outing = outings[index];
        return Padding(
          padding: EdgeInsets.only(bottom: AppSpacing.md),
          child: _buildOutingCard(outing, isUpcoming: isUpcoming)
              .animate()
              .fadeIn(delay: Duration(milliseconds: index * 100))
              .slideX(begin: 0.1, end: 0),
        );
      },
    );
  }

  Widget _buildEmptyState(bool isUpcoming) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.xl),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              isUpcoming ? Icons.calendar_today : Icons.history,
              size: 80,
              color: AppColors.inactive,
            ),
            const SizedBox(height: AppSpacing.lg),
            Text(
              isUpcoming ? 'Pas de sortie prévue' : 'Pas encore de sortie',
              style: AppTypography.headline2,
            ),
            const SizedBox(height: AppSpacing.sm),
            Text(
              isUpcoming
                  ? 'Réservez une offre pour planifier\nvotre prochaine sortie'
                  : 'Vos sorties passées\napparaîtront ici',
              style: AppTypography.bodyText.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
            if (isUpcoming) ...[
              const SizedBox(height: AppSpacing.xl),
              ElevatedButton(
                onPressed: () => context.go('/'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.primary,
                  foregroundColor: AppColors.onPrimary,
                  padding: const EdgeInsets.symmetric(
                    horizontal: AppSpacing.xxl,
                    vertical: AppSpacing.md,
                  ),
                ),
                child: const Text('Explorer les offres'),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildOutingCard(Map<String, dynamic> outing, {required bool isUpcoming}) {
    final status = outing['status'] as String;
    
    return GestureDetector(
      onTap: () => context.push('/outing/${outing['id']}'),
      child: Container(
        decoration: BoxDecoration(
          color: AppColors.surface,
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        ),
        clipBehavior: Clip.antiAlias,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Image
            Stack(
              children: [
                CachedNetworkImage(
                  imageUrl: outing['image'],
                  height: 120,
                  width: double.infinity,
                  fit: BoxFit.cover,
                  placeholder: (context, url) => Container(
                    height: 120,
                    color: AppColors.cardBackground,
                    child: const Center(
                      child: CircularProgressIndicator(color: AppColors.primary),
                    ),
                  ),
                ),
                // Status badge
                Positioned(
                  top: AppSpacing.sm,
                  right: AppSpacing.sm,
                  child: _buildStatusBadge(status),
                ),
              ],
            ),
            
            // Content
            Padding(
              padding: const EdgeInsets.all(AppSpacing.md),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Title
                  Text(
                    outing['title'],
                    style: AppTypography.headline3,
                  ),
                  const SizedBox(height: 4),
                  
                  // Partner
                  Text(
                    outing['partner'],
                    style: AppTypography.bodyText.copyWith(
                      color: AppColors.primary,
                    ),
                  ),
                  const SizedBox(height: AppSpacing.sm),
                  
                  // Date & Time
                  Row(
                    children: [
                      Icon(
                        Icons.calendar_today,
                        size: 16,
                        color: AppColors.textSecondary,
                      ),
                      const SizedBox(width: 6),
                      Text(
                        '${outing['date']} à ${outing['time']}',
                        style: AppTypography.bodyText.copyWith(
                          color: AppColors.textSecondary,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  
                  // Address
                  Row(
                    children: [
                      Icon(
                        Icons.location_on,
                        size: 16,
                        color: AppColors.textSecondary,
                      ),
                      const SizedBox(width: 6),
                      Expanded(
                        child: Text(
                          outing['address'],
                          style: AppTypography.bodyText.copyWith(
                            color: AppColors.textSecondary,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                    ],
                  ),
                  
                  // Actions
                  if (isUpcoming) ...[
                    const SizedBox(height: AppSpacing.md),
                    Row(
                      children: [
                        Expanded(
                          child: OutlinedButton.icon(
                            onPressed: () => _showQRCode(outing),
                            icon: const Icon(Icons.qr_code),
                            label: const Text('QR Code'),
                            style: OutlinedButton.styleFrom(
                              foregroundColor: AppColors.primary,
                              side: const BorderSide(color: AppColors.primary),
                            ),
                          ),
                        ),
                        const SizedBox(width: AppSpacing.sm),
                        Expanded(
                          child: OutlinedButton.icon(
                            onPressed: () => _cancelOuting(outing),
                            icon: const Icon(Icons.close),
                            label: const Text('Annuler'),
                            style: OutlinedButton.styleFrom(
                              foregroundColor: AppColors.error,
                              side: const BorderSide(color: AppColors.error),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ] else if (!(outing['hasReview'] ?? false) && status == 'checked_in') ...[
                    const SizedBox(height: AppSpacing.md),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton.icon(
                        onPressed: () => _leaveReview(outing),
                        icon: const Icon(Icons.star),
                        label: const Text('Laisser un avis'),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppColors.primary,
                          foregroundColor: AppColors.onPrimary,
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatusBadge(String status) {
    Color color;
    String label;
    IconData icon;

    switch (status) {
      case 'confirmed':
        color = AppColors.success;
        label = 'Confirmé';
        icon = Icons.check_circle;
        break;
      case 'checked_in':
        color = AppColors.primary;
        label = 'Utilisé';
        icon = Icons.done_all;
        break;
      case 'cancelled':
        color = AppColors.error;
        label = 'Annulé';
        icon = Icons.cancel;
        break;
      case 'no_show':
        color = AppColors.inactive;
        label = 'Non présenté';
        icon = Icons.person_off;
        break;
      default:
        color = AppColors.inactive;
        label = 'Inconnu';
        icon = Icons.help;
    }

    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSpacing.sm,
        vertical: 4,
      ),
      decoration: BoxDecoration(
        color: color.withOpacity(0.9),
        borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 14, color: Colors.white),
          const SizedBox(width: 4),
          Text(
            label,
            style: AppTypography.caption.copyWith(
              color: Colors.white,
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  void _showQRCode(Map<String, dynamic> outing) {
    showModalBottomSheet(
      context: context,
      backgroundColor: AppColors.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(AppSpacing.radiusLg)),
      ),
      builder: (context) => Padding(
        padding: const EdgeInsets.all(AppSpacing.xl),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Votre QR Code', style: AppTypography.headline2),
            const SizedBox(height: AppSpacing.lg),
            Container(
              width: 200,
              height: 200,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
              ),
              child: Center(
                child: Icon(
                  Icons.qr_code_2,
                  size: 180,
                  color: AppColors.background,
                ),
              ),
            ),
            const SizedBox(height: AppSpacing.md),
            Text(
              'Présentez ce code au partenaire',
              style: AppTypography.bodyText.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: AppSpacing.xxl),
          ],
        ),
      ),
    );
  }

  void _cancelOuting(Map<String, dynamic> outing) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: AppColors.surface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        ),
        title: Text('Annuler la réservation', style: AppTypography.headline2),
        content: Text(
          'Êtes-vous sûr de vouloir annuler cette réservation ?',
          style: AppTypography.bodyText.copyWith(color: AppColors.textSecondary),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(
              'Non',
              style: AppTypography.bodyText.copyWith(color: AppColors.textSecondary),
            ),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              // TODO: Call cancel API
            },
            child: Text(
              'Oui, annuler',
              style: AppTypography.bodyText.copyWith(color: AppColors.error),
            ),
          ),
        ],
      ),
    );
  }

  void _leaveReview(Map<String, dynamic> outing) {
    // TODO: Navigate to review screen
  }
}
