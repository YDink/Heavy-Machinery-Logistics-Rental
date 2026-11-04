import 'package:equatable/equatable.dart';
import 'spec_field_type.dart';

class SpecFieldDefinition extends Equatable {
  final String id;
  final String fieldKey;
  final String labelEN;
  final String labelAM;
  final SpecFieldType fieldType;
  final String? unitSymbol;
  final bool isRequired;
  final bool isLogisticsCritical;
  final double? minValue;
  final double? maxValue;
  final List<String> enumValues;

  const SpecFieldDefinition({
    required this.id,
    required this.fieldKey,
    required this.labelEN,
    required this.labelAM,
    required this.fieldType,
    this.unitSymbol,
    required this.isRequired,
    this.isLogisticsCritical = false,
    this.minValue,
    this.maxValue,
    this.enumValues = const [],
  });

  String getLocalizedLabel(String localeCode) {
    if (localeCode.startsWith('am') && labelAM.isNotEmpty) {
      return labelAM;
    }
    return labelEN;
  }

  @override
  List<Object?> get props => [
        id,
        fieldKey,
        labelEN,
        labelAM,
        fieldType,
        unitSymbol,
        isRequired,
        isLogisticsCritical,
        minValue,
        maxValue,
        enumValues,
      ];
}
