class ExternalTransport {
  final String parentBookingId;
  final String status;
  final String carrierName;
  final String driverFullName;
  final String driverPhoneNumber;
  final String primeMoverPlateNumber;
  final DateTime estimatedArrivalDateTime;

  const ExternalTransport({
    required this.parentBookingId,
    required this.status,
    required this.carrierName,
    required this.driverFullName,
    required this.driverPhoneNumber,
    required this.primeMoverPlateNumber,
    required this.estimatedArrivalDateTime,
  });
}
