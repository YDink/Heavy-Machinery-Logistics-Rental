import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_bloc.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_event.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_state.dart';
import 'package:dink_mobile/features/discovery/presentation/views/fleet_comparison_view.dart';
import 'package:dink_mobile/features/discovery/presentation/views/widgets/fleet_feed_card.dart';
import 'package:dink_mobile/features/discovery/presentation/views/widgets/multi_facet_filter_sheet.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class DiscoveryView extends StatefulWidget {
  const DiscoveryView({super.key});

  @override
  State<DiscoveryView> createState() => _DiscoveryViewState();
}

class _DiscoveryViewState extends State<DiscoveryView> {
  final TextEditingController _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    // Default search centered around Addis Ababa (9.0300, 38.7400)
    context.read<DiscoveryBloc>().add(
          const SearchFleetsRequested(
            FleetSearchFilter(
              latitude: 9.0300,
              longitude: 38.7400,
              radiusKm: 50.0,
            ),
          ),
        );
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return BlocConsumer<DiscoveryBloc, DiscoveryState>(
      listener: (context, state) {
        if (state.status == DiscoveryStatus.compareSuccess && state.compareData != null) {
          Navigator.of(context).push(
            MaterialPageRoute(
              builder: (_) => FleetComparisonView(data: state.compareData!),
            ),
          );
        }
        if (state.status == DiscoveryStatus.failure && state.errorMessage != null) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(state.errorMessage!),
              backgroundColor: const Color(0xFFDC2626),
            ),
          );
        }
      },
      builder: (context, state) {
        final results = state.searchData?.results ?? [];

        return Scaffold(
          appBar: AppBar(
            backgroundColor: Colors.white,
            elevation: 1,
            title: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: l10n.discoverySearchHint,
                prefixIcon: const Icon(Icons.search, color: Color(0xFF64748B)),
                suffixIcon: _searchController.text.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          _onSearchSubmitted('');
                        },
                      )
                    : null,
                contentPadding: const EdgeInsets.symmetric(vertical: 0, horizontal: 16),
                filled: true,
                fillColor: const Color(0xFFF1F5F9),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(24),
                  borderSide: BorderSide.none,
                ),
              ),
              textInputAction: TextInputAction.search,
              onSubmitted: _onSearchSubmitted,
            ),
            actions: [
              IconButton(
                icon: const Icon(Icons.tune, color: Color(0xFF0F172A)),
                tooltip: l10n.filterTitle,
                onPressed: () => _openFilterSheet(context, state.filter),
              ),
            ],
          ),
          body: Stack(
            children: [
              // Results Feed
              if (state.status == DiscoveryStatus.loading)
                const Center(child: CircularProgressIndicator())
              else if (results.isEmpty)
                Center(
                  child: Padding(
                    padding: const EdgeInsets.all(24.0),
                    child: Text(
                      l10n.noFleetsFound,
                      textAlign: TextAlign.center,
                      style: const TextStyle(fontSize: 16, color: Color(0xFF64748B)),
                    ),
                  ),
                )
              else
                ListView.builder(
                  padding: EdgeInsets.only(
                    bottom: state.selectedForCompare.isNotEmpty ? 90 : 20,
                    top: 8,
                  ),
                  itemCount: results.length,
                  itemBuilder: (context, index) {
                    final item = results[index];
                    final isSelected = state.selectedForCompare.any((e) => e.fleetId == item.fleetId);
                    return FleetFeedCard(
                      item: item,
                      isSelectedForCompare: isSelected,
                      onCompareToggled: () {
                        context.read<DiscoveryBloc>().add(ToggleCompareSelection(item));
                      },
                      onTap: () {
                        // Navigate to detail dossier
                      },
                    );
                  },
                ),

              // Bottom Comparison Float Bar
              if (state.selectedForCompare.isNotEmpty)
                Positioned(
                  bottom: 16,
                  left: 16,
                  right: 16,
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                    decoration: BoxDecoration(
                      color: const Color(0xFF0F172A),
                      borderRadius: BorderRadius.circular(16),
                      boxShadow: const [
                        BoxShadow(
                          color: Colors.black38,
                          blurRadius: 10,
                          offset: Offset(0, 4),
                        ),
                      ],
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        TextButton(
                          onPressed: () {
                            context.read<DiscoveryBloc>().add(const ClearCompareSelection());
                          },
                          child: Text(
                            l10n.clearCompare,
                            style: const TextStyle(color: Color(0xFF94A3B8)),
                          ),
                        ),
                        ElevatedButton(
                          style: ElevatedButton.styleFrom(
                            backgroundColor: const Color(0xFFE05638),
                            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
                            padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
                          ),
                          onPressed: state.selectedForCompare.length >= 2
                              ? () {
                                  context.read<DiscoveryBloc>().add(const ExecuteComparisonRequested());
                                }
                              : null,
                          child: Text(
                            l10n.compareSelected(state.selectedForCompare.length),
                            style: const TextStyle(
                              color: Colors.white,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
            ],
          ),
        );
      },
    );
  }

  void _onSearchSubmitted(String query) {
    final current = context.read<DiscoveryBloc>().state.filter;
    context.read<DiscoveryBloc>().add(
          SearchFleetsRequested(
            current.copyWith(query: query.trim().isEmpty ? null : query.trim()),
          ),
        );
  }

  void _openFilterSheet(BuildContext context, FleetSearchFilter currentFilter) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (_) => MultiFacetFilterSheet(
        initialFilter: currentFilter,
        onApply: (newFilter) {
          context.read<DiscoveryBloc>().add(SearchFleetsRequested(newFilter));
        },
      ),
    );
  }
}

typedef DiscoveryPage = DiscoveryView;
