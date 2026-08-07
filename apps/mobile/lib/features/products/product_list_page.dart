import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/error_widget.dart';
import '../../widgets/loading_indicator.dart';
import '../../providers/product_provider.dart';

class ProductListPage extends ConsumerStatefulWidget {
  const ProductListPage({super.key});

  @override
  ConsumerState<ProductListPage> createState() => _ProductListPageState();
}

class _ProductListPageState extends ConsumerState<ProductListPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(productListProvider.notifier).load());
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(productListProvider);
    return AppScaffold(
      title: 'Productos',
      currentIndex: 2,
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/products/new'),
        child: const Icon(Icons.add),
      ),
      child: state.isLoading
          ? const LoadingIndicator()
          : state.error != null && state.products.isEmpty
              ? AppErrorWidget(message: state.error!, onRetry: () => ref.read(productListProvider.notifier).load())
              : state.products.isEmpty
                  ? const EmptyState(
                      title: 'Sin productos',
                      subtitle: 'Crea tu primer producto',
                      icon: Icons.inventory_2_outlined,
                    )
                  : ListView.builder(
                      itemCount: state.products.length,
                      itemBuilder: (context, i) {
                        final p = state.products[i];
                        return ListTile(
                          leading: const CircleAvatar(child: Icon(Icons.inventory_2)),
                          title: Text(p['name']?.toString() ?? 'Sin nombre'),
                          subtitle: Text(p['sku']?.toString() ?? ''),
                          trailing: Text(p['sale_price']?.toString() ?? ''),
                        );
                      },
                    ),
    );
  }
}
