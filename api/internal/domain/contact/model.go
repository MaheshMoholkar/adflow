package contact

import "time"

// Contact represents a recipient of automated messages
type Contact struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ContactUpsert contains data for creating or updating a contact
type ContactUpsert struct {
	Phone string `json:"phone" validate:"required"`
	Name  string `json:"name,omitempty"`
}

// BatchRequest represents a batch contact upsert request
type BatchRequest struct {
	Contacts []ContactUpsert `json:"contacts" validate:"required,min=1"`
}
