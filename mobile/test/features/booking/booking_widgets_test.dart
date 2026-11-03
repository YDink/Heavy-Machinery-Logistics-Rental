import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/booking/domain/entities/booking.dart';
import 'package:dink_mobile/features/booking/presentation/widgets/booking_price_breakdown_card.dart';
import 'package:dink_mobile/features/booking/presentation/widgets/booking_status_badge.dart';
import 'package:dink_mobile/features/booking/presentation/widgets/hold_countdown_timer.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

Widget buildTestWidget({
  required Widget child,
  Locale locale = const Locale('en'),
}) {
  return MaterialApp(
    locale: locale,
    localizationsDelegates: const [
      AppLocalizations.delegate,
      GlobalMaterialLocalizations.delegate,
      GlobalWidgetsLocalizations.delegate,
      GlobalCupertinoLocalizations.delegate,
    ],
    supportedLocales: AppLocalizations.supportedLocales,
    home: Scaffold(body: child),
  );
}

void main() {
  group('BookingStatusBadge Widget Tests', () {
    testWidgets('renders Pending Payment in English and Amharic', (tester) async {
      late AppLocalizations l10nEn;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10nEn = AppLocalizations.of(context)!;
            return const BookingStatusBadge(status: BookingStatus.pendingPayment);
          },
        ),
        locale: const Locale('en'),
      ));
      await tester.pumpAndSettle();
      expect(find.text(l10nEn.bookingStatusPendingPayment), findsOneWidget);

      late AppLocalizations l10nAm;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10nAm = AppLocalizations.of(context)!;
            return const BookingStatusBadge(status: BookingStatus.pendingPayment);
          },
        ),
        locale: const Locale('am'),
      ));
      await tester.pumpAndSettle();
      expect(find.text(l10nAm.bookingStatusPendingPayment), findsOneWidget);
    });

    testWidgets('renders Confirmed status in English and Amharic', (tester) async {
      late AppLocalizations l10nEn;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10nEn = AppLocalizations.of(context)!;
            return const BookingStatusBadge(status: BookingStatus.confirmed);
          },
        ),
        locale: const Locale('en'),
      ));
      await tester.pumpAndSettle();
      expect(find.text(l10nEn.bookingStatusConfirmed), findsOneWidget);

      late AppLocalizations l10nAm;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10nAm = AppLocalizations.of(context)!;
            return const BookingStatusBadge(status: BookingStatus.confirmed);
          },
        ),
        locale: const Locale('am'),
      ));
      await tester.pumpAndSettle();
      expect(find.text(l10nAm.bookingStatusConfirmed), findsOneWidget);
    });

    testWidgets('renders Active status in English and Amharic', (tester) async {
      late AppLocalizations l10nEn;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10nEn = AppLocalizations.of(context)!;
            return const BookingStatusBadge(status: BookingStatus.active);
          },
        ),
        locale: const Locale('en'),
      ));
      await tester.pumpAndSettle();
      expect(find.text(l10nEn.bookingStatusActive), findsOneWidget);

      late AppLocalizations l10nAm;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10nAm = AppLocalizations.of(context)!;
            return const BookingStatusBadge(status: BookingStatus.active);
          },
        ),
        locale: const Locale('am'),
      ));
      await tester.pumpAndSettle();
      expect(find.text(l10nAm.bookingStatusActive), findsOneWidget);
    });

    testWidgets('renders Cancelled status in English and Amharic', (tester) async {
      late AppLocalizations l10nEn;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10nEn = AppLocalizations.of(context)!;
            return const BookingStatusBadge(status: BookingStatus.cancelled);
          },
        ),
        locale: const Locale('en'),
      ));
      await tester.pumpAndSettle();
      expect(find.text(l10nEn.bookingStatusCancelled), findsOneWidget);

      late AppLocalizations l10nAm;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10nAm = AppLocalizations.of(context)!;
            return const BookingStatusBadge(status: BookingStatus.cancelled);
          },
        ),
        locale: const Locale('am'),
      ));
      await tester.pumpAndSettle();
      expect(find.text(l10nAm.bookingStatusCancelled), findsOneWidget);
    });
  });

  group('BookingPriceBreakdownCard Widget Tests', () {
    testWidgets('renders price breakdown with minor units correctly formatted', (tester) async {
      late AppLocalizations l10n;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10n = AppLocalizations.of(context)!;
            return const BookingPriceBreakdownCard(
              durationDays: 3,
              dailyRateMinor: 2000000,
              baseHireFeeMinor: 6000000,
              attachmentsFeeMinor: 500000,
              grossRentalAmountMinor: 6500000,
              prepaymentAmountMinor: 3250000,
              remainingBalanceMinor: 3250000,
              currency: 'ETB',
            );
          },
        ),
        locale: const Locale('en'),
      ));
      await tester.pumpAndSettle();

      expect(find.text(l10n.bookingDetails), findsOneWidget);
      expect(find.text('3 ${l10n.daysUnit}'), findsOneWidget);
      expect(find.text('60,000.00 ETB'), findsOneWidget);
      expect(find.text('5,000.00 ETB'), findsOneWidget);
      expect(find.text('65,000.00 ETB'), findsOneWidget);
      expect(find.text('32,500.00 ETB'), findsNWidgets(2)); // Advance and Remaining
    });

    testWidgets('renders price breakdown in Amharic with zero missing keys', (tester) async {
      late AppLocalizations l10n;
      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10n = AppLocalizations.of(context)!;
            return const BookingPriceBreakdownCard(
              durationDays: 5,
              dailyRateMinor: 1000000,
              baseHireFeeMinor: 5000000,
              attachmentsFeeMinor: 0,
              grossRentalAmountMinor: 5000000,
              prepaymentAmountMinor: 2500000,
              remainingBalanceMinor: 2500000,
              currency: 'ETB',
            );
          },
        ),
        locale: const Locale('am'),
      ));
      await tester.pumpAndSettle();

      expect(find.text(l10n.bookingDetails), findsOneWidget);
      expect(find.text('5 ${l10n.daysUnit}'), findsOneWidget);
      expect(find.text(l10n.advancePrepayment), findsOneWidget);
      expect(find.text(l10n.remainingBalance), findsOneWidget);
    });
  });

  group('HoldCountdownTimer Widget Tests', () {
    testWidgets('displays remaining minutes and warning banner', (tester) async {
      late AppLocalizations l10n;
      final futureExpiry = DateTime.now().add(const Duration(minutes: 45, seconds: 30));

      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10n = AppLocalizations.of(context)!;
            return HoldCountdownTimer(expiresAt: futureExpiry);
          },
        ),
        locale: const Locale('en'),
      ));
      await tester.pump();

      expect(find.byIcon(Icons.timer_outlined), findsOneWidget);
      expect(find.textContaining(l10n.paymentHoldTimer), findsOneWidget);
    });

    testWidgets('displays hold expired state when expiresAt is in past', (tester) async {
      late AppLocalizations l10n;
      final pastExpiry = DateTime.now().subtract(const Duration(minutes: 5));

      await tester.pumpWidget(buildTestWidget(
        child: Builder(
          builder: (context) {
            l10n = AppLocalizations.of(context)!;
            return HoldCountdownTimer(expiresAt: pastExpiry);
          },
        ),
        locale: const Locale('en'),
      ));
      await tester.pump();

      expect(find.text('${l10n.paymentHoldTimer}: 00:00'), findsOneWidget);
    });
  });
}
