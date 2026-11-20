import '../../domain/entities/external_transport.dart';
import '../../domain/entities/logistics_dispatch_result.dart';
import '../../domain/entities/logistics_match_result.dart';
import '../../domain/entities/logistics_sub_booking_result.dart';
import '../../domain/entities/reverse_cascade_result.dart';
import '../../domain/entities/transshipment_alert_result.dart';
import '../../domain/entities/winching_photo.dart';
import '../../domain/repositories/logistics_repository.dart';
import '../datasources/logistics_remote_data_source.dart';
import '../models/winching_photo_model.dart';

class LogisticsRepositoryImpl implements LogisticsRepository {
  final LogisticsRemoteDataSource remoteDataSource;

  LogisticsRepositoryImpl({required this.remoteDataSource});

  @override
  Future<LogisticsMatchResult> matchCarriers({
    required String bookingId,
    required int searchRadiusKm,
    required DateTime pickupDate,
  }) async {
    final model = await remoteDataSource.matchCarriers(
      bookingId: bookingId,
      searchRadiusKm: searchRadiusKm,
      pickupDate: pickupDate,
    );
    return model.toEntity();
  }

  @override
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
  }) async {
    final model = await remoteDataSource.createSubBooking(
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
    return model.toEntity();
  }

  @override
  Future<ExternalTransport> registerExternalTransport({
    required String bookingId,
    required String carrierName,
    required String driverFullName,
    required String driverPhoneNumber,
    required String primeMoverPlateNumber,
    required DateTime estimatedArrivalDateTime,
  }) async {
    final model = await remoteDataSource.registerExternalTransport(
      bookingId: bookingId,
      carrierName: carrierName,
      driverFullName: driverFullName,
      driverPhoneNumber: driverPhoneNumber,
      primeMoverPlateNumber: primeMoverPlateNumber,
      estimatedArrivalDateTime: estimatedArrivalDateTime,
    );
    return model.toEntity();
  }

  @override
  Future<LogisticsDispatchResult> advanceDispatchStage({
    required String subBookingId,
    required String targetStage,
    List<WinchingPhoto> winchingPhotos = const [],
    String? notes,
  }) async {
    final photoModels = winchingPhotos.map(WinchingPhotoModel.fromEntity).toList();
    final model = await remoteDataSource.advanceDispatchStage(
      subBookingId: subBookingId,
      targetStage: targetStage,
      winchingPhotos: photoModels,
      notes: notes,
    );
    return model.toEntity();
  }

  @override
  Future<TransshipmentAlertResult> reportTransshipmentAlert({
    required String subBookingId,
    required double latitude,
    required double longitude,
    required String failureCategory,
    required String failureDescription,
    required List<String> photos,
  }) async {
    final model = await remoteDataSource.reportTransshipmentAlert(
      subBookingId: subBookingId,
      latitude: latitude,
      longitude: longitude,
      failureCategory: failureCategory,
      failureDescription: failureDescription,
      photos: photos,
    );
    return model.toEntity();
  }

  @override
  Future<String> acceptTransshipment({
    required String subBookingId,
    required String replacementFleetId,
  }) {
    return remoteDataSource.acceptTransshipment(
      subBookingId: subBookingId,
      replacementFleetId: replacementFleetId,
    );
  }

  @override
  Future<ReverseCascadeResult> triggerReverseCascade({
    required String parentBookingId,
    String? reason,
  }) async {
    final model = await remoteDataSource.triggerReverseCascade(
      parentBookingId: parentBookingId,
      reason: reason,
    );
    return model.toEntity();
  }
}
