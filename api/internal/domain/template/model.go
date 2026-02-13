package template

import "time"

// Template represents a message template
type Template struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`    // incoming/outgoing/missed
	Channel   string    `json:"channel"` // sms
	Language  string    `json:"language"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TemplateCreate contains data for creating a template
type TemplateCreate struct {
	Name      string `json:"name" validate:"required,max=255"`
	Body      string `json:"body" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=all incoming outgoing missed"`
	Channel   string `json:"channel" validate:"omitempty"`
	Language  string `json:"language"`
	IsDefault bool   `json:"is_default"`
}

// TemplateUpdate contains data for updating a template
type TemplateUpdate struct {
	Name      string `json:"name" validate:"required,max=255"`
	Body      string `json:"body" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=all incoming outgoing missed"`
	Channel   string `json:"channel" validate:"omitempty"`
	Language  string `json:"language"`
	IsDefault bool   `json:"is_default"`
}

// Type constants
const (
	TypeAll      = "all"
	TypeIncoming = "incoming"
	TypeOutgoing = "outgoing"
	TypeMissed   = "missed"
)

// Channel constants
const (
	ChannelSMS = "sms"
)

// SMS limits
const (
	SMSMaxChars = 918 // 6 parts × 153 chars
)
