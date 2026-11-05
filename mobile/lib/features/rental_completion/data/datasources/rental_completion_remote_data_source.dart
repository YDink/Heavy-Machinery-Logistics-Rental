import 'package:dink_mobile/core/network/api_client.dart';
import '../models/booking_review_model.dart';
import '../models/damage_claim_model.dart';
import '../models/damage_claim_response_result_model.dart';
import '../models/rental_settlement_model.dart';
import '../models/return_initiate_result_model.dart';
import '../models/return_inspection_model.dart';
import '../models/return_photo_item_model.dart';

abstract class RentalCompletionRemoteDataSource {
  Future<ReturnInitiateResultModel> initiateReturn(String bookingId);

  Future<ReturnInspectionModel> inspectReturn({
    required String bookingId,
    required double checkOutMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    required bool accelerometerSteady,
    String? notes,
    required List<ReturnPhotoItemModel> photos,
  });

  Future<DamageClaimModel> lodgeDamageClaim({
    required String bookingId,
    required String componentCategory,
    required String incidentNarrative,
    String? repairQuoteUri,
    required int claimedAmountMinor,
    required String currency,
    List<String> damagePhotos = const [],
  });

  Future<DamageClaimResponseResultModel> respondDamageClaim({
    required String bookingId,
    required String claimId,
    required String action,
    int? acceptedAmountMinor,
    String? disputeReason,
  });

  Future<RentalSettlementModel> finalReconcile(String bookingId);

  Future<BookingReviewModel> submitReview({
    required String bookingId,
    required double rating,
    String? feedback,
    Map<String, double>? dimensionRatings,
  });
}

class RentalCompletionRemoteDataSourceImpl implements RentalCompletionRemoteDataSource {
  final ApiClient apiClient;

  RentalCompletionRemoteDataSourceImpl({required this.apiClient});

  @override
  Future<ReturnInitiateResultModel> initiateReturn(String bookingId) async {
    final response = await apiClient.dio.post('/bookings/$bookingId/return-initiate');
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return ReturnInitiateResultModel.fromJson(data);
  }

  @override
  Future<ReturnInspectionModel> inspectReturn({
    required String bookingId,
    required double checkOutMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    required bool accelerometerSteady,
    String? notes,
    required List<ReturnPhotoItemModel> photos,
  }) async {
    final payload = <String, dynamic>{
      'checkOutMeterHours': checkOutMeterHours,
      'fluidLevelsOk': fluidLevelsOk,
      'safetyGlassIntact': safetyGlassIntact,
      'accelerometerSteady': accelerometerSteady,
      'inspectionPhotos': photos.map((p) => p.toJson()).toList(),
    };
    if (notes != null && notes.isNotEmpty) {
      payload['notes'] = notes;
    }

    final response = await apiClient.dio.post(
      '/bookings/$bookingId/return-inspect',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return ReturnInspectionModel.fromJson(data);
  }

  @override
  Future<DamageClaimModel> lodgeDamageClaim({
    required String bookingId,
    required String componentCategory,
    required String incidentNarrative,
    String? repairQuoteUri,
    required int claimedAmountMinor,
    required String currency,
    List<String> damagePhotos = const [],
  }) async {
    final payload = <String, dynamic>{
      'componentCategory': componentCategory,
      'incidentNarrative': incidentNarrative,
      'claimedAmountMinor': claimedAmountMinor,
      'currency': currency,
      'damagePhotos': damagePhotos,
    };
    if (repairQuoteUri != null && repairQuoteUri.isNotEmpty) {
      payload['repairQuoteUri'] = repairQuoteUri;
    }

    final response = await apiClient.dio.post(
      '/bookings/$bookingId/damage-claims',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return DamageClaimModel.fromJson(data);
  }

  @override
  Future<DamageClaimResponseResultModel> respondDamageClaim({
    required String bookingId,
    required String claimId,
    required String action,
    int? acceptedAmountMinor,
    String? disputeReason,
  }) async {
    final payload = <String, dynamic>{
      'action': action,
    };
    if (acceptedAmountMinor != null) {
      payload['acceptedAmountMinor'] = acceptedAmountMinor;
    }
    if (disputeReason != null && disputeReason.isNotEmpty) {
      payload['disputeReason'] = disputeReason;
    }

    final response = await apiClient.dio.post(
      '/bookings/$bookingId/damage-claims/$claimId/respond',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return DamageClaimResponseResultModel.fromJson(data);
  }

  @override
  Future<RentalSettlementModel> finalReconcile(String bookingId) async {
    final response = await apiClient.dio.post('/bookings/$bookingId/final-reconcile');
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return RentalSettlementModel.fromJson(data);
  }

  @override
  Future<BookingReviewModel> submitReview({
    required String bookingId,
    required double rating,
    String? feedback,
    Map<String, double>? dimensionRatings,
  }) async {
    final payload = <String, dynamic>{
      'rating': rating,
    };
    if (feedback != null && feedback.isNotEmpty) {
      payload['feedback'] = feedback;
    }
    if (dimensionRatings != null && dimensionRatings.isNotEmpty) {
      payload['dimensionRatings'] = dimensionRatings;
    }

    final response = await apiClient.dio.post(
      '/bookings/$bookingId/reviews',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BookingReviewModel.fromJson(data);
  }
}
