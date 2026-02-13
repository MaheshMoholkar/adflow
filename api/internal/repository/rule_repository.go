package repository

import (
	"context"
	"encoding/json"
	"errors"

	"callflow/internal/domain/rule"
	db "callflow/internal/sql/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RuleRepository implements rule.Repository
type RuleRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewRuleRepository creates a new rule repository
func NewRuleRepository(pool *pgxpool.Pool) *RuleRepository {
	return &RuleRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *RuleRepository) Get(ctx context.Context, userID int64) (*rule.Rule, error) {
	row, err := r.queries.GetRuleByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rule.ErrRuleNotFound
		}
		return nil, err
	}

	return &rule.Rule{
		ID:        row.ID,
		UserID:    row.UserID,
		Config:    row.Config,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *RuleRepository) Upsert(ctx context.Context, userID int64, config json.RawMessage) (*rule.Rule, error) {
	row, err := r.queries.UpsertRule(ctx, db.UpsertRuleParams{
		UserID: userID,
		Config: config,
	})
	if err != nil {
		return nil, err
	}

	return &rule.Rule{
		ID:        row.ID,
		UserID:    row.UserID,
		Config:    row.Config,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
