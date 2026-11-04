import 'package:dink_mobile/features/fleet_catalog/data/datasources/fleet_local_datasource.dart';
import 'package:dink_mobile/features/fleet_catalog/data/datasources/fleet_remote_datasource.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/category.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/spec_field_group.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/repositories/fleet_repository.dart';

class FleetRepositoryImpl implements FleetRepository {
  final FleetRemoteDataSource remoteDataSource;
  final FleetLocalDataSource localDataSource;

  FleetRepositoryImpl({
    required this.remoteDataSource,
    required this.localDataSource,
  });

  @override
  Future<List<Category>> getCategories() async {
    return remoteDataSource.getCategories();
  }

  @override
  Future<List<SpecFieldGroup>> getSpecDefinitions(String categoryId) async {
    return remoteDataSource.getSpecDefinitions(categoryId);
  }

  @override
  Future<FleetAsset> registerFleet(FleetRegistrationRequest request) async {
    final asset = await remoteDataSource.registerFleet(request);
    await localDataSource.clearDraft(request.categoryId);
    return asset;
  }

  @override
  Future<void> saveDraft(FleetRegistrationRequest draft) async {
    await localDataSource.saveDraft(draft);
  }

  @override
  Future<FleetRegistrationRequest?> getDraft(String categoryId) async {
    return localDataSource.getDraft(categoryId);
  }
}
