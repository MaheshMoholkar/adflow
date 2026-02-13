package handler

import (
	"errors"

	"callflow/internal/api/response"
	"callflow/internal/domain/rule"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RuleHandler handles HTTP requests related to rules
type RuleHandler struct {
	ruleService rule.Service
	validate    *validator.Validate
}

// NewRuleHandler creates a new rule handler instance
func NewRuleHandler(ruleService rule.Service) *RuleHandler {
	return &RuleHandler{
		ruleService: ruleService,
		validate:    validator.New(),
	}
}

// RegisterRoutes registers the rule routes
func (h *RuleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rules := rg.Group("/rules")
	{
		rules.GET("", h.Get)
		rules.PUT("", h.Update)
		rules.GET("/config", h.GetCompiledConfig)
	}
}

// Get returns the user's rules
func (h *RuleHandler) Get(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	r, err := h.ruleService.Get(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, rule.ErrRuleNotFound) {
			response.NotFound(c, response.ErrRuleNotFound, "Rules not found", "")
			return
		}
		internalError(c, response.ErrGetFailed, "Failed to get rules", err)
		return
	}

	response.Success(c, r)
}

// Update creates or updates the user's rules
func (h *RuleHandler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req rule.RuleUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, response.ErrValidationFailed, "Validation failed", err.Error())
		return
	}

	r, err := h.ruleService.Upsert(c.Request.Context(), userID, req)
	if err != nil {
		internalError(c, response.ErrUpdateFailed, "Failed to update rules", err)
		return
	}

	response.Success(c, r)
}

// GetCompiledConfig returns the compiled rule config for app sync
func (h *RuleHandler) GetCompiledConfig(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	config, err := h.ruleService.GetCompiledConfig(c.Request.Context(), userID)
	if err != nil {
		internalError(c, response.ErrGetFailed, "Failed to get compiled config", err)
		return
	}

	response.Success(c, config)
}
