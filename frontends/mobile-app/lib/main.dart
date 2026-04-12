import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'ui/search_screen.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.landscapeLeft,
    DeviceOrientation.landscapeRight,
  ]);
  runApp(const BowApp());
}

class BowApp extends StatelessWidget {
  const BowApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Bow Parts Cross-Reference',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color(0xFF2563EB),
          brightness: Brightness.light,
        ),
        useMaterial3: true,
        fontFamily: 'Roboto',
      ),
      home: const SearchScreen(),
    );
  }
}
