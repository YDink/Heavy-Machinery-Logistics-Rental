import 'package:flutter/material.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import '../../domain/entities/booking.dart';

class BookingStatusBadge extends StatelessWidget {
  final BookingStatus status;

  const BookingStatusBadge({super.key, required this.status});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final theme = Theme.of(context);

    Color color;
    String text;

    switch (status) {
      case BookingStatus.draft:
      case BookingStatus.pendingPayment:
        color = Colors.orange;
        text = l10n.bookingStatusPendingPayment;
        break;
      case BookingStatus.paymentEvidenceUploaded:
        color = Colors.blue;
        text = l10n.bookingStatusEvidenceUploaded;
        break;
      case BookingStatus.confirmed:
        color = Colors.green;
        text = l10n.bookingStatusConfirmed;
        break;
      case BookingStatus.active:
        color = Colors.teal;
        text = l10n.bookingStatusActive;
        break;
      case BookingStatus.completed:
      case BookingStatus.rated:
        color = Colors.purple;
        text = l10n.bookingStatusCompleted;
        break;
      case BookingStatus.cancelled:
      case BookingStatus.paymentDeclined:
        color = theme.colorScheme.error;
        text = l10n.bookingStatusCancelled;
        break;
      default:
        color = Colors.grey;
        text = status.name;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.15),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: color.withValues(alpha: 0.6)),
      ),
      child: Text(
        text,
        style: theme.textTheme.labelSmall?.copyWith(
          color: color,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}
