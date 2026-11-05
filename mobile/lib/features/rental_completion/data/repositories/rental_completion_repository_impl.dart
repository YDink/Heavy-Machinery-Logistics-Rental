import '../../domain/entities/booking_review.dart';
import '../../domain/entities/damage_claim.dart';
import '../../domain/entities/damage_claim_response_result.dart';
import '../../domain/entities/rental_settlement.dart';
import '../../domain/entities/return_initiate_result.dart';
import '../../domain/entities/return_inspection.dart';
import '../../domain/entities/return_photo_item.dart';
import '../../domain/repositories/rental_completion_repository.dart';
import '../datasources/rental_completion_remote_data_source.dart';
import '../models/return_photo_item_model.dart';

class RentalCompletionRepositoryImpl implements RentalCompletionRepository {
  final RentalCompletionRemoteDataSource remoteDataSource;

  RentalCompletionRepositoryImpl({required this.remoteDataSource});

  @override
  Future<ReturnInitiateResult> initiateReturn(String bookingId) {
    return remoteDataSource.initiateReturn(bookingId);
  }

  @override
  Future<ReturnInspection> inspectReturn({
    required String bookingId,
    required double checkOutMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    required bool accelerometerSteady,
    String? notes,
    required List<ReturnPhotoItem> photos,
  }) {
    return remoteDataSource.inspectReturn(
      bookingId: bookingId,
      checkOutMeterHours: checkOutMeterHours,
      fluidLevelsOk: fluidLevelsOk,
      safetyGlassIntact: safetyGlassIntact,
      accelerometerSteady: accelerometerSteady,
      notes: notes,
      photos: photos.map(ReturnPhotoItemModel.fromEntity).toList(),
    );
  }

  @override
  Future<DamageClaim> lodgeDamageClaim({
    required String bookingId,
    required String componentCategory,
    required String incidentNarrative,
    String? repairQuoteUri,
    required int claimedAmountMinor,
    required String currency,
    List<String> damagePhotos = const [],
  }) {
    return remoteDataSource.lodgeDamageClaim(
      bookingId: bookingId,
      componentCategory: componentCategory,
      incidentNarrative: incidentNarrative,
      repairQuoteUri: repairQuoteUri,
      claimedAmountMinor: claimedAmountMinor,
      currency: currency,
      damagePhotos: damagePhotos,
    );
  }

  @override
  Future<DamageClaimResponseResult> respondDamageClaim({
    required String bookingId,
    required String claimId,
    required String action,
    int? acceptedAmountMinor,
    String? disputeReason,
  }) {
    return remoteDataSource.respondDamageClaim(
      bookingId: bookingId,
      claimId: claimId,
      action: action,
      acceptedAmountMinor: acceptedAmountMinor,
      disputeReason: disputeReason,
    );
  }

  @override
  Future<RentalSettlement> finalReconcile(String bookingId) {
    return remoteDataSource.finalReconcile(bookingId);
  }

  @override
  Future<BookingReview> submitReview({
    required String bookingId,
    required double rating,
    String? feedback,
    Map<String, double>? dimensionRatings,
  }) {
    return remoteDataSource.submitReview(
      bookingId: bookingId,
      rating: rating,
      feedback: feedback,
      dimensionRatings: dimensionRatings,
    );
  }
}
