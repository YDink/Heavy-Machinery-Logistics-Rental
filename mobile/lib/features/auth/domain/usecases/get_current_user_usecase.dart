import 'package:dink_mobile/features/auth/domain/entities/auth_user.dart';
import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';

class GetCurrentUserUseCase {
  final AuthRepository repository;

  GetCurrentUserUseCase(this.repository);

  Future<AuthUser> call() {
    return repository.getCurrentUser();
  }
}
