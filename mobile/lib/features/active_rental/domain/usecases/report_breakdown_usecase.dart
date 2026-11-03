import '../entities/mechanical_breakdown.dart';
import '../repositories/active_rental_repository.dart';

class ReportBreakdownUseCase {
  final ActiveRentalRepository repository;

  ReportBreakdownUseCase(this.repository);

  Future<MechanicalBreakdown> call({
    required String bookingId,
    required String faultCategory,
    required String faultDescription,
    List<String> diagnosticTroubleCodes = const [],
    required double currentHourMeter,
    List<String> incidentPhotos = const [],
  }) {
    return repository.reportBreakdown(
      bookingId: bookingId,
      faultCategory: faultCategory,
      faultDescription: faultDescription,
      diagnosticTroubleCodes: diagnosticTroubleCodes,
      currentHourMeter: currentHourMeter,
      incidentPhotos: incidentPhotos,
    );
  }
}
