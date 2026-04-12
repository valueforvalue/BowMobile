import 'dart:async';
import 'dart:io';

import 'package:flutter/material.dart';

import '../data/database_updater.dart';

Future<void> showUpdateDbDialog(BuildContext context) async {
  final controller = TextEditingController();
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
                  'Paste a direct download URL for a parts.db file from a GitHub release or other trusted source.',
                  style: TextStyle(fontSize: 13),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: controller,
                  enabled: !loading,
                  keyboardType: TextInputType.url,
                  autocorrect: false,
                  autofocus: true,
                  decoration: InputDecoration(
                    hintText:
                        'https://github.com/.../releases/download/.../parts.db',
                    hintStyle: const TextStyle(fontSize: 12),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 10,
                    ),
                  ),
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
                        final url = controller.text.trim();
                        if (url.isEmpty) {
                          setDialogState(() {
                            errorMessage = 'Please enter a URL.';
                          });
                          return;
                        }

                        setDialogState(() {
                          loading = true;
                          errorMessage = null;
                        });

                        try {
                          await updateDatabaseFromUrl(url);
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
