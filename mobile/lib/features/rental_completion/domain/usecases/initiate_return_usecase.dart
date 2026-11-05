import '../entities/return_initiate_result.dart';
import '../repositories/rental_completion_repository.dart';

class InitiateReturnUseCase {
  final RentalCompletionRepository repository;

  InitiateReturnUseCase(this.repository);

  Future<ReturnInitiateResult> call(String bookingId) {
    return repository.initiateReturn(bookingId);
  }
}
