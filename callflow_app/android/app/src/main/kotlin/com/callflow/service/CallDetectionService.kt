package com.callflow.service

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.app.Service
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.content.IntentFilter
import android.os.Build
import android.os.IBinder
import android.os.PowerManager
import android.telephony.TelephonyManager
import android.util.Log
import com.callflow.MainActivity
import com.callflow.bridge.CallEventStreamHandler
import com.callflow.messaging.ChannelRouter
import com.callflow.messaging.SmsModule
import com.callflow.rules.LocalRuleEngine

class CallDetectionService : Service() {

    companion object {
        const val TAG = "CallDetectionService"
        const val CHANNEL_ID = "callflow_detection"
        const val NOTIFICATION_ID = 1001

        @Volatile
        var isRunning = false
            private set

        private var ruleConfigJson: String? = null
        private var instance: CallDetectionService? = null
        private const val PREFS_NAME = "callflow_rule_config"
        private const val PREFS_KEY = "rule_config_json"

        fun updateRuleConfig(json: String) {
            ruleConfigJson = json
            // Persist to SharedPreferences
            instance?.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)
                ?.edit()?.putString(PREFS_KEY, json)?.apply()
            // Update the live service instance if running
            instance?.ruleEngine?.updateConfig(json)
        }
    }

    private var wakeLock: PowerManager.WakeLock? = null
    private var lastState = TelephonyManager.CALL_STATE_IDLE
    private var callStartTime: Long = 0
    private var incomingNumber: String? = null
    private var isIncoming = false
    private val callLogReader = CallLogReader()

    private lateinit var ruleEngine: LocalRuleEngine
    private lateinit var channelRouter: ChannelRouter

    private val phoneStateReceiver = object : BroadcastReceiver() {
        override fun onReceive(context: Context?, intent: Intent?) {
            if (intent?.action == TelephonyManager.ACTION_PHONE_STATE_CHANGED) {
                val state = intent.getStringExtra(TelephonyManager.EXTRA_STATE) ?: return
                val number = intent.getStringExtra(TelephonyManager.EXTRA_INCOMING_NUMBER)

                onPhoneStateChanged(state, number)
            }
        }
    }

    override fun onCreate() {
        super.onCreate()

        ruleEngine = LocalRuleEngine()
        // Load config: prefer in-memory, fallback to SharedPreferences
        val configToLoad = ruleConfigJson
            ?: getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE).getString(PREFS_KEY, null)
        configToLoad?.let { ruleEngine.updateConfig(it) }
        Log.d(TAG, "Rule config loaded: ${configToLoad != null}")

        val smsModule = SmsModule(this)
        channelRouter = ChannelRouter(this, smsModule, ruleEngine)

        createNotificationChannel()
        acquireWakeLock()

        val filter = IntentFilter(TelephonyManager.ACTION_PHONE_STATE_CHANGED)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            registerReceiver(phoneStateReceiver, filter, RECEIVER_NOT_EXPORTED)
        } else {
            registerReceiver(phoneStateReceiver, filter)
        }

        instance = this
        isRunning = true
        Log.d(TAG, "Call detection service created")
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        val notification = buildNotification()
        startForeground(NOTIFICATION_ID, notification)
        return START_STICKY
    }

    override fun onDestroy() {
        super.onDestroy()
        instance = null
        isRunning = false
        try {
            unregisterReceiver(phoneStateReceiver)
        } catch (e: IllegalArgumentException) {
            // Receiver not registered
        }
        releaseWakeLock()
        Log.d(TAG, "Call detection service destroyed")
    }

    override fun onBind(intent: Intent?): IBinder? = null

    private fun onPhoneStateChanged(state: String, number: String?) {
        val newState = when (state) {
            TelephonyManager.EXTRA_STATE_IDLE -> TelephonyManager.CALL_STATE_IDLE
            TelephonyManager.EXTRA_STATE_RINGING -> TelephonyManager.CALL_STATE_RINGING
            TelephonyManager.EXTRA_STATE_OFFHOOK -> TelephonyManager.CALL_STATE_OFFHOOK
            else -> return
        }

        when {
            // Incoming call ringing
            lastState == TelephonyManager.CALL_STATE_IDLE
                    && newState == TelephonyManager.CALL_STATE_RINGING -> {
                isIncoming = true
                incomingNumber = number
                callStartTime = System.currentTimeMillis()
                Log.d(TAG, "Incoming call ringing: $number")
            }
            // Incoming call answered
            lastState == TelephonyManager.CALL_STATE_RINGING
                    && newState == TelephonyManager.CALL_STATE_OFFHOOK -> {
                callStartTime = System.currentTimeMillis()
                Log.d(TAG, "Incoming call answered")
            }
            // Outgoing call started
            lastState == TelephonyManager.CALL_STATE_IDLE
                    && newState == TelephonyManager.CALL_STATE_OFFHOOK -> {
                isIncoming = false
                callStartTime = System.currentTimeMillis()
                Log.d(TAG, "Outgoing call started")
            }
            // Call ended (IDLE after non-IDLE)
            newState == TelephonyManager.CALL_STATE_IDLE
                    && lastState != TelephonyManager.CALL_STATE_IDLE -> {
                onCallEnded()
            }
        }

        lastState = newState
    }

    private fun onCallEnded() {
        val wasMissed = lastState == TelephonyManager.CALL_STATE_RINGING
        val direction = when {
            wasMissed -> "missed"
            isIncoming -> "incoming"
            else -> "outgoing"
        }
        val durationSeconds = if (wasMissed) 0 else
            ((System.currentTimeMillis() - callStartTime) / 1000).toInt()

        Log.d(TAG, "Call ended: direction=$direction, duration=$durationSeconds")

        // Wait briefly for call log to be updated, then read it
        android.os.Handler(mainLooper).postDelayed({
            processCallEnd(direction, durationSeconds)
        }, 1500)
    }

    private fun processCallEnd(direction: String, durationSeconds: Int) {
        try {
            Log.d(TAG, "processCallEnd: direction=$direction, incomingNumber=$incomingNumber")
            val callInfo = callLogReader.getLatestCall(this)
            Log.d(TAG, "processCallEnd: callInfo=${callInfo?.phone ?: "null"}")
            val phone = callInfo?.phone ?: incomingNumber ?: ""
            val contactName = callInfo?.contactName ?: ""
            val actualDuration = callInfo?.duration ?: durationSeconds

            if (phone.isEmpty()) {
                Log.d(TAG, "No phone number available for call event")
                return
            }

            val eventData = mapOf<String, Any?>(
                "type" to "call_event",
                "phone" to phone,
                "contact_name" to contactName,
                "direction" to direction,
                "duration_seconds" to actualDuration,
                "call_timestamp" to System.currentTimeMillis(),
                "event_id" to java.util.UUID.randomUUID().toString()
            )

            // Stream to Flutter
            CallEventStreamHandler.getInstance().sendCallEvent(eventData)

            // Process through channel router
            val eventJson = org.json.JSONObject(eventData).toString()
            channelRouter.processCallEvent(eventJson)
        } catch (e: Exception) {
            Log.e(TAG, "Error processing call end", e)
        }

        // Reset state
        incomingNumber = null
        isIncoming = false
        callStartTime = 0
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                "Call Detection",
                NotificationManager.IMPORTANCE_LOW
            ).apply {
                description = "Monitors calls to send automated SMS"
                setShowBadge(false)
            }
            val manager = getSystemService(NotificationManager::class.java)
            manager.createNotificationChannel(channel)
        }
    }

    private fun buildNotification(): Notification {
        val intent = Intent(this, MainActivity::class.java)
        val pendingIntent = PendingIntent.getActivity(
            this, 0, intent,
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
        )

        return if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            Notification.Builder(this, CHANNEL_ID)
                .setContentTitle("CallFlow Active")
                .setContentText("Monitoring calls for automated messages")
                .setSmallIcon(android.R.drawable.ic_menu_call)
                .setContentIntent(pendingIntent)
                .setOngoing(true)
                .build()
        } else {
            @Suppress("DEPRECATION")
            Notification.Builder(this)
                .setContentTitle("CallFlow Active")
                .setContentText("Monitoring calls for automated messages")
                .setSmallIcon(android.R.drawable.ic_menu_call)
                .setContentIntent(pendingIntent)
                .setOngoing(true)
                .build()
        }
    }

    private fun acquireWakeLock() {
        val pm = getSystemService(Context.POWER_SERVICE) as PowerManager
        wakeLock = pm.newWakeLock(
            PowerManager.PARTIAL_WAKE_LOCK,
            "CallFlow::CallDetectionWakeLock"
        ).apply {
            acquire(10 * 60 * 1000L) // 10 minutes max
        }
    }

    private fun releaseWakeLock() {
        wakeLock?.let {
            if (it.isHeld) it.release()
        }
        wakeLock = null
    }
}
