import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../core/database.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class PaymentsPage extends StatefulWidget {
  const PaymentsPage({super.key});

  @override
  State<PaymentsPage> createState() => _PaymentsPageState();
}

class _PaymentsPageState extends State<PaymentsPage> {
  List<Map<String, dynamic>> _links = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final local = await db.getAll('SELECT * FROM payment_links ORDER BY created_at DESC');
      _links = local.map((r) => Map<String, dynamic>.from(r)).toList();
      try {
        _links = await paymentsApi.listLinks();
      } catch (_) {}
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Pagos',
      currentIndex: 4,
      child: _isLoading
          ? const LoadingIndicator()
          : _links.isEmpty
              ? const EmptyState(title: 'Sin enlaces de pago', icon: Icons.link_off)
              : ListView.builder(
                  itemCount: _links.length,
                  itemBuilder: (context, i) {
                    final l = _links[i];
                    return ListTile(
                      leading: const CircleAvatar(child: Icon(Icons.link)),
                      title: Text(l['provider']?.toString() ?? 'Proveedor'),
                      subtitle: Text(l['status']?.toString() ?? ''),
                      trailing: Text(l['amount']?.toString() ?? ''),
                    );
                  },
                ),
    );
  }
}
