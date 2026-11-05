import 'package:equatable/equatable.dart';

class ReturnInspection extends Equatable {
  final String bookingId;
  final double checkOutMeterHours;
  final DateTime damageAuditWindowExpiresAt;
  final DateTime inspectedAt;
  final int photosIngestedCount;
  final String status;

  const ReturnInspection({
    required this.bookingId,
    required this.checkOutMeterHours,
    required this.damageAuditWindowExpiresAt,
    required this.inspectedAt,
    required this.photosIngestedCount,
    required this.status,
  });

  @override
  List<Object?> get props => [
        bookingId,
        checkOutMeterHours,
        damageAuditWindowExpiresAt,
        inspectedAt,
        photosIngestedCount,
        status,
      ];
}
