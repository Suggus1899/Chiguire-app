import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';

class CreditNoteListPage extends StatefulWidget {
  const CreditNoteListPage({super.key});

  @override
  State<CreditNoteListPage> createState() => _CreditNoteListPageState();
}

class _CreditNoteListPageState extends State<CreditNoteListPage> {
  List<Map<String, dynamic>> _notes = [];
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
      _notes = await creditNotesApi.list();
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _void(String id) async {
    try {
      await creditNotesApi.voidNote(id);
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Notas de Crédito',
      currentIndex: 4,
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/credit-notes/new'),
        child: const Icon(Icons.add),
      ),
      child: _isLoading
          ? const LoadingIndicator()
          : _error != null
              ? AppErrorWidget(message: _error!, onRetry: _load)
              : _notes.isEmpty
                  ? const EmptyState(title: 'Sin notas de crédito', icon: Icons.receipt_long_outlined)
                  : ListView.builder(
                      itemCount: _notes.length,
                      itemBuilder: (context, i) {
                        final n = _notes[i];
                        final status = n['status']?.toString() ?? '';
                        return ListTile(
                          leading: const Icon(Icons.receipt_long),
                          title: Text(n['number']?.toString() ?? 'N/A'),
                          subtitle: Text('${n['type']?.toString() ?? ''} · Factura: ${n['invoice_id']?.toString() ?? ''} · ${n['reason']?.toString() ?? ''}'),
                          trailing: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Text(n['total']?.toString() ?? '0'),
                              if (status != 'voided')
                                IconButton(
                                  icon: const Icon(Icons.delete_outline),
                                  onPressed: () async {
                                    final ok = await showDialog<bool>(
                                      context: context,
                                      builder: (_) => AlertDialog(
                                        title: const Text('Anular'),
                                        content: const Text('¿Anular esta nota?'),
                                        actions: [
                                          TextButton(onPressed: () => Navigator.pop(context, false), child: const Text('Cancelar')),
                                          FilledButton(onPressed: () => Navigator.pop(context, true), child: const Text('Anular')),
                                        ],
                                      ),
                                    );
                                    if (ok == true) _void(n['id']?.toString() ?? '');
                                  },
                                ),
                            ],
                          ),
                        );
                      },
                    ),
    );
  }
}
