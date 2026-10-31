import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/core/config/app_config.dart';
import 'package:dink_mobile/core/di/injection.dart';
import 'package:dink_mobile/core/localization/locale_controller.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_bloc.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_event.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_state.dart';
import 'package:dink_mobile/features/auth/presentation/pages/phone_input_page.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Fallback to dev config if executed directly without a flavor entrypoint
  try {
    AppConfig.instance;
  } catch (_) {
    AppConfig.initialize(
      const AppConfig(
        environment: Environment.dev,
        appName: 'Dink (Dev)',
        apiBaseUrl: 'http://10.0.2.2:8080/api/v1',
        enableSslPinning: false,
        isarDbName: 'isar_dev',
      ),
    );
  }

  await initDependencies();
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    final config = AppConfig.instance;
    final localeController = sl<LocaleController>();

    return BlocProvider<AuthBloc>(
      create: (_) => sl<AuthBloc>()..add(CheckAuthStatusEvent()),
      child: ValueListenableBuilder<Locale>(
        valueListenable: localeController,
        builder: (context, currentLocale, _) {
          return MaterialApp(
            title: config.appName,
            locale: currentLocale,
            localizationsDelegates: AppLocalizations.localizationsDelegates,
            supportedLocales: AppLocalizations.supportedLocales,
            debugShowCheckedModeBanner: config.isDev,
            theme: ThemeData(
              colorScheme: ColorScheme.fromSeed(
                seedColor: config.isProd
                    ? const Color(0xFFE5A93C)
                    : (config.isStaging ? Colors.blue : Colors.orange),
              ),
              useMaterial3: true,
            ),
            home: MyHomePage(title: config.appName),
          );
        },
      ),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _counter = 0;

  void _incrementCounter() {
    setState(() {
      _counter++;
    });
  }

  @override
  Widget build(BuildContext context) {
    final config = AppConfig.instance;
    final localeController = sl<LocaleController>();
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
        actions: [
          // Dynamic Language Switcher (ADR-013, ADR-019)
          TextButton(
            onPressed: () => localeController.toggleLocale(),
            child: Text(
              localeController.isAmharic ? '🇬🇧 EN' : '🇪🇹 አማ',
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          BlocBuilder<AuthBloc, AuthState>(
            builder: (context, state) {
              if (state is AuthenticatedState) {
                return IconButton(
                  icon: const Icon(Icons.logout),
                  tooltip: 'Sign Out',
                  onPressed: () {
                    context.read<AuthBloc>().add(LogoutEvent());
                  },
                );
              }
              return const SizedBox.shrink();
            },
          ),
        ],
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                decoration: BoxDecoration(
                  color: config.isProd
                      ? Colors.green.shade100
                      : (config.isStaging ? Colors.blue.shade100 : Colors.amber.shade100),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  'Environment: ${config.environment.name.toUpperCase()}',
                  style: const TextStyle(fontWeight: FontWeight.bold),
                ),
              ),
              const SizedBox(height: 12),
              Text('API: ${config.apiBaseUrl}', style: TextStyle(color: Colors.grey.shade700, fontSize: 12)),
              const SizedBox(height: 24),

              // Auth Status Indicator Card
              BlocBuilder<AuthBloc, AuthState>(
                builder: (context, state) {
                  if (state is AuthenticatedState) {
                    return Card(
                      color: Colors.green.shade50,
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Column(
                          children: [
                            Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Icon(Icons.check_circle, color: Colors.green.shade700),
                                const SizedBox(width: 8),
                                Text(
                                  'Authenticated',
                                  style: TextStyle(
                                    fontWeight: FontWeight.bold,
                                    color: Colors.green.shade900,
                                  ),
                                ),
                              ],
                            ),
                            if (state.user != null) ...[
                              const SizedBox(height: 8),
                              Text('User: ${state.user!.phoneNumber} (${state.user!.role.name.toUpperCase()})'),
                            ],
                          ],
                        ),
                      ),
                    );
                  }

                  return ElevatedButton.icon(
                    onPressed: () {
                      Navigator.of(context).push(
                        MaterialPageRoute(
                          builder: (_) => BlocProvider.value(
                            value: context.read<AuthBloc>(),
                            child: const PhoneInputPage(),
                          ),
                        ),
                      );
                    },
                    icon: const Icon(Icons.lock_outline),
                    label: Text(l10n.sendOtpButton),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFFE5A93C),
                      foregroundColor: Colors.white,
                      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                    ),
                  );
                },
              ),

              const SizedBox(height: 36),
              const Text('You have pushed the button this many times:'),
              Text(
                '$_counter',
                style: Theme.of(context).textTheme.headlineMedium,
              ),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _incrementCounter,
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ),
    );
  }
}
