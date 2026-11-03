import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/repositories/active_rental_repository.dart';
import 'active_rental_event.dart';
import 'active_rental_state.dart';

class ActiveRentalBloc extends Bloc<ActiveRentalEvent, ActiveRentalState> {
  final ActiveRentalRepository repository;

  ActiveRentalBloc({required this.repository}) : super(ActiveRentalInitial()) {
    on<LoadActiveRentalHUDEvent>(_onLoadHUD);
    on<SubmitHandoverInspectionEvent>(_onSubmitInspection);
    on<GenerateHandoverChallengeEvent>(_onGenerateChallenge);
    on<ConfirmHandoverEvent>(_onConfirmHandover);
    on<RequestRentalExtensionEvent>(_onRequestExtension);
    on<ReportMechanicalBreakdownEvent>(_onReportBreakdown);
    on<ResolveMechanicalBreakdownEvent>(_onResolveBreakdown);
  }

  Future<void> _onLoadHUD(
    LoadActiveRentalHUDEvent event,
    Emitter<ActiveRentalState> emit,
  ) async {
    emit(ActiveRentalLoading());
    try {
      final hud = await repository.getHUD(event.bookingId);
      emit(ActiveRentalHUDLoaded(hud));
    } catch (e) {
      emit(ActiveRentalError(e.toString()));
    }
  }

  Future<void> _onSubmitInspection(
    SubmitHandoverInspectionEvent event,
    Emitter<ActiveRentalState> emit,
  ) async {
    emit(ActiveRentalLoading());
    try {
      await repository.inspectHandover(
        bookingId: event.bookingId,
        checkInMeterHours: event.checkInMeterHours,
        fluidLevelsOk: event.fluidLevelsOk,
        safetyGlassIntact: event.safetyGlassIntact,
        notes: event.notes,
        photos: event.photos,
      );
      emit(const ActiveRentalActionSuccess('Handover inspection submitted successfully.'));
      add(LoadActiveRentalHUDEvent(event.bookingId));
    } catch (e) {
      emit(ActiveRentalError(e.toString()));
    }
  }

  Future<void> _onGenerateChallenge(
    GenerateHandoverChallengeEvent event,
    Emitter<ActiveRentalState> emit,
  ) async {
    emit(ActiveRentalLoading());
    try {
      final challenge = await repository.generateChallenge(event.bookingId);
      emit(HandoverChallengeGenerated(challenge));
    } catch (e) {
      emit(ActiveRentalError(e.toString()));
    }
  }

  Future<void> _onConfirmHandover(
    ConfirmHandoverEvent event,
    Emitter<ActiveRentalState> emit,
  ) async {
    emit(ActiveRentalLoading());
    try {
      await repository.confirmHandover(
        bookingId: event.bookingId,
        challengePin: event.challengePin,
        touchSignatureBase64: event.touchSignatureBase64,
        latitude: event.latitude,
        longitude: event.longitude,
      );
      emit(const ActiveRentalActionSuccess('Rental activated successfully!'));
      add(LoadActiveRentalHUDEvent(event.bookingId));
    } catch (e) {
      emit(ActiveRentalError(e.toString()));
    }
  }

  Future<void> _onRequestExtension(
    RequestRentalExtensionEvent event,
    Emitter<ActiveRentalState> emit,
  ) async {
    emit(ActiveRentalLoading());
    try {
      await repository.requestExtension(
        bookingId: event.bookingId,
        newEndDate: event.newEndDate,
      );
      emit(const ActiveRentalActionSuccess('Extension requested successfully (pending owner approval).'));
      add(LoadActiveRentalHUDEvent(event.bookingId));
    } catch (e) {
      emit(ActiveRentalError(e.toString()));
    }
  }

  Future<void> _onReportBreakdown(
    ReportMechanicalBreakdownEvent event,
    Emitter<ActiveRentalState> emit,
  ) async {
    emit(ActiveRentalLoading());
    try {
      await repository.reportBreakdown(
        bookingId: event.bookingId,
        faultCategory: event.faultCategory,
        faultDescription: event.faultDescription,
        diagnosticTroubleCodes: event.diagnosticTroubleCodes,
        currentHourMeter: event.currentHourMeter,
        incidentPhotos: event.incidentPhotos,
      );
      emit(const ActiveRentalActionSuccess('Breakdown reported. Shift timer paused and 4h owner SLA activated.'));
      add(LoadActiveRentalHUDEvent(event.bookingId));
    } catch (e) {
      emit(ActiveRentalError(e.toString()));
    }
  }

  Future<void> _onResolveBreakdown(
    ResolveMechanicalBreakdownEvent event,
    Emitter<ActiveRentalState> emit,
  ) async {
    emit(ActiveRentalLoading());
    try {
      await repository.resolveBreakdown(
        bookingId: event.bookingId,
        mechanicName: event.mechanicName,
        mechanicPhone: event.mechanicPhone,
        resolutionNotes: event.resolutionNotes,
      );
      emit(const ActiveRentalActionSuccess('Breakdown resolved and shift telemetry resumed.'));
      add(LoadActiveRentalHUDEvent(event.bookingId));
    } catch (e) {
      emit(ActiveRentalError(e.toString()));
    }
  }
}
