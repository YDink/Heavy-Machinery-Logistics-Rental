import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_bloc.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_event.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_state.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/widgets/dynamic_spec_form_compiler.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/widgets/walkaround_camera_reticle.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class FleetRegistrationPage extends StatefulWidget {
  const FleetRegistrationPage({super.key});

  @override
  State<FleetRegistrationPage> createState() => _FleetRegistrationPageState();
}

class _FleetRegistrationPageState extends State<FleetRegistrationPage> {
  final _formKey = GlobalKey<FormState>();

  final _makeController = TextEditingController();
  final _modelController = TextEditingController();
  final _yearController = TextEditingController();
  final _plateController = TextEditingController();
  final _chassisController = TextEditingController();
  final _meterController = TextEditingController();
  final _dailyRateController = TextEditingController();

  String _fuelType = 'DIESEL';
  bool _selfMoveCapable = false;
  bool _operatorIncluded = false;

  @override
  void initState() {
    super.initState();
    context.read<FleetRegistrationBloc>().add(const LoadCategoriesEvent());
  }

  @override
  void dispose() {
    _makeController.dispose();
    _modelController.dispose();
    _yearController.dispose();
    _plateController.dispose();
    _chassisController.dispose();
    _meterController.dispose();
    _dailyRateController.dispose();
    super.dispose();
  }

  FleetRegistrationRequest _buildRequest(CategoriesLoaded state) {
    return FleetRegistrationRequest(
      categoryId: state.selectedCategoryId ?? '',
      make: _makeController.text.trim(),
      model: _modelController.text.trim(),
      manufacturingYear: int.tryParse(_yearController.text) ?? 2022,
      registrationNumber: _plateController.text.trim(),
      chassisVin: _chassisController.text.trim(),
      baselineMeterHours: double.tryParse(_meterController.text) ?? 0.0,
      fuelType: _fuelType,
      selfMoveCapable: _selfMoveCapable,
      dailyRateMinor:
          ((double.tryParse(_dailyRateController.text) ?? 0.0) * 100).toInt(),
      operatorIncluded: _operatorIncluded,
      latitude: 9.0300,
      longitude: 38.7400,
      addressText: 'Addis Ababa',
      photos: state.photos.values.toList(),
      specValues: state.specValues.values.toList(),
    );
  }

