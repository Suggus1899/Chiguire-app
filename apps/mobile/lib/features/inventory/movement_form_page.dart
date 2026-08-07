import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';

class MovementFormPage extends StatefulWidget {
  const MovementFormPage({super.key});

  @override
  State<MovementFormPage> createState() => _MovementFormPageState();
}

class _MovementFormPageState extends State<MovementFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _productIdController = TextEditingController();
  final _warehouseIdController = TextEditingController();
  final _quantityController = TextEditingController();
  final _referenceController = TextEditingController();
  String _type = 'in';
  bool _isLoading = false;

  @override
  void dispose() {
    _productIdController.dispose();
    _warehouseIdController.dispose();
    _quantityController.dispose();
    _referenceController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isLoading = true);
    try {
      await inventoryApi.createMovement({
        'product_id': _productIdController.text,
        'warehouse_id': _warehouseIdController.text,
        'movement_type': _type,
        'quantity': double.tryParse(_quantityController.text) ?? 0,
        'reference': _referenceController.text,
      });
      if (mounted) context.go('/inventory');
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Nuevo movimiento')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              DropdownButtonFormField<String>(
                initialValue: _type,
                decoration: const InputDecoration(labelText: 'Tipo', border: OutlineInputBorder()),
                items: const [
                  DropdownMenuItem(value: 'in', child: Text('Entrada')),
                  DropdownMenuItem(value: 'out', child: Text('Salida')),
                ],
                onChanged: (v) => setState(() => _type = v ?? 'in'),
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _productIdController,
                decoration: const InputDecoration(labelText: 'ID producto *', border: OutlineInputBorder()),
                validator: (v) => v == null || v.isEmpty ? 'Requerido' : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _warehouseIdController,
                decoration: const InputDecoration(labelText: 'ID almacén', border: OutlineInputBorder()),
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _quantityController,
                decoration: const InputDecoration(labelText: 'Cantidad *', border: OutlineInputBorder()),
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                validator: (v) => v == null || v.isEmpty ? 'Requerido' : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _referenceController,
                decoration: const InputDecoration(labelText: 'Referencia', border: OutlineInputBorder()),
              ),
              const SizedBox(height: 24),
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
