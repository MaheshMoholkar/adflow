package handler

import (
	"log"

	"callflow/internal/api/response"

	"github.com/gin-gonic/gin"
)

// getUserID extracts and validates the user ID from the gin context.
// Returns the user ID and true if successful, or sends an error response and returns false.
func getUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, response.ErrUnauthorized, "Unauthorized", "")
		return 0, false
	}
	id, ok := val.(int64)
	if !ok {
		response.InternalServerError(c, response.ErrInternalServerError, "Internal server error", "")
		return 0, false
	}
	return id, true
}

// internalError logs the actual error server-side and returns a generic 500 response
// without leaking internal error details to the client.
func internalError(c *gin.Context, code, message string, err error) {
	log.Printf("Internal error [%s]: %v", code, err)
	response.InternalServerError(c, code, message, "")
}
