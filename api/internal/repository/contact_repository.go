package repository

import (
	"context"

	"callflow/internal/domain/contact"
	db "callflow/internal/sql/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ContactRepository implements contact.Repository
type ContactRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewContactRepository creates a new contact repository
func NewContactRepository(pool *pgxpool.Pool) *ContactRepository {
	return &ContactRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *ContactRepository) GetByUserID(ctx context.Context, userID int64) ([]*contact.Contact, error) {
	rows, err := r.queries.GetContactsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []*contact.Contact{}, nil
	}
	contacts := make([]*contact.Contact, len(rows))
	for i, row := range rows {
		contacts[i] = dbContactToModel(row)
	}
	return contacts, nil
}

func (r *ContactRepository) Upsert(ctx context.Context, userID int64, data contact.ContactUpsert) (*contact.Contact, error) {
	row, err := r.queries.UpsertContact(ctx, db.UpsertContactParams{
		UserID: userID,
		Phone:  data.Phone,
		Name:   pgtype.Text{String: data.Name, Valid: data.Name != ""},
	})
	if err != nil {
		return nil, err
	}
	return dbContactToModel(row), nil
}

func (r *ContactRepository) UpsertBatch(ctx context.Context, userID int64, contacts []contact.ContactUpsert) error {
	for _, c := range contacts {
		err := r.queries.UpsertContactBatch(ctx, db.UpsertContactBatchParams{
			UserID: userID,
			Phone:  c.Phone,
			Name:   pgtype.Text{String: c.Name, Valid: c.Name != ""},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func dbContactToModel(row db.Contact) *contact.Contact {
	var name string
	if row.Name.Valid {
		name = row.Name.String
	}
	return &contact.Contact{
		ID:        row.ID,
		UserID:    row.UserID,
		Phone:     row.Phone,
		Name:      name,
		CreatedAt: row.CreatedAt.Time,
	}
}
