import 'package:equatable/equatable.dart';

class FleetSpecValueInput extends Equatable {
  final String fieldKey;
  final String? valueString;
  final double? valueNumeric;
  final bool? valueBoolean;

  const FleetSpecValueInput({
    required this.fieldKey,
    this.valueString,
    this.valueNumeric,
    this.valueBoolean,
  });

  @override
  List<Object?> get props => [fieldKey, valueString, valueNumeric, valueBoolean];
}
