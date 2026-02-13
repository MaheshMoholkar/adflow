import 'package:drift/drift.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/database/app_database.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/auth_interceptor.dart';

final authServiceProvider = Provider<AuthService>((ref) {
  return AuthService(
    ref.watch(apiClientProvider),
    ref.watch(databaseProvider),
  );
});

class AuthService {
  final ApiClient _api;
  final AppDatabase _db;

  AuthService(this._api, this._db);

  Future<void> register({
    required String phone,
    required String password,
    required String name,
    required String businessName,
    required String city,
    required String address,
  }) async {
    final response = await _api.post(
      '/auth/register',
      data: {
        'phone': phone,
        'password': password,
        'name': name,
        'business_name': businessName,
        'city': city,
        'address': address,
      },
    );

    final data = response.data['data'] as Map<String, dynamic>;
    await _saveAuthResponse(data);
  }

  Future<void> login(String phone, String password) async {
    final response = await _api.post(
      '/auth/login',
      data: {'phone': phone, 'password': password},
    );

    final data = response.data['data'] as Map<String, dynamic>;
    await _saveAuthResponse(data);
  }

  Future<void> _saveAuthResponse(Map<String, dynamic> data) async {
    await AuthInterceptor.saveToken(data['access_token'] as String);

    final user = data['user'] as Map<String, dynamic>?;
    if (user != null) {
      await _db.upsertUser(UsersCompanion(
        id: Value(user['id'] as int),
        phone: Value(user['phone'] as String? ?? ''),
        name: Value(user['name'] as String? ?? ''),
        businessName: Value(user['business_name'] as String? ?? ''),
        city: Value(user['city'] as String? ?? ''),
        address: Value(user['address'] as String? ?? ''),
        locationUrl: Value(user['location_url'] as String? ?? ''),
        plan: Value(user['plan'] as String? ?? 'none'),
        planStartedAt: Value(user['plan_started_at'] != null
            ? DateTime.parse(user['plan_started_at'] as String)
            : null),
        planExpiresAt: Value(user['plan_expires_at'] != null
            ? DateTime.parse(user['plan_expires_at'] as String)
            : null),
        status: Value(user['status'] as String? ?? 'active'),
      ));
    }
  }

  Future<void> logout() async {
    await AuthInterceptor.clearTokens();
    await _db.clearAll();
  }
}
