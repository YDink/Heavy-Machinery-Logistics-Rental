import '../entities/active_rental_hud.dart';
import '../repositories/active_rental_repository.dart';

class GetActiveRentalHUDUseCase {
  final ActiveRentalRepository repository;

  GetActiveRentalHUDUseCase(this.repository);

  Future<ActiveRentalHUD> call(String bookingId) {
    return repository.getHUD(bookingId);
  }
}
