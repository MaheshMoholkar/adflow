import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class AuthInterceptor extends Interceptor {
  AuthInterceptor(Dio dio);

  static const _storage = FlutterSecureStorage();
  static const _accessTokenKey = 'access_token';

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final publicPaths = [
      '/auth/register',
      '/auth/login',
      '/app/version',
      '/health',
    ];

    final isPublic = publicPaths.any((p) => options.path.contains(p));
    if (isPublic) {
      return handler.next(options);
    }

    final token = await _storage.read(key: _accessTokenKey);
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    return handler.next(options);
  }

  static Future<void> saveToken(String accessToken) async {
    await _storage.write(key: _accessTokenKey, value: accessToken);
  }

  static Future<void> clearTokens() async {
    await _storage.delete(key: _accessTokenKey);
  }

  static Future<String?> getAccessToken() async {
    return _storage.read(key: _accessTokenKey);
  }

  static Future<bool> hasTokens() async {
    final token = await _storage.read(key: _accessTokenKey);
    return token != null;
  }
}
