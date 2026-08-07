import 'package:flutter/material.dart';
import '../../core/api.dart';

class InvoiceDetailPage extends StatefulWidget {
  final String id;
  const InvoiceDetailPage({super.key, required this.id});

  @override
  State<InvoiceDetailPage> createState() => _InvoiceDetailPageState();
}

class _InvoiceDetailPageState extends State<InvoiceDetailPage> {
  Map<String, dynamic>? _invoice;
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      _invoice = await invoicesApi.get(widget.id);
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _emit() async {
    try {
      await invoicesApi.emit(widget.id);
      _load();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  Future<void> _void() async {
    try {
      await invoicesApi.voidInvoice(widget.id);
      _load();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Factura')),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? Center(child: Text(_error!))
              : _invoice == null
                  ? const Center(child: Text('No encontrada'))
                  : Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('N° ${_invoice!['number'] ?? ''}',
                              style: Theme.of(context).textTheme.headlineSmall),
                          const SizedBox(height: 8),
                          Text('Estado: ${_invoice!['status'] ?? ''}'),
                          Text('Moneda: ${_invoice!['currency'] ?? ''}'),
                          Text('Subtotal: ${_invoice!['subtotal'] ?? ''}'),
                          Text('Impuesto: ${_invoice!['tax'] ?? ''}'),
                          Text('Total: ${_invoice!['total'] ?? ''}'),
                          const Spacer(),
                          Row(
                            children: [
                              FilledButton(onPressed: _emit, child: const Text('Emitir')),
                              const SizedBox(width: 12),
                              OutlinedButton(onPressed: _void, child: const Text('Anular')),
                            ],
                          ),
                        ],
                      ),
                    ),
    );
  }
}
