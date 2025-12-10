import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../shared/widgets/buttons/ys_button.dart';
import '../../../../shared/widgets/feedback/ys_loader.dart';
import '../../../../shared/widgets/feedback/ys_rating.dart';
import '../../../../shared/widgets/layouts/ys_scaffold.dart';
import '../../data/models/review_model.dart';
import '../providers/reviews_provider.dart';

/// Écran de création d'un avis
class CreateReviewScreen extends ConsumerStatefulWidget {
  final String offerId;
  final String offerTitle;
  final String? outingId;

  const CreateReviewScreen({
    super.key,
    required this.offerId,
    required this.offerTitle,
    this.outingId,
  });

  @override
  ConsumerState<CreateReviewScreen> createState() => _CreateReviewScreenState();
}

class _CreateReviewScreenState extends ConsumerState<CreateReviewScreen> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _contentController = TextEditingController();

  int _rating = 0;
  final List<File> _images = [];
  final ImagePicker _picker = ImagePicker();

  @override
  void dispose() {
    _titleController.dispose();
    _contentController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final createState = ref.watch(createReviewProvider);

    // Écouter les changements d'état
    ref.listen<CreateReviewState>(createReviewProvider, (previous, next) {
      next.when(
        initial: () {},
        loading: () {},
        success: (review) {
          // Afficher un message de succès et revenir
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('Avis publié avec succès !'),
              backgroundColor: AppColors.success,
            ),
          );
          Navigator.pop(context, true);
        },
        error: (message) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Erreur : $message'),
              backgroundColor: AppColors.error,
            ),
          );
        },
      );
    });

    final isLoading = createState.maybeWhen(
      loading: () => true,
      orElse: () => false,
    );

    return YsScaffold(
      title: 'Donner votre avis',
      body: Stack(
        children: [
          Form(
            key: _formKey,
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(AppSpacing.lg),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Offre concernée
                  _buildOfferInfo(),
                  const SizedBox(height: AppSpacing.xxl),

                  // Note
                  _buildRatingSection(),
                  const SizedBox(height: AppSpacing.xxl),

                  // Titre (optionnel)
                  _buildTitleField(),
                  const SizedBox(height: AppSpacing.lg),

                  // Contenu
                  _buildContentField(),
                  const SizedBox(height: AppSpacing.lg),

                  // Photos
                  _buildPhotosSection(),
                  const SizedBox(height: AppSpacing.xxl),

                  // Conditions
                  _buildTerms(),
                  const SizedBox(height: AppSpacing.xxl),

                  // Bouton publier
                  YsButton.primary(
                    label: 'Publier mon avis',
                    onPressed: _rating > 0 ? _submitReview : null,
                    isLoading: isLoading,
                    isFullWidth: true,
                  ),

                  const SizedBox(height: AppSpacing.lg),
                ],
              ),
            ),
          ),

          // Overlay de chargement
          if (isLoading) const YsLoadingOverlay(),
        ],
      ),
    );
  }

  Widget _buildOfferInfo() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.cardBackground,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Container(
            width: 50,
            height: 50,
            decoration: BoxDecoration(
              color: AppColors.primary.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: const Icon(
              Icons.local_offer_outlined,
              color: AppColors.primary,
            ),
          ),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Votre avis sur',
                  style: TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
                Text(
                  widget.offerTitle,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRatingSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Votre note',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppSpacing.sm),
        const Text(
          'Appuyez sur une étoile pour noter',
          style: TextStyle(
            fontSize: 14,
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: AppSpacing.md),
        Center(
          child: YsRatingSelector(
            initialRating: _rating,
            onRatingChanged: (rating) {
              setState(() => _rating = rating);
            },
          ),
        ),
        if (_rating > 0)
          Center(
            child: Padding(
              padding: const EdgeInsets.only(top: AppSpacing.sm),
              child: Text(
                _getRatingLabel(_rating),
                style: const TextStyle(
                  fontSize: 16,
                  color: AppColors.primary,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ),
          ),
      ],
    );
  }

  String _getRatingLabel(int rating) {
    switch (rating) {
      case 1:
        return 'Décevant';
      case 2:
        return 'Passable';
      case 3:
        return 'Correct';
      case 4:
        return 'Bien';
      case 5:
        return 'Excellent !';
      default:
        return '';
    }
  }

  Widget _buildTitleField() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Row(
          children: [
            Text(
              'Titre',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: AppColors.textPrimary,
              ),
            ),
            SizedBox(width: AppSpacing.xs),
            Text(
              '(optionnel)',
              style: TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
        const SizedBox(height: AppSpacing.sm),
        TextFormField(
          controller: _titleController,
          style: const TextStyle(color: AppColors.textPrimary),
          maxLength: 100,
          decoration: InputDecoration(
            hintText: 'Résumez votre expérience',
            hintStyle: const TextStyle(color: AppColors.textDisabled),
            filled: true,
            fillColor: AppColors.cardBackground,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: BorderSide.none,
            ),
            counterStyle: const TextStyle(color: AppColors.textSecondary),
          ),
        ),
      ],
    );
  }

  Widget _buildContentField() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Votre avis',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
            color: AppColors.textPrimary,
          ),
        ),
        const SizedBox(height: AppSpacing.sm),
        TextFormField(
          controller: _contentController,
          style: const TextStyle(color: AppColors.textPrimary),
          maxLines: 5,
          maxLength: 1000,
          decoration: InputDecoration(
            hintText: 'Partagez votre expérience avec les autres utilisateurs...',
            hintStyle: const TextStyle(color: AppColors.textDisabled),
            filled: true,
            fillColor: AppColors.cardBackground,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(8),
              borderSide: BorderSide.none,
            ),
            counterStyle: const TextStyle(color: AppColors.textSecondary),
          ),
        ),
      ],
    );
  }

  Widget _buildPhotosSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Row(
          children: [
            Text(
              'Photos',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: AppColors.textPrimary,
              ),
            ),
            SizedBox(width: AppSpacing.xs),
            Text(
              '(optionnel)',
              style: TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
            ),
          ],
        ),
        const SizedBox(height: AppSpacing.sm),
        const Text(
          'Ajoutez jusqu\'à 5 photos',
          style: TextStyle(
            fontSize: 14,
            color: AppColors.textSecondary,
          ),
        ),
        const SizedBox(height: AppSpacing.md),
        Wrap(
          spacing: AppSpacing.md,
          runSpacing: AppSpacing.md,
          children: [
            // Photos existantes
            ..._images.asMap().entries.map((entry) {
              return _buildPhotoItem(entry.key, entry.value);
            }),

            // Bouton ajouter
            if (_images.length < 5)
              GestureDetector(
                onTap: _pickImage,
                child: Container(
                  width: 80,
                  height: 80,
                  decoration: BoxDecoration(
                    color: AppColors.cardBackground,
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(
                      color: AppColors.textDisabled.withOpacity(0.3),
                      style: BorderStyle.solid,
                    ),
                  ),
                  child: const Icon(
                    Icons.add_photo_alternate_outlined,
                    color: AppColors.textSecondary,
                    size: 32,
                  ),
                ),
              ),
          ],
        ),
      ],
    );
  }

  Widget _buildPhotoItem(int index, File file) {
    return Stack(
      children: [
        ClipRRect(
          borderRadius: BorderRadius.circular(8),
          child: Image.file(
            file,
            width: 80,
            height: 80,
            fit: BoxFit.cover,
          ),
        ),
        Positioned(
          top: 4,
          right: 4,
          child: GestureDetector(
            onTap: () {
              setState(() => _images.removeAt(index));
            },
            child: Container(
              padding: const EdgeInsets.all(4),
              decoration: const BoxDecoration(
                color: AppColors.error,
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.close,
                size: 14,
                color: Colors.white,
              ),
            ),
          ),
        ),
      ],
    );
  }

  Future<void> _pickImage() async {
    final XFile? image = await _picker.pickImage(
      source: ImageSource.gallery,
      maxWidth: 1024,
      maxHeight: 1024,
      imageQuality: 80,
    );

    if (image != null) {
      setState(() {
        _images.add(File(image.path));
      });
    }
  }

  Widget _buildTerms() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.primary.withOpacity(0.05),
        borderRadius: BorderRadius.circular(8),
      ),
      child: const Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(
            Icons.info_outline,
            size: 20,
            color: AppColors.primary,
          ),
          SizedBox(width: AppSpacing.md),
          Expanded(
            child: Text(
              'En publiant votre avis, vous acceptez qu\'il soit visible par tous les utilisateurs. '
              'Les avis doivent être respectueux et conformes à nos conditions d\'utilisation.',
              style: TextStyle(
                fontSize: 12,
                color: AppColors.textSecondary,
                height: 1.5,
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _submitReview() {
    if (_rating == 0) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Veuillez donner une note'),
          backgroundColor: AppColors.warning,
        ),
      );
      return;
    }

    final params = CreateReviewParams(
      offerId: widget.offerId,
      outingId: widget.outingId,
      rating: _rating,
      title: _titleController.text.isNotEmpty ? _titleController.text : null,
      content: _contentController.text.isNotEmpty ? _contentController.text : null,
      imagePaths: _images.map((f) => f.path).toList(),
    );

    ref.read(createReviewProvider.notifier).createReview(params);
  }
}
