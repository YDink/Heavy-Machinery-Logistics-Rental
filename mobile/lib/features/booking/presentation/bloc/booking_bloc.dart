import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/repositories/booking_repository.dart';
import 'booking_event.dart';
import 'booking_state.dart';

class BookingBloc extends Bloc<BookingEvent, BookingState> {
  final BookingRepository repository;

  BookingBloc({required this.repository}) : super(const BookingInitial()) {
    on<CreateBookingRequested>(_onCreateBookingRequested);
    on<LoadBookingDetails>(_onLoadBookingDetails);
    on<SubmitPaymentSlipRequested>(_onSubmitPaymentSlipRequested);
    on<ConfirmPaymentRequested>(_onConfirmPaymentRequested);
    on<DeclinePaymentRequested>(_onDeclinePaymentRequested);
    on<DisputeBookingRequested>(_onDisputeBookingRequested);
    on<CancelBookingRequested>(_onCancelBookingRequested);
    on<HandoverCheckInRequested>(_onHandoverCheckInRequested);
    on<ReturnCheckOutRequested>(_onReturnCheckOutRequested);
  }

  Future<void> _onCreateBookingRequested(
    CreateBookingRequested event,
    Emitter<BookingState> emit,
  ) async {
    emit(const BookingLoading(message: 'Reserving equipment hold...'));
    try {
      final booking = await repository.createBooking(
        fleetId: event.fleetId,
        startDate: event.startDate,
        endDate: event.endDate,
        siteLocationAddress: event.siteLocationAddress,
        attachmentIds: event.attachmentIds,
      );
      emit(BookingLoaded(
        booking: booking,
        remainingHoldTime: booking.expiresAt?.difference(DateTime.now()),
      ));
    } catch (e) {
      emit(BookingError(e.toString()));
    }
  }

  Future<void> _onLoadBookingDetails(
    LoadBookingDetails event,
    Emitter<BookingState> emit,
  ) async {
    emit(const BookingLoading());
    try {
      final booking = await repository.getBooking(event.bookingId);
      emit(BookingLoaded(
        booking: booking,
        remainingHoldTime: booking.expiresAt?.difference(DateTime.now()),
      ));
    } catch (e) {
      emit(BookingError(e.toString()));
    }
  }

  Future<void> _onSubmitPaymentSlipRequested(
    SubmitPaymentSlipRequested event,
    Emitter<BookingState> emit,
  ) async {
    emit(const BookingLoading(message: 'Submitting bank wire receipt...'));
    try {
      await repository.submitPaymentEvidence(
        bookingId: event.bookingId,
        storageUri: event.storageUri,
        fileHash: event.fileHash,
        transactionReference: event.transactionReference,
        amountMinor: event.amountMinor,
      );
      final refreshed = await repository.getBooking(event.bookingId);
      emit(BookingActionSuccess('Payment evidence submitted for owner audit.', refreshed));
      emit(BookingLoaded(booking: refreshed));
    } catch (e) {
      emit(BookingError(e.toString()));
    }
  }

  Future<void> _onConfirmPaymentRequested(
    ConfirmPaymentRequested event,
    Emitter<BookingState> emit,
  ) async {
    emit(const BookingLoading(message: 'Confirming bank wire receipt...'));
    try {
      final updated = await repository.confirmPayment(
        bookingId: event.bookingId,
        expectedVersion: event.expectedVersion,
        paymentReference: event.paymentReference,
        notes: event.notes,
      );
      emit(BookingActionSuccess('Payment confirmed. Booking active.', updated));
      emit(BookingLoaded(booking: updated));
    } catch (e) {
      emit(BookingError(e.toString()));
    }
  }

  Future<void> _onDeclinePaymentRequested(
    DeclinePaymentRequested event,
    Emitter<BookingState> emit,
  ) async {
    emit(const BookingLoading(message: 'Declining payment evidence...'));
    try {
      final updated = await repository.declinePayment(
        bookingId: event.bookingId,
        expectedVersion: event.expectedVersion,
        reasonCode: event.reasonCode,
        rejectionNotes: event.rejectionNotes,
      );
      emit(BookingActionSuccess('Payment evidence declined.', updated));
      emit(BookingLoaded(booking: updated));
    } catch (e) {
      emit(BookingError(e.toString()));
    }
  }

  Future<void> _onDisputeBookingRequested(
    DisputeBookingRequested event,
    Emitter<BookingState> emit,
  ) async {
    emit(const BookingLoading(message: 'Opening dispute docket...'));
    try {
      final updated = await repository.disputeBooking(
        bookingId: event.bookingId,
        expectedVersion: event.expectedVersion,
        disputeReason: event.disputeReason,
        evidenceUris: event.evidenceUris,
      );
      emit(BookingActionSuccess('Dispute submitted to compliance arbitrator.', updated));
      emit(BookingLoaded(booking: updated));
    } catch (e) {
      emit(BookingError(e.toString()));
    }
  }

  Future<void> _onCancelBookingRequested(
    CancelBookingRequested event,
    Emitter<BookingState> emit,
  ) async {
    emit(const BookingLoading(message: 'Cancelling booking...'));
    try {
      final updated = await repository.cancelBooking(
        bookingId: event.bookingId,
        expectedVersion: event.expectedVersion,
        reason: event.reason,
      );
      emit(BookingActionSuccess('Booking cancelled.', updated));
      emit(BookingLoaded(booking: updated));
    } catch (e) {
      emit(BookingError(e.toString()));
    }
  }

  Future<void> _onHandoverCheckInRequested(
    HandoverCheckInRequested event,
    Emitter<BookingState> emit,
  ) async {
    emit(const BookingLoading(message: 'Validating challenge PIN...'));
    try {
      final updated = await repository.handoverCheckIn(
        bookingId: event.bookingId,
        expectedVersion: event.expectedVersion,
        handoverPin: event.handoverPin,
        operatingHoursMeter: event.operatingHoursMeter,
        photoUris: event.photoUris,
      );
      emit(BookingActionSuccess('Handover completed. Rental active.', updated));
      emit(BookingLoaded(booking: updated));
    } catch (e) {
      emit(BookingError(e.toString()));
    }
  }

  Future<void> _onReturnCheckOutRequested(
    ReturnCheckOutRequested event,
    Emitter<BookingState> emit,
  ) async {
    emit(const BookingLoading(message: 'Completing machine return...'));
    try {
      final updated = await repository.returnCheckOut(
        bookingId: event.bookingId,
        expectedVersion: event.expectedVersion,
        operatingHoursMeter: event.operatingHoursMeter,
        damageNotes: event.damageNotes,
        photoUris: event.photoUris,
      );
      emit(BookingActionSuccess('Rental return check-out recorded.', updated));
      emit(BookingLoaded(booking: updated));
    } catch (e) {
      emit(BookingError(e.toString()));
    }
  }
}
