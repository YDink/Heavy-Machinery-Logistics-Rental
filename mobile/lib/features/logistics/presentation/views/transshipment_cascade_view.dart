import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/logistics_bloc.dart';
import '../bloc/logistics_event.dart';
import '../bloc/logistics_state.dart';

class TransshipmentCascadeView extends StatelessWidget {
  final String parentBookingId;
  final String subBookingId;
  final String replacementFleetId;
  final int urgentBountyMinor;
  final int renterRefundMinor;
  final int commissionRestitutionMinor;
  final String currency;

  const TransshipmentCascadeView({
    super.key,
    required this.parentBookingId,
    required this.subBookingId,
    required this.replacementFleetId,
    this.urgentBountyMinor = 720000,
    this.renterRefundMinor = 8750000,
    this.commissionRestitutionMinor = 1750000,
    this.currency = 'ETB',
  });

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    const bgDark = Color(0xFF111827);
    const cardDark = Color(0xFF1C1C1C);
    const borderGray = Color(0xFF3E3E3E);
    const safetyOrange = Color(0xFFFF5722);
    const primaryYellow = Color(0xFFFFB300);

    return Scaffold(
      backgroundColor: bgDark,
      appBar: AppBar(
        backgroundColor: bgDark,
        elevation: 0,
        title: Text(
          l10n.transshipmentSlaTitle.toUpperCase(),
          style: const TextStyle(
            color: Color(0xFFF9FAFB),
            fontWeight: FontWeight.w700,
            fontSize: 15,
          ),
        ),
      ),
      body: BlocBuilder<LogisticsBloc, LogisticsState>(
        builder: (context, state) {
          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Safety Orange 60m SLA Banner
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: const Color(0xFF431407),
                    border: Border.all(color: safetyOrange, width: 2),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        l10n.transshipmentSlaTitle.toUpperCase(),
                        style: const TextStyle(
                          color: safetyOrange,
                          fontWeight: FontWeight.w800,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        '+20% URGENT BOUNTY (${urgentBountyMinor ~/ 100} $currency)',
                        style: const TextStyle(
                          color: Color(0xFFF9FAFB),
                          fontWeight: FontWeight.w700,
                          fontSize: 18,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 12),

                // Parent Booking Frozen Notice
                Container(
                  padding: const EdgeInsets.all(14),
                  decoration: BoxDecoration(
                    color: cardDark,
                    border: Border.all(color: borderGray, width: 2),
                  ),
                  child: Text(
                    l10n.parentBookingFrozenNotice.toUpperCase(),
                    style: const TextStyle(
                      color: primaryYellow,
                      fontWeight: FontWeight.w700,
                      fontSize: 13,
                    ),
                  ),
                ),
                const SizedBox(height: 16),

                // Replacement Carrier Hot-Swap Card
                Container(
                  padding: const EdgeInsets.all(14),
                  decoration: BoxDecoration(
                    color: cardDark,
                    border: Border.all(color: borderGray, width: 2),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      const Text(
                        'Mercedes-Benz Actros 3340 Heavy Lowbed',
                        style: TextStyle(
                          color: Color(0xFFF9FAFB),
                          fontWeight: FontWeight.w700,
                          fontSize: 15,
                        ),
                      ),
                      const SizedBox(height: 4),
                      const Text(
                        'Samuel Odhiambo • 14.2 km • 36,000 kg Payload',
                        style: TextStyle(color: Color(0xFFD1D5DB)),
                      ),
                      const SizedBox(height: 12),
                      SizedBox(
                        height: 56,
                        child: ElevatedButton(
                          onPressed: state is LogisticsLoading
                              ? null
                              : () {
                                  context.read<LogisticsBloc>().add(
                                        AcceptTransshipmentRequested(
                                          subBookingId: subBookingId,
                                          replacementFleetId:
                                              replacementFleetId,
                                        ),
                                      );
                                },
                          style: ElevatedButton.styleFrom(
                            backgroundColor: primaryYellow,
                            foregroundColor: Colors.black,
                            shape: const RoundedRectangleBorder(
                              borderRadius: BorderRadius.zero,
                            ),
                          ),
                          child: Text(
                            l10n.acceptHotSwapTransfer.toUpperCase(),
                            style: const TextStyle(fontWeight: FontWeight.w700),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),

                // Force Majeure Reverse Cascade Ledger
                Container(
                  padding: const EdgeInsets.all(14),
                  decoration: BoxDecoration(
                    color: cardDark,
                    border: Border.all(color: borderGray, width: 2),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      const Text(
                        'FORCE MAJEURE REVERSE CASCADE LEDGER',
                        style: TextStyle(
                          color: Color(0xFFF9FAFB),
                          fontWeight: FontWeight.w700,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 10),
                      _LedgerRow(
                        label: 'Renter Refund (100%)',
                        value: '${renterRefundMinor ~/ 100} $currency',
                      ),
                      _LedgerRow(
                        label: 'Owner Penalty',
                        value: '0 $currency (EXONERATED)',
                      ),
                      _LedgerRow(
                        label: 'Commission Restitution',
                        value: '${commissionRestitutionMinor ~/ 100} $currency',
                      ),
                      const SizedBox(height: 14),
                      SizedBox(
                        height: 56,
                        child: OutlinedButton(
                          onPressed: state is LogisticsLoading
                              ? null
                              : () {
                                  context.read<LogisticsBloc>().add(
                                        TriggerReverseCascadeRequested(
                                          parentBookingId: parentBookingId,
                                          reason:
                                              '60m hot-swap SLA expired without carrier recovery',
                                        ),
                                      );
                                },
                          style: OutlinedButton.styleFrom(
                            foregroundColor: safetyOrange,
                            side:
                                const BorderSide(color: safetyOrange, width: 2),
                            shape: const RoundedRectangleBorder(
                              borderRadius: BorderRadius.zero,
                            ),
                          ),
                          child: Text(
                            l10n.executeReverseCascade.toUpperCase(),
                            style: const TextStyle(fontWeight: FontWeight.w700),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),

                if (state is ReverseCascadeExecuted)
                  Container(
                    margin: const EdgeInsets.only(top: 16),
                    padding: const EdgeInsets.all(14),
                    decoration: BoxDecoration(
                      color: const Color(0xFF064E3B),
                      border: Border.all(
                          color: const Color(0xFF10B981), width: 2),
                    ),
                    child: Text(
                      '${state.cascadeResult.resolution} • PARENT: ${state.cascadeResult.parentStatus}',
                      style: const TextStyle(
                        color: Color(0xFFF9FAFB),
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ),
              ],
            ),
          );
        },
      ),
    );
  }
}

class _LedgerRow extends StatelessWidget {
  final String label;
  final String value;

  const _LedgerRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Color(0xFF9CA3AF))),
          Text(
            value,
            style: const TextStyle(
              color: Color(0xFFF9FAFB),
              fontWeight: FontWeight.w700,
            ),
          ),
        ],
      ),
    );
  }
}
