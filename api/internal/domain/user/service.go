package user

import "context"

// Service defines the interface for user business logic
type Service interface {
	GetUser(ctx context.Context, id int64) (*User, error)
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	CreateUser(ctx context.Context, data UserCreate) (*User, error)
	UpdateUser(ctx context.Context, id int64, data UserUpdate) (*User, error)
	UpdatePlan(ctx context.Context, id int64, plan string) error
	UpdateStatus(ctx context.Context, id int64, status string) error
	ListAllUsers(ctx context.Context) ([]*User, error)
}
