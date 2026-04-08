import 'package:flutter/material.dart';
import '../data/models.dart';

typedef RowToggled = void Function(
  PartOccurrence occ,
  String basePart,
  bool checked,
);

class ResultCard extends StatelessWidget {
  final GroupedResult group;
  final Set<PartOccurrence> selectedOccurrences;
  final RowToggled onRowToggled;

  const ResultCard({
    super.key,
    required this.group,
    required this.selectedOccurrences,
    required this.onRowToggled,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 24),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.grey.shade200),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.04),
            blurRadius: 4,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      clipBehavior: Clip.antiAlias,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _CardHeader(group: group),
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: _OccurrencesTable(
              group: group,
              selectedOccurrences: selectedOccurrences,
              onRowToggled: onRowToggled,
            ),
          ),
        ],
      ),
    );
  }
}

class _CardHeader extends StatelessWidget {
  final GroupedResult group;
  const _CardHeader({required this.group});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
      decoration: BoxDecoration(
        color: Colors.grey.shade100,
        border: Border(bottom: BorderSide(color: Colors.grey.shade200)),
      ),
      child: Row(
        children: [
          Text(
            group.basePart,
            style: const TextStyle(
              fontSize: 17,
              fontWeight: FontWeight.bold,
              color: Color(0xFF374151),
            ),
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Text(
              group.description,
              style: TextStyle(fontSize: 15, color: Colors.grey.shade500),
              overflow: TextOverflow.ellipsis,
            ),
          ),
        ],
      ),
    );
  }
}

class _OccurrencesTable extends StatelessWidget {
  final GroupedResult group;
  final Set<PartOccurrence> selectedOccurrences;
  final RowToggled onRowToggled;

  const _OccurrencesTable({
    required this.group,
    required this.selectedOccurrences,
    required this.onRowToggled,
  });

  @override
  Widget build(BuildContext context) {
    const headerStyle = TextStyle(
      fontSize: 10,
      fontWeight: FontWeight.bold,
      color: Color(0xFF6B7280),
      letterSpacing: 0.8,
    );
    const monoStyle = TextStyle(
      fontFamily: 'monospace',
      fontSize: 13,
      color: Color(0xFF4B5563),
    );
    const cellPadding = EdgeInsets.symmetric(horizontal: 16, vertical: 12);

    return Table(
      defaultColumnWidth: const IntrinsicColumnWidth(),
      border: TableBorder(
        horizontalInside: BorderSide(color: Colors.grey.shade100),
      ),
      children: [
        // Header row
        TableRow(
          decoration: BoxDecoration(color: Colors.grey.shade50),
          children: [
            _headerCell('', headerStyle, cellPadding),
            _headerCell('Model / Series', headerStyle, cellPadding),
            _headerCell('Catalog Rev', headerStyle, cellPadding),
            _headerCell('Figure', headerStyle, cellPadding),
            _headerCell('Key No', headerStyle, cellPadding),
            _headerCell('Full Part Number', headerStyle, cellPadding),
            _headerCell('Remarks', headerStyle, cellPadding),
          ],
        ),
        // Data rows
        for (final occ in group.occurrences)
          TableRow(
            decoration: BoxDecoration(
              color: selectedOccurrences.contains(occ)
                  ? const Color(0xFFFEF9C3)
                  : Colors.white,
            ),
            children: [
              // Checkbox
              Padding(
                padding: cellPadding,
                child: Checkbox(
                  value: selectedOccurrences.contains(occ),
                  onChanged: (v) =>
                      onRowToggled(occ, group.basePart, v ?? false),
                  activeColor: const Color(0xFF2563EB),
                  materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                ),
              ),
              // Model / Series
              Padding(
                padding: cellPadding,
                child: Text(
                  occ.modelSeries,
                  style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13),
                ),
              ),
              // Catalog Rev badge
              Padding(
                padding: cellPadding,
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 7, vertical: 2),
                  decoration: BoxDecoration(
                    color: Colors.grey.shade600,
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    'Rev ${occ.manualRevision}',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ),
              Padding(padding: cellPadding, child: Text(occ.figureId, style: monoStyle)),
              Padding(padding: cellPadding, child: Text(occ.keyNo, style: monoStyle)),
              Padding(
                padding: cellPadding,
                child: Text(
                  occ.fullPartNumber,
                  style: monoStyle.copyWith(fontWeight: FontWeight.w500),
                ),
              ),
              Padding(
                padding: cellPadding,
                child: Text(
                  occ.remarks,
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey.shade500,
                    fontStyle: FontStyle.italic,
                  ),
                ),
              ),
            ],
          ),
      ],
    );
  }

  Widget _headerCell(String text, TextStyle style, EdgeInsets padding) {
    return Padding(
      padding: padding,
      child: Text(text.toUpperCase(), style: style),
    );
  }
}
