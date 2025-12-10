import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../models/favorite_model.dart';

/// Interface du repository des favoris
abstract class FavoritesRepository {
  /// Obtenir les favoris de l'utilisateur
  Future<PaginatedFavoritesResult> getFavorites({
    int page = 1,
    int limit = 20,
  });

  /// Ajouter une offre aux favoris
  Future<FavoriteModel> addToFavorites(String offerId);

  /// Retirer une offre des favoris
  Future<void> removeFromFavorites(String offerId);

  /// Vérifier si une offre est en favori
  Future<bool> isFavorite(String offerId);

  /// Obtenir les IDs des offres en favori (pour vérification rapide)
  Future<Set<String>> getFavoriteOfferIds();
}

/// Implémentation du repository des favoris
class FavoritesRepositoryImpl implements FavoritesRepository {
  // final Client _graphqlClient;

  FavoritesRepositoryImpl({
    // required Client graphqlClient,
  });
  // : _graphqlClient = graphqlClient;

  @override
  Future<PaginatedFavoritesResult> getFavorites({
    int page = 1,
    int limit = 20,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get favorites not implemented');
  }

  @override
  Future<FavoriteModel> addToFavorites(String offerId) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Add to favorites not implemented');
  }

  @override
  Future<void> removeFromFavorites(String offerId) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Remove from favorites not implemented');
  }

  @override
  Future<bool> isFavorite(String offerId) async {
    // TODO: Implémenter avec GraphQL ou cache local
    throw UnimplementedError('Is favorite not implemented');
  }

  @override
  Future<Set<String>> getFavoriteOfferIds() async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get favorite offer ids not implemented');
  }
}

/// Provider du repository des favoris
final favoritesRepositoryProvider = Provider<FavoritesRepository>((ref) {
  return FavoritesRepositoryImpl();
});
