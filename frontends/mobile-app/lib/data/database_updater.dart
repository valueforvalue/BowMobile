import 'dart:async';
import 'dart:io';

import 'package:http/http.dart' as http;
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';

import 'database_helper.dart';

Future<void> updateDatabaseFromUrl(String url) async {
  final uri = Uri.tryParse(url.trim());
  if (uri == null || !uri.hasScheme || !uri.hasAuthority) {
    throw const FormatException('Please enter a valid download URL.');
  }

  late final http.Response response;
  try {
    response = await http
        .get(uri, headers: {'Accept': '*/*'})
        .timeout(const Duration(seconds: 60));
  } on TimeoutException {
    throw TimeoutException(
      'Download timed out. Please check your connection and try again.',
    );
  } on SocketException catch (e) {
    throw SocketException(e.message, address: e.address, port: e.port);
  }

  if (response.statusCode != 200) {
    throw Exception('Download failed: HTTP ${response.statusCode}.');
  }

  final bytes = response.bodyBytes;
  _validateSqliteHeader(bytes);

  final docsDir = await getApplicationDocumentsDirectory();
  final dbPath = p.join(docsDir.path, 'parts.db');
  final tempPath = p.join(docsDir.path, 'parts.db.download');

  final tempFile = File(tempPath);
  await tempFile.writeAsBytes(bytes, flush: true);

  final probeDb = await openDatabase(tempPath, readOnly: true);
  await probeDb.close();

  await DatabaseHelper.instance.closeDatabase();
  await _deleteIfExists('$dbPath-wal');
  await _deleteIfExists('$dbPath-shm');
  await _deleteIfExists('$dbPath-journal');
  await _deleteIfExists(dbPath);

  await tempFile.rename(dbPath);
}

void _validateSqliteHeader(List<int> bytes) {
  const magic = [
    83,
    81,
    76,
    105,
    116,
    101,
    32,
    102,
    111,
    114,
    109,
    97,
    116,
    32,
    51,
    0,
  ];

  if (bytes.length < magic.length) {
    throw const FormatException(
      'Downloaded file is not a valid SQLite database.',
    );
  }

  for (var i = 0; i < magic.length; i++) {
    if (bytes[i] != magic[i]) {
      throw const FormatException(
        'Downloaded file is not a valid SQLite database.',
      );
    }
  }
}

Future<void> _deleteIfExists(String path) async {
  final file = File(path);
  if (await file.exists()) {
    await file.delete();
  }
}
