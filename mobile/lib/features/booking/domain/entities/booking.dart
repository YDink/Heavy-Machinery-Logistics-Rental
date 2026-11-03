import 'package:equatable/equatable.dart';

enum BookingStatus {
  draft,
  pendingPayment,
  paymentEvidenceUploaded,
  paymentDeclined,
  disputeOpened,
  disputeResolvedRenter,
  disputeResolvedOwner,
  confirmed,
  awaitingFleetReadiness,
  pendingAdditionalPayment,
  readyForPickup,
  active,
  returnInitiated,
  pendingFinalPayment,
  completed,
  rated,
  cancelled;

  static BookingStatus fromString(String val) {
    switch (val.toUpperCase()) {
      case 'DRAFT':
        return BookingStatus.draft;
      case 'PENDING_PAYMENT':
        return BookingStatus.pendingPayment;
      case 'PAYMENT_EVIDENCE_UPLOADED':
        return BookingStatus.paymentEvidenceUploaded;
      case 'PAYMENT_DECLINED':
        return BookingStatus.paymentDeclined;
      case 'DISPUTE_OPENED':
        return BookingStatus.disputeOpened;
      case 'DISPUTE_RESOLVED_RENTER':
        return BookingStatus.disputeResolvedRenter;
      case 'DISPUTE_RESOLVED_OWNER':
        return BookingStatus.disputeResolvedOwner;
      case 'CONFIRMED':
        return BookingStatus.confirmed;
      case 'AWAITING_FLEET_READINESS':
        return BookingStatus.awaitingFleetReadiness;
      case 'PENDING_ADDITIONAL_PAYMENT':
        return BookingStatus.pendingAdditionalPayment;
      case 'READY_FOR_PICKUP':
        return BookingStatus.readyForPickup;
      case 'ACTIVE':
        return BookingStatus.active;
      case 'RETURN_INITIATED':
        return BookingStatus.returnInitiated;
      case 'PENDING_FINAL_PAYMENT':
        return BookingStatus.pendingFinalPayment;
      case 'COMPLETED':
        return BookingStatus.completed;
      case 'RATED':
        return BookingStatus.rated;
      case 'CANCELLED':
        return BookingStatus.cancelled;
      default:
        return BookingStatus.pendingPayment;
    }
  }

  String toApiString() {
    switch (this) {
      case BookingStatus.draft:
        return 'DRAFT';
      case BookingStatus.pendingPayment:
        return 'PENDING_PAYMENT';
      case BookingStatus.paymentEvidenceUploaded:
        return 'PAYMENT_EVIDENCE_UPLOADED';
      case BookingStatus.paymentDeclined:
        return 'PAYMENT_DECLINED';
      case BookingStatus.disputeOpened:
        return 'DISPUTE_OPENED';
      case BookingStatus.disputeResolvedRenter:
        return 'DISPUTE_RESOLVED_RENTER';
      case BookingStatus.disputeResolvedOwner:
        return 'DISPUTE_RESOLVED_OWNER';
      case BookingStatus.confirmed:
        return 'CONFIRMED';
      case BookingStatus.awaitingFleetReadiness:
        return 'AWAITING_FLEET_READINESS';
      case BookingStatus.pendingAdditionalPayment:
        return 'PENDING_ADDITIONAL_PAYMENT';
      case BookingStatus.readyForPickup:
        return 'READY_FOR_PICKUP';
      case BookingStatus.active:
        return 'ACTIVE';
      case BookingStatus.returnInitiated:
        return 'RETURN_INITIATED';
      case BookingStatus.pendingFinalPayment:
        return 'PENDING_FINAL_PAYMENT';
      case BookingStatus.completed:
        return 'COMPLETED';
      case BookingStatus.rated:
        return 'RATED';
      case BookingStatus.cancelled:
        return 'CANCELLED';
    }
  }
}

class BookingAttachmentEntity extends Equatable {
  final String id;
  final String attachmentId;
  final String? attachmentType;
  final String? make;
  final String? model;
  final int bundledDailyRateMinor;
  final String currency;

  const BookingAttachmentEntity({
    required this.id,
    required this.attachmentId,
    this.attachmentType,
    this.make,
    this.model,
    required this.bundledDailyRateMinor,
    this.currency = 'ETB',
  });

  @override
  List<Object?> get props => [
        id,
        attachmentId,
        attachmentType,
        make,
        model,
        bundledDailyRateMinor,
        currency,
      ];
}

class BookingEntity extends Equatable {
  final String id;
  final String bookingReference;
  final String fleetId;
  final String renterId;
  final String ownerId;
  final BookingStatus status;
  final DateTime startDate;
  final DateTime endDate;
  final int durationDays;
  final int dailyRateMinor;
  final int baseHireFeeMinor;
  final int attachmentsFeeMinor;
  final int grossRentalAmountMinor;
  final int prepaymentAmountMinor;
  final int penaltyAmountMinor;
  final int refundAmountMinor;
  final String currency;
  final String siteLocationAddress;
  final String? rejectionReasonCode;
  final String? rejectionNotes;
  final String? cancellationReason;
  final String? handoverPin;
  final double? checkInMeterHours;
  final double? checkOutMeterHours;
  final DateTime? expiresAt;
  final DateTime? confirmedAt;
  final DateTime? activatedAt;
  final DateTime? completedAt;
  final DateTime? cancelledAt;
  final List<BookingAttachmentEntity> attachments;
  final int version;
  final DateTime createdAt;

