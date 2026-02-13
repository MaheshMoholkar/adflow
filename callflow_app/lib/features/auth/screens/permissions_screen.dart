import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:permission_handler/permission_handler.dart';
import '../../../core/native/native_bridge.dart';

class PermissionsScreen extends ConsumerStatefulWidget {
  const PermissionsScreen({super.key});

  @override
  ConsumerState<PermissionsScreen> createState() => _PermissionsScreenState();
}

class _PermissionsScreenState extends ConsumerState<PermissionsScreen>
    with WidgetsBindingObserver {
  bool _phoneGranted = false;
  bool _contactsGranted = false;
  bool _smsGranted = false;
  bool _notificationGranted = false;
  bool _batteryOptimizationDisabled = false;
  bool _loaded = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    _checkAll();
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state == AppLifecycleState.resumed) {
      _checkAll();
    }
  }

  Future<void> _checkAll() async {
    try {
      _phoneGranted = await Permission.phone.isGranted;
    } catch (_) {}
    try {
      _contactsGranted = await Permission.contacts.isGranted;
    } catch (_) {}
    try {
      _smsGranted = await Permission.sms.isGranted;
    } catch (_) {}
    try {
      _notificationGranted = await Permission.notification.isGranted;
    } catch (_) {}
    try {
      final bridge = ref.read(nativeBridgeProvider);
      _batteryOptimizationDisabled =
          await bridge.isBatteryOptimizationDisabled();
    } catch (_) {}

    if (mounted) setState(() => _loaded = true);
  }

  Future<void> _request(Permission p) async {
    try {
      final status = await p.request();
      if (status.isPermanentlyDenied) {
        openAppSettings();
      }
    } catch (_) {}
    await _checkAll();
  }

  bool get _allGranted =>
      _phoneGranted &&
      _contactsGranted &&
      _smsGranted &&
      _notificationGranted &&
      _batteryOptimizationDisabled;

  @override
  Widget build(BuildContext context) {
    final bridge = ref.watch(nativeBridgeProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Permissions')),
      body: !_loaded
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              padding: const EdgeInsets.all(16),
              children: [
                Text(
                  'Grant all permissions to continue',
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: Theme.of(context).colorScheme.onSurfaceVariant,
                      ),
                ),
                const SizedBox(height: 16),
                _tile(
                  icon: Icons.phone,
                  title: 'Phone Access',
                  subtitle: 'Detect incoming and outgoing calls',
                  granted: _phoneGranted,
                  onTap: () => _request(Permission.phone),
                ),
                _tile(
                  icon: Icons.contacts,
                  title: 'Contacts',
                  subtitle: 'Resolve caller names for personalized messages',
                  granted: _contactsGranted,
                  onTap: () => _request(Permission.contacts),
                ),
                _tile(
                  icon: Icons.sms,
                  title: 'SMS',
                  subtitle: 'Send automated SMS messages after calls',
                  granted: _smsGranted,
                  onTap: () => _request(Permission.sms),
                ),
                _tile(
                  icon: Icons.notifications,
                  title: 'Notifications',
                  subtitle: 'Show service status and delivery updates',
                  granted: _notificationGranted,
                  onTap: () => _request(Permission.notification),
                ),
                _tile(
                  icon: Icons.battery_saver,
                  title: 'Battery Optimization',
                  subtitle: 'Allow CallFlow to run in background',
                  granted: _batteryOptimizationDisabled,
                  onTap: () => bridge.requestBatteryOptimization(),
                ),
                const SizedBox(height: 24),
                FilledButton(
                  onPressed:
                      _allGranted ? () => context.go('/dashboard') : null,
                  child: const Text('Continue'),
                ),
              ],
            ),
    );
  }

  Widget _tile({
    required IconData icon,
    required String title,
    required String subtitle,
    required bool granted,
    required VoidCallback onTap,
  }) {
    return ListTile(
      leading: Icon(
        granted ? Icons.check_circle : icon,
        color: granted
            ? Theme.of(context).colorScheme.primary
            : Theme.of(context).colorScheme.onSurfaceVariant,
      ),
      title: Text(title),
      subtitle: Text(subtitle),
      trailing: granted
          ? null
          : const Icon(Icons.chevron_right),
      onTap: granted ? null : onTap,
    );
  }
}
