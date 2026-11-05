import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_am.dart';
import 'app_localizations_en.dart';

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'l10n/app_localizations.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale)
    : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations? of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations);
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates =
      <LocalizationsDelegate<dynamic>>[
        delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
      ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[
    Locale('am'),
    Locale('en'),
  ];

  /// The title of the application
  ///
  /// In en, this message translates to:
  /// **'Dink Machinery'**
  String get appName;

  /// Heading for the phone entry screen
  ///
  /// In en, this message translates to:
  /// **'Enter your phone number'**
  String get phoneInputTitle;

  /// Subheading on the phone entry screen
  ///
  /// In en, this message translates to:
  /// **'We will send an SMS verification code to confirm your number.'**
  String get phoneInputSubtitle;

  /// Label for phone number text field
  ///
  /// In en, this message translates to:
  /// **'Phone Number'**
  String get phoneNumberLabel;

  /// Hint text for phone number field
  ///
  /// In en, this message translates to:
  /// **'91 234 5678'**
  String get phoneNumberHint;

  /// Action button to submit phone number
  ///
  /// In en, this message translates to:
  /// **'Send Verification Code'**
  String get sendOtpButton;

  /// Prefix for mobile network operator badge
  ///
  /// In en, this message translates to:
  /// **'Carrier'**
  String get carrierPrefix;

  /// Notice for legal agreement acceptance
  ///
  /// In en, this message translates to:
  /// **'By continuing, you agree to our Terms of Service & Privacy Policy.'**
  String get termsAgreementNotice;

  /// Heading for OTP verification screen
  ///
  /// In en, this message translates to:
  /// **'Verify Phone Number'**
  String get otpVerifyTitle;

  /// Instructions for OTP input with phone placeholder
  ///
  /// In en, this message translates to:
  /// **'Enter the 6-digit verification code sent to {phone}'**
  String otpVerifySubtitle(String phone);

  /// Button to submit OTP code
  ///
  /// In en, this message translates to:
  /// **'Verify & Continue'**
  String get verifyButton;

  /// Button to trigger another OTP code
  ///
  /// In en, this message translates to:
  /// **'Resend Code'**
  String get resendCodeButton;

  /// Countdown timer for resend cooldown
  ///
  /// In en, this message translates to:
  /// **'Resend code in {seconds}s'**
  String resendCountdown(int seconds);

  /// Validation message for malformed phone number
  ///
  /// In en, this message translates to:
  /// **'Please enter a valid phone number'**
  String get errorInvalidPhone;

  /// Error when rate limit is exceeded
  ///
  /// In en, this message translates to:
  /// **'Too many attempts. Please wait before retrying.'**
  String get errorRateLimit;

  /// Generic error message
  ///
  /// In en, this message translates to:
  /// **'An unexpected error occurred. Please try again.'**
  String get errorGeneric;

  /// Generic retry button text
  ///
  /// In en, this message translates to:
  /// **'Retry'**
  String get commonRetry;

  /// Generic cancel button text
  ///
  /// In en, this message translates to:
  /// **'Cancel'**
  String get commonCancel;

  /// App bar title for machinery registration
  ///
  /// In en, this message translates to:
  /// **'Register Machinery'**
  String get fleetRegistrationTitle;

  /// Label for category selector
  ///
  /// In en, this message translates to:
  /// **'Select Category'**
  String get selectCategory;

  /// Section header for machinery core details
  ///
  /// In en, this message translates to:
  /// **'Machinery Details'**
  String get machineryDetails;

  /// Label for machine make/brand
  ///
  /// In en, this message translates to:
  /// **'Make'**
  String get make;

  /// Label for machine model
  ///
  /// In en, this message translates to:
  /// **'Model'**
  String get model;

  /// Label for machine year of manufacture
  ///
  /// In en, this message translates to:
  /// **'Manufacturing Year'**
  String get manufacturingYear;

  /// Label for plate or registration identifier
  ///
  /// In en, this message translates to:
  /// **'Plate / Registration Number'**
  String get registrationNumber;

  /// Label for chassis or VIN
  ///
  /// In en, this message translates to:
  /// **'Chassis Number / VIN'**
  String get chassisVin;

  /// Label for baseline mechanical meter reading
  ///
  /// In en, this message translates to:
  /// **'Engine Hour Meter'**
  String get baselineMeterHours;

  /// Label for fuel type
  ///
  /// In en, this message translates to:
  /// **'Fuel Type'**
  String get fuelType;

  /// Label for daily rate
  ///
  /// In en, this message translates to:
  /// **'Daily Rate'**
  String get dailyRate;

  /// Toggle indicating whether equipment can travel without transporter
  ///
  /// In en, this message translates to:
  /// **'Self-Move Capable'**
  String get selfMoveCapable;

  /// Toggle indicating whether dedicated certified operator is included
  ///
  /// In en, this message translates to:
  /// **'Operator Included'**
  String get operatorIncluded;

  /// Section title for mandatory 6-angle walkaround photos
  ///
  /// In en, this message translates to:
  /// **'Walkaround 6-Angle Photos'**
  String get walkaroundPhotos;

  /// Front angle label
  ///
  /// In en, this message translates to:
  /// **'Front'**
  String get photoFront;

  /// Rear angle label
  ///
  /// In en, this message translates to:
  /// **'Rear'**
  String get photoRear;

  /// Left angle label
  ///
  /// In en, this message translates to:
  /// **'Left Side'**
  String get photoLeft;

  /// Right angle label
  ///
  /// In en, this message translates to:
  /// **'Right Side'**
  String get photoRight;

  /// Cabin angle label
  ///
  /// In en, this message translates to:
  /// **'Cabin & Controls'**
  String get photoCabin;

  /// Meter angle label
  ///
  /// In en, this message translates to:
  /// **'Hour Meter Reading'**
  String get photoMeter;

  /// Section title for dynamic technical specs
  ///
  /// In en, this message translates to:
  /// **'Technical Specifications'**
  String get technicalSpecifications;

  /// Submit button for machinery registration
  ///
  /// In en, this message translates to:
  /// **'Submit for Verification'**
  String get submitRegistration;

  /// Save draft button
  ///
  /// In en, this message translates to:
  /// **'Save Offline Draft'**
  String get saveDraft;

  /// Toast/Snackbar message when draft is saved
  ///
  /// In en, this message translates to:
  /// **'Draft saved locally to your device'**
  String get draftSavedMessage;

  /// Validation error when some angles are not captured
  ///
  /// In en, this message translates to:
  /// **'All 6 walkaround photos are required before submission'**
  String get errorIncompletePhotos;

  /// Validation error when required fields are empty
  ///
  /// In en, this message translates to:
  /// **'Please fill in all required fields'**
  String get errorMissingFields;

  /// Success message after asset creation
  ///
  /// In en, this message translates to:
  /// **'Machinery registered successfully! Under review.'**
  String get registrationSuccessMessage;

  /// Screen title for registering a standalone equipment attachment
  ///
  /// In en, this message translates to:
  /// **'Register Attachment'**
  String get registerAttachmentTitle;

  /// Category dropdown label for attachment type
  ///
  /// In en, this message translates to:
  /// **'Attachment Type'**
  String get attachmentTypeLabel;

  /// Brand/make of the attachment
  ///
  /// In en, this message translates to:
  /// **'Make / Manufacturer'**
  String get attachmentMakeLabel;

  /// Model descriptor of the attachment
  ///
  /// In en, this message translates to:
  /// **'Model'**
  String get attachmentModelLabel;

  /// Chisel diameter or bucket size
  ///
  /// In en, this message translates to:
  /// **'Tool Sizing Specification'**
  String get attachmentSubTypeLabel;

  /// Stamped manufacturer serial PIN number
  ///
  /// In en, this message translates to:
  /// **'Serial PIN Number'**
  String get attachmentSerialLabel;

  /// Weight of implement in kilograms
  ///
  /// In en, this message translates to:
  /// **'Operating Tool Weight (kg)'**
  String get operatingWeightKgLabel;

  /// Minimum carrier class in metric tons
  ///
  /// In en, this message translates to:
  /// **'Min Host Carrier Weight (tons)'**
  String get minCarrierWeightLabel;

  /// Maximum carrier class in metric tons
  ///
  /// In en, this message translates to:
  /// **'Max Host Carrier Weight (tons)'**
  String get maxCarrierWeightLabel;

  /// Required hydraulic relief pressure
  ///
  /// In en, this message translates to:
  /// **'Operating Pressure (bar)'**
  String get requiredPressureLabel;

  /// Minimum required auxiliary oil flow
  ///
  /// In en, this message translates to:
  /// **'Min Hydraulic Flow (L/min)'**
  String get requiredFlowMinLabel;

  /// Maximum required auxiliary oil flow
  ///
  /// In en, this message translates to:
  /// **'Max Hydraulic Flow (L/min)'**
  String get requiredFlowMaxLabel;

  /// Coupler pin diameter in millimeters
  ///
  /// In en, this message translates to:
  /// **'Mounting Pin Diameter (mm)'**
  String get pinDiameterLabel;

  /// Coupler center-to-center or ear spread in millimeters
  ///
  /// In en, this message translates to:
  /// **'Pin Spread / Ear Width (mm)'**
  String get pinSpreadLabel;

  /// Daily rental tariff when hired independently
  ///
  /// In en, this message translates to:
  /// **'Standalone Daily Rate (ETB)'**
  String get standaloneDailyRateLabel;

  /// Discounted tariff when hired with host fleet
  ///
  /// In en, this message translates to:
  /// **'Bundled Daily Rate (ETB)'**
  String get bundledDailyRateLabel;

  /// Documentary deposit amount
  ///
  /// In en, this message translates to:
  /// **'Security Deposit (ETB)'**
  String get securityDepositLabel;

  /// Section header for 4 mandatory attachment photos
  ///
  /// In en, this message translates to:
  /// **'Attachment Conditioning Photos'**
  String get attachmentPhotosSection;

  /// Full profile photo of attachment
  ///
  /// In en, this message translates to:
  /// **'Overall Tool Profile'**
  String get photoAngleOverall;

  /// Coupler bracket photo
  ///
  /// In en, this message translates to:
  /// **'Mounting Bracket & Hitch'**
  String get photoAngleBracket;

  /// Close-up photo of pin and bushing play
  ///
  /// In en, this message translates to:
  /// **'Pins & Bushings Wear'**
  String get photoAnglePins;

  /// Tool tip or tooth cutting edge photo
  ///
  /// In en, this message translates to:
  /// **'Working Tip / Chisel Point'**
  String get photoAngleToolTip;

  /// Action to bundle attachment with a machinery asset
  ///
  /// In en, this message translates to:
  /// **'Link to Host Fleet'**
  String get linkAttachmentToFleet;

  /// Title of owner attachment inventory
  ///
  /// In en, this message translates to:
  /// **'My Attachments'**
  String get attachmentInventoryTitle;

  /// Active attachment status tag
  ///
  /// In en, this message translates to:
  /// **'Available'**
  String get attachmentStatusActive;

  /// Rented attachment status tag
  ///
  /// In en, this message translates to:
  /// **'Rented'**
  String get attachmentStatusRented;

  /// Validation error for attachment photos
  ///
  /// In en, this message translates to:
  /// **'All 4 attachment conditioning photos are mandatory'**
  String get errorIncompleteAttachmentPhotos;

  /// Validation error when bundled rate is too high
  ///
  /// In en, this message translates to:
  /// **'Bundled rate cannot exceed standalone rate'**
  String get errorBundledRateExceedsStandalone;

  /// Success message after attachment registration
  ///
  /// In en, this message translates to:
  /// **'Attachment registered successfully!'**
  String get attachmentRegisteredSuccess;

  /// Section header for daily rate and deposits
  ///
  /// In en, this message translates to:
  /// **'Tariff & Security Deposits'**
  String get tariffAndDepositsSection;

  /// Success message when an attachment is linked to a host fleet
  ///
  /// In en, this message translates to:
  /// **'Attachment linked to fleet asset successfully!'**
  String get attachmentLinkedSuccess;

  /// Placeholder when user has no registered attachments
  ///
  /// In en, this message translates to:
  /// **'No attachments registered yet'**
  String get emptyAttachmentInventory;

  /// Badge indicating attachment is bundled with a host machine
  ///
  /// In en, this message translates to:
  /// **'Bundled with Host'**
  String get bundledWithHost;

  /// Photo status label when captured and validated
  ///
  /// In en, this message translates to:
  /// **'Verified Clean'**
  String get verifiedClean;

  /// Prompt on empty photo reticle tile
  ///
  /// In en, this message translates to:
  /// **'Tap to Capture'**
  String get tapToCapture;

  /// Count of captured conditioning photos
  ///
  /// In en, this message translates to:
  /// **'{count} of {total} Captured'**
  String photosCapturedCount(int count, int total);

  /// Daily rental time unit
  ///
  /// In en, this message translates to:
  /// **'day'**
  String get perDay;

  /// Placeholder for discovery search input
  ///
  /// In en, this message translates to:
  /// **'Search by make, model, or city...'**
  String get discoverySearchHint;

  /// Title for multi-facet filter sheet
  ///
  /// In en, this message translates to:
  /// **'Filter Machinery'**
  String get filterTitle;

  /// Distance indicator for fleet card
  ///
  /// In en, this message translates to:
  /// **'{km} km away'**
  String distanceAway(String km);

  /// Badge for immediate availability
  ///
  /// In en, this message translates to:
  /// **'Available Now'**
  String get availableNow;

  /// Badge for fleet asset with active checkout hold
  ///
  /// In en, this message translates to:
  /// **'Held (In Checkout)'**
  String get heldPendingWire;

  /// Button to open comparison matrix
  ///
  /// In en, this message translates to:
  /// **'Compare Selected ({count})'**
  String compareSelected(int count);

  /// Title for side-by-side comparison page
  ///
  /// In en, this message translates to:
  /// **'Machinery Comparison'**
  String get comparisonMatrixTitle;

  /// Toggle to filter identical spec rows in comparison
  ///
  /// In en, this message translates to:
  /// **'Differences Only'**
  String get differencesOnly;

  /// Badge for machine rental including operator
  ///
  /// In en, this message translates to:
  /// **'Crew Included'**
  String get crewIncluded;

  /// Badge for machine requiring transport trailer
  ///
  /// In en, this message translates to:
  /// **'Lowboy Required'**
  String get lowboyRequired;

  /// Button to apply selected filters
  ///
  /// In en, this message translates to:
  /// **'Apply Filters'**
  String get applyFilters;

  /// Button to reset filters
  ///
  /// In en, this message translates to:
  /// **'Reset'**
  String get resetFilters;

  /// Empty state message when no fleets match
  ///
  /// In en, this message translates to:
  /// **'No machinery found matching your search criteria.'**
  String get noFleetsFound;

  /// Button to clear comparison selection
  ///
  /// In en, this message translates to:
  /// **'Clear'**
  String get clearCompare;

  /// Button to submit rental booking request
  ///
  /// In en, this message translates to:
  /// **'Request to Book'**
  String get requestToBook;

  /// Title for booking checkout breakdown
  ///
  /// In en, this message translates to:
  /// **'Booking Details'**
  String get bookingDetails;

  /// Label for rental duration
  ///
  /// In en, this message translates to:
  /// **'Rental Duration'**
  String get rentalDuration;

  /// Fee for primary machinery asset hire
  ///
  /// In en, this message translates to:
  /// **'Base Hire Fee'**
  String get baseHireFee;

  /// Fee for attached bundled implements
  ///
  /// In en, this message translates to:
  /// **'Attachments Fee'**
  String get attachmentsFee;

  /// Total combined rental cost
  ///
  /// In en, this message translates to:
  /// **'Gross Rental Amount'**
  String get grossRentalAmount;

  /// Required 50 percent direct bank transfer advance
  ///
  /// In en, this message translates to:
  /// **'Advance Prepayment (50%)'**
  String get advancePrepayment;

  /// Balance due upon rental return
  ///
  /// In en, this message translates to:
  /// **'Estimated Remaining Balance'**
  String get remainingBalance;

  /// Physical destination of rented machinery
  ///
  /// In en, this message translates to:
  /// **'Site Location'**
  String get siteLocation;

  /// Radial countdown for 60-minute Valkey lease
  ///
  /// In en, this message translates to:
  /// **'Payment Hold Timer'**
  String get paymentHoldTimer;

  /// Button to upload wire payment slip
  ///
  /// In en, this message translates to:
  /// **'Upload Bank Transfer Slip'**
  String get uploadBankSlip;

  /// Fleet owner button to confirm receipt of bank wire
  ///
  /// In en, this message translates to:
  /// **'Confirm Payment Received'**
  String get confirmPaymentReceived;

  /// Fleet owner button to decline unverified bank wire
  ///
  /// In en, this message translates to:
  /// **'Decline Payment'**
  String get declinePayment;

  /// Contractor button to dispute declined payment
  ///
  /// In en, this message translates to:
  /// **'Dispute Decline'**
  String get disputeDecline;

  /// Action to cancel booking
  ///
  /// In en, this message translates to:
  /// **'Cancel Booking'**
  String get cancelBooking;

  /// 6-digit challenge PIN for site handover
  ///
  /// In en, this message translates to:
  /// **'Handover PIN'**
  String get handoverPIN;

  /// Baseline cabin engine hour meter reading
  ///
  /// In en, this message translates to:
  /// **'Check-In Meter Hours'**
  String get checkInMeter;

  /// Return check-out cabin engine hour meter reading
  ///
  /// In en, this message translates to:
  /// **'Check-Out Meter Hours'**
  String get checkOutMeter;

  /// Status badge for pending payment
  ///
  /// In en, this message translates to:
  /// **'Pending Payment'**
  String get bookingStatusPendingPayment;

  /// Status badge for uploaded slip
  ///
  /// In en, this message translates to:
  /// **'Payment Evidence Uploaded'**
  String get bookingStatusEvidenceUploaded;

  /// Status badge for confirmed booking
  ///
  /// In en, this message translates to:
  /// **'Confirmed'**
  String get bookingStatusConfirmed;

  /// Status badge for active in-field operation
  ///
  /// In en, this message translates to:
  /// **'Active Rental'**
  String get bookingStatusActive;

  /// Status badge for successfully concluded rental
  ///
  /// In en, this message translates to:
  /// **'Completed'**
  String get bookingStatusCompleted;

  /// Status badge for cancelled booking
  ///
  /// In en, this message translates to:
  /// **'Cancelled'**
  String get bookingStatusCancelled;

  /// Unit label for duration in days
  ///
  /// In en, this message translates to:
  /// **'days'**
  String get daysUnit;

  /// Title for active rental telemetry HUD
  ///
  /// In en, this message translates to:
  /// **'Active Rental Telemetry'**
  String get activeRentalTitle;

  /// Title for physical handover inspection
  ///
  /// In en, this message translates to:
  /// **'Pre-Dispatch Handover Inspection'**
  String get handoverInspectionTitle;

  /// Label for hour meter reading
  ///
  /// In en, this message translates to:
  /// **'Certified Hour Meter'**
  String get meterReadingLabel;

  /// Checkbox for fluid levels
  ///
  /// In en, this message translates to:
  /// **'Fluid & Lubricant Levels Verified'**
  String get fluidLevelsCheck;

  /// Checkbox for safety glass
  ///
  /// In en, this message translates to:
  /// **'Cabin Glass & Mirrors Intact'**
  String get safetyGlassCheck;

  /// Label for inspection reticle photos
  ///
  /// In en, this message translates to:
  /// **'Handover Reticle Photos'**
  String get inspectionPhotosLabel;

  /// Button to generate owner challenge PIN
  ///
  /// In en, this message translates to:
  /// **'Generate Handover PIN'**
  String get generatePINButton;

  /// Button to verify handover and activate rental
  ///
  /// In en, this message translates to:
  /// **'Verify Handover & Activate'**
  String get verifyHandoverButton;

  /// Prompt for entering 6-digit challenge PIN
  ///
  /// In en, this message translates to:
  /// **'Enter 6-Digit Owner PIN'**
  String get pinChallengePrompt;

  /// Label for touch signature canvas
  ///
  /// In en, this message translates to:
  /// **'Contractor Touch Signature'**
  String get touchSignatureLabel;

  /// Title for daily shift telemetry card
  ///
  /// In en, this message translates to:
  /// **'Daily Shift Telemetry'**
  String get shiftProgressTitle;

  /// Label for 8-hour shift quota
  ///
  /// In en, this message translates to:
  /// **'8.0 Hr Standard Shift Quota'**
  String get shiftQuotaLabel;

  /// Label for hours accumulated today
  ///
  /// In en, this message translates to:
  /// **'Accumulated Today'**
  String get accumulatedTodayLabel;

  /// Label for excess billable hours
  ///
  /// In en, this message translates to:
  /// **'Excess Hours (Tier 2)'**
  String get excessHoursLabel;

  /// Notice when shift timer is paused due to breakdown
  ///
  /// In en, this message translates to:
  /// **'Shift Timer Paused (Breakdown Active)'**
  String get shiftTimerPausedNotice;

  /// Quick action button to extend rental
  ///
  /// In en, this message translates to:
  /// **'Extend Rental'**
  String get quickActionExtend;

  /// Quick action button to report breakdown
  ///
  /// In en, this message translates to:
  /// **'Report Breakdown'**
  String get quickActionBreakdown;

  /// Quick action button to initiate return
  ///
  /// In en, this message translates to:
  /// **'Initiate Return'**
  String get quickActionReturn;

  /// Title for breakdown reporting modal
  ///
  /// In en, this message translates to:
  /// **'Report Mechanical Breakdown'**
  String get breakdownModalTitle;

  /// Label for fault category picker
  ///
  /// In en, this message translates to:
  /// **'Fault Category'**
  String get faultCategoryLabel;

  /// Label for fault description field
  ///
  /// In en, this message translates to:
  /// **'Fault Description & Symptoms'**
  String get faultDescriptionLabel;

  /// Label for diagnostic trouble codes
  ///
  /// In en, this message translates to:
  /// **'Diagnostic Trouble Codes (DTC)'**
  String get dtcLabel;

  /// Button to submit breakdown
  ///
  /// In en, this message translates to:
  /// **'Submit Breakdown & Pause Timer'**
  String get submitBreakdownButton;

  /// Label for owner 4-hour SLA deadline
  ///
  /// In en, this message translates to:
  /// **'Owner SLA Deadline'**
  String get slaDeadlineLabel;

  /// Title for extension request bottom sheet
  ///
  /// In en, this message translates to:
  /// **'Request Rental Extension'**
  String get extensionSheetTitle;

  /// Label for new rental end date
  ///
  /// In en, this message translates to:
  /// **'New End Date'**
  String get newEndDateLabel;

  /// Label for extension quote and prepayment
  ///
  /// In en, this message translates to:
  /// **'Incremental Quote & 50% Prepay'**
  String get extensionQuoteLabel;

  /// Button to confirm extension request
  ///
  /// In en, this message translates to:
  /// **'Confirm Extension Request'**
  String get submitExtensionButton;

  /// Title for breakdown resolution modal
  ///
  /// In en, this message translates to:
  /// **'Certify Repair & Resume Shift'**
  String get resolveBreakdownTitle;

  /// Label for mechanic name
  ///
  /// In en, this message translates to:
  /// **'Lead Mechanic Name'**
  String get mechanicNameLabel;

  /// Label for mechanic phone
  ///
  /// In en, this message translates to:
  /// **'Mechanic Contact Phone'**
  String get mechanicPhoneLabel;

  /// Label for repair resolution notes
  ///
  /// In en, this message translates to:
  /// **'Resolution Notes & Replaced Parts'**
  String get resolutionNotesLabel;

  /// Button to confirm resolution and resume shift
  ///
  /// In en, this message translates to:
  /// **'Resume Shift Telemetry'**
  String get confirmResolutionButton;

  /// Title for off-hire return inspection
  ///
  /// In en, this message translates to:
  /// **'Return Inspection & Meter Reading'**
  String get returnInspectionTitle;

  /// Label for checkout engine hour meter reading
  ///
  /// In en, this message translates to:
  /// **'Return Check-Out Meter Reading'**
  String get checkoutMeterLabel;

  /// Notice for valid non-regressive meter floor
  ///
  /// In en, this message translates to:
  /// **'Monotonic verification passed'**
  String get monotonicVerificationPassed;

  /// Label for 6 mandatory photo angles
  ///
  /// In en, this message translates to:
  /// **'Mandatory 6-Angle Reticle Inspection'**
  String get mandatoryAnglesLabel;

  /// Button to submit return inspection
  ///
  /// In en, this message translates to:
  /// **'Submit Return Inspection'**
  String get submitReturnInspection;

  /// Title for final balance settlement statement
  ///
  /// In en, this message translates to:
  /// **'Final Settlement Statement'**
  String get settlementStatementTitle;

  /// Label for final net balance due
  ///
  /// In en, this message translates to:
  /// **'Net Balance Settlement Due'**
  String get netBalanceDueLabel;

  /// Title for post-rental damage docket
  ///
  /// In en, this message translates to:
  /// **'Damage Claim Docket'**
  String get damageDocketTitle;

  /// Button to accept damage claim
  ///
  /// In en, this message translates to:
  /// **'Accept Liability'**
  String get acceptDamageLiability;

  /// Button to dispute damage claim
  ///
  /// In en, this message translates to:
  /// **'Dispute Claim & Escalate'**
  String get disputeDamageClaim;

  /// Title for double-blind reputation review
  ///
  /// In en, this message translates to:
  /// **'Mutual Double-Blind Review'**
  String get doubleBlindReviewTitle;

  /// Button to submit confidential review
  ///
  /// In en, this message translates to:
  /// **'Submit Confidential Review'**
  String get submitConfidentialReview;
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(lookupAppLocalizations(locale));
  }

  @override
  bool isSupported(Locale locale) =>
      <String>['am', 'en'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'am':
      return AppLocalizationsAm();
    case 'en':
      return AppLocalizationsEn();
  }

  throw FlutterError(
    'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
    'an issue with the localizations generation tool. Please file an issue '
    'on GitHub with a reproducible sample app and the gen-l10n configuration '
    'that was used.',
  );
}
