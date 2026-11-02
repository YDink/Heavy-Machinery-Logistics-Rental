import 'package:dink_mobile/features/discovery/domain/entities/fleet_compare_matrix.dart';
import 'package:dink_mobile/features/discovery/domain/repositories/discovery_repository.dart';

class CompareFleetsUseCase {
  final DiscoveryRepository repository;

  CompareFleetsUseCase(this.repository);

  Future<FleetCompareData> call(List<String> fleetIds) {
    return repository.compareFleets(fleetIds);
  }
}
