import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment_inventory_item.dart';
import 'package:dink_mobile/features/attachment/domain/entities/bundled_attachment_link.dart';
import 'package:dink_mobile/features/attachment/domain/repositories/attachment_repository.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/get_attachment_inventory_usecase.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/link_attachment_usecase.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/register_attachment_usecase.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_bloc.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_event.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_state.dart';

class MockAttachmentRepository implements AttachmentRepository {
  AttachmentRegistrationRequest? savedDraft;
  bool shouldFail = false;
  String errorMessage = 'Simulated repository error';

  final List<AttachmentInventoryItem> items = [
    const AttachmentInventoryItem(
      id: 'att-001',
      attachmentType: 'HYDRAULIC_BREAKER',
      make: 'Soosan',
      model: 'SB70',
      subTypeSize: '135mm Chisel',
      standaloneDailyRateMinor: 1500000,
      status: 'ACTIVE',
      isLinkedToFleet: false,
    ),
  ];

  @override
  Future<Attachment> registerAttachment(AttachmentRegistrationRequest request) async {
    if (shouldFail) throw Exception(errorMessage);
    return Attachment(
      id: 'att-created-123',
      ownerId: 'owner-456',
      attachmentType: request.attachmentType,
      make: request.make,
      model: request.model,
      subTypeSize: request.subTypeSize,
      serialNumber: request.serialNumber,
      operatingWeightKg: request.operatingWeightKg,
      minCarrierWeightTonnes: request.minCarrierWeightTonnes,
      maxCarrierWeightTonnes: request.maxCarrierWeightTonnes,
      requiredPressureBar: request.requiredPressureBar,
      requiredFlowMinLmin: request.requiredFlowMinLmin,
      requiredFlowMaxLmin: request.requiredFlowMaxLmin,
      pinDiameterMm: request.pinDiameterMm,
      pinSpreadMm: request.pinSpreadMm,
      compatibleSubcategories: request.compatibleSubcategories,
      standaloneDailyRateMinor: request.standaloneDailyRateMinor,
      securityDepositMinor: request.securityDepositMinor,
      minRentalDays: request.minRentalDays,
      status: 'AVAILABLE',
      photos: request.photos,
      createdAt: DateTime.utc(2026, 3, 31),
    );
  }

  @override
  Future<List<AttachmentInventoryItem>> getInventory({String? status}) async {
    if (shouldFail) throw Exception(errorMessage);
    return items;
  }

  @override
  Future<BundledAttachmentLink> linkAttachmentToFleet({
    required String fleetId,
    required String attachmentId,
    required int bundledDailyRateMinor,
  }) async {
    if (shouldFail) throw Exception(errorMessage);
    return BundledAttachmentLink(
      fleetId: fleetId,
      attachmentId: attachmentId,
      bundledDailyRateMinor: bundledDailyRateMinor,
      discountPercentage: 20.0,
      status: 'ACTIVE',
    );
  }

  @override
  Future<void> saveDraft(AttachmentRegistrationRequest request) async {
    savedDraft = request;
  }

  @override
  Future<AttachmentRegistrationRequest?> getDraft() async {
    return savedDraft;
  }

  @override
  Future<void> clearDraft() async {
    savedDraft = null;
  }
}

