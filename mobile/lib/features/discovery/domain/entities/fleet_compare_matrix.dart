import 'package:equatable/equatable.dart';

class FleetCompareItemHeader extends Equatable {
  final String id;
  final String title;
  final int dailyRateMinor;
  final int excessHourRateMinor;
  final String? primaryPhotoUri;
  final bool selfMoveCapable;
  final bool operatorIncluded;
  final String fuelType;

  const FleetCompareItemHeader({
    required this.id,
    required this.title,
    required this.dailyRateMinor,
    required this.excessHourRateMinor,
    this.primaryPhotoUri,
    required this.selfMoveCapable,
    required this.operatorIncluded,
    required this.fuelType,
  });

  @override
  List<Object?> get props => [
        id,
        title,
        dailyRateMinor,
        excessHourRateMinor,
        primaryPhotoUri,
        selfMoveCapable,
        operatorIncluded,
        fuelType,
      ];
}

class FleetCompareSpecRow extends Equatable {
  final String fieldKey;
  final String labelEn;
  final String labelAm;
  final String? unitSymbol;
  final String fieldType;
  final Map<String, String> values; // fleetId -> value
  final bool isDifferent;

  const FleetCompareSpecRow({
    required this.fieldKey,
    required this.labelEn,
    required this.labelAm,
    this.unitSymbol,
    required this.fieldType,
    required this.values,
    required this.isDifferent,
  });

  @override
  List<Object?> get props => [
        fieldKey,
        labelEn,
        labelAm,
        unitSymbol,
        fieldType,
        values,
        isDifferent,
      ];
}

class FleetCompareData extends Equatable {
  final String subcategoryId;
  final String subcategoryNameEn;
  final String subcategoryNameAm;
  final List<FleetCompareItemHeader> fleets;
  final List<FleetCompareSpecRow> comparisonMatrix;

  const FleetCompareData({
    required this.subcategoryId,
    required this.subcategoryNameEn,
    required this.subcategoryNameAm,
    required this.fleets,
    required this.comparisonMatrix,
  });

  @override
  List<Object?> get props => [
        subcategoryId,
        subcategoryNameEn,
        subcategoryNameAm,
        fleets,
        comparisonMatrix,
      ];
}
