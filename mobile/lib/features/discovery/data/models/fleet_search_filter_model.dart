import 'package:intl/intl.dart';
import '../../domain/entities/discovery_sort_by.dart';
import '../../domain/entities/fleet_search_filter.dart';

class FleetSearchFilterModel {
  final double latitude;
  final double longitude;
  final double radiusKm;
  final String? categoryId;
  final String? subcategoryId;
  final DateTime? startDate;
  final DateTime? endDate;
  final int? minRateMinor;
  final int? maxRateMinor;
  final bool? operatorIncluded;
  final String? fuelType;
  final bool? selfMoveCapable;
  final String? query;
  final String sortBy;

  const FleetSearchFilterModel({
    required this.latitude,
    required this.longitude,
    this.radiusKm = 50.0,
    this.categoryId,
    this.subcategoryId,
    this.startDate,
    this.endDate,
    this.minRateMinor,
    this.maxRateMinor,
    this.operatorIncluded,
    this.fuelType,
    this.selfMoveCapable,
    this.query,
    this.sortBy = 'RECOMMENDED',
  });

  factory FleetSearchFilterModel.fromEntity(FleetSearchFilter filter) {
    return FleetSearchFilterModel(
      latitude: filter.latitude,
      longitude: filter.longitude,
      radiusKm: filter.radiusKm,
      categoryId: filter.categoryId,
      subcategoryId: filter.subcategoryId,
      startDate: filter.startDate,
      endDate: filter.endDate,
      minRateMinor: filter.minRateMinor,
      maxRateMinor: filter.maxRateMinor,
      operatorIncluded: filter.operatorIncluded,
      fuelType: filter.fuelType,
      selfMoveCapable: filter.selfMoveCapable,
      query: filter.query,
      sortBy: filter.sortBy.toParam(),
    );
  }

  Map<String, dynamic> toQueryParams() {
    final dateFormat = DateFormat('yyyy-MM-dd');
    final queryParams = <String, dynamic>{
      'latitude': latitude,
      'longitude': longitude,
      'radiusKm': radiusKm.toInt(),
      'sortBy': sortBy,
    };

    if (categoryId != null) queryParams['categoryId'] = categoryId;
    if (subcategoryId != null) queryParams['subcategoryId'] = subcategoryId;
    if (startDate != null) queryParams['startDate'] = dateFormat.format(startDate!);
    if (endDate != null) queryParams['endDate'] = dateFormat.format(endDate!);
    if (minRateMinor != null) queryParams['minRateMinor'] = minRateMinor;
    if (maxRateMinor != null) queryParams['maxRateMinor'] = maxRateMinor;
    if (operatorIncluded != null) queryParams['operatorIncluded'] = operatorIncluded;
    if (fuelType != null) queryParams['fuelType'] = fuelType;
    if (selfMoveCapable != null) queryParams['selfMoveCapable'] = selfMoveCapable;
    if (query != null && query!.isNotEmpty) queryParams['query'] = query;

    return queryParams;
  }
}
