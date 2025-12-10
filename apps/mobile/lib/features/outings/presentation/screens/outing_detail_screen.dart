import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:share_plus/share_plus.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';

/// Écran de détail d'une sortie (Outing)
class OutingDetailScreen extends ConsumerStatefulWidget {
  final String outingId;

  const OutingDetailScreen({
    super.key,
    required this.outingId,
  });

  @override
  ConsumerState<OutingDetailScreen> createState() => _OutingDetailScreenState();
}

class _OutingDetailScreenState extends ConsumerState<OutingDetailScreen> {
  // Mock data
  late Map<String, dynamic> _outing;

  @override
  void initState() {
    super.initState();
    _loadOuting();
  }

  void _loadOuting() {
    // TODO: Fetch from API
    _outing = {
      'id': widget.outingId,
      'title': '-30% sur les cocktails',
      'description': 'Profitez de 30% de réduction sur tous nos cocktails signature. Cette offre est valable du lundi au jeudi, de 18h à 22h.',
      'partner': {
        'name': 'Le Petit Bistrot',
        'logo': 'https://picsum.photos/seed/partner1/100',
        'rating': 4.5,
        'reviewCount': 128,
      },
      'establishment': {
        'name': 'Le Petit Bistrot - Marais',
        'address': '12 Rue de la Paix, 75002 Paris',
        'phone': '+33 1 23 45 67 89',
      },
      'date': 'Aujourd\'hui, 14 décembre 2024',
      'time': '20:00',
      'status': 'confirmed',
      'qrCode': 'YOUSOON-ABC123-XYZ789',
      'image': 'https://picsum.photos/seed/outing1/800/400',
      'bookedAt': 'Il y a 2 heures',
      'expiresAt': '20:30',
      'conditions': [
        'Valable sur les cocktails uniquement',
        'Non cumulable avec d\'autres offres',
        'Présenter le QR code à l\'arrivée',
      ],
    };
  }

