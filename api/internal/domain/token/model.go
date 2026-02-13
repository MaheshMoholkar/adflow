package token

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Token represents a token entity
type Token struct {
	ID         int64              `json:"id"`
	UserID     int64              `json:"user_id"`
	Token      string             `json:"token"`
	ExpiresAt  time.Time          `json:"expires_at"`
	IsRevoked  bool               `json:"is_revoked"`
	TokenType  string             `json:"token_type"`
	ClientIP   pgtype.Text        `json:"client_ip,omitempty"`
	UserAgent  pgtype.Text        `json:"user_agent,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	LastUsedAt pgtype.Timestamptz `json:"last_used_at,omitempty"`
}

// TokenCreate is the data structure for creating a new token
type TokenCreate struct {
	UserID    int64       `json:"user_id" validate:"required"`
	Token     string      `json:"token" validate:"required"`
	ExpiresAt time.Time   `json:"expires_at" validate:"required"`
	TokenType string      `json:"token_type" validate:"required,oneof=refresh access"`
	ClientIP  pgtype.Text `json:"client_ip,omitempty"`
	UserAgent pgtype.Text `json:"user_agent,omitempty"`
}
