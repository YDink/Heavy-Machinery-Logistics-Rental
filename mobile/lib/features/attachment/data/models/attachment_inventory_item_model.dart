import 'package:dink_mobile/features/attachment/domain/entities/attachment_inventory_item.dart';

class AttachmentInventoryItemModel extends AttachmentInventoryItem {
  const AttachmentInventoryItemModel({
    required super.id,
    required super.attachmentType,
    required super.make,
    required super.model,
    required super.subTypeSize,
    required super.status,
    required super.standaloneDailyRateMinor,
    required super.isLinkedToFleet,
    super.linkedFleetId,
    super.primaryPhotoUri,
  });

  factory AttachmentInventoryItemModel.fromJson(Map<String, dynamic> json) {
    return AttachmentInventoryItemModel(
      id: json['attachmentId'] as String? ?? json['id'] as String? ?? '',
      attachmentType: json['attachmentType'] as String? ?? '',
      make: json['make'] as String? ?? '',
      model: json['model'] as String? ?? '',
      subTypeSize: json['subTypeSize'] as String? ?? '',
      status: json['status'] as String? ?? 'ACTIVE',
      standaloneDailyRateMinor: (json['standaloneDailyRateMinor'] as num?)?.toInt() ?? 0,
      isLinkedToFleet: json['isLinkedToFleet'] as bool? ?? false,
      linkedFleetId: json['linkedFleetId'] as String?,
      primaryPhotoUri: json['primaryPhotoUri'] as String?,
    );
  }
}
