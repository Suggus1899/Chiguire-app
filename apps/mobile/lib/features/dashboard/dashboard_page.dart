import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';
import '../../providers/invoice_provider.dart';
import '../../providers/customer_provider.dart';
import '../../providers/product_provider.dart';

class DashboardPage extends ConsumerStatefulWidget {
  const DashboardPage({super.key});

  @override
  ConsumerState<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends ConsumerState<DashboardPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() {
      ref.read(invoiceListProvider.notifier).load();
      ref.read(customerListProvider.notifier).load();
      ref.read(productListProvider.notifier).load();
    });
  }

  @override
  Widget build(BuildContext context) {
    final invoices = ref.watch(invoiceListProvider);
    final customers = ref.watch(customerListProvider);
    final products = ref.watch(productListProvider);

    return AppScaffold(
      title: 'Inicio',
      currentIndex: 0,
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          _StatCard(
            icon: Icons.receipt_long,
            label: 'Facturas',
            value: '${invoices.invoices.length}',
            onTap: () => context.go('/invoices'),
          ),
          const SizedBox(height: 12),
          _StatCard(
            icon: Icons.people,
            label: 'Clientes',
            value: '${customers.customers.length}',
            onTap: () => context.go('/customers'),
          ),
          const SizedBox(height: 12),
          _StatCard(
            icon: Icons.inventory_2,
            label: 'Productos',
            value: '${products.products.length}',
            onTap: () => context.go('/products'),
          ),
          const SizedBox(height: 24),
          Text('Facturas recientes', style: Theme.of(context).textTheme.titleMedium),
          const SizedBox(height: 8),
          if (invoices.isLoading)
            const LoadingIndicator()
          else if (invoices.invoices.isEmpty)
            const EmptyState(title: 'Sin facturas', icon: Icons.receipt_long_outlined)
          else
            ...invoices.invoices.take(5).map((inv) => ListTile(
                  leading: const Icon(Icons.receipt),
                  title: Text(inv['number']?.toString() ?? 'Sin número'),
                  subtitle: Text(inv['status']?.toString() ?? ''),
                  trailing: Text(inv['total']?.toString() ?? ''),
                  onTap: () => context.go('/invoices/${inv['id']}'),
                )),
        ],
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;
  final VoidCallback? onTap;

  const _StatCard({required this.icon, required this.label, required this.value, this.onTap});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        leading: Icon(icon, size: 32, color: Theme.of(context).colorScheme.primary),
        title: Text(value, style: Theme.of(context).textTheme.headlineSmall),
        subtitle: Text(label),
        onTap: onTap,
      ),
    );
  }
}
