import '../entities/booking.dart';
import '../repositories/booking_repository.dart';

class CreateBookingParams {
  final String fleetId;
  final DateTime startDate;
  final DateTime endDate;
  final String siteLocationAddress;
  final List<String> attachmentIds;

  const CreateBookingParams({
    required this.fleetId,
    required this.startDate,
    required this.endDate,
    required this.siteLocationAddress,
    this.attachmentIds = const [],
  });
}

class CreateBookingUseCase {
  final BookingRepository repository;

  CreateBookingUseCase(this.repository);

  Future<BookingEntity> call(CreateBookingParams params) {
    return repository.createBooking(
      fleetId: params.fleetId,
      startDate: params.startDate,
      endDate: params.endDate,
      siteLocationAddress: params.siteLocationAddress,
      attachmentIds: params.attachmentIds,
    );
  }
}
