import 'matched_carrier.dart';

class LogisticsMatchResult {
  final String parentBookingId;
  final bool selfMoveCapable;
  final String recommendedTrailerClass;
  final bool abnormalLoadEscortRequired;
  final double operatingWeightKg;
  final int transportLengthMm;
  final int transportWidthMm;
  final int transportHeightMm;
  final double requiredMinimumPayloadKg;
  final List<MatchedCarrier> availableLogisticsFleets;

  const LogisticsMatchResult({
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
}
