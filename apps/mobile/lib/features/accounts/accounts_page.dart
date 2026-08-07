import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class AccountsPage extends StatefulWidget {
  const AccountsPage({super.key});

  @override
  State<AccountsPage> createState() => _AccountsPageState();
}

class _AccountsPageState extends State<AccountsPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Cuentas',
      currentIndex: 4,
      child: Column(
        children: [
          TabBar(
            controller: _tabController,
            tabs: const [Tab(text: 'Por Cobrar'), Tab(text: 'Por Pagar')],
          ),
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                _AccountsList(load: accountsPayableApi.listReceivable, type: 'receivable'),
                _AccountsList(load: accountsPayableApi.listPayable, type: 'payable'),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _AccountsList extends StatefulWidget {
  final Future<List<Map<String, dynamic>>> Function() load;
  final String type;

  const _AccountsList({required this.load, required this.type});

  @override
  State<_AccountsList> createState() => _AccountsListState();
}

class _AccountsListState extends State<_AccountsList> with AutomaticKeepAliveClientMixin {
  List<Map<String, dynamic>> _items = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _isLoading = true);
    try {
      _items = await widget.load();
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _addPayment(Map<String, dynamic> account) async {
    final amount = TextEditingController(text: account['balance']?.toString() ?? '0');
    final method = TextEditingController();
    final reference = TextEditingController();
    final notes = TextEditingController();
    final ok = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      builder: (ctx) => Padding(
        padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom, left: 20, right: 20, top: 20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Agregar pago', style: Theme.of(ctx).textTheme.titleMedium),
            const SizedBox(height: 12),
            TextField(controller: amount, decoration: const InputDecoration(labelText: 'Monto', border: OutlineInputBorder()), keyboardType: const TextInputType.numberWithOptions(decimal: true)),
            const SizedBox(height: 12),
            TextField(controller: method, decoration: const InputDecoration(labelText: 'Método', border: OutlineInputBorder())),
            const SizedBox(height: 12),
            TextField(controller: reference, decoration: const InputDecoration(labelText: 'Referencia', border: OutlineInputBorder())),
            const SizedBox(height: 12),
            TextField(controller: notes, decoration: const InputDecoration(labelText: 'Notas', border: OutlineInputBorder()), maxLines: 2),
            const SizedBox(height: 16),
            FilledButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('Guardar pago')),
            const SizedBox(height: 20),
          ],
        ),
      ),
    );
    if (ok != true) return;
    try {
      await accountsPayableApi.createPayment({
        'account_id': account['id'],
        'amount': double.tryParse(amount.text) ?? 0,
        'method': method.text,
        'reference': reference.text,
        'notes': notes.text,
      });
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  Future<void> _sendReminder(Map<String, dynamic> account) async {
    final channel = await showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Enviar recordatorio'),
        content: const Text('Elige el canal'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, 'sms'), child: const Text('SMS')),
          TextButton(onPressed: () => Navigator.pop(ctx, 'email'), child: const Text('Email')),
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Cancelar')),
        ],
      ),
    );
    if (channel == null) return;
    try {
      await accountsPayableApi.sendReminder(account['id']?.toString() ?? '', channel);
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Recordatorio enviado')));
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    if (_isLoading) return const LoadingIndicator();
    if (_items.isEmpty) return const EmptyState(title: 'Sin cuentas', icon: Icons.account_balance_outlined);
    return ListView.builder(
      itemCount: _items.length,
      itemBuilder: (context, i) {
        final a = _items[i];
        return ListTile(
          leading: const Icon(Icons.account_balance),
          title: Text(a['name']?.toString() ?? ''),
          subtitle: Text('Monto: ${a['amount']?.toString() ?? '0'} · Balance: ${a['balance']?.toString() ?? '0'} · Vence: ${a['due_date']?.toString() ?? ''}'),
          trailing: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              IconButton(icon: const Icon(Icons.notifications), onPressed: () => _sendReminder(a), tooltip: 'Recordatorio'),
              IconButton(icon: const Icon(Icons.payments), onPressed: () => _addPayment(a), tooltip: 'Pago'),
            ],
          ),
        );
      },
    );
  }

  @override
  bool get wantKeepAlive => true;
}
