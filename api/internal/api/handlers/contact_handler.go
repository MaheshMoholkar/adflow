package handler

import (
	"callflow/internal/api/response"
	"callflow/internal/domain/contact"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ContactHandler handles HTTP requests related to contacts
type ContactHandler struct {
	contactService contact.Service
	validate       *validator.Validate
}

// NewContactHandler creates a new contact handler instance
func NewContactHandler(contactService contact.Service) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
		validate:       validator.New(),
	}
}

// RegisterRoutes registers the contact routes
func (h *ContactHandler) RegisterRoutes(rg *gin.RouterGroup) {
	contacts := rg.Group("/contacts")
	{
		contacts.GET("", h.Get)
		contacts.POST("/batch", h.BatchUpsert)
	}
}

// Get returns the contacts for the authenticated user
func (h *ContactHandler) Get(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	contacts, err := h.contactService.Get(c.Request.Context(), userID)
	if err != nil {
		internalError(c, response.ErrGetFailed, "Failed to get contacts", err)
		return
	}

	response.Success(c, contacts)
}

// BatchUpsert creates or updates contacts in batch
func (h *ContactHandler) BatchUpsert(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req contact.BatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, response.ErrValidationFailed, "Validation failed", err.Error())
		return
	}

	if err := h.contactService.UpsertBatch(c.Request.Context(), userID, req.Contacts); err != nil {
		internalError(c, response.ErrCreateFailed, "Failed to save contacts", err)
		return
	}

	response.Success(c, gin.H{
		"message": "Contacts saved successfully",
		"count":   len(req.Contacts),
	})
}
