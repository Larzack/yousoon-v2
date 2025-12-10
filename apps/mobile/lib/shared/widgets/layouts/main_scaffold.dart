import 'package:flutter/material.dart';

import '../../core/theme/app_colors.dart';
import 'bottom_nav_bar.dart';

/// Scaffold principal avec navigation bottom
class MainScaffold extends StatelessWidget {
  final Widget child;

  const MainScaffold({
    super.key,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: child,
      bottomNavigationBar: const BottomNavBar(),
    );
  }
}
