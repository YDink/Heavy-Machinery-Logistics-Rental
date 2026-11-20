import 'package:dink_mobile/features/logistics/domain/entities/external_transport.dart';
import 'package:dink_mobile/features/logistics/domain/entities/logistics_dispatch_result.dart';
import 'package:dink_mobile/features/logistics/domain/entities/logistics_match_result.dart';
import 'package:dink_mobile/features/logistics/domain/entities/logistics_sub_booking_result.dart';
import 'package:dink_mobile/features/logistics/domain/entities/matched_carrier.dart';
import 'package:dink_mobile/features/logistics/domain/entities/reverse_cascade_result.dart';
import 'package:dink_mobile/features/logistics/domain/entities/transshipment_alert_result.dart';
import 'package:dink_mobile/features/logistics/domain/entities/winching_photo.dart';
import 'package:dink_mobile/features/logistics/domain/repositories/logistics_repository.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/accept_transshipment_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/advance_logistics_dispatch_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/create_logistics_sub_booking_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/match_logistics_carriers_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/register_external_transport_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/report_transshipment_alert_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/trigger_reverse_cascade_usecase.dart';
import 'package:dink_mobile/features/logistics/presentation/bloc/logistics_bloc.dart';
import 'package:dink_mobile/features/logistics/presentation/bloc/logistics_event.dart';
import 'package:dink_mobile/features/logistics/presentation/bloc/logistics_state.dart';
import 'package:flutter_test/flutter_test.dart';

class FakeLogisticsRepository implements LogisticsRepository {
  bool shouldFail = false;

  @override
  Future<LogisticsMatchResult> matchCarriers({
    required String bookingId,
    required int searchRadiusKm,
    required DateTime pickupDate,
  }) async {
    if (shouldFail) throw Exception('match error');
    return const LogisticsMatchResult(
      parentBookingId: 'bk-1',
      selfMoveCapable: false,
      recommendedTrailerClass: 'Two-Axle Lowboy Trailer (25t Payload Capacity)',
      abnormalLoadEscortRequired: false,
      operatingWeightKg: 21800,
      transportLengthMm: 9530,
      transportWidthMm: 2980,
      transportHeightMm: 3150,
      requiredMinimumPayloadKg: 23980,
      availableLogisticsFleets: [
        MatchedCarrier(
          fleetId: 'flt-1',
          make: 'Sinotruk',
          model: 'Howo Tri-Axle Lowbed',
          ownerTradingName: 'Girma Heavy Haulage PLC',
          payloadCapacityKg: 34000,
          deckWellLengthMm: 11000,
          deckWidthMm: 3100,
          tareWeightKg: 12500,
          axleCount: 3,
          dailyRateMinor: 1800000,
          currency: 'ETB',
          distanceFromDepotMeters: 8420,
          averageRating: 4.95,
        ),
      ],
    );
  }

  @override
  Future<LogisticsSubBookingResult> createSubBooking({
    required String bookingId,
    required String logisticsFleetId,
    required DateTime pickupDateTime,
    required double originLat,
    required double originLng,
    required String originName,
    required double destLat,
    required double destLng,
    required String destName,
  }) async {
    if (shouldFail) throw Exception('sub-booking failed');
    return const LogisticsSubBookingResult(
      subBookingId: 'sub-1',
      parentBookingId: 'bk-1',
      status: 'CONFIRMED',
      freightHireFeeMinor: 3600000,
      platformCommissionBps: 250,
      platformCommissionRate: '0.025',
      commissionOwedMinor: 90000,
      currency: 'ETB',
      abnormalLoadEscortRequired: false,
      escortFeeMinor: 0,
      carrierFullName: 'Ato Girma Wolde',
      carrierPhoneNumber: '+251911449120',
    );
  }

  @override
  Future<ExternalTransport> registerExternalTransport({
    required String bookingId,
    required String carrierName,
    required String driverFullName,
    required String driverPhoneNumber,
    required String primeMoverPlateNumber,
    required DateTime estimatedArrivalDateTime,
  }) async {
    return ExternalTransport(
      parentBookingId: bookingId,
      status: 'READY_FOR_PICKUP',
      carrierName: carrierName,
      driverFullName: driverFullName,
      driverPhoneNumber: driverPhoneNumber,
      primeMoverPlateNumber: primeMoverPlateNumber,
      estimatedArrivalDateTime: estimatedArrivalDateTime,
    );
  }

  @override
  Future<LogisticsDispatchResult> advanceDispatchStage({
    required String subBookingId,
    required String targetStage,
    List<WinchingPhoto> winchingPhotos = const [],
    String? notes,
  }) async {
    return LogisticsDispatchResult(
      subBookingId: subBookingId,
      parentBookingId: 'bk-1',
      subBookingStatus: targetStage,
      parentBookingStatus: 'CONFIRMED',
      winchingPhotosCount: winchingPhotos.length,
      updatedAt: DateTime.utc(2026, 11, 20),
    );
  }

