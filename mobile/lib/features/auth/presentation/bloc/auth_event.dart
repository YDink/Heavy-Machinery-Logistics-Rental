import 'package:equatable/equatable.dart';

abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object?> get props => [];
}

class CheckAuthStatusEvent extends AuthEvent {}

class RequestOtpEvent extends AuthEvent {
  final String phoneNumber;
  final String countryCode;

  const RequestOtpEvent({
    required this.phoneNumber,
    this.countryCode = 'ET',
  });

  @override
  List<Object?> get props => [phoneNumber, countryCode];
}

class VerifyOtpEvent extends AuthEvent {
  final String code;

  const VerifyOtpEvent({required this.code});

  @override
  List<Object?> get props => [code];
}

class LogoutEvent extends AuthEvent {}
