import '../entities/external_transport.dart';
import '../entities/logistics_dispatch_result.dart';
import '../entities/logistics_match_result.dart';
import '../entities/logistics_sub_booking_result.dart';
import '../entities/reverse_cascade_result.dart';
import '../entities/transshipment_alert_result.dart';
import '../entities/winching_photo.dart';

abstract class LogisticsRepository {
  Future<LogisticsMatchResult> matchCarriers({
    required String bookingId,
    required int searchRadiusKm,
    required DateTime pickupDate,
  });

  Future<LogisticsSubBookingResult> createSubBooking({
    required String bookingId,
    required String logisticsFleetId,
    required DateTime pickupDateTime,
    required double originLat,
    required double originLng,
    required String originName,
    required double destLat,
    required double destLng,
    required String destName,
  });

  Future<ExternalTransport> registerExternalTransport({
    required String bookingId,
    required String carrierName,
    required String driverFullName,
    required String driverPhoneNumber,
    required String primeMoverPlateNumber,
    required DateTime estimatedArrivalDateTime,
  });

  Future<LogisticsDispatchResult> advanceDispatchStage({
    required String subBookingId,
    required String targetStage,
    List<WinchingPhoto> winchingPhotos = const [],
    String? notes,
  });

  Future<TransshipmentAlertResult> reportTransshipmentAlert({
    required String subBookingId,
    required double latitude,
    required double longitude,
    required String failureCategory,
    required String failureDescription,
    required List<String> photos,
  });

  Future<String> acceptTransshipment({
    required String subBookingId,
    required String replacementFleetId,
  });

  Future<ReverseCascadeResult> triggerReverseCascade({
    required String parentBookingId,
    String? reason,
  });
}
