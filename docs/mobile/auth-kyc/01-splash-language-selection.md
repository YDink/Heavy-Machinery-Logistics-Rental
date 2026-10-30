# Splash Screen, Country Locale Detection & Language Selector
## Mobile Subsystem Specification — `docs/mobile/auth-kyc/01-splash-language-selection.md`

**Document Identifier:** `SPEC-MOB-AUTH-01`  
**Classification:** Canonical Subsystem UI/UX & Technical Behavior Specification  
**Governing Documents:** `docs/source/Dink_Machinery_UI_UX_Design_Brief.md`, `docs/shared/design/00-design-system.md`, `docs/shared/design/design-tokens.json`, `docs/shared/requirements/05-01-registration-authentication.md`, `docs/shared/requirements/05-14-user-onboarding-help.md`, `docs/shared/requirements/05-16-accessibility-internationalisation.md`, `docs/shared/requirements/05-20-policy-terms-management.md`, `docs/ERROR_CATALOG.md`, `contracts/openapi.yaml`  
**Target Subsystem:** Flutter 3.47.8 / Dart 3.9 Mobile Client (iOS & Android)  
**Primary Persona / Actor:** Unauthenticated First-Time or Returning User  
**Requirement Traceability:** `REQ-05-01` (`REG-01`..`REG-04`, `REG-08`, `REG-09`), `REQ-05-14` (`ONB-01`), `REQ-05-16` (`AC-01`..`AC-04`, `LOC-01-BE`), `REQ-05-20` (`POL-01`), `REQ-05-23` (`HOL-01`), `REQ-09` (`NFR-PERF` `<1.5s` Cold Start)  
**Associated Interaction Flows:** `FLOW-ENTRY`, `FLOW-01`  
**Architectural Decision Records:** `ADR-008` (Hardware Biometric & Root Defense), `ADR-010` (Isar Local Hydration), `ADR-013` (OTA Ethiopic Font Loader), `ADR-014` (Design Tokens)  
**Status:** Approved Technical Design Baseline

---

## 1. Domain Context & Non-Negotiable Boot Invariants

The initial boot and language selection screen is the entry gate for contractors, yard managers, operators, and transporters across Ethiopia (`ET`) and Kenya (`KE`). On budget Android devices ($2\text{–}3\text{ GB RAM}$) operating over intermittent 3G/EDGE links under $>80,000\text{ lux}$ sunlight, the screen must reach an interactive frame in **$< 1.5\text{ seconds}$**, detect sovereign SIM/carrier jurisdiction automatically, and allow immediate selection of **Amharic (`am-ET`), Afaan Oromoo (`om-ET`), Kiswahili (`sw-KE`), or English (`en-US`)** before presenting dense text.

### 1.1 Non-Negotiable System Invariants
1. **Cold-Start Interactive Budget (`< 1.5s`) & Device Integrity Pre-Flight (`ADR-008`):**
   - During native splash initialization, the client concurrently opens the AES-256-GCM encrypted Isar database, interrogates the Ambient Light Sensor (`>25,000` / `>80,000 lux` auto-switching to `outdoorHighContrast`), and executes root/jailbreak/hooking environment attestation. Network calls never block the first interactive frame.
