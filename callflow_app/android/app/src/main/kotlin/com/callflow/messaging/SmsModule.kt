package com.callflow.messaging

import android.app.Activity
import android.app.PendingIntent
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.content.IntentFilter
import android.os.Build
import android.telephony.SmsManager
import android.telephony.SubscriptionManager
import android.util.Log

class SmsModule(private val context: Context) {

    companion object {
        const val TAG = "SmsModule"
        const val SMS_SENT_ACTION = "com.callflow.SMS_SENT"
        const val SMS_DELIVERED_ACTION = "com.callflow.SMS_DELIVERED"
    }

    fun sendSms(
        phone: String,
        message: String,
        simSlot: Int = 0,
        callback: (success: Boolean, error: String?) -> Unit
    ) {
        try {
            val smsManager = getSmsManager(simSlot)

            val sentIntent = PendingIntent.getBroadcast(
                context, 0,
                Intent(SMS_SENT_ACTION),
                PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
            )

            val deliveryIntent = PendingIntent.getBroadcast(
                context, 0,
                Intent(SMS_DELIVERED_ACTION),
                PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
            )

            // Register sent callback
            val sentReceiver = object : BroadcastReceiver() {
                override fun onReceive(ctx: Context?, intent: Intent?) {
                    try {
                        context.unregisterReceiver(this)
                    } catch (_: Exception) {}
                    when (resultCode) {
                        Activity.RESULT_OK -> {
                            Log.d(TAG, "SMS sent successfully to $phone")
                            callback(true, null)
                        }
                        else -> {
                            val error = "SMS send failed with code: $resultCode"
                            Log.e(TAG, error)
                            callback(false, error)
                        }
                    }
                }
            }

            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
                context.registerReceiver(
                    sentReceiver,
                    IntentFilter(SMS_SENT_ACTION),
                    Context.RECEIVER_NOT_EXPORTED
                )
            } else {
                context.registerReceiver(sentReceiver, IntentFilter(SMS_SENT_ACTION))
            }

            if (message.length > 160) {
                // Multipart SMS
                val parts = smsManager.divideMessage(message)
                val sentIntents = ArrayList<PendingIntent>(parts.size)
                val deliveryIntents = ArrayList<PendingIntent>(parts.size)
                for (i in parts.indices) {
                    sentIntents.add(sentIntent)
                    deliveryIntents.add(deliveryIntent)
                }
                smsManager.sendMultipartTextMessage(
                    phone, null, parts, sentIntents, deliveryIntents
                )
                Log.d(TAG, "Sending multipart SMS (${parts.size} parts) to $phone")
            } else {
                smsManager.sendTextMessage(
                    phone, null, message, sentIntent, deliveryIntent
                )
                Log.d(TAG, "Sending SMS to $phone")
            }
        } catch (e: Exception) {
            Log.e(TAG, "Error sending SMS", e)
            callback(false, e.message)
        }
    }

    fun getSmsParts(message: String): Int {
        return if (message.length <= 160) 1
        else {
            val smsManager = getSmsManager(0)
            smsManager.divideMessage(message).size
        }
    }

    @Suppress("DEPRECATION")
    private fun getSmsManager(simSlot: Int): SmsManager {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
            // Android 12+
            if (simSlot > 0) {
                try {
                    val subManager = context.getSystemService(Context.TELEPHONY_SUBSCRIPTION_SERVICE) as SubscriptionManager
                    val subscriptions = subManager.activeSubscriptionInfoList
                    if (subscriptions != null && simSlot < subscriptions.size) {
                        val subId = subscriptions[simSlot].subscriptionId
                        return context.getSystemService(SmsManager::class.java)
                            .createForSubscriptionId(subId)
                    }
                } catch (e: SecurityException) {
                    Log.w(TAG, "No permission for subscription info", e)
                }
            }
            return context.getSystemService(SmsManager::class.java)
        } else {
            // Below Android 12
            if (simSlot > 0) {
                try {
                    val subManager = context.getSystemService(Context.TELEPHONY_SUBSCRIPTION_SERVICE) as SubscriptionManager
                    val subscriptions = subManager.activeSubscriptionInfoList
                    if (subscriptions != null && simSlot < subscriptions.size) {
                        val subId = subscriptions[simSlot].subscriptionId
                        return SmsManager.getSmsManagerForSubscriptionId(subId)
                    }
                } catch (e: SecurityException) {
                    Log.w(TAG, "No permission for subscription info", e)
                }
            }
            return SmsManager.getDefault()
        }
    }
}
