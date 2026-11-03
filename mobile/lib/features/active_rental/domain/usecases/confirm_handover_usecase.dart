import '../repositories/active_rental_repository.dart';

class ConfirmHandoverUseCase {
  final ActiveRentalRepository repository;

  ConfirmHandoverUseCase(this.repository);

  Future<void> call({
    required String bookingId,
    required String challengePin,
    String? touchSignatureBase64,
    double? latitude,
    double? longitude,
  }) {
    return repository.confirmHandover(
      bookingId: bookingId,
      challengePin: challengePin,
      touchSignatureBase64: touchSignatureBase64,
      latitude: latitude,
      longitude: longitude,
    );
  }
}
