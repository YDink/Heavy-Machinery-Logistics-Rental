import 'package:equatable/equatable.dart';
import 'fleet_compare_item_header.dart';
import 'fleet_compare_spec_row.dart';

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
