import 'package:equatable/equatable.dart';

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
