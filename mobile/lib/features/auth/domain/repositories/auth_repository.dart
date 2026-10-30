import 'package:dink_mobile/features/auth/domain/entities/auth_tokens.dart';
import 'package:dink_mobile/features/auth/domain/entities/auth_user.dart';

abstract class AuthRepository {
  Future<String> requestOtp({required String phoneNumber, required String purpose, String? countryCode});
  Future<AuthTokens> verifyOtp({required String verificationSessionId, required String code});
  Future<AuthTokens> refreshSession({required String refreshToken});
  Future<AuthUser> getCurrentUser();
  Future<void> logout();
  Future<bool> isAuthenticated();
}