  const BookingEntity({
    required this.id,
    required this.bookingReference,
    required this.fleetId,
    required this.renterId,
    required this.ownerId,
    required this.status,
    required this.startDate,
    required this.endDate,
    required this.durationDays,
    required this.dailyRateMinor,
    required this.baseHireFeeMinor,
    this.attachmentsFeeMinor = 0,
    required this.grossRentalAmountMinor,
    required this.prepaymentAmountMinor,
    this.penaltyAmountMinor = 0,
    this.refundAmountMinor = 0,
    this.currency = 'ETB',
    required this.siteLocationAddress,
    this.rejectionReasonCode,
    this.rejectionNotes,
    this.cancellationReason,
    this.handoverPin,
    this.checkInMeterHours,
    this.checkOutMeterHours,
    this.expiresAt,
    this.confirmedAt,
    this.activatedAt,
    this.completedAt,
    this.cancelledAt,
    this.attachments = const [],
    required this.version,
    required this.createdAt,
  });

  @override
  List<Object?> get props => [
        id,
        bookingReference,
        fleetId,
        renterId,
        ownerId,
        status,
        startDate,
        endDate,
        durationDays,
        dailyRateMinor,
        baseHireFeeMinor,
        attachmentsFeeMinor,
        grossRentalAmountMinor,
        prepaymentAmountMinor,
        penaltyAmountMinor,
        refundAmountMinor,
        currency,
        siteLocationAddress,
        rejectionReasonCode,
        rejectionNotes,
        cancellationReason,
        handoverPin,
        checkInMeterHours,
        checkOutMeterHours,
        expiresAt,
        confirmedAt,
        activatedAt,
        completedAt,
        cancelledAt,
        attachments,
        version,
        createdAt,
      ];

  BookingEntity copyWith({
    String? id,
    String? bookingReference,
    String? fleetId,
    String? renterId,
    String? ownerId,
    BookingStatus? status,
    DateTime? startDate,
    DateTime? endDate,
    int? durationDays,
    int? dailyRateMinor,
    int? baseHireFeeMinor,
    int? attachmentsFeeMinor,
    int? grossRentalAmountMinor,
    int? prepaymentAmountMinor,
    int? penaltyAmountMinor,
    int? refundAmountMinor,
    String? currency,
    String? siteLocationAddress,
    String? rejectionReasonCode,
    String? rejectionNotes,
    String? cancellationReason,
    String? handoverPin,
    double? checkInMeterHours,
    double? checkOutMeterHours,
    DateTime? expiresAt,
    DateTime? confirmedAt,
    DateTime? activatedAt,
    DateTime? completedAt,
    DateTime? cancelledAt,
    List<BookingAttachmentEntity>? attachments,
    int? version,
    DateTime? createdAt,
  }) {
    return BookingEntity(
      id: id ?? this.id,
      bookingReference: bookingReference ?? this.bookingReference,
      fleetId: fleetId ?? this.fleetId,
      renterId: renterId ?? this.renterId,
      ownerId: ownerId ?? this.ownerId,
      status: status ?? this.status,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
      durationDays: durationDays ?? this.durationDays,
      dailyRateMinor: dailyRateMinor ?? this.dailyRateMinor,
      baseHireFeeMinor: baseHireFeeMinor ?? this.baseHireFeeMinor,
      attachmentsFeeMinor: attachmentsFeeMinor ?? this.attachmentsFeeMinor,
      grossRentalAmountMinor: grossRentalAmountMinor ?? this.grossRentalAmountMinor,
      prepaymentAmountMinor: prepaymentAmountMinor ?? this.prepaymentAmountMinor,
      penaltyAmountMinor: penaltyAmountMinor ?? this.penaltyAmountMinor,
      refundAmountMinor: refundAmountMinor ?? this.refundAmountMinor,
      currency: currency ?? this.currency,
      siteLocationAddress: siteLocationAddress ?? this.siteLocationAddress,
      rejectionReasonCode: rejectionReasonCode ?? this.rejectionReasonCode,
      rejectionNotes: rejectionNotes ?? this.rejectionNotes,
      cancellationReason: cancellationReason ?? this.cancellationReason,
      handoverPin: handoverPin ?? this.handoverPin,
      checkInMeterHours: checkInMeterHours ?? this.checkInMeterHours,
      checkOutMeterHours: checkOutMeterHours ?? this.checkOutMeterHours,
      expiresAt: expiresAt ?? this.expiresAt,
      confirmedAt: confirmedAt ?? this.confirmedAt,
      activatedAt: activatedAt ?? this.activatedAt,
      completedAt: completedAt ?? this.completedAt,
      cancelledAt: cancelledAt ?? this.cancelledAt,
      attachments: attachments ?? this.attachments,
      version: version ?? this.version,
      createdAt: createdAt ?? this.createdAt,
    );
  }
}

class BankSlipEntity extends Equatable {
  final String id;
  final String bookingId;
  final String uploaderId;
  final String storageUri;
  final String fileHash;
  final String? transactionReference;
  final int amountMinor;
  final String currency;
  final String status;
  final DateTime createdAt;

  const BankSlipEntity({
    required this.id,
    required this.bookingId,
    required this.uploaderId,
    required this.storageUri,
    required this.fileHash,
    this.transactionReference,
    required this.amountMinor,
    this.currency = 'ETB',
    required this.status,
    required this.createdAt,
  });

  @override
  List<Object?> get props => [
        id,
        bookingId,
        uploaderId,
        storageUri,
        fileHash,
        transactionReference,
        amountMinor,
        currency,
        status,
        createdAt,
      ];
}
