import 'package:drift/drift.dart';

class Users extends Table {
  IntColumn get id => integer()();
  TextColumn get phone => text()();
  TextColumn get name => text().withDefault(const Constant(''))();
  TextColumn get businessName => text().withDefault(const Constant(''))();
  TextColumn get city => text().withDefault(const Constant(''))();
  TextColumn get address => text().withDefault(const Constant(''))();
  TextColumn get locationUrl => text().withDefault(const Constant(''))();
  TextColumn get plan => text().withDefault(const Constant('none'))();
  DateTimeColumn get planStartedAt => dateTime().nullable()();
  DateTimeColumn get planExpiresAt => dateTime().nullable()();
  TextColumn get status => text().withDefault(const Constant('active'))();

  @override
  Set<Column> get primaryKey => {id};
}

class Templates extends Table {
  IntColumn get id => integer().autoIncrement()();
  IntColumn get serverId => integer().nullable()();
  TextColumn get name => text()();
  TextColumn get body => text()();
  TextColumn get type => text()();
  TextColumn get channel => text()();
  TextColumn get imagePath => text().nullable()();
  TextColumn get language => text().withDefault(const Constant('en'))();
  BoolColumn get isDefault => boolean().withDefault(const Constant(false))();
  TextColumn get source =>
      text().withDefault(const Constant('local'))(); // 'server' or 'local'
  BoolColumn get isSynced => boolean().withDefault(const Constant(false))();
  DateTimeColumn get createdAt =>
      dateTime().withDefault(currentDateAndTime)();
  DateTimeColumn get updatedAt =>
      dateTime().withDefault(currentDateAndTime)();
}

class Rules extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get configJson => text()();
  BoolColumn get isSynced => boolean().withDefault(const Constant(false))();
  DateTimeColumn get updatedAt =>
      dateTime().withDefault(currentDateAndTime)();
}

class CallEvents extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get eventId => text()();
  TextColumn get phone => text()();
  TextColumn get contactName => text().withDefault(const Constant(''))();
  TextColumn get direction => text()();
  IntColumn get durationSeconds => integer().withDefault(const Constant(0))();
  DateTimeColumn get callTimestamp => dateTime()();
  BoolColumn get isSynced => boolean().withDefault(const Constant(false))();
  DateTimeColumn get createdAt =>
      dateTime().withDefault(currentDateAndTime)();
}

class MessageLogs extends Table {
  IntColumn get id => integer().autoIncrement()();
  IntColumn get callEventId =>
      integer().references(CallEvents, #id)();
  IntColumn get templateId => integer().nullable()();
  TextColumn get channel => text()();
  TextColumn get status => text()();
  TextColumn get sendMethod => text().withDefault(const Constant(''))();
  IntColumn get simSlot => integer().nullable()();
  IntColumn get smsParts => integer().nullable()();
  TextColumn get errorMessage => text().withDefault(const Constant(''))();
  DateTimeColumn get sentAt => dateTime().nullable()();
  DateTimeColumn get createdAt =>
      dateTime().withDefault(currentDateAndTime)();
}
