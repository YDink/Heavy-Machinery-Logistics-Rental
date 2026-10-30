import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class SecureStorageService {
  final FlutterSecureStorage _storage;

  SecureStorageService({FlutterSecureStorage? storage})
      : _storage = storage ?? const FlutterSecureStorage();

  static const String _accessTokenKey = 'dink_access_token';
  static const String _refreshTokenKey = 'dink_refresh_token';
  static const String _activeCountryKey = 'dink_active_country';

  Future<void> saveTokens({required String accessToken, required String refreshToken}) async {
    await _storage.write(key: _accessTokenKey, value: accessToken);
    await _storage.write(key: _refreshTokenKey, value: refreshToken);
  }

  Future<String?> getAccessToken() async {
    return _storage.read(key: _accessTokenKey);
  }

  Future<String?> getRefreshToken() async {
    return _storage.read(key: _refreshTokenKey);
  }

  Future<void> clearTokens() async {
    await _storage.delete(key: _accessTokenKey);
    await _storage.delete(key: _refreshTokenKey);
  }

  Future<void> saveCountryCode(String countryCode) async {
    await _storage.write(key: _activeCountryKey, value: countryCode);
  }

  Future<String> getCountryCode() async {
    final code = await _storage.read(key: _activeCountryKey);
    return code ?? 'ET';
  }
}
