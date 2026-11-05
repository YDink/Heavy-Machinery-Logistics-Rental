import '../entities/booking_review.dart';
import '../entities/damage_claim.dart';
import '../entities/damage_claim_response_result.dart';
import '../entities/rental_settlement.dart';
import '../entities/return_initiate_result.dart';
import '../entities/return_inspection.dart';
import '../entities/return_photo_item.dart';

abstract class RentalCompletionRepository {
  Future<ReturnInitiateResult> initiateReturn(String bookingId);

  Future<ReturnInspection> inspectReturn({
    required String bookingId,
    required double checkOutMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    required bool accelerometerSteady,
    String? notes,
    required List<ReturnPhotoItem> photos,
  });

  Future<DamageClaim> lodgeDamageClaim({
    required String bookingId,
    required String componentCategory,
    required String incidentNarrative,
    String? repairQuoteUri,
    required int claimedAmountMinor,
    required String currency,
    List<String> damagePhotos = const [],
  });

  Future<DamageClaimResponseResult> respondDamageClaim({
    required String bookingId,
    required String claimId,
    required String action,
    int? acceptedAmountMinor,
    String? disputeReason,
  });

  Future<RentalSettlement> finalReconcile(String bookingId);

  Future<BookingReview> submitReview({
    required String bookingId,
    required double rating,
    String? feedback,
    Map<String, double>? dimensionRatings,
  });
}
