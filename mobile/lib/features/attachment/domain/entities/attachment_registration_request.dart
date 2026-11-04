import 'attachment_photo.dart';
import 'attachment_type.dart';

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
