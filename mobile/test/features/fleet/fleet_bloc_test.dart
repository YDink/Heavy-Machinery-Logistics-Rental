import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/category.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/spec_field_group.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/repositories/fleet_repository.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_categories_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_spec_definitions_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/register_fleet_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_bloc.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_event.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_state.dart';

class MockFleetRepository implements FleetRepository {
  List<Category> categories = [
    const Category(
      id: 'cat-1',
      code: 'EXCAVATOR',
      name: 'Crawler Excavator',
      iconUrl: 'https://cdn.example.com/exc.png',
    ),
  ];

  List<SpecFieldGroup> specGroups = [
    const SpecFieldGroup(
      id: 'grp-1',
      categoryId: 'cat-1',
      sectionNameEN: 'Dimensions',
      sectionNameAM: 'ስፋት',
      displayOrder: 1,
      fields: [
        SpecFieldDefinition(
          id: 'def-1',
          fieldKey: 'operating_weight_kg',
          labelEN: 'Operating Weight',
          labelAM: 'የስራ ክብደት',
          fieldType: SpecFieldType.decimal,
          unitSymbol: 'kg',
          isRequired: true,
          isLogisticsCritical: true,
        ),
      ],
    ),
  ];

  FleetRegistrationRequest? savedDraft;

  @override
  Future<List<Category>> getCategories() async => categories;

  @override
  Future<List<SpecFieldGroup>> getSpecDefinitions(String categoryId) async =>
      specGroups;

  @override
  Future<FleetAsset> registerFleet(FleetRegistrationRequest request) async {
    return FleetAsset(
      id: 'fleet-123',
      ownerId: 'owner-456',
      categoryId: request.categoryId,
      make: request.make,
      model: request.model,
      manufacturingYear: request.manufacturingYear,
      registrationNumber: request.registrationNumber,
      dailyRateMinor: request.dailyRateMinor,
      status: 'PENDING_APPROVAL',
      latitude: request.latitude,
      longitude: request.longitude,
      addressText: request.addressText,
      version: 1,
    );
  }

  @override
  Future<void> saveDraft(FleetRegistrationRequest draft) async {
    savedDraft = draft;
  }

  @override
  Future<FleetRegistrationRequest?> getDraft(String categoryId) async =>
      savedDraft;
}

void main() {
  late MockFleetRepository mockRepo;
  late FleetRegistrationBloc bloc;

  setUp(() {
    mockRepo = MockFleetRepository();
    bloc = FleetRegistrationBloc(
      getCategoriesUseCase: GetCategoriesUseCase(mockRepo),
      getSpecDefinitionsUseCase: GetSpecDefinitionsUseCase(mockRepo),
      registerFleetUseCase: RegisterFleetUseCase(mockRepo),
      fleetRepository: mockRepo,
    );
  });

  tearDown(() {
    bloc.close();
  });

  test('initial state is FleetRegistrationInitial', () {
    expect(bloc.state, isA<FleetRegistrationInitial>());
  });

  test('LoadCategoriesEvent emits CategoriesLoading then CategoriesLoaded',
      () async {
    final expectedStates = [
      isA<CategoriesLoading>(),
      isA<CategoriesLoaded>(),
    ];

    expectLater(bloc.stream, emitsInOrder(expectedStates));

    bloc.add(const LoadCategoriesEvent());
  });

  test('SelectCategoryEvent loads dynamic spec definitions and draft',
      () async {
    bloc.add(const LoadCategoriesEvent());
    await bloc.stream.firstWhere((s) => s is CategoriesLoaded);

    bloc.add(const SelectCategoryEvent('cat-1'));

    await expectLater(
      bloc.stream,
      emits(predicate<CategoriesLoaded>((s) {
        return s.selectedCategoryId == 'cat-1' && s.specGroups.isNotEmpty;
      })),
    );
  });

  test('UpdatePhotoEvent records walkaround photo in state', () async {
    bloc.add(const LoadCategoriesEvent());
    await bloc.stream.firstWhere((s) => s is CategoriesLoaded);

    const photo = FleetPhoto(
      angle: PhotoAngle.front,
      storageUri: 's3://fleet/front.jpg',
      fileHash: 'hash1',
    );
    bloc.add(const UpdatePhotoEvent(photo));

    await expectLater(
      bloc.stream,
      emits(predicate<CategoriesLoaded>((s) {
        return s.photos[PhotoAngle.front] == photo;
      })),
    );
  });

  test('SubmitFleetRegistrationEvent registers fleet and emits success',
      () async {
    const request = FleetRegistrationRequest(
      categoryId: 'cat-1',
      make: 'Caterpillar',
      model: '320D',
      manufacturingYear: 2022,
      registrationNumber: 'ET-3-12345',
      chassisVin: 'VIN1234567890',
      baselineMeterHours: 500.0,
      fuelType: 'DIESEL',
      selfMoveCapable: false,
      dailyRateMinor: 4500000,
      latitude: 9.0300,
      longitude: 38.7400,
      addressText: 'Addis Ababa',
      photos: [
        FleetPhoto(
            angle: PhotoAngle.front, storageUri: 'uri', fileHash: 'hash'),
        FleetPhoto(angle: PhotoAngle.rear, storageUri: 'uri', fileHash: 'hash'),
        FleetPhoto(angle: PhotoAngle.left, storageUri: 'uri', fileHash: 'hash'),
        FleetPhoto(
            angle: PhotoAngle.right, storageUri: 'uri', fileHash: 'hash'),
        FleetPhoto(
            angle: PhotoAngle.cabin, storageUri: 'uri', fileHash: 'hash'),
        FleetPhoto(
            angle: PhotoAngle.meter, storageUri: 'uri', fileHash: 'hash'),
      ],
      specValues: [],
    );

    final expectedStates = [
      isA<FleetSubmitting>(),
      isA<FleetSubmitSuccess>(),
    ];

    expectLater(bloc.stream, emitsInOrder(expectedStates));

    bloc.add(const SubmitFleetRegistrationEvent(request));
  });
}
