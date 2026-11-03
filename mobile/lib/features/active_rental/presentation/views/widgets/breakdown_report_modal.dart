import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:dink_mobile/features/active_rental/presentation/bloc/active_rental_bloc.dart';
import 'package:dink_mobile/features/active_rental/presentation/bloc/active_rental_event.dart';

class BreakdownReportModal extends StatefulWidget {
  final String bookingId;
  final double currentMeterHours;

  const BreakdownReportModal({
    super.key,
    required this.bookingId,
    required this.currentMeterHours,
  });

  @override
  State<BreakdownReportModal> createState() => _BreakdownReportModalState();
}

class _BreakdownReportModalState extends State<BreakdownReportModal> {
  final _formKey = GlobalKey<FormState>();
  String _category = 'HYDRAULIC_CIRCUIT';
  late final TextEditingController _descController;
  late final TextEditingController _meterController;
  late final TextEditingController _dtcController;

  final List<String> _categories = [
    'ENGINE',
    'HYDRAULIC_CIRCUIT',
    'ELECTRICAL',
    'UNDERCARRIAGE',
    'STRUCTURAL',
    'BOOM_BUCKET',
    'OTHER',
  ];

  @override
  void initState() {
    super.initState();
    _descController = TextEditingController();
    _meterController = TextEditingController(text: widget.currentMeterHours.toStringAsFixed(1));
    _dtcController = TextEditingController();
  }

  @override
  void dispose() {
    _descController.dispose();
    _meterController.dispose();
    _dtcController.dispose();
    super.dispose();
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;

    final meter = double.tryParse(_meterController.text.trim()) ?? widget.currentMeterHours;
    final dtcs = _dtcController.text.trim().isNotEmpty
        ? _dtcController.text.trim().split(RegExp(r'[\s,]+'))
        : <String>[];

    context.read<ActiveRentalBloc>().add(
          ReportMechanicalBreakdownEvent(
            bookingId: widget.bookingId,
            faultCategory: _category,
            faultDescription: _descController.text.trim(),
            diagnosticTroubleCodes: dtcs,
            currentHourMeter: meter,
          ),
        );

    Navigator.of(context).pop();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Padding(
      padding: EdgeInsets.only(
        left: 20,
        right: 20,
        top: 24,
        bottom: MediaQuery.of(context).viewInsets.bottom + 24,
      ),
      child: Form(
        key: _formKey,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(
                l10n.breakdownModalTitle,
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
              const SizedBox(height: 16),
              // Category Dropdown
              DropdownButtonFormField<String>(
                initialValue: _category,
                dropdownColor: const Color(0xFF1F2937),
                decoration: InputDecoration(
                  labelText: l10n.faultCategoryLabel,
                  filled: true,
                  fillColor: const Color(0xFF374151),
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                ),
                items: _categories.map((c) {
                  return DropdownMenuItem(
                    value: c,
                    child: Text(c, style: const TextStyle(color: Colors.white)),
                  );
                }).toList(),
                onChanged: (val) {
                  if (val != null) setState(() => _category = val);
                },
              ),
              const SizedBox(height: 16),
              // Hour Meter Input
              TextFormField(
                controller: _meterController,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                decoration: InputDecoration(
                  labelText: l10n.meterReadingLabel,
                  filled: true,
                  fillColor: const Color(0xFF374151),
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                ),
                style: const TextStyle(color: Colors.white),
                validator: (val) {
                  if (val == null || val.trim().isEmpty) return 'Please enter hour meter';
                  if (double.tryParse(val.trim()) == null) return 'Must be a valid number';
                  return null;
                },
              ),
              const SizedBox(height: 16),
              // Fault Description
              TextFormField(
                controller: _descController,
                maxLines: 3,
                decoration: InputDecoration(
                  labelText: l10n.faultDescriptionLabel,
                  filled: true,
                  fillColor: const Color(0xFF374151),
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                ),
                style: const TextStyle(color: Colors.white),
                validator: (val) {
                  if (val == null || val.trim().isEmpty) return 'Description is required';
                  return null;
                },
              ),
              const SizedBox(height: 16),
              // DTC Codes
              TextFormField(
                controller: _dtcController,
                decoration: InputDecoration(
                  labelText: l10n.dtcLabel,
                  hintText: 'e.g. HYD-042, ERR-901',
                  hintStyle: TextStyle(color: Colors.grey.shade400),
                  filled: true,
                  fillColor: const Color(0xFF374151),
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                ),
                style: const TextStyle(color: Colors.white),
              ),
              const SizedBox(height: 24),
              // Submit button (56dp height per accessibility rules)
              SizedBox(
                height: 56,
                child: ElevatedButton.icon(
                  onPressed: _submit,
                  icon: const Icon(Icons.warning_amber_rounded, color: Colors.white),
                  label: Text(
                    l10n.submitBreakdownButton,
                    style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                  ),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.redAccent.shade700,
                    foregroundColor: Colors.white,
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
