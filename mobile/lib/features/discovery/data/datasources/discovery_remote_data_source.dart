import 'package:dio/dio.dart';
import 'package:intl/intl.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import '../models/fleet_availability_model.dart';
import '../models/fleet_compare_model.dart';
import '../models/fleet_search_filter_model.dart';
import '../models/fleet_search_result_model.dart';

abstract class DiscoveryRemoteDataSource {
  Future<FleetSearchDataModel> searchFleets(FleetSearchFilterModel filter);
  Future<FleetCompareDataModel> compareFleets(List<String> fleetIds);
  Future<FleetAvailabilityDataModel> getFleetAvailability(String fleetId, DateTime startDate, DateTime endDate);
}

class DiscoveryRemoteDataSourceImpl implements DiscoveryRemoteDataSource {
  final ApiClient apiClient;

  DiscoveryRemoteDataSourceImpl({required this.apiClient});

  final DateFormat _dateFormat = DateFormat('yyyy-MM-dd');

  @override
  Future<FleetSearchDataModel> searchFleets(FleetSearchFilterModel filter) async {
    try {
      final response = await apiClient.dio.get(
        '/fleets/search',
        queryParameters: filter.toQueryParams(),
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
