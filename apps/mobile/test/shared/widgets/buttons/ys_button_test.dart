import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:yousoon/shared/widgets/buttons/ys_button.dart';

void main() {
  group('YsButton', () {
    testWidgets('should render primary button with label', (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Se connecter',
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.text('Se connecter'), findsOneWidget);
      expect(find.byType(ElevatedButton), findsOneWidget);
    });

    testWidgets('should call onPressed when tapped', (tester) async {
      var pressed = false;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Tap me',
              onPressed: () => pressed = true,
            ),
          ),
        ),
      );

      await tester.tap(find.text('Tap me'));
      await tester.pump();

      expect(pressed, true);
    });

    testWidgets('should be disabled when onPressed is null', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Disabled',
              onPressed: null,
            ),
          ),
        ),
      );

      final button = tester.widget<ElevatedButton>(find.byType(ElevatedButton));
      expect(button.onPressed, isNull);
    });

    testWidgets('should show loading indicator when isLoading is true',
        (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Loading',
              onPressed: () {},
              isLoading: true,
            ),
          ),
        ),
      );

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      expect(find.text('Loading'), findsNothing);
    });

    testWidgets('should not call onPressed when loading', (tester) async {
      var pressed = false;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Loading',
              onPressed: () => pressed = true,
              isLoading: true,
            ),
          ),
        ),
      );

      await tester.tap(find.byType(ElevatedButton));
      await tester.pump();

      expect(pressed, false);
    });

    testWidgets('should render secondary button', (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Secondary',
              onPressed: () {},
              variant: YsButtonVariant.secondary,
            ),
          ),
        ),
      );

      expect(find.byType(OutlinedButton), findsOneWidget);
      expect(find.text('Secondary'), findsOneWidget);
    });

    testWidgets('should render tertiary button', (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Tertiary',
              onPressed: () {},
              variant: YsButtonVariant.tertiary,
            ),
          ),
        ),
      );

      expect(find.byType(TextButton), findsOneWidget);
      expect(find.text('Tertiary'), findsOneWidget);
    });

    testWidgets('should render icon when provided', (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'With Icon',
              onPressed: () {},
              icon: Icons.login,
            ),
          ),
        ),
      );

      expect(find.byIcon(Icons.login), findsOneWidget);
      expect(find.text('With Icon'), findsOneWidget);
    });

    testWidgets('should respect small size', (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Small',
              onPressed: () {},
              size: YsButtonSize.small,
            ),
          ),
        ),
      );

      final sizedBox = tester.widget<SizedBox>(find.byType(SizedBox).first);
      expect(sizedBox.height, 40);
    });

    testWidgets('should respect large size', (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Large',
              onPressed: () {},
              size: YsButtonSize.large,
            ),
          ),
        ),
      );

      final sizedBox = tester.widget<SizedBox>(find.byType(SizedBox).first);
      expect(sizedBox.height, 50);
    });

    testWidgets('should be full width by default', (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: YsButton(
              label: 'Full Width',
              onPressed: () {},
            ),
          ),
        ),
      );

      final sizedBox = tester.widget<SizedBox>(find.byType(SizedBox).first);
      expect(sizedBox.width, double.infinity);
    });

    testWidgets('should respect isFullWidth = false', (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Center(
              child: YsButton(
                label: 'Compact',
                onPressed: () {},
                isFullWidth: false,
              ),
            ),
          ),
        ),
      );

      final sizedBox = tester.widget<SizedBox>(find.byType(SizedBox).first);
      expect(sizedBox.width, isNull);
    });
  });

  group('YsButtonSize', () {
    test('small should have correct values', () {
      expect(YsButtonSize.small.height, 40);
      expect(YsButtonSize.small.fontSize, 14);
    });

    test('medium should have correct values', () {
      expect(YsButtonSize.medium.height, 44);
      expect(YsButtonSize.medium.fontSize, 14);
    });

    test('large should have correct values', () {
      expect(YsButtonSize.large.height, 50);
      expect(YsButtonSize.large.fontSize, 16);
    });
  });
}
