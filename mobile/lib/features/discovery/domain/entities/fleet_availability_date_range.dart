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
