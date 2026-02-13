package handler

import (
	"errors"
	"strconv"

	"callflow/internal/api/response"
	"callflow/internal/domain/template"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TemplateHandler handles HTTP requests related to templates
type TemplateHandler struct {
	templateService template.Service
	validate        *validator.Validate
}

// NewTemplateHandler creates a new template handler instance
func NewTemplateHandler(templateService template.Service) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
		validate:        validator.New(),
	}
}

// RegisterRoutes registers the template routes
func (h *TemplateHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tmpl := rg.Group("/template")
	{
		tmpl.GET("", h.Get)
		tmpl.POST("", h.Create)
		tmpl.PUT("/:id", h.Update)
		tmpl.DELETE("/:id", h.Delete)
	}
}

// Get returns the templates for the authenticated user
func (h *TemplateHandler) Get(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	templates, err := h.templateService.Get(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, template.ErrTemplateNotFound) {
			response.Success(c, []interface{}{})
			return
		}
		internalError(c, response.ErrGetFailed, "Failed to get templates", err)
		return
	}

	response.Success(c, templates)
}

// Create creates a new template for the authenticated user
func (h *TemplateHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req template.TemplateCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, response.ErrValidationFailed, "Validation failed", err.Error())
		return
	}

	t, err := h.templateService.Create(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, template.ErrSMSTooLong) {
			response.BadRequest(c, response.ErrSMSTooLong, "SMS body exceeds maximum character limit", "")
			return
		}
		internalError(c, response.ErrCreateFailed, "Failed to create template", err)
		return
	}

	response.SuccessWithStatus(c, 201, t)
}

// Update updates an existing template for the authenticated user
func (h *TemplateHandler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidID, "Invalid template ID", err.Error())
		return
	}

	var req template.TemplateUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, response.ErrValidationFailed, "Validation failed", err.Error())
		return
	}

	t, err := h.templateService.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		if errors.Is(err, template.ErrTemplateNotFound) {
			response.NotFound(c, response.ErrTemplateNotFound, "Template not found", "")
			return
		}
		if errors.Is(err, template.ErrSMSTooLong) {
			response.BadRequest(c, response.ErrSMSTooLong, "SMS body exceeds maximum character limit", "")
			return
		}
		internalError(c, response.ErrUpdateFailed, "Failed to update template", err)
		return
	}

	response.Success(c, t)
}

// Delete deletes a template for the authenticated user
func (h *TemplateHandler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, response.ErrInvalidID, "Invalid template ID", err.Error())
		return
	}

	if err := h.templateService.Delete(c.Request.Context(), id, userID); err != nil {
		internalError(c, response.ErrDeleteFailed, "Failed to delete template", err)
		return
	}

	response.Success(c, gin.H{"message": "Template deleted successfully"})
}
