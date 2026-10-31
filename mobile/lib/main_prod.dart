import 'package:flutter/material.dart';
import 'package:dink_mobile/core/config/app_config.dart';
import 'package:dink_mobile/core/di/injection.dart';
import 'package:dink_mobile/main.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  AppConfig.initialize(
    const AppConfig(
      environment: Environment.prod,
      appName: 'Dink Machinery',
      apiBaseUrl: 'https://api.dinkmachinery.com/api/v1',
      enableSslPinning: true,
      isarDbName: 'isar_prod',
    ),
  );

  await initDependencies();
  runApp(const MyApp());
}
