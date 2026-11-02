import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_availability.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_compare_matrix.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_result.dart';
import 'package:dink_mobile/features/discovery/domain/repositories/discovery_repository.dart';
import 'package:dink_mobile/features/discovery/domain/usecases/compare_fleets_usecase.dart';
import 'package:dink_mobile/features/discovery/domain/usecases/get_fleet_availability_usecase.dart';
import 'package:dink_mobile/features/discovery/domain/usecases/search_fleets_usecase.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_bloc.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_event.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_state.dart';

class MockDiscoveryRepository implements DiscoveryRepository {
  bool shouldFail = false;
  String errorMessage = 'Simulated search error';

  final testResultItem = const FleetSearchResult(
    fleetId: 'flt-001',
    subcategoryId: 'sub-001',
    subcategoryNameEn: 'Crawler Excavator',
    subcategoryNameAm: 'ክራውለር ኤክስካቫተር',
    make: 'Caterpillar',
    model: '320D',
    manufactureYear: 2024,
    dailyRateMinor: 2250000,
    excessHourRateMinor: 280000,
    distanceMeters: 12500.0,
    availabilityStatus: FleetAvailabilityStatus.availableNow,
    selfMoveCapable: false,
    operatorIncluded: true,
    fuelType: 'DIESEL',
    primaryPhotoUri: 'https://cdn.dinkmachinery.com/photo.webp',
    owner: FleetSearchOwnerSummary(
      ownerId: 'own-001',
      tradingName: 'Tesfaye Transport',
      averageRating: 4.95,
      totalCompletedRentals: 42,
      isVerified: true,
    ),
    keySpecs: FleetSearchKeySpecs(
      operatingWeightKg: 21800.0,
      payloadCapacityKg: null,
      grossPowerHp: 148.0,
    ),
  );

  @override
  Future<FleetSearchData> searchFleets(FleetSearchFilter filter) async {
    if (shouldFail) throw Exception(errorMessage);
    return FleetSearchData(
      totalMatches: 1,
      latitude: filter.latitude,
      longitude: filter.longitude,
      radiusKm: filter.radiusKm,
      results: [testResultItem],
    );
  }

  @override
  Future<FleetCompareData> compareFleets(List<String> fleetIds) async {
    if (shouldFail) throw Exception(errorMessage);
    return FleetCompareData(
      subcategoryId: 'sub-001',
      subcategoryNameEn: 'Crawler Excavator',
      subcategoryNameAm: 'ክራውለር ኤክስካቫተር',
      fleets: [
        const FleetCompareItemHeader(
          id: 'flt-001',
          title: 'CAT 320D',
          dailyRateMinor: 2250000,
          excessHourRateMinor: 280000,
          selfMoveCapable: false,
          operatorIncluded: true,
          fuelType: 'DIESEL',
        ),
        const FleetCompareItemHeader(
          id: 'flt-002',
          title: 'Komatsu PC200',
          dailyRateMinor: 2100000,
          excessHourRateMinor: 250000,
          selfMoveCapable: false,
          operatorIncluded: true,
          fuelType: 'DIESEL',
        ),
      ],
      comparisonMatrix: const [
        FleetCompareSpecRow(
          fieldKey: 'operating_weight_kg',
          labelEn: 'Operating Weight',
          labelAm: 'የስራ ክብደት',
          unitSymbol: 'kg',
          fieldType: 'DECIMAL',
          values: {'flt-001': '21800', 'flt-002': '20010'},
          isDifferent: true,
        ),
      ],
    );
  }

  @override
  Future<FleetAvailabilityData> getFleetAvailability(
    String fleetId,
    DateTime startDate,
    DateTime endDate,
  ) async {
    if (shouldFail) throw Exception(errorMessage);
    return FleetAvailabilityData(
      fleetId: fleetId,
      isAvailable: true,
      blackoutRanges: const [],
      reservedRanges: const [],
    );
  }
}

