import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/core/errors/failure.dart';
import 'package:dink_mobile/features/auth/domain/entities/auth_tokens.dart';
import 'package:dink_mobile/features/auth/domain/entities/auth_user.dart';
import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';
import 'package:dink_mobile/features/auth/domain/usecases/get_current_user_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/request_otp_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/verify_otp_usecase.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_bloc.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_event.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_state.dart';

class MockAuthRepository implements AuthRepository {
  bool shouldFail = false;
  String failMessage = 'Error';
  String failCode = 'ERR_TEST';
  bool authenticated = false;

  @override
  Future<String> requestOtp({required String phoneNumber, required String purpose, String? countryCode}) async {
    if (shouldFail) {
      throw ServerFailure(message: failMessage, code: failCode);
    }
    return 'session-uuid-1234';
  }

  @override
  Future<AuthTokens> verifyOtp({required String verificationSessionId, required String code}) async {
    if (shouldFail) {
      throw ServerFailure(message: failMessage, code: failCode);
    }
    authenticated = true;
    return const AuthTokens(
      accessToken: 'test-access-token',
      refreshToken: 'test-refresh-token',
      tokenType: 'Bearer',
      expiresInSeconds: 900,
    );
  }

  @override
  Future<AuthTokens> refreshSession({required String refreshToken}) async {
    return const AuthTokens(
      accessToken: 'refreshed-access-token',
      refreshToken: 'refreshed-refresh-token',
      tokenType: 'Bearer',
      expiresInSeconds: 900,
    );
  }

  @override
  Future<AuthUser> getCurrentUser() async {
    return const AuthUser(
      id: 'user-uuid-5678',
      phoneNumber: '+251911442910',
      fullName: 'Abebe Bikila',
      role: UserRole.renter,
      status: AccountStatus.active,
    );
  }

  @override
  Future<void> logout() async {
    authenticated = false;
  }

  @override
  Future<bool> isAuthenticated() async {
    return authenticated;
  }
}

void main() {
  late MockAuthRepository mockRepo;
  late RequestOtpUseCase requestOtpUseCase;
  late VerifyOtpUseCase verifyOtpUseCase;
  late GetCurrentUserUseCase getCurrentUserUseCase;
  late AuthBloc authBloc;

  setUp(() {
    mockRepo = MockAuthRepository();
    requestOtpUseCase = RequestOtpUseCase(mockRepo);
    verifyOtpUseCase = VerifyOtpUseCase(mockRepo);
    getCurrentUserUseCase = GetCurrentUserUseCase(mockRepo);
    authBloc = AuthBloc(
      requestOtpUseCase: requestOtpUseCase,
      verifyOtpUseCase: verifyOtpUseCase,
      getCurrentUserUseCase: getCurrentUserUseCase,
      authRepository: mockRepo,
    );
  });

  tearDown(() {
    authBloc.close();
  });

  test('initial state is AuthInitialState', () {
    expect(authBloc.state, equals(AuthInitialState()));
  });

  test('emits [AuthLoadingState, OtpSentState] when RequestOtp succeeds', () async {
    final expectedStates = [
      AuthLoadingState(),
      const OtpSentState(
        verificationSessionId: 'session-uuid-1234',
        phoneNumber: '+251911442910',
      ),
    ];

    expectLater(authBloc.stream, emitsInOrder(expectedStates));

    authBloc.add(const RequestOtpEvent(
      phoneNumber: '+251911442910',
      countryCode: 'ET',
    ));
  });

  test('emits [AuthLoadingState, AuthErrorState] when RequestOtp fails', () async {
    mockRepo.shouldFail = true;
    mockRepo.failMessage = 'Rate limit exceeded';
    mockRepo.failCode = 'ERR_RATE_LIMIT_EXCEEDED';

    final expectedStates = [
      AuthLoadingState(),
      const AuthErrorState(
        message: 'Rate limit exceeded',
        code: 'ERR_RATE_LIMIT_EXCEEDED',
      ),
    ];

    expectLater(authBloc.stream, emitsInOrder(expectedStates));

    authBloc.add(const RequestOtpEvent(
      phoneNumber: '+251911442910',
      countryCode: 'ET',
    ));
  });

  test('emits [AuthLoadingState, AuthenticatedState] when VerifyOtp succeeds', () async {
    // First request OTP to establish session ID
    authBloc.add(const RequestOtpEvent(
      phoneNumber: '+251911442910',
      countryCode: 'ET',
    ));
    await Future.delayed(const Duration(milliseconds: 50));

    final expectedStates = [
      AuthLoadingState(),
      const AuthenticatedState(
        tokens: AuthTokens(
          accessToken: 'test-access-token',
          refreshToken: 'test-refresh-token',
          tokenType: 'Bearer',
          expiresInSeconds: 900,
        ),
      ),
    ];

    expectLater(authBloc.stream, emitsInOrder(expectedStates));

    authBloc.add(const VerifyOtpEvent(code: '123456'));
  });

  test('emits [AuthLoadingState, UnauthenticatedState] on LogoutEvent', () async {
    final expectedStates = [
      AuthLoadingState(),
      UnauthenticatedState(),
    ];

    expectLater(authBloc.stream, emitsInOrder(expectedStates));

    authBloc.add(LogoutEvent());
  });
}
