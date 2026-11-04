import 'package:equatable/equatable.dart';
import 'fleet_search_result.dart';

class FleetSearchData extends Equatable {
  final int totalMatches;
  final double latitude;
  final double longitude;
  final double radiusKm;
  final List<FleetSearchResult> results;

  const FleetSearchData({
    required this.totalMatches,
    required this.latitude,
    required this.longitude,
    required this.radiusKm,
    required this.results,
  });

  @override
  List<Object?> get props => [totalMatches, latitude, longitude, radiusKm, results];
}
