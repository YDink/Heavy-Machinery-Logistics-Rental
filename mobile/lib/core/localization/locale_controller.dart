import 'package:flutter/material.dart';

/// Manages dynamic runtime locale switching conforming to ADR-013 & ADR-019.
class LocaleController extends ValueNotifier<Locale> {
  LocaleController([super.initialLocale = const Locale('en')]);

  void setLocale(Locale newLocale) {
    if (value != newLocale) {
      value = newLocale;
    }
  }

  void toggleLocale() {
    value = value.languageCode == 'en' ? const Locale('am') : const Locale('en');
  }

  bool get isAmharic => value.languageCode == 'am';
}

/// Typographical line-height metrics conforming to SRS §05.16 §2.3.
class TypographyConfig {
  static const double amharicLineHeight = 1.45;
  static const double latinLineHeight = 1.20;

  static double lineHeightFor(Locale locale) {
    return locale.languageCode == 'am' ? amharicLineHeight : latinLineHeight;
  }
}
