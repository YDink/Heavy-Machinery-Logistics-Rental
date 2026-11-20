import '../../domain/entities/logistics_dispatch_result.dart';

class LogisticsDispatchResultModel {
  final String subBookingId;
  final String parentBookingId;
  final String subBookingStatus;
  final String parentBookingStatus;
  final int winchingPhotosCount;
  final DateTime updatedAt;

  const LogisticsDispatchResultModel({
    required this.subBookingId,
    required this.parentBookingId,
    required this.subBookingStatus,
    required this.parentBookingStatus,
    required this.winchingPhotosCount,
    required this.updatedAt,
  });

  factory LogisticsDispatchResultModel.fromJson(Map<String, dynamic> json) {
    return LogisticsDispatchResultModel(
      subBookingId: json['subBookingId'] as String? ?? '',
      parentBookingId: json['parentBookingId'] as String? ?? '',
      subBookingStatus: json['subBookingStatus'] as String? ?? '',
      parentBookingStatus: json['parentBookingStatus'] as String? ?? '',
      winchingPhotosCount: (json['winchingPhotosCount'] as num?)?.toInt() ?? 0,
      updatedAt: DateTime.tryParse(json['updatedAt'] as String? ?? '') ?? DateTime.now().toUtc(),
    );
  }

  LogisticsDispatchResult toEntity() {
    return LogisticsDispatchResult(
      subBookingId: subBookingId,
      parentBookingId: parentBookingId,
      subBookingStatus: subBookingStatus,
      parentBookingStatus: parentBookingStatus,
      winchingPhotosCount: winchingPhotosCount,
      updatedAt: updatedAt,
    );
  }
}
