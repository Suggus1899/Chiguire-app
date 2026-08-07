import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/app_scaffold.dart';

class SettingsPage extends ConsumerWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final auth = ref.watch(authProvider);
    return AppScaffold(
      title: 'Más',
      currentIndex: 4,
      child: ListView(
        children: [
          if (auth.email != null)
            ListTile(leading: const Icon(Icons.person), title: Text(auth.email!)),
          const Divider(),
          ListTile(
            leading: const Icon(Icons.inventory),
            title: const Text('Inventario'),
            onTap: () => context.go('/inventory'),
          ),
          ListTile(
            leading: const Icon(Icons.description),
            title: const Text('Cotizaciones'),
            onTap: () => context.go('/quotations'),
          ),
          ListTile(
            leading: const Icon(Icons.shopping_cart),
            title: const Text('Compras'),
            onTap: () => context.go('/purchases'),
          ),
          ListTile(
            leading: const Icon(Icons.account_balance),
            title: const Text('Fiscal'),
            onTap: () => context.go('/fiscal'),
          ),
          ListTile(
            leading: const Icon(Icons.link),
            title: const Text('Pagos'),
            onTap: () => context.go('/payments'),
          ),
          ListTile(
            leading: const Icon(Icons.route),
            title: const Text('Rutas de entrega'),
            onTap: () => context.go('/routes'),
          ),
          const Divider(),
          ListTile(
            leading: const Icon(Icons.logout, color: Colors.red),
            title: const Text('Cerrar sesión', style: TextStyle(color: Colors.red)),
            onTap: () async {
              await ref.read(authProvider.notifier).logout();
              if (context.mounted) context.go('/login');
            },
          ),
          const SizedBox(height: 16),
          const Center(child: Text('Chiguire ERP v1.0.0', style: TextStyle(color: Colors.grey))),
        ],
      ),
    );
  }
}
