import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';

class TransferFormPage extends StatefulWidget {
  const TransferFormPage({super.key});

  @override
  State<TransferFormPage> createState() => _TransferFormPageState();
}

class _TransferFormPageState extends State<TransferFormPage> {
  final _formKey = GlobalKey<FormState>();
  List<Map<String, dynamic>> _warehouses = [];
  List<Map<String, dynamic>> _products = [];
  String? _fromWarehouse;
  String? _toWarehouse;
  bool _isLoading = false;
  final _items = <Map<String, dynamic>>[];

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    try {
      _warehouses = await inventoryApi.listWarehouses();
      _products = await productsApi.list();
    } catch (_) {}
    if (mounted) setState(() {});
  }

  void _addItem() {
    setState(() {
      _items.add({'product_id': _products.isNotEmpty ? _products.first['id']?.toString() : '', 'quantity': 1});
    });
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    if (_fromWarehouse == null || _toWarehouse == null) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Selecciona almacenes')));
      return;
    }
    if (_fromWarehouse == _toWarehouse) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Los almacenes deben ser diferentes')));
      return;
    }
    setState(() => _isLoading = true);
    try {
      await transfersApi.create({
        'from_warehouse_id': _fromWarehouse,
        'to_warehouse_id': _toWarehouse,
        'items': _items,
      });
      if (mounted) context.go('/transfers');
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Nueva transferencia')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              DropdownButtonFormField<String>(
                value: _fromWarehouse,
                decoration: const InputDecoration(labelText: 'Almacén origen *', border: OutlineInputBorder()),
                items: _warehouses.map((w) => DropdownMenuItem(
                      value: w['id']?.toString(),
                      child: Text(w['name']?.toString() ?? ''),
                    )).toList(),
                onChanged: (v) => setState(() => _fromWarehouse = v),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: _toWarehouse,
                decoration: const InputDecoration(labelText: 'Almacén destino *', border: OutlineInputBorder()),
                items: _warehouses.map((w) => DropdownMenuItem(
                      value: w['id']?.toString(),
                      child: Text(w['name']?.toString() ?? ''),
                    )).toList(),
                onChanged: (v) => setState(() => _toWarehouse = v),
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
                return Padding(
                  padding: const EdgeInsets.only(bottom: 8),
                  child: Row(
                    children: [
                      Expanded(
                        flex: 3,
                        child: DropdownButton<String>(
                          value: _items[i]['product_id']?.toString(),
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
                          initialValue: _items[i]['quantity']?.toString() ?? '1',
                          decoration: const InputDecoration(labelText: 'Cant', border: OutlineInputBorder(), isDense: true),
                          keyboardType: const TextInputType.numberWithOptions(decimal: true),
                          onChanged: (v) => _items[i]['quantity'] = double.tryParse(v) ?? 0,
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
