import '../entities/booking.dart';
import '../repositories/booking_repository.dart';

class GetBookingUseCase {
  final BookingRepository repository;

  GetBookingUseCase(this.repository);

  Future<BookingEntity> call(String id) {
    return repository.getBooking(id);
  }
}
