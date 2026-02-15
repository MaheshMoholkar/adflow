package com.callflow.bridge

import android.os.Handler
import android.os.Looper
import io.flutter.plugin.common.EventChannel
import org.json.JSONObject

class CallEventStreamHandler private constructor() : EventChannel.StreamHandler {

    companion object {
        @Volatile
        private var instance: CallEventStreamHandler? = null

        fun getInstance(): CallEventStreamHandler {
            return instance ?: synchronized(this) {
                instance ?: CallEventStreamHandler().also { instance = it }
            }
        }
    }

    private var eventSink: EventChannel.EventSink? = null
    private val mainHandler = Handler(Looper.getMainLooper())

    override fun onListen(arguments: Any?, events: EventChannel.EventSink?) {
        eventSink = events
    }

    override fun onCancel(arguments: Any?) {
        eventSink = null
    }

    fun sendCallEvent(eventData: Map<String, Any?>) {
        mainHandler.post {
            eventSink?.success(eventData)
        }
    }

    fun sendMessageLog(logData: Map<String, Any?>) {
        mainHandler.post {
            eventSink?.success(logData)
        }
    }

    fun sendError(code: String, message: String) {
        mainHandler.post {
            eventSink?.error(code, message, null)
        }
    }
}
