import 'package:flutter/material.dart';
import 'package:dink_mobile/core/config/app_config.dart';
import 'package:dink_mobile/core/di/injection.dart';
import 'package:dink_mobile/main.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  AppConfig.initialize(
    const AppConfig(
      environment: Environment.staging,
      appName: 'Dink (Staging)',
      apiBaseUrl: 'https://staging-api.dinkmachinery.com/api/v1',
      enableSslPinning: false,
      isarDbName: 'isar_staging',
    ),
  );

  await initDependencies();
  runApp(const MyApp());
}
