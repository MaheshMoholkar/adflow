package user

import "context"

// Repository defines the interface for user data access
type Repository interface {
	Create(ctx context.Context, data UserCreate) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	Update(ctx context.Context, id int64, data UserUpdate) (*User, error)
	UpdatePlan(ctx context.Context, id int64, plan string) error
	UpdateStatus(ctx context.Context, id int64, status string) error
	ListAll(ctx context.Context) ([]*User, error)
}
