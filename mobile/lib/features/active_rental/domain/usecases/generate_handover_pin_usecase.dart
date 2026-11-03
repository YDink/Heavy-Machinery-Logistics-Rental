import '../entities/handover_challenge.dart';
import '../repositories/active_rental_repository.dart';

class GenerateHandoverPINUseCase {
  final ActiveRentalRepository repository;

  GenerateHandoverPINUseCase(this.repository);

  Future<HandoverChallenge> call(String bookingId) {
    return repository.generateChallenge(bookingId);
  }
}
