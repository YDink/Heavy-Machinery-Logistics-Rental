import 'package:equatable/equatable.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_result.dart';

abstract class DiscoveryEvent extends Equatable {
  const DiscoveryEvent();

  @override
  List<Object?> get props => [];
}

class SearchFleetsRequested extends DiscoveryEvent {
  final FleetSearchFilter filter;

  const SearchFleetsRequested(this.filter);

  @override
  List<Object?> get props => [filter];
}

class ToggleViewModeRequested extends DiscoveryEvent {
  final bool isMapView;

  const ToggleViewModeRequested(this.isMapView);

  @override
  List<Object?> get props => [isMapView];
}

class ToggleCompareSelection extends DiscoveryEvent {
  final FleetSearchResult item;

  const ToggleCompareSelection(this.item);

  @override
  List<Object?> get props => [item];
}

class ClearCompareSelection extends DiscoveryEvent {
  const ClearCompareSelection();
}

class ExecuteComparisonRequested extends DiscoveryEvent {
  const ExecuteComparisonRequested();
}

class CheckAvailabilityRequested extends DiscoveryEvent {
  final String fleetId;
  final DateTime startDate;
  final DateTime endDate;

  const CheckAvailabilityRequested({
    required this.fleetId,
    required this.startDate,
    required this.endDate,
  });

  @override
  List<Object?> get props => [fleetId, startDate, endDate];
}
