import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/rental_completion_bloc.dart';
import '../bloc/rental_completion_event.dart';
import '../bloc/rental_completion_state.dart';

class DamageDocketView extends StatefulWidget {
  final String bookingId;
  final String claimId;
  final String bookingReference;
  final String equipmentTitle;
  final String componentCategory;
  final int claimedAmountMinor;
  final String currency;
  final String narrative;

  const DamageDocketView({
    super.key,
    required this.bookingId,
    required this.claimId,
    this.bookingReference = 'DNK-8812',
    this.equipmentTitle = 'CAT 320D Hydraulic Excavator',
    this.componentCategory = 'CABIN SAFETY GLASS & ROPS FRAME',
    this.claimedAmountMinor = 850000,
    this.currency = 'ETB',
    this.narrative = 'Rock blast spall cracked front upper laminated windshield and distorted wiper pivot assembly during night trenching at Sector 4-B.',
  });

  @override
  State<DamageDocketView> createState() => _DamageDocketViewState();
}

class _DamageDocketViewState extends State<DamageDocketView> {
  final _disputeReasonController = TextEditingController();
  bool _showDisputeForm = false;

  void _acceptClaim() {
    context.read<RentalCompletionBloc>().add(
          RespondDamageClaimEvent(
            bookingId: widget.bookingId,
            claimId: widget.claimId,
            action: 'ACCEPT',
            acceptedAmountMinor: widget.claimedAmountMinor,
          ),
        );
  }

