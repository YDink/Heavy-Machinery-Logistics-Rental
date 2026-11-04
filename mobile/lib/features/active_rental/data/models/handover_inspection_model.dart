import '../../domain/entities/handover_inspection.dart';

class HandoverInspectionModel extends HandoverInspection {
  const HandoverInspectionModel({
    required super.bookingId,
    required super.inspectedAt,
    required super.meterVerified,
    required super.photosIngestedCount,
    required super.status,
  });

  factory HandoverInspectionModel.fromJson(Map<String, dynamic> json) {
    return HandoverInspectionModel(
      bookingId: json['bookingId'] as String? ?? '',
      inspectedAt: DateTime.tryParse(json['inspectedAt'] as String? ?? '') ?? DateTime.now(),
      meterVerified: json['meterVerified'] as bool? ?? false,
      photosIngestedCount: (json['photosIngestedCount'] as num?)?.toInt() ?? 0,
      status: json['status'] as String? ?? '',
    );
  }
}
