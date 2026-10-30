import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';

class RequestOtpUseCase {
  final AuthRepository repository;

  RequestOtpUseCase(this.repository);

  Future<String> call({
    required String phoneNumber,
    String purpose = 'LOGIN',
    String? countryCode,
  }) {
    return repository.requestOtp(
      phoneNumber: phoneNumber,
      purpose: purpose,
      countryCode: countryCode,
    );
  }
}
