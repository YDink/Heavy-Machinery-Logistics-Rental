import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/active_rental/domain/entities/active_rental_hud.dart';
import 'package:dink_mobile/features/active_rental/domain/entities/handover_challenge.dart';
import 'package:dink_mobile/features/active_rental/domain/entities/handover_inspection.dart';
import 'package:dink_mobile/features/active_rental/domain/entities/handover_photo_item.dart';
import 'package:dink_mobile/features/active_rental/domain/entities/mechanical_breakdown.dart';
import 'package:dink_mobile/features/active_rental/domain/entities/rental_extension.dart';
import 'package:dink_mobile/features/active_rental/domain/repositories/active_rental_repository.dart';
import 'package:dink_mobile/features/active_rental/presentation/bloc/active_rental_bloc.dart';
import 'package:dink_mobile/features/active_rental/presentation/bloc/active_rental_event.dart';
import 'package:dink_mobile/features/active_rental/presentation/bloc/active_rental_state.dart';

class MockActiveRentalRepository implements ActiveRentalRepository {
  bool shouldFail = false;
  String errorMessage = 'Simulated rental failure';

  late ActiveRentalHUD dummyHUD;
  late HandoverChallenge dummyChallenge;
  late HandoverInspection dummyInspection;
  late RentalExtension dummyExtension;
  late MechanicalBreakdown dummyBreakdown;

  MockActiveRentalRepository() {
    dummyHUD = const ActiveRentalHUD(
      bookingId: 'bk-123',
      status: 'ACTIVE',
      fleetId: 'flt-001',
      fleetTitle: 'Caterpillar 320D Excavator',
      plateNumber: 'ET-3-12345',
      operator: HUDOperator(
        id: 'op-01',
        fullName: 'Kassahun Bekele',
        phoneNumber: '+251911223344',
        licenseClass: 'HEAVY_MACHINERY_A',
      ),
      currentRentalDay: 2,
      totalRentalDays: 5,
      handoverBaselineHours: 1250.0,
      accumulatedHoursToday: 4.5,
      allowedHoursPerDay: 8.0,
      excessHoursTotal: 0.0,
      shiftTimerPaused: false,
      canExtend: true,
      canReportBreakdown: true,
      canInitiateReturn: true,
    );

    dummyChallenge = HandoverChallenge(
      challengePin: '819204',
      bookingId: 'bk-123',
      expiresAt: DateTime.now().add(const Duration(minutes: 15)),
      maxAttempts: 3,
      attemptCount: 0,
    );

    dummyInspection = HandoverInspection(
      bookingId: 'bk-123',
      inspectedAt: DateTime.now(),
      meterVerified: true,
      photosIngestedCount: 2,
      status: 'INSPECTED',
    );

    dummyExtension = const RentalExtension(
      bookingId: 'bk-123',
      extensionDays: 3,
      incrementalPrepaymentRequiredMinor: 3000000,
      additionalGrossMinor: 6000000,
      currency: 'ETB',
      status: 'REQUESTED',
      ownerApprovalRequired: true,
    );

    dummyBreakdown = MechanicalBreakdown(
      incidentId: 'inc-001',
      bookingId: 'bk-123',
      status: 'ACTIVE_BREAKDOWN_PAUSED',
      ownerSlaDeadline: DateTime.now().add(const Duration(hours: 4)),
      shiftTimerPausedAt: DateTime.now(),
    );
  }

