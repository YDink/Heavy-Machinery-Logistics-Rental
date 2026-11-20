import 'package:equatable/equatable.dart';
import '../../domain/entities/external_transport.dart';
import '../../domain/entities/logistics_dispatch_result.dart';
import '../../domain/entities/logistics_match_result.dart';
import '../../domain/entities/logistics_sub_booking_result.dart';
import '../../domain/entities/reverse_cascade_result.dart';
import '../../domain/entities/transshipment_alert_result.dart';

abstract class LogisticsState extends Equatable {
  const LogisticsState();

  @override
  List<Object?> get props => [];
}

class LogisticsInitial extends LogisticsState {
  const LogisticsInitial();
}

class LogisticsLoading extends LogisticsState {
  const LogisticsLoading();
}

class LogisticsCarriersMatched extends LogisticsState {
  final LogisticsMatchResult result;

  const LogisticsCarriersMatched(this.result);

  @override
  List<Object?> get props => [result];
}

class LogisticsSubBookingCreated extends LogisticsState {
  final LogisticsSubBookingResult subBooking;

  const LogisticsSubBookingCreated(this.subBooking);

  @override
  List<Object?> get props => [subBooking];
}

class ExternalTransportRegistered extends LogisticsState {
  final ExternalTransport externalTransport;

  const ExternalTransportRegistered(this.externalTransport);

  @override
  List<Object?> get props => [externalTransport];
}

class LogisticsDispatchUpdated extends LogisticsState {
  final LogisticsDispatchResult dispatchResult;

  const LogisticsDispatchUpdated(this.dispatchResult);

  @override
  List<Object?> get props => [dispatchResult];
}

class TransshipmentAlertOpened extends LogisticsState {
  final TransshipmentAlertResult alertResult;

  const TransshipmentAlertOpened(this.alertResult);

  @override
  List<Object?> get props => [alertResult];
}

class TransshipmentHotSwapAccepted extends LogisticsState {
  final String subBookingStatus;

  const TransshipmentHotSwapAccepted(this.subBookingStatus);

  @override
  List<Object?> get props => [subBookingStatus];
}

class ReverseCascadeExecuted extends LogisticsState {
  final ReverseCascadeResult cascadeResult;

  const ReverseCascadeExecuted(this.cascadeResult);

  @override
  List<Object?> get props => [cascadeResult];
}

class LogisticsFailure extends LogisticsState {
  final String message;

  const LogisticsFailure(this.message);

  @override
  List<Object?> get props => [message];
}
