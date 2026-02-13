package repository

import (
	"context"
	"errors"

	"callflow/internal/domain/user"
	db "callflow/internal/sql/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository implements user.Repository
type UserRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *UserRepository) Create(ctx context.Context, data user.UserCreate) (*user.User, error) {
	params := db.CreateUserParams{
		Phone:         data.Phone,
		PhoneVerified: data.PhoneVerified,
		PasswordHash:  data.PasswordHash,
	}
	if data.Name != "" {
		params.Name = pgtype.Text{String: data.Name, Valid: true}
	}
	if data.BusinessName != "" {
		params.BusinessName = pgtype.Text{String: data.BusinessName, Valid: true}
	}
	if data.City != "" {
		params.City = pgtype.Text{String: data.City, Valid: true}
	}
	if data.Address != "" {
		params.Address = pgtype.Text{String: data.Address, Valid: true}
	}
	row, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return dbUserToModel(row), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return dbUserToModel(row), nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	row, err := r.queries.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return dbUserToModel(row), nil
}

func (r *UserRepository) Update(ctx context.Context, id int64, data user.UserUpdate) (*user.User, error) {
	params := db.UpdateUserParams{ID: id}
	if data.Name != nil {
		params.Name = pgtype.Text{String: *data.Name, Valid: true}
	}
	if data.BusinessName != nil {
		params.BusinessName = pgtype.Text{String: *data.BusinessName, Valid: true}
	}
	if data.City != nil {
		params.City = pgtype.Text{String: *data.City, Valid: true}
	}
	if data.Address != nil {
		params.Address = pgtype.Text{String: *data.Address, Valid: true}
	}
	if data.LocationURL != nil {
		params.LocationUrl = pgtype.Text{String: *data.LocationURL, Valid: true}
	}

	row, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return dbUserToModel(row), nil
}

func (r *UserRepository) UpdatePlan(ctx context.Context, id int64, plan string) error {
	return r.queries.UpdateUserPlan(ctx, db.UpdateUserPlanParams{
		ID:   id,
		Plan: plan,
	})
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.queries.UpdateUserStatus(ctx, db.UpdateUserStatusParams{
		ID:     id,
		Status: status,
	})
}

func (r *UserRepository) ListAll(ctx context.Context) ([]*user.User, error) {
	rows, err := r.queries.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*user.User, len(rows))
	for i, row := range rows {
		users[i] = dbUserToModel(row)
	}
	return users, nil
}

func dbUserToModel(row db.User) *user.User {
	u := &user.User{
		ID:            row.ID,
		Phone:         row.Phone,
		PasswordHash:  row.PasswordHash,
		PhoneVerified: row.PhoneVerified,
		Plan:          row.Plan,
		Status:        row.Status,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
	if row.Name.Valid {
		u.Name = row.Name.String
	}
	if row.BusinessName.Valid {
		u.BusinessName = row.BusinessName.String
	}
	if row.City.Valid {
		u.City = row.City.String
	}
	if row.Address.Valid {
		u.Address = row.Address.String
	}
	if row.LocationUrl.Valid {
		u.LocationURL = row.LocationUrl.String
	}
	if row.PlanStartedAt.Valid {
		t := row.PlanStartedAt.Time
		if t.Year() > 0 && t.Year() <= 9999 {
			u.PlanStartedAt = &t
		}
	}
	if row.PlanExpiresAt.Valid {
		t := row.PlanExpiresAt.Time
		if t.Year() > 0 && t.Year() <= 9999 {
			u.PlanExpiresAt = &t
		}
	}
	return u
}
