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
