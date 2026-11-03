import '../entities/booking.dart';
import '../repositories/booking_repository.dart';

class SubmitPaymentSlipParams {
  final String bookingId;
  final String storageUri;
  final String fileHash;
  final String? transactionReference;
  final int amountMinor;
  final String currency;

  const SubmitPaymentSlipParams({
    required this.bookingId,
    required this.storageUri,
    required this.fileHash,
    this.transactionReference,
    required this.amountMinor,
    this.currency = 'ETB',
  });
}

class SubmitPaymentSlipUseCase {
  final BookingRepository repository;

  SubmitPaymentSlipUseCase(this.repository);

  Future<BankSlipEntity> call(SubmitPaymentSlipParams params) {
    return repository.submitPaymentEvidence(
      bookingId: params.bookingId,
      storageUri: params.storageUri,
      fileHash: params.fileHash,
      transactionReference: params.transactionReference,
      amountMinor: params.amountMinor,
      currency: params.currency,
    );
  }
}
