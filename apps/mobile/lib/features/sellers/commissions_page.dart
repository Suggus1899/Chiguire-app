import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';

class CommissionsPage extends StatefulWidget {
  const CommissionsPage({super.key});

  @override
  State<CommissionsPage> createState() => _CommissionsPageState();
}

class _CommissionsPageState extends State<CommissionsPage> {
  List<Map<String, dynamic>> _commissions = [];
  List<Map<String, dynamic>> _sellers = [];
  bool _isLoading = true;
  String? _error;
  String? _sellerFilter;
  String _statusFilter = 'all';
  final _selected = <int>{};

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
      _commissions = await sellersApi.listCommissions();
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  List<Map<String, dynamic>> get _filtered {
    return _commissions.where((c) {
      if (_sellerFilter != null && c['salesperson_id']?.toString() != _sellerFilter) return false;
      if (_statusFilter != 'all') {
        final status = c['status']?.toString() ?? 'pending';
        if (_statusFilter == 'paid' && status != 'paid') return false;
        if (_statusFilter == 'pending' && status == 'paid') return false;
      }
      return true;
    }).toList();
  }

  double get _selectedTotal {
    return _selected.fold(0.0, (sum, i) {
      if (i < _filtered.length) {
        return sum + ((filtered[i]['amount'] as num?)?.toDouble() ?? 0);
      }
      return sum;
    });
  }

  List<Map<String, dynamic>> get filtered => _filtered;

  String _sellerName(String? id) {
    final s = _sellers.where((e) => e['id']?.toString() == id?.toString()).toList();
    return s.isNotEmpty ? (s.first['name']?.toString() ?? 'N/A') : 'N/A';
  }

  Future<void> _markPaid() async {
    final refController = TextEditingController();
    final notesController = TextEditingController();
    final total = _selectedTotal;
    final ok = await showModalBottomSheet<bool>(
      context: context,
      builder: (ctx) => Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Marcar como Pagado', style: Theme.of(ctx).textTheme.titleMedium),
            const SizedBox(height: 12),
            Text('Total: ${total.toStringAsFixed(2)}'),
            const SizedBox(height: 12),
            TextField(
              controller: refController,
              decoration: const InputDecoration(labelText: 'Referencia', border: OutlineInputBorder()),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: notesController,
              decoration: const InputDecoration(labelText: 'Notas', border: OutlineInputBorder()),
              maxLines: 2,
            ),
            const SizedBox(height: 16),
            FilledButton(
              onPressed: () => Navigator.pop(ctx, true),
              child: const Text('Confirmar pago'),
            ),
          ],
        ),
      ),
    );
    if (ok != true) return;
    final ids = _selected.where((i) => i < _filtered.length).map((i) => _filtered[i]['id']?.toString()).toList();
    try {
      await sellersApi.markCommissionsPaid({
        'commission_ids': ids,
        'referencia': refController.text,
        'notas': notesController.text,
      });
      _selected.clear();
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
      title: 'Comisiones',
      currentIndex: 4,
      actions: [
        if (_selected.isNotEmpty)
          IconButton(
            icon: const Icon(Icons.payment),
            onPressed: _markPaid,
            tooltip: 'Marcar como Pagado',
          ),
      ],
      child: _isLoading
          ? const LoadingIndicator()
          : _error != null
              ? AppErrorWidget(message: _error!, onRetry: _load)
              : Column(
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(12),
                      child: Row(
                        children: [
                          Expanded(
                            child: DropdownButtonFormField<String>(
                              value: _sellerFilter,
                              decoration: const InputDecoration(labelText: 'Vendedor', border: OutlineInputBorder(), isDense: true),
                              items: [
                                const DropdownMenuItem(value: null, child: Text('Todos')),
                                ..._sellers.map((s) => DropdownMenuItem(
                                      value: s['id']?.toString(),
                                      child: Text(s['name']?.toString() ?? ''),
                                    )),
                              ],
                              onChanged: (v) => setState(() {
                                _sellerFilter = v;
                                _selected.clear();
                              }),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: DropdownButtonFormField<String>(
                              initialValue: _statusFilter,
                              decoration: const InputDecoration(labelText: 'Estado', border: OutlineInputBorder(), isDense: true),
                              items: const [
                                DropdownMenuItem(value: 'all', child: Text('Todos')),
                                DropdownMenuItem(value: 'pending', child: Text('Pendiente')),
                                DropdownMenuItem(value: 'paid', child: Text('Pagado')),
                              ],
                              onChanged: (v) => setState(() {
                                _statusFilter = v ?? 'all';
                                _selected.clear();
                              }),
                            ),
                          ),
                        ],
                      ),
                    ),
                    Expanded(
                      child: _filtered.isEmpty
                          ? const EmptyState(title: 'Sin comisiones', icon: Icons.percent)
                          : ListView.builder(
                              itemCount: _filtered.length,
                              itemBuilder: (context, i) {
                                final c = _filtered[i];
                                final paid = c['status']?.toString() == 'paid';
                                return ListTile(
                                  title: Text(_sellerName(c['salesperson_id']?.toString())),
                                  subtitle: Text(
                                    '${c['type']?.toString() ?? ''} · Base: ${c['base_amount'] ?? 0} · ${c['rate'] ?? 0}%',
                                  ),
                                  trailing: Row(
                                    mainAxisSize: MainAxisSize.min,
                                    children: [
                                      Column(
                                        mainAxisAlignment: MainAxisAlignment.center,
                                        crossAxisAlignment: CrossAxisAlignment.end,
                                        children: [
                                          Text(c['amount']?.toString() ?? '0'),
                                          Text(
                                            paid ? 'Pagado' : 'Pendiente',
                                            style: TextStyle(color: paid ? Colors.green : Colors.orange, fontSize: 12),
                                          ),
                                        ],
                                      ),
                                      Checkbox(
                                        value: _selected.contains(i),
                                        onChanged: paid
                                            ? null
                                            : (v) => setState(() {
                                                  if (v == true) {
                                                    _selected.add(i);
                                                  } else {
                                                    _selected.remove(i);
                                                  }
                                                }),
                                      ),
                                    ],
                                  ),
                                );
                              },
                            ),
                    ),
                  ],
                ),
    );
  }
}
