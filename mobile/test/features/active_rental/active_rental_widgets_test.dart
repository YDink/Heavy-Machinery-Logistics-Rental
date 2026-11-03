import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/active_rental/domain/entities/active_rental_hud.dart';
import 'package:dink_mobile/features/active_rental/presentation/bloc/active_rental_bloc.dart';
import 'package:dink_mobile/features/active_rental/presentation/views/active_rental_hud_view.dart';
import 'package:dink_mobile/features/active_rental/presentation/views/widgets/breakdown_report_modal.dart';
import 'package:dink_mobile/features/active_rental/presentation/views/widgets/extension_request_sheet.dart';
import 'package:dink_mobile/features/active_rental/presentation/views/widgets/shift_telemetry_progress_card.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

import 'active_rental_bloc_test.dart';

Widget buildTestWidget({
  required Widget child,
  Locale locale = const Locale('en'),
  ActiveRentalBloc? bloc,
}) {
  Widget content = child;
  if (bloc != null) {
    content = BlocProvider<ActiveRentalBloc>.value(
      value: bloc,
      child: child,
    );
  }

  return MaterialApp(
    locale: locale,
    localizationsDelegates: const [
      AppLocalizations.delegate,
      GlobalMaterialLocalizations.delegate,
      GlobalWidgetsLocalizations.delegate,
      GlobalCupertinoLocalizations.delegate,
    ],
    supportedLocales: AppLocalizations.supportedLocales,
    home: Scaffold(body: content),
  );
}

