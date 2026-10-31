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

  @override
  String get fleetRegistrationTitle => 'ማሽነሪ መመዝገቢያ';

  @override
  String get selectCategory => 'የማሽነሪ ምድብ ይምረጡ';

  @override
  String get machineryDetails => 'የማሽነሪው ዝርዝር መረጃ';

  @override
  String get make => 'አምራች ድርጅት';

  @override
  String get model => 'ሞዴል';

  @override
  String get manufacturingYear => 'የተመረተበት ዓመት';

  @override
  String get registrationNumber => 'የሰሌዳ / የታርጋ ቁጥር';

  @override
  String get chassisVin => 'የሻሲ ቁጥር (VIN)';

  @override
  String get baselineMeterHours => 'የሞተር ሰዓት ቆጣሪ';

  @override
  String get fuelType => 'የነዳጅ ዓይነት';

  @override
  String get dailyRate => 'የቀን ኪራይ ዋጋ (ብር)';

  @override
  String get selfMoveCapable => 'ራሱን ችሎ መጓዝ የሚችል';

  @override
  String get operatorIncluded => 'ከኦፕሬተር ጋር የቀረበ';

  @override
  String get walkaroundPhotos => 'የማሽነሪው 6 አቅጣጫ ፎቶዎች';

  @override
  String get photoFront => 'የፊት ገጽታ';

  @override
  String get photoRear => 'የኋላ ገጽታ';

  @override
  String get photoLeft => 'የግራ ገጽታ';

  @override
  String get photoRight => 'የቀኝ ገጽታ';

  @override
  String get photoCabin => 'ካቢን እና መቆጣጠሪያ';

  @override
  String get photoMeter => 'የሰዓት ቆጣሪ ገጽታ';

  @override
  String get technicalSpecifications => 'ቴክኒካዊ መረጃዎች';

  @override
  String get submitRegistration => 'ማሽነሪውን ለምርመራ አስገባ';

  @override
  String get saveDraft => 'ረቂቁን በስልኩ አስቀምጥ';

  @override
  String get draftSavedMessage => 'ረቂቁ በስልክዎ ላይ ተቀምጧል';

  @override
  String get errorIncompletePhotos => 'ሁሉንም 6 ፎቶዎች ማንሳት ግዴታ ነው';

  @override
  String get errorMissingFields => 'እባክዎ የተጠየቁትን መረጃዎች በሙሉ ይሙሉ';

  @override
  String get registrationSuccessMessage =>
      'ማሽነሪው በተሳካ ሁኔታ ተመዝግቧል! በግምገማ ላይ ነው።';
}