  void _submit(CategoriesLoaded state) {
    final l10n = AppLocalizations.of(context)!;
    if (!_formKey.currentState!.validate()) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(l10n.errorMissingFields)),
      );
      return;
    }

    if (!state.hasCompleteWalkaroundPhotos) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(l10n.errorIncompletePhotos),
          backgroundColor: Colors.red.shade700,
        ),
      );
      return;
    }

    final request = _buildRequest(state);
    context
        .read<FleetRegistrationBloc>()
        .add(SubmitFleetRegistrationEvent(request));
  }

  void _saveDraft(CategoriesLoaded state) {
    final request = _buildRequest(state);
    context.read<FleetRegistrationBloc>().add(SaveDraftEvent(request));
  }

  void _simulatePhotoCapture(PhotoAngle angle) {
    final photo = FleetPhoto(
      angle: angle,
      storageUri: 's3://fleet/${angle.code.toLowerCase()}.jpg',
      fileHash: 'mock_hash_${angle.code}',
    );
    context.read<FleetRegistrationBloc>().add(UpdatePhotoEvent(photo));
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.fleetRegistrationTitle),
      ),
      body: BlocConsumer<FleetRegistrationBloc, FleetRegistrationState>(
        listener: (context, state) {
          if (state is FleetSubmitSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(l10n.registrationSuccessMessage),
                backgroundColor: Colors.green,
              ),
            );
          } else if (state is FleetSubmitFailure) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: Colors.red,
              ),
            );
          } else if (state is CategoriesLoaded && state.isDraftSaved) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(l10n.draftSavedMessage),
                duration: const Duration(seconds: 2),
              ),
            );
          }
        },
        builder: (context, state) {
          if (state is CategoriesLoading || state is FleetRegistrationInitial) {
            return const Center(child: CircularProgressIndicator());
          }

          if (state is FleetSubmitting) {
            return const Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  CircularProgressIndicator(),
                  SizedBox(height: 16),
                ],
              ),
            );
          }

          if (state is CategoriesLoaded) {
            return SingleChildScrollView(
              padding: const EdgeInsets.all(16.0),
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    // Category Selector
                    DropdownButtonFormField<String>(
                      decoration: InputDecoration(
                        labelText: l10n.selectCategory,
                        border: const OutlineInputBorder(),
                      ),
                      initialValue: state.selectedCategoryId,
                      items: state.categories.map((c) {
                        return DropdownMenuItem(
                          value: c.id,
                          child: Text(c.name),
                        );
                      }).toList(),
                      onChanged: (id) {
                        if (id != null) {
                          context
                              .read<FleetRegistrationBloc>()
                              .add(SelectCategoryEvent(id));
                        }
                      },
                      validator: (val) =>
                          val == null ? l10n.errorMissingFields : null,
                    ),
                    const SizedBox(height: 20),

                    // Core Machinery Details Section
                    Text(
                      l10n.machineryDetails,
                      style: theme.textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 12),
                    TextFormField(
                      controller: _makeController,
                      decoration: InputDecoration(
                        labelText: '${l10n.make} *',
                        border: const OutlineInputBorder(),
                      ),
                      validator: (v) =>
                          v == null || v.isEmpty ? l10n.errorMissingFields : null,
                    ),
                    const SizedBox(height: 12),
                    TextFormField(
                      controller: _modelController,
                      decoration: InputDecoration(
                        labelText: '${l10n.model} *',
                        border: const OutlineInputBorder(),
                      ),
                      validator: (v) =>
                          v == null || v.isEmpty ? l10n.errorMissingFields : null,
                    ),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        Expanded(
                          child: TextFormField(
                            controller: _yearController,
                            keyboardType: TextInputType.number,
                            decoration: InputDecoration(
                              labelText: '${l10n.manufacturingYear} *',
                              border: const OutlineInputBorder(),
                            ),
                            validator: (v) => v == null || v.isEmpty
                                ? l10n.errorMissingFields
                                : null,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: TextFormField(
                            controller: _plateController,
                            decoration: InputDecoration(
                              labelText: '${l10n.registrationNumber} *',
                              border: const OutlineInputBorder(),
                            ),
                            validator: (v) => v == null || v.isEmpty
                                ? l10n.errorMissingFields
                                : null,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    TextFormField(
                      controller: _chassisController,
                      decoration: InputDecoration(
                        labelText: '${l10n.chassisVin} *',
                        border: const OutlineInputBorder(),
                      ),
                      validator: (v) =>
                          v == null || v.isEmpty ? l10n.errorMissingFields : null,
                    ),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        Expanded(
                          child: TextFormField(
                            controller: _meterController,
                            keyboardType: const TextInputType.numberWithOptions(
                                decimal: true),
                            decoration: InputDecoration(
                              labelText: '${l10n.baselineMeterHours} *',
                              suffixText: 'hrs',
                              border: const OutlineInputBorder(),
                            ),
                            validator: (v) {
                              if (v == null || v.isEmpty) {
                                return l10n.errorMissingFields;
                              }
                              final numVal = double.tryParse(v);
                              if (numVal == null || numVal < 0) {
                                return l10n.errorGeneric;
                              }
                              return null;
                            },
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: TextFormField(
                            controller: _dailyRateController,
                            keyboardType: const TextInputType.numberWithOptions(
                                decimal: true),
                            decoration: InputDecoration(
                              labelText: '${l10n.dailyRate} *',
                              suffixText: 'ETB',
                              border: const OutlineInputBorder(),
                            ),
                            validator: (v) => v == null || v.isEmpty
                                ? l10n.errorMissingFields
                                : null,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    DropdownButtonFormField<String>(
                      decoration: InputDecoration(
                        labelText: l10n.fuelType,
                        border: const OutlineInputBorder(),
                      ),
                      initialValue: _fuelType,
                      items: const [
                        DropdownMenuItem(
                            value: 'DIESEL', child: Text('Diesel')),
                        DropdownMenuItem(
                            value: 'ELECTRIC', child: Text('Electric')),
                        DropdownMenuItem(
                            value: 'PETROL', child: Text('Petrol')),
                        DropdownMenuItem(
                            value: 'HYBRID', child: Text('Hybrid')),
                      ],
                      onChanged: (val) {
                        if (val != null) setState(() => _fuelType = val);
                      },
                    ),
                    const SizedBox(height: 8),
                    SwitchListTile(
                      contentPadding: EdgeInsets.zero,
                      title: Text(l10n.selfMoveCapable),
                      value: _selfMoveCapable,
                      onChanged: (v) => setState(() => _selfMoveCapable = v),
                    ),
                    SwitchListTile(
                      contentPadding: EdgeInsets.zero,
                      title: Text(l10n.operatorIncluded),
                      value: _operatorIncluded,
                      onChanged: (v) => setState(() => _operatorIncluded = v),
                    ),
                    const SizedBox(height: 20),

                    // Walkaround 6-Angle Photos Reticle
                    WalkaroundCameraReticle(
                      photos: state.photos,
                      onAngleSelected: _simulatePhotoCapture,
                    ),
                    const SizedBox(height: 20),

                    // Dynamic Technical Specifications
                    DynamicSpecFormCompiler(
                      specGroups: state.specGroups,
                      currentValues: state.specValues,
                      onValueChanged: (key, strVal, numVal, boolVal) {
                        context.read<FleetRegistrationBloc>().add(
                              UpdateSpecValueEvent(
                                FleetSpecValueInput(
                                  fieldKey: key,
                                  valueString: strVal,
                                  valueNumeric: numVal,
                                  valueBoolean: boolVal,
                                ),
                              ),
                            );
                      },
                    ),
                    const SizedBox(height: 24),

                    // Action Buttons
                    ElevatedButton(
                      key: const Key('submitFleetButton'),
                      onPressed: () => _submit(state),
                      style: ElevatedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      child: Text(
                        l10n.submitRegistration,
                        style: const TextStyle(
                            fontSize: 16, fontWeight: FontWeight.bold),
                      ),
                    ),
                    const SizedBox(height: 12),
                    OutlinedButton(
                      key: const Key('saveDraftButton'),
                      onPressed: () => _saveDraft(state),
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      child: Text(l10n.saveDraft),
                    ),
                    const SizedBox(height: 32),
                  ],
                ),
              ),
            );
          }

          return const SizedBox.shrink();
        },
      ),
    );
  }
}
