import 'package:flutter/material.dart';
import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart' as pw;
import 'package:share_plus/share_plus.dart';
import 'package:path_provider/path_provider.dart';
import 'package:path/path.dart' as p;
import 'dart:io';

import '../data/models.dart';

/// Generates a PDF from selected parts and triggers the system share sheet.
Future<void> shareSelectedParts(
  BuildContext context,
  Map<PartOccurrence, String> selected,
) async {
  if (selected.isEmpty) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Select at least one part to share.')),
    );
    return;
  }

  final pdf = pw.Document();

  pdf.addPage(
    pw.MultiPage(
      pageFormat: PdfPageFormat.letter,
      margin: const pw.EdgeInsets.all(40),
      build: (ctx) => [
        pw.Header(
          level: 0,
          child: pw.Text(
            'Selected Parts List',
            style: pw.TextStyle(fontSize: 20, fontWeight: pw.FontWeight.bold),
          ),
        ),
        pw.SizedBox(height: 16),
        pw.Table.fromTextArray(
          headers: [
            'Base Part',
            'Description',
            'Model / Series',
            'Figure',
            'Key No',
            'Full Part Number',
            'Remarks',
          ],
          headerStyle: pw.TextStyle(
            fontWeight: pw.FontWeight.bold,
            fontSize: 9,
          ),
          cellStyle: const pw.TextStyle(fontSize: 8),
          headerDecoration: const pw.BoxDecoration(color: PdfColors.grey300),
          cellAlignments: {
            0: pw.Alignment.centerLeft,
            1: pw.Alignment.centerLeft,
            2: pw.Alignment.centerLeft,
            3: pw.Alignment.centerLeft,
            4: pw.Alignment.centerLeft,
            5: pw.Alignment.centerLeft,
            6: pw.Alignment.centerLeft,
          },
          data: selected.entries
              .map(
                (e) => [
                  e.value,
                  e.key.description,
                  e.key.modelSeries,
                  e.key.figureId,
                  e.key.keyNo,
                  e.key.fullPartNumber,
                  e.key.remarks,
                ],
              )
              .toList(),
        ),
      ],
    ),
  );

  final bytes = await pdf.save();
  final tmpDir = await getTemporaryDirectory();
  final file = File(p.join(tmpDir.path, 'bow_parts_selection.pdf'));
  await file.writeAsBytes(bytes);

  await Share.shareXFiles(
    [XFile(file.path)],
    subject: 'Bow Parts Selection',
  );
}
