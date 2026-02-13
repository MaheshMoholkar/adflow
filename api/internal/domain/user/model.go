package user

import "time"

// User represents a user in the system
type User struct {
	ID            int64      `json:"id"`
	Phone         string     `json:"phone"`
	PasswordHash  string     `json:"-"`
	PhoneVerified bool       `json:"phone_verified"`
	Name          string     `json:"name,omitempty"`
	BusinessName  string     `json:"business_name,omitempty"`
	City          string     `json:"city,omitempty"`
	Address       string     `json:"address,omitempty"`
	LocationURL   string     `json:"location_url,omitempty"`
	Plan          string     `json:"plan"`
	PlanStartedAt *time.Time `json:"plan_started_at,omitempty"`
	PlanExpiresAt *time.Time `json:"plan_expires_at,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// UserCreate contains data for creating a new user
type UserCreate struct {
	Phone         string `json:"phone" validate:"required"`
	PasswordHash  string `json:"-"`
	PhoneVerified bool   `json:"phone_verified"`
	Name          string `json:"name"`
	BusinessName  string `json:"business_name"`
	City          string `json:"city"`
	Address       string `json:"address"`
}

// UserUpdate contains data for updating a user
type UserUpdate struct {
	Name         *string `json:"name,omitempty"`
	BusinessName *string `json:"business_name,omitempty"`
	City         *string `json:"city,omitempty"`
	Address      *string `json:"address,omitempty"`
	LocationURL  *string `json:"location_url,omitempty"`
}

// Plan constants
const (
	PlanNone = "none"
	PlanSMS  = "sms"
)

// Status constants
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// HasChannel checks if user's plan includes the given channel
func (u *User) HasChannel(channel string) bool {
	switch u.Plan {
	case PlanSMS:
		return channel == "sms"
	default:
		return false
	}
}
