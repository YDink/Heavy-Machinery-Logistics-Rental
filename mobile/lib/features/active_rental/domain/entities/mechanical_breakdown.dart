import 'package:equatable/equatable.dart';

class MechanicalBreakdown extends Equatable {
  final String incidentId;
  final String bookingId;
  final String status;
  final DateTime? ownerSlaDeadline;
  final DateTime? shiftTimerPausedAt;

  const MechanicalBreakdown({
    required this.incidentId,
    required this.bookingId,
    required this.status,
    this.ownerSlaDeadline,
    this.shiftTimerPausedAt,
  });

  @override
  List<Object?> get props => [incidentId, bookingId, status, ownerSlaDeadline, shiftTimerPausedAt];
}
