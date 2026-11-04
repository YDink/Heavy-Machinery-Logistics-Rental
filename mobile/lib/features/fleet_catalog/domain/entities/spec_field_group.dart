import 'package:equatable/equatable.dart';
import 'spec_field_definition.dart';

export 'spec_field_definition.dart';
export 'spec_field_type.dart';

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
