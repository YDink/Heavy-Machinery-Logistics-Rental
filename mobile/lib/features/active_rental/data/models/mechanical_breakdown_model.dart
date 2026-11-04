import '../../domain/entities/mechanical_breakdown.dart';

class MechanicalBreakdownModel extends MechanicalBreakdown {
  const MechanicalBreakdownModel({
    required super.incidentId,
    required super.bookingId,
    required super.status,
    super.ownerSlaDeadline,
    super.shiftTimerPausedAt,
  });

  factory MechanicalBreakdownModel.fromJson(Map<String, dynamic> json) {
    return MechanicalBreakdownModel(
      incidentId: json['incidentId'] as String? ?? '',
      bookingId: json['bookingId'] as String? ?? '',
      status: json['status'] as String? ?? '',
      ownerSlaDeadline: json['ownerSlaDeadline'] != null
          ? DateTime.tryParse(json['ownerSlaDeadline'] as String)
          : null,
      shiftTimerPausedAt: json['shiftTimerPausedAt'] != null
          ? DateTime.tryParse(json['shiftTimerPausedAt'] as String)
          : null,
    );
  }
}
