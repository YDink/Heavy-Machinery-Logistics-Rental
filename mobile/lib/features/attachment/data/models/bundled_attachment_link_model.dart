import '../../domain/entities/bundled_attachment_link.dart';

class BundledAttachmentLinkModel extends BundledAttachmentLink {
  const BundledAttachmentLinkModel({
    required super.fleetId,
    required super.attachmentId,
    required super.bundledDailyRateMinor,
    required super.discountPercentage,
    required super.status,
  });

  factory BundledAttachmentLinkModel.fromJson(Map<String, dynamic> json) {
    return BundledAttachmentLinkModel(
      fleetId: json['fleetId'] as String? ?? '',
      attachmentId: json['attachmentId'] as String? ?? '',
      bundledDailyRateMinor: (json['bundledDailyRateMinor'] as num?)?.toInt() ?? 0,
      discountPercentage: (json['discountPercentage'] as num?)?.toDouble() ?? 0.0,
      status: json['status'] as String? ?? 'LINKED',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'fleetId': fleetId,
      'attachmentId': attachmentId,
      'bundledDailyRateMinor': bundledDailyRateMinor,
      'discountPercentage': discountPercentage,
      'status': status,
    };
  }

  factory BundledAttachmentLinkModel.fromEntity(BundledAttachmentLink entity) {
    return BundledAttachmentLinkModel(
      fleetId: entity.fleetId,
      attachmentId: entity.attachmentId,
      bundledDailyRateMinor: entity.bundledDailyRateMinor,
      discountPercentage: entity.discountPercentage,
      status: entity.status,
    );
  }
}
