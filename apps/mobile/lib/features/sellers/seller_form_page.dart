import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../core/api.dart';

class SellerFormPage extends StatefulWidget {
  final String? id;

  const SellerFormPage({super.key, this.id});

  @override
  State<SellerFormPage> createState() => _SellerFormPageState();
}

class _SellerFormPageState extends State<SellerFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  final _commissionController = TextEditingController(text: '0');
  bool _isActive = true;
  bool _isLoading = false;
  bool _isEdit = false;

  @override
  void initState() {
    super.initState();
    if (widget.id != null) {
      _isEdit = true;
      _load();
    }
  }

  Future<void> _load() async {
    try {
      final s = await sellersApi.get(widget.id!);
      _nameController.text = s['name']?.toString() ?? '';
      _emailController.text = s['email']?.toString() ?? '';
      _phoneController.text = s['phone']?.toString() ?? '';
      _commissionController.text = s['commission_pct']?.toString() ?? '0';
      _isActive = s['is_active'] == true || s['is_active'] == 1;
      if (mounted) setState(() {});
    } catch (_) {}
  }

  @override
  void dispose() {
    _nameController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _commissionController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isLoading = true);
    try {
      final data = {
        'name': _nameController.text,
        'email': _emailController.text,
        'phone': _phoneController.text,
        'commission_pct': double.tryParse(_commissionController.text) ?? 0,
        'is_active': _isActive,
      };
      if (_isEdit) {
        await sellersApi.update(widget.id!, data);
      } else {
        await sellersApi.create(data);
      }
      if (mounted) context.go('/sellers');
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
      appBar: AppBar(title: Text(_isEdit ? 'Editar vendedor' : 'Nuevo vendedor')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(labelText: 'Nombre *', border: OutlineInputBorder()),
                validator: (v) => v == null || v.isEmpty ? 'Requerido' : null,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _emailController,
                decoration: const InputDecoration(labelText: 'Correo', border: OutlineInputBorder()),
                keyboardType: TextInputType.emailAddress,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _phoneController,
                decoration: const InputDecoration(labelText: 'Teléfono', border: OutlineInputBorder()),
                keyboardType: TextInputType.phone,
              ),
              const SizedBox(height: 12),
              TextFormField(
                controller: _commissionController,
                decoration: const InputDecoration(labelText: '% Comisión', border: OutlineInputBorder(), suffixText: '%'),
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
              ),
              const SizedBox(height: 12),
              SwitchListTile(
                title: const Text('Activo'),
                value: _isActive,
                onChanged: (v) => setState(() => _isActive = v),
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
