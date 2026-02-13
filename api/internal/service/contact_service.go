package service

import (
	"context"

	"callflow/internal/domain/contact"
)

// ContactService provides contact business logic
type ContactService struct {
	contactRepo contact.Repository
}

// NewContactService creates a new contact service instance
func NewContactService(contactRepo contact.Repository) *ContactService {
	return &ContactService{contactRepo: contactRepo}
}

func (s *ContactService) Get(ctx context.Context, userID int64) ([]*contact.Contact, error) {
	return s.contactRepo.GetByUserID(ctx, userID)
}

func (s *ContactService) Upsert(ctx context.Context, userID int64, data contact.ContactUpsert) (*contact.Contact, error) {
	return s.contactRepo.Upsert(ctx, userID, data)
}

func (s *ContactService) UpsertBatch(ctx context.Context, userID int64, contacts []contact.ContactUpsert) error {
	return s.contactRepo.UpsertBatch(ctx, userID, contacts)
}
