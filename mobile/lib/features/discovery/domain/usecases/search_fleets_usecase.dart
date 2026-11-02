import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_result.dart';
import 'package:dink_mobile/features/discovery/domain/repositories/discovery_repository.dart';

class SearchFleetsUseCase {
  final DiscoveryRepository repository;

  SearchFleetsUseCase(this.repository);

  Future<FleetSearchData> call(FleetSearchFilter filter) {
    return repository.searchFleets(filter);
  }
}
