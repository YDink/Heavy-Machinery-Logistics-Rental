import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/winching_photo.dart';
import '../bloc/logistics_bloc.dart';
import '../bloc/logistics_event.dart';
import '../bloc/logistics_state.dart';

class CarrierDispatchView extends StatefulWidget {
  final String subBookingId;
  final String routeSummary;
  final VoidCallback? onOpenTransshipment;

  const CarrierDispatchView({
    super.key,
    required this.subBookingId,
    this.routeSummary = 'Kality Industrial Depot -> Mojo Highway Site B',
    this.onOpenTransshipment,
  });

  @override
  State<CarrierDispatchView> createState() => _CarrierDispatchViewState();
}

class _CarrierDispatchViewState extends State<CarrierDispatchView> {
  final Map<String, String> _capturedWinching = {
    'RIG_PROFILE': 's3://dink-media/winch/rig_profile.webp',
    'CHAIN_BINDERS': 's3://dink-media/winch/chain_binders.webp',
    'BOOM_PINNED': 's3://dink-media/winch/boom_pinned.webp',
    'RAMPS_LOCKED': 's3://dink-media/winch/ramps_locked.webp',
  };

  static const _angles = [
    {'key': 'RIG_PROFILE', 'label': '01. RIG PROFILE'},
    {'key': 'CHAIN_BINDERS', 'label': '02. CHAIN BINDERS (4-PT)'},
    {'key': 'BOOM_PINNED', 'label': '03. BOOM PINNED'},
    {'key': 'RAMPS_LOCKED', 'label': '04. RAMPS LOCKED'},
  ];

  void _confirmLoadAndStartTransit() {
    final photos = _capturedWinching.entries
        .map((e) => WinchingPhoto(angle: e.key, storageUri: e.value))
        .toList();

    context.read<LogisticsBloc>().add(
          AdvanceLogisticsDispatchRequested(
            subBookingId: widget.subBookingId,
            targetStage: 'LOADING_INSPECTION',
            winchingPhotos: photos,
            notes: '4-point lashing and boom pin verified.',
          ),
        );
  }

  void _triggerBreakdownAlert() {
    if (widget.onOpenTransshipment != null) {
      widget.onOpenTransshipment!();
      return;
    }
    context.read<LogisticsBloc>().add(
          ReportTransshipmentAlertRequested(
            subBookingId: widget.subBookingId,
            latitude: 8.7500,
            longitude: 38.9000,
            failureCategory: 'PRIME_MOVER_ENGINE_BLOWN',
            failureDescription:
                'Turbocharger shaft sheared on expressway shoulder.',
            photos: const ['s3://dink-media/transshipment/engine.webp'],
          ),
        );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    const bgDark = Color(0xFF111827);
    const cardDark = Color(0xFF1C1C1C);
    const borderGray = Color(0xFF3E3E3E);
    const primaryYellow = Color(0xFFFFB300);
    const safetyOrange = Color(0xFFFF5722);

    return Scaffold(
      backgroundColor: bgDark,
      appBar: AppBar(
        backgroundColor: bgDark,
        elevation: 0,
        title: Text(
          l10n.carrierDispatchTitle.toUpperCase(),
          style: const TextStyle(
            color: Color(0xFFF9FAFB),
            fontWeight: FontWeight.w700,
          ),
        ),
      ),
      body: BlocBuilder<LogisticsBloc, LogisticsState>(
        builder: (context, state) {
          final statusText = state is LogisticsDispatchUpdated
              ? state.dispatchResult.subBookingStatus
              : 'LOADING_INSPECTION';

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Route & Status Card
                Container(
                  padding: const EdgeInsets.all(14),
                  decoration: BoxDecoration(
                    color: cardDark,
                    border: Border.all(color: borderGray, width: 2),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'STAGE: $statusText',
                        style: const TextStyle(
                          color: primaryYellow,
                          fontWeight: FontWeight.w700,
                          fontSize: 13,
                        ),
                      ),
                      const SizedBox(height: 6),
                      Text(
                        widget.routeSummary,
                        style: const TextStyle(
                          color: Color(0xFFF9FAFB),
                          fontWeight: FontWeight.w600,
                          fontSize: 15,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),

                // 4-Angle Winching Reticle Header
                Text(
                  l10n.winchingReticleHeader.toUpperCase(),
                  style: const TextStyle(
                    color: Color(0xFFF9FAFB),
                    fontWeight: FontWeight.w700,
                    fontSize: 13,
                    letterSpacing: 0.8,
                  ),
                ),
                const SizedBox(height: 10),

                // 2x2 Reticle Grid
                GridView.count(
                  crossAxisCount: 2,
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  mainAxisSpacing: 8,
                  crossAxisSpacing: 8,
                  childAspectRatio: 1.4,
                  children: _angles.map((angle) {
                    final key = angle['key']!;
                    final captured = _capturedWinching.containsKey(key);
                    return InkWell(
                      onTap: () {
                        setState(() {
                          _capturedWinching[key] =
                              's3://dink-media/winch/${key.toLowerCase()}.webp';
                        });
                      },
                      child: Container(
                        padding: const EdgeInsets.all(10),
                        decoration: BoxDecoration(
                          color: cardDark,
                          border: Border.all(
                            color: captured
                                ? const Color(0xFF10B981)
                                : primaryYellow,
                            width: 2,
                          ),
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(
                              captured
                                  ? Icons.verified
                                  : Icons.camera_alt_outlined,
                              color: captured
                                  ? const Color(0xFF10B981)
                                  : primaryYellow,
                            ),
                            const SizedBox(height: 6),
                            Text(
                              angle['label']!,
                              textAlign: TextAlign.center,
                              style: const TextStyle(
                                color: Color(0xFFF9FAFB),
                                fontSize: 11,
                                fontWeight: FontWeight.w700,
                              ),
                            ),
                          ],
                        ),
                      ),
                    );
                  }).toList(),
                ),
                const SizedBox(height: 20),

                // Primary 56dp CTA: Confirm Load & Start Transit
                SizedBox(
                  height: 56,
                  child: ElevatedButton(
                    onPressed: state is LogisticsLoading
                        ? null
                        : _confirmLoadAndStartTransit,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: primaryYellow,
                      foregroundColor: Colors.black,
                      shape: const RoundedRectangleBorder(
                        borderRadius: BorderRadius.zero,
                      ),
                    ),
                    child: state is LogisticsLoading
                        ? const CircularProgressIndicator(color: Colors.black)
                        : Text(
                            l10n.confirmLoadStartTransit.toUpperCase(),
                            style: const TextStyle(
                              fontWeight: FontWeight.w700,
                              letterSpacing: 0.8,
                            ),
                          ),
                  ),
                ),
                const SizedBox(height: 12),

                // Secondary 56dp Emergency CTA: Report Roadside Breakdown
                SizedBox(
                  height: 56,
                  child: OutlinedButton(
                    onPressed: state is LogisticsLoading
                        ? null
                        : _triggerBreakdownAlert,
                    style: OutlinedButton.styleFrom(
                      foregroundColor: safetyOrange,
                      side: const BorderSide(color: safetyOrange, width: 2),
                      shape: const RoundedRectangleBorder(
                        borderRadius: BorderRadius.zero,
                      ),
                    ),
                    child: Text(
                      l10n.reportRoadsideBreakdown.toUpperCase(),
                      style: const TextStyle(
                        fontWeight: FontWeight.w700,
                        letterSpacing: 0.8,
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
}
