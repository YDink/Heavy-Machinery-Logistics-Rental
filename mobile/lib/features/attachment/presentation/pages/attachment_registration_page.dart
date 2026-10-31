import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_bloc.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_event.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_state.dart';
import 'package:dink_mobile/features/attachment/presentation/widgets/attachment_photo_reticle.dart';

class AttachmentRegistrationPage extends StatefulWidget {
  const AttachmentRegistrationPage({super.key});

  @override
  State<AttachmentRegistrationPage> createState() => _AttachmentRegistrationPageState();
}

class _AttachmentRegistrationPageState extends State<AttachmentRegistrationPage> {
  final _formKey = GlobalKey<FormState>();

  AttachmentType _selectedType = AttachmentType.hydraulicBreaker;
  final _makeController = TextEditingController();
  final _modelController = TextEditingController();
  final _subTypeController = TextEditingController();
  final _serialController = TextEditingController();
  final _weightController = TextEditingController();
  final _minCarrierController = TextEditingController();
  final _maxCarrierController = TextEditingController();
  final _pressureController = TextEditingController();
  final _minFlowController = TextEditingController();
  final _maxFlowController = TextEditingController();
  final _pinDiameterController = TextEditingController();
  final _pinSpreadController = TextEditingController();
  final _dailyRateController = TextEditingController();
  final _securityDepositController = TextEditingController();

  @override
  void initState() {
    super.initState();
    context.read<AttachmentBloc>().add(const LoadAttachmentDraftEvent());
  }

  @override
  void dispose() {
    _makeController.dispose();
    _modelController.dispose();
    _subTypeController.dispose();
    _serialController.dispose();
    _weightController.dispose();
    _minCarrierController.dispose();
    _maxCarrierController.dispose();
    _pressureController.dispose();
    _minFlowController.dispose();
    _maxFlowController.dispose();
    _pinDiameterController.dispose();
    _pinSpreadController.dispose();
    _dailyRateController.dispose();
    _securityDepositController.dispose();
    super.dispose();
  }

  void _populateFromDraft(AttachmentRegistrationRequest draft) {
    setState(() {
      _selectedType = draft.attachmentType;
      _makeController.text = draft.make;
      _modelController.text = draft.model;
      _subTypeController.text = draft.subTypeSize;
      _serialController.text = draft.serialNumber;
      _weightController.text = draft.operatingWeightKg > 0 ? draft.operatingWeightKg.toString() : '';
      _minCarrierController.text = draft.minCarrierWeightTonnes > 0 ? draft.minCarrierWeightTonnes.toString() : '';
      _maxCarrierController.text = draft.maxCarrierWeightTonnes > 0 ? draft.maxCarrierWeightTonnes.toString() : '';
      _pressureController.text = draft.requiredPressureBar > 0 ? draft.requiredPressureBar.toString() : '';
      _minFlowController.text = draft.requiredFlowMinLmin > 0 ? draft.requiredFlowMinLmin.toString() : '';
      _maxFlowController.text = draft.requiredFlowMaxLmin > 0 ? draft.requiredFlowMaxLmin.toString() : '';
      _pinDiameterController.text = draft.pinDiameterMm > 0 ? draft.pinDiameterMm.toString() : '';
      _pinSpreadController.text = draft.pinSpreadMm > 0 ? draft.pinSpreadMm.toString() : '';
      _dailyRateController.text = draft.standaloneDailyRateMinor > 0 ? (draft.standaloneDailyRateMinor ~/ 100).toString() : '';
      _securityDepositController.text = draft.securityDepositMinor > 0 ? (draft.securityDepositMinor ~/ 100).toString() : '';
    });
  }

