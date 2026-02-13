package com.callflow.service

import android.content.Context
import android.database.Cursor
import android.provider.CallLog
import android.provider.ContactsContract
import android.util.Log

data class CallInfo(
    val phone: String,
    val contactName: String,
    val duration: Int,
    val type: Int,
    val date: Long
)

class CallLogReader {

    companion object {
        const val TAG = "CallLogReader"
    }

    fun getLatestCall(context: Context): CallInfo? {
        try {
            val cursor: Cursor? = context.contentResolver.query(
                CallLog.Calls.CONTENT_URI,
                arrayOf(
                    CallLog.Calls.NUMBER,
                    CallLog.Calls.CACHED_NAME,
                    CallLog.Calls.DURATION,
                    CallLog.Calls.TYPE,
                    CallLog.Calls.DATE
                ),
                null,
                null,
                "${CallLog.Calls.DATE} DESC"
            )

            if (cursor == null) {
                Log.d(TAG, "Call log cursor is null")
                return null
            }

            cursor.use {
                if (it.moveToFirst()) {
                    val phone = it.getString(0) ?: ""
                    val cachedName = it.getString(1) ?: ""
                    val duration = it.getInt(2)
                    val type = it.getInt(3)
                    val date = it.getLong(4)

                    Log.d(TAG, "Latest call: phone=$phone, type=$type, duration=$duration")

                    val contactName = cachedName.ifEmpty {
                        resolveContactName(context, phone)
                    }

                    return CallInfo(
                        phone = phone,
                        contactName = contactName,
                        duration = duration,
                        type = type,
                        date = date
                    )
                } else {
                    Log.d(TAG, "Call log is empty")
                }
            }
        } catch (e: SecurityException) {
            Log.d(TAG, "Permission denied for call log: ${e.message}")
        } catch (e: Exception) {
            Log.d(TAG, "Error reading call log: ${e.message}")
        }
        return null
    }

    fun resolveContactName(context: Context, phone: String): String {
        if (phone.isEmpty()) return ""
        try {
            val uri = android.net.Uri.withAppendedPath(
                ContactsContract.PhoneLookup.CONTENT_FILTER_URI,
                android.net.Uri.encode(phone)
            )
            val cursor = context.contentResolver.query(
                uri,
                arrayOf(ContactsContract.PhoneLookup.DISPLAY_NAME),
                null, null, null
            )
            cursor?.use {
                if (it.moveToFirst()) {
                    return it.getString(0) ?: ""
                }
            }
        } catch (e: Exception) {
            Log.e(TAG, "Error resolving contact name", e)
        }
        return ""
    }

    fun isContact(context: Context, phone: String): Boolean {
        return resolveContactName(context, phone).isNotEmpty()
    }
}
