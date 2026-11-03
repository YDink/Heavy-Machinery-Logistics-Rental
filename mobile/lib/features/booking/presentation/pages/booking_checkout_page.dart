import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import '../bloc/booking_bloc.dart';
import '../bloc/booking_event.dart';
import '../bloc/booking_state.dart';
import '../widgets/booking_price_breakdown_card.dart';

class BookingCheckoutPage extends StatefulWidget {
  final String fleetId;
  final String machineTitle;
  final int dailyRateMinor;
  final int minRentalDays;
  final List<String> bundledAttachmentIds;
  final int bundledAttachmentsDailyRateMinor;

  const BookingCheckoutPage({
    super.key,
    required this.fleetId,
    required this.machineTitle,
    required this.dailyRateMinor,
    this.minRentalDays = 1,
    this.bundledAttachmentIds = const [],
    this.bundledAttachmentsDailyRateMinor = 0,
  });

  @override
  State<BookingCheckoutPage> createState() => _BookingCheckoutPageState();
}

class _BookingCheckoutPageState extends State<BookingCheckoutPage> {
  final _formKey = GlobalKey<FormState>();
  final _addressController = TextEditingController();

  late DateTime _startDate;
  late DateTime _endDate;
  int _durationDays = 1;

  @override
  void initState() {
    super.initState();
    _startDate = DateTime.now().add(const Duration(days: 1));
    _endDate = _startDate.add(Duration(days: widget.minRentalDays > 0 ? widget.minRentalDays : 1));
    _calculateDuration();
  }

  void _calculateDuration() {
    final diff = _endDate.difference(_startDate).inDays;
    _durationDays = diff > 0 ? diff : 1;
  }

  @override
  void dispose() {
    _addressController.dispose();
    super.dispose();
  }

  Future<void> _selectDateRange() async {
    final picked = await showDateRangePicker(
      context: context,
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 180)),
      initialDateRange: DateTimeRange(start: _startDate, end: _endDate),
    );

    if (picked != null) {
      setState(() {
        _startDate = picked.start;
        _endDate = picked.end;
        _calculateDuration();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final theme = Theme.of(context);

    final baseFee = _durationDays * widget.dailyRateMinor;
    final attsFee = _durationDays * widget.bundledAttachmentsDailyRateMinor;
    final gross = baseFee + attsFee;
    final prepay = gross ~/ 2;
    final remain = gross - prepay;

    return Scaffold(
      appBar: AppBar(
        title: Text(widget.machineTitle),
      ),
      body: BlocConsumer<BookingBloc, BookingState>(
        listener: (context, state) {
          if (state is BookingError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: theme.colorScheme.error,
              ),
            );
          } else if (state is BookingLoaded) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text('Reservation hold created: ${state.booking.bookingReference}'),
                backgroundColor: Colors.green,
              ),
            );
          }
        },
        builder: (context, state) {
          final isLoading = state is BookingLoading;

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  // Date Selection Card
                  Card(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: ListTile(
                      contentPadding: const EdgeInsets.all(16),
                      leading: const Icon(Icons.date_range, size: 32),
                      title: Text(l10n.rentalDuration),
                      subtitle: Text(
                        '${_startDate.toLocal().toString().split(' ')[0]} → ${_endDate.toLocal().toString().split(' ')[0]} ($_durationDays ${l10n.daysUnit})',
                        style: theme.textTheme.bodyMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      trailing: TextButton(
                        onPressed: isLoading ? null : _selectDateRange,
                        child: Text(l10n.applyFilters),
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Site Location Address Input
                  TextFormField(
                    controller: _addressController,
                    decoration: InputDecoration(
                      labelText: l10n.siteLocation,
                      hintText: 'e.g. Bole Lemi Industrial Zone Phase 2',
                      prefixIcon: const Icon(Icons.location_on_outlined),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    validator: (val) {
                      if (val == null || val.trim().isEmpty) {
                        return 'Site address is required';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 16),

                  // Price Breakdown Card
                  BookingPriceBreakdownCard(
                    durationDays: _durationDays,
                    dailyRateMinor: widget.dailyRateMinor,
                    baseHireFeeMinor: baseFee,
                    attachmentsFeeMinor: attsFee,
                    grossRentalAmountMinor: gross,
                    prepaymentAmountMinor: prepay,
                    remainingBalanceMinor: remain,
                  ),
                  const SizedBox(height: 24),

                  // Request to Book Button
                  SizedBox(
                    height: 56,
                    child: ElevatedButton(
                      onPressed: isLoading
                          ? null
                          : () {
                              if (_formKey.currentState?.validate() ?? false) {
                                context.read<BookingBloc>().add(
                                      CreateBookingRequested(
                                        fleetId: widget.fleetId,
                                        startDate: _startDate,
                                        endDate: _endDate,
                                        siteLocationAddress: _addressController.text.trim(),
                                        attachmentIds: widget.bundledAttachmentIds,
                                      ),
                                    );
                              }
                            },
                      style: ElevatedButton.styleFrom(
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      child: isLoading
                          ? const CircularProgressIndicator.adaptive()
                          : Text(
                              l10n.requestToBook,
                              style: theme.textTheme.titleMedium?.copyWith(
                                color: theme.colorScheme.onPrimary,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}
