/// Multi-environment configuration for Dink Machinery mobile application.
/// Strictly governed by ADR-017 (Multi-Environment Architecture & Flavor Strategy).
enum Environment {
  dev,
  staging,
  prod,
}

/// Holds active runtime configuration parameters for the mobile client.
class AppConfig {
  final Environment environment;
  final String appName;
  final String apiBaseUrl;
  final bool enableSslPinning;
  final String isarDbName;

  const AppConfig({
    required this.environment,
    required this.appName,
    required this.apiBaseUrl,
    required this.enableSslPinning,
    required this.isarDbName,
  });

  static AppConfig? _instance;

  /// Returns the singleton configuration instance.
  static AppConfig get instance {
    final current = _instance;
    if (current == null) {
      throw StateError(
        'AppConfig has not been initialized. Call AppConfig.initialize() before accessing instance.',
      );
    }
    return current;
  }

  /// Initializes the singleton configuration.
  static void initialize(AppConfig config) {
    _instance = config;
  }

  /// Resets the configuration for testing purposes.
  static void reset() {
    _instance = null;
  }

  bool get isDev => environment == Environment.dev;
  bool get isStaging => environment == Environment.staging;
  bool get isProd => environment == Environment.prod;
}
