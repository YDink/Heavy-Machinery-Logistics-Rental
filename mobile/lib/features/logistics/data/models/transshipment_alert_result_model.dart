import '../../domain/entities/transshipment_alert_result.dart';

class TransshipmentAlertResultModel {
  final String transshipmentId;
  final String subBookingId;
  final String parentBookingId;
  final String subBookingStatus;
  final String parentBookingStatus;
  final DateTime hotSwapSlaDeadline;
  final int transshipmentBountyMinor;
  final String currency;
  final int candidateCarriersCount;

  const TransshipmentAlertResultModel({
    required this.transshipmentId,
    required this.subBookingId,
    required this.parentBookingId,
    required this.subBookingStatus,
    required this.parentBookingStatus,
    required this.hotSwapSlaDeadline,
    required this.transshipmentBountyMinor,
    required this.currency,
    required this.candidateCarriersCount,
  });

  factory TransshipmentAlertResultModel.fromJson(Map<String, dynamic> json) {
    return TransshipmentAlertResultModel(
      transshipmentId: json['transshipmentId'] as String? ?? '',
      subBookingId: json['subBookingId'] as String? ?? '',
      parentBookingId: json['parentBookingId'] as String? ?? '',
      subBookingStatus: json['subBookingStatus'] as String? ?? 'TRANSSHIPMENT_REQUIRED',
      parentBookingStatus: json['parentBookingStatus'] as String? ?? 'AWAITING_FLEET_READINESS',
      hotSwapSlaDeadline: DateTime.tryParse(
            json['hotSwapSlaDeadline'] as String? ?? '',
          ) ??
          DateTime.now().toUtc().add(const Duration(minutes: 60)),
      transshipmentBountyMinor: (json['transshipmentBountyMinor'] as num?)?.toInt() ?? 0,
      currency: json['currency'] as String? ?? 'ETB',
      candidateCarriersCount: (json['candidateCarriersCount'] as num?)?.toInt() ?? 0,
    );
  }

  TransshipmentAlertResult toEntity() {
    return TransshipmentAlertResult(
      transshipmentId: transshipmentId,
      subBookingId: subBookingId,
      parentBookingId: parentBookingId,
      subBookingStatus: subBookingStatus,
      parentBookingStatus: parentBookingStatus,
      hotSwapSlaDeadline: hotSwapSlaDeadline,
      transshipmentBountyMinor: transshipmentBountyMinor,
      currency: currency,
      candidateCarriersCount: candidateCarriersCount,
    );
  }
}
