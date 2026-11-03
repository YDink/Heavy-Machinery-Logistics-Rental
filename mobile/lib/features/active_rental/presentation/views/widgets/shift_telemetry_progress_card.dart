import 'package:flutter/material.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:dink_mobile/features/active_rental/domain/entities/active_rental_hud.dart';

class ShiftTelemetryProgressCard extends StatelessWidget {
  final ActiveRentalHUD hud;

  const ShiftTelemetryProgressCard({super.key, required this.hud});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final pacingRatio = hud.shiftPacingRatio.clamp(0.0, 1.0);
    final isExcess = hud.accumulatedHoursToday > hud.allowedHoursPerDay;

    return Card(
      color: const Color(0xFF1F2937), // Dark slate gray
      elevation: 4,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(20.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    l10n.shiftProgressTitle,
                    style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: Colors.white,
                    ),
                  ),
                ),
                if (hud.shiftTimerPaused) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.amber.shade900.withValues(alpha: 0.8),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      l10n.shiftTimerPausedNotice,
                      style: const TextStyle(
                        fontSize: 11,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                    ),
                  ),
                ],
              ],
            ),
            const SizedBox(height: 16),
            // Progress Bar
            ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: LinearProgressIndicator(
                value: pacingRatio,
                minHeight: 14,
                backgroundColor: Colors.grey.shade800,
                valueColor: AlwaysStoppedAnimation<Color>(
                  isExcess ? Colors.deepOrangeAccent : Colors.lightGreenAccent.shade400,
                ),
              ),
            ),
            const SizedBox(height: 14),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      l10n.accumulatedTodayLabel,
                      style: TextStyle(fontSize: 12, color: Colors.grey.shade400),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      '${hud.accumulatedHoursToday.toStringAsFixed(1)} / ${hud.allowedHoursPerDay.toStringAsFixed(1)} hrs',
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                    ),
                  ],
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      l10n.excessHoursLabel,
                      style: TextStyle(
                        fontSize: 12,
                        color: isExcess ? Colors.deepOrangeAccent : Colors.grey.shade400,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      '${hud.excessHoursTotal.toStringAsFixed(1)} hrs',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                        color: isExcess ? Colors.deepOrangeAccent : Colors.white,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
