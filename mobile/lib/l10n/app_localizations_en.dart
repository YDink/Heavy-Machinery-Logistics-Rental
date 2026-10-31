// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appName => 'Dink Machinery';

  @override
  String get phoneInputTitle => 'Enter your phone number';

  @override
  String get phoneInputSubtitle =>
      'We will send an SMS verification code to confirm your number.';

  @override
  String get phoneNumberLabel => 'Phone Number';

  @override
  String get phoneNumberHint => '91 234 5678';

  @override
  String get sendOtpButton => 'Send Verification Code';

  @override
  String get carrierPrefix => 'Carrier';

  @override
  String get termsAgreementNotice =>
      'By continuing, you agree to our Terms of Service & Privacy Policy.';

  @override
  String get otpVerifyTitle => 'Verify Phone Number';

  @override
  String otpVerifySubtitle(String phone) {
    return 'Enter the 6-digit verification code sent to $phone';
  }

  @override
  String get verifyButton => 'Verify & Continue';

  @override
  String get resendCodeButton => 'Resend Code';

  @override
  String resendCountdown(int seconds) {
    return 'Resend code in ${seconds}s';
  }

  @override
  String get errorInvalidPhone => 'Please enter a valid phone number';

  @override
  String get errorRateLimit =>
      'Too many attempts. Please wait before retrying.';

  @override
  String get errorGeneric => 'An unexpected error occurred. Please try again.';

  @override
  String get commonRetry => 'Retry';

  @override
  String get commonCancel => 'Cancel';
}
