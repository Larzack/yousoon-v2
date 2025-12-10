import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../data/models/favorite_model.dart';
import '../../data/repositories/favorites_repository.dart';

part 'favorites_provider.freezed.dart';

/// État des favoris
@freezed
class FavoritesState with _$FavoritesState {
  const factory FavoritesState({
    @Default([]) List<FavoriteModel> favorites,
    @Default({}) Set<String> favoriteOfferIds, // Pour vérification rapide
    @Default(false) bool isLoading,
    @Default(false) bool isLoadingMore,
    String? error,
    @Default(1) int currentPage,
    @Default(false) bool hasMore,
    @Default(0) int totalCount,
  }) = _FavoritesState;
}

/// Notifier pour les favoris
class FavoritesNotifier extends StateNotifier<FavoritesState> {
  final FavoritesRepository _repository;

  FavoritesNotifier(this._repository) : super(const FavoritesState());

  /// Charger les favoris initiaux
  Future<void> loadFavorites() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final result = await _repository.getFavorites(page: 1);
      final offerIds = result.favorites.map((f) => f.offerId).toSet();

      state = state.copyWith(
        favorites: result.favorites,
        favoriteOfferIds: offerIds,
        isLoading: false,
        currentPage: result.page,
        hasMore: result.hasMore,
        totalCount: result.totalCount,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Charger plus de favoris (pagination)
  Future<void> loadMore() async {
    if (state.isLoadingMore || !state.hasMore) return;

    state = state.copyWith(isLoadingMore: true);

    try {
      final nextPage = state.currentPage + 1;
      final result = await _repository.getFavorites(page: nextPage);

      final newOfferIds = result.favorites.map((f) => f.offerId).toSet();

      state = state.copyWith(
        favorites: [...state.favorites, ...result.favorites],
        favoriteOfferIds: {...state.favoriteOfferIds, ...newOfferIds},
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

  /// Rafraîchir les favoris
  Future<void> refresh() async {
    state = state.copyWith(currentPage: 1);
    await loadFavorites();
  }

  /// Ajouter une offre aux favoris
  Future<void> addToFavorites(String offerId) async {
    // Optimistic update
    state = state.copyWith(
      favoriteOfferIds: {...state.favoriteOfferIds, offerId},
    );

    try {
      final favorite = await _repository.addToFavorites(offerId);
      state = state.copyWith(
        favorites: [favorite, ...state.favorites],
        totalCount: state.totalCount + 1,
      );
    } catch (e) {
      // Rollback
      state = state.copyWith(
        favoriteOfferIds: state.favoriteOfferIds..remove(offerId),
        error: e.toString(),
      );
    }
  }

  /// Retirer une offre des favoris
  Future<void> removeFromFavorites(String offerId) async {
    // Sauvegarde pour rollback
    final previousFavorites = state.favorites;
    final previousIds = state.favoriteOfferIds;

    // Optimistic update
    state = state.copyWith(
      favorites: state.favorites.where((f) => f.offerId != offerId).toList(),
      favoriteOfferIds: state.favoriteOfferIds..remove(offerId),
      totalCount: state.totalCount - 1,
    );

    try {
      await _repository.removeFromFavorites(offerId);
    } catch (e) {
      // Rollback
      state = state.copyWith(
        favorites: previousFavorites,
        favoriteOfferIds: previousIds,
        totalCount: state.totalCount + 1,
        error: e.toString(),
      );
    }
  }

  /// Toggle favori
  Future<void> toggleFavorite(String offerId) async {
    if (state.favoriteOfferIds.contains(offerId)) {
      await removeFromFavorites(offerId);
    } else {
      await addToFavorites(offerId);
    }
  }

  /// Vérifier si une offre est en favori
  bool isFavorite(String offerId) {
    return state.favoriteOfferIds.contains(offerId);
  }
}

/// Provider pour les favoris
final favoritesProvider =
    StateNotifierProvider<FavoritesNotifier, FavoritesState>((ref) {
  final repository = ref.watch(favoritesRepositoryProvider);
  return FavoritesNotifier(repository);
});

/// Provider pour vérifier si une offre est en favori
final isFavoriteProvider = Provider.family<bool, String>((ref, offerId) {
  final state = ref.watch(favoritesProvider);
  return state.favoriteOfferIds.contains(offerId);
});

/// Provider pour le nombre de favoris
final favoritesCountProvider = Provider<int>((ref) {
  final state = ref.watch(favoritesProvider);
  return state.totalCount;
});
