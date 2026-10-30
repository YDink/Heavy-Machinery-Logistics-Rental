import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/core/config/app_config.dart';

void main() {
  group('AppConfig Tests', () {
    setUp(() {
      AppConfig.reset();
    });

    test('throws StateError when accessed before initialization', () {
      expect(() => AppConfig.instance, throwsA(isA<StateError>()));
    });

    test('initializes correctly for Dev environment', () {
      const devConfig = AppConfig(
        environment: Environment.dev,
        appName: 'Dink (Dev)',
        apiBaseUrl: 'http://10.0.2.2:8080/api/v1',
        enableSslPinning: false,
        isarDbName: 'isar_dev',
      );

      AppConfig.initialize(devConfig);

      expect(AppConfig.instance.environment, equals(Environment.dev));
      expect(AppConfig.instance.appName, equals('Dink (Dev)'));
      expect(AppConfig.instance.apiBaseUrl, equals('http://10.0.2.2:8080/api/v1'));
      expect(AppConfig.instance.enableSslPinning, isFalse);
      expect(AppConfig.instance.isarDbName, equals('isar_dev'));
      expect(AppConfig.instance.isDev, isTrue);
      expect(AppConfig.instance.isStaging, isFalse);
      expect(AppConfig.instance.isProd, isFalse);
    });

    test('initializes correctly for Staging environment', () {
      const stagingConfig = AppConfig(
        environment: Environment.staging,
        appName: 'Dink (Staging)',
        apiBaseUrl: 'https://staging-api.dinkmachinery.com/api/v1',
        enableSslPinning: false,
        isarDbName: 'isar_staging',
      );

      AppConfig.initialize(stagingConfig);

      expect(AppConfig.instance.environment, equals(Environment.staging));
      expect(AppConfig.instance.appName, equals('Dink (Staging)'));
      expect(
        AppConfig.instance.apiBaseUrl,
        equals('https://staging-api.dinkmachinery.com/api/v1'),
      );
      expect(AppConfig.instance.enableSslPinning, isFalse);
      expect(AppConfig.instance.isarDbName, equals('isar_staging'));
      expect(AppConfig.instance.isDev, isFalse);
      expect(AppConfig.instance.isStaging, isTrue);
      expect(AppConfig.instance.isProd, isFalse);
    });

    test('initializes correctly for Prod environment with SSL pinning enabled', () {
      const prodConfig = AppConfig(
        environment: Environment.prod,
        appName: 'Dink Machinery',
        apiBaseUrl: 'https://api.dinkmachinery.com/api/v1',
        enableSslPinning: true,
        isarDbName: 'isar_prod',
      );

      AppConfig.initialize(prodConfig);

      expect(AppConfig.instance.environment, equals(Environment.prod));
      expect(AppConfig.instance.appName, equals('Dink Machinery'));
      expect(
        AppConfig.instance.apiBaseUrl,
        equals('https://api.dinkmachinery.com/api/v1'),
      );
      expect(AppConfig.instance.enableSslPinning, isTrue);
      expect(AppConfig.instance.isarDbName, equals('isar_prod'));
      expect(AppConfig.instance.isDev, isFalse);
      expect(AppConfig.instance.isStaging, isFalse);
      expect(AppConfig.instance.isProd, isTrue);
    });
  });
}
