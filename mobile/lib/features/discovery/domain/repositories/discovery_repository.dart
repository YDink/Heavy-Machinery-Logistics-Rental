import 'package:dink_mobile/features/discovery/domain/entities/fleet_availability.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_compare_matrix.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_result.dart';

abstract class DiscoveryRepository {
  Future<FleetSearchData> searchFleets(FleetSearchFilter filter);
  Future<FleetCompareData> compareFleets(List<String> fleetIds);
  Future<FleetAvailabilityData> getFleetAvailability(String fleetId, DateTime startDate, DateTime endDate);
}
