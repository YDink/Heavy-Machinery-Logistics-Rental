import '../../domain/entities/fleet_spec_value_input.dart';

class FleetSpecValueInputModel extends FleetSpecValueInput {
  const FleetSpecValueInputModel({
    required super.fieldKey,
    super.valueString,
    super.valueNumeric,
    super.valueBoolean,
  });

  factory FleetSpecValueInputModel.fromJson(Map<String, dynamic> json) {
    return FleetSpecValueInputModel(
      fieldKey: json['fieldKey'] as String? ?? '',
      valueString: json['valueString'] as String?,
      valueNumeric: (json['valueNumeric'] as num?)?.toDouble(),
      valueBoolean: json['valueBoolean'] as bool?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'fieldKey': fieldKey,
      if (valueString != null) 'valueString': valueString,
      if (valueNumeric != null) 'valueNumeric': valueNumeric,
      if (valueBoolean != null) 'valueBoolean': valueBoolean,
    };
  }

  factory FleetSpecValueInputModel.fromEntity(FleetSpecValueInput entity) {
    return FleetSpecValueInputModel(
      fieldKey: entity.fieldKey,
      valueString: entity.valueString,
      valueNumeric: entity.valueNumeric,
      valueBoolean: entity.valueBoolean,
    );
  }
}
