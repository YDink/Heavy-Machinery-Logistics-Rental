import 'package:equatable/equatable.dart';

export 'fleet_photo.dart';
export 'fleet_registration_request.dart';
export 'fleet_spec_value_input.dart';
export 'photo_angle.dart';

class FleetAsset extends Equatable {
  final String id;
  final String ownerId;
  final String categoryId;
  final String make;
  final String model;
  final int manufacturingYear;
  final String registrationNumber;
  final int dailyRateMinor;
  final String status;
  final double latitude;
  final double longitude;
  final String addressText;
  final int version;

  const FleetAsset({
    required this.id,
    required this.ownerId,
    required this.categoryId,
    required this.make,
    required this.model,
    required this.manufacturingYear,
    required this.registrationNumber,
    required this.dailyRateMinor,
    required this.status,
    required this.latitude,
    required this.longitude,
    required this.addressText,
    required this.version,
  });

  @override
  List<Object?> get props => [
        id,
        ownerId,
        categoryId,
        make,
        model,
        manufacturingYear,
        registrationNumber,
        dailyRateMinor,
        status,
        latitude,
        longitude,
        addressText,
        version,
      ];
}
