package template

import "context"

// Repository defines the interface for template data access
type Repository interface {
	GetByUserID(ctx context.Context, userID int64) ([]*Template, error)
	GetByID(ctx context.Context, id int64, userID int64) (*Template, error)
	Create(ctx context.Context, userID int64, data TemplateCreate) (*Template, error)
	Update(ctx context.Context, id int64, userID int64, data TemplateUpdate) (*Template, error)
	Delete(ctx context.Context, id int64, userID int64) error
}