void main() {
  late MockDiscoveryRepository repository;
  late SearchFleetsUseCase searchFleetsUseCase;
  late CompareFleetsUseCase compareFleetsUseCase;
  late GetFleetAvailabilityUseCase getFleetAvailabilityUseCase;
  late DiscoveryBloc bloc;

  const defaultFilter = FleetSearchFilter(latitude: 9.0300, longitude: 38.7400);

  setUp(() {
    repository = MockDiscoveryRepository();
    searchFleetsUseCase = SearchFleetsUseCase(repository);
    compareFleetsUseCase = CompareFleetsUseCase(repository);
    getFleetAvailabilityUseCase = GetFleetAvailabilityUseCase(repository);

    bloc = DiscoveryBloc(
      searchFleetsUseCase: searchFleetsUseCase,
      compareFleetsUseCase: compareFleetsUseCase,
      getFleetAvailabilityUseCase: getFleetAvailabilityUseCase,
      initialState: const DiscoveryState(filter: defaultFilter),
    );
  });

  tearDown(() {
    bloc.close();
  });

  test('Initial state is correct', () {
    expect(bloc.state.status, DiscoveryStatus.initial);
    expect(bloc.state.filter, defaultFilter);
    expect(bloc.state.selectedForCompare, isEmpty);
    expect(bloc.state.isMapView, isFalse);
  });

  test('SearchFleetsRequested emits [loading, success] on valid search', () async {
    final expectedStates = [
      const DiscoveryState(
        status: DiscoveryStatus.loading,
        filter: defaultFilter,
      ),
      DiscoveryState(
        status: DiscoveryStatus.success,
        filter: defaultFilter,
        searchData: FleetSearchData(
          totalMatches: 1,
          latitude: 9.0300,
          longitude: 38.7400,
          radiusKm: 50.0,
          results: [repository.testResultItem],
        ),
      ),
    ];

    expectLater(bloc.stream, emitsInOrder(expectedStates));

    bloc.add(const SearchFleetsRequested(defaultFilter));
  });

  test('SearchFleetsRequested emits [loading, failure] on error', () async {
    repository.shouldFail = true;

    final expectedStates = [
      const DiscoveryState(
        status: DiscoveryStatus.loading,
        filter: defaultFilter,
      ),
      const DiscoveryState(
        status: DiscoveryStatus.failure,
        filter: defaultFilter,
        errorMessage: 'Exception: Simulated search error',
      ),
    ];

    expectLater(bloc.stream, emitsInOrder(expectedStates));

    bloc.add(const SearchFleetsRequested(defaultFilter));
  });

  test('ToggleViewModeRequested toggles map/list mode', () async {
    expect(bloc.state.isMapView, isFalse);
    bloc.add(const ToggleViewModeRequested(true));
    await Future.delayed(const Duration(milliseconds: 20));
    expect(bloc.state.isMapView, isTrue);
    bloc.add(const ToggleViewModeRequested(false));
    await Future.delayed(const Duration(milliseconds: 20));
    expect(bloc.state.isMapView, isFalse);
  });

  test('ToggleCompareSelection adds and removes fleets up to max 3 items', () async {
    final item1 = repository.testResultItem;
    const item2 = FleetSearchResult(
      fleetId: 'flt-002',
      subcategoryId: 'sub-001',
      subcategoryNameEn: 'Crawler Excavator',
      subcategoryNameAm: 'ክራውለር ኤክስካቫተር',
      make: 'Komatsu',
      model: 'PC200',
      manufactureYear: 2022,
      dailyRateMinor: 2100000,
      excessHourRateMinor: 250000,
      distanceMeters: 15000.0,
      availabilityStatus: FleetAvailabilityStatus.availableNow,
      selfMoveCapable: false,
      operatorIncluded: true,
      fuelType: 'DIESEL',
      owner: FleetSearchOwnerSummary(
        ownerId: 'own-002',
        tradingName: 'Almaz PLC',
        averageRating: 4.8,
        totalCompletedRentals: 30,
        isVerified: true,
      ),
      keySpecs: FleetSearchKeySpecs(),
    );
    const item3 = FleetSearchResult(
      fleetId: 'flt-003',
      subcategoryId: 'sub-001',
      subcategoryNameEn: 'Crawler Excavator',
      subcategoryNameAm: 'ክራውለር ኤክስካቫተር',
      make: 'Hitachi',
      model: 'ZX200',
      manufactureYear: 2023,
      dailyRateMinor: 2000000,
      excessHourRateMinor: 240000,
      distanceMeters: 18000.0,
      availabilityStatus: FleetAvailabilityStatus.availableNow,
      selfMoveCapable: false,
      operatorIncluded: true,
      fuelType: 'DIESEL',
      owner: FleetSearchOwnerSummary(
        ownerId: 'own-003',
        tradingName: 'Bekele Heavy',
        averageRating: 4.7,
        totalCompletedRentals: 25,
        isVerified: true,
      ),
      keySpecs: FleetSearchKeySpecs(),
    );
    const item4 = FleetSearchResult(
      fleetId: 'flt-004',
      subcategoryId: 'sub-001',
      subcategoryNameEn: 'Crawler Excavator',
      subcategoryNameAm: 'ክራውለር ኤክስካቫተር',
      make: 'Volvo',
      model: 'EC210',
      manufactureYear: 2021,
      dailyRateMinor: 2150000,
      excessHourRateMinor: 260000,
      distanceMeters: 22000.0,
      availabilityStatus: FleetAvailabilityStatus.availableNow,
      selfMoveCapable: false,
      operatorIncluded: true,
      fuelType: 'DIESEL',
      owner: FleetSearchOwnerSummary(
        ownerId: 'own-004',
        tradingName: 'Volvo Plant',
        averageRating: 4.9,
        totalCompletedRentals: 18,
        isVerified: true,
      ),
      keySpecs: FleetSearchKeySpecs(),
    );

    // Add 1st item
    bloc.add(ToggleCompareSelection(item1));
    await Future.delayed(const Duration(milliseconds: 20));
    expect(bloc.state.selectedForCompare, [item1]);

    // Add 2nd item
    bloc.add(ToggleCompareSelection(item2));
    await Future.delayed(const Duration(milliseconds: 20));
    expect(bloc.state.selectedForCompare, [item1, item2]);

    // Add 3rd item
    bloc.add(ToggleCompareSelection(item3));
    await Future.delayed(const Duration(milliseconds: 20));
    expect(bloc.state.selectedForCompare, [item1, item2, item3]);

    // Add 4th item (Should be capped at 3 per SR-08)
    bloc.add(ToggleCompareSelection(item4));
    await Future.delayed(const Duration(milliseconds: 20));
    expect(bloc.state.selectedForCompare.length, 3);
    expect(bloc.state.selectedForCompare.contains(item4), isFalse);

    // Remove 2nd item
    bloc.add(ToggleCompareSelection(item2));
    await Future.delayed(const Duration(milliseconds: 20));
    expect(bloc.state.selectedForCompare, [item1, item3]);

    // Clear all
    bloc.add(const ClearCompareSelection());
    await Future.delayed(const Duration(milliseconds: 20));
    expect(bloc.state.selectedForCompare, isEmpty);
  });

  test('ExecuteComparisonRequested executes when at least 2 fleets selected', () async {
    final item1 = repository.testResultItem;
    const item2 = FleetSearchResult(
      fleetId: 'flt-002',
      subcategoryId: 'sub-001',
      subcategoryNameEn: 'Crawler Excavator',
      subcategoryNameAm: 'ክራውለር ኤክስካቫተር',
      make: 'Komatsu',
      model: 'PC200',
      manufactureYear: 2022,
      dailyRateMinor: 2100000,
      excessHourRateMinor: 250000,
      distanceMeters: 15000.0,
      availabilityStatus: FleetAvailabilityStatus.availableNow,
      selfMoveCapable: false,
      operatorIncluded: true,
      fuelType: 'DIESEL',
      owner: FleetSearchOwnerSummary(
        ownerId: 'own-002',
        tradingName: 'Almaz PLC',
        averageRating: 4.8,
        totalCompletedRentals: 30,
        isVerified: true,
      ),
      keySpecs: FleetSearchKeySpecs(),
    );

    bloc.add(ToggleCompareSelection(item1));
    bloc.add(ToggleCompareSelection(item2));

    bloc.add(const ExecuteComparisonRequested());

    await expectLater(
      bloc.stream,
      emitsThrough(predicate<DiscoveryState>((s) => s.status == DiscoveryStatus.compareSuccess)),
    );

    expect(bloc.state.compareData, isNotNull);
    expect(bloc.state.compareData!.fleets.length, 2);
  });

  test('CheckAvailabilityRequested emits availabilitySuccess', () async {
    bloc.add(CheckAvailabilityRequested(
      fleetId: 'flt-001',
      startDate: DateTime(2026, 11, 1),
      endDate: DateTime(2026, 11, 10),
    ));

    await expectLater(
      bloc.stream,
      emitsThrough(predicate<DiscoveryState>((s) => s.status == DiscoveryStatus.availabilitySuccess)),
    );

    expect(bloc.state.availabilityData, isNotNull);
    expect(bloc.state.availabilityData!.isAvailable, isTrue);
  });
}
