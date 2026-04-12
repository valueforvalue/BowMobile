import 'dart:io';

import 'package:flutter/services.dart';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';

class DatabaseHelper {
  DatabaseHelper._();
  static final DatabaseHelper instance = DatabaseHelper._();

  Database? _db;

  Future<Database> get database async {
    _db ??= await _initDb();
    return _db!;
  }

  Future<Database> _initDb() async {
    final docsDir = await getApplicationDocumentsDirectory();
    final dbPath = p.join(docsDir.path, 'parts.db');

    // Copy asset DB to documents directory if not already present
    if (!File(dbPath).existsSync()) {
      final data = await rootBundle.load('assets/parts.db');
      final bytes = data.buffer.asUint8List();
      await File(dbPath).writeAsBytes(bytes, flush: true);
    }

    return openDatabase(dbPath, readOnly: true);
  }

  /// Force-replace the on-device DB with the bundled asset (e.g. after an app update).
  Future<void> resetToAsset() async {
    _db?.close();
    _db = null;
    final docsDir = await getApplicationDocumentsDirectory();
    final dbPath = p.join(docsDir.path, 'parts.db');
    final data = await rootBundle.load('assets/parts.db');
    final bytes = data.buffer.asUint8List();
    await File(dbPath).writeAsBytes(bytes, flush: true);
  }

  /// Close the current database connection so the file can be replaced.
  Future<void> closeDatabase() async {
    await _db?.close();
    _db = null;
  }
}
