import 'package:equatable/equatable.dart';
import '../../domain/entities/winching_photo.dart';

abstract class LogisticsEvent extends Equatable {
  const LogisticsEvent();

  @override
  List<Object?> get props => [];
}

class MatchLogisticsCarriersRequested extends LogisticsEvent {
  final String bookingId;
  final int searchRadiusKm;
  final DateTime pickupDate;

  const MatchLogisticsCarriersRequested({
    required this.bookingId,
    required this.searchRadiusKm,
    required this.pickupDate,
  });

  @override
  List<Object?> get props => [bookingId, searchRadiusKm, pickupDate];
}

class CreateLogisticsSubBookingRequested extends LogisticsEvent {
  final String bookingId;
  final String logisticsFleetId;
  final DateTime pickupDateTime;
  final double originLat;
  final double originLng;
  final String originName;
  final double destLat;
  final double destLng;
  final String destName;

  const CreateLogisticsSubBookingRequested({
    required this.bookingId,
    required this.logisticsFleetId,
    required this.pickupDateTime,
    required this.originLat,
    required this.originLng,
    required this.originName,
    required this.destLat,
    required this.destLng,
    required this.destName,
  });

  @override
  List<Object?> get props => [bookingId, logisticsFleetId, pickupDateTime];
}

class RegisterExternalTransportRequested extends LogisticsEvent {
  final String bookingId;
  final String carrierName;
  final String driverFullName;
  final String driverPhoneNumber;
  final String primeMoverPlateNumber;
  final DateTime estimatedArrivalDateTime;

  const RegisterExternalTransportRequested({
    required this.bookingId,
    required this.carrierName,
    required this.driverFullName,
    required this.driverPhoneNumber,
    required this.primeMoverPlateNumber,
    required this.estimatedArrivalDateTime,
  });

  @override
  List<Object?> get props => [bookingId, driverFullName, primeMoverPlateNumber];
}

class AdvanceLogisticsDispatchRequested extends LogisticsEvent {
  final String subBookingId;
  final String targetStage;
  final List<WinchingPhoto> winchingPhotos;
  final String? notes;

  const AdvanceLogisticsDispatchRequested({
    required this.subBookingId,
    required this.targetStage,
    this.winchingPhotos = const [],
    this.notes,
  });

  @override
  List<Object?> get props => [subBookingId, targetStage, winchingPhotos];
}

class ReportTransshipmentAlertRequested extends LogisticsEvent {
  final String subBookingId;
  final double latitude;
  final double longitude;
  final String failureCategory;
  final String failureDescription;
  final List<String> photos;

  const ReportTransshipmentAlertRequested({
    required this.subBookingId,
    required this.latitude,
    required this.longitude,
    required this.failureCategory,
    required this.failureDescription,
    required this.photos,
  });

  @override
  List<Object?> get props => [subBookingId, failureCategory];
}

class AcceptTransshipmentRequested extends LogisticsEvent {
  final String subBookingId;
  final String replacementFleetId;

  const AcceptTransshipmentRequested({
    required this.subBookingId,
    required this.replacementFleetId,
  });

  @override
  List<Object?> get props => [subBookingId, replacementFleetId];
}

class TriggerReverseCascadeRequested extends LogisticsEvent {
  final String parentBookingId;
  final String? reason;

  const TriggerReverseCascadeRequested({
    required this.parentBookingId,
    this.reason,
  });

  @override
  List<Object?> get props => [parentBookingId, reason];
}
