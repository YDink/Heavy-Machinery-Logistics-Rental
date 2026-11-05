import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/rental_completion_bloc.dart';
import '../bloc/rental_completion_event.dart';
import '../bloc/rental_completion_state.dart';

class SettlementReconciliationView extends StatefulWidget {
  final String bookingId;
  final String bookingReference;
  final String equipmentTitle;
  final int baseFeeMinor;
  final int attachmentsFeeMinor;
  final int excessSurchargeMinor;
  final int acceptedDamageMinor;
  final int prepaymentMinor;
  final String currency;

  const SettlementReconciliationView({
    super.key,
    required this.bookingId,
    this.bookingReference = 'DNK-8812',
    this.equipmentTitle = 'CAT 320D Hydraulic Excavator',
    this.baseFeeMinor = 15750000,
    this.attachmentsFeeMinor = 5600000,
    this.excessSurchargeMinor = 1596000,
    this.acceptedDamageMinor = 850000,
    this.prepaymentMinor = 10675000,
    this.currency = 'ETB',
  });

  @override
  State<SettlementReconciliationView> createState() => _SettlementReconciliationViewState();
}

class _SettlementReconciliationViewState extends State<SettlementReconciliationView> {
  static const _cbeAccount = '1000289192831';
  static const _cbeBeneficiary = 'Kamau Ndegwa / Almaz Tesfaye';
  static const _paymentMemo = 'DNK-8812-SETTLE';

  int get _grossLiabilitiesMinor =>
      widget.baseFeeMinor +
      widget.attachmentsFeeMinor +
      widget.excessSurchargeMinor +
      widget.acceptedDamageMinor;

  int get _netBalanceMinor => _grossLiabilitiesMinor - widget.prepaymentMinor;

  String _formatMinor(int minor) {
    final major = minor / 100.0;
    return '${major.toStringAsFixed(2)} ${widget.currency}';
  }

