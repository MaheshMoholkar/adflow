package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

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

const maxTemplateImageBytes = 5 * 1024 * 1024

var allowedTemplateImageContentTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
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
		tmpl.POST("/upload-image", h.UploadImage)
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
		if errors.Is(err, template.ErrInvalidImageURL) || errors.Is(err, template.ErrMissingImageKey) {
			response.BadRequest(c, response.ErrValidationFailed, err.Error(), "")
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
		if errors.Is(err, template.ErrInvalidImageURL) || errors.Is(err, template.ErrMissingImageKey) {
			response.BadRequest(c, response.ErrValidationFailed, err.Error(), "")
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

// UploadImage uploads a template image and returns a public URL and storage key.
func (h *TemplateHandler) UploadImage(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, response.ErrInvalidRequest, "image file is required", err.Error())
		return
	}
	if fileHeader.Size <= 0 {
		response.BadRequest(c, response.ErrInvalidRequest, "image file is empty", "")
		return
	}
	if fileHeader.Size > maxTemplateImageBytes {
		response.BadRequest(c, response.ErrValidationFailed, "image file exceeds 5MB limit", "")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		internalError(c, response.ErrCreateFailed, "Failed to open upload", err)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, maxTemplateImageBytes+1))
	if err != nil {
		internalError(c, response.ErrCreateFailed, "Failed to read upload", err)
		return
	}
	if len(content) > maxTemplateImageBytes {
		response.BadRequest(c, response.ErrValidationFailed, "image file exceeds 5MB limit", "")
		return
	}

	contentType := detectImageContentType(fileHeader.Header.Get("Content-Type"), content)
	if _, ok := allowedTemplateImageContentTypes[contentType]; !ok {
		response.BadRequest(c, response.ErrValidationFailed, "Only JPEG, PNG, and WebP images are allowed", "")
		return
	}

	uploaded, err := h.templateService.UploadImage(
		c.Request.Context(),
		userID,
		fileHeader.Filename,
		contentType,
		content,
	)
	if err != nil {
		if errors.Is(err, template.ErrUploadDisabled) {
			internalError(c, response.ErrCreateFailed, "Image upload is not configured", err)
			return
		}
		internalError(c, response.ErrCreateFailed, "Failed to upload image", err)
		return
	}

	response.Success(c, uploaded)
}

func detectImageContentType(headerValue string, file []byte) string {
	headerType := strings.TrimSpace(strings.Split(headerValue, ";")[0])
	if _, ok := allowedTemplateImageContentTypes[headerType]; ok {
		return headerType
	}
	return http.DetectContentType(file)
}
