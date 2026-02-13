package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents the standard error response structure
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// APIResponse represents the standard API response wrapper
type APIResponse struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

// Success sends a successful response with data
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 successful response with data
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithStatus sends a successful response with custom status code
func SuccessWithStatus(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Error sends a generic error response with code for i18n
func Error(c *gin.Context, status int, code, message, detail string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
			Detail:  detail,
		},
	})
}

// BadRequest sends a 400 error
func BadRequest(c *gin.Context, code, message, detail string) {
	Error(c, http.StatusBadRequest, code, message, detail)
}

// NotFound sends a 404 error
func NotFound(c *gin.Context, code, message, detail string) {
	Error(c, http.StatusNotFound, code, message, detail)
}

// InternalServerError sends a 500 error
func InternalServerError(c *gin.Context, code, message, detail string) {
	Error(c, http.StatusInternalServerError, code, message, detail)
}

// Forbidden sends a 403 error
func Forbidden(c *gin.Context, code, message, detail string) {
	Error(c, http.StatusForbidden, code, message, detail)
}

// Unauthorized sends a 401 error
func Unauthorized(c *gin.Context, code, message, detail string) {
	Error(c, http.StatusUnauthorized, code, message, detail)
}

// Conflict sends a 409 error
func Conflict(c *gin.Context, code, message, detail string) {
	Error(c, http.StatusConflict, code, message, detail)
}
