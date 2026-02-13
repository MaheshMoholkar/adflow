package contact

import "context"

// Service defines the interface for contact business logic
type Service interface {
	Get(ctx context.Context, userID int64) ([]*Contact, error)
	Upsert(ctx context.Context, userID int64, data ContactUpsert) (*Contact, error)
	UpsertBatch(ctx context.Context, userID int64, contacts []ContactUpsert) error
}
