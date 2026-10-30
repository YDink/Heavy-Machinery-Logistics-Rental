import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/core/errors/failure.dart';
import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';
import 'package:dink_mobile/features/auth/domain/usecases/get_current_user_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/request_otp_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/verify_otp_usecase.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_event.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final RequestOtpUseCase requestOtpUseCase;
  final VerifyOtpUseCase verifyOtpUseCase;
  final GetCurrentUserUseCase getCurrentUserUseCase;
  final AuthRepository authRepository;

  String? _currentSessionId;
  String? _currentPhoneNumber;

  String? get currentSessionId => _currentSessionId;
  String? get currentPhoneNumber => _currentPhoneNumber;

  AuthBloc({
    required this.requestOtpUseCase,
    required this.verifyOtpUseCase,
    required this.getCurrentUserUseCase,
    required this.authRepository,
  }) : super(AuthInitialState()) {
    on<CheckAuthStatusEvent>(_onCheckAuthStatus);
    on<RequestOtpEvent>(_onRequestOtp);
    on<VerifyOtpEvent>(_onVerifyOtp);
    on<LogoutEvent>(_onLogout);
  }

  Future<void> _onCheckAuthStatus(
    CheckAuthStatusEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoadingState());
    try {
      final isAuth = await authRepository.isAuthenticated();
      if (!isAuth) {
        emit(UnauthenticatedState());
        return;
      }

      final user = await getCurrentUserUseCase();
      final tokens = await authRepository.refreshSession(refreshToken: '');
      emit(AuthenticatedState(user: user, tokens: tokens));
    } catch (_) {
      emit(UnauthenticatedState());
    }
  }

  Future<void> _onRequestOtp(
    RequestOtpEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoadingState());
    try {
      final sessionId = await requestOtpUseCase(
        phoneNumber: event.phoneNumber,
        countryCode: event.countryCode,
      );
      _currentSessionId = sessionId;
      _currentPhoneNumber = event.phoneNumber;
      emit(OtpSentState(
        verificationSessionId: sessionId,
        phoneNumber: event.phoneNumber,
      ));
    } on Failure catch (e) {
      emit(AuthErrorState(message: e.message, code: e.code));
    } catch (e) {
      emit(AuthErrorState(message: e.toString()));
    }
  }

  Future<void> _onVerifyOtp(
    VerifyOtpEvent event,
    Emitter<AuthState> emit,
  ) async {
    final sessionId = _currentSessionId;
    if (sessionId == null) {
      emit(const AuthErrorState(
        message: 'No active verification session found. Please request OTP first.',
        code: 'ERR_NO_SESSION',
      ));
      return;
    }

    emit(AuthLoadingState());
    try {
      final tokens = await verifyOtpUseCase(
        verificationSessionId: sessionId,
        code: event.code,
      );

      emit(AuthenticatedState(tokens: tokens));
    } on Failure catch (e) {
      emit(AuthErrorState(message: e.message, code: e.code));
    } catch (e) {
      emit(AuthErrorState(message: e.toString()));
    }
  }

  Future<void> _onLogout(
    LogoutEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoadingState());
    await authRepository.logout();
    _currentSessionId = null;
    _currentPhoneNumber = null;
    emit(UnauthenticatedState());
  }
}
