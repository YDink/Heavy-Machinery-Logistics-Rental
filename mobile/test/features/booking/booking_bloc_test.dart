import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/booking/domain/entities/booking.dart';
import 'package:dink_mobile/features/booking/domain/repositories/booking_repository.dart';
import 'package:dink_mobile/features/booking/presentation/bloc/booking_bloc.dart';
import 'package:dink_mobile/features/booking/presentation/bloc/booking_event.dart';
import 'package:dink_mobile/features/booking/presentation/bloc/booking_state.dart';

class MockBookingRepository implements BookingRepository {
  bool shouldFail = false;
  String errorMessage = 'Simulated booking failure';

  late BookingEntity dummyBooking;
  late BankSlipEntity dummySlip;

  MockBookingRepository() {
    dummyBooking = BookingEntity(
      id: 'bk-12345',
      bookingReference: 'BK-20261103-0001',
      fleetId: 'flt-001',
      renterId: 'usr-renter',
      ownerId: 'usr-owner',
      status: BookingStatus.pendingPayment,
      startDate: DateTime.now().add(const Duration(days: 1)),
      endDate: DateTime.now().add(const Duration(days: 4)),
      durationDays: 3,
      dailyRateMinor: 2000000,
      baseHireFeeMinor: 6000000,
      attachmentsFeeMinor: 0,
      grossRentalAmountMinor: 6000000,
      prepaymentAmountMinor: 3000000,
      penaltyAmountMinor: 0,
      refundAmountMinor: 0,
      currency: 'ETB',
      siteLocationAddress: 'Addis Ababa Bole Sub-city Site A',
      expiresAt: DateTime.now().add(const Duration(minutes: 60)),
      attachments: const [],
      version: 1,
      createdAt: DateTime.now(),
    );

    dummySlip = BankSlipEntity(
      id: 'slip-001',
      bookingId: 'bk-12345',
      uploaderId: 'usr-renter',
      storageUri: 'https://cdn.dinkmachinery.com/slips/slip1.png',
      fileHash: 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
      amountMinor: 3000000,
      currency: 'ETB',
      status: 'PENDING_VERIFICATION',
      createdAt: DateTime.now(),
    );
  }

  @override
  Future<BookingEntity> createBooking({
    required String fleetId,
    required DateTime startDate,
    required DateTime endDate,
    required String siteLocationAddress,
    List<String> attachmentIds = const [],
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBooking;
  }

  @override
  Future<BookingEntity> getBooking(String id) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBooking;
  }

  @override
  Future<List<BookingEntity>> listBookings({String? status}) async {
    if (shouldFail) throw Exception(errorMessage);
    return [dummyBooking];
  }

