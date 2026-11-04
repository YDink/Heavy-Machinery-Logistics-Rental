import '../datasources/fleet_local_data_source.dart';
import '../datasources/fleet_remote_data_source.dart';
import '../models/fleet_registration_request_model.dart';
import '../../domain/entities/category.dart';
import '../../domain/entities/fleet_asset.dart';
import '../../domain/entities/spec_field_group.dart';
import '../../domain/repositories/fleet_repository.dart';

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
    final model = FleetRegistrationRequestModel.fromEntity(request);
    final asset = await remoteDataSource.registerFleet(model);
    await localDataSource.clearDraft(request.categoryId);
    return asset;
  }

  @override
  Future<void> saveDraft(FleetRegistrationRequest draft) async {
    final model = FleetRegistrationRequestModel.fromEntity(draft);
    await localDataSource.saveDraft(model);
  }

  @override
  Future<FleetRegistrationRequest?> getDraft(String categoryId) async {
    final model = await localDataSource.getDraft(categoryId);
    if (model == null) return null;
    return FleetRegistrationRequest(
      categoryId: model.categoryId,
      make: model.make,
      model: model.model,
      manufacturingYear: model.manufacturingYear,
      registrationNumber: model.registrationNumber,
      chassisVin: model.chassisVin,
      baselineMeterHours: model.baselineMeterHours,
      fuelType: model.fuelType,
      selfMoveCapable: model.selfMoveCapable,
      dailyRateMinor: model.dailyRateMinor,
      excessHourRateMinor: model.excessHourRateMinor,
      minRentalDays: model.minRentalDays,
      paymentTimeoutHours: model.paymentTimeoutHours,
      operatorIncluded: model.operatorIncluded,
      latitude: model.latitude,
      longitude: model.longitude,
      addressText: model.addressText,
      photos: model.photos
          .map((p) => FleetPhoto(
                angle: p.angle,
                storageUri: p.storageUri,
                fileHash: p.fileHash,
                localFilePath: p.localFilePath,
              ))
          .toList(),
      specValues: model.specValues
          .map((v) => FleetSpecValueInput(
                fieldKey: v.fieldKey,
                valueString: v.valueString,
                valueNumeric: v.valueNumeric,
                valueBoolean: v.valueBoolean,
              ))
          .toList(),
    );
  }
}
