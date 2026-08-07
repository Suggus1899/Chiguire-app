import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class AppScaffold extends StatelessWidget {
  final Widget child;
  final String title;
  final int currentIndex;
  final List<Widget>? actions;
  final Widget? floatingActionButton;

  const AppScaffold({
    super.key,
    required this.child,
    required this.title,
    this.currentIndex = 0,
    this.actions,
    this.floatingActionButton,
  });

  static const _destinations = [
    (icon: Icons.dashboard_outlined, selectedIcon: Icons.dashboard, label: 'Inicio'),
    (icon: Icons.people_outline, selectedIcon: Icons.people, label: 'Clientes'),
    (icon: Icons.inventory_2_outlined, selectedIcon: Icons.inventory_2, label: 'Productos'),
    (icon: Icons.receipt_long_outlined, selectedIcon: Icons.receipt_long, label: 'Facturas'),
    (icon: Icons.more_horiz, selectedIcon: Icons.more_horiz, label: 'Más'),
  ];

  static const _moreMenuItems = [
    (icon: Icons.badge_outlined, label: 'Vendedores', route: '/sellers'),
    (icon: Icons.payments_outlined, label: 'Métodos de Pago', route: '/payment-methods'),
    (icon: Icons.print_outlined, label: 'Dispositivos Fiscales', route: '/fiscal-devices'),
    (icon: Icons.account_balance_outlined, label: 'Cuentas', route: '/accounts'),
    (icon: Icons.assessment_outlined, label: 'Informes', route: '/reports'),
    (icon: Icons.key_outlined, label: 'API Tokens', route: '/api-tokens'),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(title), actions: actions),
      body: child,
      floatingActionButton: floatingActionButton,
      bottomNavigationBar: NavigationBar(
        selectedIndex: currentIndex,
        onDestinationSelected: (i) => _navigate(context, i),
        destinations: [
          for (final d in _destinations)
            NavigationDestination(
              icon: Icon(d.icon),
              selectedIcon: Icon(d.selectedIcon),
              label: d.label,
            ),
        ],
      ),
    );
  }

  void _navigate(BuildContext context, int index) {
    if (index == 4) {
      _showMoreMenu(context);
      return;
    }
    final routes = ['/', '/customers', '/products', '/invoices', '/settings'];
    context.go(routes[index]);
  }

  void _showMoreMenu(BuildContext context) {
    showModalBottomSheet(
      context: context,
      builder: (ctx) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Text('Más', style: Theme.of(ctx).textTheme.titleMedium),
            ),
            for (final item in _moreMenuItems)
              ListTile(
                leading: Icon(item.icon),
                title: Text(item.label),
                onTap: () {
                  Navigator.pop(ctx);
                  context.go(item.route);
                },
              ),
            const SizedBox(height: 8),
          ],
        ),
      ),
    );
  }
}
