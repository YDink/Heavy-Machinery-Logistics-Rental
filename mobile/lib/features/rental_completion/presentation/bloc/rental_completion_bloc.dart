import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/usecases/final_reconcile_usecase.dart';
import '../../domain/usecases/initiate_return_usecase.dart';
import '../../domain/usecases/inspect_return_usecase.dart';
import '../../domain/usecases/lodge_damage_claim_usecase.dart';
import '../../domain/usecases/respond_damage_claim_usecase.dart';
import '../../domain/usecases/submit_review_usecase.dart';
import 'rental_completion_event.dart';
import 'rental_completion_state.dart';

class RentalCompletionBloc extends Bloc<RentalCompletionEvent, RentalCompletionState> {
  final InitiateReturnUseCase initiateReturnUseCase;
  final InspectReturnUseCase inspectReturnUseCase;
  final LodgeDamageClaimUseCase lodgeDamageClaimUseCase;
  final RespondDamageClaimUseCase respondDamageClaimUseCase;
  final FinalReconcileUseCase finalReconcileUseCase;
  final SubmitReviewUseCase submitReviewUseCase;

  RentalCompletionBloc({
    required this.initiateReturnUseCase,
    required this.inspectReturnUseCase,
    required this.lodgeDamageClaimUseCase,
    required this.respondDamageClaimUseCase,
    required this.finalReconcileUseCase,
    required this.submitReviewUseCase,
  }) : super(const RentalCompletionInitial()) {
    on<InitiateReturnEvent>(_onInitiateReturn);
    on<SubmitReturnInspectionEvent>(_onSubmitReturnInspection);
    on<LodgeDamageClaimEvent>(_onLodgeDamageClaim);
    on<RespondDamageClaimEvent>(_onRespondDamageClaim);
    on<ReconcileFinalSettlementEvent>(_onReconcileFinalSettlement);
    on<SubmitBookingReviewEvent>(_onSubmitBookingReview);
  }

  Future<void> _onInitiateReturn(
    InitiateReturnEvent event,
    Emitter<RentalCompletionState> emit,
  ) async {
    emit(const RentalCompletionLoading());
    try {
      final result = await initiateReturnUseCase(event.bookingId);
      emit(ReturnInitiatedState(result));
    } catch (e) {
      emit(RentalCompletionError(e.toString()));
    }
  }

  Future<void> _onSubmitReturnInspection(
    SubmitReturnInspectionEvent event,
    Emitter<RentalCompletionState> emit,
  ) async {
    emit(const RentalCompletionLoading());
    try {
      final inspection = await inspectReturnUseCase(InspectReturnParams(
        bookingId: event.bookingId,
        checkOutMeterHours: event.checkOutMeterHours,
        fluidLevelsOk: event.fluidLevelsOk,
        safetyGlassIntact: event.safetyGlassIntact,
        accelerometerSteady: event.accelerometerSteady,
        notes: event.notes,
        photos: event.photos,
      ));
      emit(ReturnInspectionSuccessState(inspection));
    } catch (e) {
      emit(RentalCompletionError(e.toString()));
    }
  }

  Future<void> _onLodgeDamageClaim(
    LodgeDamageClaimEvent event,
    Emitter<RentalCompletionState> emit,
  ) async {
    emit(const RentalCompletionLoading());
    try {
      final claim = await lodgeDamageClaimUseCase(LodgeDamageClaimParams(
        bookingId: event.bookingId,
        componentCategory: event.componentCategory,
        incidentNarrative: event.incidentNarrative,
        repairQuoteUri: event.repairQuoteUri,
        claimedAmountMinor: event.claimedAmountMinor,
        currency: event.currency,
        damagePhotos: event.damagePhotos,
      ));
      emit(DamageClaimLodgedState(claim));
    } catch (e) {
      emit(RentalCompletionError(e.toString()));
    }
  }

  Future<void> _onRespondDamageClaim(
    RespondDamageClaimEvent event,
    Emitter<RentalCompletionState> emit,
  ) async {
    emit(const RentalCompletionLoading());
    try {
      final result = await respondDamageClaimUseCase(RespondDamageClaimParams(
        bookingId: event.bookingId,
        claimId: event.claimId,
        action: event.action,
        acceptedAmountMinor: event.acceptedAmountMinor,
        disputeReason: event.disputeReason,
      ));
      emit(DamageClaimRespondedState(result));
    } catch (e) {
      emit(RentalCompletionError(e.toString()));
    }
  }

  Future<void> _onReconcileFinalSettlement(
    ReconcileFinalSettlementEvent event,
    Emitter<RentalCompletionState> emit,
  ) async {
    emit(const RentalCompletionLoading());
    try {
      final settlement = await finalReconcileUseCase(event.bookingId);
      emit(SettlementReconciledState(settlement));
    } catch (e) {
      emit(RentalCompletionError(e.toString()));
    }
  }

  Future<void> _onSubmitBookingReview(
    SubmitBookingReviewEvent event,
    Emitter<RentalCompletionState> emit,
  ) async {
    emit(const RentalCompletionLoading());
    try {
      final review = await submitReviewUseCase(SubmitReviewParams(
        bookingId: event.bookingId,
        rating: event.rating,
        feedback: event.feedback,
        dimensionRatings: event.dimensionRatings,
      ));
      emit(BookingReviewSubmittedState(review));
    } catch (e) {
      emit(RentalCompletionError(e.toString()));
    }
  }
}
