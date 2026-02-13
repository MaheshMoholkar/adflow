package com.callflow.service

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.util.Log

class BootReceiver : BroadcastReceiver() {
    override fun onReceive(context: Context?, intent: Intent?) {
        if (intent?.action == Intent.ACTION_BOOT_COMPLETED && context != null) {
            Log.d("BootReceiver", "Boot completed, restarting call detection service")
            val prefs = context.getSharedPreferences("callflow_prefs", Context.MODE_PRIVATE)
            val serviceEnabled = prefs.getBoolean("service_enabled", false)
            if (serviceEnabled) {
                ForegroundServiceManager(context).startService()
            }
        }
    }
}
