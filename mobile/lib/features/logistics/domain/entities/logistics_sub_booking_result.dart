class LogisticsSubBookingResult {
  final String subBookingId;
  final String parentBookingId;
  final String status;
  final int freightHireFeeMinor;
  final int platformCommissionBps;
  final String platformCommissionRate;
  final int commissionOwedMinor;
  final String currency;
  final bool abnormalLoadEscortRequired;
  final int escortFeeMinor;
  final String carrierFullName;
  final String carrierPhoneNumber;

  const LogisticsSubBookingResult({
    required this.subBookingId,
    required this.parentBookingId,
    required this.status,
    required this.freightHireFeeMinor,
    required this.platformCommissionBps,
    required this.platformCommissionRate,
    required this.commissionOwedMinor,
    required this.currency,
    required this.abnormalLoadEscortRequired,
    required this.escortFeeMinor,
    required this.carrierFullName,
    required this.carrierPhoneNumber,
  });
}