  void _disputeClaim() {
    final reason = _disputeReasonController.text.trim();
    if (reason.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please provide factual dispute rationale and baseline reference.'),
          backgroundColor: Color(0xFFEF4444),
        ),
      );
      return;
    }

    context.read<RentalCompletionBloc>().add(
          RespondDamageClaimEvent(
            bookingId: widget.bookingId,
            claimId: widget.claimId,
            action: 'DISPUTE',
            disputeReason: reason,
          ),
        );
  }

  @override
  void dispose() {
    _disputeReasonController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    const bgDark = Color(0xFF111827);
    const cardDark = Color(0xFF1C1C1C);
    const primaryYellow = Color(0xFFFFB300);
    const borderGray = Color(0xFF393939);
    const alertOrange = Color(0xFFFF5722);
    const successGreen = Color(0xFF10B981);
    const alertRed = Color(0xFFEF4444);

    final majorAmount = (widget.claimedAmountMinor / 100.0).toStringAsFixed(2);

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
              'DAMAGE CLAIM DOCKET #DMG-03',
              style: TextStyle(
                fontFamily: 'Space Grotesk',
                fontSize: 15,
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
      ),
      body: BlocConsumer<RentalCompletionBloc, RentalCompletionState>(
        listener: (context, state) {
          if (state is DamageClaimRespondedState) {
            final isAccepted = state.result.status == 'ACCEPTED';
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(isAccepted
                    ? 'Damage liability accepted! Injected into final settlement.'
                    : 'Claim disputed! Dossier escalated to Compliance Officer arbitration.'),
                backgroundColor: isAccepted ? successGreen : alertOrange,
              ),
            );
          } else if (state is RentalCompletionError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: alertRed,
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
                // 24-Hour Statutory Countdown Timer Bar
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: const Color(0xFF3B1E08),
                    border: Border.all(color: alertOrange, width: 1.5),
                    borderRadius: BorderRadius.zero,
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.timer_outlined, size: 20, color: alertOrange),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'STATUTORY DAMAGE RESPONSE DEADLINE',
                              style: TextStyle(
                                fontFamily: 'Space Grotesk',
                                fontSize: 10,
                                fontWeight: FontWeight.w700,
                                color: alertOrange,
                              ),
                            ),
                            const Text(
                              '18h 42m REMAINING BEFORE AUTO-WAIVER',
                              style: TextStyle(
                                fontFamily: 'JetBrains Mono',
                                fontSize: 14,
                                fontWeight: FontWeight.w800,
                                color: Colors.white,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),

                // Component Category & Valuation Card
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
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        color: const Color(0xFF262626),
                        child: Text(
                          widget.componentCategory,
                          style: const TextStyle(
                            fontFamily: 'JetBrains Mono',
                            fontSize: 11,
                            fontWeight: FontWeight.w700,
                            color: primaryYellow,
                          ),
                        ),
                      ),
                      const SizedBox(height: 12),
                      const Text(
                        'CLAIMED COMPENSATION AMOUNT',
                        style: TextStyle(fontSize: 11, color: Colors.white60),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '$majorAmount ${widget.currency}',
                        style: const TextStyle(
                          fontFamily: 'JetBrains Mono',
                          fontSize: 26,
                          fontWeight: FontWeight.w900,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(height: 12),
                      const Text(
                        'INCIDENT NARRATIVE',
                        style: TextStyle(
                          fontFamily: 'Space Grotesk',
                          fontSize: 11,
                          fontWeight: FontWeight.w700,
                          color: Colors.white70,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        widget.narrative,
                        style: const TextStyle(fontSize: 13, color: Colors.white, height: 1.4),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),

                // Certified Dealer Quote Card
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(14),
                  decoration: BoxDecoration(
                    color: cardDark,
                    border: Border.all(color: borderGray),
                    borderRadius: BorderRadius.zero,
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.picture_as_pdf, color: alertRed, size: 28),
                      const SizedBox(width: 12),
                      const Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'RIES ENGINEERING SHARE CO.',
                              style: TextStyle(
                                fontFamily: 'Space Grotesk',
                                fontSize: 12,
                                fontWeight: FontWeight.w700,
                                color: Colors.white,
                              ),
                            ),
                            Text(
                              'Caterpillar Authorized Sole Dealer • Quote #RE-2026-9921',
                              style: TextStyle(fontSize: 11, color: Colors.white60),
                            ),
                          ],
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.open_in_new, color: primaryYellow, size: 20),
                        onPressed: () {},
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 24),

                // Contractor Resolution Action Buttons (56dp)
                SizedBox(
                  width: double.infinity,
                  height: 56,
                  child: ElevatedButton(
                    onPressed: isLoading ? null : _acceptClaim,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: primaryYellow,
                      foregroundColor: Colors.black,
                      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.zero),
                    ),
                    child: Text(
                      'ACCEPT LIABILITY ($majorAmount ${widget.currency})',
                      style: const TextStyle(
                        fontFamily: 'Space Grotesk',
                        fontSize: 14,
                        fontWeight: FontWeight.w800,
                        letterSpacing: 0.5,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 12),

                if (!_showDisputeForm)
                  SizedBox(
                    width: double.infinity,
                    height: 56,
                    child: OutlinedButton(
                      onPressed: () => setState(() => _showDisputeForm = true),
                      style: OutlinedButton.styleFrom(
                        side: const BorderSide(color: alertRed, width: 1.5),
                        shape: const RoundedRectangleBorder(borderRadius: BorderRadius.zero),
                      ),
                      child: const Text(
                        'DISPUTE CLAIM & ESCALATE TO ARBITRATION',
                        style: TextStyle(
                          fontFamily: 'Space Grotesk',
                          fontSize: 13,
                          fontWeight: FontWeight.w700,
                          color: alertRed,
                        ),
                      ),
                    ),
                  )
                else ...[
                  const SizedBox(height: 8),
                  TextField(
                    controller: _disputeReasonController,
                    maxLines: 3,
                    style: const TextStyle(color: Colors.white, fontSize: 13),
                    decoration: const InputDecoration(
                      labelText: 'DISPUTE RATIONALE & BASELINE REFERENCE',
                      labelStyle: TextStyle(color: Colors.white70),
                      hintText: 'e.g., Pre-existing micro-fracture documented in Dispatch Reticle #IN-4890...',
                      hintStyle: TextStyle(color: Colors.white30),
                      filled: true,
                      fillColor: cardDark,
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.zero,
                        borderSide: BorderSide(color: alertRed),
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  SizedBox(
                    width: double.infinity,
                    height: 56,
                    child: ElevatedButton(
                      onPressed: isLoading ? null : _disputeClaim,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: alertRed,
                        foregroundColor: Colors.white,
                        shape: const RoundedRectangleBorder(borderRadius: BorderRadius.zero),
                      ),
                      child: const Text(
                        'SUBMIT FORMAL DISPUTE (HOLD ESCROW)',
                        style: TextStyle(
                          fontFamily: 'Space Grotesk',
                          fontSize: 14,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ),
                  ),
                ],
              ],
            ),
          );
        },
      ),
    );
  }
}
