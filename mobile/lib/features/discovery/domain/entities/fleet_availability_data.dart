import 'package:equatable/equatable.dart';
import 'fleet_availability_date_range.dart';

class FleetAvailabilityData extends Equatable {
  final String fleetId;
  final bool isAvailable;
  final String? unavailableReason;
  final int? holdExpiresInSeconds;
  final List<FleetAvailabilityDateRange> blackoutRanges;
  final List<FleetAvailabilityDateRange> reservedRanges;

  const FleetAvailabilityData({
    required this.fleetId,
    required this.isAvailable,
    this.unavailableReason,
    this.holdExpiresInSeconds,
    required this.blackoutRanges,
    required this.reservedRanges,
  });

  @override
  List<Object?> get props => [
        fleetId,
        isAvailable,
        unavailableReason,
        holdExpiresInSeconds,
        blackoutRanges,
        reservedRanges,
      ];
}
