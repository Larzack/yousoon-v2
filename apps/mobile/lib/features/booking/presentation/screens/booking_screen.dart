import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:flutter_animate/flutter_animate.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';
import '../../../../shared/widgets/buttons/ys_button.dart';
import '../../../../shared/widgets/layouts/ys_scaffold.dart';

/// Écran de réservation d'une offre
class BookingScreen extends ConsumerStatefulWidget {
  final String offerId;

  const BookingScreen({
    super.key,
    required this.offerId,
  });

  @override
  ConsumerState<BookingScreen> createState() => _BookingScreenState();
}

class _BookingScreenState extends ConsumerState<BookingScreen> {
  DateTime _selectedDate = DateTime.now();
  String? _selectedTimeSlot;
  int _guestCount = 1;
  bool _isBooking = false;

  final List<String> _timeSlots = [
    '17:00',
    '18:00',
    '19:00',
    '20:00',
  ];

  @override
  Widget build(BuildContext context) {
    return YsScaffold(
      title: 'Réserver',
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(AppSpacing.xl),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Offer summary
            _buildOfferSummary().animate().fadeIn(),
            
            const SizedBox(height: AppSpacing.xl),
            
            // Date selection
            _buildDateSection().animate().fadeIn(delay: 100.ms),
            
            const SizedBox(height: AppSpacing.xl),
            
            // Time selection
            _buildTimeSection().animate().fadeIn(delay: 200.ms),
            
            const SizedBox(height: AppSpacing.xl),
            
            // Guest count
            _buildGuestSection().animate().fadeIn(delay: 300.ms),
            
            const SizedBox(height: AppSpacing.xxl),
            
            // Book button
            YsButton(
              label: 'Confirmer la réservation',
              isLoading: _isBooking,
              onPressed: _selectedTimeSlot != null ? _confirmBooking : null,
            ).animate().fadeIn(delay: 400.ms),
            
            const SizedBox(height: AppSpacing.lg),
            
            // Disclaimer
            Text(
              'Vous avez 30 minutes après la réservation pour effectuer le check-in.',
              style: AppTypography.caption.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildOfferSummary() {
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
            width: 60,
            height: 60,
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
              color: AppColors.primary.withOpacity(0.1),
            ),
            child: const Icon(Icons.local_bar, color: AppColors.primary, size: 30),
          ),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Cocktails à moitié prix', style: AppTypography.headline3),
                Text(
                  'Le Bar à Cocktails',
                  style: AppTypography.caption.copyWith(color: AppColors.textSecondary),
                ),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(
              horizontal: AppSpacing.sm,
              vertical: AppSpacing.xs,
            ),
            decoration: BoxDecoration(
              color: AppColors.primary,
              borderRadius: BorderRadius.circular(AppSpacing.radiusXs),
            ),
            child: Text(
              '-30%',
              style: AppTypography.caption.copyWith(
                color: AppColors.onPrimary,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDateSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Date', style: AppTypography.headline2),
        const SizedBox(height: AppSpacing.sm),
        SizedBox(
          height: 80,
          child: ListView.separated(
            scrollDirection: Axis.horizontal,
            itemCount: 7,
            separatorBuilder: (_, __) => const SizedBox(width: AppSpacing.sm),
            itemBuilder: (context, index) {
              final date = DateTime.now().add(Duration(days: index));
              final isSelected = _selectedDate.day == date.day;
              
              return GestureDetector(
                onTap: () => setState(() => _selectedDate = date),
                child: Container(
                  width: 60,
                  decoration: BoxDecoration(
                    color: isSelected ? AppColors.primary : AppColors.surface,
                    borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
                    border: Border.all(
                      color: isSelected ? AppColors.primary : AppColors.border,
                    ),
                  ),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        _getDayName(date),
                        style: AppTypography.caption.copyWith(
                          color: isSelected ? AppColors.onPrimary : AppColors.textSecondary,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '${date.day}',
                        style: AppTypography.headline2.copyWith(
                          color: isSelected ? AppColors.onPrimary : AppColors.textPrimary,
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _buildTimeSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Créneau horaire', style: AppTypography.headline2),
        const SizedBox(height: AppSpacing.sm),
        Wrap(
          spacing: AppSpacing.sm,
          runSpacing: AppSpacing.sm,
          children: _timeSlots.map((slot) {
            final isSelected = _selectedTimeSlot == slot;
            return GestureDetector(
              onTap: () => setState(() => _selectedTimeSlot = slot),
              child: Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: AppSpacing.lg,
                  vertical: AppSpacing.md,
                ),
                decoration: BoxDecoration(
                  color: isSelected ? AppColors.primary : AppColors.surface,
                  borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
                  border: Border.all(
                    color: isSelected ? AppColors.primary : AppColors.border,
                  ),
                ),
                child: Text(
                  slot,
                  style: AppTypography.bodyText.copyWith(
                    color: isSelected ? AppColors.onPrimary : AppColors.textPrimary,
                    fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                  ),
                ),
              ),
            );
          }).toList(),
        ),
      ],
    );
  }

  Widget _buildGuestSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Nombre de personnes', style: AppTypography.headline2),
        const SizedBox(height: AppSpacing.sm),
        Container(
          padding: const EdgeInsets.all(AppSpacing.md),
          decoration: BoxDecoration(
            color: AppColors.surface,
            borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
            border: Border.all(color: AppColors.border),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              IconButton(
                onPressed: _guestCount > 1
                    ? () => setState(() => _guestCount--)
                    : null,
                icon: Icon(
                  Icons.remove_circle_outline,
                  color: _guestCount > 1 ? AppColors.primary : AppColors.inactive,
                ),
              ),
              Text(
                '$_guestCount',
                style: AppTypography.headline1.copyWith(fontSize: 24),
              ),
              IconButton(
                onPressed: _guestCount < 10
                    ? () => setState(() => _guestCount++)
                    : null,
                icon: Icon(
                  Icons.add_circle_outline,
                  color: _guestCount < 10 ? AppColors.primary : AppColors.inactive,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  String _getDayName(DateTime date) {
    final days = ['Dim', 'Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam'];
    return days[date.weekday % 7];
  }

  Future<void> _confirmBooking() async {
    setState(() => _isBooking = true);
    
    // Simulate API call
    await Future.delayed(const Duration(seconds: 2));
    
    setState(() => _isBooking = false);
    
    if (mounted) {
      context.go('/booking-confirmation');
    }
  }
}

/// Écran de confirmation de réservation avec QR code
class BookingConfirmationScreen extends StatelessWidget {
  const BookingConfirmationScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return YsScaffold(
      title: 'Réservation confirmée',
      showBackButton: false,
      body: Center(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(AppSpacing.xl),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              // Success icon
              Container(
                width: 80,
                height: 80,
                decoration: BoxDecoration(
                  color: AppColors.success.withOpacity(0.1),
                  shape: BoxShape.circle,
                ),
                child: const Icon(
                  Icons.check_circle,
                  color: AppColors.success,
                  size: 48,
                ),
              ).animate().scale(delay: 200.ms),
              
              const SizedBox(height: AppSpacing.xl),
              
              Text(
                'Réservation confirmée !',
                style: AppTypography.headline1.copyWith(fontSize: 24),
                textAlign: TextAlign.center,
              ).animate().fadeIn(delay: 400.ms),
              
              const SizedBox(height: AppSpacing.sm),
              
              Text(
                'Présentez ce QR code au partenaire pour valider votre réduction.',
                style: AppTypography.bodyText.copyWith(
                  color: AppColors.textSecondary,
                ),
                textAlign: TextAlign.center,
              ).animate().fadeIn(delay: 500.ms),
              
              const SizedBox(height: AppSpacing.xxl),
              
              // QR Code
              Container(
                padding: const EdgeInsets.all(AppSpacing.lg),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
                ),
                child: QrImageView(
                  data: 'yousoon:outing:abc123',
                  version: QrVersions.auto,
                  size: 200,
                  backgroundColor: Colors.white,
                ),
              ).animate().fadeIn(delay: 600.ms).scale(),
              
              const SizedBox(height: AppSpacing.lg),
              
              // Timer warning
              Container(
                padding: const EdgeInsets.all(AppSpacing.md),
                decoration: BoxDecoration(
                  color: AppColors.warning.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
                  border: Border.all(color: AppColors.warning.withOpacity(0.3)),
                ),
                child: Row(
                  children: [
                    const Icon(Icons.timer, color: AppColors.warning),
                    const SizedBox(width: AppSpacing.sm),
                    Expanded(
                      child: Text(
                        'Valide pendant 30 minutes',
                        style: AppTypography.bodyText.copyWith(
                          color: AppColors.warning,
                        ),
                      ),
                    ),
                  ],
                ),
              ).animate().fadeIn(delay: 700.ms),
              
              const SizedBox(height: AppSpacing.xxl),
              
              // Details
              _buildDetailRow('Offre', 'Cocktails à moitié prix'),
              _buildDetailRow('Lieu', 'Le Bar à Cocktails'),
              _buildDetailRow('Date', '10 décembre 2025'),
              _buildDetailRow('Heure', '18:00'),
              _buildDetailRow('Personnes', '2'),
              
              const SizedBox(height: AppSpacing.xxl),
              
              YsButton(
                label: 'Retour à l\'accueil',
                onPressed: () => context.go('/'),
              ).animate().fadeIn(delay: 800.ms),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: AppSpacing.xs),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: AppTypography.bodyText.copyWith(color: AppColors.textSecondary),
          ),
          Text(value, style: AppTypography.bodyText),
        ],
      ),
    );
  }
}
