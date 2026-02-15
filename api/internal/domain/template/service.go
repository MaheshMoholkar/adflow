package template

import "context"

// Service defines the interface for template business logic
type Service interface {
	Get(ctx context.Context, userID int64) ([]*Template, error)
	GetByID(ctx context.Context, id int64, userID int64) (*Template, error)
	Create(ctx context.Context, userID int64, data TemplateCreate) (*Template, error)
	Update(ctx context.Context, id int64, userID int64, data TemplateUpdate) (*Template, error)
	Delete(ctx context.Context, id int64, userID int64) error
	UploadImage(ctx context.Context, userID int64, filename, contentType string, file []byte) (*UploadedImage, error)
}
