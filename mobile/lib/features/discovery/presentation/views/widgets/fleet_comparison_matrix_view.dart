import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_compare_matrix.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class FleetComparisonMatrixView extends StatefulWidget {
  final FleetCompareData data;

  const FleetComparisonMatrixView({super.key, required this.data});

  @override
  State<FleetComparisonMatrixView> createState() => _FleetComparisonMatrixViewState();
}

class _FleetComparisonMatrixViewState extends State<FleetComparisonMatrixView> {
  bool _showDifferencesOnly = false;

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final currencyFormat = NumberFormat.currency(
      locale: 'en_ET',
      symbol: 'ETB ',
      decimalDigits: 0,
    );

    final filteredRows = _showDifferencesOnly
        ? widget.data.comparisonMatrix.where((r) => r.isDifferent).toList()
        : widget.data.comparisonMatrix;

    final columnCount = widget.data.fleets.length;

    return Column(
      children: [
        // Top Toolbar (Differences Only Toggle)
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          color: const Color(0xFFF8FAFC),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                widget.data.subcategoryNameEn,
                style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
              ),
              Row(
                children: [
                  Text(l10n.differencesOnly, style: const TextStyle(fontSize: 13)),
                  Switch(
                    value: _showDifferencesOnly,
                    activeTrackColor: const Color(0xFFE05638),
                    onChanged: (val) => setState(() => _showDifferencesOnly = val),
                  ),
                ],
              ),
            ],
          ),
        ),

        // Comparison Grid Table
        Expanded(
          child: SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: SizedBox(
              width: 140.0 + (columnCount * 160.0),
              child: ListView(
                children: [
                  // Fleet Header Row
                  Container(
                    color: Colors.white,
                    padding: const EdgeInsets.symmetric(vertical: 12),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const SizedBox(
                          width: 140,
                          child: Padding(
                            padding: EdgeInsets.all(8.0),
                            child: Text(
                              'Fleet Asset',
                              style: TextStyle(fontWeight: FontWeight.bold, color: Color(0xFF64748B)),
                            ),
                          ),
                        ),
                        ...widget.data.fleets.map(
                          (f) => SizedBox(
                            width: 160,
                            child: Padding(
                              padding: const EdgeInsets.symmetric(horizontal: 6.0),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  ClipRRect(
                                    borderRadius: BorderRadius.circular(8),
                                    child: AspectRatio(
                                      aspectRatio: 16 / 9,
                                      child: f.primaryPhotoUri != null && f.primaryPhotoUri!.isNotEmpty
                                          ? Image.network(f.primaryPhotoUri!, fit: BoxFit.cover)
                                          : Container(
                                              color: const Color(0xFFE2E8F0),
                                              child: const Icon(Icons.precision_manufacturing, size: 28),
                                            ),
                                    ),
                                  ),
                                  const SizedBox(height: 6),
                                  Text(
                                    f.title,
                                    style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
                                    maxLines: 2,
                                    overflow: TextOverflow.ellipsis,
                                  ),
                                  const SizedBox(height: 2),
                                  Text(
                                    currencyFormat.format(f.dailyRateMinor / 100.0),
                                    style: const TextStyle(
                                      fontWeight: FontWeight.w700,
                                      color: Color(0xFFE05638),
                                      fontSize: 13,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),

                  const Divider(height: 1, thickness: 1),

                  // Spec Rows
                  ...filteredRows.map((row) {
                    final isHighlighted = row.isDifferent;
                    return Container(
                      color: isHighlighted ? const Color(0xFFFEF3C7).withValues(alpha: 0.4) : Colors.white,
                      padding: const EdgeInsets.symmetric(vertical: 10),
                      child: Row(
                        children: [
                          // Label column
                          SizedBox(
                            width: 140,
                            child: Padding(
                              padding: const EdgeInsets.symmetric(horizontal: 10),
                              child: Text(
                                row.labelEn + (row.unitSymbol != null ? ' (${row.unitSymbol})' : ''),
                                style: const TextStyle(
                                  fontSize: 12,
                                  fontWeight: FontWeight.w600,
                                  color: Color(0xFF334155),
                                ),
                              ),
                            ),
                          ),

                          // Fleet Values columns
                          ...widget.data.fleets.map((f) {
                            final val = row.values[f.id] ?? '—';
                            return SizedBox(
                              width: 160,
                              child: Padding(
                                padding: const EdgeInsets.symmetric(horizontal: 6),
                                child: Text(
                                  val.isEmpty ? '—' : val,
                                  style: TextStyle(
                                    fontSize: 13,
                                    fontWeight: isHighlighted ? FontWeight.bold : FontWeight.normal,
                                    color: const Color(0xFF0F172A),
                                  ),
                                ),
                              ),
                            );
                          }),
                        ],
                      ),
                    );
                  }),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }
}
