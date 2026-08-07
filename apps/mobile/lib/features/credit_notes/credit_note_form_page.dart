import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';

class CreditNoteFormPage extends StatefulWidget {
  const CreditNoteFormPage({super.key});

  @override
  State<CreditNoteFormPage> createState() => _CreditNoteFormPageState();
}

class _CreditNoteFormPageState extends State<CreditNoteFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _reasonController = TextEditingController();
  List<Map<String, dynamic>> _invoices = [];
  List<Map<String, dynamic>> _products = [];
  String? _invoiceId;
  String _type = 'credit';
  bool _isLoading = false;
  final _items = <Map<String, dynamic>>[];

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    try {
      _invoices = await invoicesApi.list();
      _products = await productsApi.list();
    } catch (_) {}
    if (mounted) setState(() {});
  }

  double get _total {
    return _items.fold(0.0, (sum, item) {
      final qty = (item['quantity'] as num?)?.toDouble() ?? 0;
      final price = (item['unit_price'] as num?)?.toDouble() ?? 0;
      return sum + (qty * price);
    });
  }

  void _addItem() {
    setState(() {
      _items.add({'product_id': _products.isNotEmpty ? _products.first['id']?.toString() : '', 'quantity': 1, 'unit_price': 0});
    });
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    if (_invoiceId == null) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Selecciona una factura')));
      return;
    }
    setState(() => _isLoading = true);
    try {
      await creditNotesApi.create({
        'invoice_id': _invoiceId,
        'type': _type,
        'reason': _reasonController.text,
        'items': _items,
      });
      if (mounted) context.go('/credit-notes');
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Nueva nota de crédito')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              DropdownButtonFormField<String>(
                value: _invoiceId,
                decoration: const InputDecoration(labelText: 'Factura *', border: OutlineInputBorder()),
                items: _invoices.map((inv) => DropdownMenuItem(
                      value: inv['id']?.toString(),
                      child: Text(inv['number']?.toString() ?? inv['id']?.toString() ?? ''),
                    )).toList(),
                onChanged: (v) => setState(() => _invoiceId = v),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: _type,
                decoration: const InputDecoration(labelText: 'Tipo', border: OutlineInputBorder()),
                items: const [
                  DropdownMenuItem(value: 'credit', child: Text('Nota de crédito')),
                  DropdownMenuItem(value: 'debit', child: Text('Nota de débito')),
                ],
                onChanged: (v) => setState(() => _type = v ?? 'credit'),
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _reasonController,
                decoration: const InputDecoration(labelText: 'Motivo', border: OutlineInputBorder()),
                maxLines: 2,
              ),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text('Items', style: Theme.of(context).textTheme.titleMedium),
                  TextButton.icon(onPressed: _addItem, icon: const Icon(Icons.add), label: const Text('Agregar')),
                ],
              ),
              ..._items.asMap().entries.map((entry) {
                final i = entry.key;
                final item = entry.value;
                return Padding(
                  padding: const EdgeInsets.only(bottom: 8),
                  child: Row(
                    children: [
                      Expanded(
                        flex: 3,
                        child: DropdownButton<String>(
                          value: item['product_id']?.toString(),
                          isExpanded: true,
                          hint: const Text('Producto'),
                          items: _products.map((p) => DropdownMenuItem(
                                value: p['id']?.toString(),
                                child: Text(p['name']?.toString() ?? ''),
                              )).toList(),
                          onChanged: (v) => setState(() => _items[i]['product_id'] = v),
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: TextFormField(
                          initialValue: item['quantity']?.toString() ?? '1',
                          decoration: const InputDecoration(labelText: 'Cant', border: OutlineInputBorder(), isDense: true),
                          keyboardType: const TextInputType.numberWithOptions(decimal: true),
                          onChanged: (v) => _items[i]['quantity'] = double.tryParse(v) ?? 0,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: TextFormField(
                          initialValue: item['unit_price']?.toString() ?? '0',
                          decoration: const InputDecoration(labelText: 'Precio', border: OutlineInputBorder(), isDense: true),
                          keyboardType: const TextInputType.numberWithOptions(decimal: true),
                          onChanged: (v) => _items[i]['unit_price'] = double.tryParse(v) ?? 0,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.delete),
                        onPressed: () => setState(() => _items.removeAt(i)),
                      ),
                    ],
                  ),
                );
              }),
              const SizedBox(height: 16),
              Card(
                child: ListTile(
                  title: const Text('Total'),
                  trailing: Text(_total.toStringAsFixed(2), style: Theme.of(context).textTheme.titleLarge),
                ),
              ),
              const SizedBox(height: 16),
              FilledButton(
                onPressed: _isLoading ? null : _submit,
                child: _isLoading
                    ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2))
                    : const Text('Guardar'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
