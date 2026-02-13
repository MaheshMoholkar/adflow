package handler

import (
	"callflow/internal/api/response"
	"callflow/internal/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService user.Service
	validate    *validator.Validate
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/user")
	{
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
	}
}

// GetProfile returns the authenticated user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	u, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		internalError(c, response.ErrGetFailed, "Failed to get profile", err)
		return
	}

	response.Success(c, u)
}

// UpdateProfile updates the authenticated user's profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req user.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "Invalid request body", err.Error())
		return
	}

	u, err := h.userService.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		internalError(c, response.ErrUpdateFailed, "Failed to update profile", err)
		return
	}

	response.Success(c, u)
}
