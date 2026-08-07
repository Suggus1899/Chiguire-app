import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../core/api.dart';

class TenantNotifier extends Notifier<Map<String, dynamic>?> {
  @override
  Map<String, dynamic>? build() => null;

  Future<List<Map<String, dynamic>>> loadTenants() async {
    return tenantsApi.list();
  }

  void select(Map<String, dynamic> tenant) {
    state = tenant;
  }

  void clear() {
    state = null;
  }
}

final tenantProvider = NotifierProvider<TenantNotifier, Map<String, dynamic>?>(() => TenantNotifier());
