import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';

class ApiTokensPage extends StatefulWidget {
  const ApiTokensPage({super.key});

  @override
  State<ApiTokensPage> createState() => _ApiTokensPageState();
}

class _ApiTokensPageState extends State<ApiTokensPage> {
  List<Map<String, dynamic>> _tokens = [];
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
      _tokens = await apiTokensApi.list();
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _showCreateForm() async {
    final name = TextEditingController();
    final expires = TextEditingController();
    final ok = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      builder: (ctx) => Padding(
        padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom, left: 20, right: 20, top: 20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Nuevo token', style: Theme.of(ctx).textTheme.titleMedium),
            const SizedBox(height: 12),
            TextField(controller: name, decoration: const InputDecoration(labelText: 'Nombre', border: OutlineInputBorder())),
            const SizedBox(height: 12),
            TextField(controller: expires, decoration: const InputDecoration(labelText: 'Expira (YYYY-MM-DD)', border: OutlineInputBorder())),
            const SizedBox(height: 16),
            FilledButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('Crear')),
            const SizedBox(height: 20),
          ],
        ),
      ),
    );
    if (ok != true) return;
    try {
      final res = await apiTokensApi.create({'name': name.text, 'expires_at': expires.text});
      final token = res['token']?.toString() ?? '';
      _load();
      if (mounted) _showTokenDialog(token);
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  void _showTokenDialog(String token) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Token creado'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text('Copia tu token ahora. No se mostrará de nuevo.'),
            const SizedBox(height: 12),
            SelectableText(token, style: const TextStyle(fontFamily: 'monospace')),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () async {
              await Clipboard.setData(ClipboardData(text: token));
              if (ctx.mounted) ScaffoldMessenger.of(ctx).showSnackBar(const SnackBar(content: Text('Copiado')));
            },
            child: const Text('Copiar'),
          ),
          FilledButton(onPressed: () => Navigator.pop(ctx), child: const Text('Cerrar')),
        ],
      ),
    );
  }

  Future<void> _revoke(String id) async {
    try {
      await apiTokensApi.revoke(id);
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'API Tokens',
      currentIndex: 4,
      floatingActionButton: FloatingActionButton(
        onPressed: _showCreateForm,
        child: const Icon(Icons.add),
      ),
      child: _isLoading
          ? const LoadingIndicator()
          : _error != null
              ? AppErrorWidget(message: _error!, onRetry: _load)
              : _tokens.isEmpty
                  ? const EmptyState(title: 'Sin tokens', icon: Icons.key)
                  : ListView.builder(
                      itemCount: _tokens.length,
                      itemBuilder: (context, i) {
                        final t = _tokens[i];
                        final revoked = t['status']?.toString() == 'revoked';
                        return ListTile(
                          leading: const Icon(Icons.key),
                          title: Text(t['name']?.toString() ?? ''),
                          subtitle: Text('Último uso: ${t['last_used']?.toString() ?? 'N/A'} · Expira: ${t['expires_at']?.toString() ?? 'N/A'}'),
                          trailing: revoked
                              ? const Text('Revocado', style: TextStyle(color: Colors.grey))
                              : IconButton(
                                  icon: const Icon(Icons.delete_outline),
                                  onPressed: () async {
                                    final ok = await showDialog<bool>(
                                      context: context,
                                      builder: (_) => AlertDialog(
                                        title: const Text('Revocar token'),
                                        content: const Text('¿Revocar este token?'),
                                        actions: [
                                          TextButton(onPressed: () => Navigator.pop(context, false), child: const Text('Cancelar')),
                                          FilledButton(onPressed: () => Navigator.pop(context, true), child: const Text('Revocar')),
                                        ],
                                      ),
                                    );
                                    if (ok == true) _revoke(t['id']?.toString() ?? '');
                                  },
                                ),
                        );
                      },
                    ),
    );
  }
}
