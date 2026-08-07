import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../core/database.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class QuotationListPage extends StatefulWidget {
  const QuotationListPage({super.key});

  @override
  State<QuotationListPage> createState() => _QuotationListPageState();
}

class _QuotationListPageState extends State<QuotationListPage> {
  List<Map<String, dynamic>> _quotations = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final local = await db.getAll('SELECT * FROM quotations ORDER BY created_at DESC');
      _quotations = local.map((r) => Map<String, dynamic>.from(r)).toList();
      try {
        _quotations = await quotationsApi.list();
      } catch (_) {}
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Cotizaciones',
      currentIndex: 4,
      child: _isLoading
          ? const LoadingIndicator()
          : _quotations.isEmpty
              ? const EmptyState(title: 'Sin cotizaciones', icon: Icons.description_outlined)
              : ListView.builder(
                  itemCount: _quotations.length,
                  itemBuilder: (context, i) {
                    final q = _quotations[i];
                    return ListTile(
                      leading: const CircleAvatar(child: Icon(Icons.description)),
                      title: Text(q['number']?.toString() ?? 'Sin número'),
                      subtitle: Text(q['status']?.toString() ?? ''),
                      trailing: Text(q['total']?.toString() ?? ''),
                    );
                  },
                ),
    );
  }
}