  @override
  Future<TransshipmentAlertResult> reportTransshipmentAlert({
    required String subBookingId,
    required double latitude,
    required double longitude,
    required String failureCategory,
    required String failureDescription,
    required List<String> photos,
  }) async {
    return TransshipmentAlertResult(
      transshipmentId: 'tr-1',
      subBookingId: subBookingId,
      parentBookingId: 'bk-1',
      subBookingStatus: 'TRANSSHIPMENT_REQUIRED',
      parentBookingStatus: 'AWAITING_FLEET_READINESS',
      hotSwapSlaDeadline: DateTime.utc(2026, 11, 20, 12, 0),
      transshipmentBountyMinor: 720000,
      currency: 'ETB',
      candidateCarriersCount: 1,
    );
  }

  @override
  Future<String> acceptTransshipment({
    required String subBookingId,
    required String replacementFleetId,
  }) async {
    return 'IN_TRANSIT';
  }

  @override
  Future<ReverseCascadeResult> triggerReverseCascade({
    required String parentBookingId,
    String? reason,
  }) async {
    return ReverseCascadeResult(
      parentBookingId: parentBookingId,
      subBookingId: 'sub-1',
      resolution: 'FORCE_MAJEURE_CANCELLATION',
      renterRefundDueMinor: 8750000,
      ownerPenaltyLeviedMinor: 0,
      ownerCommissionRestitutionMinor: 1750000,
      carrierDefaultPenaltyMinor: 720000,
      currency: 'ETB',
      parentStatus: 'CANCELLED',
      subStatus: 'CANCELLED_CARRIER_DEFAULT',
    );
  }
}

void main() {
  late FakeLogisticsRepository repo;
  late LogisticsBloc bloc;

  setUp(() {
    repo = FakeLogisticsRepository();
    bloc = LogisticsBloc(
      matchCarriersUseCase: MatchLogisticsCarriersUseCase(repo),
      createSubBookingUseCase: CreateLogisticsSubBookingUseCase(repo),
      registerExternalTransportUseCase: RegisterExternalTransportUseCase(repo),
      advanceDispatchUseCase: AdvanceLogisticsDispatchUseCase(repo),
      reportTransshipmentAlertUseCase: ReportTransshipmentAlertUseCase(repo),
      acceptTransshipmentUseCase: AcceptTransshipmentUseCase(repo),
      triggerReverseCascadeUseCase: TriggerReverseCascadeUseCase(repo),
    );
  });

  tearDown(() => bloc.close());

  test(
    'emits LogisticsCarriersMatched when MatchLogisticsCarriersRequested succeeds',
    () async {
      final future = expectLater(
        bloc.stream,
        emitsInOrder([
          const LogisticsLoading(),
          isA<LogisticsCarriersMatched>().having(
            (s) => s.result.availableLogisticsFleets.length,
            'carriers count',
            1,
          ),
        ]),
      );

      bloc.add(MatchLogisticsCarriersRequested(
        bookingId: 'bk-1',
        searchRadiusKm: 75,
        pickupDate: DateTime.utc(2026, 11, 25),
      ));

      await future;
    },
  );

  test(
    'emits LogisticsSubBookingCreated with 2.5% commission rate',
    () async {
      final future = expectLater(
        bloc.stream,
        emitsInOrder([
          const LogisticsLoading(),
          isA<LogisticsSubBookingCreated>().having(
            (s) => s.subBooking.platformCommissionBps,
            'commission bps',
            250,
          ),
        ]),
      );

      bloc.add(CreateLogisticsSubBookingRequested(
        bookingId: 'bk-1',
        logisticsFleetId: 'flt-1',
        pickupDateTime: DateTime.utc(2026, 11, 25),
        originLat: 8.91,
        originLng: 38.74,
        originName: 'Kality',
        destLat: 8.54,
        destLng: 39.27,
        destName: 'Mojo',
      ));

      await future;
    },
  );

  test(
    'emits ReverseCascadeExecuted with 0 owner penalty',
    () async {
      final future = expectLater(
        bloc.stream,
        emitsInOrder([
          const LogisticsLoading(),
          isA<ReverseCascadeExecuted>().having(
            (s) => s.cascadeResult.ownerPenaltyLeviedMinor,
            'owner penalty',
            0,
          ),
        ]),
      );

      bloc.add(const TriggerReverseCascadeRequested(
        parentBookingId: 'bk-1',
        reason: 'SLA expired',
      ));

      await future;
    },
  );

  test(
    'emits LogisticsFailure when repository throws',
    () async {
      repo.shouldFail = true;
      final future = expectLater(
        bloc.stream,
        emitsInOrder([
          const LogisticsLoading(),
          isA<LogisticsFailure>(),
        ]),
      );

      bloc.add(MatchLogisticsCarriersRequested(
        bookingId: 'bk-1',
        searchRadiusKm: 75,
        pickupDate: DateTime.utc(2026, 11, 25),
      ));

      await future;
    },
  );
}
