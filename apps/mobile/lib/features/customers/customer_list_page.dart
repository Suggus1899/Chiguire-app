import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';
import '../../providers/customer_provider.dart';

class CustomerListPage extends ConsumerStatefulWidget {
  const CustomerListPage({super.key});

  @override
  ConsumerState<CustomerListPage> createState() => _CustomerListPageState();
}

class _CustomerListPageState extends ConsumerState<CustomerListPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(customerListProvider.notifier).load());
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(customerListProvider);
    return AppScaffold(
      title: 'Clientes',
      currentIndex: 1,
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/customers/new'),
        child: const Icon(Icons.add),
      ),
      child: state.isLoading
          ? const LoadingIndicator()
          : state.error != null && state.customers.isEmpty
              ? AppErrorWidget(message: state.error!, onRetry: () => ref.read(customerListProvider.notifier).load())
              : state.customers.isEmpty
                  ? const EmptyState(
                      title: 'Sin clientes',
                      subtitle: 'Crea tu primer cliente',
                      icon: Icons.people_outline,
                    )
                  : ListView.builder(
                      itemCount: state.customers.length,
                      itemBuilder: (context, i) {
                        final c = state.customers[i];
                        return ListTile(
                          leading: const CircleAvatar(child: Icon(Icons.person)),
                          title: Text(c.name.isEmpty ? 'Sin nombre' : c.name),
                          subtitle: Text(c.taxId ?? c.email ?? ''),
                          trailing: Text(c.phone ?? ''),
                        );
                      },
                    ),
    );
  }
}
