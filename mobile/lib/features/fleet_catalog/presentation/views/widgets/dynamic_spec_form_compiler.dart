import 'package:flutter/material.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/spec_field_group.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class DynamicSpecFormCompiler extends StatelessWidget {
  final List<SpecFieldGroup> specGroups;
  final Map<String, FleetSpecValueInput> currentValues;
  final Function(String fieldKey, String? strVal, double? numVal, bool? boolVal)
      onValueChanged;

  const DynamicSpecFormCompiler({
    super.key,
    required this.specGroups,
    required this.currentValues,
    required this.onValueChanged,
  });

  @override
  Widget build(BuildContext context) {
    if (specGroups.isEmpty) {
      return const SizedBox.shrink();
    }

    final locale = Localizations.localeOf(context).languageCode;
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.technicalSpecifications,
          style: theme.textTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        ...specGroups.map((group) {
          final title = group.getLocalizedSectionName(locale);
          return Card(
            margin: const EdgeInsets.only(bottom: 12),
            elevation: 0,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
              side: BorderSide(color: theme.colorScheme.outlineVariant),
            ),
            child: ExpansionTile(
              initiallyExpanded: true,
              title: Text(
                title,
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
              ),
              childrenPadding:
                  const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              children: group.fields.map((field) {
                return _buildFieldWidget(context, field, locale);
              }).toList(),
            ),
          );
        }),
      ],
    );
  }

  Widget _buildFieldWidget(
    BuildContext context,
    SpecFieldDefinition field,
    String locale,
  ) {
    final label = field.getLocalizedLabel(locale);
    final currentVal = currentValues[field.fieldKey];

    switch (field.fieldType) {
      case SpecFieldType.boolean:
        final boolVal = currentVal?.valueBoolean ?? false;
        return SwitchListTile(
          contentPadding: EdgeInsets.zero,
          title: Text(label),
          value: boolVal,
          onChanged: (val) {
            onValueChanged(field.fieldKey, null, null, val);
          },
        );

      case SpecFieldType.enumeration:
        final selectedVal = currentVal?.valueString;
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 8),
          child: DropdownButtonFormField<String>(
            decoration: InputDecoration(
              labelText: field.isRequired ? '$label *' : label,
              border: const OutlineInputBorder(),
            ),
            initialValue:
                field.enumValues.contains(selectedVal) ? selectedVal : null,
            items: field.enumValues.map((val) {
              return DropdownMenuItem(
                value: val,
                child: Text(val),
              );
            }).toList(),
            onChanged: (val) {
              if (val != null) {
                onValueChanged(field.fieldKey, val, null, null);
              }
            },
          ),
        );

      case SpecFieldType.decimal:
      case SpecFieldType.integer:
        final isDecimal = field.fieldType == SpecFieldType.decimal;
        final numVal = currentVal?.valueNumeric;
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 8),
          child: TextFormField(
            initialValue: numVal != null ? numVal.toString() : '',
            keyboardType: TextInputType.numberWithOptions(decimal: isDecimal),
            decoration: InputDecoration(
              labelText: field.isRequired ? '$label *' : label,
              suffixText: field.unitSymbol,
              border: const OutlineInputBorder(),
            ),
            onChanged: (val) {
              final parsed = double.tryParse(val);
              onValueChanged(field.fieldKey, null, parsed, null);
            },
          ),
        );

      case SpecFieldType.string:
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 8),
          child: TextFormField(
            initialValue: currentVal?.valueString ?? '',
            decoration: InputDecoration(
              labelText: field.isRequired ? '$label *' : label,
              border: const OutlineInputBorder(),
            ),
            onChanged: (val) {
              onValueChanged(field.fieldKey, val, null, null);
            },
          ),
        );
    }
  }
}
