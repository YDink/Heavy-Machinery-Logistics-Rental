import 'package:equatable/equatable.dart';

enum DiscoverySortBy {
  recommended,
  priceAsc,
  priceDesc,
  distanceAsc,
  ratingDesc,
}

extension DiscoverySortByExtension on DiscoverySortBy {
  String toParam() {
    switch (this) {
      case DiscoverySortBy.priceAsc:
        return 'PRICE_ASC';
      case DiscoverySortBy.priceDesc:
        return 'PRICE_DESC';
      case DiscoverySortBy.distanceAsc:
        return 'DISTANCE_ASC';
      case DiscoverySortBy.ratingDesc:
        return 'RATING_DESC';
      case DiscoverySortBy.recommended:
        return 'RECOMMENDED';
    }
  }
}

class FleetSearchFilter extends Equatable {
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
  final DiscoverySortBy sortBy;

  const FleetSearchFilter({
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
    this.sortBy = DiscoverySortBy.recommended,
  });

  FleetSearchFilter copyWith({
    double? latitude,
    double? longitude,
    double? radiusKm,
    String? categoryId,
    String? subcategoryId,
    DateTime? startDate,
    DateTime? endDate,
    int? minRateMinor,
    int? maxRateMinor,
    bool? operatorIncluded,
    String? fuelType,
    bool? selfMoveCapable,
    String? query,
    DiscoverySortBy? sortBy,
  }) {
    return FleetSearchFilter(
      latitude: latitude ?? this.latitude,
      longitude: longitude ?? this.longitude,
      radiusKm: radiusKm ?? this.radiusKm,
      categoryId: categoryId ?? this.categoryId,
      subcategoryId: subcategoryId ?? this.subcategoryId,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
      minRateMinor: minRateMinor ?? this.minRateMinor,
      maxRateMinor: maxRateMinor ?? this.maxRateMinor,
      operatorIncluded: operatorIncluded ?? this.operatorIncluded,
      fuelType: fuelType ?? this.fuelType,
      selfMoveCapable: selfMoveCapable ?? this.selfMoveCapable,
      query: query ?? this.query,
      sortBy: sortBy ?? this.sortBy,
    );
  }

  @override
  List<Object?> get props => [
        latitude,
        longitude,
        radiusKm,
        categoryId,
        subcategoryId,
        startDate,
        endDate,
        minRateMinor,
        maxRateMinor,
        operatorIncluded,
        fuelType,
        selfMoveCapable,
        query,
        sortBy,
      ];
}
