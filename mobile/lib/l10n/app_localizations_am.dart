// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Amharic (`am`).
class AppLocalizationsAm extends AppLocalizations {
  AppLocalizationsAm([String locale = 'am']) : super(locale);

  @override
  String get appName => 'ድንቅ ማሽነሪ';

  @override
  String get phoneInputTitle => 'ስልክ ቁጥርዎን ያስገቡ';

  @override
  String get phoneInputSubtitle =>
      'ቁጥርዎን ለማረጋገጥ ባለ 6 አሃዝ የማረጋገጫ ኮድ በኤስኤምኤስ (SMS) እንልካለን።';

  @override
  String get phoneNumberLabel => 'ስልክ ቁጥር';

  @override
  String get phoneNumberHint => '91 234 5678';

  @override
  String get sendOtpButton => 'የማረጋገጫ ኮድ ላክ';

  @override
  String get carrierPrefix => 'አገልግሎት አቅራቢ';

  @override
  String get termsAgreementNotice =>
      'በመቀጠልዎ በአገልግሎት ውል እና የግላዊነት መመሪያችን ተስማምተዋል።';

  @override
  String get otpVerifyTitle => 'ስልክ ቁጥርዎን ያረጋግጡ';

  @override
  String otpVerifySubtitle(String phone) {
    return 'የተላከውን ባለ 6 አሃዝ የማረጋገጫ ኮድ ያስገቡ ወደ $phone';
  }

  @override
  String get verifyButton => 'አረጋግጥና ቀጥል';

  @override
  String get resendCodeButton => 'ኮዱን እንደገና ላክ';

  @override
  String resendCountdown(int seconds) {
    return 'ኮዱን በድጋሚ ለመላክ $seconds ሰከንድ';
  }

  @override
  String get errorInvalidPhone => 'እባክዎን ትክክለኛ ስልክ ቁጥር ያስገቡ';

  @override
  String get errorRateLimit =>
      'የተፈቀደው የጥያቄ ገደብ አልፏል። እባክዎን ትንሽ ቆይተው እንደገና ይሞክሩ።';

  @override
  String get errorGeneric => 'ያልተጠበቀ ስህተት አጋጥሟል። እባክዎን እንደገና ይሞክሩ።';

  @override
  String get commonRetry => 'እንደገና ሞክር';

  @override
  String get commonCancel => 'ሰርዝ';
}
