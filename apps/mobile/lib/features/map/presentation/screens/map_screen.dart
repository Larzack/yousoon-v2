import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:geolocator/geolocator.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';

/// Écran carte avec les offres à proximité
class MapScreen extends ConsumerStatefulWidget {
  const MapScreen({super.key});

  @override
  ConsumerState<MapScreen> createState() => _MapScreenState();
}

class _MapScreenState extends ConsumerState<MapScreen> {
  GoogleMapController? _mapController;
  Position? _currentPosition;
  String? _selectedMarkerId;
  bool _isLoading = true;
  
  final Set<Marker> _markers = {};
  
  // Paris par défaut
  static const _defaultLocation = LatLng(48.8566, 2.3522);

  @override
  void initState() {
    super.initState();
    _getCurrentLocation();
    _loadOffers();
  }

  @override
  void dispose() {
    _mapController?.dispose();
    super.dispose();
  }

  Future<void> _getCurrentLocation() async {
    try {
      final permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        await Geolocator.requestPermission();
      }
      
      final position = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
      );
      
      setState(() {
        _currentPosition = position;
        _isLoading = false;
      });
      
      _mapController?.animateCamera(
        CameraUpdate.newLatLng(
          LatLng(position.latitude, position.longitude),
        ),
      );
    } catch (e) {
      setState(() => _isLoading = false);
    }
  }

  void _loadOffers() {
    // TODO: Load offers from API
    // For now, add some fake markers
    final fakeOffers = [
      {'id': '1', 'title': 'Bar à cocktails', 'lat': 48.8566, 'lng': 2.3522, 'discount': '-30%'},
      {'id': '2', 'title': 'Restaurant italien', 'lat': 48.8606, 'lng': 2.3376, 'discount': '-20%'},
      {'id': '3', 'title': 'Spa & Bien-être', 'lat': 48.8530, 'lng': 2.3499, 'discount': '-40%'},
      {'id': '4', 'title': 'Cinéma', 'lat': 48.8619, 'lng': 2.3478, 'discount': '-15%'},
      {'id': '5', 'title': 'Escape Game', 'lat': 48.8555, 'lng': 2.3600, 'discount': '-25%'},
    ];
    
    for (final offer in fakeOffers) {
      _markers.add(
        Marker(
          markerId: MarkerId(offer['id'] as String),
          position: LatLng(offer['lat'] as double, offer['lng'] as double),
          icon: BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueOrange),
          onTap: () => _onMarkerTap(offer['id'] as String),
        ),
      );
    }
    setState(() {});
  }

  void _onMarkerTap(String markerId) {
    setState(() => _selectedMarkerId = markerId);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: Stack(
        children: [
          // Map
          GoogleMap(
            initialCameraPosition: CameraPosition(
              target: _currentPosition != null
                  ? LatLng(_currentPosition!.latitude, _currentPosition!.longitude)
                  : _defaultLocation,
              zoom: 14,
            ),
            onMapCreated: (controller) => _mapController = controller,
            markers: _markers,
            myLocationEnabled: true,
            myLocationButtonEnabled: false,
            zoomControlsEnabled: false,
            mapToolbarEnabled: false,
            style: _mapStyle,
            onTap: (_) => setState(() => _selectedMarkerId = null),
          ),
          
          // Search bar overlay
          SafeArea(
            child: Padding(
              padding: const EdgeInsets.all(AppSpacing.xl),
              child: Column(
                children: [
                  // Search bar
                  Container(
                    decoration: BoxDecoration(
                      color: AppColors.surface,
                      borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.2),
                          blurRadius: 10,
                          offset: const Offset(0, 2),
                        ),
                      ],
                    ),
                    child: TextField(
                      style: AppTypography.bodyText,
                      decoration: InputDecoration(
                        hintText: 'Rechercher un lieu...',
                        hintStyle: AppTypography.bodyText.copyWith(
                          color: AppColors.textDisabled,
                        ),
                        prefixIcon: const Icon(Icons.search, color: AppColors.textDisabled),
                        suffixIcon: IconButton(
                          icon: const Icon(Icons.tune, color: AppColors.primary),
                          onPressed: () {
                            // TODO: Show filters
                          },
                        ),
                        border: InputBorder.none,
                        contentPadding: const EdgeInsets.symmetric(
                          horizontal: AppSpacing.md,
                          vertical: AppSpacing.md,
                        ),
                      ),
                    ),
                  ),
                  
                  // Category chips
                  const SizedBox(height: AppSpacing.md),
                  SizedBox(
                    height: 36,
                    child: ListView(
                      scrollDirection: Axis.horizontal,
                      children: [
                        _buildCategoryChip('Tout', '✨', true),
                        _buildCategoryChip('Bars', '🍹', false),
                        _buildCategoryChip('Restos', '🍽️', false),
                        _buildCategoryChip('Loisirs', '🎮', false),
                        _buildCategoryChip('Sport', '🎾', false),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          
          // My location button
          Positioned(
            right: AppSpacing.xl,
            bottom: _selectedMarkerId != null ? 220 : 120,
            child: FloatingActionButton(
              mini: true,
              backgroundColor: AppColors.surface,
              onPressed: _getCurrentLocation,
              child: const Icon(Icons.my_location, color: AppColors.primary),
            ),
          ),
          
          // Selected offer card
          if (_selectedMarkerId != null)
            Positioned(
              left: AppSpacing.xl,
              right: AppSpacing.xl,
              bottom: AppSpacing.xl,
              child: _buildOfferCard(),
            ),
        ],
      ),
    );
  }

  Widget _buildCategoryChip(String label, String emoji, bool isSelected) {
    return Padding(
      padding: const EdgeInsets.only(right: AppSpacing.sm),
      child: FilterChip(
        selected: isSelected,
        label: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(emoji),
            const SizedBox(width: 4),
            Text(label),
          ],
        ),
        labelStyle: AppTypography.caption.copyWith(
          color: isSelected ? AppColors.onPrimary : AppColors.textPrimary,
        ),
        backgroundColor: AppColors.surface,
        selectedColor: AppColors.primary,
        side: BorderSide(
          color: isSelected ? AppColors.primary : AppColors.border,
        ),
        onSelected: (_) {
          // TODO: Filter by category
        },
      ),
    );
  }

  Widget _buildOfferCard() {
    return GestureDetector(
      onTap: () => context.push('/offers/$_selectedMarkerId'),
      child: Container(
        padding: const EdgeInsets.all(AppSpacing.md),
        decoration: BoxDecoration(
          color: AppColors.surface,
          borderRadius: BorderRadius.circular(AppSpacing.radiusMd),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.3),
              blurRadius: 15,
              offset: const Offset(0, 5),
            ),
          ],
        ),
        child: Row(
          children: [
            // Image placeholder
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: AppColors.primary.withOpacity(0.1),
                borderRadius: BorderRadius.circular(AppSpacing.radiusSm),
              ),
              child: const Icon(Icons.local_bar, color: AppColors.primary, size: 36),
            ),
            
            const SizedBox(width: AppSpacing.md),
            
            // Content
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          'Bar à cocktails',
                          style: AppTypography.headline3,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: AppSpacing.sm,
                          vertical: 2,
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
                  const SizedBox(height: 4),
                  Text(
                    'Le Bar à Cocktails',
                    style: AppTypography.caption.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      const Icon(Icons.star, color: AppColors.primary, size: 14),
                      const SizedBox(width: 2),
                      Text(
                        '4.5',
                        style: AppTypography.caption,
                      ),
                      const SizedBox(width: AppSpacing.sm),
                      const Icon(Icons.location_on, color: AppColors.textSecondary, size: 14),
                      const SizedBox(width: 2),
                      Text(
                        '350m',
                        style: AppTypography.caption.copyWith(
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
      ),
    );
  }

  // Dark map style
  static const String _mapStyle = '''
[
  {
    "elementType": "geometry",
    "stylers": [{"color": "#1d2c4d"}]
  },
  {
    "elementType": "labels.text.fill",
    "stylers": [{"color": "#8ec3b9"}]
  },
  {
    "elementType": "labels.text.stroke",
    "stylers": [{"color": "#1a3646"}]
  },
  {
    "featureType": "administrative.country",
    "elementType": "geometry.stroke",
    "stylers": [{"color": "#4b6878"}]
  },
  {
    "featureType": "road",
    "elementType": "geometry",
    "stylers": [{"color": "#304a7d"}]
  },
  {
    "featureType": "road",
    "elementType": "geometry.stroke",
    "stylers": [{"color": "#255763"}]
  },
  {
    "featureType": "water",
    "elementType": "geometry",
    "stylers": [{"color": "#0e1626"}]
  }
]
''';
}
