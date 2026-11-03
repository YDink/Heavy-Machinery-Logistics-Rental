import 'package:equatable/equatable.dart';

abstract class BookingEvent extends Equatable {
  const BookingEvent();

  @override
  List<Object?> get props => [];
}

class CreateBookingRequested extends BookingEvent {
  final String fleetId;
  final DateTime startDate;
  final DateTime endDate;
  final String siteLocationAddress;
  final List<String> attachmentIds;

  const CreateBookingRequested({
    required this.fleetId,
    required this.startDate,
    required this.endDate,
    required this.siteLocationAddress,
    this.attachmentIds = const [],
  });

  @override
  List<Object?> get props => [fleetId, startDate, endDate, siteLocationAddress, attachmentIds];
}

class LoadBookingDetails extends BookingEvent {
  final String bookingId;

  const LoadBookingDetails(this.bookingId);

  @override
  List<Object?> get props => [bookingId];
}

class SubmitPaymentSlipRequested extends BookingEvent {
  final String bookingId;
  final String storageUri;
  final String fileHash;
  final String? transactionReference;
  final int amountMinor;

  const SubmitPaymentSlipRequested({
    required this.bookingId,
    required this.storageUri,
    required this.fileHash,
    this.transactionReference,
    required this.amountMinor,
  });

  @override
  List<Object?> get props => [bookingId, storageUri, fileHash, transactionReference, amountMinor];
}

class ConfirmPaymentRequested extends BookingEvent {
  final String bookingId;
  final int expectedVersion;
  final String? paymentReference;
  final String? notes;

  const ConfirmPaymentRequested({
    required this.bookingId,
    required this.expectedVersion,
    this.paymentReference,
    this.notes,
  });

  @override
  List<Object?> get props => [bookingId, expectedVersion, paymentReference, notes];
}

class DeclinePaymentRequested extends BookingEvent {
  final String bookingId;
  final int expectedVersion;
  final String reasonCode;
  final String rejectionNotes;

  const DeclinePaymentRequested({
    required this.bookingId,
    required this.expectedVersion,
    required this.reasonCode,
    required this.rejectionNotes,
  });

  @override
  List<Object?> get props => [bookingId, expectedVersion, reasonCode, rejectionNotes];
}

class DisputeBookingRequested extends BookingEvent {
  final String bookingId;
  final int expectedVersion;
  final String disputeReason;
  final List<String> evidenceUris;

  const DisputeBookingRequested({
    required this.bookingId,
    required this.expectedVersion,
    required this.disputeReason,
    this.evidenceUris = const [],
  });

  @override
  List<Object?> get props => [bookingId, expectedVersion, disputeReason, evidenceUris];
}

class CancelBookingRequested extends BookingEvent {
  final String bookingId;
  final int expectedVersion;
  final String reason;

  const CancelBookingRequested({
    required this.bookingId,
    required this.expectedVersion,
    required this.reason,
  });

  @override
  List<Object?> get props => [bookingId, expectedVersion, reason];
}

class HandoverCheckInRequested extends BookingEvent {
  final String bookingId;
  final int expectedVersion;
  final String handoverPin;
  final double operatingHoursMeter;
  final List<String> photoUris;

  const HandoverCheckInRequested({
    required this.bookingId,
    required this.expectedVersion,
    required this.handoverPin,
    required this.operatingHoursMeter,
    this.photoUris = const [],
  });

  @override
  List<Object?> get props => [bookingId, expectedVersion, handoverPin, operatingHoursMeter, photoUris];
}

class ReturnCheckOutRequested extends BookingEvent {
  final String bookingId;
  final int expectedVersion;
  final double operatingHoursMeter;
  final String? damageNotes;
  final List<String> photoUris;

  const ReturnCheckOutRequested({
    required this.bookingId,
    required this.expectedVersion,
    required this.operatingHoursMeter,
    this.damageNotes,
    this.photoUris = const [],
  });

  @override
  List<Object?> get props => [bookingId, expectedVersion, operatingHoursMeter, damageNotes, photoUris];
}
