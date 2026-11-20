import 'package:dink_mobile/core/network/api_client.dart';
import '../models/external_transport_model.dart';
import '../models/logistics_dispatch_result_model.dart';
import '../models/logistics_match_result_model.dart';
import '../models/logistics_sub_booking_result_model.dart';
import '../models/reverse_cascade_result_model.dart';
import '../models/transshipment_alert_result_model.dart';
import '../models/winching_photo_model.dart';

abstract class LogisticsRemoteDataSource {
  Future<LogisticsMatchResultModel> matchCarriers({
    required String bookingId,
    required int searchRadiusKm,
    required DateTime pickupDate,
  });

  Future<LogisticsSubBookingResultModel> createSubBooking({
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

  Future<ExternalTransportModel> registerExternalTransport({
    required String bookingId,
    required String carrierName,
    required String driverFullName,
    required String driverPhoneNumber,
    required String primeMoverPlateNumber,
    required DateTime estimatedArrivalDateTime,
  });

  Future<LogisticsDispatchResultModel> advanceDispatchStage({
    required String subBookingId,
    required String targetStage,
    List<WinchingPhotoModel> winchingPhotos = const [],
    String? notes,
  });

  Future<TransshipmentAlertResultModel> reportTransshipmentAlert({
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

  Future<ReverseCascadeResultModel> triggerReverseCascade({
    required String parentBookingId,
    String? reason,
  });
}

class LogisticsRemoteDataSourceImpl implements LogisticsRemoteDataSource {
  final ApiClient apiClient;

  LogisticsRemoteDataSourceImpl({required this.apiClient});

  @override
  Future<LogisticsMatchResultModel> matchCarriers({
    required String bookingId,
    required int searchRadiusKm,
    required DateTime pickupDate,
  }) async {
    final response = await apiClient.dio.post(
      '/bookings/$bookingId/logistics/match',
      data: {
        'searchRadiusKm': searchRadiusKm,
        'pickupDate': pickupDate.toUtc().toIso8601String(),
      },
    );
    final body = response.data as Map<String, dynamic>;
    return LogisticsMatchResultModel.fromJson(body['data'] as Map<String, dynamic>);
  }

  @override
  Future<LogisticsSubBookingResultModel> createSubBooking({
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
    final response = await apiClient.dio.post(
      '/bookings/$bookingId/logistics/sub-booking',
      data: {
        'logisticsFleetId': logisticsFleetId,
        'pickupDateTime': pickupDateTime.toUtc().toIso8601String(),
        'originCoordinates': {
          'latitude': originLat,
          'longitude': originLng,
          'locationName': originName,
        },
        'destinationCoordinates': {
          'latitude': destLat,
          'longitude': destLng,
          'locationName': destName,
        },
      },
    );
    final body = response.data as Map<String, dynamic>;
    return LogisticsSubBookingResultModel.fromJson(body['data'] as Map<String, dynamic>);
  }

  @override
  Future<ExternalTransportModel> registerExternalTransport({
    required String bookingId,
    required String carrierName,
    required String driverFullName,
    required String driverPhoneNumber,
    required String primeMoverPlateNumber,
    required DateTime estimatedArrivalDateTime,
  }) async {
    final response = await apiClient.dio.post(
      '/bookings/$bookingId/logistics/external-transport',
      data: {
        'carrierName': carrierName,
        'driverFullName': driverFullName,
        'driverPhoneNumber': driverPhoneNumber,
        'primeMoverPlateNumber': primeMoverPlateNumber,
        'estimatedArrivalDateTime': estimatedArrivalDateTime.toUtc().toIso8601String(),
      },
    );
    final body = response.data as Map<String, dynamic>;
    return ExternalTransportModel.fromJson(body['data'] as Map<String, dynamic>);
  }

  @override
  Future<LogisticsDispatchResultModel> advanceDispatchStage({
    required String subBookingId,
    required String targetStage,
    List<WinchingPhotoModel> winchingPhotos = const [],
    String? notes,
  }) async {
    final payload = <String, dynamic>{
      'targetStage': targetStage,
      if (winchingPhotos.isNotEmpty)
        'winchingPhotos': winchingPhotos.map((p) => p.toJson()).toList(),
      if (notes != null && notes.isNotEmpty) 'notes': notes,
    };
    final response = await apiClient.dio.post(
      '/bookings/$subBookingId/logistics/dispatch',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    return LogisticsDispatchResultModel.fromJson(body['data'] as Map<String, dynamic>);
  }

  @override
  Future<TransshipmentAlertResultModel> reportTransshipmentAlert({
    required String subBookingId,
    required double latitude,
    required double longitude,
    required String failureCategory,
    required String failureDescription,
    required List<String> photos,
  }) async {
    final response = await apiClient.dio.post(
      '/bookings/$subBookingId/logistics/transshipment-alert',
      data: {
        'currentLocation': {
          'latitude': latitude,
          'longitude': longitude,
        },
        'failureCategory': failureCategory,
        'failureDescription': failureDescription,
        'photos': photos,
      },
    );
    final body = response.data as Map<String, dynamic>;
    return TransshipmentAlertResultModel.fromJson(body['data'] as Map<String, dynamic>);
  }

  @override
  Future<String> acceptTransshipment({
    required String subBookingId,
    required String replacementFleetId,
  }) async {
    final response = await apiClient.dio.post(
      '/bookings/$subBookingId/logistics/transshipment-accept',
      data: {
        'replacementFleetId': replacementFleetId,
      },
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return data['subBookingStatus'] as String? ?? 'IN_TRANSIT';
  }

  @override
  Future<ReverseCascadeResultModel> triggerReverseCascade({
    required String parentBookingId,
    String? reason,
  }) async {
    final response = await apiClient.dio.post(
      '/bookings/$parentBookingId/logistics/reverse-cascade',
      data: {
        if (reason != null && reason.isNotEmpty) 'reason': reason,
      },
    );
    final body = response.data as Map<String, dynamic>;
    return ReverseCascadeResultModel.fromJson(body['data'] as Map<String, dynamic>);
  }
}
