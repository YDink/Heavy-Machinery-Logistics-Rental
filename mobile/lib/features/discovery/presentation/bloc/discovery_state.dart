import 'package:equatable/equatable.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_availability.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_compare_matrix.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_result.dart';

enum DiscoveryStatus {
  initial,
  loading,
  success,
  comparing,
  compareSuccess,
  checkingAvailability,
  availabilitySuccess,
  failure,
}

class DiscoveryState extends Equatable {
  final DiscoveryStatus status;
  final FleetSearchFilter filter;
  final FleetSearchData? searchData;
  final bool isMapView;
  final List<FleetSearchResult> selectedForCompare;
  final FleetCompareData? compareData;
  final FleetAvailabilityData? availabilityData;
  final String? errorMessage;

  const DiscoveryState({
    this.status = DiscoveryStatus.initial,
    required this.filter,
    this.searchData,
    this.isMapView = false,
    this.selectedForCompare = const [],
    this.compareData,
    this.availabilityData,
    this.errorMessage,
  });

  DiscoveryState copyWith({
    DiscoveryStatus? status,
    FleetSearchFilter? filter,
    FleetSearchData? searchData,
    bool? isMapView,
    List<FleetSearchResult>? selectedForCompare,
    FleetCompareData? compareData,
    FleetAvailabilityData? availabilityData,
    String? errorMessage,
  }) {
    return DiscoveryState(
      status: status ?? this.status,
      filter: filter ?? this.filter,
      searchData: searchData ?? this.searchData,
      isMapView: isMapView ?? this.isMapView,
      selectedForCompare: selectedForCompare ?? this.selectedForCompare,
      compareData: compareData ?? this.compareData,
      availabilityData: availabilityData ?? this.availabilityData,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }

  @override
  List<Object?> get props => [
        status,
        filter,
        searchData,
        isMapView,
        selectedForCompare,
        compareData,
        availabilityData,
        errorMessage,
      ];
}
