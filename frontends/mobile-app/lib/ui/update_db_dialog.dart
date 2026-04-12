import 'dart:async';
import 'dart:io';

import 'package:flutter/material.dart';

import '../data/database_updater.dart';

Future<void> showUpdateDbDialog(BuildContext context) async {
  var loading = false;
  String? errorMessage;

  await showDialog<void>(
    context: context,
    builder: (dialogContext) {
      return StatefulBuilder(
        builder: (dialogContext, setDialogState) {
          return AlertDialog(
            title: const Text('Update Database from GitHub'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Download the latest parts.db from the official BowDB GitHub repository and replace the local database.',
                  style: TextStyle(fontSize: 13),
                ),
                const SizedBox(height: 12),
                SelectableText(
                  officialDatabaseUrl,
                  style: const TextStyle(fontSize: 12),
                ),
                if (errorMessage != null) ...[
                  const SizedBox(height: 8),
                  Text(
                    errorMessage!,
                    style: const TextStyle(color: Colors.red, fontSize: 12),
                  ),
                ],
                if (loading) ...[
                  const SizedBox(height: 16),
                  const Row(
                    children: [
                      SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      ),
                      SizedBox(width: 12),
                      Text('Downloading...', style: TextStyle(fontSize: 13)),
                    ],
                  ),
                ],
              ],
            ),
            actions: [
              TextButton(
                onPressed: loading
                    ? null
                    : () => Navigator.of(dialogContext).pop(),
                child: const Text('Cancel'),
              ),
              FilledButton(
                onPressed: loading
                    ? null
                    : () async {
                        setDialogState(() {
                          loading = true;
                          errorMessage = null;
                        });

                        try {
                          await updateDatabaseFromOfficialSource();
                          if (!dialogContext.mounted) {
                            return;
                          }
                          Navigator.of(dialogContext).pop();
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('Database updated successfully.'),
                              backgroundColor: Color(0xFF059669),
                            ),
                          );
                        } catch (e) {
                          setDialogState(() {
                            loading = false;
                            errorMessage = _friendlyError(e);
                          });
                        }
                      },
                child: const Text('Download & Apply'),
              ),
            ],
          );
        },
      );
    },
  );
}

String _friendlyError(Object error) {
  if (error is TimeoutException) {
    return 'Download timed out. Please check your connection and try again.';
  }
  if (error is SocketException) {
    return 'Network error. Please check your connection and try again.';
  }
  if (error is FormatException) {
    return error.message;
  }

  final message = error.toString();
  const prefix = 'Exception: ';
  return message.startsWith(prefix) ? message.substring(prefix.length) : message;
}
