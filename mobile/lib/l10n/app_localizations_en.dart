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

  @override
  String get registerAttachmentTitle => 'Register Attachment';

  @override
  String get attachmentTypeLabel => 'Attachment Type';

  @override
  String get attachmentMakeLabel => 'Make / Manufacturer';

  @override
  String get attachmentModelLabel => 'Model';

  @override
  String get attachmentSubTypeLabel => 'Tool Sizing Specification';

  @override
  String get attachmentSerialLabel => 'Serial PIN Number';

  @override
  String get operatingWeightKgLabel => 'Operating Tool Weight (kg)';

  @override
  String get minCarrierWeightLabel => 'Min Host Carrier Weight (tons)';

  @override
  String get maxCarrierWeightLabel => 'Max Host Carrier Weight (tons)';

  @override
  String get requiredPressureLabel => 'Operating Pressure (bar)';

  @override
  String get requiredFlowMinLabel => 'Min Hydraulic Flow (L/min)';

  @override
  String get requiredFlowMaxLabel => 'Max Hydraulic Flow (L/min)';

  @override
  String get pinDiameterLabel => 'Mounting Pin Diameter (mm)';

  @override
  String get pinSpreadLabel => 'Pin Spread / Ear Width (mm)';

  @override
  String get standaloneDailyRateLabel => 'Standalone Daily Rate (ETB)';

  @override
  String get bundledDailyRateLabel => 'Bundled Daily Rate (ETB)';

  @override
  String get securityDepositLabel => 'Security Deposit (ETB)';

  @override
  String get attachmentPhotosSection => 'Attachment Conditioning Photos';

  @override
  String get photoAngleOverall => 'Overall Tool Profile';

  @override
  String get photoAngleBracket => 'Mounting Bracket & Hitch';

  @override
  String get photoAnglePins => 'Pins & Bushings Wear';

  @override
  String get photoAngleToolTip => 'Working Tip / Chisel Point';

  @override
  String get linkAttachmentToFleet => 'Link to Host Fleet';

  @override
  String get attachmentInventoryTitle => 'My Attachments';

  @override
  String get attachmentStatusActive => 'Available';

  @override
  String get attachmentStatusRented => 'Rented';

  @override
  String get errorIncompleteAttachmentPhotos =>
      'All 4 attachment conditioning photos are mandatory';

  @override
  String get errorBundledRateExceedsStandalone =>
      'Bundled rate cannot exceed standalone rate';

  @override
  String get attachmentRegisteredSuccess =>
      'Attachment registered successfully!';

  @override
  String get tariffAndDepositsSection => 'Tariff & Security Deposits';

  @override
  String get attachmentLinkedSuccess =>
      'Attachment linked to fleet asset successfully!';

  @override
  String get emptyAttachmentInventory => 'No attachments registered yet';

  @override
  String get bundledWithHost => 'Bundled with Host';

  @override
  String get verifiedClean => 'Verified Clean';

  @override
  String get tapToCapture => 'Tap to Capture';

  @override
  String photosCapturedCount(int count, int total) {
    return '$count of $total Captured';
  }

  @override
  String get perDay => 'day';
}
