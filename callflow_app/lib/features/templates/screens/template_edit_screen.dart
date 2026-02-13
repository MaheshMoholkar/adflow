import 'dart:io';

import 'package:drift/drift.dart' as drift;
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';
import '../../../core/database/app_database.dart';

class TemplateEditScreen extends ConsumerStatefulWidget {
  final int? templateId;

  const TemplateEditScreen({super.key, this.templateId});

  @override
  ConsumerState<TemplateEditScreen> createState() => _TemplateEditScreenState();
}

class _TemplateEditScreenState extends ConsumerState<TemplateEditScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _bodyController = TextEditingController();
  String _channel = 'sms';
  String _type = 'incoming';
  String? _imagePath;
  bool _isLoading = false;
  Template? _existing;

  static const _variables = [
    '{contact_name}',
    '{business_name}',
    '{phone_number}',
    '{call_duration}',
    '{date}',
    '{time}',
  ];

  @override
  void initState() {
    super.initState();
    if (widget.templateId != null) {
      _loadTemplate();
    }
  }

  Future<void> _loadTemplate() async {
    final db = ref.read(databaseProvider);
    final templates = await db.getTemplates();
    final template =
        templates.where((t) => t.id == widget.templateId).firstOrNull;
    if (template != null && mounted) {
      setState(() {
        _existing = template;
        _nameController.text = template.name;
        _bodyController.text = template.body;
        _channel = template.channel;
        _type = template.type;
        _imagePath = template.imagePath;
      });
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    _bodyController.dispose();
    super.dispose();
  }

  void _insertVariable(String variable) {
    final text = _bodyController.text;
    final selection = _bodyController.selection;
    final newText = text.replaceRange(
      selection.start,
      selection.end,
      variable,
    );
    _bodyController.value = TextEditingValue(
      text: newText,
      selection: TextSelection.collapsed(
        offset: selection.start + variable.length,
      ),
    );
  }

  int get _smsCharCount {
    final body = _bodyController.text;
    return body.length;
  }

  int get _smsParts {
    final len = _smsCharCount;
    if (len <= 160) return 1;
    return (len / 153).ceil();
  }

  bool get _showImagePicker => false; // SMS only, no image support

  Future<void> _pickImage() async {
    try {
      final picker = ImagePicker();
      final picked = await picker.pickImage(source: ImageSource.gallery);
      if (picked == null) return;

      final oldPath = _imagePath;

      // Copy to app's local storage so it persists
      final appDir = await getApplicationDocumentsDirectory();
      final imagesDir = Directory(p.join(appDir.path, 'template_images'));
      if (!imagesDir.existsSync()) {
        imagesDir.createSync(recursive: true);
      }
      final ext = p.extension(picked.path);
      final fileName = 'template_${DateTime.now().millisecondsSinceEpoch}$ext';
      final savedFile = await File(picked.path).copy(
        p.join(imagesDir.path, fileName),
      );

      // Delete the previous image file if it was replaced
      _deleteImageFile(oldPath);

      if (mounted) setState(() => _imagePath = savedFile.path);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to pick image: $e')),
        );
      }
    }
  }

  void _removeImage() {
    _deleteImageFile(_imagePath);
    setState(() => _imagePath = null);
  }

  void _deleteImageFile(String? path) {
    if (path == null) return;
    try {
      final file = File(path);
      if (file.existsSync()) file.deleteSync();
    } catch (_) {}
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isLoading = true);

    try {
      final db = ref.read(databaseProvider);
      // Clear imagePath if channel is SMS-only
      final effectiveImagePath = _channel == 'sms' ? null : _imagePath;

      if (_existing != null) {
        await db.updateTemplate(TemplatesCompanion(
          id: drift.Value(_existing!.id),
          name: drift.Value(_nameController.text.trim()),
          body: drift.Value(_bodyController.text),
          type: drift.Value(_type),
          channel: drift.Value(_channel),
          imagePath: drift.Value(effectiveImagePath),
          updatedAt: drift.Value(DateTime.now()),
        ));
      } else {
        await db.insertTemplate(TemplatesCompanion.insert(
          name: _nameController.text.trim(),
          body: _bodyController.text,
          type: _type,
          channel: _channel,
          imagePath: drift.Value(effectiveImagePath),
          source: const drift.Value('local'),
        ));
      }
      if (mounted) context.pop();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _delete() async {
    if (_existing == null) return;

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Template'),
        content: const Text('Are you sure you want to delete this template?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Delete'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final db = ref.read(databaseProvider);
      await db.deleteTemplate(_existing!.id);
      if (mounted) context.pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    final isEditing = _existing != null;
    const showSmsCounter = true; // Always show for SMS-only
    return Scaffold(
      appBar: AppBar(
        title: Text(isEditing ? 'Edit Template' : 'New Template'),
        actions: [
          if (isEditing && _existing?.source != 'server')
            IconButton(
              icon: const Icon(Icons.delete_outline),
              onPressed: _delete,
            ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            TextFormField(
              controller: _nameController,
              decoration: const InputDecoration(
                labelText: 'Template Name',
                hintText: 'e.g. Follow-up After Incoming Call',
              ),
              validator: (v) =>
                  v == null || v.trim().isEmpty ? 'Name is required' : null,
            ),
            const SizedBox(height: 16),

            const SizedBox(height: 16),

            // Image picker (disabled for SMS-only)
            if (_showImagePicker) ...[
              Text('Attach Image (optional)',
                  style: Theme.of(context).textTheme.labelLarge),
              const SizedBox(height: 8),
              if (_imagePath != null && File(_imagePath!).existsSync())
                Stack(
                  children: [
                    ClipRRect(
                      borderRadius: BorderRadius.circular(8),
                      child: Image.file(
                        File(_imagePath!),
                        height: 150,
                        width: double.infinity,
                        fit: BoxFit.cover,
                      ),
                    ),
                    Positioned(
                      top: 4,
                      right: 4,
                      child: IconButton.filled(
                        icon: const Icon(Icons.close, size: 18),
                        onPressed: _removeImage,
                        style: IconButton.styleFrom(
                          backgroundColor:
                              Theme.of(context).colorScheme.errorContainer,
                          foregroundColor:
                              Theme.of(context).colorScheme.onErrorContainer,
                        ),
                      ),
                    ),
                  ],
                )
              else
                OutlinedButton.icon(
                  onPressed: _pickImage,
                  icon: const Icon(Icons.image_outlined),
                  label: const Text('Pick Image'),
                ),
              const SizedBox(height: 16),
            ],

            // Type selector
            Text('Call Type', style: Theme.of(context).textTheme.labelLarge),
            const SizedBox(height: 8),
            SegmentedButton<String>(
              segments: const [
                ButtonSegment(value: 'all', label: Text('All')),
                ButtonSegment(value: 'incoming', label: Text('Incoming')),
                ButtonSegment(value: 'outgoing', label: Text('Outgoing')),
                ButtonSegment(value: 'missed', label: Text('Missed')),
              ],
              selected: {_type},
              onSelectionChanged: (s) => setState(() => _type = s.first),
            ),
            const SizedBox(height: 16),

            // Message body
            TextFormField(
              controller: _bodyController,
              maxLines: 6,
              decoration: InputDecoration(
                labelText: 'Message Body',
                hintText: 'Type your message template here...',
                alignLabelWithHint: true,
                counterText:
                    '$_smsCharCount/918 chars, $_smsParts part${_smsParts > 1 ? 's' : ''}',
              ),
              validator: (v) =>
                  v == null || v.trim().isEmpty ? 'Body is required' : null,
              onChanged: (_) => setState(() {}),
            ),
            if (showSmsCounter && _smsCharCount > 918)
              Padding(
                padding: const EdgeInsets.only(top: 4),
                child: Text(
                  'SMS body exceeds 918 character limit',
                  style: TextStyle(
                    color: Theme.of(context).colorScheme.error,
                    fontSize: 12,
                  ),
                ),
              ),
            const SizedBox(height: 12),

            // Variable chips
            Text('Insert Variable',
                style: Theme.of(context).textTheme.labelMedium),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: _variables
                  .map((v) => ActionChip(
                        label: Text(v),
                        onPressed: () => _insertVariable(v),
                      ))
                  .toList(),
            ),
            const SizedBox(height: 32),

            FilledButton(
              onPressed: _isLoading ? null : _save,
              child: _isLoading
                  ? const SizedBox(
                      height: 20,
                      width: 20,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : Text(isEditing ? 'Save Changes' : 'Create Template'),
            ),
          ],
        ),
      ),
    );
  }
}
