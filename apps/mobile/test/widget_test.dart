import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:chiguire_mobile/main.dart';

void main() {
  testWidgets('App renders Chiguire title', (WidgetTester tester) async {
    await tester.pumpWidget(
      const ProviderScope(child: ChiguireApp()),
    );
    expect(find.text('Chiguire'), findsOneWidget);
  });
}
