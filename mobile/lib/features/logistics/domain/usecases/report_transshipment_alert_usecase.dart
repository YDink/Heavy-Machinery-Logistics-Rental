import '../entities/transshipment_alert_result.dart';
import '../repositories/logistics_repository.dart';

class ReportTransshipmentAlertUseCase {
  final LogisticsRepository repository;

  ReportTransshipmentAlertUseCase(this.repository);

  Future<TransshipmentAlertResult> call({
    required String subBookingId,
    required double latitude,
    required double longitude,
    required String failureCategory,
    required String failureDescription,
    required List<String> photos,
  }) {
    return repository.reportTransshipmentAlert(
      subBookingId: subBookingId,
      latitude: latitude,
      longitude: longitude,
      failureCategory: failureCategory,
      failureDescription: failureDescription,
      photos: photos,
    );
  }
}
