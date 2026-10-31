import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/core/config/app_config.dart';
import 'package:dink_mobile/features/auth/domain/entities/auth_tokens.dart';
import 'package:dink_mobile/features/auth/domain/entities/auth_user.dart';
import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';
import 'package:dink_mobile/features/auth/domain/usecases/get_current_user_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/request_otp_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/verify_otp_usecase.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_bloc.dart';
import 'package:dink_mobile/features/auth/presentation/pages/phone_input_page.dart';
import 'package:dink_mobile/features/auth/presentation/pages/otp_verification_page.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class MockAuthRepo implements AuthRepository {
  @override
  Future<String> requestOtp({required String phoneNumber, required String purpose, String? countryCode}) async {
    return 'mock-session-1234';
  }

  @override
  Future<AuthTokens> verifyOtp({required String verificationSessionId, required String code}) async {
    return const AuthTokens(
      accessToken: 'test-token',
      refreshToken: 'test-refresh',
      tokenType: 'Bearer',
      expiresInSeconds: 900,
    );
  }

  @override
  Future<AuthTokens> refreshSession({required String refreshToken}) async {
    return const AuthTokens(
      accessToken: 'refreshed-token',
      refreshToken: 'refreshed-refresh',
      tokenType: 'Bearer',
      expiresInSeconds: 900,
    );
  }

  @override
  Future<AuthUser> getCurrentUser() async {
    return const AuthUser(
      id: 'uid',
      phoneNumber: '+251911223344',
      fullName: 'Test User',
      role: UserRole.renter,
      status: AccountStatus.active,
    );
  }

  @override
  Future<void> logout() async {}

  @override
  Future<bool> isAuthenticated() async => false;
}

void main() {
  late MockAuthRepo mockRepo;
  late AuthBloc authBloc;

  setUp(() {
    AppConfig.initialize(
      const AppConfig(
        environment: Environment.dev,
        appName: 'Dink (Dev)',
        apiBaseUrl: 'http://10.0.2.2:8080/api/v1',
        enableSslPinning: false,
        isarDbName: 'isar_dev',
      ),
    );
    mockRepo = MockAuthRepo();
    authBloc = AuthBloc(
      requestOtpUseCase: RequestOtpUseCase(mockRepo),
      verifyOtpUseCase: VerifyOtpUseCase(mockRepo),
      getCurrentUserUseCase: GetCurrentUserUseCase(mockRepo),
      authRepository: mockRepo,
    );
  });

  tearDown(() {
    authBloc.close();
  });

  testWidgets('PhoneInputPage renders and detects carrier dynamically', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: BlocProvider.value(
          value: authBloc,
          child: const PhoneInputPage(),
        ),
      ),
    );

    expect(find.text('Dink Machinery'), findsOneWidget);
    expect(find.text('Send Verification Code'), findsOneWidget);

    // Enter Ethiopian phone number starting with 9
    final phoneField = find.byType(TextField);
    await tester.enterText(phoneField, '911442910');
    await tester.pump();

    // Verify dynamic carrier badge appears (REG-02)
    expect(find.text('Carrier: Ethio Telecom'), findsOneWidget);

    // Clear and enter Safaricom number starting with 7
    await tester.enterText(phoneField, '711223344');
    await tester.pump();
    expect(find.text('Carrier: Safaricom ET'), findsOneWidget);
  });

  testWidgets('OtpVerificationPage displays 6-digit reticle and countdown', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: BlocProvider.value(
          value: authBloc,
          child: const OtpVerificationPage(phoneNumber: '+251911442910'),
        ),
      ),
    );

    expect(find.text('Verify Phone Number'), findsNWidgets(2));
    expect(find.textContaining('+251911442910'), findsOneWidget);
    expect(find.text('Verify & Continue'), findsOneWidget);
    expect(find.textContaining('Resend code in'), findsOneWidget);
  });
}
