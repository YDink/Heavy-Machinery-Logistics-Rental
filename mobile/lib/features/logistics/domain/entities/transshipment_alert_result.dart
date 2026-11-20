class TransshipmentAlertResult {
  final String transshipmentId;
  final String subBookingId;
  final String parentBookingId;
  final String subBookingStatus;
  final String parentBookingStatus;
  final DateTime hotSwapSlaDeadline;
  final int transshipmentBountyMinor;
  final String currency;
  final int candidateCarriersCount;

  const TransshipmentAlertResult({
    required this.transshipmentId,
    required this.subBookingId,
    required this.parentBookingId,
    required this.subBookingStatus,
    required this.parentBookingStatus,
    required this.hotSwapSlaDeadline,
    required this.transshipmentBountyMinor,
    required this.currency,
    required this.candidateCarriersCount,
  });
}
