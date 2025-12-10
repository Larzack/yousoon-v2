import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../models/outing_model.dart';

/// Interface du repository des sorties
abstract class OutingsRepository {
  /// Réserver une offre
  Future<OutingModel> bookOffer(String offerId);

  /// Obtenir une sortie par ID
  Future<OutingModel> getOutingById(String id);

  /// Obtenir les sorties de l'utilisateur
  Future<PaginatedOutingsResult> getUserOutings(OutingsFilterParams params);

  /// Obtenir les sorties à venir
  Future<List<OutingModel>> getUpcomingOutings({int limit = 10});

  /// Obtenir l'historique des sorties
  Future<PaginatedOutingsResult> getOutingHistory({
    int page = 1,
    int limit = 20,
  });

  /// Annuler une sortie
  Future<OutingModel> cancelOuting(String id, {String? reason});

  /// Effectuer le check-in
  Future<OutingModel> checkIn(String qrCode);

  /// Vérifier si une offre peut être réservée
  Future<bool> canBookOffer(String offerId);
}

/// Implémentation du repository des sorties
class OutingsRepositoryImpl implements OutingsRepository {
  // final Client _graphqlClient;

  OutingsRepositoryImpl({
    // required Client graphqlClient,
  });
  // : _graphqlClient = graphqlClient;

  @override
  Future<OutingModel> bookOffer(String offerId) async {
    // TODO: Implémenter avec GraphQL
    // final result = await _graphqlClient.execute(
    //   BookOfferMutation(
    //     variables: BookOfferArguments(offerId: offerId),
    //   ),
    // );
    throw UnimplementedError('Book offer not implemented');
  }

  @override
  Future<OutingModel> getOutingById(String id) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get outing by id not implemented');
  }

  @override
  Future<PaginatedOutingsResult> getUserOutings(
      OutingsFilterParams params) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get user outings not implemented');
  }

  @override
  Future<List<OutingModel>> getUpcomingOutings({int limit = 10}) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get upcoming outings not implemented');
  }

  @override
  Future<PaginatedOutingsResult> getOutingHistory({
    int page = 1,
    int limit = 20,
  }) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Get outing history not implemented');
  }

  @override
  Future<OutingModel> cancelOuting(String id, {String? reason}) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Cancel outing not implemented');
  }

  @override
  Future<OutingModel> checkIn(String qrCode) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Check in not implemented');
  }

  @override
  Future<bool> canBookOffer(String offerId) async {
    // TODO: Implémenter avec GraphQL
    throw UnimplementedError('Can book offer not implemented');
  }
}

/// Provider du repository des sorties
final outingsRepositoryProvider = Provider<OutingsRepository>((ref) {
  return OutingsRepositoryImpl();
});
