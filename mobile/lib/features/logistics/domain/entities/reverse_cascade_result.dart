class ReverseCascadeResult {
  final String parentBookingId;
  final String subBookingId;
  final String resolution;
  final int renterRefundDueMinor;
  final int ownerPenaltyLeviedMinor;
  final int ownerCommissionRestitutionMinor;
  final int carrierDefaultPenaltyMinor;
  final String currency;
  final String parentStatus;
  final String subStatus;

  const ReverseCascadeResult({
    required this.parentBookingId,
    required this.subBookingId,
    required this.resolution,
    required this.renterRefundDueMinor,
    required this.ownerPenaltyLeviedMinor,
    required this.ownerCommissionRestitutionMinor,
    required this.carrierDefaultPenaltyMinor,
    required this.currency,
    required this.parentStatus,
    required this.subStatus,
  });
}
