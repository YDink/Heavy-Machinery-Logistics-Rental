import 'package:flutter/material.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class MultiFacetFilterSheet extends StatefulWidget {
  final FleetSearchFilter initialFilter;
  final ValueChanged<FleetSearchFilter> onApply;

  const MultiFacetFilterSheet({
    super.key,
    required this.initialFilter,
    required this.onApply,
  });

  @override
  State<MultiFacetFilterSheet> createState() => _MultiFacetFilterSheetState();
}

class _MultiFacetFilterSheetState extends State<MultiFacetFilterSheet> {
  late double _radiusKm;
  late RangeValues _priceRange;
  late bool? _operatorIncluded;
  late bool? _selfMoveCapable;
  late String? _selectedFuelType;
  late DiscoverySortBy _sortBy;

  final double _minAllowedPrice = 1000.0;
  final double _maxAllowedPrice = 100000.0;

  @override
  void initState() {
    super.initState();
    _radiusKm = widget.initialFilter.radiusKm;
    final minP = (widget.initialFilter.minRateMinor != null)
        ? widget.initialFilter.minRateMinor! / 100.0
        : _minAllowedPrice;
    final maxP = (widget.initialFilter.maxRateMinor != null)
        ? widget.initialFilter.maxRateMinor! / 100.0
        : _maxAllowedPrice;
    _priceRange = RangeValues(minP.clamp(_minAllowedPrice, _maxAllowedPrice), maxP.clamp(_minAllowedPrice, _maxAllowedPrice));
    _operatorIncluded = widget.initialFilter.operatorIncluded;
    _selfMoveCapable = widget.initialFilter.selfMoveCapable;
    _selectedFuelType = widget.initialFilter.fuelType;
    _sortBy = widget.initialFilter.sortBy;
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Container(
      padding: EdgeInsets.only(
        top: 20,
        left: 20,
        right: 20,
        bottom: MediaQuery.of(context).viewInsets.bottom + 20,
      ),
      decoration: const BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      child: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Center(
              child: Container(
                width: 40,
                height: 4,
                margin: const EdgeInsets.only(bottom: 16),
                decoration: BoxDecoration(
                  color: const Color(0xFFCBD5E1),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  l10n.filterTitle,
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: Color(0xFF0F172A),
                  ),
                ),
                TextButton(
                  onPressed: _resetFilters,
                  child: Text(
                    l10n.resetFilters,
                    style: const TextStyle(color: Color(0xFFE05638)),
                  ),
                ),
              ],
            ),

            const SizedBox(height: 16),

            Text(
              'Search Radius: ${_radiusKm.toInt()} km',
              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
            ),
            Slider(
              value: _radiusKm,
              min: 5.0,
              max: 500.0,
              divisions: 99,
              activeColor: const Color(0xFFE05638),
              onChanged: (val) => setState(() => _radiusKm = val),
            ),

            const SizedBox(height: 12),

            Text(
              'Daily Rate: ETB ${_priceRange.start.toInt()} - ETB ${_priceRange.end.toInt()}',
              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
            ),
            RangeSlider(
              values: _priceRange,
              min: _minAllowedPrice,
              max: _maxAllowedPrice,
              divisions: 100,
              activeColor: const Color(0xFFE05638),
              onChanged: (val) => setState(() => _priceRange = val),
            ),

            const SizedBox(height: 12),

            SwitchListTile(
              contentPadding: EdgeInsets.zero,
              title: Text(l10n.crewIncluded, style: const TextStyle(fontWeight: FontWeight.w600)),
              value: _operatorIncluded ?? false,
              activeTrackColor: const Color(0xFFE05638),
              onChanged: (val) => setState(() => _operatorIncluded = val ? true : null),
            ),

            SwitchListTile(
              contentPadding: EdgeInsets.zero,
              title: const Text('Self-Move Capable (No Lowboy Required)', style: TextStyle(fontWeight: FontWeight.w600)),
              value: _selfMoveCapable ?? false,
              activeTrackColor: const Color(0xFFE05638),
              onChanged: (val) => setState(() => _selfMoveCapable = val ? true : null),
            ),

            const SizedBox(height: 12),

            const Text('Fuel System', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              children: ['DIESEL', 'ELECTRIC', 'HYBRID', 'PETROL'].map((fuel) {
                final isSelected = _selectedFuelType == fuel;
                return ChoiceChip(
                  label: Text(fuel),
                  selected: isSelected,
                  selectedColor: const Color(0xFFE05638).withValues(alpha: 0.2),
                  onSelected: (val) {
                    setState(() {
                      _selectedFuelType = val ? fuel : null;
                    });
                  },
                );
              }).toList(),
            ),

            const SizedBox(height: 24),

            SizedBox(
              width: double.infinity,
              height: 52,
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFFE05638),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
                ),
                onPressed: _applyFilters,
                child: Text(
                  l10n.applyFilters,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _resetFilters() {
    setState(() {
      _radiusKm = 50.0;
      _priceRange = RangeValues(_minAllowedPrice, _maxAllowedPrice);
      _operatorIncluded = null;
      _selfMoveCapable = null;
      _selectedFuelType = null;
      _sortBy = DiscoverySortBy.recommended;
    });
  }

  void _applyFilters() {
    final updated = widget.initialFilter.copyWith(
      radiusKm: _radiusKm,
      minRateMinor: (_priceRange.start * 100.0).toInt(),
      maxRateMinor: (_priceRange.end * 100.0).toInt(),
      operatorIncluded: _operatorIncluded,
      selfMoveCapable: _selfMoveCapable,
      fuelType: _selectedFuelType,
      sortBy: _sortBy,
    );
    widget.onApply(updated);
    Navigator.of(context).pop();
  }
}
