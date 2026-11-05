import '../../domain/entities/rental_settlement.dart';

class RentalSettlementModel extends RentalSettlement {
  const RentalSettlementModel({
    required super.bookingId,
    required super.baseMachineFeeMinor,
    required super.bundledAttachmentsFeeMinor,
    required super.contractualAllowedHours,
    required super.actualOperatingHours,
    required super.creditedBreakdownHours,
    required super.billableExcessHours,
    required super.excessHourlyRateMinor,
    required super.excessHoursSurchargeMinor,
    required super.acceptedDamageMinor,
    required super.grossLiabilitiesMinor,
    required super.prepaymentSettledMinor,
    required super.netBalanceMinor,
    required super.currency,
    required super.status,
    required super.reconciledAt,
  });

  factory RentalSettlementModel.fromJson(Map<String, dynamic> json) {
    return RentalSettlementModel(
      bookingId: json['bookingId'] as String? ?? '',
      baseMachineFeeMinor: (json['baseMachineFeeMinor'] as num?)?.toInt() ?? 0,
      bundledAttachmentsFeeMinor: (json['bundledAttachmentsFeeMinor'] as num?)?.toInt() ?? 0,
      contractualAllowedHours: (json['contractualAllowedHours'] as num?)?.toDouble() ?? 0.0,
      actualOperatingHours: (json['actualOperatingHours'] as num?)?.toDouble() ?? 0.0,
      creditedBreakdownHours: (json['creditedBreakdownHours'] as num?)?.toDouble() ?? 0.0,
      billableExcessHours: (json['billableExcessHours'] as num?)?.toDouble() ?? 0.0,
      excessHourlyRateMinor: (json['excessHourlyRateMinor'] as num?)?.toInt() ?? 0,
      excessHoursSurchargeMinor: (json['excessHoursSurchargeMinor'] as num?)?.toInt() ?? 0,
      acceptedDamageMinor: (json['acceptedDamageMinor'] as num?)?.toInt() ?? 0,
      grossLiabilitiesMinor: (json['grossLiabilitiesMinor'] as num?)?.toInt() ?? 0,
      prepaymentSettledMinor: (json['prepaymentSettledMinor'] as num?)?.toInt() ?? 0,
      netBalanceMinor: (json['netBalanceMinor'] as num?)?.toInt() ?? 0,
      currency: json['currency'] as String? ?? 'ETB',
      status: json['status'] as String? ?? '',
      reconciledAt: DateTime.tryParse(json['reconciledAt'] as String? ?? '') ?? DateTime.now(),
    );
  }
}
