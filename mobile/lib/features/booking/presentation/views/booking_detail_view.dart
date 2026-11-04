import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:dink_mobile/features/booking/domain/entities/booking.dart';
import 'package:dink_mobile/features/booking/presentation/bloc/booking_bloc.dart';
import 'package:dink_mobile/features/booking/presentation/bloc/booking_event.dart';
import 'package:dink_mobile/features/booking/presentation/bloc/booking_state.dart';
import 'package:dink_mobile/features/booking/presentation/views/widgets/booking_price_breakdown_card.dart';
import 'package:dink_mobile/features/booking/presentation/views/widgets/booking_status_badge.dart';
import 'package:dink_mobile/features/booking/presentation/views/widgets/hold_countdown_timer.dart';

class BookingDetailView extends StatefulWidget {
  final String bookingId;
  final bool isOwner;

  const BookingDetailView({
    super.key,
    required this.bookingId,
    this.isOwner = false,
  });

  @override
  State<BookingDetailView> createState() => _BookingDetailViewState();
}

class _BookingDetailViewState extends State<BookingDetailView> {
  @override
  void initState() {
    super.initState();
    context.read<BookingBloc>().add(LoadBookingDetails(widget.bookingId));
  }

  void _showCancelDialog(BuildContext context, BookingEntity booking) {
    final l10n = AppLocalizations.of(context)!;
    final reasonController = TextEditingController();

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(l10n.cancelBooking),
        content: TextField(
          controller: reasonController,
          decoration: const InputDecoration(
            labelText: 'Reason for cancellation',
            hintText: 'e.g. Site delay or geotechnical impediment',
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Dismiss'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(ctx);
              context.read<BookingBloc>().add(
                    CancelBookingRequested(
                      bookingId: booking.id,
                      expectedVersion: booking.version,
                      reason: reasonController.text.trim().isEmpty
                          ? 'Contractor requested cancellation'
                          : reasonController.text.trim(),
                    ),
                  );
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: Theme.of(context).colorScheme.error,
            ),
            child: Text(l10n.cancelBooking),
          ),
        ],
      ),
    );
  }

  void _showHandoverPinDialog(BuildContext context, BookingEntity booking) {
    final l10n = AppLocalizations.of(context)!;
    final pinController = TextEditingController();
    final meterController = TextEditingController();

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(l10n.handoverPIN),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: pinController,
              keyboardType: TextInputType.number,
              maxLength: 6,
              decoration: InputDecoration(
                labelText: l10n.handoverPIN,
                hintText: '6-digit PIN',
              ),
            ),
            TextField(
              controller: meterController,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: InputDecoration(
                labelText: l10n.checkInMeter,
                hintText: 'e.g. 1205.5',
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              final pin = pinController.text.trim();
              final meter = double.tryParse(meterController.text.trim()) ?? 0.0;
              Navigator.pop(ctx);
              context.read<BookingBloc>().add(
                    HandoverCheckInRequested(
                      bookingId: booking.id,
                      expectedVersion: booking.version,
                      handoverPin: pin,
                      operatingHoursMeter: meter,
                    ),
                  );
            },
            child: const Text('Verify & Start Rental'),
          ),
        ],
      ),
    );
  }

  void _showReturnDialog(BuildContext context, BookingEntity booking) {
    final l10n = AppLocalizations.of(context)!;
    final meterController = TextEditingController();
    final damageController = TextEditingController();

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Return Check-Out'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: meterController,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: InputDecoration(
                labelText: l10n.checkOutMeter,
                hintText: 'e.g. 1230.0',
              ),
            ),
            TextField(
              controller: damageController,
              decoration: const InputDecoration(
                labelText: 'Damage or Discrepancy Notes (Optional)',
                hintText: 'Leave empty if machine is clean and intact',
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              final meter = double.tryParse(meterController.text.trim()) ?? 0.0;
              Navigator.pop(ctx);
              context.read<BookingBloc>().add(
                    ReturnCheckOutRequested(
                      bookingId: booking.id,
                      expectedVersion: booking.version,
                      operatingHoursMeter: meter,
                      damageNotes: damageController.text.trim().isEmpty
                          ? null
                          : damageController.text.trim(),
                    ),
                  );
            },
            child: const Text('Confirm Return'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.bookingDetails),
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
          } else if (state is BookingActionSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: Colors.green,
              ),
            );
          }
        },
        builder: (context, state) {
          if (state is BookingLoading) {
            return const Center(child: CircularProgressIndicator.adaptive());
          }

          if (state is BookingLoaded) {
            final b = state.booking;

            return SingleChildScrollView(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Card(
                    elevation: 1,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                b.bookingReference,
                                style: theme.textTheme.headlineSmall?.copyWith(
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                              BookingStatusBadge(status: b.status),
                            ],
                          ),
                          const SizedBox(height: 8),
                          Text(
                            '${b.startDate.toLocal().toString().split(' ')[0]} → ${b.endDate.toLocal().toString().split(' ')[0]} (${b.durationDays} ${l10n.daysUnit})',
                            style: theme.textTheme.bodyMedium?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
                            ),
                          ),
                          const SizedBox(height: 4),
                          Text(
                            '${l10n.siteLocation}: ${b.siteLocationAddress}',
                            style: theme.textTheme.bodyMedium,
                          ),
                          if (b.status == BookingStatus.pendingPayment &&
                              b.expiresAt != null) ...[
                            const SizedBox(height: 12),
                            HoldCountdownTimer(expiresAt: b.expiresAt!),
                          ],
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),
                  BookingPriceBreakdownCard(
                    durationDays: b.durationDays,
                    dailyRateMinor: b.dailyRateMinor,
                    baseHireFeeMinor: b.baseHireFeeMinor,
                    attachmentsFeeMinor: b.attachmentsFeeMinor,
                    grossRentalAmountMinor: b.grossRentalAmountMinor,
                    prepaymentAmountMinor: b.prepaymentAmountMinor,
                    remainingBalanceMinor: b.grossRentalAmountMinor - b.prepaymentAmountMinor,
                    currency: b.currency,
                  ),
                  const SizedBox(height: 24),
                  if (b.status == BookingStatus.pendingPayment ||
                      b.status == BookingStatus.paymentDeclined) ...[
                    SizedBox(
                      height: 56,
                      child: ElevatedButton.icon(
                        onPressed: () {
                          context.read<BookingBloc>().add(
                                SubmitPaymentSlipRequested(
                                  bookingId: b.id,
                                  storageUri: 's3://payment-vault/slip_${b.id}.webp',
                                  fileHash: 'sha256_mock_slip_file_hash_123',
                                  amountMinor: b.prepaymentAmountMinor,
                                  transactionReference: 'WIRE-AUTOGEN-${b.bookingReference}',
                                ),
                              );
                        },
                        icon: const Icon(Icons.receipt_long),
                        label: Text(l10n.uploadBankSlip),
                      ),
                    ),
                    const SizedBox(height: 12),
                    OutlinedButton(
                      onPressed: () => _showCancelDialog(context, b),
                      child: Text(l10n.cancelBooking),
                    ),
                  ],

                  if (b.status == BookingStatus.paymentEvidenceUploaded && widget.isOwner) ...[
                    SizedBox(
                      height: 56,
                      child: ElevatedButton.icon(
                        onPressed: () {
                          context.read<BookingBloc>().add(
                                ConfirmPaymentRequested(
                                  bookingId: b.id,
                                  expectedVersion: b.version,
                                  paymentReference: 'VERIFIED-OWNER-${b.bookingReference}',
                                ),
                              );
                        },
                        icon: const Icon(Icons.check_circle_outline),
                        label: Text(l10n.confirmPaymentReceived),
                        style: ElevatedButton.styleFrom(backgroundColor: Colors.green),
                      ),
                    ),
                    const SizedBox(height: 12),
                    OutlinedButton(
                      onPressed: () {
                        context.read<BookingBloc>().add(
                              DeclinePaymentRequested(
                                bookingId: b.id,
                                expectedVersion: b.version,
                                reasonCode: 'FUNDS_NOT_RECEIVED_IN_ACCOUNT',
                                rejectionNotes: 'Bank wire statement unverified',
                              ),
                            );
                      },
                      child: Text(l10n.declinePayment),
                    ),
                  ],

                  if (b.status == BookingStatus.paymentDeclined && !widget.isOwner) ...[
                    ElevatedButton(
                      onPressed: () {
                        context.read<BookingBloc>().add(
                              DisputeBookingRequested(
                                bookingId: b.id,
                                expectedVersion: b.version,
                                disputeReason: 'Funds transferred via bank app reference',
                              ),
                            );
                      },
                      child: Text(l10n.disputeDecline),
                    ),
                  ],

                  if (b.status == BookingStatus.confirmed) ...[
                    Card(
                      color: Colors.amber.shade50,
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Column(
                          children: [
                            Text(
                              widget.isOwner
                                  ? 'Share this PIN with the contractor upon machine delivery:'
                                  : 'Provide this challenge PIN to the machine operator upon site inspection:',
                              style: theme.textTheme.bodyMedium,
                            ),
                            const SizedBox(height: 8),
                            Text(
                              b.handoverPin ?? '------',
                              style: theme.textTheme.headlineMedium?.copyWith(
                                fontWeight: FontWeight.bold,
                                letterSpacing: 4,
                                color: Colors.blueGrey.shade900,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),
                    if (!widget.isOwner)
                      SizedBox(
                        height: 56,
                        child: ElevatedButton(
                          onPressed: () => _showHandoverPinDialog(context, b),
                          child: const Text('Enter PIN & Complete Handover'),
                        ),
                      ),
                  ],

                  if (b.status == BookingStatus.active && widget.isOwner) ...[
                    SizedBox(
                      height: 56,
                      child: ElevatedButton.icon(
                        onPressed: () => _showReturnDialog(context, b),
                        icon: const Icon(Icons.assignment_turned_in),
                        label: const Text('Record Machine Return Check-Out'),
                      ),
                    ),
                  ],
                ],
              ),
            );
          }

          return const SizedBox.shrink();
        },
      ),
    );
  }
}

typedef BookingDetailPage = BookingDetailView;
