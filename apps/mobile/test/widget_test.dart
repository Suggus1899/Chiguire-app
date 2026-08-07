import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:chiguire_mobile/main.dart';

void main() {
  testWidgets('App renders login page by default', (WidgetTester tester) async {
    await tester.pumpWidget(
      const ProviderScope(child: ChiguireApp()),
    );
    await tester.pumpAndSettle();
    expect(find.text('Iniciar sesión'), findsOneWidget);
  });
}
