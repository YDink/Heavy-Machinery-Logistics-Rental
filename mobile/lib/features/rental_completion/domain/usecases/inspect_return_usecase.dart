import '../entities/return_inspection.dart';
import '../entities/return_photo_item.dart';
import '../repositories/rental_completion_repository.dart';

class InspectReturnParams {
  final String bookingId;
  final double checkOutMeterHours;
  final bool fluidLevelsOk;
  final bool safetyGlassIntact;
  final bool accelerometerSteady;
  final String? notes;
  final List<ReturnPhotoItem> photos;

  const InspectReturnParams({
    required this.bookingId,
    required this.checkOutMeterHours,
    required this.fluidLevelsOk,
    required this.safetyGlassIntact,
    required this.accelerometerSteady,
    this.notes,
    required this.photos,
  });
}

class InspectReturnUseCase {
  final RentalCompletionRepository repository;

  InspectReturnUseCase(this.repository);

  Future<ReturnInspection> call(InspectReturnParams params) {
    return repository.inspectReturn(
      bookingId: params.bookingId,
      checkOutMeterHours: params.checkOutMeterHours,
      fluidLevelsOk: params.fluidLevelsOk,
      safetyGlassIntact: params.safetyGlassIntact,
      accelerometerSteady: params.accelerometerSteady,
      notes: params.notes,
      photos: params.photos,
    );
  }
}
