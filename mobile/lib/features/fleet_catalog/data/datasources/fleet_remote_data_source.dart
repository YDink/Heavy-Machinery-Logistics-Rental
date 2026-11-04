import 'package:dio/dio.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import '../models/category_model.dart';
import '../models/fleet_asset_model.dart';
import '../models/fleet_registration_request_model.dart';
import '../models/spec_definition_model.dart';

abstract class FleetRemoteDataSource {
  Future<List<CategoryModel>> getCategories();
  Future<List<SpecFieldGroupModel>> getSpecDefinitions(String categoryId);
  Future<FleetAssetModel> registerFleet(FleetRegistrationRequestModel request);
}

class FleetRemoteDataSourceImpl implements FleetRemoteDataSource {
  final ApiClient apiClient;

  FleetRemoteDataSourceImpl({required this.apiClient});

  @override
  Future<List<CategoryModel>> getCategories() async {
    try {
      final response = await apiClient.dio.get('/categories');
      final data = response.data['data'] as List<dynamic>;
      return data
          .map((item) => CategoryModel.fromJson(item as Map<String, dynamic>))
          .toList();
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }

  @override
  Future<List<SpecFieldGroupModel>> getSpecDefinitions(String categoryId) async {
    try {
      final response = await apiClient.dio.get('/categories/$categoryId/spec-definitions');
      final data = response.data['data'] as List<dynamic>;
      return data
          .map((item) => SpecFieldGroupModel.fromJson(item as Map<String, dynamic>))
          .toList();
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }

  @override
  Future<FleetAssetModel> registerFleet(FleetRegistrationRequestModel request) async {
    try {
      final payload = {
        'categoryId': request.categoryId,
        'make': request.make,
        'model': request.model,
        'manufacturingYear': request.manufacturingYear,
        'registrationNumber': request.registrationNumber,
        'chassisVin': request.chassisVin,
        'baselineMeterHours': request.baselineMeterHours,
        'fuelType': request.fuelType,
        'selfMoveCapable': request.selfMoveCapable,
        'dailyRate': {
          'amountMinor': request.dailyRateMinor,
          'currency': 'ETB',
        },
        'excessHourRate': {
          'amountMinor': request.excessHourRateMinor,
          'currency': 'ETB',
        },
        'minRentalDays': request.minRentalDays,
        'paymentTimeoutHours': request.paymentTimeoutHours,
        'operatorIncluded': request.operatorIncluded,
        'location': {
          'latitude': request.latitude,
          'longitude': request.longitude,
          'addressText': request.addressText,
        },
        'photos': request.photos
            .map((p) => {
                  'angle': p.angle.code,
                  'storageUri': p.storageUri,
                  'fileHash': p.fileHash,
                })
            .toList(),
        'specValues': request.specValues.map((v) => v.toJson()).toList(),
      };

      final response = await apiClient.dio.post(
        '/fleets',
        data: payload,
      );

      final data = response.data['data'] as Map<String, dynamic>;
      return FleetAssetModel.fromJson(data);
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }
}
