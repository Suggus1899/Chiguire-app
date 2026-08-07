import 'package:go_router/go_router.dart';
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

final router = GoRouter(
  initialLocation: '/login',
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
  ],
);
