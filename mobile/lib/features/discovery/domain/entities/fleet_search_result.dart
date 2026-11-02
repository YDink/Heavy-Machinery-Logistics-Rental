import 'package:equatable/equatable.dart';

enum FleetAvailabilityStatus {
  availableNow,
  heldPendingWire,
  unavailable,
}

class FleetSearchOwnerSummary extends Equatable {
  final String ownerId;
  final String tradingName;
  final double averageRating;
  final int totalCompletedRentals;
  final bool isVerified;

  const FleetSearchOwnerSummary({
    required this.ownerId,
    required this.tradingName,
    required this.averageRating,
    required this.totalCompletedRentals,
    required this.isVerified,
  });

  @override
  List<Object?> get props => [
        ownerId,
        tradingName,
        averageRating,
        totalCompletedRentals,
        isVerified,
      ];
}

class FleetSearchKeySpecs extends Equatable {
  final double? operatingWeightKg;
  final double? payloadCapacityKg;
  final double? grossPowerHp;

  const FleetSearchKeySpecs({
    this.operatingWeightKg,
    this.payloadCapacityKg,
    this.grossPowerHp,
  });

  @override
  List<Object?> get props => [
        operatingWeightKg,
        payloadCapacityKg,
        grossPowerHp,
      ];
}

class FleetSearchResult extends Equatable {
  final String fleetId;
  final String subcategoryId;
  final String subcategoryNameEn;
  final String subcategoryNameAm;
  final String make;
  final String model;
  final int manufactureYear;
  final int dailyRateMinor;
  final int excessHourRateMinor;
  final double distanceMeters;
  final FleetAvailabilityStatus availabilityStatus;
  final bool selfMoveCapable;
  final bool operatorIncluded;
  final String fuelType;
  final String? primaryPhotoUri;
  final FleetSearchOwnerSummary owner;
  final FleetSearchKeySpecs keySpecs;

  const FleetSearchResult({
    required this.fleetId,
    required this.subcategoryId,
    required this.subcategoryNameEn,
    required this.subcategoryNameAm,
    required this.make,
    required this.model,
    required this.manufactureYear,
    required this.dailyRateMinor,
    required this.excessHourRateMinor,
    required this.distanceMeters,
    required this.availabilityStatus,
    required this.selfMoveCapable,
    required this.operatorIncluded,
    required this.fuelType,
    this.primaryPhotoUri,
    required this.owner,
    required this.keySpecs,
  });

  @override
  List<Object?> get props => [
        fleetId,
        subcategoryId,
        subcategoryNameEn,
        subcategoryNameAm,
        make,
        model,
        manufactureYear,
        dailyRateMinor,
        excessHourRateMinor,
        distanceMeters,
        availabilityStatus,
        selfMoveCapable,
        operatorIncluded,
        fuelType,
        primaryPhotoUri,
        owner,
        keySpecs,
      ];
}

class FleetSearchData extends Equatable {
  final int totalMatches;
  final double latitude;
  final double longitude;
  final double radiusKm;
  final List<FleetSearchResult> results;

  const FleetSearchData({
    required this.totalMatches,
    required this.latitude,
    required this.longitude,
    required this.radiusKm,
    required this.results,
  });

  @override
  List<Object?> get props => [totalMatches, latitude, longitude, radiusKm, results];
}