  @override
  Future<BankSlipEntity> submitPaymentEvidence({
    required String bookingId,
    required String storageUri,
    required String fileHash,
    String? transactionReference,
    required int amountMinor,
    String currency = 'ETB',
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummySlip;
  }

  @override
  Future<BookingEntity> confirmPayment({
    required String bookingId,
    required int expectedVersion,
    String? paymentReference,
    String? notes,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBooking.copyWith(
      status: BookingStatus.confirmed,
      version: expectedVersion + 1,
    );
  }

  @override
  Future<BookingEntity> declinePayment({
    required String bookingId,
    required int expectedVersion,
    required String reasonCode,
    required String rejectionNotes,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBooking.copyWith(
      status: BookingStatus.paymentDeclined,
      rejectionReasonCode: reasonCode,
      rejectionNotes: rejectionNotes,
      version: expectedVersion + 1,
    );
  }

  @override
  Future<BookingEntity> disputeBooking({
    required String bookingId,
    required int expectedVersion,
    required String disputeReason,
    List<String> evidenceUris = const [],
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBooking.copyWith(
      status: BookingStatus.disputeOpened,
      version: expectedVersion + 1,
    );
  }

  @override
  Future<BookingEntity> cancelBooking({
    required String bookingId,
    required int expectedVersion,
    required String reason,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBooking.copyWith(
      status: BookingStatus.cancelled,
      cancellationReason: reason,
      version: expectedVersion + 1,
    );
  }

  @override
  Future<BookingEntity> handoverCheckIn({
    required String bookingId,
    required int expectedVersion,
    required String handoverPin,
    required double operatingHoursMeter,
    List<String> photoUris = const [],
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBooking.copyWith(
      status: BookingStatus.active,
      handoverPin: handoverPin,
      checkInMeterHours: operatingHoursMeter,
      version: expectedVersion + 1,
    );
  }

  @override
  Future<BookingEntity> returnCheckOut({
    required String bookingId,
    required int expectedVersion,
    required double operatingHoursMeter,
    String? damageNotes,
    List<String> photoUris = const [],
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return dummyBooking.copyWith(
      status: BookingStatus.returnInitiated,
      checkOutMeterHours: operatingHoursMeter,
      version: expectedVersion + 1,
    );
  }
}

void main() {
  group('BookingBloc Tests', () {
    late MockBookingRepository repository;
    late BookingBloc bloc;

    setUp(() {
      repository = MockBookingRepository();
      bloc = BookingBloc(repository: repository);
    });

    tearDown(() {
      bloc.close();
    });

    test('initial state is BookingInitial', () {
      expect(bloc.state, equals(const BookingInitial()));
    });

    test('CreateBookingRequested emits [BookingLoading, BookingLoaded] on success', () async {
      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingLoaded>().having((s) => s.booking.id, 'booking id', 'bk-12345'),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(CreateBookingRequested(
        fleetId: 'flt-001',
        startDate: DateTime.now().add(const Duration(days: 1)),
        endDate: DateTime.now().add(const Duration(days: 4)),
        siteLocationAddress: 'Addis Ababa Bole Sub-city Site A',
      ));
    });

    test('CreateBookingRequested emits [BookingLoading, BookingError] on failure', () async {
      repository.shouldFail = true;
      repository.errorMessage = 'Active lease hold collision';

      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingError>().having((s) => s.message, 'error message', contains('collision')),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(CreateBookingRequested(
        fleetId: 'flt-001',
        startDate: DateTime.now().add(const Duration(days: 1)),
        endDate: DateTime.now().add(const Duration(days: 4)),
        siteLocationAddress: 'Addis Ababa Bole Sub-city Site A',
      ));
    });

    test('LoadBookingDetails emits [BookingLoading, BookingLoaded] on success', () async {
      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingLoaded>().having((s) => s.booking.id, 'id', 'bk-12345'),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(const LoadBookingDetails('bk-12345'));
    });

    test('SubmitPaymentSlipRequested emits [BookingLoading, BookingActionSuccess, BookingLoaded]', () async {
      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingActionSuccess>().having((s) => s.message, 'success msg', contains('Payment evidence submitted')),
        isA<BookingLoaded>(),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(const SubmitPaymentSlipRequested(
        bookingId: 'bk-12345',
        storageUri: 'https://cdn.dinkmachinery.com/slips/slip1.png',
        fileHash: 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
        amountMinor: 3000000,
      ));
    });

    test('ConfirmPaymentRequested transitions to confirmed and reloads booking', () async {
      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingActionSuccess>(),
        isA<BookingLoaded>().having((s) => s.booking.status, 'status', BookingStatus.confirmed),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(const ConfirmPaymentRequested(
        bookingId: 'bk-12345',
        expectedVersion: 1,
      ));
    });

    test('DeclinePaymentRequested transitions to paymentDeclined', () async {
      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingActionSuccess>(),
        isA<BookingLoaded>().having((s) => s.booking.status, 'status', BookingStatus.paymentDeclined),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(const DeclinePaymentRequested(
        bookingId: 'bk-12345',
        expectedVersion: 1,
        reasonCode: 'ILLEGIBLE_SLIP',
        rejectionNotes: 'Blurry photo, transaction number unreadable',
      ));
    });

    test('DisputeBookingRequested transitions to disputeOpened', () async {
      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingActionSuccess>(),
        isA<BookingLoaded>().having((s) => s.booking.status, 'status', BookingStatus.disputeOpened),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(const DisputeBookingRequested(
        bookingId: 'bk-12345',
        expectedVersion: 1,
        disputeReason: 'Payment verified with CBE branch manager, funds cleared',
      ));
    });

    test('CancelBookingRequested transitions to cancelled', () async {
      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingActionSuccess>(),
        isA<BookingLoaded>().having((s) => s.booking.status, 'status', BookingStatus.cancelled),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(const CancelBookingRequested(
        bookingId: 'bk-12345',
        expectedVersion: 1,
        reason: 'Project delayed due to municipal permit',
      ));
    });

    test('HandoverCheckInRequested verifies PIN and updates engine meter', () async {
      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingActionSuccess>(),
        isA<BookingLoaded>().having((s) => s.booking.status, 'status', BookingStatus.active),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(const HandoverCheckInRequested(
        bookingId: 'bk-12345',
        expectedVersion: 1,
        handoverPin: '9482',
        operatingHoursMeter: 1250.5,
      ));
    });

    test('ReturnCheckOutRequested initiates return inspection', () async {
      final expectedStates = [
        isA<BookingLoading>(),
        isA<BookingActionSuccess>(),
        isA<BookingLoaded>().having((s) => s.booking.status, 'status', BookingStatus.returnInitiated),
      ];

      expectLater(bloc.stream, emitsInOrder(expectedStates));

      bloc.add(const ReturnCheckOutRequested(
        bookingId: 'bk-12345',
        expectedVersion: 1,
        operatingHoursMeter: 1274.5,
        damageNotes: 'None',
      ));
    });
  });
}
