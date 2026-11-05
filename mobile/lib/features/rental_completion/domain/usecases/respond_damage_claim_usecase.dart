import '../entities/damage_claim_response_result.dart';
import '../repositories/rental_completion_repository.dart';

class RespondDamageClaimParams {
  final String bookingId;
  final String claimId;
  final String action;
  final int? acceptedAmountMinor;
  final String? disputeReason;

  const RespondDamageClaimParams({
    required this.bookingId,
    required this.claimId,
    required this.action,
    this.acceptedAmountMinor,
    this.disputeReason,
  });
}

class RespondDamageClaimUseCase {
  final RentalCompletionRepository repository;

  RespondDamageClaimUseCase(this.repository);

  Future<DamageClaimResponseResult> call(RespondDamageClaimParams params) {
    return repository.respondDamageClaim(
      bookingId: params.bookingId,
      claimId: params.claimId,
      action: params.action,
      acceptedAmountMinor: params.acceptedAmountMinor,
      disputeReason: params.disputeReason,
    );
  }
}
