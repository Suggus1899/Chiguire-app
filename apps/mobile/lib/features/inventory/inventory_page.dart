import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';
import '../../core/database.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class InventoryPage extends ConsumerStatefulWidget {
  const InventoryPage({super.key});

  @override
  ConsumerState<InventoryPage> createState() => _InventoryPageState();
}

class _InventoryPageState extends ConsumerState<InventoryPage> {
  List<Map<String, dynamic>> _movements = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _isLoading = true);
    try {
      final local = await db.getAll('SELECT * FROM stock_movements ORDER BY created_at DESC');
      _movements = local.map((r) => Map<String, dynamic>.from(r)).toList();
      try {
        _movements = await inventoryApi.listMovements();
      } catch (_) {}
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Inventario',
      currentIndex: 2,
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/inventory/new'),
        child: const Icon(Icons.add),
      ),
      child: _isLoading
          ? const LoadingIndicator()
          : _movements.isEmpty
              ? const EmptyState(title: 'Sin movimientos', icon: Icons.swap_horiz)
              : ListView.builder(
                  itemCount: _movements.length,
                  itemBuilder: (context, i) {
                    final m = _movements[i];
                    final qty = m['quantity']?.toString() ?? '0';
                    final type = m['movement_type']?.toString() ?? '';
                    return ListTile(
                      leading: Icon(
                        type == 'in' ? Icons.arrow_downward : Icons.arrow_upward,
                        color: type == 'in' ? Colors.green : Colors.red,
                      ),
                      title: Text(m['reference']?.toString() ?? 'Movimiento'),
                      subtitle: Text('Producto: ${m['product_id']?.toString() ?? ''}'),
                      trailing: Text(qty),
                    );
                  },
                ),
    );
  }
}
