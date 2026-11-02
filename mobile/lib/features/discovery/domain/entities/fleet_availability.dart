import 'package:equatable/equatable.dart';

class FleetAvailabilityDateRange extends Equatable {
  final DateTime startDate;
  final DateTime endDate;
  final String? reason;
  final String? state;

  const FleetAvailabilityDateRange({
    required this.startDate,
    required this.endDate,
    this.reason,
    this.state,
  });

  @override
  List<Object?> get props => [startDate, endDate, reason, state];
}

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
