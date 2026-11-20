import '../entities/logistics_sub_booking_result.dart';
import '../repositories/logistics_repository.dart';

class CreateLogisticsSubBookingUseCase {
  final LogisticsRepository repository;

  CreateLogisticsSubBookingUseCase(this.repository);

  Future<LogisticsSubBookingResult> call({
    required String bookingId,
    required String logisticsFleetId,
    required DateTime pickupDateTime,
    required double originLat,
    required double originLng,
    required String originName,
    required double destLat,
    required double destLng,
    required String destName,
  }) {
    return repository.createSubBooking(
      bookingId: bookingId,
      logisticsFleetId: logisticsFleetId,
      pickupDateTime: pickupDateTime,
      originLat: originLat,
      originLng: originLng,
      originName: originName,
      destLat: destLat,
      destLng: destLng,
      destName: destName,
    );
  }
}
