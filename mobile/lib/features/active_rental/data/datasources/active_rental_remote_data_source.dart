import 'package:intl/intl.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import '../models/handover_photo_item_model.dart';
import '../models/active_rental_hud_model.dart';
import '../models/handover_challenge_model.dart';
import '../models/handover_inspection_model.dart';
import '../models/mechanical_breakdown_model.dart';
import '../models/rental_extension_model.dart';

abstract class ActiveRentalRemoteDataSource {
  Future<ActiveRentalHUDModel> getHUD(String bookingId);

  Future<HandoverInspectionModel> inspectHandover({
    required String bookingId,
    required double checkInMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    String? notes,
    required List<HandoverPhotoItemModel> photos,
  });

  Future<HandoverChallengeModel> generateChallenge(String bookingId);

  Future<void> confirmHandover({
    required String bookingId,
    required String challengePin,
    String? touchSignatureBase64,
    double? latitude,
    double? longitude,
  });

  Future<RentalExtensionModel> requestExtension({
    required String bookingId,
    required DateTime newEndDate,
  });

  Future<MechanicalBreakdownModel> reportBreakdown({
    required String bookingId,
    required String faultCategory,
    required String faultDescription,
    List<String> diagnosticTroubleCodes = const [],
    required double currentHourMeter,
    List<String> incidentPhotos = const [],
  });

  Future<MechanicalBreakdownModel> resolveBreakdown({
    required String bookingId,
    required String mechanicName,
    required String mechanicPhone,
    required String resolutionNotes,
  });
}

class ActiveRentalRemoteDataSourceImpl implements ActiveRentalRemoteDataSource {
  final ApiClient apiClient;

  ActiveRentalRemoteDataSourceImpl({required this.apiClient});

  final DateFormat _dateFormat = DateFormat('yyyy-MM-dd');

  @override
  Future<ActiveRentalHUDModel> getHUD(String bookingId) async {
    final response = await apiClient.dio.get('/bookings/$bookingId/hud');
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return ActiveRentalHUDModel.fromJson(data);
  }

  @override
  Future<HandoverInspectionModel> inspectHandover({
    required String bookingId,
    required double checkInMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    String? notes,
    required List<HandoverPhotoItemModel> photos,
  }) async {
    final payload = <String, dynamic>{
      'checkInMeterHours': checkInMeterHours,
      'fluidLevelsOk': fluidLevelsOk,
      'safetyGlassIntact': safetyGlassIntact,
      'inspectionPhotos': photos
          .map((p) => {
                'angle': p.angle,
                'storageUri': p.storageUri,
              })
          .toList(),
    };
    if (notes != null && notes.isNotEmpty) {
      payload['notes'] = notes;
    }

    final response = await apiClient.dio.post(
      '/bookings/$bookingId/handover-inspect',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return HandoverInspectionModel.fromJson(data);
  }

  @override
  Future<HandoverChallengeModel> generateChallenge(String bookingId) async {
    final response = await apiClient.dio.post(
      '/bookings/$bookingId/handover-challenge',
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return HandoverChallengeModel.fromJson(data);
  }

  @override
  Future<void> confirmHandover({
    required String bookingId,
    required String challengePin,
    String? touchSignatureBase64,
    double? latitude,
    double? longitude,
  }) async {
    final payload = <String, dynamic>{
      'challengePin': challengePin,
    };
    if (touchSignatureBase64 != null) {
      payload['touchSignatureBase64'] = touchSignatureBase64;
    }
    if (latitude != null && longitude != null) {
      payload['contractorLocation'] = {
        'latitude': latitude,
        'longitude': longitude,
      };
    }

    await apiClient.dio.post(
      '/bookings/$bookingId/handover-confirm',
      data: payload,
    );
  }

  @override
  Future<RentalExtensionModel> requestExtension({
    required String bookingId,
    required DateTime newEndDate,
  }) async {
    final payload = <String, dynamic>{
      'newEndDate': _dateFormat.format(newEndDate),
    };

    final response = await apiClient.dio.post(
      '/bookings/$bookingId/extend',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return RentalExtensionModel.fromJson(data);
  }

  @override
  Future<MechanicalBreakdownModel> reportBreakdown({
    required String bookingId,
    required String faultCategory,
    required String faultDescription,
    List<String> diagnosticTroubleCodes = const [],
    required double currentHourMeter,
    List<String> incidentPhotos = const [],
  }) async {
    final payload = <String, dynamic>{
      'faultCategory': faultCategory,
      'faultDescription': faultDescription,
      'currentHourMeter': currentHourMeter,
    };
    if (diagnosticTroubleCodes.isNotEmpty) {
      payload['diagnosticTroubleCodes'] = diagnosticTroubleCodes;
    }
    if (incidentPhotos.isNotEmpty) {
      payload['incidentPhotos'] = incidentPhotos;
    }

    final response = await apiClient.dio.post(
      '/bookings/$bookingId/breakdown',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return MechanicalBreakdownModel.fromJson(data);
  }

  @override
  Future<MechanicalBreakdownModel> resolveBreakdown({
    required String bookingId,
    required String mechanicName,
    required String mechanicPhone,
    required String resolutionNotes,
  }) async {
    final payload = <String, dynamic>{
      'mechanicName': mechanicName,
      'mechanicPhone': mechanicPhone,
      'resolutionNotes': resolutionNotes,
    };

    final response = await apiClient.dio.post(
      '/bookings/$bookingId/breakdown/resolve',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return MechanicalBreakdownModel.fromJson(data);
  }
}
