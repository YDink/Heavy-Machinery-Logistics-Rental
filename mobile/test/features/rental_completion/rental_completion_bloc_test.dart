import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/rental_completion/domain/entities/booking_review.dart';
import 'package:dink_mobile/features/rental_completion/domain/entities/damage_claim.dart';
import 'package:dink_mobile/features/rental_completion/domain/entities/damage_claim_response_result.dart';
import 'package:dink_mobile/features/rental_completion/domain/entities/rental_settlement.dart';
import 'package:dink_mobile/features/rental_completion/domain/entities/return_initiate_result.dart';
import 'package:dink_mobile/features/rental_completion/domain/entities/return_inspection.dart';
import 'package:dink_mobile/features/rental_completion/domain/entities/return_photo_item.dart';
import 'package:dink_mobile/features/rental_completion/domain/repositories/rental_completion_repository.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/final_reconcile_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/initiate_return_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/inspect_return_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/lodge_damage_claim_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/respond_damage_claim_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/submit_review_usecase.dart';
import 'package:dink_mobile/features/rental_completion/presentation/bloc/rental_completion_bloc.dart';
import 'package:dink_mobile/features/rental_completion/presentation/bloc/rental_completion_event.dart';
import 'package:dink_mobile/features/rental_completion/presentation/bloc/rental_completion_state.dart';

class MockRentalCompletionRepository implements RentalCompletionRepository {
  bool shouldFail = false;
  String errorMessage = 'Simulated completion failure';

  late ReturnInitiateResult dummyInitiate;
  late ReturnInspection dummyInspection;
  late DamageClaim dummyClaim;
  late DamageClaimResponseResult dummyClaimResponse;
  late RentalSettlement dummySettlement;
  late BookingReview dummyReview;

  MockRentalCompletionRepository() {
    final now = DateTime.now();

    dummyInitiate = ReturnInitiateResult(
      bookingId: 'bk-8812',
      status: 'RETURN_INITIATED',
      checkInMeterHours: 1250.0,
      allowedShiftHours: 24.0,
      initiatedAt: now,
    );

    dummyInspection = ReturnInspection(
      bookingId: 'bk-8812',
      checkOutMeterHours: 1314.80,
      damageAuditWindowExpiresAt: now.add(const Duration(hours: 24)),
      inspectedAt: now,
      photosIngestedCount: 6,
      status: 'DAMAGE_PENDING_REVIEW',
    );

    dummyClaim = DamageClaim(
      claimId: 'clm-001',
      bookingId: 'bk-8812',
      componentCategory: 'CABIN SAFETY GLASS & ROPS FRAME',
      claimedAmountMinor: 850000,
      acceptedAmountMinor: 0,
      currency: 'ETB',
      status: 'PENDING_REVIEW',
      createdAt: now,
    );

    dummyClaimResponse = DamageClaimResponseResult(
      claimId: 'clm-001',
      bookingId: 'bk-8812',
      acceptedAmountMinor: 850000,
      status: 'ACCEPTED',
      updatedAt: now,
    );

    dummySettlement = RentalSettlement(
      bookingId: 'bk-8812',
      baseMachineFeeMinor: 15750000,
      bundledAttachmentsFeeMinor: 5600000,
      contractualAllowedHours: 56.0,
      actualOperatingHours: 61.70,
      creditedBreakdownHours: 0.0,
      billableExcessHours: 5.70,
      excessHourlyRateMinor: 280000,
      excessHoursSurchargeMinor: 1596000,
      acceptedDamageMinor: 850000,
      grossLiabilitiesMinor: 23796000,
      prepaymentSettledMinor: 10675000,
      netBalanceMinor: 13121000,
      currency: 'ETB',
      status: 'PENDING_FINAL_PAYMENT',
      reconciledAt: now,
    );

    dummyReview = BookingReview(
      reviewId: 'rev-001',
      bookingId: 'bk-8812',
      authorId: 'user-renter',
      targetId: 'user-owner',
      role: 'RENTER',
      rating: 4.9,
      feedback: 'Excellent condition',
      isPublished: false,
      createdAt: now,
    );
  }

