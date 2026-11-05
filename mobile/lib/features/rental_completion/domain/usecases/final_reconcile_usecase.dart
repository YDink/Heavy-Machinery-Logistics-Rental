import '../entities/rental_settlement.dart';
import '../repositories/rental_completion_repository.dart';

class FinalReconcileUseCase {
  final RentalCompletionRepository repository;

  FinalReconcileUseCase(this.repository);

  Future<RentalSettlement> call(String bookingId) {
    return repository.finalReconcile(bookingId);
  }
}
