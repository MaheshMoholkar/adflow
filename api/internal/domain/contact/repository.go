package contact

import "context"

// Repository defines the interface for contact data access
type Repository interface {
	GetByUserID(ctx context.Context, userID int64) ([]*Contact, error)
	Upsert(ctx context.Context, userID int64, data ContactUpsert) (*Contact, error)
	UpsertBatch(ctx context.Context, userID int64, contacts []ContactUpsert) error
}
