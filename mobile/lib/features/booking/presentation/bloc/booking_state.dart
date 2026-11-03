import 'package:equatable/equatable.dart';
import '../../domain/entities/booking.dart';

abstract class BookingState extends Equatable {
  const BookingState();

  @override
  List<Object?> get props => [];
}

class BookingInitial extends BookingState {
  const BookingInitial();
}

class BookingLoading extends BookingState {
  final String? message;

  const BookingLoading({this.message});

  @override
  List<Object?> get props => [message];
}

class BookingLoaded extends BookingState {
  final BookingEntity booking;
  final Duration? remainingHoldTime;

  const BookingLoaded({
    required this.booking,
    this.remainingHoldTime,
  });

  @override
  List<Object?> get props => [booking, remainingHoldTime];
}

class BookingError extends BookingState {
  final String message;
  final String? code;

  const BookingError(this.message, {this.code});

  @override
  List<Object?> get props => [message, code];
}

class BookingActionSuccess extends BookingState {
  final String message;
  final BookingEntity booking;

  const BookingActionSuccess(this.message, this.booking);

  @override
  List<Object?> get props => [message, booking];
}
