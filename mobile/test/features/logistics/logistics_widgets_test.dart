import 'package:dink_mobile/features/logistics/domain/usecases/accept_transshipment_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/advance_logistics_dispatch_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/create_logistics_sub_booking_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/match_logistics_carriers_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/register_external_transport_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/report_transshipment_alert_usecase.dart';
import 'package:dink_mobile/features/logistics/domain/usecases/trigger_reverse_cascade_usecase.dart';
import 'package:dink_mobile/features/logistics/presentation/bloc/logistics_bloc.dart';
import 'package:dink_mobile/features/logistics/presentation/views/carrier_dispatch_view.dart';
import 'package:dink_mobile/features/logistics/presentation/views/logistics_mobilization_view.dart';
import 'package:dink_mobile/features/logistics/presentation/views/transshipment_cascade_view.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';

import 'logistics_bloc_test.dart';

Widget _wrapWithApp(Widget child, {Locale locale = const Locale('en')}) {
  final repo = FakeLogisticsRepository();
  final bloc = LogisticsBloc(
    matchCarriersUseCase: MatchLogisticsCarriersUseCase(repo),
    createSubBookingUseCase: CreateLogisticsSubBookingUseCase(repo),
    registerExternalTransportUseCase: RegisterExternalTransportUseCase(repo),
    advanceDispatchUseCase: AdvanceLogisticsDispatchUseCase(repo),
    reportTransshipmentAlertUseCase: ReportTransshipmentAlertUseCase(repo),
    acceptTransshipmentUseCase: AcceptTransshipmentUseCase(repo),
    triggerReverseCascadeUseCase: TriggerReverseCascadeUseCase(repo),
  );

  return MaterialApp(
    locale: locale,
    localizationsDelegates: const [
      AppLocalizations.delegate,
      GlobalMaterialLocalizations.delegate,
      GlobalWidgetsLocalizations.delegate,
      GlobalCupertinoLocalizations.delegate,
    ],
    supportedLocales: const [Locale('en'), Locale('am')],
    home: BlocProvider<LogisticsBloc>.value(
      value: bloc,
      child: child,
    ),
  );
}

void main() {
  group('Logistics Mobile Widget Tests', () {
    testWidgets(
      'LogisticsMobilizationView renders matched carrier and >=56dp Book CTA',
      (tester) async {
        await tester.pumpWidget(
          _wrapWithApp(const LogisticsMobilizationView(bookingId: 'bk-1')),
        );
        await tester.pumpAndSettle();

        expect(find.textContaining('Sinotruk'), findsOneWidget);
        expect(find.text('BOOK LOWBOY CARRIER'), findsOneWidget);

        final btnSize = tester.getSize(
          find.widgetWithText(ElevatedButton, 'BOOK LOWBOY CARRIER'),
        );
        expect(btnSize.height, greaterThanOrEqualTo(56.0));
      },
    );

    testWidgets(
      'LogisticsMobilizationView renders bilingual Amharic labels',
      (tester) async {
        await tester.pumpWidget(
          _wrapWithApp(
            const LogisticsMobilizationView(bookingId: 'bk-1'),
            locale: const Locale('am'),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('የሎውቦይ ማጓጓዣ ዝግጅት'), findsOneWidget);
        expect(find.text('ሎውቦይ አጓጓዥ ይዘዙ'), findsOneWidget);
      },
    );

    testWidgets(
      'CarrierDispatchView renders 4 winching reticle angles and 56dp buttons',
      (tester) async {
        await tester.pumpWidget(
          _wrapWithApp(const CarrierDispatchView(subBookingId: 'sub-1')),
        );
        await tester.pumpAndSettle();

        expect(find.text('01. RIG PROFILE'), findsOneWidget);
        expect(find.text('02. CHAIN BINDERS (4-PT)'), findsOneWidget);
        expect(find.text('03. BOOM PINNED'), findsOneWidget);
        expect(find.text('04. RAMPS LOCKED'), findsOneWidget);

        final confirmBtn = tester.getSize(
          find.widgetWithText(ElevatedButton, 'CONFIRM LOAD & START TRANSIT'),
        );
        expect(confirmBtn.height, greaterThanOrEqualTo(56.0));

        final alertBtn = tester.getSize(
          find.widgetWithText(OutlinedButton, 'REPORT ROADSIDE BREAKDOWN'),
        );
        expect(alertBtn.height, greaterThanOrEqualTo(56.0));
      },
    );

    testWidgets(
      'TransshipmentCascadeView renders 60m SLA bounty and Force Majeure exoneration',
      (tester) async {
        await tester.pumpWidget(
          _wrapWithApp(
            const TransshipmentCascadeView(
              parentBookingId: 'bk-1',
              subBookingId: 'sub-1',
              replacementFleetId: 'flt-2',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.textContaining('+20% URGENT BOUNTY'), findsOneWidget);
        expect(find.text('0 ETB (EXONERATED)'), findsOneWidget);

        final acceptBtn = tester.getSize(
          find.widgetWithText(ElevatedButton, 'ACCEPT HOT-SWAP TRANSFER'),
        );
        expect(acceptBtn.height, greaterThanOrEqualTo(56.0));

        final cascadeBtn = tester.getSize(
          find.widgetWithText(OutlinedButton, 'EXECUTE REVERSE CASCADE'),
        );
        expect(cascadeBtn.height, greaterThanOrEqualTo(56.0));
      },
    );
  });
}
