import '../../domain/entities/booking.dart';
import '../../domain/repositories/booking_repository.dart';
import '../datasources/booking_remote_data_source.dart';

class BookingRepositoryImpl implements BookingRepository {
  final BookingRemoteDataSource remoteDataSource;

  BookingRepositoryImpl({required this.remoteDataSource});

  @override
  Future<BookingEntity> createBooking({
    required String fleetId,
    required DateTime startDate,
    required DateTime endDate,
    required String siteLocationAddress,
    List<String> attachmentIds = const [],
  }) {
    return remoteDataSource.createBooking(
      fleetId: fleetId,
      startDate: startDate,
      endDate: endDate,
      siteLocationAddress: siteLocationAddress,
      attachmentIds: attachmentIds,
    );
  }

  @override
  Future<BookingEntity> getBooking(String id) {
    return remoteDataSource.getBooking(id);
  }

  @override
  Future<List<BookingEntity>> listBookings({String? status}) {
    return remoteDataSource.listBookings(status: status);
  }

  @override
  Future<BankSlipEntity> submitPaymentEvidence({
    required String bookingId,
    required String storageUri,
    required String fileHash,
    String? transactionReference,
    required int amountMinor,
    String currency = 'ETB',
  }) {
    return remoteDataSource.submitPaymentEvidence(
      bookingId: bookingId,
      storageUri: storageUri,
      fileHash: fileHash,
      transactionReference: transactionReference,
      amountMinor: amountMinor,
      currency: currency,
    );
  }

  @override
  Future<BookingEntity> confirmPayment({
    required String bookingId,
    required int expectedVersion,
    String? paymentReference,
    String? notes,
  }) {
    return remoteDataSource.confirmPayment(
      bookingId: bookingId,
      expectedVersion: expectedVersion,
      paymentReference: paymentReference,
      notes: notes,
    );
  }

  @override
  Future<BookingEntity> declinePayment({
    required String bookingId,
    required int expectedVersion,
    required String reasonCode,
    required String rejectionNotes,
  }) {
    return remoteDataSource.declinePayment(
      bookingId: bookingId,
      expectedVersion: expectedVersion,
      reasonCode: reasonCode,
      rejectionNotes: rejectionNotes,
    );
  }

  @override
  Future<BookingEntity> disputeBooking({
    required String bookingId,
    required int expectedVersion,
    required String disputeReason,
    List<String> evidenceUris = const [],
  }) {
    return remoteDataSource.disputeBooking(
      bookingId: bookingId,
      expectedVersion: expectedVersion,
      disputeReason: disputeReason,
      evidenceUris: evidenceUris,
    );
  }

  @override
  Future<BookingEntity> cancelBooking({
    required String bookingId,
    required int expectedVersion,
    required String reason,
  }) {
    return remoteDataSource.cancelBooking(
      bookingId: bookingId,
      expectedVersion: expectedVersion,
      reason: reason,
    );
  }

  @override
  Future<BookingEntity> handoverCheckIn({
    required String bookingId,
    required int expectedVersion,
    required String handoverPin,
    required double operatingHoursMeter,
    List<String> photoUris = const [],
  }) {
    return remoteDataSource.handoverCheckIn(
      bookingId: bookingId,
      expectedVersion: expectedVersion,
      handoverPin: handoverPin,
      operatingHoursMeter: operatingHoursMeter,
      photoUris: photoUris,
    );
  }

  @override
  Future<BookingEntity> returnCheckOut({
    required String bookingId,
    required int expectedVersion,
    required double operatingHoursMeter,
    String? damageNotes,
    List<String> photoUris = const [],
  }) {
    return remoteDataSource.returnCheckOut(
      bookingId: bookingId,
      expectedVersion: expectedVersion,
      operatingHoursMeter: operatingHoursMeter,
      damageNotes: damageNotes,
      photoUris: photoUris,
    );
  }
}
