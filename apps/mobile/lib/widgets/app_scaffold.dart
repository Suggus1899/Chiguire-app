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
    final routes = ['/', '/customers', '/products', '/invoices', '/settings'];
    context.go(routes[index]);
  }
}
