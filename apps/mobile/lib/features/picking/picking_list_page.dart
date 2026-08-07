import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';

class PickingListPage extends StatefulWidget {
  const PickingListPage({super.key});

  @override
  State<PickingListPage> createState() => _PickingListPageState();
}

class _PickingListPageState extends State<PickingListPage> {
  List<Map<String, dynamic>> _lists = [];
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
      _lists = await pickingApi.list();
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Picking',
      currentIndex: 2,
      child: _isLoading
          ? const LoadingIndicator()
          : _error != null
              ? AppErrorWidget(message: _error!, onRetry: _load)
              : _lists.isEmpty
                  ? const EmptyState(title: 'Sin listas de picking', icon: Icons.checklist_outlined)
                  : ListView.builder(
                      itemCount: _lists.length,
                      itemBuilder: (context, i) {
                        final p = _lists[i];
                        final itemCount = (p['items'] as List?)?.length ?? p['item_count'] ?? 0;
                        return ListTile(
                          leading: const Icon(Icons.checklist),
                          title: Text('Almacén: ${p['warehouse_id']?.toString() ?? ''}'),
                          subtitle: Text('Ruta: ${p['route']?.toString() ?? ''} · Estado: ${p['status']?.toString() ?? ''} · Items: $itemCount'),
                          onTap: () => context.go('/picking/${p['id']}'),
                        );
                      },
                    ),
    );
  }
}
