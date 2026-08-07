import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../core/api.dart';
import '../core/database.dart';
import '../models/customer.dart';

class CustomerListState {
  final List<Customer> customers;
  final bool isLoading;
  final String? error;

  const CustomerListState({
    this.customers = const [],
    this.isLoading = false,
    this.error,
  });

  CustomerListState copyWith({
    List<Customer>? customers,
    bool? isLoading,
    String? error,
  }) {
    return CustomerListState(
      customers: customers ?? this.customers,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

class CustomerListNotifier extends Notifier<CustomerListState> {
  @override
  CustomerListState build() => const CustomerListState();

  Future<void> load() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final tenantId = authApi.tenantId;
      final local = tenantId != null
          ? await db.getAll(
              'SELECT * FROM customers WHERE tenant_id = ? ORDER BY name',
              [tenantId],
            )
          : await db.getAll('SELECT * FROM customers ORDER BY name');
      final rows = local
          .map((r) => Customer.fromJson(Map<String, dynamic>.from(r)))
          .toList();
      state = CustomerListState(customers: rows);
      try {
        final remote = await customersApi.list();
        state = CustomerListState(
          customers: remote.map(Customer.fromJson).toList(),
        );
      } catch (_) {}
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> create(Map<String, dynamic> data) async {
    try {
      await customersApi.create(data);
      await load();
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }
}

final customerListProvider =
    NotifierProvider<CustomerListNotifier, CustomerListState>(() => CustomerListNotifier());
