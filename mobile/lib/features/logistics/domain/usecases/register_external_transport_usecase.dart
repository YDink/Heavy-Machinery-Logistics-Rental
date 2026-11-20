import '../entities/external_transport.dart';
import '../repositories/logistics_repository.dart';

class RegisterExternalTransportUseCase {
  final LogisticsRepository repository;

  RegisterExternalTransportUseCase(this.repository);

  Future<ExternalTransport> call({
    required String bookingId,
    required String carrierName,
    required String driverFullName,
    required String driverPhoneNumber,
    required String primeMoverPlateNumber,
    required DateTime estimatedArrivalDateTime,
  }) {
    return repository.registerExternalTransport(
      bookingId: bookingId,
      carrierName: carrierName,
      driverFullName: driverFullName,
      driverPhoneNumber: driverPhoneNumber,
      primeMoverPlateNumber: primeMoverPlateNumber,
      estimatedArrivalDateTime: estimatedArrivalDateTime,
    );
  }
}
