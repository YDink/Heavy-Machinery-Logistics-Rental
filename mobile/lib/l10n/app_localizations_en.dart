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
  String get dailyRate => 'Daily Rate';

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

  @override
  String get discoverySearchHint => 'Search by make, model, or city...';

  @override
  String get filterTitle => 'Filter Machinery';

  @override
  String distanceAway(String km) {
    return '$km km away';
  }

  @override
  String get availableNow => 'Available Now';

  @override
  String get heldPendingWire => 'Held (In Checkout)';

  @override
  String compareSelected(int count) {
    return 'Compare Selected ($count)';
  }

  @override
  String get comparisonMatrixTitle => 'Machinery Comparison';

  @override
  String get differencesOnly => 'Differences Only';

  @override
  String get crewIncluded => 'Crew Included';

  @override
  String get lowboyRequired => 'Lowboy Required';

  @override
  String get applyFilters => 'Apply Filters';

  @override
  String get resetFilters => 'Reset';

  @override
  String get noFleetsFound =>
      'No machinery found matching your search criteria.';

  @override
  String get clearCompare => 'Clear';

  @override
  String get requestToBook => 'Request to Book';

  @override
  String get bookingDetails => 'Booking Details';

  @override
  String get rentalDuration => 'Rental Duration';

  @override
  String get baseHireFee => 'Base Hire Fee';

  @override
  String get attachmentsFee => 'Attachments Fee';

  @override
  String get grossRentalAmount => 'Gross Rental Amount';

  @override
  String get advancePrepayment => 'Advance Prepayment (50%)';

  @override
  String get remainingBalance => 'Estimated Remaining Balance';

  @override
  String get siteLocation => 'Site Location';

  @override
  String get paymentHoldTimer => 'Payment Hold Timer';

  @override
  String get uploadBankSlip => 'Upload Bank Transfer Slip';

  @override
  String get confirmPaymentReceived => 'Confirm Payment Received';

  @override
  String get declinePayment => 'Decline Payment';

  @override
  String get disputeDecline => 'Dispute Decline';

  @override
  String get cancelBooking => 'Cancel Booking';

  @override
  String get handoverPIN => 'Handover PIN';

  @override
  String get checkInMeter => 'Check-In Meter Hours';

  @override
  String get checkOutMeter => 'Check-Out Meter Hours';

  @override
  String get bookingStatusPendingPayment => 'Pending Payment';

  @override
  String get bookingStatusEvidenceUploaded => 'Payment Evidence Uploaded';

  @override
  String get bookingStatusConfirmed => 'Confirmed';

  @override
  String get bookingStatusActive => 'Active Rental';

  @override
  String get bookingStatusCompleted => 'Completed';

  @override
  String get bookingStatusCancelled => 'Cancelled';

  @override
  String get daysUnit => 'days';

  @override
  String get activeRentalTitle => 'Active Rental Telemetry';

  @override
  String get handoverInspectionTitle => 'Pre-Dispatch Handover Inspection';

  @override
  String get meterReadingLabel => 'Certified Hour Meter';

  @override
  String get fluidLevelsCheck => 'Fluid & Lubricant Levels Verified';

  @override
  String get safetyGlassCheck => 'Cabin Glass & Mirrors Intact';

  @override
  String get inspectionPhotosLabel => 'Handover Reticle Photos';

  @override
  String get generatePINButton => 'Generate Handover PIN';

  @override
  String get verifyHandoverButton => 'Verify Handover & Activate';

  @override
  String get pinChallengePrompt => 'Enter 6-Digit Owner PIN';

  @override
  String get touchSignatureLabel => 'Contractor Touch Signature';

  @override
  String get shiftProgressTitle => 'Daily Shift Telemetry';

  @override
  String get shiftQuotaLabel => '8.0 Hr Standard Shift Quota';

  @override
  String get accumulatedTodayLabel => 'Accumulated Today';

  @override
  String get excessHoursLabel => 'Excess Hours (Tier 2)';

  @override
  String get shiftTimerPausedNotice => 'Shift Timer Paused (Breakdown Active)';

  @override
  String get quickActionExtend => 'Extend Rental';

  @override
  String get quickActionBreakdown => 'Report Breakdown';

  @override
  String get quickActionReturn => 'Initiate Return';

  @override
  String get breakdownModalTitle => 'Report Mechanical Breakdown';

  @override
  String get faultCategoryLabel => 'Fault Category';

  @override
  String get faultDescriptionLabel => 'Fault Description & Symptoms';

  @override
  String get dtcLabel => 'Diagnostic Trouble Codes (DTC)';

  @override
  String get submitBreakdownButton => 'Submit Breakdown & Pause Timer';

  @override
  String get slaDeadlineLabel => 'Owner SLA Deadline';

  @override
  String get extensionSheetTitle => 'Request Rental Extension';

  @override
  String get newEndDateLabel => 'New End Date';

  @override
  String get extensionQuoteLabel => 'Incremental Quote & 50% Prepay';

  @override
  String get submitExtensionButton => 'Confirm Extension Request';

  @override
  String get resolveBreakdownTitle => 'Certify Repair & Resume Shift';

  @override
  String get mechanicNameLabel => 'Lead Mechanic Name';

  @override
  String get mechanicPhoneLabel => 'Mechanic Contact Phone';

  @override
  String get resolutionNotesLabel => 'Resolution Notes & Replaced Parts';

  @override
  String get confirmResolutionButton => 'Resume Shift Telemetry';

  @override
  String get returnInspectionTitle => 'Return Inspection & Meter Reading';

  @override
  String get checkoutMeterLabel => 'Return Check-Out Meter Reading';

  @override
  String get monotonicVerificationPassed => 'Monotonic verification passed';

  @override
  String get mandatoryAnglesLabel => 'Mandatory 6-Angle Reticle Inspection';

  @override
  String get submitReturnInspection => 'Submit Return Inspection';

  @override
  String get settlementStatementTitle => 'Final Settlement Statement';

  @override
  String get netBalanceDueLabel => 'Net Balance Settlement Due';

  @override
  String get damageDocketTitle => 'Damage Claim Docket';

  @override
  String get acceptDamageLiability => 'Accept Liability';

  @override
  String get disputeDamageClaim => 'Dispute Claim & Escalate';

  @override
  String get doubleBlindReviewTitle => 'Mutual Double-Blind Review';

  @override
  String get submitConfidentialReview => 'Submit Confidential Review';
}
