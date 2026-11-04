import 'package:equatable/equatable.dart';

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
