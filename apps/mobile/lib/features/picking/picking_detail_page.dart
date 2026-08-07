import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';

class PickingDetailPage extends StatefulWidget {
  final String id;

  const PickingDetailPage({super.key, required this.id});

  @override
  State<PickingDetailPage> createState() => _PickingDetailPageState();
}

class _PickingDetailPageState extends State<PickingDetailPage> {
  Map<String, dynamic>? _picking;
  List<Map<String, dynamic>> _items = [];
  bool _isLoading = true;
  String? _error;
  final _qtyControllers = <String, TextEditingController>{};

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
      _picking = await pickingApi.get(widget.id);
      _items = List<Map<String, dynamic>>.from(_picking?['items'] ?? []);
      for (final item in _items) {
        final id = item['id']?.toString() ?? '';
        _qtyControllers[id] = TextEditingController(text: item['qty_picked']?.toString() ?? '0');
      }
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _verify(String itemId) async {
    final qty = double.tryParse(_qtyControllers[itemId]?.text ?? '0') ?? 0;
    try {
      await pickingApi.verifyItem(widget.id, itemId, qty);
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  Future<void> _complete() async {
    try {
      await pickingApi.complete(widget.id);
      if (mounted) context.go('/picking');
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  @override
  void dispose() {
    for (final c in _qtyControllers.values) {
      c.dispose();
    }
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Detalle de picking')),
      body: _isLoading
          ? const LoadingIndicator()
          : _error != null
              ? AppErrorWidget(message: _error!, onRetry: _load)
              : _items.isEmpty
                  ? const EmptyState(title: 'Sin items', icon: Icons.checklist)
                  : ListView.builder(
                      padding: const EdgeInsets.all(16),
                      itemCount: _items.length,
                      itemBuilder: (context, i) {
                        final item = _items[i];
                        final id = item['id']?.toString() ?? '';
                        final verified = item['verified'] == true;
                        return Card(
                          child: Padding(
                            padding: const EdgeInsets.all(12),
                            child: Row(
                              children: [
                                Expanded(
                                  flex: 3,
                                  child: Column(
                                    crossAxisAlignment: CrossAxisAlignment.start,
                                    children: [
                                      Text(item['product_id']?.toString() ?? '', style: Theme.of(context).textTheme.titleSmall),
                                      Text('Esperado: ${item['quantity']?.toString() ?? '0'}'),
                                    ],
                                  ),
                                ),
                                const SizedBox(width: 8),
                                Expanded(
                                  child: TextField(
                                    controller: _qtyControllers[id],
                                    decoration: const InputDecoration(labelText: 'Pickeado', border: OutlineInputBorder(), isDense: true),
                                    keyboardType: const TextInputType.numberWithOptions(decimal: true),
                                  ),
                                ),
                                const SizedBox(width: 8),
                                Checkbox(value: verified, onChanged: verified ? null : (_) => _verify(id)),
                              ],
                            ),
                          ),
                        );
                      },
                    ),
      floatingActionButton: _isLoading || _items.isEmpty
          ? null
          : FloatingActionButton.extended(
              onPressed: _complete,
              icon: const Icon(Icons.check),
              label: const Text('Completar'),
            ),
    );
  }
}
