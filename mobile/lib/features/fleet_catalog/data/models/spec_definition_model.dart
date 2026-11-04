import 'package:dink_mobile/features/fleet_catalog/domain/entities/spec_field_group.dart';

class SpecFieldDefinitionModel extends SpecFieldDefinition {
  const SpecFieldDefinitionModel({
    required super.id,
    required super.fieldKey,
    required super.labelEN,
    required super.labelAM,
    required super.fieldType,
    super.unitSymbol,
    required super.isRequired,
    super.isLogisticsCritical,
    super.minValue,
    super.maxValue,
    super.enumValues,
  });

  factory SpecFieldDefinitionModel.fromJson(Map<String, dynamic> json) {
    final enumList = (json['enumValues'] as List<dynamic>?)
            ?.map((e) => e.toString())
            .toList() ??
        const [];

    return SpecFieldDefinitionModel(
      id: json['id'] as String? ?? json['definitionId'] as String,
      fieldKey: json['fieldKey'] as String,
      labelEN: json['labelEn'] as String? ?? '',
      labelAM: json['labelAm'] as String? ?? '',
      fieldType: SpecFieldType.fromString(json['fieldType'] as String? ?? 'STRING'),
      unitSymbol: json['unitSymbol'] as String?,
      isRequired: json['isRequired'] as bool? ?? false,
      isLogisticsCritical: json['isLogisticsCritical'] as bool? ?? false,
      minValue: (json['minValue'] as num?)?.toDouble(),
      maxValue: (json['maxValue'] as num?)?.toDouble(),
      enumValues: enumList,
    );
  }
}

class SpecFieldGroupModel extends SpecFieldGroup {
  const SpecFieldGroupModel({
    required super.id,
    required super.categoryId,
    required super.sectionNameEN,
    required super.sectionNameAM,
    required super.displayOrder,
    required super.fields,
  });

  factory SpecFieldGroupModel.fromJson(Map<String, dynamic> json) {
    final fieldsRaw = (json['fields'] as List<dynamic>?) ?? [];
    final fieldsList = fieldsRaw
        .map((f) => SpecFieldDefinitionModel.fromJson(f as Map<String, dynamic>))
        .toList();

    return SpecFieldGroupModel(
      id: json['id'] as String? ?? json['groupId'] as String,
      categoryId: json['categoryId'] as String? ?? '',
      sectionNameEN: json['sectionNameEn'] as String? ?? '',
      sectionNameAM: json['sectionNameAm'] as String? ?? '',
      displayOrder: json['displayOrder'] as int? ?? 1,
      fields: fieldsList,
    );
  }
}
