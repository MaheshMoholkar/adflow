package com.callflow.service

import android.content.Context
import android.content.Intent
import android.os.Build

class ForegroundServiceManager(private val context: Context) {

    fun startService() {
        val intent = Intent(context, CallDetectionService::class.java)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            context.startForegroundService(intent)
        } else {
            context.startService(intent)
        }
    }

    fun stopService() {
        val intent = Intent(context, CallDetectionService::class.java)
        context.stopService(intent)
    }

    fun isServiceRunning(): Boolean {
        return CallDetectionService.isRunning
    }
}
