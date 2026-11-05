import '../../domain/entities/damage_claim_response_result.dart';

class DamageClaimResponseResultModel extends DamageClaimResponseResult {
  const DamageClaimResponseResultModel({
    required super.claimId,
    required super.bookingId,
    required super.acceptedAmountMinor,
    required super.status,
    required super.updatedAt,
  });

  factory DamageClaimResponseResultModel.fromJson(Map<String, dynamic> json) {
    return DamageClaimResponseResultModel(
      claimId: json['claimId'] as String? ?? '',
      bookingId: json['bookingId'] as String? ?? '',
      acceptedAmountMinor: (json['acceptedAmountMinor'] as num?)?.toInt() ?? 0,
      status: json['status'] as String? ?? '',
      updatedAt: DateTime.tryParse(json['updatedAt'] as String? ?? '') ?? DateTime.now(),
    );
  }
}
