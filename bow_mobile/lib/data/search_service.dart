import 'dart:io';

import 'package:path_provider/path_provider.dart';
import 'package:path/path.dart' as p;

import 'database_helper.dart';
import 'models.dart';

class SearchService {
  SearchService._();
  static final SearchService instance = SearchService._();

  /// Smart search matching the Go performSmartSearch logic.
  /// Searches part_number, base_part (with and without hyphens), description, remarks.
  Future<List<GroupedResult>> search(String query) async {
    if (query.trim().isEmpty) return [];

    final db = await DatabaseHelper.instance.database;
    final q = query.trim().toUpperCase();
    final likeQ = '%$q%';
    final normalizedQ = q.replaceAll('-', '');
    final likeNormalizedQ = '%$normalizedQ%';

    const sql = '''
      SELECT p.base_part, p.description, m.model_series, m.revision,
             p.figure_id, p.key_no, p.part_number, p.revision AS part_rev, p.remarks
      FROM parts p
      JOIN manuals m ON p.manual_id = m.id
      WHERE p.part_number LIKE ?
         OR REPLACE(p.part_number, '-', '') LIKE ?
         OR p.base_part LIKE ?
         OR REPLACE(p.base_part, '-', '') LIKE ?
         OR p.description LIKE ?
         OR p.remarks LIKE ?
      ORDER BY p.base_part, m.model_series
    ''';

    final rows = await db.rawQuery(sql, [
      likeQ,
      likeNormalizedQ,
      likeQ,
      likeNormalizedQ,
      likeQ,
      likeQ,
    ]);

    // Group by base_part preserving insertion order
    final order = <String>[];
    final groups = <String, GroupedResult>{};

    for (final row in rows) {
      final base = row['base_part'] as String? ?? '';
      final desc = row['description'] as String? ?? '';

      if (!groups.containsKey(base)) {
        order.add(base);
        groups[base] = GroupedResult(
          basePart: base,
          description: desc,
          occurrences: [],
        );
      }

      final existing = groups[base]!;
      groups[base] = GroupedResult(
        basePart: existing.basePart,
        description: existing.description,
        occurrences: [
          ...existing.occurrences,
          PartOccurrence(
            modelSeries: row['model_series'] as String? ?? '',
            manualRevision: row['revision'] as String? ?? '',
            figureId: row['figure_id'] as String? ?? '',
            keyNo: row['key_no'] as String? ?? '',
            fullPartNumber: row['part_number'] as String? ?? '',
            revision: row['part_rev'] as String? ?? '',
            description: desc,
            remarks: row['remarks'] as String? ?? '',
          ),
        ],
      );
    }

    return order.map((b) => groups[b]!).toList();
  }

  Future<List<ManualInfo>> getModels() async {
    final db = await DatabaseHelper.instance.database;
    final rows = await db.rawQuery(
      'SELECT model_series, revision, filename FROM manuals ORDER BY model_series',
    );
    return rows
        .map(
          (r) => ManualInfo(
            modelSeries: r['model_series'] as String? ?? '',
            revision: r['revision'] as String? ?? '',
            filename: r['filename'] as String? ?? '',
          ),
        )
        .toList();
  }

  Future<DbInfo> getDbInfo() async {
    final db = await DatabaseHelper.instance.database;

    String lastUpdated = 'Unknown';
    final metaRows = await db.rawQuery(
      'SELECT last_updated FROM metadata WHERE id = 1',
    );
    if (metaRows.isNotEmpty) {
      lastUpdated = metaRows.first['last_updated'] as String? ?? 'Unknown';
    }

    int count = 0;
    final countRows = await db.rawQuery('SELECT COUNT(*) as c FROM parts');
    if (countRows.isNotEmpty) {
      count = (countRows.first['c'] as int?) ?? 0;
    }

    int sizeKb = 0;
    try {
      final docsDir = await getApplicationDocumentsDirectory();
      final file = File(p.join(docsDir.path, 'parts.db'));
      if (file.existsSync()) {
        sizeKb = file.lengthSync() ~/ 1024;
      }
    } catch (_) {}

    return DbInfo(lastUpdated: lastUpdated, partCount: count, sizeKb: sizeKb);
  }
}
