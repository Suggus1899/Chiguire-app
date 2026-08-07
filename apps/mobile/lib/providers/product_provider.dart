import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../core/api.dart';
import '../core/database.dart';

class ProductListState {
  final List<Map<String, dynamic>> products;
  final bool isLoading;
  final String? error;

  const ProductListState({
    this.products = const [],
    this.isLoading = false,
    this.error,
  });

  ProductListState copyWith({
    List<Map<String, dynamic>>? products,
    bool? isLoading,
    String? error,
  }) {
    return ProductListState(
      products: products ?? this.products,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

class ProductListNotifier extends Notifier<ProductListState> {
  @override
  ProductListState build() => const ProductListState();

  Future<void> load() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final local = await db.getAll('SELECT * FROM products ORDER BY name');
      final rows = local.map((r) => Map<String, dynamic>.from(r)).toList();
      state = ProductListState(products: rows);
      try {
        final remote = await productsApi.list();
        state = ProductListState(products: remote);
      } catch (_) {}
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> create(Map<String, dynamic> data) async {
    try {
      await productsApi.create(data);
      await load();
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }
}

final productListProvider =
    NotifierProvider<ProductListNotifier, ProductListState>(() => ProductListNotifier());
