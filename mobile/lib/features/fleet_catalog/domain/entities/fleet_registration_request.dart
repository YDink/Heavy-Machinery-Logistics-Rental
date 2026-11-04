import 'package:equatable/equatable.dart';
import 'fleet_photo.dart';
import 'fleet_spec_value_input.dart';
import 'photo_angle.dart';

class FleetRegistrationRequest extends Equatable {
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
  final List<FleetPhoto> photos;
  final List<FleetSpecValueInput> specValues;

  const FleetRegistrationRequest({
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

  bool get hasCompleteWalkaroundPhotos {
    final set = photos.map((p) => p.angle).toSet();
    return set.containsAll(PhotoAngle.values);
  }

  @override
  List<Object?> get props => [
        categoryId,
        make,
        model,
        manufacturingYear,
        registrationNumber,
        chassisVin,
        baselineMeterHours,
        fuelType,
        selfMoveCapable,
        dailyRateMinor,
        excessHourRateMinor,
        minRentalDays,
        paymentTimeoutHours,
        operatorIncluded,
        latitude,
        longitude,
        addressText,
        photos,
        specValues,
      ];
}
