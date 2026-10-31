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
