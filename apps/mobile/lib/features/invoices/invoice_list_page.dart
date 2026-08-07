import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';
import '../../providers/invoice_provider.dart';

class InvoiceListPage extends ConsumerStatefulWidget {
  const InvoiceListPage({super.key});

  @override
  ConsumerState<InvoiceListPage> createState() => _InvoiceListPageState();
}

class _InvoiceListPageState extends ConsumerState<InvoiceListPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(invoiceListProvider.notifier).load());
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(invoiceListProvider);
    return AppScaffold(
      title: 'Facturas',
      currentIndex: 3,
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/invoices/new'),
        child: const Icon(Icons.add),
      ),
      child: state.isLoading
          ? const LoadingIndicator()
          : state.error != null && state.invoices.isEmpty
              ? AppErrorWidget(message: state.error!, onRetry: () => ref.read(invoiceListProvider.notifier).load())
              : state.invoices.isEmpty
                  ? const EmptyState(
                      title: 'Sin facturas',
                      subtitle: 'Crea tu primera factura',
                      icon: Icons.receipt_long_outlined,
                    )
                  : ListView.builder(
                      itemCount: state.invoices.length,
                      itemBuilder: (context, i) {
                        final inv = state.invoices[i];
                        return ListTile(
                          leading: const CircleAvatar(child: Icon(Icons.receipt)),
                          title: Text(inv['number']?.toString() ?? 'Sin número'),
                          subtitle: Text(inv['status']?.toString() ?? ''),
                          trailing: Text(inv['total']?.toString() ?? ''),
                          onTap: () => context.go('/invoices/${inv['id']}'),
                        );
                      },
                    ),
    );
  }
}
