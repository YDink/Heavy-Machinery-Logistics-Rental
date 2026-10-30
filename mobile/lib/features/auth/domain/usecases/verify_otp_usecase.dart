import 'package:dink_mobile/features/auth/domain/entities/auth_tokens.dart';
import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';

class VerifyOtpUseCase {
  final AuthRepository repository;

  VerifyOtpUseCase(this.repository);

  Future<AuthTokens> call({
    required String verificationSessionId,
    required String code,
  }) {
    return repository.verifyOtp(
      verificationSessionId: verificationSessionId,
      code: code,
    );
  }
}
