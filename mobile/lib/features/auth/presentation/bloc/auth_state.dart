import 'package:equatable/equatable.dart';
import 'package:dink_mobile/features/auth/domain/entities/auth_tokens.dart';
import 'package:dink_mobile/features/auth/domain/entities/auth_user.dart';

abstract class AuthState extends Equatable {
  const AuthState();

  @override
  List<Object?> get props => [];
}

class AuthInitialState extends AuthState {}

class AuthLoadingState extends AuthState {}

class OtpSentState extends AuthState {
  final String verificationSessionId;
  final String phoneNumber;

  const OtpSentState({
    required this.verificationSessionId,
    required this.phoneNumber,
  });

  @override
  List<Object?> get props => [verificationSessionId, phoneNumber];
}

class AuthenticatedState extends AuthState {
  final AuthUser? user;
  final AuthTokens tokens;

  const AuthenticatedState({
    this.user,
    required this.tokens,
  });

  @override
  List<Object?> get props => [user, tokens];
}

class UnauthenticatedState extends AuthState {}

class AuthErrorState extends AuthState {
  final String message;
  final String code;

  const AuthErrorState({
    required this.message,
    this.code = 'ERR_AUTH',
  });

  @override
  List<Object?> get props => [message, code];
}
