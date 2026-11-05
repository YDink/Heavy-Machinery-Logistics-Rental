import '../../domain/entities/damage_claim.dart';

class DamageClaimModel extends DamageClaim {
  const DamageClaimModel({
    required super.claimId,
    required super.bookingId,
    required super.componentCategory,
    required super.claimedAmountMinor,
    required super.acceptedAmountMinor,
    required super.currency,
    required super.status,
    required super.createdAt,
  });

  factory DamageClaimModel.fromJson(Map<String, dynamic> json) {
    return DamageClaimModel(
      claimId: json['claimId'] as String? ?? '',
      bookingId: json['bookingId'] as String? ?? '',
      componentCategory: json['componentCategory'] as String? ?? '',
      claimedAmountMinor: (json['claimedAmountMinor'] as num?)?.toInt() ?? 0,
      acceptedAmountMinor: (json['acceptedAmountMinor'] as num?)?.toInt() ?? 0,
      currency: json['currency'] as String? ?? 'ETB',
      status: json['status'] as String? ?? '',
      createdAt: DateTime.tryParse(json['createdAt'] as String? ?? '') ?? DateTime.now(),
    );
  }
}
