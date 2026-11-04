import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_result.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class FleetFeedCard extends StatelessWidget {
  final FleetSearchResult item;
  final bool isSelectedForCompare;
  final VoidCallback onCompareToggled;
  final VoidCallback onTap;

  const FleetFeedCard({
    super.key,
    required this.item,
    required this.isSelectedForCompare,
    required this.onCompareToggled,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final currencyFormat = NumberFormat.currency(
      locale: 'en_ET',
      symbol: 'ETB ',
      decimalDigits: 2,
    );
    final rateFormatted = currencyFormat.format(item.dailyRateMinor / 100.0);
    final distKm = (item.distanceMeters / 1000.0).toStringAsFixed(1);

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Top Image / Badges Container
            Stack(
              children: [
                AspectRatio(
                  aspectRatio: 16 / 9,
                  child: item.primaryPhotoUri != null && item.primaryPhotoUri!.isNotEmpty
                      ? Image.network(
                          item.primaryPhotoUri!,
                          fit: BoxFit.cover,
                          errorBuilder: (ctx, err, stack) => _buildImagePlaceholder(),
                        )
                      : _buildImagePlaceholder(),
                ),
                // Availability Badge
                Positioned(
                  top: 12,
                  right: 12,
                  child: _buildAvailabilityBadge(context, item.availabilityStatus, l10n),
                ),
                // Compare Checkbox (56x56dp Hit Target for WCAG 2.2 AAA Glove Mode)
                Positioned(
                  top: 4,
                  left: 4,
                  child: Semantics(
                    label: 'Toggle compare for ${item.make} ${item.model}',
                    child: InkWell(
                      onTap: onCompareToggled,
                      borderRadius: BorderRadius.circular(28),
                      child: Container(
                        width: 56,
                        height: 56,
                        alignment: Alignment.center,
                        child: Container(
                          width: 28,
                          height: 28,
                          decoration: BoxDecoration(
                            color: isSelectedForCompare
                                ? const Color(0xFFE05638)
                                : Colors.black.withValues(alpha: 0.6),
                            borderRadius: BorderRadius.circular(6),
                            border: Border.all(color: Colors.white, width: 2),
                          ),
                          child: isSelectedForCompare
                              ? const Icon(Icons.check, size: 20, color: Colors.white)
                              : null,
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ),

            // Card Body Details
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Title & Rating Row
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Expanded(
                        child: Text(
                          '${item.make} ${item.model} (${item.manufactureYear})',
                          style: const TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                            color: Color(0xFF1E293B),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      Row(
                        children: [
                          const Icon(Icons.star, size: 18, color: Color(0xFFF59E0B)),
                          const SizedBox(width: 4),
                          Text(
                            item.owner.averageRating.toStringAsFixed(1),
                            style: const TextStyle(
                              fontWeight: FontWeight.bold,
                              fontSize: 14,
                              color: Color(0xFF1E293B),
                            ),
                          ),
                          Text(
                            ' (${item.owner.totalCompletedRentals})',
                            style: const TextStyle(
                              fontSize: 12,
                              color: Color(0xFF64748B),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),

                  const SizedBox(height: 6),

                  // Distance & Depot Location
                  Row(
                    children: [
                      const Icon(Icons.location_on_outlined, size: 16, color: Color(0xFF64748B)),
                      const SizedBox(width: 4),
                      Text(
                        l10n.distanceAway(distKm),
                        style: const TextStyle(
                          fontSize: 13,
                          color: Color(0xFF64748B),
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ],
                  ),

                  const SizedBox(height: 10),

                  // Key Spec Badges
                  Wrap(
                    spacing: 6,
                    runSpacing: 6,
                    children: [
                      if (item.keySpecs.operatingWeightKg != null)
                        _buildSpecChip('${(item.keySpecs.operatingWeightKg! / 1000.0).toStringAsFixed(1)}t Weight'),
                      if (item.keySpecs.grossPowerHp != null)
                        _buildSpecChip('${item.keySpecs.grossPowerHp!.toInt()} HP'),
                      if (item.operatorIncluded)
                        _buildSpecChip(l10n.crewIncluded, isHighlight: true),
                      if (!item.selfMoveCapable)
                        _buildSpecChip(l10n.lowboyRequired),
                    ],
                  ),

                  const Divider(height: 24, thickness: 1),

                  // Price Footer
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        item.owner.tradingName,
                        style: const TextStyle(
                          fontSize: 13,
                          color: Color(0xFF475569),
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      RichText(
                        text: TextSpan(
                          children: [
                            TextSpan(
                              text: rateFormatted,
                              style: const TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.bold,
                                color: Color(0xFF0F172A),
                              ),
                            ),
                            TextSpan(
                              text: ' / ${l10n.perDay}',
                              style: const TextStyle(
                                fontSize: 13,
                                color: Color(0xFF64748B),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildImagePlaceholder() {
    return Container(
      color: const Color(0xFFE2E8F0),
      alignment: Alignment.center,
      child: const Icon(
        Icons.precision_manufacturing_outlined,
        size: 48,
        color: Color(0xFF94A3B8),
      ),
    );
  }

  Widget _buildAvailabilityBadge(
    BuildContext context,
    FleetAvailabilityStatus status,
    AppLocalizations l10n,
  ) {
    Color bg;
    Color fg;
    String text;

    switch (status) {
      case FleetAvailabilityStatus.availableNow:
        bg = const Color(0xFF16A34A);
        fg = Colors.white;
        text = l10n.availableNow;
        break;
      case FleetAvailabilityStatus.heldPendingWire:
        bg = const Color(0xFFD97706);
        fg = Colors.white;
        text = l10n.heldPendingWire;
        break;
      case FleetAvailabilityStatus.unavailable:
        bg = const Color(0xFFDC2626);
        fg = Colors.white;
        text = 'Unavailable';
        break;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: bg,
        borderRadius: BorderRadius.circular(6),
        boxShadow: const [
          BoxShadow(
            color: Colors.black26,
            blurRadius: 4,
            offset: Offset(0, 2),
          ),
        ],
      ),
      child: Text(
        text,
        style: TextStyle(
          color: fg,
          fontSize: 12,
          fontWeight: FontWeight.bold,
          letterSpacing: 0.3,
        ),
      ),
    );
  }

  Widget _buildSpecChip(String label, {bool isHighlight = false}) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: isHighlight ? const Color(0xFFEFF6FF) : const Color(0xFFF1F5F9),
        borderRadius: BorderRadius.circular(4),
        border: Border.all(
          color: isHighlight ? const Color(0xFF93C5FD) : const Color(0xFFCBD5E1),
          width: 1,
        ),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w600,
          color: isHighlight ? const Color(0xFF1D4ED8) : const Color(0xFF334155),
        ),
      ),
    );
  }
}
