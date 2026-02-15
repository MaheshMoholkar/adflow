import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/database/app_database.dart';

final ruleConfigProvider = StreamProvider<Rule?>((ref) {
  final db = ref.watch(databaseProvider);
  return db.watchRule();
});

final smsTemplatesProvider = StreamProvider<List<Template>>((ref) {
  final db = ref.watch(databaseProvider);
  return db.watchTemplates().map(
        (templates) => templates
            .where((t) => t.channel == 'sms' || t.channel == 'both')
            .toList(),
      );
});
