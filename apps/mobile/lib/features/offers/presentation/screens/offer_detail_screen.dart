import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:go_router/go_router.dart';
import 'package:share_plus/share_plus.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';
import '../../../../shared/widgets/buttons/ys_button.dart';

/// Écran de détail d'une offre
class OfferDetailScreen extends ConsumerStatefulWidget {
  final String offerId;

  const OfferDetailScreen({
    super.key,
    required this.offerId,
  });

  @override
  ConsumerState<OfferDetailScreen> createState() => _OfferDetailScreenState();
}

class _OfferDetailScreenState extends ConsumerState<OfferDetailScreen> {
  bool _isFavorite = false;
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: CustomScrollView(
        slivers: [
          // Image header with back button
          _buildImageHeader(),
          
          // Content
          SliverToBoxAdapter(
            child: _buildContent(),
          ),
        ],
      ),
      bottomNavigationBar: _buildBottomBar(),
    );
  }

  Widget _buildImageHeader() {
    return SliverAppBar(
      expandedHeight: 300,
      pinned: true,
      backgroundColor: AppColors.background,
      leading: _buildCircleButton(
        icon: Icons.arrow_back,
        onTap: () => context.pop(),
      ),
      actions: [
        _buildCircleButton(
          icon: _isFavorite ? Icons.favorite : Icons.favorite_border,
          iconColor: _isFavorite ? AppColors.error : AppColors.textPrimary,
          onTap: _toggleFavorite,
        ),
        _buildCircleButton(
          icon: Icons.share,
          onTap: _shareOffer,
        ),
        const SizedBox(width: AppSpacing.sm),
      ],
      flexibleSpace: FlexibleSpaceBar(
        background: Stack(
          fit: StackFit.expand,
          children: [
            CachedNetworkImage(
              imageUrl: 'https://picsum.photos/800/600?random=${widget.offerId}',
              fit: BoxFit.cover,
            ),
            Container(
              decoration: const BoxDecoration(
                gradient: AppColors.cardOverlayGradient,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCircleButton({
    required IconData icon,
    Color? iconColor,
    required VoidCallback onTap,
  }) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.sm),
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          width: 40,
          height: 40,
          decoration: BoxDecoration(
            color: AppColors.background.withOpacity(0.7),
            shape: BoxShape.circle,
          ),
          child: Icon(
            icon,
            color: iconColor ?? AppColors.textPrimary,
            size: 20,
          ),
        ),
      ),
    );
  }

  Widget _buildContent() {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.xl),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Discount badge
          Container(
            padding: const EdgeInsets.symmetric(
              horizontal: AppSpacing.md,
              vertical: AppSpacing.xs,
            ),
            decoration: BoxDecoration(
              color: AppColors.primary,
              borderRadius: BorderRadius.circular(AppSpacing.radiusXs),
            ),
            child: Text(
              '-30%',
              style: AppTypography.headline3.copyWith(
                color: AppColors.onPrimary,
                fontWeight: FontWeight.bold,
              ),
            ),
          ).animate().fadeIn().slideX(begin: -0.2),
          
          const SizedBox(height: AppSpacing.md),
          
          // Title
          Text(
            'Cocktails à moitié prix',
            style: AppTypography.headline1.copyWith(fontSize: 24),
          ).animate().fadeIn(delay: 100.ms),
          
          const SizedBox(height: AppSpacing.sm),
          
          // Partner info
          Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
                  border: Border.all(color: AppColors.border),
                ),
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
                  child: CachedNetworkImage(
                    imageUrl: 'https://picsum.photos/80/80',
                    fit: BoxFit.cover,
                  ),
                ),
              ),
              const SizedBox(width: AppSpacing.sm),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('Le Bar à Cocktails', style: AppTypography.headline3),
                  Row(
                    children: [
                      const Icon(Icons.star, color: AppColors.primary, size: 14),
                      const SizedBox(width: 2),
                      Text(
                        '4.5 (128 avis)',
                        style: AppTypography.caption.copyWith(
                          color: AppColors.textSecondary,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ],
          ).animate().fadeIn(delay: 200.ms),
          
          const SizedBox(height: AppSpacing.xl),
          
          // Description
          Text('Description', style: AppTypography.headline2),
          const SizedBox(height: AppSpacing.sm),
          Text(
            'Profitez de 30% de réduction sur tous nos cocktails signatures. Une sélection raffinée de créations originales préparées par nos mixologistes experts.',
            style: AppTypography.bodyText.copyWith(color: AppColors.textSecondary),
          ).animate().fadeIn(delay: 300.ms),
          
          const SizedBox(height: AppSpacing.xl),
          
          // Validity
          _buildInfoSection(
            icon: Icons.calendar_today,
            title: 'Validité',
            content: 'Du 10 au 31 décembre 2025',
          ).animate().fadeIn(delay: 400.ms),
          
          const SizedBox(height: AppSpacing.md),
          
          // Schedule
          _buildInfoSection(
            icon: Icons.schedule,
            title: 'Horaires',
            content: 'Du lundi au vendredi, 17h-20h',
          ).animate().fadeIn(delay: 500.ms),
          
          const SizedBox(height: AppSpacing.md),
          
          // Location
          _buildInfoSection(
            icon: Icons.location_on,
            title: 'Adresse',
            content: '123 Rue de la Soif, 75001 Paris',
          ).animate().fadeIn(delay: 600.ms),
          
          const SizedBox(height: AppSpacing.xl),
          
          // Conditions
          Text('Conditions', style: AppTypography.headline2),
          const SizedBox(height: AppSpacing.sm),
          _buildConditionItem('Minimum 2 cocktails par personne'),
          _buildConditionItem('Non cumulable avec d\'autres offres'),
          _buildConditionItem('Sur présentation du QR code'),
          
          const SizedBox(height: AppSpacing.xxxl),
        ],
      ),
    );
  }

  Widget _buildInfoSection({
    required IconData icon,
    required String title,
    required String content,
  }) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: AppColors.primary.withOpacity(0.1),
              borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
            ),
            child: Icon(icon, color: AppColors.primary, size: 20),
          ),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: AppTypography.caption.copyWith(
                  color: AppColors.textSecondary,
                )),
                Text(content, style: AppTypography.bodyText),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildConditionItem(String text) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSpacing.xs),
      child: Row(
        children: [
          const Icon(Icons.check_circle, color: AppColors.success, size: 16),
          const SizedBox(width: AppSpacing.sm),
          Expanded(
            child: Text(
              text,
              style: AppTypography.bodyText.copyWith(
                color: AppColors.textSecondary,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomBar() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.xl),
      decoration: BoxDecoration(
        color: AppColors.background,
        border: Border(
          top: BorderSide(color: AppColors.border),
        ),
      ),
      child: SafeArea(
        child: YsButton(
          label: 'Réserver cette offre',
          onPressed: _bookOffer,
        ),
      ),
    );
  }

  void _toggleFavorite() {
    setState(() {
      _isFavorite = !_isFavorite;
    });
    // TODO: Call API to toggle favorite
  }

  void _shareOffer() {
    Share.share(
      'Découvre cette super offre sur Yousoon : Cocktails à moitié prix ! https://yousoon.com/offers/${widget.offerId}',
    );
  }

  void _bookOffer() {
    context.push('/booking/${widget.offerId}');
  }
}
