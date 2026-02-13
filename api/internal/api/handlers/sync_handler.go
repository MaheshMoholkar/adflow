package handler

import (
	"callflow/internal/api/response"
	"callflow/internal/domain/rule"
	"callflow/internal/domain/template"
	"callflow/internal/domain/user"

	"github.com/gin-gonic/gin"
)

// SyncHandler handles HTTP requests related to app configuration sync
type SyncHandler struct {
	userService     user.Service
	templateService template.Service
	ruleService     rule.Service
}

// NewSyncHandler creates a new sync handler instance
func NewSyncHandler(
	userService user.Service,
	templateService template.Service,
	ruleService rule.Service,
) *SyncHandler {
	return &SyncHandler{
		userService:     userService,
		templateService: templateService,
		ruleService:     ruleService,
	}
}

// RegisterRoutes registers the sync routes
func (h *SyncHandler) RegisterRoutes(rg *gin.RouterGroup) {
	sync := rg.Group("/sync")
	{
		sync.GET("/config", h.GetConfig)
	}
}

// GetConfig returns the unified configuration payload for the app
func (h *SyncHandler) GetConfig(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	// Fetch user profile
	u, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		internalError(c, response.ErrGetFailed, "Failed to get user", err)
		return
	}

	// Fetch templates for user
	templates, err := h.templateService.Get(c.Request.Context(), userID)
	if err != nil {
		// Templates might not exist yet, that's ok
		templates = nil
	}

	// Fetch compiled rule config
	ruleConfig, err := h.ruleService.GetCompiledConfig(c.Request.Context(), userID)
	if err != nil {
		// Rules might not exist yet, that's ok
		ruleConfig = nil
	}

	response.Success(c, gin.H{
		"user": gin.H{
			"id":              u.ID,
			"phone":           u.Phone,
			"business_name":   u.BusinessName,
			"plan":            u.Plan,
			"plan_started_at": u.PlanStartedAt,
			"plan_expires_at": u.PlanExpiresAt,
			"status":          u.Status,
		},
		"templates": templates,
		"rules":     ruleConfig,
	})
}