  @override
  Future<ActiveRentalHUD> getHUD(String bookingId) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyHUD;
  }

  @override
  Future<HandoverInspection> inspectHandover({
    required String bookingId,
    required double checkInMeterHours,
    required bool fluidLevelsOk,
    required bool safetyGlassIntact,
    String? notes,
    required List<HandoverPhotoItem> photos,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyInspection;
  }

  @override
  Future<HandoverChallenge> generateChallenge(String bookingId) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyChallenge;
  }

  @override
  Future<void> confirmHandover({
    required String bookingId,
    required String challengePin,
    String? touchSignatureBase64,
    double? latitude,
    double? longitude,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
  }

  @override
  Future<RentalExtension> requestExtension({
    required String bookingId,
    required DateTime newEndDate,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyExtension;
  }

  @override
  Future<MechanicalBreakdown> reportBreakdown({
    required String bookingId,
    required String faultCategory,
    required String faultDescription,
    List<String> diagnosticTroubleCodes = const [],
    required double currentHourMeter,
    List<String> incidentPhotos = const [],
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBreakdown;
  }

  @override
  Future<MechanicalBreakdown> resolveBreakdown({
    required String bookingId,
    required String mechanicName,
    required String mechanicPhone,
    required String resolutionNotes,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBreakdown;
  }
}

void main() {
  group('ActiveRentalBloc Unit Tests', () {
    late MockActiveRentalRepository repo;
    late ActiveRentalBloc bloc;

    setUp(() {
      repo = MockActiveRentalRepository();
      bloc = ActiveRentalBloc(repository: repo);
    });

    tearDown(() {
      bloc.close();
    });

    test('initial state is ActiveRentalInitial', () {
      expect(bloc.state, equals(ActiveRentalInitial()));
    });

    test('LoadActiveRentalHUDEvent emits [Loading, Loaded] on success', () async {
      final expected = [
        ActiveRentalLoading(),
        ActiveRentalHUDLoaded(repo.dummyHUD),
      ];

      expectLater(bloc.stream, emitsInOrder(expected));
      bloc.add(const LoadActiveRentalHUDEvent('bk-123'));
    });

    test('LoadActiveRentalHUDEvent emits [Loading, Error] on failure', () async {
      repo.shouldFail = true;
      repo.errorMessage = 'Network connection dropped';

      final expected = [
        ActiveRentalLoading(),
        const ActiveRentalError('Exception: Network connection dropped'),
      ];

      expectLater(bloc.stream, emitsInOrder(expected));
      bloc.add(const LoadActiveRentalHUDEvent('bk-123'));
    });

    test('GenerateHandoverChallengeEvent emits [Loading, HandoverChallengeGenerated]', () async {
      final expected = [
        ActiveRentalLoading(),
        HandoverChallengeGenerated(repo.dummyChallenge),
      ];

      expectLater(bloc.stream, emitsInOrder(expected));
      bloc.add(const GenerateHandoverChallengeEvent('bk-123'));
    });

    test('ConfirmHandoverEvent verifies PIN, emits [Success, Loaded]', () async {
      final expected = [
        ActiveRentalLoading(),
        const ActiveRentalActionSuccess('Rental activated successfully!'),
        ActiveRentalLoading(),
        ActiveRentalHUDLoaded(repo.dummyHUD),
      ];

      expectLater(bloc.stream, emitsInOrder(expected));
      bloc.add(const ConfirmHandoverEvent(
        bookingId: 'bk-123',
        challengePin: '819204',
        touchSignatureBase64: 'sig...',
      ));
    });

    test('ReportMechanicalBreakdownEvent reports incident, pauses shift, and reloads HUD', () async {
      final expected = [
        ActiveRentalLoading(),
        const ActiveRentalActionSuccess('Breakdown reported. Shift timer paused and 4h owner SLA activated.'),
        ActiveRentalLoading(),
        ActiveRentalHUDLoaded(repo.dummyHUD),
      ];

      expectLater(bloc.stream, emitsInOrder(expected));
      bloc.add(const ReportMechanicalBreakdownEvent(
        bookingId: 'bk-123',
        faultCategory: 'HYDRAULIC_CIRCUIT',
        faultDescription: 'Main ram pressure drop',
        currentHourMeter: 1255.0,
      ));
    });

    test('RequestRentalExtensionEvent submits extension, emits [Success, Loaded]', () async {
      final expected = [
        ActiveRentalLoading(),
        const ActiveRentalActionSuccess('Extension requested successfully (pending owner approval).'),
        ActiveRentalLoading(),
        ActiveRentalHUDLoaded(repo.dummyHUD),
      ];

      expectLater(bloc.stream, emitsInOrder(expected));
      bloc.add(RequestRentalExtensionEvent(
        bookingId: 'bk-123',
        newEndDate: DateTime.now().add(const Duration(days: 3)),
      ));
    });
  });
}
