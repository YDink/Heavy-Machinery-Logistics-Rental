import 'package:equatable/equatable.dart';
import '../../domain/entities/handover_inspection.dart';

abstract class ActiveRentalEvent extends Equatable {
  const ActiveRentalEvent();

  @override
  List<Object?> get props => [];
}

class LoadActiveRentalHUDEvent extends ActiveRentalEvent {
  final String bookingId;

  const LoadActiveRentalHUDEvent(this.bookingId);

  @override
  List<Object?> get props => [bookingId];
}

class SubmitHandoverInspectionEvent extends ActiveRentalEvent {
  final String bookingId;
  final double checkInMeterHours;
  final bool fluidLevelsOk;
  final bool safetyGlassIntact;
  final String? notes;
  final List<HandoverPhotoItem> photos;

  const SubmitHandoverInspectionEvent({
    required this.bookingId,
    required this.checkInMeterHours,
    required this.fluidLevelsOk,
    required this.safetyGlassIntact,
    this.notes,
    required this.photos,
  });

  @override
  List<Object?> get props => [
        bookingId,
        checkInMeterHours,
        fluidLevelsOk,
        safetyGlassIntact,
        notes,
        photos,
      ];
}

class GenerateHandoverChallengeEvent extends ActiveRentalEvent {
  final String bookingId;

  const GenerateHandoverChallengeEvent(this.bookingId);

  @override
  List<Object?> get props => [bookingId];
}

class ConfirmHandoverEvent extends ActiveRentalEvent {
  final String bookingId;
  final String challengePin;
  final String? touchSignatureBase64;
  final double? latitude;
  final double? longitude;

  const ConfirmHandoverEvent({
    required this.bookingId,
    required this.challengePin,
    this.touchSignatureBase64,
    this.latitude,
    this.longitude,
  });

  @override
  List<Object?> get props => [
        bookingId,
        challengePin,
        touchSignatureBase64,
        latitude,
        longitude,
      ];
}

class RequestRentalExtensionEvent extends ActiveRentalEvent {
  final String bookingId;
  final DateTime newEndDate;

  const RequestRentalExtensionEvent({
    required this.bookingId,
    required this.newEndDate,
  });

  @override
  List<Object?> get props => [bookingId, newEndDate];
}

class ReportMechanicalBreakdownEvent extends ActiveRentalEvent {
  final String bookingId;
  final String faultCategory;
  final String faultDescription;
  final List<String> diagnosticTroubleCodes;
  final double currentHourMeter;
  final List<String> incidentPhotos;

  const ReportMechanicalBreakdownEvent({
    required this.bookingId,
    required this.faultCategory,
    required this.faultDescription,
    this.diagnosticTroubleCodes = const [],
    required this.currentHourMeter,
    this.incidentPhotos = const [],
  });

  @override
  List<Object?> get props => [
        bookingId,
        faultCategory,
        faultDescription,
        diagnosticTroubleCodes,
        currentHourMeter,
        incidentPhotos,
      ];
}

class ResolveMechanicalBreakdownEvent extends ActiveRentalEvent {
  final String bookingId;
  final String mechanicName;
  final String mechanicPhone;
  final String resolutionNotes;

  const ResolveMechanicalBreakdownEvent({
    required this.bookingId,
    required this.mechanicName,
    required this.mechanicPhone,
    required this.resolutionNotes,
  });

  @override
  List<Object?> get props => [
        bookingId,
        mechanicName,
        mechanicPhone,
        resolutionNotes,
      ];
}
