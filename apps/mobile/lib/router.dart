import 'package:go_router/go_router.dart';
import 'core/api.dart';
import 'features/auth/login_page.dart';
import 'features/auth/register_page.dart';
import 'features/dashboard/dashboard_page.dart';
import 'features/customers/customer_list_page.dart';
import 'features/customers/customer_form_page.dart';
import 'features/products/product_list_page.dart';
import 'features/products/product_form_page.dart';
import 'features/inventory/inventory_page.dart';
import 'features/inventory/movement_form_page.dart';
import 'features/invoices/invoice_list_page.dart';
import 'features/invoices/invoice_form_page.dart';
import 'features/invoices/invoice_detail_page.dart';
import 'features/quotations/quotation_list_page.dart';
import 'features/purchases/purchase_list_page.dart';
import 'features/fiscal/fiscal_page.dart';
import 'features/payments/payments_page.dart';
import 'features/delivery/delivery_routes_page.dart';
import 'features/settings/settings_page.dart';
import 'features/sellers/seller_list_page.dart';
import 'features/sellers/seller_form_page.dart';
import 'features/sellers/commissions_page.dart';
import 'features/payment_methods/payment_methods_page.dart';
import 'features/fiscal_devices/fiscal_devices_page.dart';
import 'features/credit_notes/credit_note_list_page.dart';
import 'features/credit_notes/credit_note_form_page.dart';
import 'features/transfers/transfer_list_page.dart';
import 'features/transfers/transfer_form_page.dart';
import 'features/manufacturing/manufacturing_list_page.dart';
import 'features/manufacturing/manufacturing_form_page.dart';
import 'features/picking/picking_list_page.dart';
import 'features/picking/picking_detail_page.dart';
import 'features/accounts/accounts_page.dart';
import 'features/reports/reports_page.dart';
import 'features/api_tokens/api_tokens_page.dart';
import 'features/scanner/scanner_page.dart';

final router = GoRouter(
  initialLocation: '/login',
  redirect: (context, state) {
    final isLoggedIn = authApi.accessToken != null;
    final isAuthRoute = state.matchedLocation == '/login' || state.matchedLocation == '/register';
    if (!isLoggedIn && !isAuthRoute) return '/login';
    if (isLoggedIn && isAuthRoute) return '/';
    return null;
  },
  routes: [
    GoRoute(path: '/login', builder: (context, state) => const LoginPage()),
    GoRoute(path: '/register', builder: (context, state) => const RegisterPage()),
    GoRoute(path: '/', builder: (context, state) => const DashboardPage()),
    GoRoute(path: '/customers', builder: (context, state) => const CustomerListPage()),
    GoRoute(path: '/customers/new', builder: (context, state) => const CustomerFormPage()),
    GoRoute(path: '/products', builder: (context, state) => const ProductListPage()),
    GoRoute(path: '/products/new', builder: (context, state) => const ProductFormPage()),
    GoRoute(path: '/inventory', builder: (context, state) => const InventoryPage()),
    GoRoute(path: '/inventory/new', builder: (context, state) => const MovementFormPage()),
    GoRoute(path: '/invoices', builder: (context, state) => const InvoiceListPage()),
    GoRoute(path: '/invoices/new', builder: (context, state) => const InvoiceFormPage()),
    GoRoute(path: '/invoices/:id', builder: (context, state) => InvoiceDetailPage(id: state.pathParameters['id']!)),
    GoRoute(path: '/quotations', builder: (context, state) => const QuotationListPage()),
    GoRoute(path: '/purchases', builder: (context, state) => const PurchaseListPage()),
    GoRoute(path: '/fiscal', builder: (context, state) => const FiscalPage()),
    GoRoute(path: '/payments', builder: (context, state) => const PaymentsPage()),
    GoRoute(path: '/routes', builder: (context, state) => const DeliveryRoutesPage()),
    GoRoute(path: '/settings', builder: (context, state) => const SettingsPage()),
    GoRoute(path: '/sellers', builder: (context, state) => const SellerListPage()),
    GoRoute(path: '/sellers/new', builder: (context, state) => SellerFormPage(id: state.uri.queryParameters['id'])),
    GoRoute(path: '/sellers/commissions', builder: (context, state) => const CommissionsPage()),
    GoRoute(path: '/payment-methods', builder: (context, state) => const PaymentMethodsPage()),
    GoRoute(path: '/fiscal-devices', builder: (context, state) => const FiscalDevicesPage()),
    GoRoute(path: '/credit-notes', builder: (context, state) => const CreditNoteListPage()),
    GoRoute(path: '/credit-notes/new', builder: (context, state) => const CreditNoteFormPage()),
    GoRoute(path: '/transfers', builder: (context, state) => const TransferListPage()),
    GoRoute(path: '/transfers/new', builder: (context, state) => const TransferFormPage()),
    GoRoute(path: '/manufacturing', builder: (context, state) => const ManufacturingListPage()),
    GoRoute(path: '/manufacturing/new', builder: (context, state) => const ManufacturingFormPage()),
    GoRoute(path: '/picking', builder: (context, state) => const PickingListPage()),
    GoRoute(path: '/picking/:id', builder: (context, state) => PickingDetailPage(id: state.pathParameters['id']!)),
    GoRoute(path: '/accounts', builder: (context, state) => const AccountsPage()),
    GoRoute(path: '/reports', builder: (context, state) => const ReportsPage()),
    GoRoute(path: '/api-tokens', builder: (context, state) => const ApiTokensPage()),
    GoRoute(path: '/scanner', builder: (context, state) => const ScannerPage()),
  ],
);
