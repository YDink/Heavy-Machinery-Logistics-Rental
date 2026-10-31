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

  @override
  String get fleetRegistrationTitle => 'Register Machinery';

  @override
  String get selectCategory => 'Select Category';

  @override
  String get machineryDetails => 'Machinery Details';

  @override
  String get make => 'Make';

  @override
  String get model => 'Model';

  @override
  String get manufacturingYear => 'Manufacturing Year';

  @override
  String get registrationNumber => 'Plate / Registration Number';

  @override
  String get chassisVin => 'Chassis Number / VIN';

  @override
  String get baselineMeterHours => 'Engine Hour Meter';

  @override
  String get fuelType => 'Fuel Type';

  @override
  String get dailyRate => 'Daily Rental Rate (ETB)';

  @override
  String get selfMoveCapable => 'Self-Move Capable';

  @override
  String get operatorIncluded => 'Operator Included';

  @override
  String get walkaroundPhotos => 'Walkaround 6-Angle Photos';

  @override
  String get photoFront => 'Front';

  @override
  String get photoRear => 'Rear';

  @override
  String get photoLeft => 'Left Side';

  @override
  String get photoRight => 'Right Side';

  @override
  String get photoCabin => 'Cabin & Controls';

  @override
  String get photoMeter => 'Hour Meter Reading';

  @override
  String get technicalSpecifications => 'Technical Specifications';

  @override
  String get submitRegistration => 'Submit for Verification';

  @override
  String get saveDraft => 'Save Offline Draft';

  @override
  String get draftSavedMessage => 'Draft saved locally to your device';

  @override
  String get errorIncompletePhotos =>
      'All 6 walkaround photos are required before submission';

  @override
  String get errorMissingFields => 'Please fill in all required fields';

  @override
  String get registrationSuccessMessage =>
      'Machinery registered successfully! Under review.';
}