  @override
  Widget build(BuildContext context) {
    final status = _outing['status'] as String;
    final isUpcoming = status == 'confirmed' || status == 'pending';

    return Scaffold(
      backgroundColor: AppColors.background,
      body: CustomScrollView(
        slivers: [
          // App Bar avec image
          _buildSliverAppBar(),

          // Content
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.all(AppSpacing.xl),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Status badge
                  _buildStatusSection(status),
                  const SizedBox(height: AppSpacing.lg),

                  // Title
                  Text(
                    _outing['title'],
                    style: AppTypography.headline1.copyWith(fontSize: 24),
                  ).animate().fadeIn().slideY(begin: 0.2, end: 0),
                  const SizedBox(height: AppSpacing.md),

                  // Partner info
                  _buildPartnerSection(),
                  const SizedBox(height: AppSpacing.xl),

                  // QR Code section (only for upcoming)
                  if (isUpcoming) ...[
                    _buildQRCodeSection(),
                    const SizedBox(height: AppSpacing.xl),
                  ],

                  // Date & Time
                  _buildInfoSection(
                    icon: Icons.calendar_today,
                    title: 'Date et heure',
                    content: '${_outing['date']} à ${_outing['time']}',
                  ),
                  const SizedBox(height: AppSpacing.md),

                  // Location
                  _buildInfoSection(
                    icon: Icons.location_on,
                    title: 'Adresse',
                    content: _outing['establishment']['address'],
                    action: TextButton(
                      onPressed: _openMaps,
                      child: Text(
                        'Itinéraire',
                        style: AppTypography.bodyText.copyWith(
                          color: AppColors.primary,
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(height: AppSpacing.md),

                  // Phone
                  _buildInfoSection(
                    icon: Icons.phone,
                    title: 'Téléphone',
                    content: _outing['establishment']['phone'],
                    action: TextButton(
                      onPressed: _callPartner,
                      child: Text(
                        'Appeler',
                        style: AppTypography.bodyText.copyWith(
                          color: AppColors.primary,
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(height: AppSpacing.xl),

                  // Description
                  Text('Description', style: AppTypography.headline3),
                  const SizedBox(height: AppSpacing.sm),
                  Text(
                    _outing['description'],
                    style: AppTypography.bodyText.copyWith(
                      color: AppColors.textSecondary,
                      height: 1.6,
                    ),
                  ),
                  const SizedBox(height: AppSpacing.xl),

                  // Conditions
                  Text('Conditions', style: AppTypography.headline3),
                  const SizedBox(height: AppSpacing.sm),
                  ...(_outing['conditions'] as List<String>).map(
                    (condition) => Padding(
                      padding: const EdgeInsets.only(bottom: AppSpacing.xs),
                      child: Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text('• ', style: TextStyle(color: AppColors.primary)),
                          Expanded(
                            child: Text(
                              condition,
                              style: AppTypography.bodyText.copyWith(
                                color: AppColors.textSecondary,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: AppSpacing.xxl),

                  // Actions
                  if (isUpcoming) _buildActions(),

                  // Leave review (for past outings)
                  if (!isUpcoming && status == 'checked_in')
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton.icon(
                        onPressed: _leaveReview,
                        icon: const Icon(Icons.star),
                        label: const Text('Laisser un avis'),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppColors.primary,
                          foregroundColor: AppColors.onPrimary,
                          padding: const EdgeInsets.symmetric(vertical: AppSpacing.md),
                        ),
                      ),
                    ),

                  const SizedBox(height: AppSpacing.xl),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSliverAppBar() {
    return SliverAppBar(
      expandedHeight: 200,
      pinned: true,
      backgroundColor: AppColors.background,
      leading: IconButton(
        icon: Container(
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: AppColors.overlay,
            shape: BoxShape.circle,
          ),
          child: const Icon(Icons.arrow_back, color: Colors.white),
        ),
        onPressed: () => context.pop(),
      ),
      actions: [
        IconButton(
          icon: Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: AppColors.overlay,
              shape: BoxShape.circle,
            ),
            child: const Icon(Icons.share, color: Colors.white),
          ),
          onPressed: _shareOuting,
        ),
      ],
      flexibleSpace: FlexibleSpaceBar(
        background: CachedNetworkImage(
          imageUrl: _outing['image'],
          fit: BoxFit.cover,
          placeholder: (context, url) => Container(
            color: AppColors.surface,
            child: const Center(
              child: CircularProgressIndicator(color: AppColors.primary),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildStatusSection(String status) {
    Color color;
    String label;
    IconData icon;
    String subtitle;

    switch (status) {
      case 'confirmed':
        color = AppColors.success;
        label = 'Réservation confirmée';
        icon = Icons.check_circle;
        subtitle = 'Réservé ${_outing['bookedAt']}';
        break;
      case 'pending':
        color = AppColors.warning;
        label = 'En attente de confirmation';
        icon = Icons.hourglass_empty;
        subtitle = 'Le partenaire doit valider';
        break;
      case 'checked_in':
        color = AppColors.primary;
        label = 'Sortie effectuée';
        icon = Icons.done_all;
        subtitle = 'Check-in validé';
        break;
      case 'cancelled':
        color = AppColors.error;
        label = 'Réservation annulée';
        icon = Icons.cancel;
        subtitle = '';
        break;
      case 'no_show':
        color = AppColors.inactive;
        label = 'Non présenté';
        icon = Icons.person_off;
        subtitle = '';
        break;
      default:
        color = AppColors.inactive;
        label = 'Statut inconnu';
        icon = Icons.help;
        subtitle = '';
    }

    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        border: Border.all(color: color.withOpacity(0.3)),
      ),
      child: Row(
        children: [
          Icon(icon, color: color, size: 28),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: AppTypography.headline3.copyWith(color: color),
                ),
                if (subtitle.isNotEmpty)
                  Text(
                    subtitle,
                    style: AppTypography.caption.copyWith(
                      color: color.withOpacity(0.8),
                    ),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPartnerSection() {
    final partner = _outing['partner'] as Map<String, dynamic>;

    return GestureDetector(
      onTap: () {
        // TODO: Navigate to partner detail
      },
      child: Row(
        children: [
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: CachedNetworkImage(
              imageUrl: partner['logo'],
              width: 48,
              height: 48,
              fit: BoxFit.cover,
            ),
          ),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(partner['name'], style: AppTypography.headline3),
                Row(
                  children: [
                    const Icon(Icons.star, color: AppColors.primary, size: 16),
                    const SizedBox(width: 4),
                    Text(
                      '${partner['rating']} (${partner['reviewCount']} avis)',
                      style: AppTypography.bodyText.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          const Icon(Icons.chevron_right, color: AppColors.textSecondary),
        ],
      ),
    );
  }

  Widget _buildQRCodeSection() {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
      ),
      child: Column(
        children: [
          Text(
            'Votre QR Code',
            style: AppTypography.headline3.copyWith(color: Colors.black),
          ),
          const SizedBox(height: AppSpacing.md),
          // Placeholder for QR code
          Container(
            width: 180,
            height: 180,
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
            ),
            child: Icon(
              Icons.qr_code_2,
              size: 160,
              color: Colors.black,
            ),
          ),
          const SizedBox(height: AppSpacing.md),
          Text(
            'Présentez ce code au partenaire',
            style: AppTypography.bodyText.copyWith(color: Colors.black54),
          ),
          const SizedBox(height: AppSpacing.sm),
          Text(
            'Valide jusqu\'à ${_outing['expiresAt']}',
            style: AppTypography.caption.copyWith(color: Colors.black38),
          ),
        ],
      ),
    ).animate().scale(begin: const Offset(0.95, 0.95), end: const Offset(1, 1));
  }

  Widget _buildInfoSection({
    required IconData icon,
    required String title,
    required String content,
    Widget? action,
  }) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Container(
          width: 40,
          height: 40,
          decoration: BoxDecoration(
            color: AppColors.surface,
            borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
          ),
          child: Icon(icon, color: AppColors.primary, size: 20),
        ),
        const SizedBox(width: AppSpacing.md),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                title,
                style: AppTypography.caption.copyWith(
                  color: AppColors.textSecondary,
                ),
              ),
              const SizedBox(height: 2),
              Text(content, style: AppTypography.bodyText),
            ],
          ),
        ),
        if (action != null) action,
      ],
    );
  }

  Widget _buildActions() {
    return Column(
      children: [
        SizedBox(
          width: double.infinity,
          child: ElevatedButton.icon(
            onPressed: _showFullQRCode,
            icon: const Icon(Icons.qr_code),
            label: const Text('Afficher le QR Code'),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.primary,
              foregroundColor: AppColors.onPrimary,
              padding: const EdgeInsets.symmetric(vertical: AppSpacing.md),
            ),
          ),
        ),
        const SizedBox(height: AppSpacing.sm),
        SizedBox(
          width: double.infinity,
          child: OutlinedButton.icon(
            onPressed: _cancelOuting,
            icon: const Icon(Icons.close),
            label: const Text('Annuler la réservation'),
            style: OutlinedButton.styleFrom(
              foregroundColor: AppColors.error,
              side: const BorderSide(color: AppColors.error),
              padding: const EdgeInsets.symmetric(vertical: AppSpacing.md),
            ),
          ),
        ),
      ],
    );
  }

  void _showFullQRCode() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(AppSpacing.radiusLg)),
      ),
      builder: (context) => SizedBox(
        height: MediaQuery.of(context).size.height * 0.7,
        child: Column(
          children: [
            const SizedBox(height: AppSpacing.lg),
            Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Colors.grey.shade300,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            const SizedBox(height: AppSpacing.xl),
            Text(
              'Votre QR Code',
              style: AppTypography.headline1.copyWith(color: Colors.black),
            ),
            const SizedBox(height: AppSpacing.sm),
            Text(
              _outing['partner']['name'],
              style: AppTypography.bodyText.copyWith(color: Colors.black54),
            ),
            const Spacer(),
            Container(
              width: 250,
              height: 250,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.1),
                    blurRadius: 10,
                    spreadRadius: 5,
                  ),
                ],
              ),
              child: Icon(
                Icons.qr_code_2,
                size: 230,
                color: Colors.black,
              ),
            ),
            const SizedBox(height: AppSpacing.lg),
            Text(
              _outing['qrCode'],
              style: AppTypography.caption.copyWith(
                color: Colors.black38,
                letterSpacing: 2,
              ),
            ),
            const Spacer(),
            Padding(
              padding: const EdgeInsets.all(AppSpacing.xl),
              child: Text(
                'Présentez ce QR code au partenaire\npour valider votre sortie',
                style: AppTypography.bodyText.copyWith(color: Colors.black54),
                textAlign: TextAlign.center,
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _shareOuting() {
    Share.share(
      'Je vais profiter de "${_outing['title']}" chez ${_outing['partner']['name']} ! 🎉\n\nDécouvre Yousoon pour des sorties à prix réduit.',
      subject: 'Ma sortie Yousoon',
    );
  }

  void _openMaps() {
    // TODO: Open maps with address
  }

  void _callPartner() {
    // TODO: Call partner phone
  }

  void _cancelOuting() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: AppColors.surface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        ),
        title: Text('Annuler la réservation', style: AppTypography.headline2),
        content: Text(
          'Êtes-vous sûr de vouloir annuler cette réservation ? Cette action est irréversible.',
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
              context.go('/outings');
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

  void _leaveReview() {
    // TODO: Navigate to review screen
  }
}
