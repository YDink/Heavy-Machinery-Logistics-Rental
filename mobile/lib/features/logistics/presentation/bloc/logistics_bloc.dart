import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/usecases/accept_transshipment_usecase.dart';
import '../../domain/usecases/advance_logistics_dispatch_usecase.dart';
import '../../domain/usecases/create_logistics_sub_booking_usecase.dart';
import '../../domain/usecases/match_logistics_carriers_usecase.dart';
import '../../domain/usecases/register_external_transport_usecase.dart';
import '../../domain/usecases/report_transshipment_alert_usecase.dart';
import '../../domain/usecases/trigger_reverse_cascade_usecase.dart';
import 'logistics_event.dart';
import 'logistics_state.dart';

class LogisticsBloc extends Bloc<LogisticsEvent, LogisticsState> {
  final MatchLogisticsCarriersUseCase matchCarriersUseCase;
  final CreateLogisticsSubBookingUseCase createSubBookingUseCase;
  final RegisterExternalTransportUseCase registerExternalTransportUseCase;
  final AdvanceLogisticsDispatchUseCase advanceDispatchUseCase;
  final ReportTransshipmentAlertUseCase reportTransshipmentAlertUseCase;
  final AcceptTransshipmentUseCase acceptTransshipmentUseCase;
  final TriggerReverseCascadeUseCase triggerReverseCascadeUseCase;

  LogisticsBloc({
    required this.matchCarriersUseCase,
    required this.createSubBookingUseCase,
    required this.registerExternalTransportUseCase,
    required this.advanceDispatchUseCase,
    required this.reportTransshipmentAlertUseCase,
    required this.acceptTransshipmentUseCase,
    required this.triggerReverseCascadeUseCase,
  }) : super(const LogisticsInitial()) {
    on<MatchLogisticsCarriersRequested>(_onMatchCarriers);
    on<CreateLogisticsSubBookingRequested>(_onCreateSubBooking);
    on<RegisterExternalTransportRequested>(_onRegisterExternalTransport);
    on<AdvanceLogisticsDispatchRequested>(_onAdvanceDispatch);
    on<ReportTransshipmentAlertRequested>(_onReportTransshipment);
    on<AcceptTransshipmentRequested>(_onAcceptTransshipment);
    on<TriggerReverseCascadeRequested>(_onTriggerReverseCascade);
  }

  Future<void> _onMatchCarriers(
    MatchLogisticsCarriersRequested event,
    Emitter<LogisticsState> emit,
  ) async {
    emit(const LogisticsLoading());
    try {
      final result = await matchCarriersUseCase(
        bookingId: event.bookingId,
        searchRadiusKm: event.searchRadiusKm,
        pickupDate: event.pickupDate,
      );
      emit(LogisticsCarriersMatched(result));
    } catch (e) {
      emit(LogisticsFailure(e.toString()));
    }
  }

  Future<void> _onCreateSubBooking(
    CreateLogisticsSubBookingRequested event,
    Emitter<LogisticsState> emit,
  ) async {
    emit(const LogisticsLoading());
    try {
      final result = await createSubBookingUseCase(
        bookingId: event.bookingId,
        logisticsFleetId: event.logisticsFleetId,
        pickupDateTime: event.pickupDateTime,
        originLat: event.originLat,
        originLng: event.originLng,
        originName: event.originName,
        destLat: event.destLat,
        destLng: event.destLng,
        destName: event.destName,
      );
      emit(LogisticsSubBookingCreated(result));
    } catch (e) {
      emit(LogisticsFailure(e.toString()));
    }
  }

  Future<void> _onRegisterExternalTransport(
    RegisterExternalTransportRequested event,
    Emitter<LogisticsState> emit,
  ) async {
    emit(const LogisticsLoading());
    try {
      final result = await registerExternalTransportUseCase(
        bookingId: event.bookingId,
        carrierName: event.carrierName,
        driverFullName: event.driverFullName,
        driverPhoneNumber: event.driverPhoneNumber,
        primeMoverPlateNumber: event.primeMoverPlateNumber,
        estimatedArrivalDateTime: event.estimatedArrivalDateTime,
      );
      emit(ExternalTransportRegistered(result));
    } catch (e) {
      emit(LogisticsFailure(e.toString()));
    }
  }

  Future<void> _onAdvanceDispatch(
    AdvanceLogisticsDispatchRequested event,
    Emitter<LogisticsState> emit,
  ) async {
    emit(const LogisticsLoading());
    try {
      final result = await advanceDispatchUseCase(
        subBookingId: event.subBookingId,
        targetStage: event.targetStage,
        winchingPhotos: event.winchingPhotos,
        notes: event.notes,
      );
      emit(LogisticsDispatchUpdated(result));
    } catch (e) {
      emit(LogisticsFailure(e.toString()));
    }
  }

  Future<void> _onReportTransshipment(
    ReportTransshipmentAlertRequested event,
    Emitter<LogisticsState> emit,
  ) async {
    emit(const LogisticsLoading());
    try {
      final result = await reportTransshipmentAlertUseCase(
        subBookingId: event.subBookingId,
        latitude: event.latitude,
        longitude: event.longitude,
        failureCategory: event.failureCategory,
        failureDescription: event.failureDescription,
        photos: event.photos,
      );
      emit(TransshipmentAlertOpened(result));
    } catch (e) {
      emit(LogisticsFailure(e.toString()));
    }
  }

  Future<void> _onAcceptTransshipment(
    AcceptTransshipmentRequested event,
    Emitter<LogisticsState> emit,
  ) async {
    emit(const LogisticsLoading());
    try {
      final status = await acceptTransshipmentUseCase(
        subBookingId: event.subBookingId,
        replacementFleetId: event.replacementFleetId,
      );
      emit(TransshipmentHotSwapAccepted(status));
    } catch (e) {
      emit(LogisticsFailure(e.toString()));
    }
  }

  Future<void> _onTriggerReverseCascade(
    TriggerReverseCascadeRequested event,
    Emitter<LogisticsState> emit,
  ) async {
    emit(const LogisticsLoading());
    try {
      final result = await triggerReverseCascadeUseCase(
        parentBookingId: event.parentBookingId,
        reason: event.reason,
      );
      emit(ReverseCascadeExecuted(result));
    } catch (e) {
      emit(LogisticsFailure(e.toString()));
    }
  }
}
