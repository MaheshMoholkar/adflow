import 'dart:async';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class NativeBridge {
  static const _methodChannel = MethodChannel('com.callflow/methods');
  static const _eventChannel = EventChannel('com.callflow/call_events');

  Stream<Map<String, dynamic>>? _eventStream;

  // --- Service control ---

  Future<bool> startCallDetection() async {
    final result =
        await _methodChannel.invokeMethod<bool>('startCallDetection');
    return result ?? false;
  }

  Future<bool> stopCallDetection() async {
    final result = await _methodChannel.invokeMethod<bool>('stopCallDetection');
    return result ?? false;
  }

  Future<bool> isServiceRunning() async {
    final result = await _methodChannel.invokeMethod<bool>('isServiceRunning');
    return result ?? false;
  }

  // --- Rule config ---

  Future<bool> updateRuleConfig(String configJson) async {
    final result = await _methodChannel.invokeMethod<bool>(
      'updateRuleConfig',
      {'config': configJson},
    );
    return result ?? false;
  }

  // --- SMS ---

  Future<Map<String, dynamic>> sendSms({
    required String phone,
    required String message,
    int simSlot = 0,
  }) async {
    final result = await _methodChannel.invokeMethod<Map>(
      'sendSms',
      {'phone': phone, 'message': message, 'simSlot': simSlot},
    );
    return Map<String, dynamic>.from(result ?? {});
  }

  // --- SIM cards ---

  Future<List<Map<String, dynamic>>> getSimCards() async {
    final result = await _methodChannel.invokeMethod<List>('getSimCards');
    return result?.map((e) => Map<String, dynamic>.from(e as Map)).toList() ??
        [];
  }

  // --- Battery optimization ---

  Future<bool> isBatteryOptimizationDisabled() async {
    final result = await _methodChannel
        .invokeMethod<bool>('isBatteryOptimizationDisabled');
    return result ?? false;
  }

  Future<void> requestBatteryOptimization() async {
    await _methodChannel.invokeMethod('requestBatteryOptimization');
  }

  // --- Call event stream ---

  Stream<Map<String, dynamic>> get callEventStream {
    _eventStream ??= _eventChannel
        .receiveBroadcastStream()
        .map((event) => Map<String, dynamic>.from(event as Map));
    return _eventStream!;
  }
}

final nativeBridgeProvider = Provider<NativeBridge>((ref) {
  return NativeBridge();
});
