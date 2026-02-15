import 'dart:io';

import 'package:drift/drift.dart';
import 'package:drift/native.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';

import 'tables.dart';

part 'app_database.g.dart';

@DriftDatabase(tables: [Users, Templates, Rules, CallEvents, MessageLogs])
class AppDatabase extends _$AppDatabase {
  AppDatabase() : super(_openConnection());

  @override
  int get schemaVersion => 1;

  // --- User queries ---

  Future<User?> getUser() async {
    return (select(users)..limit(1)).getSingleOrNull();
  }

  Future<void> upsertUser(UsersCompanion user) async {
    await into(users).insertOnConflictUpdate(user);
  }

  Future<void> clearUser() async {
    await delete(users).go();
  }

  // --- Template queries ---

  Stream<List<Template>> watchTemplates() {
    return (select(templates)
          ..orderBy([
            (t) => OrderingTerm.desc(t.createdAt),
          ]))
        .watch();
  }

  Future<List<Template>> getTemplates() {
    return (select(templates)
          ..orderBy([
            (t) => OrderingTerm.desc(t.createdAt),
          ]))
        .get();
  }

  Future<List<Template>> getTemplatesByChannel(String channel) {
    return (select(templates)
          ..where((t) => t.channel.equals(channel) | t.channel.equals('both')))
        .get();
  }

  Future<Template> insertTemplate(TemplatesCompanion template) {
    return into(templates).insertReturning(template);
  }

  Future<bool> updateTemplate(TemplatesCompanion template) {
    return (update(templates)..where((t) => t.id.equals(template.id.value)))
        .write(template)
        .then((rows) => rows > 0);
  }

  Future<int> deleteTemplate(int id) {
    return (delete(templates)..where((t) => t.id.equals(id))).go();
  }

  Future<void> replaceServerTemplates(
      List<TemplatesCompanion> serverTemplates) async {
    await transaction(() async {
      await (delete(templates)..where((t) => t.source.equals('server'))).go();
      for (final tmpl in serverTemplates) {
        await into(templates).insert(tmpl);
      }
    });
  }

  // --- Rule queries ---

  Future<Rule?> getRule() async {
    return (select(rules)..limit(1)).getSingleOrNull();
  }

  Stream<Rule?> watchRule() {
    return (select(rules)..limit(1)).watchSingleOrNull();
  }

  Future<void> upsertRule(RulesCompanion rule) async {
    final existing = await getRule();
    if (existing != null) {
      await (update(rules)..where((r) => r.id.equals(existing.id))).write(rule);
    } else {
      await into(rules).insert(rule);
    }
  }

  // --- CallEvent queries ---

  Stream<List<CallEvent>> watchCallEvents({int limit = 50, int offset = 0}) {
    return (select(callEvents)
          ..orderBy([(e) => OrderingTerm.desc(e.callTimestamp)])
          ..limit(limit, offset: offset))
        .watch();
  }

  Future<List<CallEvent>> getCallEvents({
    int limit = 50,
    int offset = 0,
    String? direction,
  }) {
    final query = select(callEvents)
      ..orderBy([(e) => OrderingTerm.desc(e.callTimestamp)])
      ..limit(limit, offset: offset);
    if (direction != null) {
      query.where((e) => e.direction.equals(direction));
    }
    return query.get();
  }

  Future<CallEvent> insertCallEvent(CallEventsCompanion event) {
    return into(callEvents).insertReturning(event);
  }

  Future<List<CallEvent>> getUnsyncedEvents({int limit = 100}) {
    return (select(callEvents)
          ..where((e) => e.isSynced.equals(false))
          ..limit(limit))
        .get();
  }

  Future<void> markEventsSynced(List<int> ids) async {
    await (update(callEvents)..where((e) => e.id.isIn(ids)))
        .write(const CallEventsCompanion(isSynced: Value(true)));
  }

  // --- MessageLog queries ---

  Future<List<MessageLog>> getMessageLogsForEvent(int callEventId) {
    return (select(messageLogs)
          ..where((m) => m.callEventId.equals(callEventId)))
        .get();
  }

  Stream<List<MessageLog>> watchMessageLogsForEvent(int callEventId) {
    return (select(messageLogs)
          ..where((m) => m.callEventId.equals(callEventId)))
        .watch();
  }

  Future<MessageLog> insertMessageLog(MessageLogsCompanion log) {
    return into(messageLogs).insertReturning(log);
  }

  // --- Stats queries ---

  Future<int> countEventsToday() async {
    final now = DateTime.now();
    final startOfDay = DateTime(now.year, now.month, now.day);
    final result = await (select(callEvents)
          ..where((e) => e.callTimestamp.isBiggerOrEqualValue(startOfDay)))
        .get();
    return result.length;
  }

  Stream<int> watchEventsTodayCount() {
    final now = DateTime.now();
    final startOfDay = DateTime(now.year, now.month, now.day);
    return (select(callEvents)
          ..where((e) => e.callTimestamp.isBiggerOrEqualValue(startOfDay)))
        .watch()
        .map((rows) => rows.length);
  }

  Future<int> countMessagesByChannel(String channel) async {
    final now = DateTime.now();
    final startOfDay = DateTime(now.year, now.month, now.day);
    final result = await (select(messageLogs)
          ..where((m) =>
              m.channel.equals(channel) &
              m.sentAt.isBiggerOrEqualValue(startOfDay)))
        .get();
    return result.length;
  }

  Stream<int> watchMessagesByChannelCount(String channel) {
    final now = DateTime.now();
    final startOfDay = DateTime(now.year, now.month, now.day);
    return (select(messageLogs)
          ..where((m) =>
              m.channel.equals(channel) &
              m.sentAt.isBiggerOrEqualValue(startOfDay)))
        .watch()
        .map((rows) => rows.length);
  }

  Future<double> successRate() async {
    final now = DateTime.now();
    final startOfDay = DateTime(now.year, now.month, now.day);
    final all = await (select(messageLogs)
          ..where((m) => m.sentAt.isBiggerOrEqualValue(startOfDay)))
        .get();
    if (all.isEmpty) return 0.0;
    final sent =
        all.where((m) => m.status == 'sent' || m.status == 'delivered');
    return sent.length / all.length;
  }

  Stream<double> watchSuccessRate() {
    final now = DateTime.now();
    final startOfDay = DateTime(now.year, now.month, now.day);
    return (select(messageLogs)
          ..where((m) => m.sentAt.isBiggerOrEqualValue(startOfDay)))
        .watch()
        .map((all) {
      if (all.isEmpty) return 0.0;
      final sent =
          all.where((m) => m.status == 'sent' || m.status == 'delivered');
      return sent.length / all.length;
    });
  }

  // --- Cleanup ---

  Future<void> clearAll() async {
    await transaction(() async {
      await delete(messageLogs).go();
      await delete(callEvents).go();
      await delete(templates).go();
      await delete(rules).go();
      await delete(users).go();
    });
  }
}

LazyDatabase _openConnection() {
  return LazyDatabase(() async {
    final dbFolder = await getApplicationDocumentsDirectory();
    final file = File(p.join(dbFolder.path, 'callflow.db'));
    return NativeDatabase.createInBackground(file);
  });
}

final databaseProvider = Provider<AppDatabase>((ref) {
  final db = AppDatabase();
  ref.onDispose(() => db.close());
  return db;
});
