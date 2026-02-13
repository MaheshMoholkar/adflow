package middleware

import (
	"net/http"
	"strings"

	"callflow/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles authentication for protected routes
type AuthMiddleware struct {
	authService auth.Service
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(authService auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// authenticate performs the authentication check and sets user context.
// Returns true if authentication succeeded, false if the request was aborted.
// Does NOT call c.Next().
func (m *AuthMiddleware) authenticate(c *gin.Context) bool {
	var tokenString string
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "ERR_UNAUTHORIZED",
					"message": "Invalid authorization header format",
				},
			})
			c.Abort()
			return false
		}
		tokenString = headerParts[1]
	}

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ERR_UNAUTHORIZED",
				"message": "Access token required",
			},
		})
		c.Abort()
		return false
	}

	// Verify token and extract claims
	claims, err := m.authService.VerifyToken(c.Request.Context(), tokenString)
	if err != nil {
		message := "Invalid token"
		if err == auth.ErrExpiredToken {
			message = "Token has expired"
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ERR_UNAUTHORIZED",
				"message": message,
			},
		})
		c.Abort()
		return false
	}

	// Set user context — simplified for CallFlow (no org context)
	c.Set("userID", claims.UserID)
	c.Set("phone", claims.Phone)
	c.Set("plan", claims.Plan)

	return true
}

// RequireAuth middleware ensures that requests have a valid JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.authenticate(c) {
			c.Next()
		}
	}
}
