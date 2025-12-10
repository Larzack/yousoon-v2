import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../data/models/offer_model.dart';
import '../../data/repositories/offers_repository.dart';

part 'offers_provider.freezed.dart';

/// État de chargement des offres
@freezed
class OffersState with _$OffersState {
  const factory OffersState({
    @Default([]) List<OfferModel> offers,
    @Default(false) bool isLoading,
    @Default(false) bool isLoadingMore,
    String? error,
    @Default(1) int currentPage,
    @Default(false) bool hasMore,
    OfferSearchParams? searchParams,
  }) = _OffersState;
}

/// Notifier pour les offres (liste/recherche)
class OffersNotifier extends StateNotifier<OffersState> {
  final OffersRepository _repository;

  OffersNotifier(this._repository) : super(const OffersState());

  /// Charger les offres initiales
  Future<void> loadOffers(OfferSearchParams params) async {
    state = state.copyWith(
      isLoading: true,
      error: null,
      searchParams: params,
    );

    try {
      final result = await _repository.searchOffers(params);
      state = state.copyWith(
        offers: result.offers,
        isLoading: false,
        currentPage: result.page,
        hasMore: result.hasMore,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Charger plus d'offres (pagination)
  Future<void> loadMore() async {
    if (state.isLoadingMore || !state.hasMore || state.searchParams == null) {
      return;
    }

    state = state.copyWith(isLoadingMore: true);

    try {
      final nextPage = state.currentPage + 1;
      final params = state.searchParams!.copyWith(page: nextPage);
      final result = await _repository.searchOffers(params);

      state = state.copyWith(
        offers: [...state.offers, ...result.offers],
        isLoadingMore: false,
        currentPage: result.page,
        hasMore: result.hasMore,
      );
    } catch (e) {
      state = state.copyWith(
        isLoadingMore: false,
        error: e.toString(),
      );
    }
  }

  /// Rafraîchir les offres
  Future<void> refresh() async {
    if (state.searchParams == null) return;

    final params = state.searchParams!.copyWith(page: 1);
    await loadOffers(params);
  }
}

/// Provider pour les offres (recherche/liste)
final offersProvider =
    StateNotifierProvider<OffersNotifier, OffersState>((ref) {
  final repository = ref.watch(offersRepositoryProvider);
  return OffersNotifier(repository);
});

/// État du détail d'une offre
@freezed
class OfferDetailState with _$OfferDetailState {
  const factory OfferDetailState.loading() = _Loading;
  const factory OfferDetailState.loaded(OfferModel offer) = _Loaded;
  const factory OfferDetailState.error(String message) = _Error;
}

/// Notifier pour le détail d'une offre
class OfferDetailNotifier extends StateNotifier<OfferDetailState> {
  final OffersRepository _repository;
  final String offerId;

  OfferDetailNotifier(this._repository, this.offerId)
      : super(const OfferDetailState.loading()) {
    _loadOffer();
  }

  Future<void> _loadOffer() async {
    try {
      final offer = await _repository.getOfferById(offerId);
      state = OfferDetailState.loaded(offer);
    } catch (e) {
      state = OfferDetailState.error(e.toString());
    }
  }

  Future<void> refresh() async {
    state = const OfferDetailState.loading();
    await _loadOffer();
  }
}

/// Provider famille pour le détail d'une offre
final offerDetailProvider = StateNotifierProvider.family<OfferDetailNotifier,
    OfferDetailState, String>((ref, offerId) {
  final repository = ref.watch(offersRepositoryProvider);
  return OfferDetailNotifier(repository, offerId);
});

/// Provider pour les offres recommandées ("Pour vous")
final recommendedOffersProvider =
    FutureProvider.family<List<OfferModel>, ({double lat, double lng})>(
        (ref, location) async {
  final repository = ref.watch(offersRepositoryProvider);
  return repository.getRecommendedOffers(
    latitude: location.lat,
    longitude: location.lng,
  );
});

/// Provider pour les offres à proximité
final nearbyOffersProvider =
    FutureProvider.family<List<OfferModel>, ({double lat, double lng, double radius})>(
        (ref, params) async {
  final repository = ref.watch(offersRepositoryProvider);
  return repository.getNearbyOffers(
    latitude: params.lat,
    longitude: params.lng,
    radius: params.radius,
  );
});

/// Provider pour les catégories
final categoriesProvider = FutureProvider<List<CategoryModel>>((ref) async {
  final repository = ref.watch(offersRepositoryProvider);
  return repository.getCategories();
});

/// Provider pour les offres par catégorie
final offersByCategoryProvider = StateNotifierProvider.family<OffersNotifier,
    OffersState, String>((ref, categoryId) {
  final repository = ref.watch(offersRepositoryProvider);
  return OffersNotifier(repository);
});
