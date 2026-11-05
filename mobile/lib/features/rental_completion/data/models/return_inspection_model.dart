import '../../domain/entities/return_inspection.dart';

class ReturnInspectionModel extends ReturnInspection {
  const ReturnInspectionModel({
    required super.bookingId,
    required super.checkOutMeterHours,
    required super.damageAuditWindowExpiresAt,
    required super.inspectedAt,
    required super.photosIngestedCount,
    required super.status,
  });

  factory ReturnInspectionModel.fromJson(Map<String, dynamic> json) {
    return ReturnInspectionModel(
      bookingId: json['bookingId'] as String? ?? '',
      checkOutMeterHours: (json['checkOutMeterHours'] as num?)?.toDouble() ?? 0.0,
      damageAuditWindowExpiresAt: DateTime.tryParse(json['damageAuditWindowExpiresAt'] as String? ?? '') ?? DateTime.now(),
      inspectedAt: DateTime.tryParse(json['inspectedAt'] as String? ?? '') ?? DateTime.now(),
      photosIngestedCount: (json['photosIngestedCount'] as num?)?.toInt() ?? 0,
      status: json['status'] as String? ?? '',
    );
  }
}
