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
  String get dailyRate => 'የቀን ኪራይ';

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

  @override
  String get registerAttachmentTitle => 'አክሰሰሪ/አታችመንት መዝግብ';

  @override
  String get attachmentTypeLabel => 'የአታችመንት ዓይነት';

  @override
  String get attachmentMakeLabel => 'አምራች / ሜክ';

  @override
  String get attachmentModelLabel => 'ሞዴል';

  @override
  String get attachmentSubTypeLabel => 'መጠን / የዕቃው ዝርዝር';

  @override
  String get attachmentSerialLabel => 'የሴሪያል ቁጥር';

  @override
  String get operatingWeightKgLabel => 'የሥራ ክብደት (ኪ.ግ)';

  @override
  String get minCarrierWeightLabel => 'አነስተኛ የተሸካሚ ክብደት (ቶን)';

  @override
  String get maxCarrierWeightLabel => 'ከፍተኛ የተሸካሚ ክብደት (ቶን)';

  @override
  String get requiredPressureLabel => 'የሚያስፈልግ ግፊት (ባር)';

  @override
  String get requiredFlowMinLabel => 'አነስተኛ ዘይት ፍሰት (ሊ/ደቂቃ)';

  @override
  String get requiredFlowMaxLabel => 'ከፍተኛ ዘይት ፍሰት (ሊ/ደቂቃ)';

  @override
  String get pinDiameterLabel => 'የፒን ዲያሜትር (ሚ.ሜ)';

  @override
  String get pinSpreadLabel => 'የፒን ርቀት / ስፋት (ሚ.ሜ)';

  @override
  String get standaloneDailyRateLabel => 'የብቻ ዕለታዊ ኪራይ (ብር)';

  @override
  String get bundledDailyRateLabel => 'ከማሽነሪ ጋር ጥምር ዕለታዊ ኪራይ (ብር)';

  @override
  String get securityDepositLabel => 'የዋስትና ገንዘብ (ብር)';

  @override
  String get attachmentPhotosSection => 'የአታችመንት 4 አቅጣጫ ፎቶዎች';

  @override
  String get photoAngleOverall => 'አጠቃላይ የአታችመንት ገጽታ';

  @override
  String get photoAngleBracket => 'የመገጠሚያ ብሬኬት';

  @override
  String get photoAnglePins => 'የፒኖች እና ቡሺንጎች ሁኔታ';

  @override
  String get photoAngleToolTip => 'የመሥሪያው ጫፍ / ጥርስ';

  @override
  String get linkAttachmentToFleet => 'ከዋና ማሽነሪ ጋር አጣምር';

  @override
  String get attachmentInventoryTitle => 'የእኔ አታችመንቶች';

  @override
  String get attachmentStatusActive => 'ክፍት';

  @override
  String get attachmentStatusRented => 'የተከራየ';

  @override
  String get errorIncompleteAttachmentPhotos =>
      'ሁለቱም 4ቱ የአታችመንት ፎቶዎች የግድ ያስፈልጋሉ';

  @override
  String get errorBundledRateExceedsStandalone =>
      'የጥምር ኪራይ ዋጋ ከብቻ ዋጋው ሊበልጥ አይችልም';

  @override
  String get attachmentRegisteredSuccess => 'አታችመንት በተሳካ ሁኔታ ተመዝግቧል!';

  @override
  String get tariffAndDepositsSection => 'የኪራይ ተመን እና የዋስትና ተቀማጭ';

  @override
  String get attachmentLinkedSuccess => 'አባሪ መሳሪያው ከዋናው ማሽን ጋር በተሳካ ሁኔታ ተያይዟል!';

  @override
  String get emptyAttachmentInventory => 'እስካሁን የተመዘገበ አባሪ መሳሪያ የለም';

  @override
  String get bundledWithHost => 'ከዋናው ማሽን ጋር የተጣመረ';

  @override
  String get verifiedClean => 'የተረጋገጠ';

  @override
  String get tapToCapture => 'ፎቶ ለማንሳት ይንኩ';

  @override
  String photosCapturedCount(int count, int total) {
    return '$count ከ $total ተነስተዋል';
  }

  @override
  String get perDay => 'ቀን';

  @override
  String get discoverySearchHint => 'በሞዴል ወይም በከተማ ይፈልጉ...';

  @override
  String get filterTitle => 'የማሽን ማጣሪያ';

  @override
  String distanceAway(String km) {
    return '$km ኪ.ሜ ርቀት';
  }

  @override
  String get availableNow => 'አሁን ክፍት ነው';

  @override
  String get heldPendingWire => 'በይዞታ ላይ';

  @override
  String compareSelected(int count) {
    return 'የተመረጡትን አወዳድር ($count)';
  }

  @override
  String get comparisonMatrixTitle => 'የማሽነሪ ማነፃፀሪያ';

  @override
  String get differencesOnly => 'ልዩነቶችን ብቻ አሳይ';

  @override
  String get crewIncluded => 'ኦፕሬተር ያለው';

  @override
  String get lowboyRequired => 'ሎውቦይ የሚያስፈልገው';

  @override
  String get applyFilters => 'አጣራ';

  @override
  String get resetFilters => 'ዳግም አስጀምር';

  @override
  String get noFleetsFound => 'በመረጡት መስፈርት የተገኘ ማሽን የለም';

  @override
  String get clearCompare => 'አጽዳ';

  @override
  String get requestToBook => 'ይከራዩ';

  @override
  String get bookingDetails => 'የኪራይ ዝርዝር መረጃ';

  @override
  String get rentalDuration => 'የኪራይ ቆይታ';

  @override
  String get baseHireFee => 'መሰረታዊ የኪራይ ክፍያ';

  @override
  String get attachmentsFee => 'የተጨማሪ መገልገያዎች ክፍያ';

  @override
  String get grossRentalAmount => 'ጠቅላላ የኪራይ ክፍያ';

  @override
  String get advancePrepayment => 'የቅድመ ክፍያ (50%)';

  @override
  String get remainingBalance => 'የቀረው ክፍያ';

  @override
  String get siteLocation => 'የስራ ቦታ አድራሻ';

  @override
  String get paymentHoldTimer => 'የክፍያ ማቆያ ጊዜ';

  @override
  String get uploadBankSlip => 'የባንክ ዝውውር ደረሰኝ ያስገቡ';

  @override
  String get confirmPaymentReceived => 'ክፍያ መቀበሉን ያረጋግጡ';

  @override
  String get declinePayment => 'ክፍያውን አይቀበሉ';

  @override
  String get disputeDecline => 'የተቃውሞ አቤቱታ ያቅርቡ';

  @override
  String get cancelBooking => 'ኪራዩን ይሰርዙ';

  @override
  String get handoverPIN => 'የማስረከቢያ ፒን';

  @override
  String get checkInMeter => 'የመጀመሪያ የሰዓት ቆጣሪ';

  @override
  String get checkOutMeter => 'የመጨረሻ የሰዓት ቆጣሪ';

  @override
  String get bookingStatusPendingPayment => 'ክፍያ በመጠባበቅ ላይ';

  @override
  String get bookingStatusEvidenceUploaded => 'የክፍያ ደረሰኝ ተጭኗል';

  @override
  String get bookingStatusConfirmed => 'ተረጋግጧል';

  @override
  String get bookingStatusActive => 'በስራ ላይ';

  @override
  String get bookingStatusCompleted => 'ተጠናቋል';

  @override
  String get bookingStatusCancelled => 'ተሰርዟል';

  @override
  String get daysUnit => 'ቀናት';

  @override
  String get activeRentalTitle => 'የነቃ ኪራይ መረጃ';

  @override
  String get handoverInspectionTitle => 'ማሽን ከማስረከብ በፊት የሚደረግ ፍተሻ';

  @override
  String get meterReadingLabel => 'የማሽን የስራ ሰዓት ቆጣሪ';

  @override
  String get fluidLevelsCheck => 'የዘይት እና ፈሳሽ መጠን ተረጋግጧል';

  @override
  String get safetyGlassCheck => 'የካቢን መስታወቶች ያልተሰበሩ';

  @override
  String get inspectionPhotosLabel => 'የማስረከቢያ ፎቶዎች';

  @override
  String get generatePINButton => 'የማስረከቢያ ፒን አፍልቅ';

  @override
  String get verifyHandoverButton => 'ማስረከቢያውን አረጋግጥ እና አስጀምር';

  @override
  String get pinChallengePrompt => 'የባለቤቱን 6 አሃዝ ፒን ያስገቡ';

  @override
  String get touchSignatureLabel => 'የተከራይ ፊርማ';

  @override
  String get shiftProgressTitle => 'የዕለት ፈረቃ የስራ ሂደት';

  @override
  String get shiftQuotaLabel => '8.0 ሰዓት መደበኛ የፈረቃ ድልድል';

  @override
  String get accumulatedTodayLabel => 'ዛሬ የተሰራ ሰዓት';

  @override
  String get excessHoursLabel => 'ትርፍ ሰዓታት (ደረጃ 2)';

  @override
  String get shiftTimerPausedNotice => 'የፈረቃ ቆጣሪ ቆሟል (ብልሽት ተመዝግቧል)';

  @override
  String get quickActionExtend => 'ኪራይ ያራዝሙ';

  @override
  String get quickActionBreakdown => 'ብልሽት ያመልክቱ';

  @override
  String get quickActionReturn => 'ማስመለስ ይጀምሩ';

  @override
  String get breakdownModalTitle => 'የማሽን ብልሽት ማሳወቂያ';

  @override
  String get faultCategoryLabel => 'የብልሽቱ አይነት';

  @override
  String get faultDescriptionLabel => 'የብልሽቱ ዝርዝር መግለጫ';

  @override
  String get dtcLabel => 'የዲያግኖስቲክስ ኮዶች (DTC)';

  @override
  String get submitBreakdownButton => 'ብልሽቱን አስመዝግብ እና ቆጣሪውን አቁም';

  @override
  String get slaDeadlineLabel => 'የባለቤቱ የጥገና ምላሽ ገደብ';

  @override
  String get extensionSheetTitle => 'የኪራይ ማራዘሚያ ጥያቄ';

  @override
  String get newEndDateLabel => 'አዲሱ የማብቂያ ቀን';

  @override
  String get extensionQuoteLabel => 'የተጨማሪ ክፍያ እና 50% ቅድመ ክፍያ';

  @override
  String get submitExtensionButton => 'የማራዘሚያ ጥያቄውን አረጋግጥ';

  @override
  String get resolveBreakdownTitle => 'ጥገናውን አረጋግጥ እና ፈረቃውን አስጀምር';

  @override
  String get mechanicNameLabel => 'የመካኒኩ ስም';

  @override
  String get mechanicPhoneLabel => 'የመካኒኩ ስልክ ቁጥር';

  @override
  String get resolutionNotesLabel => 'የተከናወኑ ጥገናዎች እና የተተኩ እቃዎች';

  @override
  String get confirmResolutionButton => 'የፈረቃ ቆጣሪውን ቀጥል';

  @override
  String get returnInspectionTitle => 'የማሽነሪ መመለሻ ምርመራ እና የሰዓት ቆጣሪ ንባብ';

  @override
  String get checkoutMeterLabel => 'የማሽነሪ መመለሻ የሞተር ሰዓት ቆጣሪ';

  @override
  String get monotonicVerificationPassed => 'የሰዓት ቆጣሪ ቅደም ተከተል ተረጋግጧል';

  @override
  String get mandatoryAnglesLabel => 'ግዴታ የሆኑ 6 የማሽን አቅጣጫ ፎቶዎች';

  @override
  String get submitReturnInspection => 'የመመለሻ ምርመራውን መዝግብ';

  @override
  String get settlementStatementTitle => 'የመጨረሻ ሂሳብ ማጠቃለያ መግለጫ';

  @override
  String get netBalanceDueLabel => 'የሚከፈል የተጣራ ቀሪ ሂሳብ';

  @override
  String get damageDocketTitle => 'የማሽን ብልሽት ካሳ ጥያቄ ማህደር';

  @override
  String get acceptDamageLiability => 'የብልሽት ካሳ ተቀበል';

  @override
  String get disputeDamageClaim => 'የካሳ ጥያቄውን ተቃወም እና ለሽምግልና አስተላልፍ';

  @override
  String get doubleBlindReviewTitle => 'የሁለትዮሽ ሚስጥራዊ አስተያየት መስጫ';

  @override
  String get submitConfidentialReview => 'ሚስጥራዊ ግምገማ አስገባ';

  @override
  String get logisticsMobilizationTitle => 'የሎውቦይ ማጓጓዣ ዝግጅት';

  @override
  String get platformLowboyTab => 'የመድረክ ሎውቦይ (2.5% ክፍያ)';

  @override
  String get externalTruckTab => 'የራስ የግል መኪና';

  @override
  String get bookLowboyCarrier => 'ሎውቦይ አጓጓዥ ይዘዙ';

  @override
  String get declareExternalDriver => 'የውጭ ሹፌር መረጃ ያስገቡ';

  @override
  String get carrierDispatchTitle => 'የአጓጓዥ ጉዞ እና የመጫኛ ፍተሻ';

  @override
  String get winchingReticleHeader => 'ግዴታ የሆኑ 4 የመጫኛ ማረጋገጫ ፎቶዎች';

  @override
  String get confirmLoadStartTransit => 'መጫኑን አረጋግጥ እና ጉዞ ጀምር';

  @override
  String get reportRoadsideBreakdown => 'የመንገድ ላይ ብልሽት አመልክት';

  @override
  String get transshipmentSlaTitle => 'የአስቸኳይ ቅብብሎሽ ጊዜ (60 ደቂቃ)';

  @override
  String get acceptHotSwapTransfer => 'አስቸኳይ ቅብብሎሹን ተቀበል';

  @override
  String get executeReverseCascade => 'የፎርስ ማዥር ስረዛን አስፈጽም';

  @override
  String get parentBookingFrozenNotice => 'ዋናው ኪራይ ታግዷል፡ 0.0 ሰዓት ክፍያ';

  @override
  String get driverFullNameLabel => 'የሹፌር ሙሉ ስም';

  @override
  String get driverPhoneLabel => 'የሹፌር ስልክ ቁጥር';

  @override
  String get primeMoverPlateLabel => 'የመኪናው ታርጋ ቁጥር';

  @override
  String get carrierCompanyLabel => 'የአጓጓዥ ድርጅት ስም';
}
