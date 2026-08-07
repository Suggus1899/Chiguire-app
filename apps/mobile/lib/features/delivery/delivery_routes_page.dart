import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../core/database.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class DeliveryRoutesPage extends StatefulWidget {
  const DeliveryRoutesPage({super.key});

  @override
  State<DeliveryRoutesPage> createState() => _DeliveryRoutesPageState();
}

class _DeliveryRoutesPageState extends State<DeliveryRoutesPage> {
  List<Map<String, dynamic>> _routes = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final local = await db.getAll('SELECT * FROM delivery_routes ORDER BY date DESC');
      _routes = local.map((r) => Map<String, dynamic>.from(r)).toList();
      try {
        _routes = await deliveryApi.listRoutes();
      } catch (_) {}
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Rutas',
      currentIndex: 4,
      child: _isLoading
          ? const LoadingIndicator()
          : _routes.isEmpty
              ? const EmptyState(title: 'Sin rutas de entrega', icon: Icons.route_outlined)
              : ListView.builder(
                  itemCount: _routes.length,
                  itemBuilder: (context, i) {
                    final r = _routes[i];
                    return ListTile(
                      leading: const CircleAvatar(child: Icon(Icons.route)),
                      title: Text(r['driver_name']?.toString() ?? 'Conductor'),
                      subtitle: Text(r['date']?.toString() ?? ''),
                      trailing: Text(r['status']?.toString() ?? ''),
                    );
                  },
                ),
    );
  }
}
