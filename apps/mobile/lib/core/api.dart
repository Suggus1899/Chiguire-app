import 'package:dio/dio.dart';

const _base = String.fromEnvironment(
  'API_URL',
  defaultValue: 'http://10.0.2.2:3001',
);

class _TokenStore {
  String? accessToken;
  String? refreshToken;
  String? tenantId;
}

final _tokens = _TokenStore();

final _dio = Dio(BaseOptions(baseUrl: _base))
  ..interceptors.add(InterceptorsWrapper(
    onRequest: (options, handler) {
      if (_tokens.accessToken != null) {
        options.headers['Authorization'] = 'Bearer ${_tokens.accessToken}';
      }
      handler.next(options);
    },
  ));

class AuthApi {
  Future<void> login(String email, String password, {String? tenantId}) async {
    final res = await _dio.post('/auth/login', data: {
      'email': email,
      'password': password,
      if (tenantId case final id?) 'tenant_id': id,
    });
    _tokens.accessToken = res.data['access_token'];
    _tokens.refreshToken = res.data['refresh_token'];
    _tokens.tenantId = tenantId;
  }

  Future<void> register(String email, String password, String fullName) async {
    await _dio.post('/auth/register', data: {
      'email': email,
      'password': password,
      'full_name': fullName,
    });
  }

  void logout() {
    _tokens.accessToken = null;
    _tokens.refreshToken = null;
    _tokens.tenantId = null;
  }

  String? get accessToken => _tokens.accessToken;
  String? get tenantId => _tokens.tenantId;
}

class TenantsApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/tenants');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(String name, String slug) async {
    final res = await _dio.post('/tenants', data: {'name': name, 'slug': slug});
    return Map<String, dynamic>.from(res.data);
  }
}

class PowerSyncApi {
  Future<String> getToken() async {
    final res = await _dio.get('/powersync/token');
    return res.data['token'] as String;
  }
}

class CustomersApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/customers');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/customers', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> get(String id) async {
    final res = await _dio.get('/customers/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> update(String id, Map<String, dynamic> data) async {
    final res = await _dio.put('/customers/$id', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<void> delete(String id) async {
    await _dio.delete('/customers/$id');
  }
}

class VendorsApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/vendors');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/vendors', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> get(String id) async {
    final res = await _dio.get('/vendors/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> update(String id, Map<String, dynamic> data) async {
    final res = await _dio.put('/vendors/$id', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<void> delete(String id) async {
    await _dio.delete('/vendors/$id');
  }
}

class ProductsApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/products');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/products', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> get(String id) async {
    final res = await _dio.get('/products/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> update(String id, Map<String, dynamic> data) async {
    final res = await _dio.put('/products/$id', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<void> delete(String id) async {
    await _dio.delete('/products/$id');
  }
}

class InventoryApi {
  Future<List<Map<String, dynamic>>> listMovements() async {
    final res = await _dio.get('/inventory/movements');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> createMovement(Map<String, dynamic> data) async {
    final res = await _dio.post('/inventory/movements', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> getStock(String productId) async {
    final res = await _dio.get('/inventory/stock/$productId');
    return Map<String, dynamic>.from(res.data);
  }

  Future<List<Map<String, dynamic>>> listWarehouses() async {
    final res = await _dio.get('/inventory/warehouses');
    return List<Map<String, dynamic>>.from(res.data);
  }
}

class InvoicesApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/invoices');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/invoices', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> get(String id) async {
    final res = await _dio.get('/invoices/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> emit(String id) async {
    final res = await _dio.post('/invoices/$id/emit');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> addPayment(String id, Map<String, dynamic> data) async {
    final res = await _dio.post('/invoices/$id/payments', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<void> voidInvoice(String id) async {
    await _dio.post('/invoices/$id/void');
  }
}

class FiscalApi {
  Future<Map<String, dynamic>> getExchangeRate(String currency) async {
    final res = await _dio.get('/fiscal/exchange-rate/$currency');
    return Map<String, dynamic>.from(res.data);
  }

  Future<List<Map<String, dynamic>>> listWithholdings() async {
    final res = await _dio.get('/fiscal/withholdings');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> generateFiscalBook(String period, String bookType) async {
    final res = await _dio.post('/fiscal/books', data: {
      'period': period,
      'book_type': bookType,
    });
    return Map<String, dynamic>.from(res.data);
  }
}

class PaymentsApi {
  Future<Map<String, dynamic>> createLink(Map<String, dynamic> data) async {
    final res = await _dio.post('/payments/links', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> getLink(String id) async {
    final res = await _dio.get('/payments/links/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<List<Map<String, dynamic>>> listLinks() async {
    final res = await _dio.get('/payments/links');
    return List<Map<String, dynamic>>.from(res.data);
  }
}

class PurchasesApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/purchases');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/purchases', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> approve(String id) async {
    final res = await _dio.post('/purchases/$id/approve');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> receiveItem(String id, String itemId, double quantity) async {
    final res = await _dio.post('/purchases/$id/receive', data: {
      'item_id': itemId,
      'quantity': quantity,
    });
    return Map<String, dynamic>.from(res.data);
  }
}

class QuotationsApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/quotations');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/quotations', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> convertToInvoice(String id) async {
    final res = await _dio.post('/quotations/$id/convert');
    return Map<String, dynamic>.from(res.data);
  }
}

class CommissionsApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/commissions');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> calculate(String invoiceId) async {
    final res = await _dio.post('/commissions/calculate', data: {
      'invoice_id': invoiceId,
    });
    return Map<String, dynamic>.from(res.data);
  }
}

class DeliveryApi {
  Future<List<Map<String, dynamic>>> listRoutes() async {
    final res = await _dio.get('/delivery/routes');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> createRoute(Map<String, dynamic> data) async {
    final res = await _dio.post('/delivery/routes', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> updateStopStatus(String routeId, String stopId, String status) async {
    final res = await _dio.patch('/delivery/routes/$routeId/stops/$stopId', data: {
      'status': status,
    });
    return Map<String, dynamic>.from(res.data);
  }
}

class SaasApi {
  Future<Map<String, dynamic>> getSubscription() async {
    final res = await _dio.get('/saas/subscription');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> getPlanLimits() async {
    final res = await _dio.get('/saas/plan-limits');
    return Map<String, dynamic>.from(res.data);
  }
}

class SellersApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/sellers');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/sellers', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> get(String id) async {
    final res = await _dio.get('/sellers/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> update(String id, Map<String, dynamic> data) async {
    final res = await _dio.put('/sellers/$id', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<void> delete(String id) async {
    await _dio.delete('/sellers/$id');
  }

  Future<List<Map<String, dynamic>>> listCommissions({Map<String, dynamic>? query}) async {
    final res = await _dio.get('/sellers/commissions', queryParameters: query);
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> markCommissionsPaid(Map<String, dynamic> data) async {
    final res = await _dio.post('/sellers/commissions/pay', data: data);
    return Map<String, dynamic>.from(res.data);
  }
}

class PaymentMethodsApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/payment-methods');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/payment-methods', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> update(String id, Map<String, dynamic> data) async {
    final res = await _dio.put('/payment-methods/$id', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<void> delete(String id) async {
    await _dio.delete('/payment-methods/$id');
  }
}

class FiscalDevicesApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/fiscal-devices');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/fiscal-devices', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> update(String id, Map<String, dynamic> data) async {
    final res = await _dio.put('/fiscal-devices/$id', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<void> delete(String id) async {
    await _dio.delete('/fiscal-devices/$id');
  }

  Future<List<Map<String, dynamic>>> listSequences(String deviceId) async {
    final res = await _dio.get('/fiscal-devices/$deviceId/sequences');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> createSequence(String deviceId, Map<String, dynamic> data) async {
    final res = await _dio.post('/fiscal-devices/$deviceId/sequences', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<List<Map<String, dynamic>>> listContingency(String deviceId) async {
    final res = await _dio.get('/fiscal-devices/$deviceId/contingency');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> createContingency(String deviceId, Map<String, dynamic> data) async {
    final res = await _dio.post('/fiscal-devices/$deviceId/contingency', data: data);
    return Map<String, dynamic>.from(res.data);
  }
}

class CreditNotesApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/credit-notes');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/credit-notes', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> get(String id) async {
    final res = await _dio.get('/credit-notes/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<void> voidNote(String id) async {
    await _dio.post('/credit-notes/$id/void');
  }
}

class TransfersApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/transfers');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/transfers', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> get(String id) async {
    final res = await _dio.get('/transfers/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> ship(String id) async {
    final res = await _dio.post('/transfers/$id/ship');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> receive(String id) async {
    final res = await _dio.post('/transfers/$id/receive');
    return Map<String, dynamic>.from(res.data);
  }
}

class ManufacturingApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/manufacturing');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/manufacturing', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> get(String id) async {
    final res = await _dio.get('/manufacturing/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> start(String id) async {
    final res = await _dio.post('/manufacturing/$id/start');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> complete(String id) async {
    final res = await _dio.post('/manufacturing/$id/complete');
    return Map<String, dynamic>.from(res.data);
  }
}

class PickingApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/picking');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/picking', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> get(String id) async {
    final res = await _dio.get('/picking/$id');
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> verifyItem(String id, String itemId, double qtyPicked) async {
    final res = await _dio.post('/picking/$id/items/$itemId/verify', data: {
      'qty_picked': qtyPicked,
    });
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> complete(String id) async {
    final res = await _dio.post('/picking/$id/complete');
    return Map<String, dynamic>.from(res.data);
  }
}

class AccountsPayableApi {
  Future<List<Map<String, dynamic>>> listReceivable() async {
    final res = await _dio.get('/accounts/receivable');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<List<Map<String, dynamic>>> listPayable() async {
    final res = await _dio.get('/accounts/payable');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> createPayment(Map<String, dynamic> data) async {
    final res = await _dio.post('/accounts/payments', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<List<Map<String, dynamic>>> listPayments() async {
    final res = await _dio.get('/accounts/payments');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> sendReminder(String id, String channel) async {
    final res = await _dio.post('/accounts/$id/remind', data: {'channel': channel});
    return Map<String, dynamic>.from(res.data);
  }
}

class ReportsApi {
  Future<Map<String, dynamic>> salesBook(Map<String, dynamic> params) async {
    final res = await _dio.get('/reports/sales-book', queryParameters: params);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> purchasesBook(Map<String, dynamic> params) async {
    final res = await _dio.get('/reports/purchases-book', queryParameters: params);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> inventoryCurrent(Map<String, dynamic> params) async {
    final res = await _dio.get('/reports/inventory-current', queryParameters: params);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> inventoryValued(Map<String, dynamic> params) async {
    final res = await _dio.get('/reports/inventory-valued', queryParameters: params);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> kardex(Map<String, dynamic> params) async {
    final res = await _dio.get('/reports/kardex', queryParameters: params);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> art177(Map<String, dynamic> params) async {
    final res = await _dio.get('/reports/art177', queryParameters: params);
    return Map<String, dynamic>.from(res.data);
  }

  Future<Map<String, dynamic>> igtfReport(Map<String, dynamic> params) async {
    final res = await _dio.get('/reports/igtf', queryParameters: params);
    return Map<String, dynamic>.from(res.data);
  }
}

class ApiTokensApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/api-tokens');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(Map<String, dynamic> data) async {
    final res = await _dio.post('/api-tokens', data: data);
    return Map<String, dynamic>.from(res.data);
  }

  Future<void> revoke(String id) async {
    await _dio.delete('/api-tokens/$id');
  }
}

final authApi = AuthApi();
final tenantsApi = TenantsApi();
final powerSyncApi = PowerSyncApi();
final customersApi = CustomersApi();
final vendorsApi = VendorsApi();
final productsApi = ProductsApi();
final inventoryApi = InventoryApi();
final invoicesApi = InvoicesApi();
final fiscalApi = FiscalApi();
final paymentsApi = PaymentsApi();
final purchasesApi = PurchasesApi();
final quotationsApi = QuotationsApi();
final commissionsApi = CommissionsApi();
final deliveryApi = DeliveryApi();
final saasApi = SaasApi();
final sellersApi = SellersApi();
final paymentMethodsApi = PaymentMethodsApi();
final fiscalDevicesApi = FiscalDevicesApi();
final creditNotesApi = CreditNotesApi();
final transfersApi = TransfersApi();
final manufacturingApi = ManufacturingApi();
final pickingApi = PickingApi();
final accountsPayableApi = AccountsPayableApi();
final reportsApi = ReportsApi();
final apiTokensApi = ApiTokensApi();
