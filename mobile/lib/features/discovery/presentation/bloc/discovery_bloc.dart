import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/core/errors/failure.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/features/discovery/domain/usecases/compare_fleets_usecase.dart';
import 'package:dink_mobile/features/discovery/domain/usecases/get_fleet_availability_usecase.dart';
import 'package:dink_mobile/features/discovery/domain/usecases/search_fleets_usecase.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_event.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_state.dart';

class DiscoveryBloc extends Bloc<DiscoveryEvent, DiscoveryState> {
  final SearchFleetsUseCase searchFleetsUseCase;
  final CompareFleetsUseCase compareFleetsUseCase;
  final GetFleetAvailabilityUseCase getFleetAvailabilityUseCase;

  DiscoveryBloc({
    required this.searchFleetsUseCase,
    required this.compareFleetsUseCase,
    required this.getFleetAvailabilityUseCase,
    DiscoveryState initialState = const DiscoveryState(
      filter: FleetSearchFilter(latitude: 9.0300, longitude: 38.7400),
    ),
  }) : super(initialState) {
    on<SearchFleetsRequested>(_onSearchFleets);
    on<ToggleViewModeRequested>(_onToggleViewMode);
    on<ToggleCompareSelection>(_onToggleCompareSelection);
    on<ClearCompareSelection>(_onClearCompareSelection);
    on<ExecuteComparisonRequested>(_onExecuteComparison);
    on<CheckAvailabilityRequested>(_onCheckAvailability);
  }

  Future<void> _onSearchFleets(
    SearchFleetsRequested event,
    Emitter<DiscoveryState> emit,
  ) async {
    emit(state.copyWith(
      status: DiscoveryStatus.loading,
      filter: event.filter,
    ));

    try {
      final data = await searchFleetsUseCase(event.filter);
      emit(state.copyWith(
        status: DiscoveryStatus.success,
        searchData: data,
        errorMessage: null,
      ));
    } catch (e) {
      final msg = e is Failure ? e.message : e.toString();
      emit(state.copyWith(
        status: DiscoveryStatus.failure,
        errorMessage: msg,
      ));
    }
  }

  void _onToggleViewMode(
    ToggleViewModeRequested event,
    Emitter<DiscoveryState> emit,
  ) {
    emit(state.copyWith(isMapView: event.isMapView));
  }

  void _onToggleCompareSelection(
    ToggleCompareSelection event,
    Emitter<DiscoveryState> emit,
  ) {
    final current = List.of(state.selectedForCompare);
    final idx = current.indexWhere((item) => item.fleetId == event.item.fleetId);

    if (idx >= 0) {
      current.removeAt(idx);
    } else {
      // Enforce comparison limit of 3 (SR-08)
      if (current.length >= 3) {
        return;
      }
      current.add(event.item);
    }

    emit(state.copyWith(selectedForCompare: current));
  }

  void _onClearCompareSelection(
    ClearCompareSelection event,
    Emitter<DiscoveryState> emit,
  ) {
    emit(state.copyWith(selectedForCompare: const []));
  }

  Future<void> _onExecuteComparison(
    ExecuteComparisonRequested event,
    Emitter<DiscoveryState> emit,
  ) async {
    if (state.selectedForCompare.length < 2) {
      return;
    }

    emit(state.copyWith(status: DiscoveryStatus.comparing));

    final ids = state.selectedForCompare.map((e) => e.fleetId).toList();

    try {
      final compareResult = await compareFleetsUseCase(ids);
      emit(state.copyWith(
        status: DiscoveryStatus.compareSuccess,
        compareData: compareResult,
        errorMessage: null,
      ));
    } catch (e) {
      final msg = e is Failure ? e.message : e.toString();
      emit(state.copyWith(
        status: DiscoveryStatus.failure,
        errorMessage: msg,
      ));
    }
  }

  Future<void> _onCheckAvailability(
    CheckAvailabilityRequested event,
    Emitter<DiscoveryState> emit,
  ) async {
    emit(state.copyWith(status: DiscoveryStatus.checkingAvailability));

    try {
      final result = await getFleetAvailabilityUseCase(
        event.fleetId,
        event.startDate,
        event.endDate,
      );
      emit(state.copyWith(
        status: DiscoveryStatus.availabilitySuccess,
        availabilityData: result,
        errorMessage: null,
      ));
    } catch (e) {
      final msg = e is Failure ? e.message : e.toString();
      emit(state.copyWith(
        status: DiscoveryStatus.failure,
        errorMessage: msg,
      ));
    }
  }
}
