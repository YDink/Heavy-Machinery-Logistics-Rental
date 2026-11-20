import '../repositories/logistics_repository.dart';

class AcceptTransshipmentUseCase {
  final LogisticsRepository repository;

  AcceptTransshipmentUseCase(this.repository);

  Future<String> call({
    required String subBookingId,
    required String replacementFleetId,
  }) {
    return repository.acceptTransshipment(
      subBookingId: subBookingId,
      replacementFleetId: replacementFleetId,
    );
  }
}
