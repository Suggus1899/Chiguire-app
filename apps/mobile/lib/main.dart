import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/database.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await openDatabase();
  runApp(const ProviderScope(child: ChiguireApp()));
}

class ChiguireApp extends StatelessWidget {
  const ChiguireApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Chiguire',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF2563EB)),
        useMaterial3: true,
      ),
      home: const _Home(),
    );
  }
}

// Phase 0 placeholder — full routing and auth in Phase 1
class _Home extends StatelessWidget {
  const _Home();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Chiguire',
              style: Theme.of(context).textTheme.headlineLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: const Color(0xFF2563EB),
                  ),
            ),
            const SizedBox(height: 8),
            const Text('ERP offline-first · Fase 0'),
          ],
        ),
      ),
    );
  }
}
