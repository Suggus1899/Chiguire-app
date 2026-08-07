import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../core/database.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class PurchaseListPage extends StatefulWidget {
  const PurchaseListPage({super.key});

  @override
  State<PurchaseListPage> createState() => _PurchaseListPageState();
}

class _PurchaseListPageState extends State<PurchaseListPage> {
  List<Map<String, dynamic>> _orders = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final local = await db.getAll('SELECT * FROM purchase_orders ORDER BY created_at DESC');
      _orders = local.map((r) => Map<String, dynamic>.from(r)).toList();
      try {
        _orders = await purchasesApi.list();
      } catch (_) {}
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Compras',
      currentIndex: 4,
      child: _isLoading
          ? const LoadingIndicator()
          : _orders.isEmpty
              ? const EmptyState(title: 'Sin órdenes de compra', icon: Icons.shopping_cart_outlined)
              : ListView.builder(
                  itemCount: _orders.length,
                  itemBuilder: (context, i) {
                    final o = _orders[i];
                    return ListTile(
                      leading: const CircleAvatar(child: Icon(Icons.shopping_cart)),
                      title: Text(o['number']?.toString() ?? 'Sin número'),
                      subtitle: Text(o['status']?.toString() ?? ''),
                      trailing: Text(o['total']?.toString() ?? ''),
                    );
                  },
                ),
    );
  }
}
