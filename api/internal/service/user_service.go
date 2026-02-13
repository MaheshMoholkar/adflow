package service

import (
	"context"
	"fmt"

	"callflow/internal/domain/user"
)

// UserService provides user business logic
type UserService struct {
	userRepo user.Repository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo user.Repository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*user.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetUserByPhone(ctx context.Context, phone string) (*user.User, error) {
	return s.userRepo.GetByPhone(ctx, phone)
}

func (s *UserService) CreateUser(ctx context.Context, data user.UserCreate) (*user.User, error) {
	return s.userRepo.Create(ctx, data)
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, data user.UserUpdate) (*user.User, error) {
	return s.userRepo.Update(ctx, id, data)
}

func (s *UserService) UpdatePlan(ctx context.Context, id int64, plan string) error {
	switch plan {
	case user.PlanNone, user.PlanSMS:
		// valid
	default:
		return fmt.Errorf("invalid plan: %s", plan)
	}
	return s.userRepo.UpdatePlan(ctx, id, plan)
}

func (s *UserService) UpdateStatus(ctx context.Context, id int64, status string) error {
	switch status {
	case user.StatusActive, user.StatusInactive:
		// valid
	default:
		return fmt.Errorf("invalid status: %s", status)
	}
	return s.userRepo.UpdateStatus(ctx, id, status)
}

func (s *UserService) ListAllUsers(ctx context.Context) ([]*user.User, error) {
	return s.userRepo.ListAll(ctx)
}
