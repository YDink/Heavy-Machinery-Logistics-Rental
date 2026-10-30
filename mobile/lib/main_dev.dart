import 'package:flutter/material.dart';
import 'package:dink_mobile/core/config/app_config.dart';
import 'package:dink_mobile/main.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();

  AppConfig.initialize(
    const AppConfig(
      environment: Environment.dev,
      appName: 'Dink (Dev)',
      apiBaseUrl: 'http://10.0.2.2:8080/api/v1',
      enableSslPinning: false,
      isarDbName: 'isar_dev',
    ),
  );

  runApp(const MyApp());
}
