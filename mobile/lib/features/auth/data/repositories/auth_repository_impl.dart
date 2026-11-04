import 'package:dink_mobile/core/storage/secure_storage_service.dart';
import 'package:dink_mobile/features/auth/data/datasources/auth_remote_datasource.dart';
import 'package:dink_mobile/features/auth/domain/entities/auth_tokens.dart';
import 'package:dink_mobile/features/auth/domain/entities/auth_user.dart';
import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource remoteDataSource;
  final SecureStorageService storageService;

  AuthRepositoryImpl({
    required this.remoteDataSource,
    required this.storageService,
  });

  @override
  Future<String> requestOtp({
    required String phoneNumber,
    required String purpose,
    String? countryCode,
  }) {
    return remoteDataSource.requestOtp(
      phoneNumber: phoneNumber,
      purpose: purpose,
      countryCode: countryCode,
    );
  }

  @override
  Future<AuthTokens> verifyOtp({
    required String verificationSessionId,
    required String code,
  }) async {
    final tokens = await remoteDataSource.verifyOtp(
      verificationSessionId: verificationSessionId,
      code: code,
    );

    await storageService.saveTokens(
      accessToken: tokens.accessToken,
      refreshToken: tokens.refreshToken,
    );

    return tokens;
  }

  @override
  Future<AuthTokens> refreshSession({required String refreshToken}) async {
    final tokens = await remoteDataSource.refreshSession(
      refreshToken: refreshToken,
    );

    await storageService.saveTokens(
      accessToken: tokens.accessToken,
      refreshToken: tokens.refreshToken,
    );

    return tokens;
  }

  @override
  Future<AuthUser> getCurrentUser() {
    return remoteDataSource.getCurrentUser();
  }

  @override
  Future<void> logout() {
    return storageService.clearTokens();
  }

  @override
  Future<bool> isAuthenticated() async {
    final token = await storageService.getAccessToken();
    return token != null && token.isNotEmpty;
  }
}
