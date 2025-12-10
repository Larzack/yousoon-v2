import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../models/offer_model.dart';

/// Interface du repository des offres
abstract class OffersRepository {
  /// Rechercher des offres
  Future<PaginatedOffersResult> searchOffers(OfferSearchParams params);

  /// Obtenir une offre par ID
  Future<OfferModel> getOfferById(String id);

  /// Obtenir les offres "Pour vous" (recommandations)
  Future<List<OfferModel>> getRecommendedOffers({
    required double latitude,
    required double longitude,
    int limit = 20,
  });

  /// Obtenir les offres à proximité
  Future<List<OfferModel>> getNearbyOffers({
    required double latitude,
    required double longitude,
    double radius = 10,
    int limit = 50,
  });

  /// Obtenir les offres par catégorie
  Future<PaginatedOffersResult> getOffersByCategory({
    required String categoryId,
    double? latitude,
    double? longitude,
    int page = 1,
    int limit = 20,
  });

  /// Obtenir les catégories
  Future<List<CategoryModel>> getCategories();
}

/// Implémentation du repository des offres
class OffersRepositoryImpl implements OffersRepository {
  // final Client _graphqlClient;

  OffersRepositoryImpl({
    // required Client graphqlClient,
  });
  // : _graphqlClient = graphqlClient;

  @override
  Future<PaginatedOffersResult> searchOffers(OfferSearchParams params) async {
    // TODO: Implémenter avec GraphQL
    // final result = await _graphqlClient.execute(
    //   SearchOffersQuery(
    //     variables: SearchOffersArguments(
    //       query: params.query,
    //       latitude: params.latitude,
    //       longitude: params.longitude,
    //       radius: params.radius,
    //       categoryIds: params.categoryIds,
    //       sortBy: params.sortBy.name,
    //       page: params.page,
    //       limit: params.limit,
    //     ),
    //   ),
    // );
    throw UnimplementedError('Search offers not implemented');
  }

  @override
  Future<OfferModel> getOfferById(String id) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get offer by id not implemented');
  }

  @override
  Future<List<OfferModel>> getRecommendedOffers({
    required double latitude,
    required double longitude,
    int limit = 20,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get recommended offers not implemented');
  }

  @override
  Future<List<OfferModel>> getNearbyOffers({
    required double latitude,
    required double longitude,
    double radius = 10,
    int limit = 50,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get nearby offers not implemented');
  }

  @override
  Future<PaginatedOffersResult> getOffersByCategory({
    required String categoryId,
    double? latitude,
    double? longitude,
    int page = 1,
    int limit = 20,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get offers by category not implemented');
  }

  @override
  Future<List<CategoryModel>> getCategories() async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get categories not implemented');
  }
}

/// Provider du repository des offres
final offersRepositoryProvider = Provider<OffersRepository>((ref) {
  return OffersRepositoryImpl();
});
