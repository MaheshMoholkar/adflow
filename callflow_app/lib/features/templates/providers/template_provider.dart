import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/database/app_database.dart';

final templatesStreamProvider = StreamProvider<List<Template>>((ref) {
  final db = ref.watch(databaseProvider);
  return db.watchTemplates();
});

final templatesByChannelProvider =
    FutureProvider.family<List<Template>, String>((ref, channel) {
  final db = ref.watch(databaseProvider);
  return db.getTemplatesByChannel(channel);
});
