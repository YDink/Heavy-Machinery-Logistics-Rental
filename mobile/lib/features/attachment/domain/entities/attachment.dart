import 'attachment_photo.dart';
import 'attachment_type.dart';

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
