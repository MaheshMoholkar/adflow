package com.callflow.service

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.util.Log

class PhoneStateReceiver : BroadcastReceiver() {
    override fun onReceive(context: Context?, intent: Intent?) {
        // The foreground service handles phone state changes directly.
        // This receiver exists as a fallback and to ensure the service
        // is restarted if it was killed.
        if (context != null && !CallDetectionService.isRunning) {
            val prefs = context.getSharedPreferences("callflow_prefs", Context.MODE_PRIVATE)
            val serviceEnabled = prefs.getBoolean("service_enabled", false)
            if (serviceEnabled) {
                Log.d("PhoneStateReceiver", "Service not running, restarting")
                ForegroundServiceManager(context).startService()
            }
        }
    }
}
