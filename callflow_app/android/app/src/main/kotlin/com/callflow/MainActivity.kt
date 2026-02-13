package com.callflow

import android.os.Bundle
import com.callflow.bridge.CallEventStreamHandler
import com.callflow.bridge.NativeMethodHandler
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.EventChannel
import io.flutter.plugin.common.MethodChannel

class MainActivity : FlutterActivity() {
    companion object {
        const val METHOD_CHANNEL = "com.callflow/methods"
        const val EVENT_CHANNEL = "com.callflow/call_events"
    }

    private lateinit var methodHandler: NativeMethodHandler
    private lateinit var eventStreamHandler: CallEventStreamHandler

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)

        eventStreamHandler = CallEventStreamHandler.getInstance()

        methodHandler = NativeMethodHandler(this)

        MethodChannel(
            flutterEngine.dartExecutor.binaryMessenger,
            METHOD_CHANNEL
        ).setMethodCallHandler(methodHandler)

        EventChannel(
            flutterEngine.dartExecutor.binaryMessenger,
            EVENT_CHANNEL
        ).setStreamHandler(eventStreamHandler)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
    }
}
