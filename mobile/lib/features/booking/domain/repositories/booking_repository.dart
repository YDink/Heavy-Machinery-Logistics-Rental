import '../entities/booking.dart';

abstract class BookingRepository {
  Future<BookingEntity> createBooking({
    required String fleetId,
    required DateTime startDate,
    required DateTime endDate,
    required String siteLocationAddress,
    List<String> attachmentIds = const [],
  });

  Future<BookingEntity> getBooking(String id);

  Future<List<BookingEntity>> listBookings({String? status});

  Future<BankSlipEntity> submitPaymentEvidence({
    required String bookingId,
    required String storageUri,
    required String fileHash,
    String? transactionReference,
    required int amountMinor,
    String currency = 'ETB',
  });

  Future<BookingEntity> confirmPayment({
    required String bookingId,
    required int expectedVersion,
    String? paymentReference,
    String? notes,
  });

  Future<BookingEntity> declinePayment({
    required String bookingId,
    required int expectedVersion,
    required String reasonCode,
    required String rejectionNotes,
  });

  Future<BookingEntity> disputeBooking({
    required String bookingId,
    required int expectedVersion,
    required String disputeReason,
    List<String> evidenceUris = const [],
  });

  Future<BookingEntity> cancelBooking({
    required String bookingId,
    required int expectedVersion,
    required String reason,
  });

  Future<BookingEntity> handoverCheckIn({
    required String bookingId,
    required int expectedVersion,
    required String handoverPin,
    required double operatingHoursMeter,
    List<String> photoUris = const [],
  });

  Future<BookingEntity> returnCheckOut({
    required String bookingId,
    required int expectedVersion,
    required double operatingHoursMeter,
    String? damageNotes,
    List<String> photoUris = const [],
  });
}
