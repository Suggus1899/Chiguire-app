import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../core/api.dart';
import '../core/database.dart';

class InvoiceListState {
  final List<Map<String, dynamic>> invoices;
  final bool isLoading;
  final String? error;

  const InvoiceListState({
    this.invoices = const [],
    this.isLoading = false,
    this.error,
  });

  InvoiceListState copyWith({
    List<Map<String, dynamic>>? invoices,
    bool? isLoading,
    String? error,
  }) {
    return InvoiceListState(
      invoices: invoices ?? this.invoices,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

class InvoiceListNotifier extends Notifier<InvoiceListState> {
  @override
  InvoiceListState build() => const InvoiceListState();

  Future<void> load() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final local = await db.getAll('SELECT * FROM invoices ORDER BY created_at DESC');
      final rows = local.map((r) => Map<String, dynamic>.from(r)).toList();
      state = InvoiceListState(invoices: rows);
      try {
        final remote = await invoicesApi.list();
        state = InvoiceListState(invoices: remote);
      } catch (_) {}
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> create(Map<String, dynamic> data) async {
    try {
      await invoicesApi.create(data);
      await load();
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }
}

final invoiceListProvider =
    NotifierProvider<InvoiceListNotifier, InvoiceListState>(() => InvoiceListNotifier());
