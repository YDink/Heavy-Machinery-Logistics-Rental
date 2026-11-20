import '../../domain/entities/logistics_sub_booking_result.dart';

class LogisticsSubBookingResultModel {
  final String subBookingId;
  final String parentBookingId;
  final String status;
  final int freightHireFeeMinor;
  final int platformCommissionBps;
  final String platformCommissionRate;
  final int commissionOwedMinor;
  final String currency;
  final bool abnormalLoadEscortRequired;
  final int escortFeeMinor;
  final String carrierFullName;
  final String carrierPhoneNumber;

  const LogisticsSubBookingResultModel({
    required this.subBookingId,
    required this.parentBookingId,
    required this.status,
    required this.freightHireFeeMinor,
    required this.platformCommissionBps,
    required this.platformCommissionRate,
    required this.commissionOwedMinor,
    required this.currency,
    required this.abnormalLoadEscortRequired,
    required this.escortFeeMinor,
    required this.carrierFullName,
    required this.carrierPhoneNumber,
  });

  factory LogisticsSubBookingResultModel.fromJson(Map<String, dynamic> json) {
    final contact = json['carrierContact'] as Map<String, dynamic>? ?? const {};
    return LogisticsSubBookingResultModel(
      subBookingId: json['subBookingId'] as String? ?? '',
      parentBookingId: json['parentBookingId'] as String? ?? '',
      status: json['status'] as String? ?? 'CONFIRMED',
      freightHireFeeMinor: (json['freightHireFeeMinor'] as num?)?.toInt() ?? 0,
      platformCommissionBps: (json['platformCommissionBps'] as num?)?.toInt() ?? 250,
      platformCommissionRate: json['platformCommissionRate'] as String? ?? '0.025',
      commissionOwedMinor: (json['commissionOwedMinor'] as num?)?.toInt() ?? 0,
      currency: json['currency'] as String? ?? 'ETB',
      abnormalLoadEscortRequired: json['abnormalLoadEscortRequired'] as bool? ?? false,
      escortFeeMinor: (json['escortFeeMinor'] as num?)?.toInt() ?? 0,
      carrierFullName: contact['fullName'] as String? ?? '',
      carrierPhoneNumber: contact['phoneNumber'] as String? ?? '',
    );
  }

  LogisticsSubBookingResult toEntity() {
    return LogisticsSubBookingResult(
      subBookingId: subBookingId,
      parentBookingId: parentBookingId,
      status: status,
      freightHireFeeMinor: freightHireFeeMinor,
      platformCommissionBps: platformCommissionBps,
      platformCommissionRate: platformCommissionRate,
      commissionOwedMinor: commissionOwedMinor,
      currency: currency,
      abnormalLoadEscortRequired: abnormalLoadEscortRequired,
      escortFeeMinor: escortFeeMinor,
      carrierFullName: carrierFullName,
      carrierPhoneNumber: carrierPhoneNumber,
    );
  }
}
