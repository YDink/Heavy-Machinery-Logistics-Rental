import '../../domain/entities/logistics_match_result.dart';
import 'matched_carrier_model.dart';

class LogisticsMatchResultModel {
  final String parentBookingId;
  final bool selfMoveCapable;
  final String recommendedTrailerClass;
  final bool abnormalLoadEscortRequired;
  final double operatingWeightKg;
  final int transportLengthMm;
  final int transportWidthMm;
  final int transportHeightMm;
  final double requiredMinimumPayloadKg;
  final List<MatchedCarrierModel> availableLogisticsFleets;

  const LogisticsMatchResultModel({
    required this.parentBookingId,
    required this.selfMoveCapable,
    required this.recommendedTrailerClass,
    required this.abnormalLoadEscortRequired,
    required this.operatingWeightKg,
    required this.transportLengthMm,
    required this.transportWidthMm,
    required this.transportHeightMm,
    required this.requiredMinimumPayloadKg,
    required this.availableLogisticsFleets,
  });

  factory LogisticsMatchResultModel.fromJson(Map<String, dynamic> json) {
    final specs = json['machineSpecs'] as Map<String, dynamic>? ?? const {};
    final fleetsRaw = json['availableLogisticsFleets'] as List<dynamic>? ?? const [];
    return LogisticsMatchResultModel(
      parentBookingId: json['parentBookingId'] as String? ?? '',
      selfMoveCapable: json['selfMoveCapable'] as bool? ?? false,
      recommendedTrailerClass: json['recommendedTrailerClass'] as String? ?? '',
      abnormalLoadEscortRequired: json['abnormalLoadEscortRequired'] as bool? ?? false,
      operatingWeightKg: (specs['operatingWeightKg'] as num?)?.toDouble() ?? 0.0,
      transportLengthMm: (specs['transportLengthMm'] as num?)?.toInt() ?? 0,
      transportWidthMm: (specs['transportWidthMm'] as num?)?.toInt() ?? 0,
      transportHeightMm: (specs['transportHeightMm'] as num?)?.toInt() ?? 0,
      requiredMinimumPayloadKg: (json['requiredMinimumPayloadKg'] as num?)?.toDouble() ?? 0.0,
      availableLogisticsFleets: fleetsRaw
          .map((e) => MatchedCarrierModel.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  LogisticsMatchResult toEntity() {
    return LogisticsMatchResult(
      parentBookingId: parentBookingId,
      selfMoveCapable: selfMoveCapable,
      recommendedTrailerClass: recommendedTrailerClass,
      abnormalLoadEscortRequired: abnormalLoadEscortRequired,
      operatingWeightKg: operatingWeightKg,
      transportLengthMm: transportLengthMm,
      transportWidthMm: transportWidthMm,
      transportHeightMm: transportHeightMm,
      requiredMinimumPayloadKg: requiredMinimumPayloadKg,
      availableLogisticsFleets:
          availableLogisticsFleets.map((e) => e.toEntity()).toList(),
    );
  }
}
