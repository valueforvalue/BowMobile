class PartOccurrence {
  final String modelSeries;
  final String manualRevision;
  final String figureId;
  final String keyNo;
  final String fullPartNumber;
  final String revision;
  final String description;
  final String remarks;

  const PartOccurrence({
    required this.modelSeries,
    required this.manualRevision,
    required this.figureId,
    required this.keyNo,
    required this.fullPartNumber,
    required this.revision,
    required this.description,
    required this.remarks,
  });
}

class GroupedResult {
  final String basePart;
  final String description;
  final List<PartOccurrence> occurrences;

  const GroupedResult({
    required this.basePart,
    required this.description,
    required this.occurrences,
  });
}

class ManualInfo {
  final String modelSeries;
  final String revision;
  final String filename;

  const ManualInfo({
    required this.modelSeries,
    required this.revision,
    required this.filename,
  });
}

class DbInfo {
  final String lastUpdated;
  final int partCount;
  final int sizeKb;

  const DbInfo({
    required this.lastUpdated,
    required this.partCount,
    required this.sizeKb,
  });

  @override
  String toString() => '$lastUpdated | Parts: $partCount | Size: $sizeKb KB';
}
