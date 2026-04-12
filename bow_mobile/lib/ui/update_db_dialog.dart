import 'dart:async';
import 'dart:io';

import 'package:flutter/material.dart';

import '../data/database_updater.dart';

/// Shows a dialog that lets the user enter a GitHub raw URL for a parts.db
/// file and downloads it to replace the current on-device database.
Future<void> showUpdateDbDialog(BuildContext context) async {
  final controller = TextEditingController();
  bool loading = false;
  String? errorMessage;

  await showDialog<void>(
    context: context,
    barrierDismissible: !loading,
    builder: (ctx) {
      return StatefulBuilder(
        builder: (ctx, setDialogState) {
          return AlertDialog(
            title: const Text('Update Database from GitHub'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Paste a direct download URL for a parts.db file '
                  '(e.g. a GitHub release asset raw URL).',
                  style: TextStyle(fontSize: 13),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: controller,
                  enabled: !loading,
                  decoration: InputDecoration(
                    hintText: 'https://github.com/.../releases/download/.../parts.db',
                    hintStyle: const TextStyle(fontSize: 12),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 10,
                    ),
                  ),
                  keyboardType: TextInputType.url,
                  autocorrect: false,
                  autofocus: true,
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
                      Text('Downloading…', style: TextStyle(fontSize: 13)),
                    ],
                  ),
                ],
              ],
            ),
            actions: [
              TextButton(
                onPressed: loading ? null : () => Navigator.of(ctx).pop(),
                child: const Text('Cancel'),
              ),
              FilledButton(
                onPressed: loading
                    ? null
                    : () async {
                        final url = controller.text.trim();
                        if (url.isEmpty) {
                          setDialogState(() => errorMessage = 'Please enter a URL.');
                          return;
                        }
                        setDialogState(() {
                          loading = true;
                          errorMessage = null;
                        });
                        try {
                          await updateDatabaseFromUrl(url);
                          if (ctx.mounted) {
                            Navigator.of(ctx).pop();
                            ScaffoldMessenger.of(ctx).showSnackBar(
                              const SnackBar(
                                content: Text('Database updated successfully!'),
                                backgroundColor: Color(0xFF059669),
                              ),
                            );
                          }
                        } catch (e) {
                          setDialogState(() {
                            loading = false;
                            errorMessage = _friendlyError(e);
                          });
                        }
                      },
                style: FilledButton.styleFrom(
                  backgroundColor: const Color(0xFF2563EB),
                ),
                child: const Text('Download & Apply'),
              ),
            ],
          );
        },
      );
    },
  );
}

/// Converts a caught error to a short, user-readable message.
String _friendlyError(Object e) {
  if (e is TimeoutException) {
    return 'Download timed out. Please check your connection and try again.';
  }
  if (e is SocketException) {
    return 'Network error. Please check your connection and try again.';
  }
  if (e is FormatException) {
    return e.message;
  }
  // Exception messages from updateDatabaseFromUrl are already user-friendly.
  final msg = e.toString();
  final prefix = 'Exception: ';
  return msg.startsWith(prefix) ? msg.substring(prefix.length) : msg;
}
