import 'package:flutter/material.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';

/// Widget d'état de chargement
class YsLoader extends StatelessWidget {
  final double size;
  final Color? color;
  final double strokeWidth;

  const YsLoader({
    super.key,
    this.size = 40,
    this.color,
    this.strokeWidth = 3,
  });

  /// Loader petit (inline)
  const YsLoader.small({super.key})
      : size = 20,
        color = null,
        strokeWidth = 2;

  /// Loader moyen
  const YsLoader.medium({super.key})
      : size = 40,
        color = null,
        strokeWidth = 3;

  /// Loader large (plein écran)
  const YsLoader.large({super.key})
      : size = 60,
        color = null,
        strokeWidth = 4;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: size,
      height: size,
      child: CircularProgressIndicator(
        strokeWidth: strokeWidth,
        valueColor: AlwaysStoppedAnimation<Color>(
          color ?? AppColors.primary,
        ),
      ),
    );
  }
}

/// Widget de chargement plein écran avec message optionnel
class YsLoadingOverlay extends StatelessWidget {
  final String? message;
  final bool isLoading;
  final Widget child;

  const YsLoadingOverlay({
    super.key,
    this.message,
    required this.isLoading,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        child,
        if (isLoading)
          Positioned.fill(
            child: Container(
              color: Colors.black54,
              child: Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const YsLoader.large(),
                    if (message != null) ...[
                      const SizedBox(height: AppSpacing.lg),
                      Text(
                        message!,
                        style: const TextStyle(
                          color: AppColors.textPrimary,
                          fontSize: 16,
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
          ),
      ],
    );
  }
}

/// Widget de chargement centré
class YsLoadingCenter extends StatelessWidget {
  final String? message;

  const YsLoadingCenter({
    super.key,
    this.message,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const YsLoader.medium(),
          if (message != null) ...[
            const SizedBox(height: AppSpacing.lg),
            Text(
              message!,
              style: const TextStyle(
                color: AppColors.textSecondary,
                fontSize: 14,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ],
      ),
    );
  }
}

/// Widget Shimmer pour les loading states
class YsShimmer extends StatefulWidget {
  final double width;
  final double height;
  final double borderRadius;

  const YsShimmer({
    super.key,
    required this.width,
    required this.height,
    this.borderRadius = 8,
  });

  /// Shimmer pour une ligne de texte
  const YsShimmer.text({
    super.key,
    this.width = 100,
    this.height = 16,
  }) : borderRadius = 4;

  /// Shimmer pour un avatar
  const YsShimmer.avatar({
    super.key,
    double size = 40,
  })  : width = size,
        height = size,
        borderRadius = 100;

  /// Shimmer pour une card
  const YsShimmer.card({
    super.key,
    this.width = double.infinity,
    this.height = 200,
  }) : borderRadius = 12;

  @override
  State<YsShimmer> createState() => _YsShimmerState();
}

class _YsShimmerState extends State<YsShimmer>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _animation;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1500),
    )..repeat();

    _animation = Tween<double>(begin: -2, end: 2).animate(
      CurvedAnimation(parent: _controller, curve: Curves.easeInOutSine),
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _animation,
      builder: (context, child) {
        return Container(
          width: widget.width,
          height: widget.height,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(widget.borderRadius),
            gradient: LinearGradient(
              begin: Alignment(_animation.value - 1, 0),
              end: Alignment(_animation.value + 1, 0),
              colors: const [
                Color(0xFF2A2A2A),
                Color(0xFF3A3A3A),
                Color(0xFF2A2A2A),
              ],
            ),
          ),
        );
      },
    );
  }
}

/// Widget skeleton pour une card d'offre
class YsOfferCardSkeleton extends StatelessWidget {
  const YsOfferCardSkeleton({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.cardBackground,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Image
          const YsShimmer.card(height: 150),
          const SizedBox(height: AppSpacing.md),
          
          // Titre
          const YsShimmer.text(width: 180, height: 20),
          const SizedBox(height: AppSpacing.sm),
          
          // Sous-titre
          const YsShimmer.text(width: 120),
          const SizedBox(height: AppSpacing.md),
          
          // Ligne du bas
          Row(
            children: [
              const YsShimmer.avatar(size: 32),
              const SizedBox(width: AppSpacing.sm),
              const Expanded(child: YsShimmer.text(width: 80)),
              const YsShimmer.text(width: 50),
            ],
          ),
        ],
      ),
    );
  }
}
