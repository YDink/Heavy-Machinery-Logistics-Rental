import 'package:dink_mobile/features/discovery/domain/entities/fleet_availability.dart';
import 'package:dink_mobile/features/discovery/domain/repositories/discovery_repository.dart';

class GetFleetAvailabilityUseCase {
  final DiscoveryRepository repository;

  GetFleetAvailabilityUseCase(this.repository);

  Future<FleetAvailabilityData> call(String fleetId, DateTime startDate, DateTime endDate) {
    return repository.getFleetAvailability(fleetId, startDate, endDate);
  }
}
