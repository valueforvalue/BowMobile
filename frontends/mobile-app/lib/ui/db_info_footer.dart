import 'package:flutter/material.dart';
import '../data/models.dart';
import '../data/search_service.dart';

class DbInfoFooter extends StatefulWidget {
  const DbInfoFooter({super.key});

  @override
  State<DbInfoFooter> createState() => _DbInfoFooterState();
}

class _DbInfoFooterState extends State<DbInfoFooter> {
  late final Future<DbInfo> _future;

  @override
  void initState() {
    super.initState();
    _future = SearchService.instance.getDbInfo();
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<DbInfo>(
      future: _future,
      builder: (context, snapshot) {
        final info = snapshot.data?.toString() ?? 'Loading...';
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 32),
          child: Column(
            children: [
              const Divider(),
              const SizedBox(height: 12),
              Text(
                'Bow Mobile v1.1.3',
                style: TextStyle(color: Colors.grey.shade400, fontSize: 11),
              ),
              const SizedBox(height: 4),
              Text(
                'Database Last Updated: $info',
                style: TextStyle(color: Colors.grey.shade400, fontSize: 11),
              ),
              const SizedBox(height: 4),
              Text(
                '© 2026',
                style: TextStyle(color: Colors.grey.shade400, fontSize: 11),
              ),
            ],
          ),
        );
      },
    );
  }
}
