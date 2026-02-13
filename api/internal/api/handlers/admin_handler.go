package handler

import (
	"strconv"

	"callflow/internal/api/response"
	"callflow/internal/domain/user"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin HTTP requests
type AdminHandler struct {
	userService user.Service
}

// NewAdminHandler creates a new admin handler instance
func NewAdminHandler(userService user.Service) *AdminHandler {
	return &AdminHandler{userService: userService}
}

// RegisterRoutes registers the admin routes
func (h *AdminHandler) RegisterRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	{
		admin.GET("/users", h.ListUsers)
		admin.PUT("/users/:id/plan", h.UpdatePlan)
		admin.PUT("/users/:id/status", h.UpdateStatus)
	}
}

// ListUsers returns all users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.userService.ListAllUsers(c.Request.Context())
	if err != nil {
		internalError(c, response.ErrListFailed, "Failed to list users", err)
		return
	}
	response.Success(c, users)
}

// UpdatePlan updates a user's plan
func (h *AdminHandler) UpdatePlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidID, "Invalid user ID", err.Error())
		return
	}

	var req struct {
		Plan string `json:"plan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.userService.UpdatePlan(c.Request.Context(), id, req.Plan); err != nil {
		internalError(c, response.ErrUpdateFailed, "Failed to update plan", err)
		return
	}

	response.Success(c, gin.H{"message": "Plan updated successfully"})
}

// UpdateStatus updates a user's status
func (h *AdminHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidID, "Invalid user ID", err.Error())
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.userService.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		internalError(c, response.ErrUpdateFailed, "Failed to update status", err)
		return
	}

	response.Success(c, gin.H{"message": "Status updated successfully"})
}
