import 'package:ferry/ferry.dart';
import 'package:gql_http_link/gql_http_link.dart';
import 'package:gql_websocket_link/gql_websocket_link.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive/hive.dart';

import '../config/env_config.dart';

/// Provider pour le client GraphQL Ferry
final graphqlClientProvider = Provider<Client>((ref) {
  final httpLink = HttpLink(
    EnvConfig.apiBaseUrl,
    defaultHeaders: {
      'Content-Type': 'application/json',
    },
  );

  // TODO: Ajouter un AuthLink pour les tokens JWT
  // final authLink = AuthLink(getToken: () async {
  //   final token = await ref.read(authTokenProvider);
  //   return token != null ? 'Bearer $token' : null;
  // });

  final link = httpLink;

  final cache = Cache(
    store: HiveStore(),
  );

  return Client(
    link: link,
    cache: cache,
    defaultFetchPolicies: {
      OperationType.query: FetchPolicy.CacheFirst,
      OperationType.mutation: FetchPolicy.NetworkOnly,
    },
  );
});

/// Provider pour le client WebSocket (Subscriptions)
final graphqlWsClientProvider = Provider<WebSocketLink>((ref) {
  return WebSocketLink(
    EnvConfig.apiWsUrl,
    reconnectInterval: const Duration(seconds: 5),
  );
});

/// Store Hive pour le cache GraphQL
class HiveStore extends Store {
  final Box<Map<dynamic, dynamic>> _box;

  HiveStore._(this._box);

  static Future<HiveStore> open() async {
    final box = await Hive.openBox<Map<dynamic, dynamic>>('graphql_cache');
    return HiveStore._(box);
  }

  @override
  Map<String, dynamic>? get(String dataId) {
    final data = _box.get(dataId);
    if (data == null) return null;
    return Map<String, dynamic>.from(data);
  }

  @override
  void put(String dataId, Map<String, dynamic>? value) {
    if (value == null) {
      _box.delete(dataId);
    } else {
      _box.put(dataId, value);
    }
  }

  @override
  void putAll(Map<String, Map<String, dynamic>?> data) {
    for (final entry in data.entries) {
      put(entry.key, entry.value);
    }
  }

  @override
  void delete(String dataId) {
    _box.delete(dataId);
  }

  @override
  Map<String, Map<String, dynamic>?> toMap() {
    return _box.toMap().map(
          (key, value) => MapEntry(
            key.toString(),
            value != null ? Map<String, dynamic>.from(value) : null,
          ),
        );
  }

  @override
  void clear() {
    _box.clear();
  }
}
