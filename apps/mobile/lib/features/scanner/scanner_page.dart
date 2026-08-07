import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class ScannerPage extends StatefulWidget {
  const ScannerPage({super.key});

  @override
  State<ScannerPage> createState() => _ScannerPageState();
}

class _ScannerPageState extends State<ScannerPage> {
  final _barcodeController = TextEditingController();
  Map<String, dynamic>? _product;
  bool _isLoading = false;
  bool _notFound = false;

  @override
  void dispose() {
    _barcodeController.dispose();
    super.dispose();
  }

  Future<void> _search() async {
    final code = _barcodeController.text.trim();
    if (code.isEmpty) return;
    setState(() {
      _isLoading = true;
      _product = null;
      _notFound = false;
    });
    try {
      final products = await productsApi.list();
      final found = products.where((p) {
        final sku = p['sku']?.toString().toLowerCase() ?? '';
        final barcode = p['barcode']?.toString().toLowerCase() ?? '';
        return sku == code.toLowerCase() || barcode == code.toLowerCase();
      }).toList();
      if (found.isNotEmpty) {
        _product = found.first;
      } else {
        _notFound = true;
      }
    } catch (_) {
      _notFound = true;
    }
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Escáner',
      currentIndex: 4,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Card(
              child: Padding(
                padding: const EdgeInsets.all(20),
                child: Column(
                  children: [
                    Icon(Icons.qr_code_scanner, size: 80, color: Theme.of(context).colorScheme.outline),
                    const SizedBox(height: 12),
                    const Text('Escáner de código de barras', style: TextStyle(fontWeight: FontWeight.bold)),
                    const SizedBox(height: 8),
                    const Text('Ingresa el código manualmente o usa la cámara cuando esté disponible', textAlign: TextAlign.center),
                    const SizedBox(height: 16),
                    TextField(
                      controller: _barcodeController,
                      decoration: const InputDecoration(
                        labelText: 'Código de barras',
                        border: OutlineInputBorder(),
                        suffixIcon: Icon(Icons.search),
                      ),
                      onSubmitted: (_) => _search(),
                    ),
                    const SizedBox(height: 12),
                    FilledButton(onPressed: _search, child: const Text('Buscar')),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),
            if (_isLoading)
              const LoadingIndicator()
            else if (_notFound)
              const EmptyState(title: 'Producto no encontrado', icon: Icons.search_off)
            else if (_product != null)
              Card(
                child: ListTile(
                  leading: const Icon(Icons.inventory_2),
                  title: Text(_product!['name']?.toString() ?? ''),
                  subtitle: Text('SKU: ${_product!['sku']?.toString() ?? ''} · Precio: ${_product!['sale_price']?.toString() ?? '0'}'),
                  trailing: Text(_product!['cost_price']?.toString() ?? ''),
                ),
              ),
          ],
        ),
      ),
    );
  }
}
