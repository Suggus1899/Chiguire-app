import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';

class TransferListPage extends StatefulWidget {
  const TransferListPage({super.key});

  @override
  State<TransferListPage> createState() => _TransferListPageState();
}

class _TransferListPageState extends State<TransferListPage> {
  List<Map<String, dynamic>> _transfers = [];
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
      _transfers = await transfersApi.list();
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _ship(String id) async {
    try {
      await transfersApi.ship(id);
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  Future<void> _receive(String id) async {
    try {
      await transfersApi.receive(id);
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Transferencias',
      currentIndex: 2,
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/transfers/new'),
        child: const Icon(Icons.add),
      ),
      child: _isLoading
          ? const LoadingIndicator()
          : _error != null
              ? AppErrorWidget(message: _error!, onRetry: _load)
              : _transfers.isEmpty
                  ? const EmptyState(title: 'Sin transferencias', icon: Icons.swap_horiz)
                  : ListView.builder(
                      itemCount: _transfers.length,
                      itemBuilder: (context, i) {
                        final t = _transfers[i];
                        final status = t['status']?.toString() ?? 'draft';
                        return ListTile(
                          leading: const Icon(Icons.swap_horiz),
                          title: Text('${t['from_warehouse_id']?.toString() ?? ''} → ${t['to_warehouse_id']?.toString() ?? ''}'),
                          subtitle: Text('Estado: $status · ${t['date']?.toString() ?? ''}'),
                          trailing: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              if (status == 'draft')
                                TextButton(onPressed: () => _ship(t['id']?.toString() ?? ''), child: const Text('Enviar')),
                              if (status == 'shipped')
                                TextButton(onPressed: () => _receive(t['id']?.toString() ?? ''), child: const Text('Recibir')),
                            ],
                          ),
                        );
                      },
                    ),
    );
  }
}
