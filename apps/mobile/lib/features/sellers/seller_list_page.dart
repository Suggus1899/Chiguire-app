import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';

class SellerListPage extends StatefulWidget {
  const SellerListPage({super.key});

  @override
  State<SellerListPage> createState() => _SellerListPageState();
}

class _SellerListPageState extends State<SellerListPage> {
  List<Map<String, dynamic>> _sellers = [];
  bool _isLoading = true;
  String? _error;
  String _statusFilter = 'all';

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
      _sellers = await sellersApi.list();
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  List<Map<String, dynamic>> get _filtered {
    if (_statusFilter == 'all') return _sellers;
    return _sellers.where((s) {
      final active = s['is_active'] == true || s['is_active'] == 1;
      return _statusFilter == 'active' ? active : !active;
    }).toList();
  }

  Future<void> _delete(String id) async {
    try {
      await sellersApi.delete(id);
      _load();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Vendedores',
      currentIndex: 4,
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/sellers/new'),
        child: const Icon(Icons.add),
      ),
      child: _isLoading
          ? const LoadingIndicator()
          : _error != null && _sellers.isEmpty
              ? AppErrorWidget(message: _error!, onRetry: _load)
              : Column(
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(12),
                      child: SegmentedButton<String>(
                        segments: const [
                          ButtonSegment(value: 'all', label: Text('Todos')),
                          ButtonSegment(value: 'active', label: Text('Activos')),
                          ButtonSegment(value: 'inactive', label: Text('Inactivos')),
                        ],
                        selected: {_statusFilter},
                        onSelectionChanged: (s) => setState(() => _statusFilter = s.first),
                      ),
                    ),
                    Expanded(
                      child: _sellers.isEmpty
                          ? const EmptyState(title: 'Sin vendedores', icon: Icons.badge_outlined)
                          : ListView.builder(
                              itemCount: _filtered.length,
                              itemBuilder: (context, i) {
                                final s = _filtered[i];
                                final active = s['is_active'] == true || s['is_active'] == 1;
                                return ListTile(
                                  leading: const CircleAvatar(child: Icon(Icons.person)),
                                  title: Text(s['name']?.toString() ?? 'Sin nombre'),
                                  subtitle: Text(s['email']?.toString() ?? ''),
                                  trailing: Row(
                                    mainAxisSize: MainAxisSize.min,
                                    children: [
                                      Text('${s['commission_pct'] ?? 0}%'),
                                      const SizedBox(width: 8),
                                      Icon(
                                        active ? Icons.check_circle : Icons.cancel,
                                        color: active ? Colors.green : Colors.grey,
                                      ),
                                    ],
                                  ),
                                  onTap: () => context.go('/sellers/new?id=${s['id']}'),
                                  onLongPress: () async {
                                    final ok = await showDialog<bool>(
                                      context: context,
                                      builder: (_) => AlertDialog(
                                        title: const Text('Eliminar'),
                                        content: const Text('¿Eliminar este vendedor?'),
                                        actions: [
                                          TextButton(onPressed: () => Navigator.pop(context, false), child: const Text('Cancelar')),
                                          FilledButton(onPressed: () => Navigator.pop(context, true), child: const Text('Eliminar')),
                                        ],
                                      ),
                                    );
                                    if (ok == true) _delete(s['id']?.toString() ?? '');
                                  },
                                );
                              },
                            ),
                    ),
                  ],
                ),
    );
  }
}
