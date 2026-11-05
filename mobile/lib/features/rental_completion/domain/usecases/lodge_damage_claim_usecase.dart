import '../entities/damage_claim.dart';
import '../repositories/rental_completion_repository.dart';

class LodgeDamageClaimParams {
  final String bookingId;
  final String componentCategory;
  final String incidentNarrative;
  final String? repairQuoteUri;
  final int claimedAmountMinor;
  final String currency;
  final List<String> damagePhotos;

  const LodgeDamageClaimParams({
    required this.bookingId,
    required this.componentCategory,
    required this.incidentNarrative,
    this.repairQuoteUri,
    required this.claimedAmountMinor,
    required this.currency,
    this.damagePhotos = const [],
  });
}

class LodgeDamageClaimUseCase {
  final RentalCompletionRepository repository;

  LodgeDamageClaimUseCase(this.repository);

  Future<DamageClaim> call(LodgeDamageClaimParams params) {
    return repository.lodgeDamageClaim(
      bookingId: params.bookingId,
      componentCategory: params.componentCategory,
      incidentNarrative: params.incidentNarrative,
      repairQuoteUri: params.repairQuoteUri,
      claimedAmountMinor: params.claimedAmountMinor,
      currency: params.currency,
      damagePhotos: params.damagePhotos,
    );
  }
}