  void _copyAccount() {
    Clipboard.setData(const ClipboardData(text: _cbeAccount));
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('CBE Account Number copied to clipboard!'),
        backgroundColor: Color(0xFF10B981),
      ),
    );
  }

  void _reconcileOnline() {
    context.read<RentalCompletionBloc>().add(
          ReconcileFinalSettlementEvent(bookingId: widget.bookingId),
        );
  }

  @override
  Widget build(BuildContext context) {
    const bgDark = Color(0xFF111827);
    const cardDark = Color(0xFF1C1C1C);
    const primaryYellow = Color(0xFFFFB300);
    const borderGray = Color(0xFF393939);
    const successGreen = Color(0xFF10B981);
    const alertOrange = Color(0xFFFF5722);

    return Scaffold(
      backgroundColor: bgDark,
      appBar: AppBar(
        backgroundColor: bgDark,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'SETTLEMENT STATEMENT',
              style: TextStyle(
                fontFamily: 'Space Grotesk',
                fontSize: 16,
                fontWeight: FontWeight.w700,
                color: Colors.white,
              ),
            ),
            Text(
              '${widget.bookingReference} • ${widget.equipmentTitle}',
              style: const TextStyle(fontSize: 11, color: Colors.white70),
            ),
          ],
        ),
        actions: [
          Container(
            margin: const EdgeInsets.only(right: 16, top: 12, bottom: 12),
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: alertOrange.withOpacity(0.2),
              border: Border.all(color: alertOrange),
              borderRadius: BorderRadius.zero,
            ),
            child: const Center(
              child: Text(
                'PENDING PAYMENT',
                style: TextStyle(
                  fontFamily: 'JetBrains Mono',
                  fontSize: 10,
                  fontWeight: FontWeight.w700,
                  color: alertOrange,
                ),
              ),
            ),
          ),
        ],
      ),
      body: BlocConsumer<RentalCompletionBloc, RentalCompletionState>(
        listener: (context, state) {
          if (state is SettlementReconciledState) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Settlement statement reconciled and ledger finalized!'),
                backgroundColor: successGreen,
              ),
            );
          } else if (state is RentalCompletionError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: const Color(0xFFEF4444),
              ),
            );
          }
        },
        builder: (context, state) {
          final isLoading = state is RentalCompletionLoading;

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Itemized Financial Ledger Card
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: cardDark,
                    border: Border.all(color: borderGray, width: 2),
                    borderRadius: BorderRadius.zero,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text(
                            'ITEMIZED SETTLEMENT LEDGER',
                            style: TextStyle(
                              fontFamily: 'Space Grotesk',
                              fontSize: 13,
                              fontWeight: FontWeight.w700,
                              color: Colors.white,
                              letterSpacing: 0.5,
                            ),
                          ),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                            color: const Color(0xFF262626),
                            child: const Text(
                              'RC-05 FORMULA',
                              style: TextStyle(
                                fontFamily: 'JetBrains Mono',
                                fontSize: 9,
                                color: Colors.white70,
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),

                      _buildLedgerRow('Base Machine Hire (7 Days @ 22,500.00 ETB)', widget.baseFeeMinor),
                      _buildLedgerRow('Bundled Implement (Rock Breaker)', widget.attachmentsFeeMinor),
                      _buildLedgerRow('Excess Shift Overtime (5.70 Excess Hrs)', widget.excessSurchargeMinor,
                          subtitle: 'Gross 61.70h - 56.00h shift quota (0h breakdown)'),
                      _buildLedgerRow('Accepted Damage (Cabin Safety Glass)', widget.acceptedDamageMinor,
                          subtitle: 'Docket #DMG-03 • Ries Engineering CAT OEM'),
                      const Divider(color: borderGray, height: 24, thickness: 1),

                      _buildLedgerRow('GROSS LIABILITIES TOTAL', _grossLiabilitiesMinor, isBold: true),
                      _buildLedgerRow('Less: 50% Prepayment Confirmed', -widget.prepaymentMinor,
                          isCredit: true),
                      const SizedBox(height: 16),

                      // Net Balance Due Box
                      Container(
                        width: double.infinity,
                        padding: const EdgeInsets.all(14),
                        decoration: BoxDecoration(
                          color: const Color(0xFF2A2207),
                          border: Border.all(color: primaryYellow, width: 2),
                          borderRadius: BorderRadius.zero,
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'NET BALANCE SETTLEMENT DUE TO OWNER',
                              style: TextStyle(
                                fontFamily: 'Space Grotesk',
                                fontSize: 11,
                                fontWeight: FontWeight.w700,
                                color: primaryYellow,
                                letterSpacing: 0.5,
                              ),
                            ),
                            const SizedBox(height: 6),
                            Text(
                              _formatMinor(_netBalanceMinor),
                              style: const TextStyle(
                                fontFamily: 'JetBrains Mono',
                                fontSize: 26,
                                fontWeight: FontWeight.w900,
                                color: primaryYellow,
                              ),
                            ),
                            const SizedBox(height: 4),
                            const Text(
                              'Direct P2P wire payment required upon off-hire handover sign-off.',
                              style: TextStyle(fontSize: 10, color: Colors.white70),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 20),

                // Owner Direct Settlement Coordinates Card
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: cardDark,
                    border: Border.all(color: borderGray),
                    borderRadius: BorderRadius.zero,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Row(
                        children: [
                          Icon(Icons.account_balance, size: 16, color: primaryYellow),
                          SizedBox(width: 8),
                          Text(
                            'OWNER DIRECT WIRE COORDINATES',
                            style: TextStyle(
                              fontFamily: 'Space Grotesk',
                              fontSize: 12,
                              fontWeight: FontWeight.w700,
                              color: Colors.white,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      const Text(
                        'Commercial Bank of Ethiopia (CBE)',
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        'Beneficiary: $_cbeBeneficiary',
                        style: const TextStyle(fontSize: 12, color: Colors.white70),
                      ),
                      const SizedBox(height: 8),
                      Row(
                        children: [
                          const Expanded(
                            child: Text(
                              'A/C: $_cbeAccount',
                              style: TextStyle(
                                fontFamily: 'JetBrains Mono',
                                fontSize: 15,
                                fontWeight: FontWeight.w700,
                                color: primaryYellow,
                              ),
                            ),
                          ),
                          SizedBox(
                            height: 36,
                            child: OutlinedButton.icon(
                              onPressed: _copyAccount,
                              icon: const Icon(Icons.copy, size: 14, color: primaryYellow),
                              label: const Text(
                                'COPY',
                                style: TextStyle(
                                  fontFamily: 'JetBrains Mono',
                                  fontSize: 11,
                                  color: primaryYellow,
                                ),
                              ),
                              style: OutlinedButton.styleFrom(
                                side: const BorderSide(color: primaryYellow),
                                shape: const RoundedRectangleBorder(borderRadius: BorderRadius.zero),
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 8),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        color: const Color(0xFF1E293B),
                        child: const Text(
                          'Payment Reference: $_paymentMemo (Mandatory)',
                          style: TextStyle(
                            fontFamily: 'JetBrains Mono',
                            fontSize: 11,
                            color: Colors.white70,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 20),

                // Wire Transfer Slip Upload Section
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: const Color(0xFF131313),
                    border: Border.all(color: borderGray, style: BorderStyle.solid),
                    borderRadius: BorderRadius.zero,
                  ),
                  child: Column(
                    children: [
                      const Icon(Icons.cloud_upload_outlined, size: 32, color: primaryYellow),
                      const SizedBox(height: 8),
                      const Text(
                        'ATTACH WIRE TRANSFER CONFIRMATION SLIP',
                        style: TextStyle(
                          fontFamily: 'Space Grotesk',
                          fontSize: 12,
                          fontWeight: FontWeight.w700,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(height: 4),
                      const Text(
                        'Supports PDF, JPEG, PNG (Max 15MB) • ClamAV Anti-Malware Scanned',
                        style: TextStyle(fontSize: 10, color: Colors.white60),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 24),

                // 56dp Primary Button
                SizedBox(
                  width: double.infinity,
                  height: 56,
                  child: ElevatedButton(
                    onPressed: isLoading ? null : _reconcileOnline,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: primaryYellow,
                      foregroundColor: Colors.black,
                      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.zero),
                    ),
                    child: isLoading
                        ? const CircularProgressIndicator(color: Colors.black)
                        : const Text(
                            'CONFIRM & RECONCILE SETTLEMENT →',
                            style: TextStyle(
                              fontFamily: 'Space Grotesk',
                              fontSize: 15,
                              fontWeight: FontWeight.w700,
                              letterSpacing: 0.5,
                            ),
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

  Widget _buildLedgerRow(String title, int amountMinor,
      {bool isBold = false, bool isCredit = false, String? subtitle}) {
    final amountText = _formatMinor(amountMinor.abs());
    final displayAmount = isCredit ? '- $amountText' : amountText;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Expanded(
                child: Text(
                  title,
                  style: TextStyle(
                    fontSize: isBold ? 13 : 12,
                    fontWeight: isBold ? FontWeight.w700 : FontWeight.w400,
                    color: isCredit ? const Color(0xFF10B981) : Colors.white,
                  ),
                ),
              ),
              Text(
                displayAmount,
                style: TextStyle(
                  fontFamily: 'JetBrains Mono',
                  fontSize: isBold ? 13 : 12,
                  fontWeight: isBold ? FontWeight.w700 : FontWeight.w500,
                  color: isCredit ? const Color(0xFF10B981) : Colors.white,
                ),
              ),
            ],
          ),
          if (subtitle != null)
            Text(
              subtitle,
              style: const TextStyle(fontSize: 10, color: Colors.white60),
            ),
        ],
      ),
    );
  }
}
