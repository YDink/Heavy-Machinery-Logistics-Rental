import '../../domain/entities/booking.dart';

class BookingAttachmentModel extends BookingAttachmentEntity {
  const BookingAttachmentModel({
    required super.id,
    required super.attachmentId,
    super.attachmentType,
    super.make,
    super.model,
    required super.bundledDailyRateMinor,
    super.currency = 'ETB',
  });

  factory BookingAttachmentModel.fromJson(Map<String, dynamic> json) {
    final rateMap = json['bundledDailyRate'] as Map<String, dynamic>? ?? {};
    return BookingAttachmentModel(
      id: json['id'] as String? ?? '',
      attachmentId: json['attachmentId'] as String? ?? '',
      attachmentType: json['attachmentType'] as String?,
      make: json['make'] as String?,
      model: json['model'] as String?,
      bundledDailyRateMinor: (rateMap['amountMinor'] as num?)?.toInt() ?? 0,
      currency: rateMap['currency'] as String? ?? 'ETB',
    );
  }
}

class BookingModel extends BookingEntity {
  const BookingModel({
    required super.id,
    required super.bookingReference,
    required super.fleetId,
    required super.renterId,
    required super.ownerId,
    required super.status,
    required super.startDate,
    required super.endDate,
    required super.durationDays,
    required super.dailyRateMinor,
    required super.baseHireFeeMinor,
    super.attachmentsFeeMinor = 0,
    required super.grossRentalAmountMinor,
    required super.prepaymentAmountMinor,
    super.penaltyAmountMinor = 0,
    super.refundAmountMinor = 0,
    super.currency = 'ETB',
    required super.siteLocationAddress,
    super.rejectionReasonCode,
    super.rejectionNotes,
    super.cancellationReason,
    super.handoverPin,
    super.checkInMeterHours,
    super.checkOutMeterHours,
    super.expiresAt,
    super.confirmedAt,
    super.activatedAt,
    super.completedAt,
    super.cancelledAt,
    super.attachments = const [],
    required super.version,
    required super.createdAt,
  });

  factory BookingModel.fromJson(Map<String, dynamic> json) {
    final dailyRateMap = json['dailyRate'] as Map<String, dynamic>? ?? {};
    final baseHireMap = json['baseHireFee'] as Map<String, dynamic>? ?? {};
    final attsFeeMap = json['attachmentsFee'] as Map<String, dynamic>? ?? {};
    final grossMap = json['grossRentalAmount'] as Map<String, dynamic>? ?? {};
    final prepayMap = json['prepaymentAmount'] as Map<String, dynamic>? ?? {};
    final penaltyMap = json['penaltyAmount'] as Map<String, dynamic>?;
    final refundMap = json['refundAmount'] as Map<String, dynamic>?;

    final attList = (json['attachments'] as List<dynamic>?)
            ?.map((e) => BookingAttachmentModel.fromJson(e as Map<String, dynamic>))
            .toList() ??
        [];

    DateTime? parseDate(dynamic v) {
      if (v == null) return null;
      if (v is String) return DateTime.tryParse(v);
      return null;
    }

    return BookingModel(
      id: json['id'] as String? ?? '',
      bookingReference: json['bookingReference'] as String? ?? '',
      fleetId: json['fleetId'] as String? ?? '',
      renterId: json['renterId'] as String? ?? '',
      ownerId: json['ownerId'] as String? ?? '',
      status: BookingStatus.fromString(json['status'] as String? ?? 'PENDING_PAYMENT'),
      startDate: DateTime.tryParse(json['startDate'] as String? ?? '') ?? DateTime.now(),
      endDate: DateTime.tryParse(json['endDate'] as String? ?? '') ?? DateTime.now(),
      durationDays: (json['durationDays'] as num?)?.toInt() ?? 1,
      dailyRateMinor: (dailyRateMap['amountMinor'] as num?)?.toInt() ?? 0,
      baseHireFeeMinor: (baseHireMap['amountMinor'] as num?)?.toInt() ?? 0,
      attachmentsFeeMinor: (attsFeeMap['amountMinor'] as num?)?.toInt() ?? 0,
      grossRentalAmountMinor: (grossMap['amountMinor'] as num?)?.toInt() ?? 0,
      prepaymentAmountMinor: (prepayMap['amountMinor'] as num?)?.toInt() ?? 0,
      penaltyAmountMinor: (penaltyMap?['amountMinor'] as num?)?.toInt() ?? 0,
      refundAmountMinor: (refundMap?['amountMinor'] as num?)?.toInt() ?? 0,
      currency: dailyRateMap['currency'] as String? ?? 'ETB',
      siteLocationAddress: json['siteLocationAddress'] as String? ?? '',
      rejectionReasonCode: json['rejectionReasonCode'] as String?,
      rejectionNotes: json['rejectionNotes'] as String?,
      cancellationReason: json['cancellationReason'] as String?,
      handoverPin: json['handoverPin'] as String?,
      checkInMeterHours: (json['checkInMeterHours'] as num?)?.toDouble(),
      checkOutMeterHours: (json['checkOutMeterHours'] as num?)?.toDouble(),
      expiresAt: parseDate(json['expiresAt']),
      confirmedAt: parseDate(json['confirmedAt']),
      activatedAt: parseDate(json['activatedAt']),
      completedAt: parseDate(json['completedAt']),
      cancelledAt: parseDate(json['cancelledAt']),
      attachments: attList,
      version: (json['version'] as num?)?.toInt() ?? 1,
      createdAt: DateTime.tryParse(json['createdAt'] as String? ?? '') ?? DateTime.now(),
    );
  }
}

class BankSlipModel extends BankSlipEntity {
  const BankSlipModel({
    required super.id,
    required super.bookingId,
    required super.uploaderId,
    required super.storageUri,
    required super.fileHash,
    super.transactionReference,
    required super.amountMinor,
    super.currency = 'ETB',
    required super.status,
    required super.createdAt,
  });

  factory BankSlipModel.fromJson(Map<String, dynamic> json) {
    final amountMap = json['amount'] as Map<String, dynamic>? ?? {};
    return BankSlipModel(
      id: json['id'] as String? ?? '',
      bookingId: json['bookingId'] as String? ?? '',
      uploaderId: json['uploaderId'] as String? ?? '',
      storageUri: json['storageUri'] as String? ?? '',
      fileHash: json['fileHash'] as String? ?? '',
      transactionReference: json['transactionReference'] as String?,
      amountMinor: (amountMap['amountMinor'] as num?)?.toInt() ?? 0,
      currency: amountMap['currency'] as String? ?? 'ETB',
      status: json['status'] as String? ?? 'PENDING_AUDIT',
      createdAt: DateTime.tryParse(json['createdAt'] as String? ?? '') ?? DateTime.now(),
    );
  }
}
