import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';

class PaymentMethodsPage extends StatefulWidget {
  const PaymentMethodsPage({super.key});

  @override
  State<PaymentMethodsPage> createState() => _PaymentMethodsPageState();
}

class _PaymentMethodsPageState extends State<PaymentMethodsPage> {
  List<Map<String, dynamic>> _methods = [];
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
      _methods = await paymentMethodsApi.list();
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _showForm([Map<String, dynamic>? existing]) async {
    final name = TextEditingController(text: existing?['name']?.toString() ?? '');
    String type = existing?['type']?.toString() ?? 'cash';
    String currency = existing?['currency']?.toString() ?? 'VES';
    bool forInvoices = existing?['for_invoices'] == true;
    bool forChange = existing?['for_change'] == true;
    bool forRefunds = existing?['for_refunds'] == true;
    final ok = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setS) => Padding(
          padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom, left: 20, right: 20, top: 20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(existing == null ? 'Nuevo método de pago' : 'Editar método', style: Theme.of(ctx).textTheme.titleMedium),
              const SizedBox(height: 12),
              TextField(
                controller: name,
                decoration: const InputDecoration(labelText: 'Nombre', border: OutlineInputBorder()),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: type,
                decoration: const InputDecoration(labelText: 'Tipo', border: OutlineInputBorder()),
                items: const [
                  DropdownMenuItem(value: 'cash', child: Text('Efectivo')),
                  DropdownMenuItem(value: 'card', child: Text('Tarjeta')),
                  DropdownMenuItem(value: 'transfer', child: Text('Transferencia')),
                  DropdownMenuItem(value: 'mobile_payment', child: Text('Pago móvil')),
                  DropdownMenuItem(value: 'crypto', child: Text('Cripto')),
                ],
                onChanged: (v) => setS(() => type = v ?? 'cash'),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: currency,
                decoration: const InputDecoration(labelText: 'Moneda', border: OutlineInputBorder()),
                items: const [
                  DropdownMenuItem(value: 'VES', child: Text('VES')),
                  DropdownMenuItem(value: 'USD', child: Text('USD')),
                  DropdownMenuItem(value: 'EUR', child: Text('EUR')),
                ],
                onChanged: (v) => setS(() => currency = v ?? 'VES'),
              ),
              const SizedBox(height: 8),
              SwitchListTile(title: const Text('Para facturas'), value: forInvoices, onChanged: (v) => setS(() => forInvoices = v)),
              SwitchListTile(title: const Text('Para cambio'), value: forChange, onChanged: (v) => setS(() => forChange = v)),
              SwitchListTile(title: const Text('Para reembolsos'), value: forRefunds, onChanged: (v) => setS(() => forRefunds = v)),
              const SizedBox(height: 12),
              FilledButton(
                onPressed: () => Navigator.pop(ctx, true),
                child: const Text('Guardar'),
              ),
              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
    if (ok != true) return;
    final data = {
      'name': name.text,
      'type': type,
      'currency': currency,
      'for_invoices': forInvoices,
      'for_change': forChange,
      'for_refunds': forRefunds,
    };
    try {
      if (existing != null) {
        await paymentMethodsApi.update(existing['id']?.toString() ?? '', data);
      } else {
        await paymentMethodsApi.create(data);
      }
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
      title: 'Métodos de Pago',
      currentIndex: 4,
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showForm(),
        child: const Icon(Icons.add),
      ),
      child: _isLoading
          ? const LoadingIndicator()
          : _error != null
              ? AppErrorWidget(message: _error!, onRetry: _load)
              : _methods.isEmpty
                  ? const EmptyState(title: 'Sin métodos de pago', icon: Icons.payments_outlined)
                  : ListView.builder(
                      itemCount: _methods.length,
                      itemBuilder: (context, i) {
                        final m = _methods[i];
                        return ListTile(
                          leading: const Icon(Icons.payments),
                          title: Text(m['name']?.toString() ?? ''),
                          subtitle: Text('${m['type']?.toString() ?? ''} · ${m['currency']?.toString() ?? ''}'),
                          trailing: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              if (m['for_invoices'] == true) const Icon(Icons.receipt, size: 16),
                              if (m['for_change'] == true) const Icon(Icons.swap_horiz, size: 16),
                              if (m['for_refunds'] == true) const Icon(Icons.undo, size: 16),
                            ],
                          ),
                          onTap: () => _showForm(m),
                        );
                      },
                    ),
    );
  }
}
