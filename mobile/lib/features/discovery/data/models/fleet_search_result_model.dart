import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_result.dart';

class FleetSearchResultModel extends FleetSearchResult {
  const FleetSearchResultModel({
    required super.fleetId,
    required super.subcategoryId,
    required super.subcategoryNameEn,
    required super.subcategoryNameAm,
    required super.make,
    required super.model,
    required super.manufactureYear,
    required super.dailyRateMinor,
    required super.excessHourRateMinor,
    required super.distanceMeters,
    required super.availabilityStatus,
    required super.selfMoveCapable,
    required super.operatorIncluded,
    required super.fuelType,
    super.primaryPhotoUri,
    required super.owner,
    required super.keySpecs,
  });

  factory FleetSearchResultModel.fromJson(Map<String, dynamic> json) {
    FleetAvailabilityStatus status;
    final statusStr = json['availabilityStatus'] as String? ?? 'AVAILABLE_NOW';
    if (statusStr == 'HELD_PENDING_WIRE') {
      status = FleetAvailabilityStatus.heldPendingWire;
    } else if (statusStr == 'UNAVAILABLE') {
      status = FleetAvailabilityStatus.unavailable;
    } else {
      status = FleetAvailabilityStatus.availableNow;
    }

    final ownerMap = json['owner'] as Map<String, dynamic>? ?? {};
    final owner = FleetSearchOwnerSummary(
      ownerId: ownerMap['ownerId'] as String? ?? '',
      tradingName: ownerMap['tradingName'] as String? ?? '',
      averageRating: (ownerMap['averageRating'] as num?)?.toDouble() ?? 5.0,
      totalCompletedRentals: (ownerMap['totalCompletedRentals'] as num?)?.toInt() ?? 0,
      isVerified: ownerMap['isVerified'] as bool? ?? false,
    );

    final specsMap = json['keySpecs'] as Map<String, dynamic>? ?? {};
    final specs = FleetSearchKeySpecs(
      operatingWeightKg: (specsMap['operatingWeightKg'] as num?)?.toDouble(),
      payloadCapacityKg: (specsMap['payloadCapacityKg'] as num?)?.toDouble(),
      grossPowerHp: (specsMap['grossPowerHp'] as num?)?.toDouble(),
    );

    return FleetSearchResultModel(
      fleetId: json['fleetId'] as String? ?? '',
      subcategoryId: json['subcategoryId'] as String? ?? '',
      subcategoryNameEn: json['subcategoryNameEn'] as String? ?? '',
      subcategoryNameAm: json['subcategoryNameAm'] as String? ?? '',
      make: json['make'] as String? ?? '',
      model: json['model'] as String? ?? '',
      manufactureYear: (json['manufactureYear'] as num?)?.toInt() ?? 2024,
      dailyRateMinor: (json['dailyRateMinor'] as num?)?.toInt() ?? 0,
      excessHourRateMinor: (json['excessHourRateMinor'] as num?)?.toInt() ?? 0,
      distanceMeters: (json['distanceMeters'] as num?)?.toDouble() ?? 0.0,
      availabilityStatus: status,
      selfMoveCapable: json['selfMoveCapable'] as bool? ?? false,
      operatorIncluded: json['operatorIncluded'] as bool? ?? false,
      fuelType: json['fuelType'] as String? ?? 'DIESEL',
      primaryPhotoUri: json['primaryPhotoUri'] as String?,
      owner: owner,
      keySpecs: specs,
    );
  }
}

class FleetSearchDataModel extends FleetSearchData {
  const FleetSearchDataModel({
    required super.totalMatches,
    required super.latitude,
    required super.longitude,
    required super.radiusKm,
    required super.results,
  });

  factory FleetSearchDataModel.fromJson(Map<String, dynamic> json) {
    final center = json['searchCenter'] as Map<String, dynamic>? ?? {};
    final rawResults = json['results'] as List<dynamic>? ?? [];
    final items = rawResults
        .map((e) => FleetSearchResultModel.fromJson(e as Map<String, dynamic>))
        .toList();

    return FleetSearchDataModel(
      totalMatches: (json['totalMatches'] as num?)?.toInt() ?? items.length,
      latitude: (center['latitude'] as num?)?.toDouble() ?? 0.0,
      longitude: (center['longitude'] as num?)?.toDouble() ?? 0.0,
      radiusKm: (center['radiusKm'] as num?)?.toDouble() ?? 50.0,
      results: items,
    );
  }
}
