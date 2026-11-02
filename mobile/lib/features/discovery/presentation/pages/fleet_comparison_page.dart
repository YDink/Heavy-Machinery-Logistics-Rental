import 'package:flutter/material.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_compare_matrix.dart';
import 'package:dink_mobile/features/discovery/presentation/widgets/fleet_comparison_matrix_view.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class FleetComparisonPage extends StatelessWidget {
  final FleetCompareData data;

  const FleetComparisonPage({super.key, required this.data});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.comparisonMatrixTitle),
        backgroundColor: Colors.white,
        foregroundColor: const Color(0xFF0F172A),
        elevation: 1,
      ),
      body: FleetComparisonMatrixView(data: data),
    );
  }
}
