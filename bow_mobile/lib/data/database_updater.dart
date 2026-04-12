import 'dart:io';

import 'package:http/http.dart' as http;
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';

import 'database_helper.dart';

/// Downloads a new parts.db from [url] (raw GitHub URL or any direct download
/// link) and replaces the on-device database.
///
/// Throws on HTTP error, IO error, or if the downloaded file does not look
/// like a valid SQLite database.
Future<void> updateDatabaseFromUrl(String url) async {
  final uri = Uri.parse(url.trim());
  final response = await http
      .get(uri, headers: {'Accept': '*/*'})
      .timeout(const Duration(seconds: 60));

  if (response.statusCode != 200) {
    throw Exception(
      'Download failed: HTTP ${response.statusCode}',
    );
  }

  final bytes = response.bodyBytes;
  _validateSqliteHeader(bytes);

  final docsDir = await getApplicationDocumentsDirectory();
  final dbPath = p.join(docsDir.path, 'parts.db');

  // Close existing connection before replacing the file
  await DatabaseHelper.instance.closeDatabase();

  await File(dbPath).writeAsBytes(bytes, flush: true);

  // Re-open to verify the new DB is readable
  final db = await openDatabase(dbPath, readOnly: true);
  await db.close();
}

/// Checks that [bytes] starts with the SQLite magic header.
void _validateSqliteHeader(List<int> bytes) {
  const magic = [83, 81, 76, 105, 116, 101, 32, 102, 111, 114, 109, 97, 116, 32, 51, 0];
  if (bytes.length < 16) {
    throw const FormatException('Downloaded file is not a valid SQLite database.');
  }
  for (int i = 0; i < 16; i++) {
    if (bytes[i] != magic[i]) {
      throw const FormatException('Downloaded file is not a valid SQLite database.');
    }
  }
}
