import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/repositories/fleet_repository.dart';

class RegisterFleetUseCase {
  final FleetRepository repository;

  RegisterFleetUseCase(this.repository);

  Future<FleetAsset> call(FleetRegistrationRequest request) {
    return repository.registerFleet(request);
  }
}
