import 'package:equatable/equatable.dart';

class ReturnInitiateResult extends Equatable {
  final String bookingId;
  final String status;
  final double checkInMeterHours;
  final double allowedShiftHours;
  final DateTime initiatedAt;

  const ReturnInitiateResult({
    required this.bookingId,
    required this.status,
    required this.checkInMeterHours,
    required this.allowedShiftHours,
    required this.initiatedAt,
  });

  @override
  List<Object?> get props => [
        bookingId,
        status,
        checkInMeterHours,
        allowedShiftHours,
        initiatedAt,
      ];
}
