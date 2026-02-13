package middleware

import (
	"net/http"

	"callflow/internal/domain/user"

	"github.com/gin-gonic/gin"
)

// RequirePlan middleware checks if the user has an active plan
func RequirePlan() gin.HandlerFunc {
	return func(c *gin.Context) {
		planVal, exists := c.Get("plan")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ERR_PLAN_REQUIRED",
					"message": "An active plan is required to use this feature",
				},
			})
			c.Abort()
			return
		}

		plan, ok := planVal.(string)
		if !ok || plan == user.PlanNone {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ERR_PLAN_REQUIRED",
					"message": "An active plan is required to use this feature",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireChannel middleware checks if the user's plan includes the specified channel
func RequireChannel(channel string) gin.HandlerFunc {
	return func(c *gin.Context) {
		planVal, exists := c.Get("plan")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ERR_PLAN_REQUIRED",
					"message": "An active plan is required",
				},
			})
			c.Abort()
			return
		}

		planStr, ok := planVal.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ERR_PLAN_REQUIRED",
					"message": "An active plan is required",
				},
			})
			c.Abort()
			return
		}

		hasChannel := false
		switch planStr {
		case user.PlanSMS:
			hasChannel = channel == "sms"
		}

		if !hasChannel {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ERR_CHANNEL_NOT_IN_PLAN",
					"message": "Your plan does not include the " + channel + " channel",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
