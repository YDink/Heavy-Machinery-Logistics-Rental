import '../../domain/entities/rental_extension.dart';

class RentalExtensionModel extends RentalExtension {
  const RentalExtensionModel({
    required super.bookingId,
    required super.extensionDays,
    required super.incrementalPrepaymentRequiredMinor,
    required super.additionalGrossMinor,
    required super.currency,
    required super.status,
    required super.ownerApprovalRequired,
  });

  factory RentalExtensionModel.fromJson(Map<String, dynamic> json) {
    return RentalExtensionModel(
      bookingId: json['bookingId'] as String? ?? '',
      extensionDays: (json['extensionDays'] as num?)?.toInt() ?? 0,
      incrementalPrepaymentRequiredMinor:
          (json['incrementalPrepaymentRequiredMinor'] as num?)?.toInt() ?? 0,
      additionalGrossMinor: (json['additionalGrossMinor'] as num?)?.toInt() ?? 0,
      currency: json['currency'] as String? ?? 'ETB',
      status: json['status'] as String? ?? 'REQUESTED',
      ownerApprovalRequired: json['ownerApprovalRequired'] as bool? ?? true,
    );
  }
}
