class LogisticsDispatchResult {
  final String subBookingId;
  final String parentBookingId;
  final String subBookingStatus;
  final String parentBookingStatus;
  final int winchingPhotosCount;
  final DateTime updatedAt;

  const LogisticsDispatchResult({
    required this.subBookingId,
    required this.parentBookingId,
    required this.subBookingStatus,
    required this.parentBookingStatus,
    required this.winchingPhotosCount,
    required this.updatedAt,
  });
}
