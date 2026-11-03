import '../entities/rental_extension.dart';
import '../repositories/active_rental_repository.dart';

class RequestExtensionUseCase {
  final ActiveRentalRepository repository;

  RequestExtensionUseCase(this.repository);

  Future<RentalExtension> call({
    required String bookingId,
    required DateTime newEndDate,
  }) {
    return repository.requestExtension(
      bookingId: bookingId,
      newEndDate: newEndDate,
    );
  }
}
