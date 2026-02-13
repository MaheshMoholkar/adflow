import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:permission_handler/permission_handler.dart';
import '../../core/native/native_bridge.dart';
import '../../core/network/api_client.dart';
import '../../core/network/auth_interceptor.dart';
import '../../core/database/app_database.dart';

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen> {
  @override
  void initState() {
    super.initState();
    _initialize();
  }

  Future<void> _initialize() async {
    await Future.delayed(const Duration(milliseconds: 500));

    // Check version
    await _checkVersion();

    // Check auth — need both a token and a local user
    final hasToken = await AuthInterceptor.hasTokens();
    final db = ref.read(databaseProvider);
    final user = await db.getUser();
    if (!hasToken || user == null) {
      if (mounted) context.go('/auth/phone');
      return;
    }

    // Check required permissions
    final allGranted = await _checkAllPermissions();
    if (!allGranted) {
      if (mounted) context.go('/auth/permissions');
      return;
    }

    if (mounted) context.go('/dashboard');
  }

  Future<bool> _checkAllPermissions() async {
    try {
      final permissions = [
        Permission.phone,
        Permission.contacts,
        Permission.sms,
        Permission.notification,
      ];
      for (final p in permissions) {
        if (!await p.isGranted) return false;
      }
      final bridge = ref.read(nativeBridgeProvider);
      if (!await bridge.isBatteryOptimizationDisabled()) return false;
      return true;
    } catch (_) {
      return false;
    }
  }

  Future<void> _checkVersion() async {
    try {
      final api = ref.read(apiClientProvider);
      final response = await api.get('/app/version');
      final data = response.data['data'] as Map<String, dynamic>?;

      if (data != null) {
        final forceUpdate = data['force_update'] as bool? ?? false;
        final serverVersionCode = int.tryParse(
                data['version_code']?.toString() ?? '0') ??
            0;
        const currentVersionCode = 1; // From pubspec

        if (forceUpdate && serverVersionCode > currentVersionCode && mounted) {
          await _showUpdateDialog(
            force: true,
            downloadUrl: data['download_url'] as String? ?? '',
            releaseNotes: data['release_notes'] as String? ?? '',
          );
        }
      }
    } catch (_) {
      // Offline — continue
    }
  }

  Future<void> _showUpdateDialog({
    required bool force,
    required String downloadUrl,
    required String releaseNotes,
  }) async {
    await showDialog(
      context: context,
      barrierDismissible: !force,
      builder: (context) => AlertDialog(
        title: const Text('Update Available'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('A new version of CallFlow is available.'),
            if (releaseNotes.isNotEmpty) ...[
              const SizedBox(height: 12),
              Text(releaseNotes, style: Theme.of(context).textTheme.bodySmall),
            ],
          ],
        ),
        actions: [
          if (!force)
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Later'),
            ),
          FilledButton(
            onPressed: () {
              // Open download URL
              Navigator.pop(context);
            },
            child: const Text('Update'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.phone_in_talk,
              size: 80,
              color: Theme.of(context).colorScheme.primary,
            ),
            const SizedBox(height: 24),
            Text(
              'CallFlow',
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: Theme.of(context).colorScheme.primary,
                  ),
            ),
            const SizedBox(height: 8),
            Text(
              'Automated Call Follow-up',
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                  ),
            ),
            const SizedBox(height: 48),
            const CircularProgressIndicator(),
          ],
        ),
      ),
    );
  }
}
