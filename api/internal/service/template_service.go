package service

import (
	"context"
	"log"
	"net/url"
	"strings"

	"callflow/internal/domain/template"
)

type TemplateImageStore interface {
	UploadTemplateImage(ctx context.Context, filename, contentType string, file []byte) (*template.UploadedImage, error)
	DeleteTemplateImage(ctx context.Context, imageKey string) error
}

// TemplateService provides template business logic
type TemplateService struct {
	templateRepo template.Repository
	imageStore   TemplateImageStore
}

// NewTemplateService creates a new template service instance
func NewTemplateService(templateRepo template.Repository, imageStore TemplateImageStore) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
		imageStore:   imageStore,
	}
}

func (s *TemplateService) Get(ctx context.Context, userID int64) ([]*template.Template, error) {
	return s.templateRepo.GetByUserID(ctx, userID)
}

func (s *TemplateService) GetByID(ctx context.Context, id int64, userID int64) (*template.Template, error) {
	return s.templateRepo.GetByID(ctx, id, userID)
}

func (s *TemplateService) Create(ctx context.Context, userID int64, data template.TemplateCreate) (*template.Template, error) {
	data.ImageURL = normalizeURL(data.ImageURL)
	data.ImageKey = normalizeStringPtr(data.ImageKey)

	if err := validateImageFields(data.ImageURL, data.ImageKey, true); err != nil {
		return nil, err
	}
	if err := validateSMSLength(data.Channel, data.Body, data.ImageURL); err != nil {
		return nil, err
	}

	return s.templateRepo.Create(ctx, userID, data)
}

func (s *TemplateService) Update(ctx context.Context, id int64, userID int64, data template.TemplateUpdate) (*template.Template, error) {
	existing, err := s.templateRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	data.ImageURL = normalizeURL(data.ImageURL)
	data.ImageKey = normalizeStringPtr(data.ImageKey)

	// Preserve stored image key when URL is unchanged and client does not resend key.
	if data.ImageURL != nil && existing.ImageURL != nil &&
		*data.ImageURL == *existing.ImageURL && data.ImageKey == nil {
		data.ImageKey = existing.ImageKey
	}

	if data.ImageURL == nil {
		data.ImageKey = nil
	}

	requiresImageKey := data.ImageURL != nil && (existing.ImageURL == nil || *data.ImageURL != *existing.ImageURL)
	if err := validateImageFields(data.ImageURL, data.ImageKey, requiresImageKey); err != nil {
		return nil, err
	}
	if err := validateSMSLength(data.Channel, data.Body, data.ImageURL); err != nil {
		return nil, err
	}

	updated, err := s.templateRepo.Update(ctx, id, userID, data)
	if err != nil {
		return nil, err
	}

	if shouldDeleteOldImage(existing.ImageKey, updated.ImageKey) {
		s.deleteImageKeyAsync(existing.ImageKey)
	}

	return updated, nil
}

func (s *TemplateService) Delete(ctx context.Context, id int64, userID int64) error {
	existing, err := s.templateRepo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	if err := s.templateRepo.Delete(ctx, id, userID); err != nil {
		return err
	}

	if existing.ImageKey != nil {
		s.deleteImageKeyAsync(existing.ImageKey)
	}

	return nil
}

func (s *TemplateService) UploadImage(ctx context.Context, _ int64, filename, contentType string, file []byte) (*template.UploadedImage, error) {
	if s.imageStore == nil {
		return nil, template.ErrUploadDisabled
	}
	return s.imageStore.UploadTemplateImage(ctx, filename, contentType, file)
}

func validateSMSLength(channel string, body string, imageURL *string) error {
	if channel != "" && channel != template.ChannelSMS {
		return nil
	}

	finalText := body
	if imageURL != nil {
		finalText = *imageURL + "\n" + finalText
	}
	if len([]rune(finalText)) > template.SMSMaxChars {
		return template.ErrSMSTooLong
	}
	return nil
}

func validateImageFields(imageURL, imageKey *string, requireImageKey bool) error {
	if imageURL == nil {
		return nil
	}
	parsed, err := url.Parse(*imageURL)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return template.ErrInvalidImageURL
	}
	if requireImageKey && imageKey == nil {
		return template.ErrMissingImageKey
	}
	return nil
}

func normalizeURL(v *string) *string {
	value := normalizeStringPtr(v)
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeStringPtr(v *string) *string {
	if v == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*v)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func shouldDeleteOldImage(oldKey, newKey *string) bool {
	if oldKey == nil || *oldKey == "" {
		return false
	}
	if newKey == nil || *newKey == "" {
		return true
	}
	return *oldKey != *newKey
}

func (s *TemplateService) deleteImageKeyAsync(imageKey *string) {
	if imageKey == nil || s.imageStore == nil {
		return
	}
	if err := s.imageStore.DeleteTemplateImage(context.Background(), *imageKey); err != nil {
		log.Printf("failed to delete image key %s: %v", *imageKey, err)
	}
}
