import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../core/database.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class FiscalPage extends StatefulWidget {
  const FiscalPage({super.key});

  @override
  State<FiscalPage> createState() => _FiscalPageState();
}

class _FiscalPageState extends State<FiscalPage> {
  List<Map<String, dynamic>> _withholdings = [];
  Map<String, dynamic>? _rate;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final local = await db.getAll('SELECT * FROM tax_withholdings ORDER BY created_at DESC');
      _withholdings = local.map((r) => Map<String, dynamic>.from(r)).toList();
      try {
        _withholdings = await fiscalApi.listWithholdings();
        _rate = await fiscalApi.getExchangeRate('USD');
      } catch (_) {}
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Fiscal',
      currentIndex: 4,
      child: _isLoading
          ? const LoadingIndicator()
          : ListView(
              padding: const EdgeInsets.all(16),
              children: [
                if (_rate != null) ...[
                  Card(
                    child: ListTile(
                      leading: const Icon(Icons.currency_exchange),
                      title: const Text('Tasa USD → VES'),
                      trailing: Text(_rate!['rate_to_ves']?.toString() ?? 'N/A'),
                    ),
                  ),
                  const SizedBox(height: 16),
                ],
                Text('Retenciones', style: Theme.of(context).textTheme.titleMedium),
                const SizedBox(height: 8),
                if (_withholdings.isEmpty)
                  const EmptyState(title: 'Sin retenciones', icon: Icons.receipt_outlined)
                else
                  ..._withholdings.map((w) => ListTile(
                        leading: const Icon(Icons.percent),
                        title: Text(w['type']?.toString() ?? ''),
                        trailing: Text(w['amount']?.toString() ?? ''),
                      )),
              ],
            ),
    );
  }
}
