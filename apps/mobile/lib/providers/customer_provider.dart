import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../core/api.dart';
import '../core/database.dart';

class CustomerListState {
  final List<Map<String, dynamic>> customers;
  final bool isLoading;
  final String? error;

  const CustomerListState({
    this.customers = const [],
    this.isLoading = false,
    this.error,
  });

  CustomerListState copyWith({
    List<Map<String, dynamic>>? customers,
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
      final local = await db.getAll('SELECT * FROM customers ORDER BY name');
      final rows = local.map((r) => Map<String, dynamic>.from(r)).toList();
      state = CustomerListState(customers: rows);
      try {
        final remote = await customersApi.list();
        state = CustomerListState(customers: remote);
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
