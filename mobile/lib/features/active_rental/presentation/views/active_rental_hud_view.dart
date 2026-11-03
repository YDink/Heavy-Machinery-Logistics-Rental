import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import '../bloc/active_rental_bloc.dart';
import '../bloc/active_rental_event.dart';
import '../bloc/active_rental_state.dart';
import 'widgets/breakdown_report_modal.dart';
import 'widgets/extension_request_sheet.dart';
import 'widgets/shift_telemetry_progress_card.dart';

class ActiveRentalHUDView extends StatefulWidget {
  final String bookingId;

  const ActiveRentalHUDView({super.key, required this.bookingId});

  @override
  State<ActiveRentalHUDView> createState() => _ActiveRentalHUDViewState();
}

class _ActiveRentalHUDViewState extends State<ActiveRentalHUDView> {
  @override
  void initState() {
    super.initState();
    context.read<ActiveRentalBloc>().add(LoadActiveRentalHUDEvent(widget.bookingId));
  }

  void _openBreakdownModal(double meter) {
    showModalBottomSheet(
      context: context,
      backgroundColor: const Color(0xFF1F2937),
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (_) => BlocProvider.value(
        value: context.read<ActiveRentalBloc>(),
        child: BreakdownReportModal(
          bookingId: widget.bookingId,
          currentMeterHours: meter,
        ),
      ),
    );
  }

  void _openExtensionSheet() {
    showModalBottomSheet(
      context: context,
      backgroundColor: const Color(0xFF1F2937),
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (_) => BlocProvider.value(
        value: context.read<ActiveRentalBloc>(),
        child: ExtensionRequestSheet(bookingId: widget.bookingId),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      backgroundColor: const Color(0xFF111827), // Dark high-contrast background
      appBar: AppBar(
        title: Text(
          l10n.activeRentalTitle,
          style: const TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
        ),
        backgroundColor: const Color(0xFF1F2937),
        elevation: 0,
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh, color: Colors.white),
            tooltip: 'Refresh',
            onPressed: () {
              context.read<ActiveRentalBloc>().add(LoadActiveRentalHUDEvent(widget.bookingId));
            },
          ),
        ],
      ),
      body: BlocConsumer<ActiveRentalBloc, ActiveRentalState>(
        listener: (context, state) {
          if (state is ActiveRentalActionSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: Colors.green.shade700,
              ),
            );
          } else if (state is ActiveRentalError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: Colors.red.shade700,
              ),
            );
          }
        },
        builder: (context, state) {
          if (state is ActiveRentalLoading) {
            return const Center(
              child: CircularProgressIndicator(color: Color(0xFF2563EB)),
            );
          }

          if (state is ActiveRentalError) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24.0),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Icon(Icons.error_outline, size: 48, color: Colors.redAccent),
                    const SizedBox(height: 16),
                    Text(
                      state.message,
                      textAlign: TextAlign.center,
                      style: const TextStyle(color: Colors.white, fontSize: 16),
                    ),
                    const SizedBox(height: 20),
                    ElevatedButton(
                      onPressed: () {
                        context
                            .read<ActiveRentalBloc>()
                            .add(LoadActiveRentalHUDEvent(widget.bookingId));
                      },
                      child: const Text('Retry'),
                    ),
                  ],
                ),
              ),
            );
          }

          if (state is ActiveRentalHUDLoaded) {
            final hud = state.hud;

            return SingleChildScrollView(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  // Fleet Header Card
                  Card(
                    color: const Color(0xFF1F2937),
                    elevation: 2,
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Expanded(
                                child: Text(
                                  hud.fleetTitle,
                                  style: const TextStyle(
                                    fontSize: 20,
                                    fontWeight: FontWeight.bold,
                                    color: Colors.white,
                                  ),
                                ),
                              ),
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                                decoration: BoxDecoration(
                                  color: const Color(0xFF2563EB),
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                child: Text(
                                  'Day ${hud.currentRentalDay} / ${hud.totalRentalDays}',
                                  style: const TextStyle(
                                    color: Colors.white,
                                    fontWeight: FontWeight.bold,
                                    fontSize: 12,
                                  ),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 6),
                          Text(
                            hud.plateNumber,
                            style: TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w500,
                              color: Colors.grey.shade400,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Operator Card (if present)
                  if (hud.operator != null) ...[
                    Card(
                      color: const Color(0xFF1F2937),
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                      child: ListTile(
                        leading: const CircleAvatar(
                          backgroundColor: Color(0xFF374151),
                          child: Icon(Icons.person, color: Colors.white70),
                        ),
                        title: Text(
                          hud.operator!.fullName,
                          style: const TextStyle(fontWeight: FontWeight.bold, color: Colors.white),
                        ),
                        subtitle: Text(
                          '${hud.operator!.licenseClass} • ${hud.operator!.phoneNumber}',
                          style: TextStyle(color: Colors.grey.shade400),
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),
                  ],

                  // Shift Progress Card (8.0 hr Daily Quota)
                  ShiftTelemetryProgressCard(hud: hud),
                  const SizedBox(height: 24),

                  // Quick Action Buttons (56dp minimum height per accessibility rule)
                  Row(
                    children: [
                      // Report Breakdown
                      Expanded(
                        child: SizedBox(
                          height: 56,
                          child: ElevatedButton.icon(
                            onPressed: hud.canReportBreakdown
                                ? () => _openBreakdownModal(
                                      hud.handoverBaselineHours + hud.accumulatedHoursToday,
                                    )
                                : null,
                            icon: const Icon(Icons.warning_amber_rounded, size: 20),
                            label: Text(
                              l10n.quickActionBreakdown,
                              style: const TextStyle(fontWeight: FontWeight.bold),
                            ),
                            style: ElevatedButton.styleFrom(
                              backgroundColor: Colors.amber.shade900,
                              foregroundColor: Colors.white,
                              disabledBackgroundColor: Colors.grey.shade800,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      // Extend Rental
                      Expanded(
                        child: SizedBox(
                          height: 56,
                          child: ElevatedButton.icon(
                            onPressed: hud.canExtend ? _openExtensionSheet : null,
                            icon: const Icon(Icons.update, size: 20),
                            label: Text(
                              l10n.quickActionExtend,
                              style: const TextStyle(fontWeight: FontWeight.bold),
                            ),
                            style: ElevatedButton.styleFrom(
                              backgroundColor: const Color(0xFF2563EB),
                              foregroundColor: Colors.white,
                              disabledBackgroundColor: Colors.grey.shade800,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
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
