import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
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
import 'package:dink_mobile/features/rental_completion/presentation/views/damage_docket_view.dart';
import 'package:dink_mobile/features/rental_completion/presentation/views/return_inspection_view.dart';
import 'package:dink_mobile/features/rental_completion/presentation/views/settlement_reconciliation_view.dart';
import 'package:dink_mobile/features/rental_completion/presentation/views/widgets/booking_review_modal.dart';

class StubRentalCompletionRepository implements RentalCompletionRepository {
  @override
  Future<ReturnInitiateResult> initiateReturn(String bookingId) async {
    return ReturnInitiateResult(
      bookingId: bookingId,
      status: 'RETURN_INITIATED',
      checkInMeterHours: 1250.0,
      allowedShiftHours: 24.0,
      initiatedAt: DateTime.now(),
    );
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
    return ReturnInspection(
      bookingId: bookingId,
      checkOutMeterHours: checkOutMeterHours,
      damageAuditWindowExpiresAt: DateTime.now().add(const Duration(hours: 24)),
      inspectedAt: DateTime.now(),
      photosIngestedCount: photos.length,
      status: 'DAMAGE_PENDING_REVIEW',
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
  }) async {
    return DamageClaim(
      claimId: 'clm-001',
      bookingId: bookingId,
      componentCategory: componentCategory,
      claimedAmountMinor: claimedAmountMinor,
      acceptedAmountMinor: 0,
      currency: currency,
      status: 'PENDING_REVIEW',
      createdAt: DateTime.now(),
    );
  }

  @override
  Future<DamageClaimResponseResult> respondDamageClaim({
    required String bookingId,
    required String claimId,
    required String action,
    int? acceptedAmountMinor,
    String? disputeReason,
  }) async {
    return DamageClaimResponseResult(
      claimId: claimId,
      bookingId: bookingId,
      acceptedAmountMinor: acceptedAmountMinor ?? 0,
      status: action == 'ACCEPT' ? 'ACCEPTED' : 'DISPUTED',
      updatedAt: DateTime.now(),
    );
  }

  @override
  Future<RentalSettlement> finalReconcile(String bookingId) async {
    return RentalSettlement(
      bookingId: bookingId,
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
      reconciledAt: DateTime.now(),
    );
  }

  @override
  Future<BookingReview> submitReview({
    required String bookingId,
    required double rating,
    String? feedback,
    Map<String, double>? dimensionRatings,
  }) async {
    return BookingReview(
      reviewId: 'rev-001',
      bookingId: bookingId,
      authorId: 'u-renter',
      targetId: 'u-owner',
      role: 'RENTER',
      rating: rating,
      feedback: feedback,
      isPublished: true,
      createdAt: DateTime.now(),
    );
  }
}

Widget createTestApp(Widget child, RentalCompletionBloc bloc) {
  return MaterialApp(
    home: BlocProvider<RentalCompletionBloc>.value(
      value: bloc,
      child: child,
    ),
  );
}

