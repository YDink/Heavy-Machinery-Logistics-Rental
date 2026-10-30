import 'package:flutter/material.dart';
import 'package:dink_mobile/core/config/app_config.dart';

void main() {
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

  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    final config = AppConfig.instance;

    return MaterialApp(
      title: config.appName,
      debugShowCheckedModeBanner: config.isDev,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: config.isProd
              ? const Color(0xFFE5A93C)
              : (config.isStaging ? Colors.blue : Colors.orange),
        ),
      ),
      home: MyHomePage(title: config.appName),
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

    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Center(
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
            const SizedBox(height: 16),
            Text('API Endpoint: ${config.apiBaseUrl}'),
            const SizedBox(height: 24),
            const Text('You have pushed the button this many times:'),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
          ],
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
