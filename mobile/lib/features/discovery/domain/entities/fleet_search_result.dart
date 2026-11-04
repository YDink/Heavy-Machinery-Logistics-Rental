import 'package:equatable/equatable.dart';
import 'fleet_availability_status.dart';
import 'fleet_search_key_specs.dart';
import 'fleet_search_owner_summary.dart';

export 'fleet_availability_status.dart';
export 'fleet_search_data.dart';
export 'fleet_search_key_specs.dart';
export 'fleet_search_owner_summary.dart';

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
