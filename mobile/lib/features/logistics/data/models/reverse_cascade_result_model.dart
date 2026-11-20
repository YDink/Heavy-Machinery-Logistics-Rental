import '../../domain/entities/reverse_cascade_result.dart';

class ReverseCascadeResultModel {
  final String parentBookingId;
  final String subBookingId;
  final String resolution;
  final int renterRefundDueMinor;
  final int ownerPenaltyLeviedMinor;
  final int ownerCommissionRestitutionMinor;
  final int carrierDefaultPenaltyMinor;
  final String currency;
  final String parentStatus;
  final String subStatus;

  const ReverseCascadeResultModel({
    required this.parentBookingId,
    required this.subBookingId,
    required this.resolution,
    required this.renterRefundDueMinor,
    required this.ownerPenaltyLeviedMinor,
    required this.ownerCommissionRestitutionMinor,
    required this.carrierDefaultPenaltyMinor,
    required this.currency,
    required this.parentStatus,
    required this.subStatus,
  });

  factory ReverseCascadeResultModel.fromJson(Map<String, dynamic> json) {
    return ReverseCascadeResultModel(
      parentBookingId: json['parentBookingId'] as String? ?? '',
      subBookingId: json['subBookingId'] as String? ?? '',
      resolution: json['resolution'] as String? ?? 'FORCE_MAJEURE_CANCELLATION',
      renterRefundDueMinor: (json['renterRefundDueMinor'] as num?)?.toInt() ?? 0,
      ownerPenaltyLeviedMinor: (json['ownerPenaltyLeviedMinor'] as num?)?.toInt() ?? 0,
      ownerCommissionRestitutionMinor:
          (json['ownerCommissionRestitutionMinor'] as num?)?.toInt() ?? 0,
      carrierDefaultPenaltyMinor:
          (json['carrierDefaultPenaltyMinor'] as num?)?.toInt() ?? 0,
      currency: json['currency'] as String? ?? 'ETB',
      parentStatus: json['parentStatus'] as String? ?? 'CANCELLED',
      subStatus: json['subStatus'] as String? ?? 'CANCELLED_CARRIER_DEFAULT',
    );
  }

  ReverseCascadeResult toEntity() {
    return ReverseCascadeResult(
      parentBookingId: parentBookingId,
      subBookingId: subBookingId,
      resolution: resolution,
      renterRefundDueMinor: renterRefundDueMinor,
      ownerPenaltyLeviedMinor: ownerPenaltyLeviedMinor,
      ownerCommissionRestitutionMinor: ownerCommissionRestitutionMinor,
      carrierDefaultPenaltyMinor: carrierDefaultPenaltyMinor,
      currency: currency,
      parentStatus: parentStatus,
      subStatus: subStatus,
    );
  }
}
