import 'package:dink_mobile/features/discovery/domain/entities/fleet_compare_matrix.dart';

class FleetCompareDataModel extends FleetCompareData {
  const FleetCompareDataModel({
    required super.subcategoryId,
    required super.subcategoryNameEn,
    required super.subcategoryNameAm,
    required super.fleets,
    required super.comparisonMatrix,
  });

  factory FleetCompareDataModel.fromJson(Map<String, dynamic> json) {
    final rawFleets = json['fleets'] as List<dynamic>? ?? [];
    final fleets = rawFleets.map((e) {
      final m = e as Map<String, dynamic>;
      return FleetCompareItemHeader(
        id: m['id'] as String? ?? '',
        title: m['title'] as String? ?? '',
        dailyRateMinor: (m['dailyRateMinor'] as num?)?.toInt() ?? 0,
        excessHourRateMinor: (m['excessHourRateMinor'] as num?)?.toInt() ?? 0,
        primaryPhotoUri: m['primaryPhotoUri'] as String?,
        selfMoveCapable: m['selfMoveCapable'] as bool? ?? false,
        operatorIncluded: m['operatorIncluded'] as bool? ?? false,
        fuelType: m['fuelType'] as String? ?? 'DIESEL',
      );
    }).toList();

    final rawMatrix = json['comparisonMatrix'] as List<dynamic>? ?? [];
    final matrix = rawMatrix.map((e) {
      final m = e as Map<String, dynamic>;
      final rawValues = m['values'] as Map<String, dynamic>? ?? {};
      final values = rawValues.map((k, v) => MapEntry(k, v?.toString() ?? ''));
      return FleetCompareSpecRow(
        fieldKey: m['fieldKey'] as String? ?? '',
        labelEn: m['labelEn'] as String? ?? '',
        labelAm: m['labelAm'] as String? ?? '',
        unitSymbol: m['unitSymbol'] as String?,
        fieldType: m['fieldType'] as String? ?? 'STRING',
        values: values,
        isDifferent: m['isDifferent'] as bool? ?? false,
      );
    }).toList();

    return FleetCompareDataModel(
      subcategoryId: json['subcategoryId'] as String? ?? '',
      subcategoryNameEn: json['subcategoryNameEn'] as String? ?? '',
      subcategoryNameAm: json['subcategoryNameAm'] as String? ?? '',
      fleets: fleets,
      comparisonMatrix: matrix,
    );
  }
}
