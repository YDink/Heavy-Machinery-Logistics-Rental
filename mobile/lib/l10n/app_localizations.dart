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

  /// Label for daily rental price in ETB
  ///
  /// In en, this message translates to:
  /// **'Daily Rental Rate (ETB)'**
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
