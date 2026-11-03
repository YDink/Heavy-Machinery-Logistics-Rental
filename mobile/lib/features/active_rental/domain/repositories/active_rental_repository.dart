import '../entities/active_rental_hud.dart';
import '../entities/handover_challenge.dart';
import '../entities/handover_inspection.dart';
import '../entities/mechanical_breakdown.dart';
import '../entities/rental_extension.dart';

abstract class ActiveRentalRepository {
  Future<ActiveRentalHUD> getHUD(String bookingId);

  Future<HandoverInspection> inspectHandover({
    required String bookingId,
    required double checkInMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    String? notes,
    required List<HandoverPhotoItem> photos,
  });

  Future<HandoverChallenge> generateChallenge(String bookingId);

  Future<void> confirmHandover({
    required String bookingId,
    required String challengePin,
    String? touchSignatureBase64,
    double? latitude,
    double? longitude,
  });

  Future<RentalExtension> requestExtension({
    required String bookingId,
    required DateTime newEndDate,
  });

  Future<MechanicalBreakdown> reportBreakdown({
    required String bookingId,
    required String faultCategory,
    required String faultDescription,
    List<String> diagnosticTroubleCodes = const [],
    required double currentHourMeter,
    List<String> incidentPhotos = const [],
  });

  Future<MechanicalBreakdown> resolveBreakdown({
    required String bookingId,
    required String mechanicName,
    required String mechanicPhone,
    required String resolutionNotes,
  });
}
