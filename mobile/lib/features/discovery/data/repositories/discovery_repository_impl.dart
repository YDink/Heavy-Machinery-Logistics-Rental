import 'package:dink_mobile/features/discovery/data/datasources/discovery_remote_datasource.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_availability.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_compare_matrix.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_result.dart';
import 'package:dink_mobile/features/discovery/domain/repositories/discovery_repository.dart';

class DiscoveryRepositoryImpl implements DiscoveryRepository {
  final DiscoveryRemoteDataSource remoteDataSource;

  DiscoveryRepositoryImpl({required this.remoteDataSource});

  @override
  Future<FleetSearchData> searchFleets(FleetSearchFilter filter) {
    return remoteDataSource.searchFleets(filter);
  }

  @override
  Future<FleetCompareData> compareFleets(List<String> fleetIds) {
    return remoteDataSource.compareFleets(fleetIds);
  }

  @override
  Future<FleetAvailabilityData> getFleetAvailability(
    String fleetId,
    DateTime startDate,
    DateTime endDate,
  ) {
    return remoteDataSource.getFleetAvailability(fleetId, startDate, endDate);
  }
}
