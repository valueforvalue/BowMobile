import 'package:flutter/material.dart';
import '../data/models.dart';
import 'result_card.dart';

/// Called when the set of selected [PartOccurrence]s changes.
typedef SelectionChanged = void Function(
  Map<PartOccurrence, String> selected, // occurrence → basePart
);

class ResultsList extends StatefulWidget {
  final List<GroupedResult> results;
  final String query;
  final SelectionChanged onSelectionChanged;

  const ResultsList({
    super.key,
    required this.results,
    required this.query,
    required this.onSelectionChanged,
  });

  @override
  State<ResultsList> createState() => _ResultsListState();
}

class _ResultsListState extends State<ResultsList> {
  // basePart+index → selected occurrence
  final Map<PartOccurrence, String> _selected = {};

  void _onRowToggled(PartOccurrence occ, String basePart, bool checked) {
    setState(() {
      if (checked) {
        _selected[occ] = basePart;
      } else {
        _selected.remove(occ);
      }
    });
    widget.onSelectionChanged(Map.unmodifiable(_selected));
  }

  @override
  Widget build(BuildContext context) {
    if (widget.results.isEmpty && widget.query.isNotEmpty) {
      return Container(
        margin: const EdgeInsets.only(top: 16),
        padding: const EdgeInsets.symmetric(vertical: 60),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: Colors.grey.shade300, style: BorderStyle.solid),
        ),
        child: Center(
          child: Text(
            'No parts found matching "${widget.query}"',
            style: TextStyle(color: Colors.grey.shade500, fontSize: 16),
          ),
        ),
      );
    }

    return ListView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      itemCount: widget.results.length,
      itemBuilder: (context, i) {
        final group = widget.results[i];
        return ResultCard(
          group: group,
          selectedOccurrences: _selected.keys.toSet(),
          onRowToggled: _onRowToggled,
        );
      },
    );
  }
}
