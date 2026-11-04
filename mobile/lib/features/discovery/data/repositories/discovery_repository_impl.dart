import '../datasources/discovery_remote_data_source.dart';
import '../models/fleet_search_filter_model.dart';
import '../../domain/entities/fleet_availability.dart';
import '../../domain/entities/fleet_compare_matrix.dart';
import '../../domain/entities/fleet_search_filter.dart';
import '../../domain/entities/fleet_search_result.dart';
import '../../domain/repositories/discovery_repository.dart';

class DiscoveryRepositoryImpl implements DiscoveryRepository {
  final DiscoveryRemoteDataSource remoteDataSource;

  DiscoveryRepositoryImpl({required this.remoteDataSource});

  @override
  Future<FleetSearchData> searchFleets(FleetSearchFilter filter) {
    return remoteDataSource.searchFleets(FleetSearchFilterModel.fromEntity(filter));
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