  @override
  Future<ReturnInitiateResult> initiateReturn(String bookingId) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyInitiate;
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
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyInspection;
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
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyClaim;
  }

  @override
  Future<DamageClaimResponseResult> respondDamageClaim({
    required String bookingId,
    required String claimId,
    required String action,
    int? acceptedAmountMinor,
    String? disputeReason,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyClaimResponse;
  }

  @override
  Future<RentalSettlement> finalReconcile(String bookingId) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummySettlement;
  }

  @override
  Future<BookingReview> submitReview({
    required String bookingId,
    required double rating,
    String? feedback,
    Map<String, double>? dimensionRatings,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyReview;
  }
}

void main() {
  late MockRentalCompletionRepository repository;
  late RentalCompletionBloc bloc;

  setUp(() {
    repository = MockRentalCompletionRepository();
    bloc = RentalCompletionBloc(
      initiateReturnUseCase: InitiateReturnUseCase(repository),
      inspectReturnUseCase: InspectReturnUseCase(repository),
      lodgeDamageClaimUseCase: LodgeDamageClaimUseCase(repository),
      respondDamageClaimUseCase: RespondDamageClaimUseCase(repository),
      finalReconcileUseCase: FinalReconcileUseCase(repository),
      submitReviewUseCase: SubmitReviewUseCase(repository),
    );
  });

  tearDown(() {
    bloc.close();
  });

  group('RentalCompletionBloc Tests', () {
    test('initial state is RentalCompletionInitial', () {
      expect(bloc.state, equals(const RentalCompletionInitial()));
    });

    test('emits [RentalCompletionLoading, ReturnInitiatedState] on InitiateReturnEvent', () async {
      final expectedStates = [
        const RentalCompletionLoading(),
        ReturnInitiatedState(repository.dummyInitiate),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));
      bloc.add(const InitiateReturnEvent(bookingId: 'bk-8812'));
    });

    test('emits [RentalCompletionLoading, ReturnInspectionSuccessState] on SubmitReturnInspectionEvent', () async {
      final expectedStates = [
        const RentalCompletionLoading(),
        ReturnInspectionSuccessState(repository.dummyInspection),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));
      bloc.add(const SubmitReturnInspectionEvent(
        bookingId: 'bk-8812',
        checkOutMeterHours: 1314.80,
        fluidLevelsOk: true,
        safetyGlassIntact: true,
        accelerometerSteady: true,
        photos: [],
      ));
    });

    test('emits [RentalCompletionLoading, DamageClaimLodgedState] on LodgeDamageClaimEvent', () async {
      final expectedStates = [
        const RentalCompletionLoading(),
        DamageClaimLodgedState(repository.dummyClaim),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));
      bloc.add(const LodgeDamageClaimEvent(
        bookingId: 'bk-8812',
        componentCategory: 'CABIN SAFETY GLASS & ROPS FRAME',
        incidentNarrative: 'Spiderweb crack from rock blast',
        claimedAmountMinor: 850000,
        currency: 'ETB',
      ));
    });

    test('emits [RentalCompletionLoading, DamageClaimRespondedState] on RespondDamageClaimEvent', () async {
      final expectedStates = [
        const RentalCompletionLoading(),
        DamageClaimRespondedState(repository.dummyClaimResponse),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));
      bloc.add(const RespondDamageClaimEvent(
        bookingId: 'bk-8812',
        claimId: 'clm-001',
        action: 'ACCEPT',
        acceptedAmountMinor: 850000,
      ));
    });

    test('emits [RentalCompletionLoading, SettlementReconciledState] on ReconcileFinalSettlementEvent', () async {
      final expectedStates = [
        const RentalCompletionLoading(),
        SettlementReconciledState(repository.dummySettlement),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));
      bloc.add(const ReconcileFinalSettlementEvent(bookingId: 'bk-8812'));
    });

    test('emits [RentalCompletionLoading, BookingReviewSubmittedState] on SubmitBookingReviewEvent', () async {
      final expectedStates = [
        const RentalCompletionLoading(),
        BookingReviewSubmittedState(repository.dummyReview),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));
      bloc.add(const SubmitBookingReviewEvent(
        bookingId: 'bk-8812',
        rating: 4.9,
        feedback: 'Excellent condition',
      ));
    });

    test('emits [RentalCompletionLoading, RentalCompletionError] on repository failure', () async {
      repository.shouldFail = true;
      final expectedStates = [
        const RentalCompletionLoading(),
        const RentalCompletionError('Exception: Simulated completion failure'),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));
      bloc.add(const InitiateReturnEvent(bookingId: 'bk-8812'));
    });
  });
}
