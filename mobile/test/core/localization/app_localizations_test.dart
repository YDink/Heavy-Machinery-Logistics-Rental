import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/core/localization/locale_controller.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

void main() {
  test('LocaleController toggles between English and Amharic', () {
    final controller = LocaleController();
    expect(controller.value, const Locale('en'));
    expect(controller.isAmharic, false);

    controller.toggleLocale();
    expect(controller.value, const Locale('am'));
    expect(controller.isAmharic, true);

    controller.toggleLocale();
    expect(controller.value, const Locale('en'));
    expect(controller.isAmharic, false);
  });

  test('TypographyConfig enforces Amharic line height factor of 1.45', () {
    expect(TypographyConfig.lineHeightFor(const Locale('am')), 1.45);
    expect(TypographyConfig.lineHeightFor(const Locale('en')), 1.20);
  });

  testWidgets('AppLocalizations renders Amharic strings correctly', (tester) async {
    late AppLocalizations localizations;

    await tester.pumpWidget(
      MaterialApp(
        locale: const Locale('am'),
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: Builder(
          builder: (context) {
            localizations = AppLocalizations.of(context)!;
            return Text(localizations.appName);
          },
        ),
      ),
    );

    expect(find.text('ድንቅ ማሽነሪ'), findsOneWidget);
    expect(localizations.phoneInputTitle, 'ስልክ ቁጥርዎን ያስገቡ');
    expect(localizations.sendOtpButton, 'የማረጋገጫ ኮድ ላክ');
    expect(localizations.carrierPrefix, 'አገልግሎት አቅራቢ');
    expect(localizations.verifyButton, 'አረጋግጥና ቀጥል');
  });
}
