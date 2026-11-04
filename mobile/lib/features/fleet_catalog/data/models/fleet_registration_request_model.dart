import '../../domain/entities/fleet_registration_request.dart';
import 'fleet_photo_model.dart';
import 'fleet_spec_value_input_model.dart';

class FleetRegistrationRequestModel {
  final String categoryId;
  final String make;
  final String model;
  final int manufacturingYear;
  final String registrationNumber;
  final String chassisVin;
  final double baselineMeterHours;
  final String fuelType;
  final bool selfMoveCapable;
  final int dailyRateMinor;
  final int excessHourRateMinor;
  final int minRentalDays;
  final int paymentTimeoutHours;
  final bool operatorIncluded;
  final double latitude;
  final double longitude;
  final String addressText;
  final List<FleetPhotoModel> photos;
  final List<FleetSpecValueInputModel> specValues;

  const FleetRegistrationRequestModel({
    required this.categoryId,
    required this.make,
    required this.model,
    required this.manufacturingYear,
    required this.registrationNumber,
    required this.chassisVin,
    required this.baselineMeterHours,
    required this.fuelType,
    required this.selfMoveCapable,
    required this.dailyRateMinor,
    this.excessHourRateMinor = 0,
    this.minRentalDays = 1,
    this.paymentTimeoutHours = 24,
    this.operatorIncluded = false,
    required this.latitude,
    required this.longitude,
    required this.addressText,
    required this.photos,
    required this.specValues,
  });

  factory FleetRegistrationRequestModel.fromJson(Map<String, dynamic> json) {
    final photosRaw = json['photos'] as List<dynamic>? ?? [];
    final specsRaw = json['specValues'] as List<dynamic>? ?? [];
    return FleetRegistrationRequestModel(
      categoryId: json['categoryId'] as String? ?? '',
      make: json['make'] as String? ?? '',
      model: json['model'] as String? ?? '',
      manufacturingYear: (json['manufacturingYear'] as num?)?.toInt() ?? 0,
      registrationNumber: json['registrationNumber'] as String? ?? '',
      chassisVin: json['chassisVin'] as String? ?? '',
      baselineMeterHours: (json['baselineMeterHours'] as num?)?.toDouble() ?? 0.0,
      fuelType: json['fuelType'] as String? ?? 'DIESEL',
      selfMoveCapable: json['selfMoveCapable'] as bool? ?? false,
      dailyRateMinor: (json['dailyRateMinor'] as num?)?.toInt() ?? 0,
      excessHourRateMinor: (json['excessHourRateMinor'] as num?)?.toInt() ?? 0,
      minRentalDays: (json['minRentalDays'] as num?)?.toInt() ?? 1,
      paymentTimeoutHours: (json['paymentTimeoutHours'] as num?)?.toInt() ?? 24,
      operatorIncluded: json['operatorIncluded'] as bool? ?? false,
      latitude: (json['latitude'] as num?)?.toDouble() ?? 0.0,
      longitude: (json['longitude'] as num?)?.toDouble() ?? 0.0,
      addressText: json['addressText'] as String? ?? '',
      photos: photosRaw.map((p) => FleetPhotoModel.fromJson(p as Map<String, dynamic>)).toList(),
      specValues: specsRaw.map((s) => FleetSpecValueInputModel.fromJson(s as Map<String, dynamic>)).toList(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'categoryId': categoryId,
      'make': make,
      'model': model,
      'manufacturingYear': manufacturingYear,
      'registrationNumber': registrationNumber,
      'chassisVin': chassisVin,
      'baselineMeterHours': baselineMeterHours,
      'fuelType': fuelType,
      'selfMoveCapable': selfMoveCapable,
      'dailyRateMinor': dailyRateMinor,
      'excessHourRateMinor': excessHourRateMinor,
      'minRentalDays': minRentalDays,
      'paymentTimeoutHours': paymentTimeoutHours,
      'operatorIncluded': operatorIncluded,
      'latitude': latitude,
      'longitude': longitude,
      'addressText': addressText,
      'photos': photos.map((p) => p.toJson()).toList(),
      'specValues': specValues.map((v) => v.toJson()).toList(),
    };
  }

  factory FleetRegistrationRequestModel.fromEntity(FleetRegistrationRequest entity) {
    return FleetRegistrationRequestModel(
      categoryId: entity.categoryId,
      make: entity.make,
      model: entity.model,
      manufacturingYear: entity.manufacturingYear,
      registrationNumber: entity.registrationNumber,
      chassisVin: entity.chassisVin,
      baselineMeterHours: entity.baselineMeterHours,
      fuelType: entity.fuelType,
      selfMoveCapable: entity.selfMoveCapable,
      dailyRateMinor: entity.dailyRateMinor,
      excessHourRateMinor: entity.excessHourRateMinor,
      minRentalDays: entity.minRentalDays,
      paymentTimeoutHours: entity.paymentTimeoutHours,
      operatorIncluded: entity.operatorIncluded,
      latitude: entity.latitude,
      longitude: entity.longitude,
      addressText: entity.addressText,
      photos: entity.photos.map(FleetPhotoModel.fromEntity).toList(),
      specValues: entity.specValues.map(FleetSpecValueInputModel.fromEntity).toList(),
    );
  }
}
