import 'package:dink_mobile/features/fleet_catalog/domain/entities/category.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/spec_field_group.dart';

abstract class FleetRepository {
  Future<List<Category>> getCategories();
  Future<List<SpecFieldGroup>> getSpecDefinitions(String categoryId);
  Future<FleetAsset> registerFleet(FleetRegistrationRequest request);
  Future<void> saveDraft(FleetRegistrationRequest draft);
  Future<FleetRegistrationRequest?> getDraft(String categoryId);
}
