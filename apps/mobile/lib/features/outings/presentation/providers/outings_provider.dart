import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../data/models/outing_model.dart';
import '../../data/repositories/outings_repository.dart';

part 'outings_provider.freezed.dart';

/// État des sorties
@freezed
class OutingsState with _$OutingsState {
  const factory OutingsState({
    @Default([]) List<OutingModel> upcoming,
    @Default([]) List<OutingModel> past,
    @Default([]) List<OutingModel> cancelled,
    @Default(false) bool isLoading,
    @Default(false) bool isLoadingMore,
    String? error,
    @Default(1) int currentPage,
    @Default(false) bool hasMore,
    @Default(OutingFilter.upcoming) OutingFilter currentFilter,
  }) = _OutingsState;
}

/// Notifier pour les sorties
class OutingsNotifier extends StateNotifier<OutingsState> {
  final OutingsRepository _repository;

  OutingsNotifier(this._repository) : super(const OutingsState());

  /// Charger toutes les sorties
  Future<void> loadOutings() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      // Charger les sorties à venir
      final upcoming = await _repository.getUpcomingOutings();

      // Charger l'historique
      final pastResult = await _repository.getOutingHistory(page: 1);
      final past = pastResult.outings
          .where((o) => o.status != 'cancelled')
          .toList();
      final cancelled = pastResult.outings
          .where((o) => o.status == 'cancelled')
          .toList();

      state = state.copyWith(
        upcoming: upcoming,
        past: past,
        cancelled: cancelled,
        isLoading: false,
        hasMore: pastResult.hasMore,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Changer le filtre actif
  void setFilter(OutingFilter filter) {
    state = state.copyWith(currentFilter: filter);
  }

  /// Charger plus de sorties (pagination)
  Future<void> loadMore() async {
    if (state.isLoadingMore || !state.hasMore) return;

    state = state.copyWith(isLoadingMore: true);

    try {
      final nextPage = state.currentPage + 1;
      final result = await _repository.getOutingHistory(page: nextPage);

      final newPast = result.outings
          .where((o) => o.status != 'cancelled')
          .toList();
      final newCancelled = result.outings
          .where((o) => o.status == 'cancelled')
          .toList();

      state = state.copyWith(
        past: [...state.past, ...newPast],
        cancelled: [...state.cancelled, ...newCancelled],
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

  /// Rafraîchir les sorties
  Future<void> refresh() async {
    state = state.copyWith(currentPage: 1);
    await loadOutings();
  }

  /// Annuler une sortie
  Future<void> cancelOuting(String id, {String? reason}) async {
    try {
      final updated = await _repository.cancelOuting(id, reason: reason);

      // Mettre à jour la liste
      final newUpcoming = state.upcoming.where((o) => o.id != id).toList();
      final newCancelled = [...state.cancelled, updated];

      state = state.copyWith(
        upcoming: newUpcoming,
        cancelled: newCancelled,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }
}

/// Provider pour les sorties
final outingsProvider =
    StateNotifierProvider<OutingsNotifier, OutingsState>((ref) {
  final repository = ref.watch(outingsRepositoryProvider);
  return OutingsNotifier(repository);
});

/// Liste filtrée selon le filtre actif
final filteredOutingsProvider = Provider<List<OutingModel>>((ref) {
  final state = ref.watch(outingsProvider);

  switch (state.currentFilter) {
    case OutingFilter.upcoming:
      return state.upcoming;
    case OutingFilter.past:
      return state.past;
    case OutingFilter.cancelled:
      return state.cancelled;
    case OutingFilter.all:
      return [...state.upcoming, ...state.past, ...state.cancelled];
  }
});

/// État du détail d'une sortie
@freezed
class OutingDetailState with _$OutingDetailState {
  const factory OutingDetailState.loading() = _Loading;
  const factory OutingDetailState.loaded(OutingModel outing) = _Loaded;
  const factory OutingDetailState.error(String message) = _Error;
}

/// Notifier pour le détail d'une sortie
class OutingDetailNotifier extends StateNotifier<OutingDetailState> {
  final OutingsRepository _repository;
  final String outingId;

  OutingDetailNotifier(this._repository, this.outingId)
      : super(const OutingDetailState.loading()) {
    _loadOuting();
  }

  Future<void> _loadOuting() async {
    try {
      final outing = await _repository.getOutingById(outingId);
      state = OutingDetailState.loaded(outing);
    } catch (e) {
      state = OutingDetailState.error(e.toString());
    }
  }

  Future<void> refresh() async {
    state = const OutingDetailState.loading();
    await _loadOuting();
  }

  Future<void> cancel({String? reason}) async {
    final currentState = state;
    if (currentState is! _Loaded) return;

    try {
      final updated = await _repository.cancelOuting(outingId, reason: reason);
      state = OutingDetailState.loaded(updated);
    } catch (e) {
      state = OutingDetailState.error(e.toString());
    }
  }
}

/// Provider famille pour le détail d'une sortie
final outingDetailProvider = StateNotifierProvider.family<OutingDetailNotifier,
    OutingDetailState, String>((ref, outingId) {
  final repository = ref.watch(outingsRepositoryProvider);
  return OutingDetailNotifier(repository, outingId);
});

/// État de réservation
@freezed
class BookingState with _$BookingState {
  const factory BookingState.initial() = _Initial;
  const factory BookingState.loading() = _BookingLoading;
  const factory BookingState.success(OutingModel outing) = _Success;
  const factory BookingState.error(String message) = _BookingError;
}

/// Notifier pour la réservation
class BookingNotifier extends StateNotifier<BookingState> {
  final OutingsRepository _repository;

  BookingNotifier(this._repository) : super(const BookingState.initial());

  /// Réserver une offre
  Future<void> bookOffer(String offerId) async {
    state = const BookingState.loading();

    try {
      final outing = await _repository.bookOffer(offerId);
      state = BookingState.success(outing);
    } catch (e) {
      state = BookingState.error(e.toString());
    }
  }

  /// Réinitialiser l'état
  void reset() {
    state = const BookingState.initial();
  }
}

/// Provider pour la réservation
final bookingProvider =
    StateNotifierProvider<BookingNotifier, BookingState>((ref) {
  final repository = ref.watch(outingsRepositoryProvider);
  return BookingNotifier(repository);
});

/// Provider pour vérifier si une offre peut être réservée
final canBookOfferProvider = FutureProvider.family<bool, String>((ref, offerId) async {
  final repository = ref.watch(outingsRepositoryProvider);
  return repository.canBookOffer(offerId);
});