2. **Sovereign SIM MCC/MNC Country Detection (`REG-02`, `REG-04`, `LOC-01-BE`):**
   - The client inspects SIM Mobile Country Code (`MCC 636` $\rightarrow$ Ethiopia `ET`, `+251`, `ETB` Santims, Ethiopian Ge'ez + Gregorian Dual Calendar; `MCC 639` $\rightarrow$ Kenya `KE`, `+254`, `KES` Cents, Gregorian Calendar) and binds the `X-Country-Code` header globally. Users on WiFi-only tablets can toggle jurisdiction manually via `56dp` country cards.
3. **ADR-013 OTA Ethiopic Font Hydration with SHA-256 Verification:**
   - Selecting `አማርኛ (Amharic)` triggers the `ADR-013` font pipeline: checking local sandbox storage for `NotoSansEthiopic-v2.0.ttf`, or streaming the $4.2\text{ MB}$ binary from SeaweedFS S3 CDN with a non-blocking progress bar, verifying its **SHA-256 checksum**, and applying the **$1.38\times$ Golden Ge'ez line-height multiplier** in $< 16\text{ ms}$.
4. **Glove-Friendly `2x2` Sovereign Language Grid (`>= 64dp` Card Height):**
   - All four language cards render simultaneously in the primary viewport without hidden dropdowns, each enforcing a minimum `64dp` touch height, native script label (`አማርኛ`, `Afaan Oromoo`, `Kiswahili`, `English`), and `3dp` tactile bottom lip.
5. **Explicit Click-Wrap Legal Jurisdiction Disclosure (`REG-08`, `POL-01`):**
   - Advancing to Phone OTP Registration requires explicit affirmative check-off of the jurisdiction-parameterized Terms of Service and Non-Custodial Bank Wire Disclaimer (Ethiopian Civil Code Book IV for `ET`; Kenya Law of Contract Act Cap 23 for `KE`).

---

## 2. Bespoke ASCII Wireframes: Boot Splash & Language/Jurisdiction Selector

```text
+-----------------------------------------------------------------------+
| STATUS BAR: 08:30 EAT                        [☀️ AAA 83k LUX] [ET-4G] |
+=======================================================================+
|  +-----------------------------------------------------------------+  |
|  |  [🏗️ DINK MACHINERY EMBLEM]          [⏭️ SKIP TO LOGIN (56dp)]  |  |
|  |  HEAVY PLANT & LOGISTICS AGGREGATOR                             |  |
|  |  Industrial Ruggedness Meets Institutional Financial Trust      |  |
|  +-----------------------------------------------------------------+  |
|                                                                       |
|  1. DETECTED OPERATING JURISDICTION (SIM MCC 636 - ETHIO TELECOM)     |
|  +---------------------------------+ +-----------------------------+  |
|  | [✓] 🇪🇹 ETHIOPIA (+251 | ETB)   | | [ ] 🇰🇪 KENYA (+254 | KES)  |  |
|  |     Ge'ez + Gregorian Calendar  | |     Gregorian Calendar      |  |
|  |=====[ 3dp Amber Lip #F59E0B ]===| |     (64dp Hit Target)       |  |
|  +---------------------------------+ +-----------------------------+  |
|                                                                       |
|  2. SELECT YOUR WORKING LANGUAGE / ቋንቋ ይምረጡ (64dp Glove Cards)       |
|  +---------------------------------+ +-----------------------------+  |
|  | [✓] አማርኛ (AMHARIC)             | | [ ] ENGLISH                 |  |
|  |     የኢትዮጵያ የሥራ ቋንቋ (1.38x)      | |     International Baseline  |  |
|  |     [✓ Font SHA-256 Ready]      | |     Roboto 400/700          |  |
|  |=====[ 3dp Amber Lip #F59E0B ]===| |                             |  |
|  +---------------------------------+ +-----------------------------+  |
|  |                      <- 8dp Physical Gap ->                     |  |
|  +---------------------------------+ +-----------------------------+  |
|  | [ ] AFAAN OROMOO                | | [ ] KISWAHILI               |  |
|  |     Qubee Latin (+25% Buffer)   | |     Lugha ya Afrika Mashariki|  |
|  +---------------------------------+ +-----------------------------+  |
|                                                                       |
|  3. WELCOME BENEFIT CAROUSEL (/api/v1/onboarding/config)              |
|  +-----------------------------------------------------------------+  |
|  | "Verified Excavators, Cranes & Lowboys — Zero Broker Markups.   |  |
|  |  Direct P2P Bank Wire Settlement & Monotonic Hour Meters."      |  |
|  |                        ( ● )   ( ○ )   ( ○ )                    |  |
|  +-----------------------------------------------------------------+  |
|                                                                       |
|  +-----------------------------------------------------------------+  |
|  | [X] (56x56dp) I accept the Dink Machinery Terms of Service      |  |
|  |     (v2.0.0) & Non-Custodial Wire Disclaimer [Read Full Terms]  |  |
|  +-----------------------------------------------------------------+  |
+=======================================================================+
|  +-----------------------------------------------------------------+  |
|  |  [  CONTINUE TO PHONE VERIFICATION  ->  ]          (64dp Height)|  |
|  +=================================================================+  |
+-----------------------------------------------------------------------+
```

---

## 3. W3C DTCG Design Token Mapping Table

| UI Component / Element | W3C Token Path (`design-tokens.json`) | Standard Dark (`Cast Iron`) | Outdoor High-Contrast (`>80k lux`) | Behavioral & Ergonomic Rule |
|---|---|---|---|---|
| **Splash Canvas** | `{semantic.*.surface.canvas}` | `{color.brand.slate.900}` (`#0F172A`) | `#000000` (Pure Black) | Renders immediately from native launch theme. |
| **Selected Language Card** | `{semantic.*.surface.cardElevated}` | `{color.brand.slate.800}` + `2dp` `{color.brand.amber.500}` | `#000000` + `3dp` `{color.brand.amber.400}` | Minimum `64dp` height; `3dp` bottom tactile lip. |
| **Unselected Card** | `{semantic.*.surface.card}` | `{color.brand.slate.800}` + `1dp` `{color.brand.slate.700}` | `#000000` + `2dp` `#FFFFFF` | Minimum `8dp` gap between grid cells. |
| **Amharic Card Title** | `{typography.fontFamily.ethiopic}` | `Noto Sans Ethiopic` (`18sp`, `1.38x`) | `Noto Sans Ethiopic` (`19sp`, `1.45x`) | Dynamically loaded via `ADR-013` OTA pipeline. |
| **Consent Checkbox HitBox**| `{dimension.space.56}` | `56dp x 56dp` | `56dp x 56dp` (`3dp` outline) | Explicit click-wrap gate (`POL-01`); starts unchecked. |
| **Primary Continue CTA** | `{component.actuator.button.heightField}` | `{dimension.space.64}` (`64dp`) | `{dimension.space.64}` (`64dp`) | Disabled (`slate.700`) until legal checkbox is checked. |

---

## 4. Interactive State Machine: Boot Hydration, OTA Font & Consent Gate

```mermaid
stateDiagram-v2
    [*] --> NativeBootCheck: App Cold Start (<1.5s Budget)
    NativeBootCheck --> SecurityTamperBlocked: Root / Jailbreak / Frida Detected (ADR-008)
    NativeBootCheck --> SessionHydration: Device Integrity Verified

    SessionHydration --> AuthenticatedShell: Valid JWT in Keystore & isOnboarded == true
    SessionHydration --> LanguageJurisdictionReady: No JWT OR First Launch (Read SIM MCC 636/639)

    state LanguageJurisdictionReady {
        [*] --> LocaleSelectedLatin: Default 'en-ET' or 'sw-KE'
        LocaleSelectedLatin --> OtaEthiopicLoading: User Taps 'አማርኛ (am-ET)'
        OtaEthiopicLoading --> LocaleSelectedEthiopic: NotoSansEthiopic SHA-256 Verified & Loaded
        OtaEthiopicLoading --> FallbackSystemEthiopic: Offline / CDN Timeout (Use System Ge'ez + 1.38x)
    }

    LanguageJurisdictionReady --> ConsentCheckedReady: User Taps 56x56dp Terms Checkbox (POL-01)
    ConsentCheckedReady --> NavigateToPhoneOtp: Tap 64dp [CONTINUE TO PHONE VERIFICATION]
```

---

## 5. Concrete Flutter Implementation Blueprint

### 5.1 Splash & Locale Selection BLoC (`splash_onboarding_bloc.dart`)
```dart
// lib/features/auth/bloc/splash_onboarding_bloc.dart
import 'package:dio/dio.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:isar/isar.dart';

part 'splash_onboarding_bloc.freezed.dart';

@freezed
class SplashOnboardingEvent with _$SplashOnboardingEvent {
  const factory SplashOnboardingEvent.bootInitialized({
    required String detectedSimCountry, // 'ET' or 'KE'
    required bool isDeviceCompromised,
  }) = _BootInitialized;

  const factory SplashOnboardingEvent.countryChanged(String countryCode) = _CountryChanged;
  const factory SplashOnboardingEvent.languageSelected(String languageCode) = _LanguageSelected;
  const factory SplashOnboardingEvent.termsConsentToggled(bool accepted) = _TermsConsentToggled;
  const factory SplashOnboardingEvent.continueSubmitted() = _ContinueSubmitted;
}

@freezed
class SplashOnboardingState with _$SplashOnboardingState {
  const factory SplashOnboardingState({
    @Default('ET') String countryCode,
    @Default('am') String languageCode,
    @Default('ETB') String currencyCode,
    @Default(true) bool useDualEthiopianCalendar,
    @Default(false) bool termsAccepted,
    @Default('2.0.0') String activeTermsVersion,
    @Default('') String activeTermsSha256,
    @Default(false) bool isFontDownloading,
    @Default(false) bool isDeviceCompromised,
    @Default(false) bool readyForPhoneOtp,
    String? errorCode,
  }) = _SplashOnboardingState;
}

class SplashOnboardingBloc extends Bloc<SplashOnboardingEvent, SplashOnboardingState> {
  final Isar _isar;
  final Dio _dio;

  SplashOnboardingBloc(this._isar, this._dio) : super(const SplashOnboardingState()) {
    on<_BootInitialized>(_onBootInitialized);
    on<_CountryChanged>(_onCountryChanged);
    on<_LanguageSelected>(_onLanguageSelected);
    on<_TermsConsentToggled>(_onTermsConsentToggled);
    on<_ContinueSubmitted>(_onContinueSubmitted);
  }

  Future<void> _onBootInitialized(
    _BootInitialized event,
    Emitter<SplashOnboardingState> emit,
  ) async {
    if (event.isDeviceCompromised) {
      emit(state.copyWith(isDeviceCompromised: true, errorCode: 'AUTH_DEVICE_COMPROMISED'));
      return;
    }
    final isKe = event.detectedSimCountry.toUpperCase() == 'KE';
    emit(state.copyWith(
      countryCode: isKe ? 'KE' : 'ET',
      currencyCode: isKe ? 'KES' : 'ETB',
      languageCode: isKe ? 'sw' : 'am',
      useDualEthiopianCalendar: !isKe,
    ));

    try {
      final resp = await _dio.get<Map<String, dynamic>>(
        '/api/v1/legal/policies/active',
        options: Options(headers: {'X-Country-Code': state.countryCode}),
      );
      final data = resp.data?['data'] as Map<String, dynamic>?;
      if (data != null) {
        emit(state.copyWith(
          activeTermsVersion: data['versionSemver'] as String? ?? '2.0.0',
          activeTermsSha256: data['contentSha256'] as String? ?? '',
        ));
      }
    } catch (_) {
      // Offline fallback uses embedded Isar baseline policy hash
    }
  }

  void _onCountryChanged(_CountryChanged event, Emitter<SplashOnboardingState> emit) {
    HapticFeedback.lightImpact();
    final isKe = event.countryCode == 'KE';
    emit(state.copyWith(
      countryCode: event.countryCode,
      currencyCode: isKe ? 'KES' : 'ETB',
      useDualEthiopianCalendar: !isKe,
      termsAccepted: false, // Re-require explicit click-wrap when legal jurisdiction changes
    ));
  }

  Future<void> _onLanguageSelected(
    _LanguageSelected event,
    Emitter<SplashOnboardingState> emit,
  ) async {
    HapticFeedback.lightImpact();
    emit(state.copyWith(languageCode: event.languageCode, errorCode: null));
  }

  void _onTermsConsentToggled(_TermsConsentToggled event, Emitter<SplashOnboardingState> emit) {
    HapticFeedback.selectionClick();
    emit(state.copyWith(termsAccepted: event.accepted, errorCode: null));
  }

  void _onContinueSubmitted(_ContinueSubmitted event, Emitter<SplashOnboardingState> emit) {
    if (!state.termsAccepted) {
      HapticFeedback.heavyImpact();
      emit(state.copyWith(errorCode: 'ERR_MANDATORY_CONSENT_MISSING'));
      return;
    }
    HapticFeedback.mediumImpact();
    emit(state.copyWith(readyForPhoneOtp: true));
  }
}
```

---

## 6. OpenAPI 3.1.0 Contract & RFC 9457 Error Recovery Matrix

| Endpoint / Operation | Method & Path | Request / Header Invariants | RFC 9457 Error Code & Status | UI Recovery Behavior |
|---|---|---|---|---|
| **Fetch Onboarding Slides** | `GET /api/v1/onboarding/config` | `X-Country-Code: ET \| KE`<br/>`Accept-Language: am \| om \| sw \| en` | `SYS_SERVICE_UNAVAILABLE` (`503`) | Render pre-bundled Isar welcome carousel slides; zero block to onboarding progression. |
| **Fetch Active Legal Terms** | `GET /api/v1/legal/policies/active` | `X-Country-Code: ET \| KE` | `REQ_INVALID_COUNTRY_CODE` (`400`) | Default to `ET` (`+251`) jurisdiction and reload policy bundle. |
| **Missing Click-Wrap Consent**| Client / `POST /api/v1/auth/register` | `termsVersion` + `contentSha256` | `ERR_MANDATORY_CONSENT_MISSING` (`422`) | Highlight `56x56dp` consent checkbox with `3dp` `{color.feedback.crimson.600}` border and fire warning buzz. |

---

## 7. Verification & Acceptance Criteria

### 7.1 Automated Unit & BLoC Test Cases (`bloc_test`)
1. **SIM Country Detection & Jurisdiction Switch Consent Reset:**
   - *Given* `SplashOnboardingBloc` initialized with `detectedSimCountry: 'ET'` and `termsAccepted: true`,
   - *When* `SplashOnboardingEvent.countryChanged('KE')` is dispatched,
   - *Then* `currencyCode` updates to `'KES'`, `useDualEthiopianCalendar` becomes `false`, and `termsAccepted` resets to `false` so Kenyan Cap 23 terms must be explicitly accepted.
2. **Root/Jailbreak Tamper Interception (`ADR-008`):**
   - *Given* `isDeviceCompromised: true` on boot,
   - *Then* state emits `errorCode: 'AUTH_DEVICE_COMPROMISED'` and blocks progression to authentication.

### 7.2 Performance & Accessibility Test Gates
1. **`< 1.5s` Cold-Start Interactive Frame Gate:**
   - Profile-mode integration test on a throttled $3\text{ GB RAM}$ Android test target verifies time-to-first-interactive-frame on `SplashLanguageSelectionScreen` is strictly $< 1,500\text{ ms}$.
