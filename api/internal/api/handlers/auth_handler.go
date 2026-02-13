package handler

import (
	"errors"

	"callflow/internal/api/middleware"
	"callflow/internal/api/response"
	"callflow/internal/domain/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthHandler handles HTTP requests related to authentication
type AuthHandler struct {
	authService auth.Service
	validate    *validator.Validate
}

// NewAuthHandler creates a new authentication handler instance
func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// RegisterRoutes registers the public authentication routes
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", middleware.RateLimitAuth(), h.Register)
		authGroup.POST("/login", middleware.RateLimitAuth(), h.Login)
	}
}

// Register handles user registration with phone and password
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, response.ErrValidationFailed, "Validation failed", err.Error())
		return
	}

	tokenResponse, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrPhoneTaken):
			response.Conflict(c, response.ErrPhoneTaken, "Phone number already registered", "")
		default:
			internalError(c, response.ErrAuthFailed, "Registration failed", err)
		}
		return
	}

	response.Created(c, tokenResponse)
}

// Login handles user login with phone and password
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, response.ErrValidationFailed, "Validation failed", err.Error())
		return
	}

	tokenResponse, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			response.Unauthorized(c, response.ErrInvalidCredentials, "Invalid phone or password", "")
		case errors.Is(err, auth.ErrUserInactive):
			response.Forbidden(c, response.ErrUserInactive, "User account is inactive", "")
		default:
			internalError(c, response.ErrAuthFailed, "Login failed", err)
		}
		return
	}

	response.Success(c, tokenResponse)
}
