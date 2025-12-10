import 'package:flutter/material.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/theme/app_typography.dart';

/// Scaffold personnalisé Yousoon
/// Fournit une structure cohérente pour tous les écrans
class YsScaffold extends StatelessWidget {
  final String? title;
  final Widget body;
  final List<Widget>? actions;
  final Widget? floatingActionButton;
  final FloatingActionButtonLocation? floatingActionButtonLocation;
  final Widget? bottomNavigationBar;
  final bool showBackButton;
  final VoidCallback? onBackPressed;
  final bool extendBodyBehindAppBar;
  final Color? backgroundColor;
  final PreferredSizeWidget? customAppBar;

  const YsScaffold({
    super.key,
    this.title,
    required this.body,
    this.actions,
    this.floatingActionButton,
    this.floatingActionButtonLocation,
    this.bottomNavigationBar,
    this.showBackButton = true,
    this.onBackPressed,
    this.extendBodyBehindAppBar = false,
    this.backgroundColor,
    this.customAppBar,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: backgroundColor ?? AppColors.background,
      extendBodyBehindAppBar: extendBodyBehindAppBar,
      appBar: customAppBar ?? _buildAppBar(context),
      body: body,
      floatingActionButton: floatingActionButton,
      floatingActionButtonLocation: floatingActionButtonLocation,
      bottomNavigationBar: bottomNavigationBar,
    );
  }

  PreferredSizeWidget? _buildAppBar(BuildContext context) {
    if (title == null && actions == null && !showBackButton) {
      return null;
    }

    return AppBar(
      backgroundColor: AppColors.background,
      elevation: 0,
      scrolledUnderElevation: 0,
      centerTitle: true,
      leading: showBackButton
          ? IconButton(
              icon: const Icon(Icons.arrow_back, color: AppColors.textPrimary),
              onPressed: onBackPressed ?? () => Navigator.of(context).pop(),
            )
          : null,
      automaticallyImplyLeading: showBackButton,
      title: title != null
          ? Text(
              title!,
              style: AppTypography.headline2,
            )
          : null,
      actions: actions,
    );
  }
}

/// Scaffold avec TabBar bottom navigation
class YsTabScaffold extends StatelessWidget {
  final int currentIndex;
  final ValueChanged<int> onTabChanged;
  final List<Widget> children;

  const YsTabScaffold({
    super.key,
    required this.currentIndex,
    required this.onTabChanged,
    required this.children,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: IndexedStack(
        index: currentIndex,
        children: children,
      ),
      bottomNavigationBar: _buildBottomNav(),
    );
  }

  Widget _buildBottomNav() {
    return Container(
      decoration: BoxDecoration(
        color: AppColors.surface,
        border: Border(
          top: BorderSide(color: AppColors.border, width: 0.5),
        ),
      ),
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(
            horizontal: AppSpacing.md,
            vertical: AppSpacing.sm,
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              _buildNavItem(0, Icons.calendar_today, 'Mes events'),
              _buildNavItem(1, Icons.favorite, 'Favoris'),
              _buildNavItem(2, Icons.auto_awesome, 'Pour vous'),
              _buildNavItem(3, Icons.map, 'Carte'),
              _buildNavItem(4, Icons.chat_bubble_outline, 'Messages'),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildNavItem(int index, IconData icon, String label) {
    final isSelected = currentIndex == index;
    
    return GestureDetector(
      onTap: () => onTabChanged(index),
      behavior: HitTestBehavior.opaque,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            icon,
            color: isSelected ? AppColors.primary : AppColors.inactive,
            size: 24,
          ),
          const SizedBox(height: 4),
          Text(
            label,
            style: AppTypography.caption.copyWith(
              color: isSelected ? AppColors.primary : AppColors.inactive,
              fontSize: 10,
            ),
          ),
          const SizedBox(height: 2),
          Container(
            width: 4,
            height: 4,
            decoration: BoxDecoration(
              color: isSelected ? AppColors.primary : Colors.transparent,
              shape: BoxShape.circle,
            ),
          ),
        ],
      ),
    );
  }
}
