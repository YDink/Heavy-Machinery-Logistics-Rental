import '../../domain/entities/attachment.dart';
import '../../domain/entities/attachment_inventory_item.dart';
import '../../domain/entities/attachment_photo.dart';
import '../../domain/entities/attachment_registration_request.dart';
import '../../domain/entities/attachment_type.dart';
import '../../domain/entities/bundled_attachment_link.dart';
import '../../domain/repositories/attachment_repository.dart';
import '../datasources/attachment_local_data_source.dart';
import '../datasources/attachment_remote_data_source.dart';
import '../models/attachment_registration_request_model.dart';

class AttachmentRepositoryImpl implements AttachmentRepository {
  final AttachmentRemoteDataSource remoteDataSource;
  final AttachmentLocalDataSource localDataSource;

  AttachmentRepositoryImpl({
    required this.remoteDataSource,
    required this.localDataSource,
  });

  @override
  Future<Attachment> registerAttachment(AttachmentRegistrationRequest request) async {
    final model = AttachmentRegistrationRequestModel.fromEntity(request);
    final result = await remoteDataSource.registerAttachment(model);
    await localDataSource.clearDraft();
    return result;
  }

  @override
  Future<List<AttachmentInventoryItem>> getInventory({String? status}) {
    return remoteDataSource.getInventory(status: status);
  }

  @override
  Future<BundledAttachmentLink> linkAttachmentToFleet({
    required String fleetId,
    required String attachmentId,
    required int bundledDailyRateMinor,
  }) {
    return remoteDataSource.linkAttachmentToFleet(
      fleetId: fleetId,
      attachmentId: attachmentId,
      bundledDailyRateMinor: bundledDailyRateMinor,
    );
  }

  @override
  Future<void> saveDraft(AttachmentRegistrationRequest request) {
    return localDataSource.saveDraft(
      AttachmentRegistrationRequestModel.fromEntity(request),
    );
  }

  @override
  Future<AttachmentRegistrationRequest?> getDraft() async {
    final model = await localDataSource.getDraft();
    if (model == null) return null;
    return AttachmentRegistrationRequest(
      attachmentType: AttachmentType.fromCode(model.attachmentTypeCode),
      make: model.make,
      model: model.model,
      subTypeSize: model.subTypeSize,
      serialNumber: model.serialNumber,
      operatingWeightKg: model.operatingWeightKg,
      minCarrierWeightTonnes: model.minCarrierWeightTonnes,
      maxCarrierWeightTonnes: model.maxCarrierWeightTonnes,
      requiredFlowMinLmin: model.requiredFlowMinLmin,
      requiredFlowMaxLmin: model.requiredFlowMaxLmin,
      requiredPressureBar: model.requiredPressureBar,
      pinDiameterMm: model.pinDiameterMm,
      pinSpreadMm: model.pinSpreadMm,
      compatibleSubcategories: model.compatibleSubcategories,
      standaloneDailyRateMinor: model.standaloneDailyRateMinor,
      securityDepositMinor: model.securityDepositMinor,
      minRentalDays: model.minRentalDays,
      photos: model.photos
          .map((p) => AttachmentPhoto(
                angle: p.angle,
                storageUri: p.storageUri,
                localFilePath: p.localFilePath,
              ))
          .toList(),
    );
  }

  @override
  Future<void> clearDraft() {
    return localDataSource.clearDraft();
  }
}
