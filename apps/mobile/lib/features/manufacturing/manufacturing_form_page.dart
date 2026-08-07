import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';

class ManufacturingFormPage extends StatefulWidget {
  const ManufacturingFormPage({super.key});

  @override
  State<ManufacturingFormPage> createState() => _ManufacturingFormPageState();
}

class _ManufacturingFormPageState extends State<ManufacturingFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _qtyController = TextEditingController(text: '1');
  List<Map<String, dynamic>> _products = [];
  List<Map<String, dynamic>> _warehouses = [];
  String? _product;
  String? _warehouse;
  bool _isLoading = false;
  final _bomItems = <Map<String, dynamic>>[];

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    try {
      _products = await productsApi.list();
      _warehouses = await inventoryApi.listWarehouses();
    } catch (_) {}
    if (mounted) setState(() {});
  }

  void _addBomItem() {
    setState(() {
      _bomItems.add({'component_id': _products.isNotEmpty ? _products.first['id']?.toString() : '', 'qty_per_unit': 1});
    });
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    if (_product == null || _warehouse == null) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Selecciona producto y almacén')));
      return;
    }
    setState(() => _isLoading = true);
    try {
      await manufacturingApi.create({
        'product_id': _product,
        'warehouse_id': _warehouse,
        'quantity': double.tryParse(_qtyController.text) ?? 1,
        'bom': _bomItems,
      });
      if (mounted) context.go('/manufacturing');
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  void dispose() {
    _qtyController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Nueva orden de producción')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              DropdownButtonFormField<String>(
                value: _product,
                decoration: const InputDecoration(labelText: 'Producto *', border: OutlineInputBorder()),
                items: _products.map((p) => DropdownMenuItem(
                      value: p['id']?.toString(),
                      child: Text(p['name']?.toString() ?? ''),
                    )).toList(),
                onChanged: (v) => setState(() => _product = v),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: _warehouse,
                decoration: const InputDecoration(labelText: 'Almacén *', border: OutlineInputBorder()),
                items: _warehouses.map((w) => DropdownMenuItem(
                      value: w['id']?.toString(),
                      child: Text(w['name']?.toString() ?? ''),
                    )).toList(),
                onChanged: (v) => setState(() => _warehouse = v),
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _qtyController,
                decoration: const InputDecoration(labelText: 'Cantidad', border: OutlineInputBorder()),
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                validator: (v) => v == null || v.isEmpty ? 'Requerido' : null,
              ),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text('Lista de materiales (BOM)', style: Theme.of(context).textTheme.titleMedium),
                  TextButton.icon(onPressed: _addBomItem, icon: const Icon(Icons.add), label: const Text('Agregar')),
                ],
              ),
              ..._bomItems.asMap().entries.map((entry) {
                final i = entry.key;
                return Padding(
                  padding: const EdgeInsets.only(bottom: 8),
                  child: Row(
                    children: [
                      Expanded(
                        flex: 3,
                        child: DropdownButton<String>(
                          value: _bomItems[i]['component_id']?.toString(),
                          isExpanded: true,
                          hint: const Text('Componente'),
                          items: _products.map((p) => DropdownMenuItem(
                                value: p['id']?.toString(),
                                child: Text(p['name']?.toString() ?? ''),
                              )).toList(),
                          onChanged: (v) => setState(() => _bomItems[i]['component_id'] = v),
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: TextFormField(
                          initialValue: _bomItems[i]['qty_per_unit']?.toString() ?? '1',
                          decoration: const InputDecoration(labelText: 'Cant/U', border: OutlineInputBorder(), isDense: true),
                          keyboardType: const TextInputType.numberWithOptions(decimal: true),
                          onChanged: (v) => _bomItems[i]['qty_per_unit'] = double.tryParse(v) ?? 0,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.delete),
                        onPressed: () => setState(() => _bomItems.removeAt(i)),
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
