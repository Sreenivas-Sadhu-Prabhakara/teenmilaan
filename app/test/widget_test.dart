import 'package:flutter_test/flutter_test.dart';

import 'package:teenmilaan_app/main.dart';

void main() {
  test('mismatch flags qty+rate and totals value leak', () {
    final m = Match('bolt', 100, 95, 12, 11);
    expect(m.mismatch, true);
    expect(m.valueLeak, closeTo(150, 1e-9));
  });

  test('clean match is not flagged', () {
    expect(Match('x', 10, 10, 5, 5).mismatch, false);
  });

  testWidgets('renders match title', (tester) async {
    await tester.pumpWidget(const TeenmilaanApp());
    expect(find.text('Item'), findsOneWidget);
  });
}
