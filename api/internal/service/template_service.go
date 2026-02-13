package service

import (
	"context"

	"callflow/internal/domain/template"
)

// TemplateService provides template business logic
type TemplateService struct {
	templateRepo template.Repository
}

// NewTemplateService creates a new template service instance
func NewTemplateService(templateRepo template.Repository) *TemplateService {
	return &TemplateService{templateRepo: templateRepo}
}

func (s *TemplateService) Get(ctx context.Context, userID int64) ([]*template.Template, error) {
	return s.templateRepo.GetByUserID(ctx, userID)
}

func (s *TemplateService) GetByID(ctx context.Context, id int64, userID int64) (*template.Template, error) {
	return s.templateRepo.GetByID(ctx, id, userID)
}

func (s *TemplateService) Create(ctx context.Context, userID int64, data template.TemplateCreate) (*template.Template, error) {
	if data.Channel == "" || data.Channel == template.ChannelSMS {
		if len([]rune(data.Body)) > template.SMSMaxChars {
			return nil, template.ErrSMSTooLong
		}
	}
	return s.templateRepo.Create(ctx, userID, data)
}

func (s *TemplateService) Update(ctx context.Context, id int64, userID int64, data template.TemplateUpdate) (*template.Template, error) {
	if data.Channel == "" || data.Channel == template.ChannelSMS {
		if len([]rune(data.Body)) > template.SMSMaxChars {
			return nil, template.ErrSMSTooLong
		}
	}
	return s.templateRepo.Update(ctx, id, userID, data)
}

func (s *TemplateService) Delete(ctx context.Context, id int64, userID int64) error {
	return s.templateRepo.Delete(ctx, id, userID)
}
