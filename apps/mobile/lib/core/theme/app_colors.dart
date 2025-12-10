import 'package:flutter/material.dart';

/// Palette de couleurs Yousoon
/// Basée sur le Design System Figma
abstract class AppColors {
  // Couleurs principales
  static const Color primary = Color(0xFFE99B27);      // Indian Gold
  static const Color onPrimary = Color(0xFF000000);    // Noir sur primary
  static const Color secondary = Color(0xFFFFFFFF);    // Flash White
  static const Color onSecondary = Color(0xFF000000);  // Noir sur secondary
  
  // Background & Surface
  static const Color background = Color(0xFF000000);   // Dark Black
  static const Color surface = Color(0xFF1A1A1A);      // Noir légèrement plus clair
  static const Color onSurface = Color(0xFFFFFFFF);    // Blanc sur surface
  static const Color cardBackground = Color(0xFF1A1A1A);
  
  // Texte
  static const Color textPrimary = Color(0xFFFFFFFF);  // Flash White
  static const Color textSecondary = Color(0xFFCCCCCC); // Eerie Black
  static const Color textDisabled = Color(0xFF6D6D6D); // Grey Jet
  
  // Feedback
  static const Color success = Color(0xFF5FC15C);      // Mantis Green
  static const Color error = Color(0xFFCC2936);        // Persian Red
  static const Color warning = Color(0xFFE99B27);      // Indian Gold
  static const Color info = Color(0xFF3B82F6);         // Blue
  
  // Éléments UI
  static const Color inactive = Color(0xFF6D6D6D);     // Grey Jet
  static const Color divider = Color(0xFF333333);      // Gris foncé
  static const Color border = Color(0xFF333333);
  
  // Overlay
  static const Color overlay = Color(0x80000000);      // 50% noir
  static const Color overlayLight = Color(0x40000000); // 25% noir
  
  // Gradients
  static const LinearGradient primaryGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [
      Color(0xFFE99B27),
      Color(0xFFD4890F),
    ],
  );
  
  static const LinearGradient cardOverlayGradient = LinearGradient(
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
    colors: [
      Color(0x00000000),
      Color(0xCC000000),
    ],
  );
  
  // Grades Yousooner
  static const Color gradeExplorateur = Color(0xFF6D6D6D);
  static const Color gradeAventurier = Color(0xFF5FC15C);
  static const Color gradeGrandVoyageur = Color(0xFF3B82F6);
  static const Color gradeConquerant = Color(0xFFE99B27);
}