void main() {
  late StubRentalCompletionRepository repo;
  late RentalCompletionBloc bloc;

  setUp(() {
    repo = StubRentalCompletionRepository();
    bloc = RentalCompletionBloc(
      initiateReturnUseCase: InitiateReturnUseCase(repo),
      inspectReturnUseCase: InspectReturnUseCase(repo),
      lodgeDamageClaimUseCase: LodgeDamageClaimUseCase(repo),
      respondDamageClaimUseCase: RespondDamageClaimUseCase(repo),
      finalReconcileUseCase: FinalReconcileUseCase(repo),
      submitReviewUseCase: SubmitReviewUseCase(repo),
    );
  });

  tearDown(() {
    bloc.close();
  });

  group('Rental Completion Widget Tests', () {
    testWidgets('ReturnInspectionView renders hour meter audit and reticle angles', (tester) async {
      await tester.pumpWidget(createTestApp(
        const ReturnInspectionView(
          bookingId: 'bk-8812',
          checkInMeterHours: 1250.0,
        ),
        bloc,
      ));

      expect(find.text('RETURN INSPECTION'), findsOneWidget);
      expect(find.text('RC-02 AUDIT'), findsOneWidget);
      expect(find.text('ENGINE HOUR METER AUDIT'), findsOneWidget);
      expect(find.text('1250.00 hrs'), findsOneWidget);
      expect(find.text('MANDATORY 6-ANGLE RETICLE AUDIT'), findsOneWidget);
      expect(find.text('01. FRONT IMPLEMENT'), findsOneWidget);
      expect(find.text('02. REAR COUNTERWEIGHT'), findsOneWidget);
      expect(find.text('05. HOUR METER DIAL'), findsOneWidget);
      expect(find.text('SUBMIT RETURN INSPECTION →'), findsOneWidget);

      // Verify touch target height >= 56dp
      final submitButton = tester.widget<ElevatedButton>(find.widgetWithText(ElevatedButton, 'SUBMIT RETURN INSPECTION →'));
      expect(submitButton.onPressed, isNotNull);
      final size = tester.getSize(find.widgetWithText(ElevatedButton, 'SUBMIT RETURN INSPECTION →'));
      expect(size.height, greaterThanOrEqualTo(56.0));
    });

    testWidgets('SettlementReconciliationView renders itemized ledger and bank coordinates', (tester) async {
      await tester.pumpWidget(createTestApp(
        const SettlementReconciliationView(
          bookingId: 'bk-8812',
          bookingReference: 'DNK-8812',
        ),
        bloc,
      ));

      expect(find.text('SETTLEMENT STATEMENT'), findsOneWidget);
      expect(find.text('PENDING PAYMENT'), findsOneWidget);
      expect(find.text('ITEMIZED SETTLEMENT LEDGER'), findsOneWidget);
      expect(find.text('OWNER DIRECT WIRE COORDINATES'), findsOneWidget);
      expect(find.text('Commercial Bank of Ethiopia (CBE)'), findsOneWidget);
      expect(find.text('A/C: 1000289192831'), findsOneWidget);
      expect(find.text('COPY'), findsOneWidget);
      expect(find.text('CONFIRM & RECONCILE SETTLEMENT →'), findsOneWidget);

      final size = tester.getSize(find.widgetWithText(ElevatedButton, 'CONFIRM & RECONCILE SETTLEMENT →'));
      expect(size.height, greaterThanOrEqualTo(56.0));
    });

    testWidgets('DamageDocketView renders statutory countdown and resolution buttons', (tester) async {
      await tester.pumpWidget(createTestApp(
        const DamageDocketView(
          bookingId: 'bk-8812',
          claimId: 'clm-001',
          componentCategory: 'CABIN SAFETY GLASS & ROPS FRAME',
          claimedAmountMinor: 850000,
        ),
        bloc,
      ));

      expect(find.text('DAMAGE CLAIM DOCKET #DMG-03'), findsOneWidget);
      expect(find.text('18h 42m REMAINING BEFORE AUTO-WAIVER'), findsOneWidget);
      expect(find.text('CABIN SAFETY GLASS & ROPS FRAME'), findsOneWidget);
      expect(find.text('8500.00 ETB'), findsOneWidget);
      expect(find.text('RIES ENGINEERING SHARE CO.'), findsOneWidget);
      expect(find.text('ACCEPT LIABILITY (8500.00 ETB)'), findsOneWidget);
      expect(find.text('DISPUTE CLAIM & ESCALATE TO ARBITRATION'), findsOneWidget);

      final acceptBtnSize = tester.getSize(find.widgetWithText(ElevatedButton, 'ACCEPT LIABILITY (8500.00 ETB)'));
      expect(acceptBtnSize.height, greaterThanOrEqualTo(56.0));
    });

    testWidgets('BookingReviewModal renders double-blind guarantee banner and stars', (tester) async {
      await tester.pumpWidget(createTestApp(
        const BookingReviewModal(
          bookingId: 'bk-8812',
        ),
        bloc,
      ));

      expect(find.text('MUTUAL PERFORMANCE REVIEW'), findsOneWidget);
      expect(find.textContaining('DOUBLE-BLIND GUARANTEE'), findsOneWidget);
      expect(find.text('OVERALL STAR RATING'), findsOneWidget);
      expect(find.text('SUBMIT CONFIDENTIAL REVIEW →'), findsOneWidget);

      final submitBtnSize = tester.getSize(find.widgetWithText(ElevatedButton, 'SUBMIT CONFIDENTIAL REVIEW →'));
      expect(submitBtnSize.height, greaterThanOrEqualTo(56.0));
    });
  });
}
