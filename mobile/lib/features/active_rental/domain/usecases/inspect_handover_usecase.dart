import '../entities/handover_inspection.dart';
import '../repositories/active_rental_repository.dart';

class InspectHandoverUseCase {
  final ActiveRentalRepository repository;

  InspectHandoverUseCase(this.repository);

  Future<HandoverInspection> call({
    required String bookingId,
    required double checkInMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    String? notes,
    required List<HandoverPhotoItem> photos,
  }) {
    return repository.inspectHandover(
      bookingId: bookingId,
      checkInMeterHours: checkInMeterHours,
      fluidLevelsOk: fluidLevelsOk,
      safetyGlassIntact: safetyGlassIntact,
      notes: notes,
      photos: photos,
    );
  }
}
