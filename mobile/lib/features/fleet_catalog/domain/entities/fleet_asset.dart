import 'package:equatable/equatable.dart';

enum PhotoAngle {
  front,
  rear,
  left,
  right,
  cabin,
  meter;

  String get code => name.toUpperCase();

  static PhotoAngle? fromCode(String code) {
    for (final a in PhotoAngle.values) {
      if (a.code == code.toUpperCase()) return a;
    }
    return null;
  }
}

class FleetPhoto extends Equatable {
  final PhotoAngle angle;
  final String storageUri;
  final String fileHash;
  final String? localFilePath;

  const FleetPhoto({
    required this.angle,
    required this.storageUri,
    required this.fileHash,
    this.localFilePath,
  });

  @override
  List<Object?> get props => [angle, storageUri, fileHash, localFilePath];
}

class FleetSpecValueInput extends Equatable {
  final String fieldKey;
  final String? valueString;
  final double? valueNumeric;
  final bool? valueBoolean;

  const FleetSpecValueInput({
    required this.fieldKey,
    this.valueString,
    this.valueNumeric,
    this.valueBoolean,
  });

  @override
  List<Object?> get props => [fieldKey, valueString, valueNumeric, valueBoolean];
}

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
