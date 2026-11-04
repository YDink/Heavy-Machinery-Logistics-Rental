import '../../domain/entities/active_rental_hud.dart';
import '../../domain/entities/handover_challenge.dart';
import '../../domain/entities/handover_inspection.dart';
import '../../domain/entities/handover_photo_item.dart';
import '../../domain/entities/mechanical_breakdown.dart';
import '../../domain/entities/rental_extension.dart';
import '../../domain/repositories/active_rental_repository.dart';
import '../datasources/active_rental_remote_data_source.dart';
import '../models/handover_photo_item_model.dart';

class ActiveRentalRepositoryImpl implements ActiveRentalRepository {
  final ActiveRentalRemoteDataSource remoteDataSource;

  ActiveRentalRepositoryImpl({required this.remoteDataSource});

  @override
  Future<ActiveRentalHUD> getHUD(String bookingId) {
    return remoteDataSource.getHUD(bookingId);
  }

  @override
  Future<HandoverInspection> inspectHandover({
    required String bookingId,
    required double checkInMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    String? notes,
    required List<HandoverPhotoItem> photos,
  }) {
    return remoteDataSource.inspectHandover(
      bookingId: bookingId,
      checkInMeterHours: checkInMeterHours,
      fluidLevelsOk: fluidLevelsOk,
      safetyGlassIntact: safetyGlassIntact,
      notes: notes,
      photos: photos.map(HandoverPhotoItemModel.fromEntity).toList(),
    );
  }

  @override
  Future<HandoverChallenge> generateChallenge(String bookingId) {
    return remoteDataSource.generateChallenge(bookingId);
  }

  @override
  Future<void> confirmHandover({
    required String bookingId,
    required String challengePin,
    String? touchSignatureBase64,
    double? latitude,
    double? longitude,
  }) {
    return remoteDataSource.confirmHandover(
      bookingId: bookingId,
      challengePin: challengePin,
      touchSignatureBase64: touchSignatureBase64,
      latitude: latitude,
      longitude: longitude,
    );
  }

  @override
  Future<RentalExtension> requestExtension({
    required String bookingId,
    required DateTime newEndDate,
  }) {
    return remoteDataSource.requestExtension(
      bookingId: bookingId,
      newEndDate: newEndDate,
    );
  }

  @override
  Future<MechanicalBreakdown> reportBreakdown({
    required String bookingId,
    required String faultCategory,
    required String faultDescription,
    List<String> diagnosticTroubleCodes = const [],
    required double currentHourMeter,
    List<String> incidentPhotos = const [],
  }) {
    return remoteDataSource.reportBreakdown(
      bookingId: bookingId,
      faultCategory: faultCategory,
      faultDescription: faultDescription,
      diagnosticTroubleCodes: diagnosticTroubleCodes,
      currentHourMeter: currentHourMeter,
      incidentPhotos: incidentPhotos,
    );
  }

  @override
  Future<MechanicalBreakdown> resolveBreakdown({
    required String bookingId,
    required String mechanicName,
    required String mechanicPhone,
    required String resolutionNotes,
  }) {
    return remoteDataSource.resolveBreakdown(
      bookingId: bookingId,
      mechanicName: mechanicName,
      mechanicPhone: mechanicPhone,
      resolutionNotes: resolutionNotes,
    );
  }
}
