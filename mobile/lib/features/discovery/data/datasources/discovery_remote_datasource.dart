import 'package:dio/dio.dart';
import 'package:intl/intl.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import 'package:dink_mobile/features/discovery/data/models/fleet_availability_model.dart';
import 'package:dink_mobile/features/discovery/data/models/fleet_compare_model.dart';
import 'package:dink_mobile/features/discovery/data/models/fleet_search_result_model.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';

abstract class DiscoveryRemoteDataSource {
  Future<FleetSearchDataModel> searchFleets(FleetSearchFilter filter);
  Future<FleetCompareDataModel> compareFleets(List<String> fleetIds);
  Future<FleetAvailabilityDataModel> getFleetAvailability(String fleetId, DateTime startDate, DateTime endDate);
}

class DiscoveryRemoteDataSourceImpl implements DiscoveryRemoteDataSource {
  final ApiClient apiClient;

  DiscoveryRemoteDataSourceImpl({required this.apiClient});

  final DateFormat _dateFormat = DateFormat('yyyy-MM-dd');

  @override
  Future<FleetSearchDataModel> searchFleets(FleetSearchFilter filter) async {
    final queryParams = <String, dynamic>{
      'latitude': filter.latitude,
      'longitude': filter.longitude,
      'radiusKm': filter.radiusKm.toInt(),
      'sortBy': filter.sortBy.toParam(),
    };

    if (filter.categoryId != null) queryParams['categoryId'] = filter.categoryId;
    if (filter.subcategoryId != null) queryParams['subcategoryId'] = filter.subcategoryId;
    if (filter.startDate != null) queryParams['startDate'] = _dateFormat.format(filter.startDate!);
    if (filter.endDate != null) queryParams['endDate'] = _dateFormat.format(filter.endDate!);
    if (filter.minRateMinor != null) queryParams['minRateMinor'] = filter.minRateMinor;
    if (filter.maxRateMinor != null) queryParams['maxRateMinor'] = filter.maxRateMinor;
    if (filter.operatorIncluded != null) queryParams['operatorIncluded'] = filter.operatorIncluded;
    if (filter.fuelType != null) queryParams['fuelType'] = filter.fuelType;
    if (filter.selfMoveCapable != null) queryParams['selfMoveCapable'] = filter.selfMoveCapable;
    if (filter.query != null && filter.query!.isNotEmpty) queryParams['query'] = filter.query;

    try {
      final response = await apiClient.dio.get(
        '/fleets/search',
        queryParameters: queryParams,
      );
      final body = response.data as Map<String, dynamic>;
      final data = body['data'] as Map<String, dynamic>;
      return FleetSearchDataModel.fromJson(data);
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }

  @override
  Future<FleetCompareDataModel> compareFleets(List<String> fleetIds) async {
    try {
      final response = await apiClient.dio.post(
        '/fleets/compare',
        data: {'fleetIds': fleetIds},
      );
      final body = response.data as Map<String, dynamic>;
      final data = body['data'] as Map<String, dynamic>;
      return FleetCompareDataModel.fromJson(data);
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }

  @override
  Future<FleetAvailabilityDataModel> getFleetAvailability(
    String fleetId,
    DateTime startDate,
    DateTime endDate,
  ) async {
    try {
      final response = await apiClient.dio.get(
        '/fleets/$fleetId/availability',
        queryParameters: {
          'startDate': _dateFormat.format(startDate),
          'endDate': _dateFormat.format(endDate),
        },
      );
      final body = response.data as Map<String, dynamic>;
      final data = body['data'] as Map<String, dynamic>;
      return FleetAvailabilityDataModel.fromJson(data);
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }
}