void main() {
  late MockAttachmentRepository mockRepo;
  late AttachmentBloc bloc;

  setUp(() {
    mockRepo = MockAttachmentRepository();
    bloc = AttachmentBloc(
      registerAttachmentUseCase: RegisterAttachmentUseCase(mockRepo),
      getInventoryUseCase: GetAttachmentInventoryUseCase(mockRepo),
      linkAttachmentUseCase: LinkAttachmentUseCase(mockRepo),
      attachmentRepository: mockRepo,
    );
  });

  tearDown(() {
    bloc.close();
  });

  group('AttachmentBloc', () {
    test('initial state has initial status and empty maps', () {
      expect(bloc.state.status, AttachmentStatusState.initial);
      expect(bloc.state.capturedPhotos, isEmpty);
      expect(bloc.state.areAllPhotosCaptured, isFalse);
      expect(bloc.state.inventory, isEmpty);
    });

    test('PhotoCapturedEvent adds photos and updates areAllPhotosCaptured', () async {
      final states = <AttachmentState>[];
      bloc.stream.listen(states.add);

      bloc.add(const PhotoCapturedEvent(
        angle: AttachmentPhotoAngle.overall,
        storageUri: 's3://media/overall.webp',
        localFilePath: '/tmp/overall.jpg',
      ));
      bloc.add(const PhotoCapturedEvent(
        angle: AttachmentPhotoAngle.bracket,
        storageUri: 's3://media/bracket.webp',
        localFilePath: '/tmp/bracket.jpg',
      ));
      bloc.add(const PhotoCapturedEvent(
        angle: AttachmentPhotoAngle.pins,
        storageUri: 's3://media/pins.webp',
        localFilePath: '/tmp/pins.jpg',
      ));
      bloc.add(const PhotoCapturedEvent(
        angle: AttachmentPhotoAngle.toolTip,
        storageUri: 's3://media/tooltip.webp',
        localFilePath: '/tmp/tooltip.jpg',
      ));

      await Future.delayed(const Duration(milliseconds: 50));

      expect(bloc.state.capturedPhotos.length, 4);
      expect(bloc.state.areAllPhotosCaptured, isTrue);
      expect(bloc.state.capturedPhotos[AttachmentPhotoAngle.overall]?.storageUri, 's3://media/overall.webp');
      expect(bloc.state.capturedPhotos[AttachmentPhotoAngle.toolTip]?.storageUri, 's3://media/tooltip.webp');
    });

    test('SaveAttachmentDraftEvent and LoadAttachmentDraftEvent preserve state and photos', () async {
      const sampleDraft = AttachmentRegistrationRequest(
        attachmentType: AttachmentType.hydraulicBreaker,
        make: 'Soosan',
        model: 'SB70',
        subTypeSize: '135mm',
        serialNumber: 'SN-998877',
        operatingWeightKg: 1800,
        minCarrierWeightTonnes: 18,
        maxCarrierWeightTonnes: 26,
        requiredPressureBar: 160,
        requiredFlowMinLmin: 130,
        requiredFlowMaxLmin: 150,
        pinDiameterMm: 80,
        pinSpreadMm: 450,
        compatibleSubcategories: ['cat_sub_excavator_standard'],
        standaloneDailyRateMinor: 1500000,
        securityDepositMinor: 20000000,
        minRentalDays: 1,
        photos: [
          AttachmentPhoto(angle: AttachmentPhotoAngle.overall, storageUri: 's3://draft/overall.webp'),
          AttachmentPhoto(angle: AttachmentPhotoAngle.bracket, storageUri: 's3://draft/bracket.webp'),
        ],
      );

      bloc.add(const SaveAttachmentDraftEvent(sampleDraft));
      await Future.delayed(const Duration(milliseconds: 50));

      expect(bloc.state.status, AttachmentStatusState.draftSaved);
      expect(mockRepo.savedDraft, equals(sampleDraft));

      bloc.add(const LoadAttachmentDraftEvent());
      await Future.delayed(const Duration(milliseconds: 50));

      expect(bloc.state.status, AttachmentStatusState.draftLoaded);
      expect(bloc.state.draft?.model, 'SB70');
      expect(bloc.state.capturedPhotos.length, 2);
      expect(bloc.state.capturedPhotos[AttachmentPhotoAngle.overall]?.storageUri, 's3://draft/overall.webp');
    });

    test('SubmitAttachmentRegistrationEvent emits loading then registered on success', () async {
      const request = AttachmentRegistrationRequest(
        attachmentType: AttachmentType.hydraulicBreaker,
        make: 'Soosan',
        model: 'SB70',
        subTypeSize: '135mm',
        serialNumber: 'SN-12345',
        operatingWeightKg: 1800,
        minCarrierWeightTonnes: 18,
        maxCarrierWeightTonnes: 26,
        requiredPressureBar: 160,
        requiredFlowMinLmin: 130,
        requiredFlowMaxLmin: 150,
        pinDiameterMm: 80,
        pinSpreadMm: 450,
        compatibleSubcategories: ['cat_sub_excavator_standard'],
        standaloneDailyRateMinor: 1500000,
        securityDepositMinor: 20000000,
        minRentalDays: 1,
        photos: [
          AttachmentPhoto(angle: AttachmentPhotoAngle.overall, storageUri: 'uri1'),
          AttachmentPhoto(angle: AttachmentPhotoAngle.bracket, storageUri: 'uri2'),
          AttachmentPhoto(angle: AttachmentPhotoAngle.pins, storageUri: 'uri3'),
          AttachmentPhoto(angle: AttachmentPhotoAngle.toolTip, storageUri: 'uri4'),
        ],
      );

      final expectedStatuses = <AttachmentStatusState>[];
      bloc.stream.listen((state) => expectedStatuses.add(state.status));

      bloc.add(const SubmitAttachmentRegistrationEvent(request));
      await Future.delayed(const Duration(milliseconds: 50));

      expect(expectedStatuses, containsAllInOrder([
        AttachmentStatusState.loading,
        AttachmentStatusState.registered,
      ]));
      expect(bloc.state.capturedPhotos, isEmpty);
    });

    test('SubmitAttachmentRegistrationEvent emits loading then error on failure', () async {
      mockRepo.shouldFail = true;
      mockRepo.errorMessage = 'Network failure while registering';

      const request = AttachmentRegistrationRequest(
        attachmentType: AttachmentType.hydraulicBreaker,
        make: 'Soosan',
        model: 'SB70',
        subTypeSize: '135mm',
        serialNumber: 'SN-12345',
        operatingWeightKg: 1800,
        minCarrierWeightTonnes: 18,
        maxCarrierWeightTonnes: 26,
        requiredPressureBar: 160,
        requiredFlowMinLmin: 130,
        requiredFlowMaxLmin: 150,
        pinDiameterMm: 80,
        pinSpreadMm: 450,
        compatibleSubcategories: ['cat_sub_excavator_standard'],
        standaloneDailyRateMinor: 1500000,
        securityDepositMinor: 20000000,
        minRentalDays: 1,
        photos: [],
      );

      final expectedStatuses = <AttachmentStatusState>[];
      bloc.stream.listen((state) => expectedStatuses.add(state.status));

      bloc.add(const SubmitAttachmentRegistrationEvent(request));
      await Future.delayed(const Duration(milliseconds: 50));

      expect(expectedStatuses, containsAllInOrder([
        AttachmentStatusState.loading,
        AttachmentStatusState.error,
      ]));
      expect(bloc.state.errorMessage, contains('Network failure while registering'));
    });

    test('LoadAttachmentInventoryEvent loads items successfully', () async {
      final expectedStatuses = <AttachmentStatusState>[];
      bloc.stream.listen((state) => expectedStatuses.add(state.status));

      bloc.add(const LoadAttachmentInventoryEvent());
      await Future.delayed(const Duration(milliseconds: 50));

      expect(expectedStatuses, containsAllInOrder([
        AttachmentStatusState.loading,
        AttachmentStatusState.loaded,
      ]));
      expect(bloc.state.inventory.length, 1);
      expect(bloc.state.inventory.first.model, 'SB70');
    });

    test('LinkAttachmentToFleetEvent links attachment and updates state', () async {
      final expectedStatuses = <AttachmentStatusState>[];
      bloc.stream.listen((state) => expectedStatuses.add(state.status));

      bloc.add(const LinkAttachmentToFleetEvent(
        fleetId: 'fleet-100',
        attachmentId: 'att-001',
        bundledDailyRateMinor: 1200000,
      ));
      await Future.delayed(const Duration(milliseconds: 50));

      expect(expectedStatuses, containsAllInOrder([
        AttachmentStatusState.loading,
        AttachmentStatusState.linked,
      ]));
      expect(bloc.state.lastLinked?.fleetId, 'fleet-100');
      expect(bloc.state.lastLinked?.bundledDailyRateMinor, 1200000);
    });
  });
}