  AttachmentRegistrationRequest _buildRequest(Map<AttachmentPhotoAngle, AttachmentPhoto> photos) {
    return AttachmentRegistrationRequest(
      attachmentType: _selectedType,
      make: _makeController.text.trim(),
      model: _modelController.text.trim(),
      subTypeSize: _subTypeController.text.trim(),
      serialNumber: _serialController.text.trim(),
      operatingWeightKg: double.tryParse(_weightController.text) ?? 0.0,
      minCarrierWeightTonnes: double.tryParse(_minCarrierController.text) ?? 0.0,
      maxCarrierWeightTonnes: double.tryParse(_maxCarrierController.text) ?? 0.0,
      requiredPressureBar: double.tryParse(_pressureController.text) ?? 0.0,
      requiredFlowMinLmin: double.tryParse(_minFlowController.text) ?? 0.0,
      requiredFlowMaxLmin: double.tryParse(_maxFlowController.text) ?? 0.0,
      pinDiameterMm: double.tryParse(_pinDiameterController.text) ?? 0.0,
      pinSpreadMm: double.tryParse(_pinSpreadController.text) ?? 0.0,
      compatibleSubcategories: const ['cat_sub_excavator_standard'],
      standaloneDailyRateMinor: ((double.tryParse(_dailyRateController.text) ?? 0.0) * 100).round(),
      securityDepositMinor: ((double.tryParse(_securityDepositController.text) ?? 0.0) * 100).round(),
      minRentalDays: 1,
      photos: photos.values.toList(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return BlocConsumer<AttachmentBloc, AttachmentState>(
      listener: (context, state) {
        if (state.status == AttachmentStatusState.draftLoaded && state.draft != null) {
          _populateFromDraft(state.draft!);
        } else if (state.status == AttachmentStatusState.draftSaved) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text(l10n.draftSavedMessage)),
          );
        } else if (state.status == AttachmentStatusState.registered) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(l10n.attachmentRegisteredSuccess),
              backgroundColor: Colors.green,
            ),
          );
          Navigator.of(context).pop();
        } else if (state.status == AttachmentStatusState.error) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(state.errorMessage ?? l10n.errorGeneric),
              backgroundColor: Colors.red,
            ),
          );
        }
      },
      builder: (context, state) {
        return Scaffold(
          appBar: AppBar(
            title: Text(l10n.registerAttachmentTitle),
            actions: [
              IconButton(
                icon: const Icon(Icons.save_outlined),
                tooltip: l10n.saveDraft,
                onPressed: () {
                  final draft = _buildRequest(state.capturedPhotos);
                  context.read<AttachmentBloc>().add(SaveAttachmentDraftEvent(draft));
                },
              ),
            ],
          ),
          body: Form(
            key: _formKey,
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: [
                DropdownButtonFormField<AttachmentType>(
                  initialValue: _selectedType,
                  decoration: InputDecoration(
                    labelText: l10n.attachmentTypeLabel,
                    border: const OutlineInputBorder(),
                  ),
                  items: AttachmentType.values.map((type) {
                    return DropdownMenuItem(
                      value: type,
                      child: Text(type.displayName),
                    );
                  }).toList(),
                  onChanged: (val) {
                    if (val != null) setState(() => _selectedType = val);
                  },
                ),
                const SizedBox(height: 14),
                TextFormField(
                  controller: _makeController,
                  decoration: InputDecoration(
                    labelText: l10n.attachmentMakeLabel,
                    border: const OutlineInputBorder(),
                  ),
                  validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                ),
                const SizedBox(height: 14),
                TextFormField(
                  controller: _modelController,
                  decoration: InputDecoration(
                    labelText: l10n.attachmentModelLabel,
                    border: const OutlineInputBorder(),
                  ),
                  validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                ),
                const SizedBox(height: 14),
                TextFormField(
                  controller: _subTypeController,
                  decoration: InputDecoration(
                    labelText: l10n.attachmentSubTypeLabel,
                    border: const OutlineInputBorder(),
                  ),
                  validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                ),
                const SizedBox(height: 14),
                TextFormField(
                  controller: _serialController,
                  decoration: InputDecoration(
                    labelText: l10n.attachmentSerialLabel,
                    border: const OutlineInputBorder(),
                  ),
                  validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                ),
                const SizedBox(height: 20),
                const Text(
                  'Mechanical & Carrier Specifications',
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: _weightController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.operatingWeightKgLabel,
                          border: const OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                      ),
                    ),
                    const SizedBox(width: 10),
                    Expanded(
                      child: TextFormField(
                        controller: _pinDiameterController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.pinDiameterLabel,
                          border: const OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 14),
                Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: _minCarrierController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.minCarrierWeightLabel,
                          border: const OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                      ),
                    ),
                    const SizedBox(width: 10),
                    Expanded(
                      child: TextFormField(
                        controller: _maxCarrierController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.maxCarrierWeightLabel,
                          border: const OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 14),
                Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: _minFlowController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.requiredFlowMinLabel,
                          border: const OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                      ),
                    ),
                    const SizedBox(width: 10),
                    Expanded(
                      child: TextFormField(
                        controller: _maxFlowController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.requiredFlowMaxLabel,
                          border: const OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 14),
                Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: _pressureController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.requiredPressureLabel,
                          border: const OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                      ),
                    ),
                    const SizedBox(width: 10),
                    Expanded(
                      child: TextFormField(
                        controller: _pinSpreadController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.pinSpreadLabel,
                          border: const OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 20),
                Text(
                  l10n.tariffAndDepositsSection,
                  style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: _dailyRateController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.standaloneDailyRateLabel,
                          border: const OutlineInputBorder(),
                        ),
                        validator: (v) => (v == null || v.trim().isEmpty) ? l10n.errorMissingFields : null,
                      ),
                    ),
                    const SizedBox(width: 10),
                    Expanded(
                      child: TextFormField(
                        controller: _securityDepositController,
                        keyboardType: TextInputType.number,
                        decoration: InputDecoration(
                          labelText: l10n.securityDepositLabel,
                          border: const OutlineInputBorder(),
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 24),
                AttachmentPhotoReticle(
                  capturedPhotos: state.capturedPhotos,
                  onCapturePhoto: (angle) {
                    context.read<AttachmentBloc>().add(
                          PhotoCapturedEvent(
                            angle: angle,
                            localFilePath: '/mock/path/${angle.code}.jpg',
                            storageUri: 's3://dink-fleet-media/attachments/${angle.code}.webp',
                          ),
                        );
                  },
                ),
                const SizedBox(height: 24),
                SizedBox(
                  height: 48,
                  child: ElevatedButton(
                    onPressed: state.status == AttachmentStatusState.loading
                        ? null
                        : () {
                            if (!_formKey.currentState!.validate()) return;
                            if (!state.areAllPhotosCaptured) {
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(
                                  content: Text(l10n.errorIncompleteAttachmentPhotos),
                                  backgroundColor: Colors.red,
                                ),
                              );
                              return;
                            }
                            final req = _buildRequest(state.capturedPhotos);
                            context.read<AttachmentBloc>().add(SubmitAttachmentRegistrationEvent(req));
                          },
                    child: state.status == AttachmentStatusState.loading
                        ? const CircularProgressIndicator.adaptive()
                        : Text(l10n.submitRegistration),
                  ),
                ),
                const SizedBox(height: 20),
              ],
            ),
          ),
        );
      },
    );
  }
}
