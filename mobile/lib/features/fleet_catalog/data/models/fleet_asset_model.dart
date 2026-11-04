import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';

class FleetAssetModel extends FleetAsset {
  const FleetAssetModel({
    required super.id,
    required super.ownerId,
    required super.categoryId,
    required super.make,
    required super.model,
    required super.manufacturingYear,
    required super.registrationNumber,
    required super.dailyRateMinor,
    required super.status,
    required super.latitude,
    required super.longitude,
    required super.addressText,
    required super.version,
  });

  factory FleetAssetModel.fromJson(Map<String, dynamic> json) {
    final location = (json['location'] as Map<String, dynamic>?) ?? {};
    int dailyRate = 0;
    if (json['dailyRate'] is Map<String, dynamic>) {
      dailyRate = (json['dailyRate']['amountMinor'] as num?)?.toInt() ?? 0;
    } else if (json['dailyRate'] is num) {
      dailyRate = (json['dailyRate'] as num).toInt();
    }

    return FleetAssetModel(
      id: json['id'] as String,
      ownerId: json['ownerId'] as String? ?? '',
      categoryId: json['categoryId'] as String,
      make: json['make'] as String? ?? '',
      model: json['model'] as String? ?? '',
      manufacturingYear: json['manufacturingYear'] as int? ?? 2020,
      registrationNumber: json['registrationNumber'] as String? ?? '',
      dailyRateMinor: dailyRate,
      status: json['status'] as String? ?? 'PENDING_APPROVAL',
      latitude: (location['latitude'] as num?)?.toDouble() ?? 0.0,
      longitude: (location['longitude'] as num?)?.toDouble() ?? 0.0,
      addressText: location['addressText'] as String? ?? '',
      version: json['version'] as int? ?? 1,
    );
  }
}