void main() {
  group('ShiftTelemetryProgressCard Widget Tests', () {
    const dummyHUD = ActiveRentalHUD(
      bookingId: 'bk-123',
      status: 'ACTIVE',
      fleetId: 'flt-001',
      fleetTitle: 'Caterpillar 320D Excavator',
      plateNumber: 'ET-3-12345',
      currentRentalDay: 2,
      totalRentalDays: 5,
      handoverBaselineHours: 1250.0,
      accumulatedHoursToday: 6.5,
      allowedHoursPerDay: 8.0,
      excessHoursTotal: 1.5,
      shiftTimerPaused: false,
      canExtend: true,
      canReportBreakdown: true,
      canInitiateReturn: true,
    );

    testWidgets('renders daily shift telemetry and 8.0 hr progress in English', (tester) async {
      await tester.pumpWidget(buildTestWidget(
        child: const ShiftTelemetryProgressCard(hud: dummyHUD),
        locale: const Locale('en'),
      ));
      await tester.pumpAndSettle();

      expect(find.text('Daily Shift Telemetry'), findsOneWidget);
      expect(find.text('6.5 / 8.0 hrs'), findsOneWidget);
      expect(find.text('1.5 hrs'), findsOneWidget);
      expect(find.byType(LinearProgressIndicator), findsOneWidget);
    });

    testWidgets('renders daily shift telemetry in Amharic', (tester) async {
      await tester.pumpWidget(buildTestWidget(
        child: const ShiftTelemetryProgressCard(hud: dummyHUD),
        locale: const Locale('am'),
      ));
      await tester.pumpAndSettle();

      expect(find.text('የዕለት ፈረቃ የስራ ሂደት'), findsOneWidget);
      expect(find.text('6.5 / 8.0 hrs'), findsOneWidget);
    });

    testWidgets('renders shift paused banner when breakdown is active', (tester) async {
      final pausedHUD = ActiveRentalHUD(
        bookingId: dummyHUD.bookingId,
        status: dummyHUD.status,
        fleetId: dummyHUD.fleetId,
        fleetTitle: dummyHUD.fleetTitle,
        plateNumber: dummyHUD.plateNumber,
        currentRentalDay: dummyHUD.currentRentalDay,
        totalRentalDays: dummyHUD.totalRentalDays,
        handoverBaselineHours: dummyHUD.handoverBaselineHours,
        accumulatedHoursToday: dummyHUD.accumulatedHoursToday,
        allowedHoursPerDay: dummyHUD.allowedHoursPerDay,
        excessHoursTotal: dummyHUD.excessHoursTotal,
        shiftTimerPaused: true, // Paused!
        canExtend: dummyHUD.canExtend,
        canReportBreakdown: dummyHUD.canReportBreakdown,
        canInitiateReturn: dummyHUD.canInitiateReturn,
      );

      await tester.pumpWidget(buildTestWidget(
        child: ShiftTelemetryProgressCard(hud: pausedHUD),
        locale: const Locale('en'),
      ));
      await tester.pumpAndSettle();

      expect(find.text('Shift Timer Paused (Breakdown Active)'), findsOneWidget);
    });
  });

  group('ActiveRentalHUDView Widget Tests', () {
    late MockActiveRentalRepository repo;
    late ActiveRentalBloc bloc;

    setUp(() {
      repo = MockActiveRentalRepository();
      bloc = ActiveRentalBloc(repository: repo);
    });

    tearDown(() {
      bloc.close();
    });

    testWidgets('renders HUD header, operator, progress card, and 56dp action buttons', (tester) async {
      await tester.pumpWidget(buildTestWidget(
        bloc: bloc,
        child: const ActiveRentalHUDView(bookingId: 'bk-123'),
        locale: const Locale('en'),
      ));
      await tester.runAsync(() async {
        await Future.delayed(const Duration(milliseconds: 50));
      });
      await tester.pumpAndSettle();

      // Fleet Header
      expect(find.text('Caterpillar 320D Excavator'), findsOneWidget);
      expect(find.text('ET-3-12345'), findsOneWidget);
      expect(find.text('Day 2 / 5'), findsOneWidget);

      // Operator details
      expect(find.text('Kassahun Bekele'), findsOneWidget);

      // Shift Progress
      expect(find.text('Daily Shift Telemetry'), findsOneWidget);

      // Action Buttons (56dp minimum height per accessibility rule)
      final breakdownBtn = find.text('Report Breakdown');
      final extendBtn = find.text('Extend Rental');
      expect(breakdownBtn, findsOneWidget);
      expect(extendBtn, findsOneWidget);

      final breakdownSize = tester.getSize(find.ancestor(of: breakdownBtn, matching: find.byType(SizedBox)).first);
      expect(breakdownSize.height, equals(56.0));

      final extendSize = tester.getSize(find.ancestor(of: extendBtn, matching: find.byType(SizedBox)).first);
      expect(extendSize.height, equals(56.0));
    });
  });

  group('BreakdownReportModal & ExtensionRequestSheet Tests', () {
    late MockActiveRentalRepository repo;
    late ActiveRentalBloc bloc;

    setUp(() {
      repo = MockActiveRentalRepository();
      bloc = ActiveRentalBloc(repository: repo);
    });

    tearDown(() {
      bloc.close();
    });

    testWidgets('BreakdownReportModal renders fault category and submit button', (tester) async {
      await tester.pumpWidget(buildTestWidget(
        bloc: bloc,
        child: const BreakdownReportModal(bookingId: 'bk-123', currentMeterHours: 1255.0),
        locale: const Locale('en'),
      ));
      await tester.pumpAndSettle();

      expect(find.text('Report Mechanical Breakdown'), findsOneWidget);
      expect(find.text('HYDRAULIC_CIRCUIT'), findsOneWidget);
      expect(find.text('Submit Breakdown & Pause Timer'), findsOneWidget);
    });

    testWidgets('ExtensionRequestSheet renders date selection and submit button', (tester) async {
      await tester.pumpWidget(buildTestWidget(
        bloc: bloc,
        child: const ExtensionRequestSheet(bookingId: 'bk-123'),
        locale: const Locale('en'),
      ));
      await tester.pumpAndSettle();

      expect(find.text('Request Rental Extension'), findsOneWidget);
      expect(find.text('Confirm Extension Request'), findsOneWidget);
    });
  });
}
