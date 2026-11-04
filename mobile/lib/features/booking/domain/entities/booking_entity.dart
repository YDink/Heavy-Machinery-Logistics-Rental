import 'package:equatable/equatable.dart';
import 'booking_attachment_entity.dart';
import 'booking_status.dart';

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
