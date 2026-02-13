package rule

import (
	"context"
	"encoding/json"
)

// Repository defines the interface for rule data access
type Repository interface {
	Get(ctx context.Context, userID int64) (*Rule, error)
	Upsert(ctx context.Context, userID int64, config json.RawMessage) (*Rule, error)
}
