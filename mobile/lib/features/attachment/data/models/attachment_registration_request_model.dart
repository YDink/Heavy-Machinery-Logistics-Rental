import '../../domain/entities/attachment_registration_request.dart';
import 'attachment_photo_model.dart';

class AttachmentRegistrationRequestModel {
  final String attachmentTypeCode;
  final String make;
  final String model;
  final String subTypeSize;
  final String serialNumber;
  final double operatingWeightKg;
  final double minCarrierWeightTonnes;
  final double maxCarrierWeightTonnes;
  final double requiredFlowMinLmin;
  final double requiredFlowMaxLmin;
  final double requiredPressureBar;
  final double pinDiameterMm;
  final double pinSpreadMm;
  final List<String> compatibleSubcategories;
  final int standaloneDailyRateMinor;
  final int securityDepositMinor;
  final int minRentalDays;
  final List<AttachmentPhotoModel> photos;

  const AttachmentRegistrationRequestModel({
    required this.attachmentTypeCode,
    required this.make,
    required this.model,
    required this.subTypeSize,
    required this.serialNumber,
    required this.operatingWeightKg,
    required this.minCarrierWeightTonnes,
    required this.maxCarrierWeightTonnes,
    required this.requiredFlowMinLmin,
    required this.requiredFlowMaxLmin,
    required this.requiredPressureBar,
    required this.pinDiameterMm,
    required this.pinSpreadMm,
    required this.compatibleSubcategories,
    required this.standaloneDailyRateMinor,
    required this.securityDepositMinor,
    required this.minRentalDays,
    required this.photos,
  });

  factory AttachmentRegistrationRequestModel.fromJson(Map<String, dynamic> json) {
    final photosRaw = json['photos'] as List<dynamic>? ?? [];
    return AttachmentRegistrationRequestModel(
      attachmentTypeCode: json['attachmentType'] as String? ?? '',
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
      compatibleSubcategories: (json['compatibleSubcategories'] as List<dynamic>?)
              ?.map((e) => e.toString())
              .toList() ??
          [],
      standaloneDailyRateMinor: (json['standaloneDailyRateMinor'] as num?)?.toInt() ?? 0,
      securityDepositMinor: (json['securityDepositMinor'] as num?)?.toInt() ?? 0,
      minRentalDays: (json['minRentalDays'] as num?)?.toInt() ?? 1,
      photos: photosRaw
          .map((p) => AttachmentPhotoModel.fromJson(p as Map<String, dynamic>))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'attachmentType': attachmentTypeCode,
      'make': make,
      'model': model,
      'subTypeSize': subTypeSize,
      'serialNumber': serialNumber,
      'operatingWeightKg': operatingWeightKg,
      'minCarrierWeightTonnes': minCarrierWeightTonnes,
      'maxCarrierWeightTonnes': maxCarrierWeightTonnes,
      'requiredFlowMinLmin': requiredFlowMinLmin,
      'requiredFlowMaxLmin': requiredFlowMaxLmin,
      'requiredPressureBar': requiredPressureBar,
      'pinDiameterMm': pinDiameterMm,
      'pinSpreadMm': pinSpreadMm,
      'compatibleSubcategories': compatibleSubcategories,
      'standaloneDailyRateMinor': standaloneDailyRateMinor,
      'securityDepositMinor': securityDepositMinor,
      'minRentalDays': minRentalDays,
      'photos': photos.map((p) => p.toJson()).toList(),
    };
  }

  factory AttachmentRegistrationRequestModel.fromEntity(AttachmentRegistrationRequest entity) {
    return AttachmentRegistrationRequestModel(
      attachmentTypeCode: entity.attachmentType.code,
      make: entity.make,
      model: entity.model,
      subTypeSize: entity.subTypeSize,
      serialNumber: entity.serialNumber,
      operatingWeightKg: entity.operatingWeightKg,
      minCarrierWeightTonnes: entity.minCarrierWeightTonnes,
      maxCarrierWeightTonnes: entity.maxCarrierWeightTonnes,
      requiredFlowMinLmin: entity.requiredFlowMinLmin,
      requiredFlowMaxLmin: entity.requiredFlowMaxLmin,
      requiredPressureBar: entity.requiredPressureBar,
      pinDiameterMm: entity.pinDiameterMm,
      pinSpreadMm: entity.pinSpreadMm,
      compatibleSubcategories: entity.compatibleSubcategories,
      standaloneDailyRateMinor: entity.standaloneDailyRateMinor,
      securityDepositMinor: entity.securityDepositMinor,
      minRentalDays: entity.minRentalDays,
      photos: entity.photos.map(AttachmentPhotoModel.fromEntity).toList(),
    );
  }
}
