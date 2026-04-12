import 'package:flutter/material.dart';
import '../data/models.dart';
import '../data/search_service.dart';

class ModelsScreen extends StatefulWidget {
  const ModelsScreen({super.key});

  @override
  State<ModelsScreen> createState() => _ModelsScreenState();
}

class _ModelsScreenState extends State<ModelsScreen> {
  late final Future<List<ManualInfo>> _future;

  @override
  void initState() {
    super.initState();
    _future = SearchService.instance.getModels();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8F9FA),
      appBar: AppBar(
        title: const Text('Manuals in Database'),
        backgroundColor: Colors.white,
        foregroundColor: const Color(0xFF374151),
        elevation: 1,
      ),
      body: FutureBuilder<List<ManualInfo>>(
        future: _future,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          }
          final manuals = snapshot.data ?? [];
          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Container(
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
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
                    decoration: BoxDecoration(
                      color: Colors.grey.shade100,
                      border: Border(bottom: BorderSide(color: Colors.grey.shade200)),
                    ),
                    child: const Text(
                      'Manuals in Database',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                        color: Color(0xFF374151),
                      ),
                    ),
                  ),
                  SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    child: _ManualsTable(manuals: manuals),
                  ),
                  Padding(
                    padding: const EdgeInsets.all(16),
                    child: Text(
                      'Total manuals: ${manuals.length}',
                      style: TextStyle(
                        color: Colors.grey.shade400,
                        fontSize: 12,
                        fontStyle: FontStyle.italic,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class _ManualsTable extends StatelessWidget {
  final List<ManualInfo> manuals;
  const _ManualsTable({required this.manuals});

  @override
  Widget build(BuildContext context) {
    const headerStyle = TextStyle(
      fontSize: 10,
      fontWeight: FontWeight.bold,
      color: Color(0xFF6B7280),
      letterSpacing: 0.8,
    );
    const cellPadding = EdgeInsets.symmetric(horizontal: 16, vertical: 10);

    return Table(
      defaultColumnWidth: const IntrinsicColumnWidth(),
      border: TableBorder(
        horizontalInside: BorderSide(color: Colors.grey.shade100),
      ),
      children: [
        TableRow(
          decoration: BoxDecoration(color: Colors.grey.shade50),
          children: [
            Padding(padding: cellPadding, child: Text('MODEL SERIES', style: headerStyle)),
            Padding(padding: cellPadding, child: Text('REVISION', style: headerStyle)),
            Padding(padding: cellPadding, child: Text('SOURCE FILENAME', style: headerStyle)),
          ],
        ),
        for (final m in manuals)
          TableRow(
            decoration: const BoxDecoration(color: Colors.white),
            children: [
              Padding(
                padding: cellPadding,
                child: Text(
                  m.modelSeries,
                  style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13),
                ),
              ),
              Padding(
                padding: cellPadding,
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 7, vertical: 2),
                  decoration: BoxDecoration(
                    color: Colors.grey.shade200,
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    'Rev ${m.revision}',
                    style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w500),
                  ),
                ),
              ),
              Padding(
                padding: cellPadding,
                child: Text(
                  m.filename,
                  style: TextStyle(
                    fontFamily: 'monospace',
                    fontSize: 10,
                    color: Colors.grey.shade500,
                  ),
                ),
              ),
            ],
          ),
      ],
    );
  }
}
