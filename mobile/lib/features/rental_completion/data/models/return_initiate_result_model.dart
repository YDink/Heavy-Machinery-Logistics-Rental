import '../../domain/entities/return_initiate_result.dart';

class ReturnInitiateResultModel extends ReturnInitiateResult {
  const ReturnInitiateResultModel({
    required super.bookingId,
    required super.status,
    required super.checkInMeterHours,
    required super.allowedShiftHours,
    required super.initiatedAt,
  });

  factory ReturnInitiateResultModel.fromJson(Map<String, dynamic> json) {
    return ReturnInitiateResultModel(
      bookingId: json['bookingId'] as String? ?? '',
      status: json['status'] as String? ?? '',
      checkInMeterHours: (json['checkInMeterHours'] as num?)?.toDouble() ?? 0.0,
      allowedShiftHours: (json['allowedShiftHours'] as num?)?.toDouble() ?? 0.0,
      initiatedAt: DateTime.tryParse(json['initiatedAt'] as String? ?? '') ?? DateTime.now(),
    );
  }
}
