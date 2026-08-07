import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';

class ManufacturingListPage extends StatefulWidget {
  const ManufacturingListPage({super.key});

  @override
  State<ManufacturingListPage> createState() => _ManufacturingListPageState();
}

class _ManufacturingListPageState extends State<ManufacturingListPage> {
  List<Map<String, dynamic>> _orders = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      _orders = await manufacturingApi.list();
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _start(String id) async {
    try {
      await manufacturingApi.start(id);
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  Future<void> _complete(String id) async {
    try {
      await manufacturingApi.complete(id);
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Producción',
      currentIndex: 2,
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/manufacturing/new'),
        child: const Icon(Icons.add),
      ),
      child: _isLoading
          ? const LoadingIndicator()
          : _error != null
              ? AppErrorWidget(message: _error!, onRetry: _load)
              : _orders.isEmpty
                  ? const EmptyState(title: 'Sin órdenes de producción', icon: Icons.factory_outlined)
                  : ListView.builder(
                      itemCount: _orders.length,
                      itemBuilder: (context, i) {
                        final o = _orders[i];
                        final status = o['status']?.toString() ?? 'draft';
                        return ListTile(
                          leading: const Icon(Icons.factory),
                          title: Text(o['product_id']?.toString() ?? ''),
                          subtitle: Text('Almacén: ${o['warehouse_id']?.toString() ?? ''} · Cant: ${o['quantity']?.toString() ?? '0'} · $status'),
                          trailing: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              if (status == 'draft')
                                TextButton(onPressed: () => _start(o['id']?.toString() ?? ''), child: const Text('Iniciar')),
                              if (status == 'in_progress')
                                TextButton(onPressed: () => _complete(o['id']?.toString() ?? ''), child: const Text('Completar')),
                            ],
                          ),
                        );
                      },
                    ),
    );
  }
}
