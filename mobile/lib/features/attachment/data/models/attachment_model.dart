import '../../domain/entities/attachment.dart';
import '../../domain/entities/attachment_photo.dart';
import '../../domain/entities/attachment_photo_angle.dart';
import '../../domain/entities/attachment_type.dart';

class AttachmentModel extends Attachment {
  const AttachmentModel({
    required super.id,
    required super.ownerId,
    required super.attachmentType,
    required super.make,
    required super.model,
    required super.subTypeSize,
    required super.serialNumber,
    required super.operatingWeightKg,
    required super.minCarrierWeightTonnes,
    required super.maxCarrierWeightTonnes,
    required super.requiredFlowMinLmin,
    required super.requiredFlowMaxLmin,
    required super.requiredPressureBar,
    required super.pinDiameterMm,
    required super.pinSpreadMm,
    required super.standaloneDailyRateMinor,
    required super.securityDepositMinor,
    required super.minRentalDays,
    required super.status,
    required super.compatibleSubcategories,
    required super.photos,
    required super.createdAt,
  });

  factory AttachmentModel.fromJson(Map<String, dynamic> json) {
    return AttachmentModel(
      id: json['attachmentId'] as String? ?? json['id'] as String? ?? '',
      ownerId: json['ownerId'] as String? ?? '',
      attachmentType: AttachmentType.fromCode(json['attachmentType'] as String? ?? ''),
      make: json['make'] as String? ?? '',
      model: json['model'] as String? ?? '',
      subTypeSize: json['subTypeSize'] as String? ?? '',
      serialNumber: json['serialNumber'] as String? ?? '',
      operatingWeightKg: (json['operatingWeightKg'] as num?)?.toDouble() ?? 0.0,
      minCarrierWeightTonnes: (json['minCarrierWeightTonnes'] as num?)?.toDouble() ?? 0.0,
      maxCarrierWeightTonnes: (json['maxCarrierWeightTonnes'] as num?)?.toDouble() ?? 0.0,
      requiredFlowMinLmin: (json['requiredFlowMinLmin'] as num?)?.toDouble() ?? 0.0,
      requiredFlowMaxLmin: (json['requiredFlowMaxLmin'] as num?)?.toDouble() ?? 0.0,
      requiredPressureBar: (json['requiredPressureBar'] as num?)?.toDouble() ?? 0.0,
      pinDiameterMm: (json['pinDiameterMm'] as num?)?.toDouble() ?? 0.0,
      pinSpreadMm: (json['pinSpreadMm'] as num?)?.toDouble() ?? 0.0,
      standaloneDailyRateMinor: (json['standaloneDailyRateMinor'] as num?)?.toInt() ?? 0,
      securityDepositMinor: (json['securityDepositMinor'] as num?)?.toInt() ?? 0,
      minRentalDays: (json['minRentalDays'] as num?)?.toInt() ?? 1,
      status: json['status'] as String? ?? 'ACTIVE',
      compatibleSubcategories: (json['compatibleSubcategories'] as List<dynamic>?)
              ?.map((e) => e.toString())
              .toList() ??
          [],
      photos: (json['photos'] as List<dynamic>?)
              ?.map((p) => AttachmentPhoto(
                    angle: AttachmentPhotoAngle.fromCode((p as Map<String, dynamic>)['angle'] as String? ?? ''),
                    storageUri: p['storageUri'] as String? ?? '',
                  ))
              .toList() ??
          [],
      createdAt: json['createdAt'] != null
          ? DateTime.tryParse(json['createdAt'] as String) ?? DateTime.now()
          : DateTime.now(),
    );
  }
}
