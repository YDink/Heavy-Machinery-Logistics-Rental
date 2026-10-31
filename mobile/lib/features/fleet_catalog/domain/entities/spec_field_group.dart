import 'package:equatable/equatable.dart';

enum SpecFieldType {
  string,
  decimal,
  integer,
  boolean,
  enumeration;

  static SpecFieldType fromString(String val) {
    switch (val.toUpperCase()) {
      case 'DECIMAL':
        return SpecFieldType.decimal;
      case 'INTEGER':
        return SpecFieldType.integer;
      case 'BOOLEAN':
        return SpecFieldType.boolean;
      case 'ENUM':
        return SpecFieldType.enumeration;
      default:
        return SpecFieldType.string;
    }
  }
}

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

class SpecFieldGroup extends Equatable {
  final String id;
  final String categoryId;
  final String sectionNameEN;
  final String sectionNameAM;
  final int displayOrder;
  final List<SpecFieldDefinition> fields;

  const SpecFieldGroup({
    required this.id,
    required this.categoryId,
    required this.sectionNameEN,
    required this.sectionNameAM,
    required this.displayOrder,
    required this.fields,
  });

  String getLocalizedSectionName(String localeCode) {
    if (localeCode.startsWith('am') && sectionNameAM.isNotEmpty) {
      return sectionNameAM;
    }
    return sectionNameEN;
  }

  @override
  List<Object?> get props => [
        id,
        categoryId,
        sectionNameEN,
        sectionNameAM,
        displayOrder,
        fields,
      ];
}
