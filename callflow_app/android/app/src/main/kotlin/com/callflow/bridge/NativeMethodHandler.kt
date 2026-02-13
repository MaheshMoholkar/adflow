package com.callflow.bridge

import android.app.Activity
import android.content.Context
import android.os.Build
import android.os.PowerManager
import android.provider.Settings
import android.content.Intent
import android.telephony.SubscriptionManager
import com.callflow.messaging.ChannelRouter
import com.callflow.messaging.SmsModule
import com.callflow.rules.LocalRuleEngine
import com.callflow.service.CallDetectionService
import com.callflow.service.ForegroundServiceManager
import io.flutter.plugin.common.MethodCall
import io.flutter.plugin.common.MethodChannel

class NativeMethodHandler(private val activity: Activity) : MethodChannel.MethodCallHandler {

    private val serviceManager = ForegroundServiceManager(activity)
    private val smsModule = SmsModule(activity)
    private val ruleEngine = LocalRuleEngine()
    private val channelRouter = ChannelRouter(activity, smsModule, ruleEngine)

    override fun onMethodCall(call: MethodCall, result: MethodChannel.Result) {
        when (call.method) {
            "startCallDetection" -> {
                serviceManager.startService()
                result.success(true)
            }
            "stopCallDetection" -> {
                serviceManager.stopService()
                result.success(true)
            }
            "isServiceRunning" -> {
                result.success(serviceManager.isServiceRunning())
            }
            "updateRuleConfig" -> {
                val configJson = call.argument<String>("config")
                if (configJson != null) {
                    ruleEngine.updateConfig(configJson)
                    CallDetectionService.updateRuleConfig(configJson)
                    // Also persist directly in case service instance isn't alive yet
                    activity.getSharedPreferences("callflow_rule_config", android.content.Context.MODE_PRIVATE)
                        .edit().putString("rule_config_json", configJson).apply()
                    result.success(true)
                } else {
                    result.error("INVALID_ARGS", "config is required", null)
                }
            }
            "sendSms" -> {
                val phone = call.argument<String>("phone") ?: ""
                val message = call.argument<String>("message") ?: ""
                val simSlot = call.argument<Int>("simSlot") ?: 0
                smsModule.sendSms(phone, message, simSlot) { success, error ->
                    if (success) {
                        result.success(mapOf("status" to "sent"))
                    } else {
                        result.success(mapOf("status" to "failed", "error" to error))
                    }
                }
            }
            "getSimCards" -> {
                val sims = getSimCards()
                result.success(sims)
            }
            "isBatteryOptimizationDisabled" -> {
                val pm = activity.getSystemService(Context.POWER_SERVICE) as PowerManager
                result.success(pm.isIgnoringBatteryOptimizations(activity.packageName))
            }
            "requestBatteryOptimization" -> {
                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
                    val pm = activity.getSystemService(Context.POWER_SERVICE) as PowerManager
                    if (!pm.isIgnoringBatteryOptimizations(activity.packageName)) {
                        val intent = Intent(Settings.ACTION_REQUEST_IGNORE_BATTERY_OPTIMIZATIONS)
                        intent.data = android.net.Uri.parse("package:${activity.packageName}")
                        activity.startActivity(intent)
                    }
                }
                result.success(true)
            }
            "processCallEvent" -> {
                val eventJson = call.argument<String>("event")
                if (eventJson != null) {
                    channelRouter.processCallEvent(eventJson)
                    result.success(true)
                } else {
                    result.error("INVALID_ARGS", "event is required", null)
                }
            }
            else -> result.notImplemented()
        }
    }

    private fun getSimCards(): List<Map<String, Any>> {
        val sims = mutableListOf<Map<String, Any>>()
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP_MR1) {
            try {
                val subscriptionManager = activity.getSystemService(Context.TELEPHONY_SUBSCRIPTION_SERVICE) as SubscriptionManager
                val subscriptions = subscriptionManager.activeSubscriptionInfoList
                subscriptions?.forEachIndexed { index, info ->
                    sims.add(mapOf(
                        "slot" to info.simSlotIndex,
                        "subscriptionId" to info.subscriptionId,
                        "carrierName" to (info.carrierName?.toString() ?: "SIM ${index + 1}"),
                        "displayName" to (info.displayName?.toString() ?: "SIM ${index + 1}")
                    ))
                }
            } catch (e: SecurityException) {
                // Permission not granted
            }
        }
        return sims
    }
}
