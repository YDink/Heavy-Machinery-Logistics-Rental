import 'package:dink_mobile/features/fleet_catalog/domain/entities/spec_field_group.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/repositories/fleet_repository.dart';

class GetSpecDefinitionsUseCase {
  final FleetRepository repository;

  GetSpecDefinitionsUseCase(this.repository);

  Future<List<SpecFieldGroup>> call(String categoryId) {
    return repository.getSpecDefinitions(categoryId);
  }
}
