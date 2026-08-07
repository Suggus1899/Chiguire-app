import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class ReportsPage extends StatelessWidget {
  const ReportsPage({super.key});

  static const _reports = [
    (key: 'sales_book', title: 'Libro de Ventas', icon: Icons.menu_book),
    (key: 'purchases_book', title: 'Libro de Compras', icon: Icons.book),
    (key: 'inventory_current', title: 'Inventario Actual', icon: Icons.inventory_2),
    (key: 'inventory_valued', title: 'Inventario Valorizado', icon: Icons.monetization_on),
    (key: 'kardex', title: 'Kardex', icon: Icons.list_alt),
    (key: 'art177', title: 'Art. 177', icon: Icons.gavel),
    (key: 'igtf', title: 'IGTF', icon: Icons.percent),
  ];

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Informes',
      currentIndex: 4,
      child: GridView.count(
        crossAxisCount: 2,
        padding: const EdgeInsets.all(16),
        crossAxisSpacing: 12,
        mainAxisSpacing: 12,
        children: _reports.map((r) {
          return Card(
            clipBehavior: Clip.antiAlias,
            child: InkWell(
              onTap: () => Navigator.push(
                context,
                MaterialPageRoute(builder: (_) => _ReportDetailPage(reportKey: r.key, title: r.title)),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(r.icon, size: 48),
                  const SizedBox(height: 8),
                  Text(r.title, textAlign: TextAlign.center),
                ],
              ),
            ),
          );
        }).toList(),
      ),
    );
  }
}

class _ReportDetailPage extends StatefulWidget {
  final String reportKey;
  final String title;

  const _ReportDetailPage({required this.reportKey, required this.title});

  @override
  State<_ReportDetailPage> createState() => _ReportDetailPageState();
}

class _ReportDetailPageState extends State<_ReportDetailPage> {
  final _fromController = TextEditingController();
  final _toController = TextEditingController();
  Map<String, dynamic>? _result;
  bool _isLoading = false;

  @override
  void dispose() {
    _fromController.dispose();
    _toController.dispose();
    super.dispose();
  }

  Future<void> _generate() async {
    setState(() => _isLoading = true);
    try {
      final params = {'from': _fromController.text, 'to': _toController.text};
      final fn = switch (widget.reportKey) {
        'sales_book' => reportsApi.salesBook,
        'purchases_book' => reportsApi.purchasesBook,
        'inventory_current' => reportsApi.inventoryCurrent,
        'inventory_valued' => reportsApi.inventoryValued,
        'kardex' => reportsApi.kardex,
        'art177' => reportsApi.art177,
        'igtf' => reportsApi.igtfReport,
        _ => reportsApi.salesBook,
      };
      _result = await fn(params);
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.title)),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Row(
              children: [
                Expanded(child: TextField(controller: _fromController, decoration: const InputDecoration(labelText: 'Desde', border: OutlineInputBorder(), hintText: 'YYYY-MM-DD'))),
                const SizedBox(width: 12),
                Expanded(child: TextField(controller: _toController, decoration: const InputDecoration(labelText: 'Hasta', border: OutlineInputBorder(), hintText: 'YYYY-MM-DD'))),
              ],
            ),
            const SizedBox(height: 12),
            FilledButton(
              onPressed: _isLoading ? null : _generate,
              child: _isLoading
                  ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2))
                  : const Text('Generar'),
            ),
            const SizedBox(height: 16),
            Expanded(
              child: _isLoading
                  ? const LoadingIndicator()
                  : _result == null
                      ? const EmptyState(title: 'Genera el informe', icon: Icons.assessment)
                      : SingleChildScrollView(
                          scrollDirection: Axis.vertical,
                          child: SingleChildScrollView(
                            scrollDirection: Axis.horizontal,
                            child: _buildTable(_result!),
                          ),
                        ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTable(Map<String, dynamic> data) {
    final rows = data['rows'] as List? ?? [];
    if (rows.isEmpty) return const Text('Sin datos');
    final firstRow = Map<String, dynamic>.from(rows.first as Map);
    final columns = firstRow.keys.toList();
    return DataTable(
      columns: columns.map((c) => DataColumn(label: Text(c.toString()))).toList(),
      rows: rows.map((r) {
        final row = Map<String, dynamic>.from(r as Map);
        return DataRow(
          cells: columns.map((c) => DataCell(Text(row[c]?.toString() ?? ''))).toList(),
        );
      }).toList(),
    );
  }
}
