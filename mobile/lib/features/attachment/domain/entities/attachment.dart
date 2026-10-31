enum AttachmentType {
  hydraulicBreaker('HYDRAULIC_BREAKER', 'Hydraulic Breaker'),
  standardBucket('STANDARD_BUCKET', 'Standard Bucket'),
  trenchingBucket('TRENCHING_BUCKET', 'Trenching Bucket'),
  auger('AUGER', 'Earth Auger'),
  ripperTooth('RIPPER_TOOTH', 'Ripper Tooth'),
  sortingGrapple('SORTING_GRAPPLE', 'Sorting Grapple'),
  crusher('CRUSHER', 'Demolition Crusher');

  final String code;
  final String displayName;
  const AttachmentType(this.code, this.displayName);

  static AttachmentType fromCode(String code) {
    return AttachmentType.values.firstWhere(
      (t) => t.code == code,
      orElse: () => AttachmentType.standardBucket,
    );
  }
}

enum AttachmentPhotoAngle {
  overall('OVERALL', 'Overall Profile'),
  bracket('BRACKET', 'Mounting Bracket'),
  pins('PINS', 'Pins & Bushings'),
  toolTip('TOOL_TIP', 'Working Tip / Chisel');

  final String code;
  final String label;
  const AttachmentPhotoAngle(this.code, this.label);

  static AttachmentPhotoAngle fromCode(String code) {
    return AttachmentPhotoAngle.values.firstWhere(
      (a) => a.code == code,
      orElse: () => AttachmentPhotoAngle.overall,
    );
  }
}

class AttachmentPhoto {
  final AttachmentPhotoAngle angle;
  final String storageUri;
  final String? localFilePath;

  const AttachmentPhoto({
    required this.angle,
    required this.storageUri,
    this.localFilePath,
  });
}

class Attachment {
  final String id;
  final String ownerId;
  final AttachmentType attachmentType;
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
  final int standaloneDailyRateMinor;
  final int securityDepositMinor;
  final int minRentalDays;
  final String status;
  final List<String> compatibleSubcategories;
  final List<AttachmentPhoto> photos;
  final DateTime createdAt;

  const Attachment({
    required this.id,
    required this.ownerId,
    required this.attachmentType,
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
    required this.standaloneDailyRateMinor,
    required this.securityDepositMinor,
    required this.minRentalDays,
    required this.status,
    required this.compatibleSubcategories,
    required this.photos,
    required this.createdAt,
  });
}

class AttachmentRegistrationRequest {
  final AttachmentType attachmentType;
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
  final List<AttachmentPhoto> photos;

  const AttachmentRegistrationRequest({
    required this.attachmentType,
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
}
