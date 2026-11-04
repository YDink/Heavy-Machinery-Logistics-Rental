import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class BookingPriceBreakdownCard extends StatelessWidget {
  final int durationDays;
  final int dailyRateMinor;
  final int baseHireFeeMinor;
  final int attachmentsFeeMinor;
  final int grossRentalAmountMinor;
  final int prepaymentAmountMinor;
  final int remainingBalanceMinor;
  final String currency;

  const BookingPriceBreakdownCard({
    super.key,
    required this.durationDays,
    required this.dailyRateMinor,
    required this.baseHireFeeMinor,
    required this.attachmentsFeeMinor,
    required this.grossRentalAmountMinor,
    required this.prepaymentAmountMinor,
    required this.remainingBalanceMinor,
    this.currency = 'ETB',
  });

  String _formatMoney(int minorUnits) {
    final amount = minorUnits / 100.0;
    final formatter = NumberFormat('#,##0.00', 'en_US');
    return '${formatter.format(amount)} $currency';
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final theme = Theme.of(context);

    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              l10n.bookingDetails,
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            _buildRow(
              context,
              l10n.rentalDuration,
              '$durationDays ${l10n.daysUnit}',
            ),
            const SizedBox(height: 8),
            _buildRow(
              context,
              l10n.baseHireFee,
              _formatMoney(baseHireFeeMinor),
            ),
            if (attachmentsFeeMinor > 0) ...[
              const SizedBox(height: 8),
              _buildRow(
                context,
                l10n.attachmentsFee,
                _formatMoney(attachmentsFeeMinor),
              ),
            ],
            const Divider(height: 24),
            _buildRow(
              context,
              l10n.grossRentalAmount,
              _formatMoney(grossRentalAmountMinor),
              isBold: true,
            ),
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: theme.colorScheme.primaryContainer.withValues(alpha: 0.3),
                borderRadius: BorderRadius.circular(8),
                border: Border.all(
                  color: theme.colorScheme.primary.withValues(alpha: 0.4),
                ),
              ),
              child: Column(
                children: [
                  _buildRow(
                    context,
                    l10n.advancePrepayment,
                    _formatMoney(prepaymentAmountMinor),
                    highlightColor: theme.colorScheme.primary,
                    isBold: true,
                  ),
                  const SizedBox(height: 8),
                  _buildRow(
                    context,
                    l10n.remainingBalance,
                    _formatMoney(remainingBalanceMinor),
                    isSubtle: true,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRow(
    BuildContext context,
    String label,
    String value, {
    bool isBold = false,
    bool isSubtle = false,
    Color? highlightColor,
  }) {
    final theme = Theme.of(context);
    final textStyle = theme.textTheme.bodyMedium?.copyWith(
      fontWeight: isBold ? FontWeight.bold : FontWeight.normal,
      color: highlightColor ?? (isSubtle ? theme.colorScheme.outline : null),
    );

    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label, style: textStyle),
        Text(value, style: textStyle),
      ],
    );
  }
}
