import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:dink_mobile/core/config/app_config.dart';
import 'package:dink_mobile/main.dart';

void main() {
  setUp(() {
    AppConfig.initialize(
      const AppConfig(
        environment: Environment.dev,
        appName: 'Dink (Dev)',
        apiBaseUrl: 'http://10.0.2.2:8080/api/v1',
        enableSslPinning: false,
        isarDbName: 'isar_dev',
      ),
    );
  });

  testWidgets('Counter increments smoke test', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(const MyApp());

    // Verify that our counter starts at 0.
    expect(find.text('0'), findsOneWidget);
    expect(find.text('1'), findsNothing);

    // Tap the '+' icon and trigger a frame.
    await tester.tap(find.byIcon(Icons.add));
    await tester.pump();

    // Verify that our counter has incremented.
    expect(find.text('0'), findsNothing);
    expect(find.text('1'), findsOneWidget);
  });
}
