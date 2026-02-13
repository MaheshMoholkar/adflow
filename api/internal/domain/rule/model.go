package rule

import (
	"encoding/json"
	"time"
)

// Rule represents user's automation rules
type Rule struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Config    json.RawMessage `json:"config"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// RuleConfig represents the structured rule configuration
type RuleConfig struct {
	Enabled         bool           `json:"enabled"`
	DelaySeconds    int            `json:"delay_seconds"`
	SMS             ChannelConfig  `json:"sms"`
	ExcludedNumbers []string       `json:"excluded_numbers"`
	WorkingHours    *WorkingHours  `json:"working_hours,omitempty"`
	ContactFilter   *ContactFilter `json:"contact_filter,omitempty"`
}

// ChannelConfig represents per-channel rule settings
type ChannelConfig struct {
	Enabled            bool   `json:"enabled"`
	IncomingTemplateID *int64 `json:"incoming_template_id,omitempty"`
	OutgoingTemplateID *int64 `json:"outgoing_template_id,omitempty"`
	MissedTemplateID   *int64 `json:"missed_template_id,omitempty"`
}

// WorkingHours represents working hour constraints
type WorkingHours struct {
	Enabled   bool   `json:"enabled"`
	StartTime string `json:"start_time"` // "09:00"
	EndTime   string `json:"end_time"`   // "18:00"
	Timezone  string `json:"timezone"`
}

// ContactFilter represents contact filtering rules
type ContactFilter struct {
	Mode string `json:"mode"` // all/contacts_only/non_contacts_only
}

// RuleUpdate contains data for updating rules
type RuleUpdate struct {
	Config json.RawMessage `json:"config" validate:"required"`
}
