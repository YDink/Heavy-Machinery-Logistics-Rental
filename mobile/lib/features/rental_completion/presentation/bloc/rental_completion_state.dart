import 'package:equatable/equatable.dart';
import '../../domain/entities/booking_review.dart';
import '../../domain/entities/damage_claim.dart';
import '../../domain/entities/damage_claim_response_result.dart';
import '../../domain/entities/rental_settlement.dart';
import '../../domain/entities/return_initiate_result.dart';
import '../../domain/entities/return_inspection.dart';

abstract class RentalCompletionState extends Equatable {
  const RentalCompletionState();

  @override
  List<Object?> get props => [];
}

class RentalCompletionInitial extends RentalCompletionState {
  const RentalCompletionInitial();
}

class RentalCompletionLoading extends RentalCompletionState {
  const RentalCompletionLoading();
}

class ReturnInitiatedState extends RentalCompletionState {
  final ReturnInitiateResult result;

  const ReturnInitiatedState(this.result);

  @override
  List<Object?> get props => [result];
}

class ReturnInspectionSuccessState extends RentalCompletionState {
  final ReturnInspection inspection;

  const ReturnInspectionSuccessState(this.inspection);

  @override
  List<Object?> get props => [inspection];
}

class DamageClaimLodgedState extends RentalCompletionState {
  final DamageClaim claim;

  const DamageClaimLodgedState(this.claim);

  @override
  List<Object?> get props => [claim];
}

class DamageClaimRespondedState extends RentalCompletionState {
  final DamageClaimResponseResult result;

  const DamageClaimRespondedState(this.result);

  @override
  List<Object?> get props => [result];
}

class SettlementReconciledState extends RentalCompletionState {
  final RentalSettlement settlement;

  const SettlementReconciledState(this.settlement);

  @override
  List<Object?> get props => [settlement];
}

class BookingReviewSubmittedState extends RentalCompletionState {
  final BookingReview review;

  const BookingReviewSubmittedState(this.review);

  @override
  List<Object?> get props => [review];
}

class RentalCompletionError extends RentalCompletionState {
  final String message;

  const RentalCompletionError(this.message);

  @override
  List<Object?> get props => [message];
}
