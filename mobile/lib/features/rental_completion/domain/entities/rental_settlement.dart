import 'package:equatable/equatable.dart';

class RentalSettlement extends Equatable {
  final String bookingId;
  final int baseMachineFeeMinor;
  final int bundledAttachmentsFeeMinor;
  final double contractualAllowedHours;
  final double actualOperatingHours;
  final double creditedBreakdownHours;
  final double billableExcessHours;
  final int excessHourlyRateMinor;
  final int excessHoursSurchargeMinor;
  final int acceptedDamageMinor;
  final int grossLiabilitiesMinor;
  final int prepaymentSettledMinor;
  final int netBalanceMinor;
  final String currency;
  final String status;
  final DateTime reconciledAt;

  const RentalSettlement({
    required this.bookingId,
    required this.baseMachineFeeMinor,
    required this.bundledAttachmentsFeeMinor,
    required this.contractualAllowedHours,
    required this.actualOperatingHours,
    required this.creditedBreakdownHours,
    required this.billableExcessHours,
    required this.excessHourlyRateMinor,
    required this.excessHoursSurchargeMinor,
    required this.acceptedDamageMinor,
    required this.grossLiabilitiesMinor,
    required this.prepaymentSettledMinor,
    required this.netBalanceMinor,
    required this.currency,
    required this.status,
    required this.reconciledAt,
  });

  @override
  List<Object?> get props => [
        bookingId,
        baseMachineFeeMinor,
        bundledAttachmentsFeeMinor,
        contractualAllowedHours,
        actualOperatingHours,
        creditedBreakdownHours,
        billableExcessHours,
        excessHourlyRateMinor,
        excessHoursSurchargeMinor,
        acceptedDamageMinor,
        grossLiabilitiesMinor,
        prepaymentSettledMinor,
        netBalanceMinor,
        currency,
        status,
        reconciledAt,
      ];
}
