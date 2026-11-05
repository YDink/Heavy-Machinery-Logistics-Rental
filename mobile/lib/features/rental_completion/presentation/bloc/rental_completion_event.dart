import 'package:equatable/equatable.dart';
import '../../domain/entities/return_photo_item.dart';

abstract class RentalCompletionEvent extends Equatable {
  const RentalCompletionEvent();

  @override
  List<Object?> get props => [];
}

class InitiateReturnEvent extends RentalCompletionEvent {
  final String bookingId;

  const InitiateReturnEvent({required this.bookingId});

  @override
  List<Object?> get props => [bookingId];
}

class SubmitReturnInspectionEvent extends RentalCompletionEvent {
  final String bookingId;
  final double checkOutMeterHours;
  final bool fluidLevelsOk;
  final bool safetyGlassIntact;
  final bool accelerometerSteady;
  final String? notes;
  final List<ReturnPhotoItem> photos;

  const SubmitReturnInspectionEvent({
    required this.bookingId,
    required this.checkOutMeterHours,
    required this.fluidLevelsOk,
    required this.safetyGlassIntact,
    required this.accelerometerSteady,
    this.notes,
    required this.photos,
  });

  @override
  List<Object?> get props => [
        bookingId,
        checkOutMeterHours,
        fluidLevelsOk,
        safetyGlassIntact,
        accelerometerSteady,
        notes,
        photos,
      ];
}

class LodgeDamageClaimEvent extends RentalCompletionEvent {
  final String bookingId;
  final String componentCategory;
  final String incidentNarrative;
  final String? repairQuoteUri;
  final int claimedAmountMinor;
  final String currency;
  final List<String> damagePhotos;

  const LodgeDamageClaimEvent({
    required this.bookingId,
    required this.componentCategory,
    required this.incidentNarrative,
    this.repairQuoteUri,
    required this.claimedAmountMinor,
    required this.currency,
    this.damagePhotos = const [],
  });

  @override
  List<Object?> get props => [
        bookingId,
        componentCategory,
        incidentNarrative,
        repairQuoteUri,
        claimedAmountMinor,
        currency,
        damagePhotos,
      ];
}

class RespondDamageClaimEvent extends RentalCompletionEvent {
  final String bookingId;
  final String claimId;
  final String action;
  final int? acceptedAmountMinor;
  final String? disputeReason;

  const RespondDamageClaimEvent({
    required this.bookingId,
    required this.claimId,
    required this.action,
    this.acceptedAmountMinor,
    this.disputeReason,
  });

  @override
  List<Object?> get props => [
        bookingId,
        claimId,
        action,
        acceptedAmountMinor,
        disputeReason,
      ];
}

class ReconcileFinalSettlementEvent extends RentalCompletionEvent {
  final String bookingId;

  const ReconcileFinalSettlementEvent({required this.bookingId});

  @override
  List<Object?> get props => [bookingId];
}

class SubmitBookingReviewEvent extends RentalCompletionEvent {
  final String bookingId;
  final double rating;
  final String? feedback;
  final Map<String, double>? dimensionRatings;

  const SubmitBookingReviewEvent({
    required this.bookingId,
    required this.rating,
    this.feedback,
    this.dimensionRatings,
  });

  @override
  List<Object?> get props => [
        bookingId,
        rating,
        feedback,
        dimensionRatings,
      ];
}
