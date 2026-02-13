import 'package:drift/drift.dart' as drift;
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/database/app_database.dart';
import '../../../core/native/native_bridge.dart';
import '../../../core/network/api_client.dart';
import '../../../core/providers/core_providers.dart';
import '../../auth/providers/auth_provider.dart';

class SettingsScreen extends ConsumerStatefulWidget {
  const SettingsScreen({super.key});

  @override
  ConsumerState<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends ConsumerState<SettingsScreen> {
  final _businessNameController = TextEditingController();
  bool _editingName = false;
  bool _savingName = false;
  List<Map<String, dynamic>> _simCards = [];
  int _selectedSim = 0;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  @override
  void dispose() {
    _businessNameController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    final bridge = ref.read(nativeBridgeProvider);

    try {
      _simCards = await bridge.getSimCards();
    } catch (_) {}

    if (mounted) setState(() {});
  }

  Future<void> _saveBusinessName() async {
    setState(() => _savingName = true);
    try {
      final api = ref.read(apiClientProvider);
      await api.put('/user/profile', data: {
        'business_name': _businessNameController.text.trim(),
      });

      final db = ref.read(databaseProvider);
      final user = await db.getUser();
      if (user != null) {
        await db.upsertUser(UsersCompanion(
          id: drift.Value(user.id),
          businessName: drift.Value(_businessNameController.text.trim()),
        ));
      }

      setState(() => _editingName = false);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Business name updated')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _savingName = false);
    }
  }

  Future<void> _logout() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Logout'),
        content: const Text('Are you sure you want to logout?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Logout'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      ref.read(serviceRunningProvider.notifier).stop();
      await ref.read(authServiceProvider).logout();
      if (mounted) context.go('/auth/login');
    }
  }

  bool _isLifetime(DateTime? expiresAt) {
    if (expiresAt == null) return false;
    return expiresAt.year >= 2099;
  }

  @override
  Widget build(BuildContext context) {
    final user = ref.watch(currentUserProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Settings')),
      body: user.when(
        data: (u) {
          if (u == null) {
            return const Center(child: Text('Not logged in'));
          }

          if (!_editingName && _businessNameController.text.isEmpty) {
            _businessNameController.text = u.businessName;
          }

          final isLifetime = _isLifetime(u.planExpiresAt);

          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              // Profile section
              const _SectionHeader(icon: Icons.person_outline, title: 'Profile'),
              const SizedBox(height: 8),
              Card(
                child: Column(
                  children: [
                    ListTile(
                      leading: const Icon(Icons.phone),
                      title: const Text('Phone'),
                      subtitle: Text(u.phone),
                    ),
                    const Divider(height: 0, indent: 16, endIndent: 16),
                    ListTile(
                      leading: const Icon(Icons.business),
                      title: const Text('Business Name'),
                      subtitle: _editingName
                          ? Row(
                              children: [
                                Expanded(
                                  child: TextField(
                                    controller: _businessNameController,
                                    autofocus: true,
                                    decoration: const InputDecoration(
                                      isDense: true,
                                      hintText: 'Enter business name',
                                    ),
                                  ),
                                ),
                                IconButton(
                                  icon: _savingName
                                      ? const SizedBox(
                                          height: 16,
                                          width: 16,
                                          child: CircularProgressIndicator(
                                              strokeWidth: 2))
                                      : const Icon(Icons.check),
                                  onPressed:
                                      _savingName ? null : _saveBusinessName,
                                ),
                                IconButton(
                                  icon: const Icon(Icons.close),
                                  onPressed: () => setState(() {
                                    _editingName = false;
                                    _businessNameController.text =
                                        u.businessName;
                                  }),
                                ),
                              ],
                            )
                          : Text(u.businessName.isEmpty
                              ? 'Not set'
                              : u.businessName),
                      trailing: _editingName
                          ? null
                          : IconButton(
                              icon: const Icon(Icons.edit_outlined),
                              onPressed: () =>
                                  setState(() => _editingName = true),
                            ),
                    ),
                    if (u.city.isNotEmpty) ...[
                      const Divider(height: 0, indent: 16, endIndent: 16),
                      ListTile(
                        leading: const Icon(Icons.location_city),
                        title: const Text('City'),
                        subtitle: Text(u.city),
                      ),
                    ],
                  ],
                ),
              ),
              const SizedBox(height: 24),

              // Plan section
              const _SectionHeader(icon: Icons.workspace_premium, title: 'Plan'),
              const SizedBox(height: 8),
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 12, vertical: 6),
                        decoration: BoxDecoration(
                          color: _planColor(context, u.plan),
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: Text(
                          u.plan.toUpperCase(),
                          style: const TextStyle(
                            color: Colors.white,
                            fontWeight: FontWeight.bold,
                            fontSize: 12,
                          ),
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            if (u.plan == 'none')
                              Text(
                                'No active plan',
                                style: theme.textTheme.bodyMedium?.copyWith(
                                  color: theme.colorScheme.onSurfaceVariant,
                                ),
                              )
                            else if (isLifetime)
                              Text(
                                'Lifetime',
                                style: theme.textTheme.bodyLarge?.copyWith(
                                  fontWeight: FontWeight.w600,
                                  color: theme.colorScheme.primary,
                                ),
                              )
                            else if (u.planExpiresAt != null) ...[
                              Text(
                                'Expires ${DateFormat('dd MMM yyyy').format(u.planExpiresAt!)}',
                                style: theme.textTheme.bodyMedium,
                              ),
                            ] else
                              Text(
                                'No expiry set',
                                style: theme.textTheme.bodyMedium?.copyWith(
                                  color: theme.colorScheme.onSurfaceVariant,
                                ),
                              ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),

              // SIM selection
              if (_simCards.isNotEmpty) ...[
                const _SectionHeader(
                    icon: Icons.sim_card_outlined, title: 'SIM Card'),
                const SizedBox(height: 8),
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: DropdownButtonFormField<int>(
                      initialValue: _selectedSim,
                      decoration: const InputDecoration(
                        labelText: 'Select SIM for SMS',
                        isDense: true,
                        prefixIcon: Icon(Icons.sim_card, size: 20),
                      ),
                      items: _simCards.asMap().entries.map((entry) {
                        return DropdownMenuItem(
                          value: entry.key,
                          child: Text(
                            '${entry.value['displayName']} (${entry.value['carrierName']})',
                          ),
                        );
                      }).toList(),
                      onChanged: (v) => setState(() => _selectedSim = v ?? 0),
                    ),
                  ),
                ),
                const SizedBox(height: 24),
              ],

              // Logout
              const SizedBox(height: 8),
              OutlinedButton.icon(
                onPressed: _logout,
                icon: const Icon(Icons.logout),
                label: const Text('Logout'),
                style: OutlinedButton.styleFrom(
                  foregroundColor: theme.colorScheme.error,
                  side: BorderSide(color: theme.colorScheme.error),
                  minimumSize: const Size(double.infinity, 48),
                ),
              ),
              const SizedBox(height: 32),
            ],
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Error: $e')),
      ),
    );
  }

  Color _planColor(BuildContext context, String plan) {
    switch (plan) {
      case 'sms':
        return Theme.of(context).colorScheme.tertiary;
      default:
        return Colors.grey;
    }
  }
}

class _SectionHeader extends StatelessWidget {
  final IconData icon;
  final String title;

  const _SectionHeader({required this.icon, required this.title});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon,
            size: 20, color: Theme.of(context).colorScheme.onSurfaceVariant),
        const SizedBox(width: 8),
        Text(title, style: Theme.of(context).textTheme.titleMedium),
      ],
    );
  }
}
