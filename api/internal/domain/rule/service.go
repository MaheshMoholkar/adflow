package rule

import "context"

// Service defines the interface for rule business logic
type Service interface {
	Get(ctx context.Context, userID int64) (*Rule, error)
	Upsert(ctx context.Context, userID int64, data RuleUpdate) (*Rule, error)
	GetCompiledConfig(ctx context.Context, userID int64) (*RuleConfig, error)
}
