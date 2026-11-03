import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:intl/intl.dart';
import 'package:dink_mobile/features/active_rental/presentation/bloc/active_rental_bloc.dart';
import 'package:dink_mobile/features/active_rental/presentation/bloc/active_rental_event.dart';

class ExtensionRequestSheet extends StatefulWidget {
  final String bookingId;

  const ExtensionRequestSheet({
    super.key,
    required this.bookingId,
  });

  @override
  State<ExtensionRequestSheet> createState() => _ExtensionRequestSheetState();
}

class _ExtensionRequestSheetState extends State<ExtensionRequestSheet> {
  DateTime _selectedDate = DateTime.now().add(const Duration(days: 3));
  final DateFormat _dateFormat = DateFormat('yyyy-MM-dd');

  Future<void> _pickDate() async {
    final now = DateTime.now();
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: now.add(const Duration(days: 1)),
      lastDate: now.add(const Duration(days: 90)),
      builder: (context, child) {
        return Theme(
          data: ThemeData.dark().copyWith(
            colorScheme: const ColorScheme.dark(
              primary: Color(0xFF2563EB),
              onPrimary: Colors.white,
              surface: Color(0xFF1F2937),
              onSurface: Colors.white,
            ),
          ),
          child: child!,
        );
      },
    );

    if (picked != null) {
      setState(() => _selectedDate = picked);
    }
  }

  void _submit() {
    context.read<ActiveRentalBloc>().add(
          RequestRentalExtensionEvent(
            bookingId: widget.bookingId,
            newEndDate: _selectedDate,
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
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(
            l10n.extensionSheetTitle,
            style: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            l10n.extensionQuoteLabel,
            style: TextStyle(fontSize: 14, color: Colors.grey.shade400),
          ),
          const SizedBox(height: 20),
          InkWell(
            onTap: _pickDate,
            borderRadius: BorderRadius.circular(12),
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 18),
              decoration: BoxDecoration(
                color: const Color(0xFF374151),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: Colors.grey.shade600),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '${l10n.newEndDateLabel}: ${_dateFormat.format(_selectedDate)}',
                    style: const TextStyle(fontSize: 16, color: Colors.white),
                  ),
                  const Icon(Icons.calendar_month, color: Colors.white70),
                ],
              ),
            ),
          ),
          const SizedBox(height: 24),
          // Submit button (56dp height)
          SizedBox(
            height: 56,
            child: ElevatedButton.icon(
              onPressed: _submit,
              icon: const Icon(Icons.schedule_send, color: Colors.white),
              label: Text(
                l10n.submitExtensionButton,
                style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
              ),
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF2563EB),
                foregroundColor: Colors.white,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
